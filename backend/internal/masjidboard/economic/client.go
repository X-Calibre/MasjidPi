package economic

import (
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultAPIURL  = "https://www.jamiatsa.org/api/economic/latest"
	SourcePageURL  = "https://www.jamiatsa.org/tools/economic-indicators.html"
	defaultTimeout = 20 * time.Second
)

type Client struct {
	HTTPClient *http.Client
	APIURL     string
	Now        func() time.Time
}

type apiResponse struct {
	GregorianDate  string `json:"gregorian_date"`
	HijriDay       int    `json:"hijri_day"`
	HijriMonth     int    `json:"hijri_month"`
	HijriYear      int    `json:"hijri_year"`
	HijriMonthName string `json:"hijri_month_name"`
	USDZAR         string `json:"usd_zar"`
	Gold24K        string `json:"gold_24k"`
	Gold22K        string `json:"gold_22k"`
	Gold18K        string `json:"gold_18k"`
	Gold14K        string `json:"gold_14k"`
	Gold9K         string `json:"gold_9k"`
	Silver         string `json:"silver"`
	Nisaab         string `json:"nisaab"`
	MinimumMahr    string `json:"mahr_min"`
	MahrFaatimi    string `json:"mahr_faatimi"`
	Krugerrand     string `json:"krugerrand"`
}

func (c Client) Fetch(ctx context.Context) (Indicators, error) {
	endpoint := strings.TrimSpace(c.APIURL)
	if endpoint == "" {
		endpoint = DefaultAPIURL
	}
	if _, err := url.ParseRequestURI(endpoint); err != nil {
		return Indicators{}, fmt.Errorf("economic indicators: invalid API URL: %w", err)
	}

	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: defaultTimeout}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return Indicators{}, fmt.Errorf("economic indicators: create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "MasjidPi Islamic Economic Indicators")

	response, err := client.Do(req)
	if err != nil {
		return Indicators{}, fmt.Errorf("economic indicators: fetch: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return Indicators{}, fmt.Errorf("economic indicators: unexpected HTTP status %s", response.Status)
	}

	contentType := response.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil || mediaType != "application/json" {
		return Indicators{}, fmt.Errorf("economic indicators: unexpected content type %q", contentType)
	}

	var data apiResponse
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&data); err != nil {
		return Indicators{}, fmt.Errorf("economic indicators: decode response: %w", err)
	}

	return normalizeResponse(data, c.now().UTC())
}

func (c Client) now() time.Time {
	if c.Now != nil {
		return c.Now()
	}
	return time.Now()
}

func normalizeResponse(data apiResponse, fetchedAt time.Time) (Indicators, error) {
	effectiveDate := strings.TrimSpace(data.GregorianDate)
	parsedDate, err := time.Parse("2006-01-02", effectiveDate)
	if err != nil || parsedDate.Format("2006-01-02") != effectiveDate {
		return Indicators{}, fmt.Errorf("economic indicators: invalid gregorian_date %q", data.GregorianDate)
	}
	if data.HijriDay <= 0 || data.HijriMonth <= 0 || data.HijriMonth > 12 || data.HijriYear <= 0 || strings.TrimSpace(data.HijriMonthName) == "" {
		return Indicators{}, fmt.Errorf("economic indicators: invalid Hijri date")
	}

	result := Indicators{
		Source:        SourceName,
		SourceURL:     SourcePageURL,
		EffectiveDate: effectiveDate,
		HijriDate:     fmt.Sprintf("%d %s %d", data.HijriDay, strings.TrimSpace(data.HijriMonthName), data.HijriYear),
		FetchedAt:     fetchedAt,
	}

	fields := []struct {
		name   string
		raw    string
		target *float64
	}{
		{"usd_zar", data.USDZAR, &result.RandDollar},
		{"gold_24k", data.Gold24K, &result.Gold24Carat},
		{"gold_22k", data.Gold22K, &result.Gold22Carat},
		{"gold_18k", data.Gold18K, &result.Gold18Carat},
		{"gold_14k", data.Gold14K, &result.Gold14Carat},
		{"gold_9k", data.Gold9K, &result.Gold9Carat},
		{"silver", data.Silver, &result.Silver},
		{"nisaab", data.Nisaab, &result.Nisaab},
		{"mahr_min", data.MinimumMahr, &result.MinimumMahr},
		{"mahr_faatimi", data.MahrFaatimi, &result.MahrFaatimi},
		{"krugerrand", data.Krugerrand, &result.Krugerrand},
	}
	for _, field := range fields {
		value, parseErr := parsePositiveNumber(field.name, field.raw)
		if parseErr != nil {
			return Indicators{}, parseErr
		}
		*field.target = value
	}

	return result, nil
}

func parsePositiveNumber(name, raw string) (float64, error) {
	trimmed := strings.TrimSpace(raw)
	value, err := strconv.ParseFloat(trimmed, 64)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("economic indicators: invalid %s value %q", name, raw)
	}
	return value, nil
}
