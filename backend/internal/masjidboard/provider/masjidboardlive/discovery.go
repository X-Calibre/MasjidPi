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

const defaultDiscoveryEndpoint = "https://masjidboardlive.com/findmasjid/assets/php/live_page.php"

// DiscoveryClient retrieves structured board-directory results from the public
// MasjidBoard Live FindMasjid endpoint. Discovery is intentionally separate
// from individual-board retrieval.
type DiscoveryClient struct {
	HTTPClient *http.Client
	Endpoint   string
}

// MasjidSearch identifies one FindMasjid city-level search request.
type MasjidSearch struct {
	Search   string
	Country  string
	Province string
	City     string
}

// CatalogueEntry is the provider-level representation of a board returned by
// FindMasjid. It contains selection/discovery metadata, not the complete board
// timetable model.
type CatalogueEntry struct {
	Name             string
	City             string
	WebURL           string
	MBLID            string
	TimeZoneOffsetMS int64
	LastUpdated      string
	FajrJamaah       string
	DhuhrJamaah      string
	AsrJamaah        string
	MaghribAdhan     string
	EshaJamaah       string
	Sunset           string
	JumuahKhutbah    string
	RamadhaanActive  string
	DateAdjust       string
	MoonSeen         string
	LadiesFacility   string
}

// DiscoveryResult preserves the location label supplied by FindMasjid plus
// the normalised entries returned for the request.
type DiscoveryResult struct {
	Location string
	Entries  []CatalogueEntry
}

// SearchMasjids queries the structured FindMasjid endpoint for board records.
func (c DiscoveryClient) SearchMasjids(ctx context.Context, search MasjidSearch) (DiscoveryResult, error) {
	endpoint := c.Endpoint
	if endpoint == "" {
		endpoint = defaultDiscoveryEndpoint
	}

	form := url.Values{}
	form.Set("type", "masjid")
	form.Set("search", strings.TrimSpace(search.Search))
	form.Set("countryName", strings.TrimSpace(search.Country))
	form.Set("provinceName", strings.TrimSpace(search.Province))
	form.Set("cityName", strings.TrimSpace(search.City))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return DiscoveryResult{}, fmt.Errorf("masjidboardlive: create discovery request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}

	resp, err := client.Do(req)
	if err != nil {
		return DiscoveryResult{}, fmt.Errorf("masjidboardlive: discovery request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return DiscoveryResult{}, fmt.Errorf("masjidboardlive: unexpected discovery HTTP status %s", resp.Status)
	}

	var rows []json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&rows); err != nil {
		return DiscoveryResult{}, fmt.Errorf("masjidboardlive: decode discovery response: %w", err)
	}
	return parseMasjidDiscovery(rows)
}

func parseMasjidDiscovery(rows []json.RawMessage) (DiscoveryResult, error) {
	if len(rows) == 0 {
		return DiscoveryResult{}, fmt.Errorf("masjidboardlive: empty discovery response")
	}

	var location string
	if err := json.Unmarshal(rows[0], &location); err != nil {
		return DiscoveryResult{}, fmt.Errorf("masjidboardlive: discovery location: %w", err)
	}
	location = strings.TrimSpace(location)

	entries := make([]CatalogueEntry, 0, len(rows)-1)
	for i, raw := range rows[1:] {
		var record struct {
			Masjid          string `json:"masjid"`
			FajrJamaat      string `json:"fajr_jamaat"`
			ZuhrJamaat      string `json:"zuhr_jamaat"`
			AsrJamaat       string `json:"asr_jamaat"`
			MaghribAdhan    string `json:"maghrib_adhan"`
			EshaJamaat      string `json:"esha_jamaat"`
			LastUpdated     string `json:"last_updated"`
			MBLID           string `json:"MBL_ID"`
			City            string `json:"city"`
			Sunset          string `json:"sunset"`
			TimeZoneMilli   string `json:"time_zone_milli"`
			WebURL          string `json:"web_url"`
			JumuahKhutbah   string `json:"jumuah_khutbah"`
			RamadhaanActive string `json:"ramadhaanactive"`
			DateAdjust      string `json:"date_adjust"`
			MoonSeen        string `json:"moon_seen"`
			LadiesFacility  string `json:"ladies_facility"`
		}
		if err := json.Unmarshal(raw, &record); err != nil {
			return DiscoveryResult{}, fmt.Errorf("masjidboardlive: decode discovery record %d: %w", i+1, err)
		}

		name := strings.TrimSpace(record.Masjid)
		webURL := strings.TrimSpace(record.WebURL)
		if name == "" || webURL == "" {
			return DiscoveryResult{}, fmt.Errorf("masjidboardlive: discovery record %d missing masjid or web_url", i+1)
		}

		var offsetMS int64
		if value := strings.TrimSpace(record.TimeZoneMilli); value != "" {
			parsed, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return DiscoveryResult{}, fmt.Errorf("masjidboardlive: discovery record %d invalid time_zone_milli %q: %w", i+1, value, err)
			}
			offsetMS = parsed
		}

		entries = append(entries, CatalogueEntry{
			Name:             name,
			City:             strings.TrimSpace(record.City),
			WebURL:           webURL,
			MBLID:            strings.TrimSpace(record.MBLID),
			TimeZoneOffsetMS: offsetMS,
			LastUpdated:      strings.TrimSpace(record.LastUpdated),
			FajrJamaah:       strings.TrimSpace(record.FajrJamaat),
			DhuhrJamaah:      strings.TrimSpace(record.ZuhrJamaat),
			AsrJamaah:        strings.TrimSpace(record.AsrJamaat),
			MaghribAdhan:     strings.TrimSpace(record.MaghribAdhan),
			EshaJamaah:       strings.TrimSpace(record.EshaJamaat),
			Sunset:           strings.TrimSpace(record.Sunset),
			JumuahKhutbah:    strings.TrimSpace(record.JumuahKhutbah),
			RamadhaanActive:  strings.TrimSpace(record.RamadhaanActive),
			DateAdjust:       strings.TrimSpace(record.DateAdjust),
			MoonSeen:         strings.TrimSpace(record.MoonSeen),
			LadiesFacility:   strings.TrimSpace(record.LadiesFacility),
		})
	}

	return DiscoveryResult{Location: location, Entries: entries}, nil
}
