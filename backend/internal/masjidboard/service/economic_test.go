package service

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/economic"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

type recordingLogger struct {
	warnings atomic.Int32
}

func (l *recordingLogger) Warn(string, ...any) {
	l.warnings.Add(1)
}

func economicResponse(effectiveDate string) string {
	return fmt.Sprintf(`{
		"gregorian_date":%q,
		"hijri_day":11,
		"hijri_month":3,
		"hijri_year":1448,
		"hijri_month_name":"Rabi al-Awwal",
		"usd_zar":"16.01",
		"gold_24k":"2385.85",
		"gold_22k":"2187.03",
		"gold_21k":"2087.62",
		"gold_18k":"1789.39",
		"gold_14k":"1391.75",
		"gold_9k":"894.69",
		"silver":"35.45",
		"nisaab":"21708.16",
		"mahr_min":"1085.40",
		"mahr_faatimi":"54270.41",
		"krugerrand":"77626.36",
		"updated_at":"2026-08-25T07:30:00Z",
		"notes":null
	}`, effectiveDate)
}

func completeIndicators(effectiveDate string) *economic.Indicators {
	return &economic.Indicators{
		EffectiveDate: effectiveDate,
		RandDollar:    1, Gold24Carat: 1, Gold22Carat: 1, Gold21Carat: 1, Gold18Carat: 1,
		Gold14Carat: 1, Gold9Carat: 1, Silver: 1, Nisaab: 1,
		MinimumMahr: 1, MahrFaatimi: 1, Krugerrand: 1,
		UpdatedAt: time.Date(2026, 8, 24, 7, 30, 0, 0, time.UTC),
	}
}

func TestRefreshEconomicIndicatorsFetchesOnceForCurrentSourceDay(t *testing.T) {
	t.Parallel()
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, economicResponse("2026-08-24"))
	}))
	defer server.Close()

	cachePath := filepath.Join(t.TempDir(), "indicators.json")
	now := time.Date(2026, 8, 24, 19, 0, 0, 0, time.UTC)
	service := &Service{
		selection:      selection.State{ShowEconomicIndicators: true},
		economicClient: economic.Client{APIURL: server.URL, HTTPClient: server.Client(), Now: func() time.Time { return now }},
		economicStore:  economic.Store{Path: cachePath},
	}
	if err := service.RefreshEconomicIndicators(context.Background()); err != nil {
		t.Fatalf("first refresh error = %v", err)
	}
	if err := service.RefreshEconomicIndicators(context.Background()); err != nil {
		t.Fatalf("second refresh error = %v", err)
	}
	if got := requests.Load(); got != 1 {
		t.Fatalf("requests = %d, want 1", got)
	}
	if got := service.EconomicIndicators(); got == nil || got.Nisaab != 21708.16 || got.Krugerrand != 77626.36 || got.Gold21Carat != 2087.62 {
		t.Fatalf("EconomicIndicators() = %+v", got)
	}
	if cached, err := (economic.Store{Path: cachePath}).Load(); err != nil || cached == nil {
		t.Fatalf("cached indicators = %+v, %v", cached, err)
	}
}

func TestRefreshEconomicIndicatorsBackfillsIncompleteCurrentDayCache(t *testing.T) {
	t.Parallel()
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, economicResponse("2026-08-24"))
	}))
	defer server.Close()

	cachePath := filepath.Join(t.TempDir(), "indicators.json")
	now := time.Date(2026, 8, 24, 6, 0, 0, 0, time.UTC)
	service := &Service{
		selection: selection.State{ShowEconomicIndicators: true},
		indicators: &economic.Indicators{
			EffectiveDate: "2026-08-24", RandDollar: 16.01,
			Nisaab: 21708.16, Krugerrand: 77626.36,
		},
		economicClient: economic.Client{APIURL: server.URL, HTTPClient: server.Client(), Now: func() time.Time { return now }},
		economicStore:  economic.Store{Path: cachePath},
	}

	if err := service.RefreshEconomicIndicators(context.Background()); err != nil {
		t.Fatalf("refresh error = %v", err)
	}
	if got := requests.Load(); got != 1 {
		t.Fatalf("requests = %d, want 1", got)
	}
	if got := service.EconomicIndicators(); got == nil || got.Gold14Carat != 1391.75 || got.Gold9Carat != 894.69 || got.Gold21Carat != 2087.62 {
		t.Fatalf("EconomicIndicators() = %+v", got)
	}
}

func TestRefreshEconomicIndicatorsWaitsUntilNineInJohannesburg(t *testing.T) {
	t.Parallel()
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		fmt.Fprint(w, economicResponse("2026-08-25"))
	}))
	defer server.Close()

	now := time.Date(2026, 8, 25, 6, 59, 0, 0, time.UTC)
	current := completeIndicators("2026-08-24")
	current.FetchedAt = time.Date(2026, 8, 24, 19, 0, 0, 0, time.UTC)
	service := &Service{
		selection:      selection.State{ShowEconomicIndicators: true},
		indicators:     current,
		economicClient: economic.Client{APIURL: server.URL, HTTPClient: server.Client(), Now: func() time.Time { return now }},
	}
	if err := service.RefreshEconomicIndicators(context.Background()); err != nil {
		t.Fatalf("refresh error = %v", err)
	}
	if got := requests.Load(); got != 0 {
		t.Fatalf("requests = %d, want 0", got)
	}
}

func TestRefreshEconomicIndicatorsRetriesUnchangedSource(t *testing.T) {
	t.Parallel()
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, economicResponse("2026-08-24"))
	}))
	defer server.Close()

	originalFetchedAt := time.Date(2026, 8, 24, 19, 0, 0, 0, time.UTC)
	now := time.Date(2026, 8, 25, 7, 0, 0, 0, time.UTC)
	current := completeIndicators("2026-08-24")
	current.FetchedAt = originalFetchedAt
	service := &Service{
		selection:      selection.State{ShowEconomicIndicators: true},
		indicators:     current,
		economicClient: economic.Client{APIURL: server.URL, HTTPClient: server.Client(), Now: func() time.Time { return now }},
		economicStore:  economic.Store{Path: filepath.Join(t.TempDir(), "indicators.json")},
	}
	for range 2 {
		if err := service.RefreshEconomicIndicators(context.Background()); err != nil {
			t.Fatalf("refresh error = %v", err)
		}
	}
	if got := requests.Load(); got != 2 {
		t.Fatalf("requests = %d, want 2", got)
	}
	if got := service.EconomicIndicators().FetchedAt; !got.Equal(originalFetchedAt) {
		t.Fatalf("FetchedAt = %v, want %v", got, originalFetchedAt)
	}
}

func TestRefreshEconomicIndicatorsStopsAfterEffectiveDateAdvances(t *testing.T) {
	t.Parallel()
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, economicResponse("2026-08-25"))
	}))
	defer server.Close()

	now := time.Date(2026, 8, 25, 7, 30, 0, 0, time.UTC)
	current := completeIndicators("2026-08-24")
	service := &Service{
		selection:      selection.State{ShowEconomicIndicators: true},
		indicators:     current,
		economicClient: economic.Client{APIURL: server.URL, HTTPClient: server.Client(), Now: func() time.Time { return now }},
		economicStore:  economic.Store{Path: filepath.Join(t.TempDir(), "indicators.json")},
	}
	if err := service.RefreshEconomicIndicators(context.Background()); err != nil {
		t.Fatalf("first refresh error = %v", err)
	}
	if err := service.RefreshEconomicIndicators(context.Background()); err != nil {
		t.Fatalf("second refresh error = %v", err)
	}
	if got := requests.Load(); got != 1 {
		t.Fatalf("requests = %d, want 1", got)
	}
	if got := service.EconomicIndicators().EffectiveDate; got != "2026-08-25" {
		t.Fatalf("EffectiveDate = %q", got)
	}
}

func TestEconomicRefreshDueLimitsWeekendAttempts(t *testing.T) {
	t.Parallel()
	current := completeIndicators("2026-08-28")
	tests := []struct {
		name string
		now  time.Time
		want bool
	}{
		{"before weekend window", time.Date(2026, 8, 29, 6, 59, 0, 0, time.UTC), false},
		{"weekend first attempt", time.Date(2026, 8, 29, 7, 0, 0, 0, time.UTC), true},
		{"weekend retry window", time.Date(2026, 8, 29, 8, 29, 0, 0, time.UTC), true},
		{"weekend cutoff", time.Date(2026, 8, 29, 8, 30, 0, 0, time.UTC), false},
		{"weekday retries continue", time.Date(2026, 8, 31, 18, 0, 0, 0, time.UTC), true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := economicRefreshDue(current, test.now); got != test.want {
				t.Fatalf("economicRefreshDue() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestEconomicIndicatorsHiddenWhenDisabled(t *testing.T) {
	t.Parallel()
	service := &Service{indicators: &economic.Indicators{Source: economic.SourceName, SourceURL: "https://example.test", EffectiveDate: "2026-08-24", Nisaab: 1, Krugerrand: 2}}
	if got := service.EconomicIndicators(); got != nil {
		t.Fatalf("EconomicIndicators() = %+v, want nil", got)
	}
}

func TestRefreshEconomicIndicatorsLogsFailureAndKeepsCurrentData(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	current := completeIndicators("2026-08-31")
	log := &recordingLogger{}
	now := time.Date(2026, 9, 2, 8, 0, 0, 0, time.UTC)
	service := &Service{
		selection:      selection.State{ShowEconomicIndicators: true},
		indicators:     current,
		economicClient: economic.Client{APIURL: server.URL, HTTPClient: server.Client(), Now: func() time.Time { return now }},
		log:            log,
	}

	if err := service.RefreshEconomicIndicators(context.Background()); err == nil {
		t.Fatal("RefreshEconomicIndicators() expected error")
	}
	if got := log.warnings.Load(); got != 1 {
		t.Fatalf("warning count = %d, want 1", got)
	}
	if got := service.EconomicIndicators(); got == nil || got.EffectiveDate != "2026-08-31" {
		t.Fatalf("EconomicIndicators() = %+v, want last-known-good data", got)
	}
}
