package masjidboardlive

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// HierarchyEntry is one country, province/region, or town/city returned by the
// FindMasjid hierarchy endpoint.
type HierarchyEntry struct {
	Name  string
	Count int
}

// Countries returns the active countries exposed by FindMasjid.
func (c DiscoveryClient) Countries(ctx context.Context) ([]HierarchyEntry, error) {
	return c.fetchHierarchyPairs(ctx, "country", "", "", "", "")
}

// Regions returns the province/region buckets beneath a country. Blank region
// names are intentionally preserved because upstream has been observed to
// expose a legitimate blank bucket.
func (c DiscoveryClient) Regions(ctx context.Context, country string) ([]HierarchyEntry, error) {
	country = strings.TrimSpace(country)
	if country == "" {
		return nil, fmt.Errorf("masjidboardlive: country is required")
	}
	return c.fetchHierarchyPairs(ctx, "province", country, country, "", "")
}

// Cities returns the town/city buckets beneath a country and region. The
// FindMasjid city response contains a primary city row array plus auxiliary
// grouping metadata used only by its frontend; only the primary rows are
// returned here. The primary rows have been observed in both pair-array and
// object forms, so both upstream encodings are accepted.
func (c DiscoveryClient) Cities(ctx context.Context, country, region string) ([]HierarchyEntry, error) {
	country = strings.TrimSpace(country)
	region = strings.TrimSpace(region)
	if country == "" {
		return nil, fmt.Errorf("masjidboardlive: country is required")
	}

	body, err := c.postDiscovery(ctx, "cityProvince", region, country, region, "")
	if err != nil {
		return nil, err
	}

	var top []json.RawMessage
	if err := json.Unmarshal(body, &top); err != nil {
		return nil, fmt.Errorf("masjidboardlive: decode city hierarchy response: %w", err)
	}
	if len(top) == 0 {
		return nil, fmt.Errorf("masjidboardlive: empty city hierarchy response")
	}

	entries, err := parseCityHierarchyRows(top[0])
	if err != nil {
		return nil, fmt.Errorf("masjidboardlive: decode city hierarchy rows: %w", err)
	}
	return entries, nil
}

func (c DiscoveryClient) fetchHierarchyPairs(ctx context.Context, kind, search, country, region, city string) ([]HierarchyEntry, error) {
	body, err := c.postDiscovery(ctx, kind, search, country, region, city)
	if err != nil {
		return nil, err
	}
	return parseHierarchyPairs(body, kind == "province")
}

func (c DiscoveryClient) postDiscovery(ctx context.Context, kind, search, country, region, city string) ([]byte, error) {
	endpoint := c.Endpoint
	if endpoint == "" {
		endpoint = defaultDiscoveryEndpoint
	}

	form := url.Values{}
	form.Set("type", kind)
	form.Set("search", search)
	form.Set("countryName", country)
	form.Set("provinceName", region)
	form.Set("cityName", city)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("masjidboardlive: create hierarchy request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("masjidboardlive: hierarchy request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("masjidboardlive: unexpected hierarchy HTTP status %s", resp.Status)
	}

	var raw json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("masjidboardlive: decode hierarchy response: %w", err)
	}
	if string(raw) == "null" {
		return nil, fmt.Errorf("masjidboardlive: hierarchy response is null for %s", kind)
	}
	return raw, nil
}

func parseHierarchyPairs(raw []byte, allowBlankName bool) ([]HierarchyEntry, error) {
	var rows [][]json.RawMessage
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("empty hierarchy response")
	}

	entries := make([]HierarchyEntry, 0, len(rows))
	for i, row := range rows {
		if len(row) < 2 {
			return nil, fmt.Errorf("row %d has %d fields, want at least 2", i, len(row))
		}
		var name, countText string
		if err := json.Unmarshal(row[0], &name); err != nil {
			return nil, fmt.Errorf("row %d name: %w", i, err)
		}
		if err := json.Unmarshal(row[1], &countText); err != nil {
			return nil, fmt.Errorf("row %d count: %w", i, err)
		}
		entry, err := hierarchyEntry(name, countText, allowBlankName, i)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func parseCityHierarchyRows(raw []byte) ([]HierarchyEntry, error) {
	var rows []json.RawMessage
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("empty hierarchy response")
	}

	entries := make([]HierarchyEntry, 0, len(rows))
	for i, rawRow := range rows {
		var pair []json.RawMessage
		if err := json.Unmarshal(rawRow, &pair); err == nil && len(pair) >= 2 {
			var name, countText string
			if err := json.Unmarshal(pair[0], &name); err != nil {
				return nil, fmt.Errorf("row %d name: %w", i, err)
			}
			if err := json.Unmarshal(pair[1], &countText); err != nil {
				return nil, fmt.Errorf("row %d count: %w", i, err)
			}
			entry, err := hierarchyEntry(name, countText, false, i)
			if err != nil {
				return nil, err
			}
			entries = append(entries, entry)
			continue
		}

		var object map[string]json.RawMessage
		if err := json.Unmarshal(rawRow, &object); err != nil {
			return nil, fmt.Errorf("row %d is neither pair nor object: %w", i, err)
		}
		nameRaw, ok := object["city"]
		if !ok {
			nameRaw = object["0"]
		}
		countRaw, ok := object["COUNT(*)"]
		if !ok {
			countRaw = object["1"]
		}
		if len(nameRaw) == 0 || len(countRaw) == 0 {
			return nil, fmt.Errorf("row %d missing city or count", i)
		}
		var name, countText string
		if err := json.Unmarshal(nameRaw, &name); err != nil {
			return nil, fmt.Errorf("row %d city: %w", i, err)
		}
		if err := json.Unmarshal(countRaw, &countText); err != nil {
			return nil, fmt.Errorf("row %d count: %w", i, err)
		}
		entry, err := hierarchyEntry(name, countText, false, i)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func hierarchyEntry(name, countText string, allowBlankName bool, row int) (HierarchyEntry, error) {
	name = strings.TrimSpace(name)
	if name == "" && !allowBlankName {
		return HierarchyEntry{}, fmt.Errorf("row %d name is blank", row)
	}
	count, err := strconv.Atoi(strings.TrimSpace(countText))
	if err != nil || count < 0 {
		return HierarchyEntry{}, fmt.Errorf("row %d invalid count %q", row, countText)
	}
	return HierarchyEntry{Name: name, Count: count}, nil
}
