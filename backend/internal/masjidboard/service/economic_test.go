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

func economicResponse(modifiedGMT, effectiveDate string) string {
	return fmt.Sprintf(`[{"date":"2026-08-14T09:02:11","modified_gmt":%q,"link":"https://www.jamiatsa.org/source/","title":{"rendered":"Rabi al-Awwal 1448"},"content":{"rendered":"<table><thead><tr><th>Hijri</th><th>Date</th><th>Rand-Dollar</th><th>24 Carat</th><th>22 Carat</th><th>18 Carat</th><th>Silver</th><th>Nisaab</th><th>Min Mahr</th><th>Mahr Faatimi</th><th>Krugerrand</th></tr></thead><tbody><tr><td>11</td><td>%s</td><td>R16.01</td><td>R2385.85</td><td>R2187.03</td><td>R1789.39</td><td>R35.45</td><td>R21708.16</td><td>R1085.40</td><td>R54270.41</td><td>R77626.36</td></tr></tbody></table>"}}]`, modifiedGMT, effectiveDate)
}

func TestRefreshEconomicIndicatorsFetchesOnceForCurrentSourceDay(t *testing.T) {
	t.Parallel()
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, economicResponse("2026-08-24T07:02:11", "24 Aug"))
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
	if got := service.EconomicIndicators(); got == nil || got.Nisaab != 21708.16 || got.Krugerrand != 77626.36 {
		t.Fatalf("EconomicIndicators() = %+v", got)
	}
	if cached, err := (economic.Store{Path: cachePath}).Load(); err != nil || cached == nil {
		t.Fatalf("cached indicators = %+v, %v", cached, err)
	}
}

func TestRefreshEconomicIndicatorsWaitsUntilNineInJohannesburg(t *testing.T) {
	t.Parallel()
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		fmt.Fprint(w, economicResponse("2026-08-25T07:02:11", "25 Aug"))
	}))
	defer server.Close()

	now := time.Date(2026, 8, 25, 6, 59, 0, 0, time.UTC)
	service := &Service{
		selection: selection.State{ShowEconomicIndicators: true},
		indicators: &economic.Indicators{
			EffectiveDate:   "2026-08-24",
			SourceUpdatedAt: time.Date(2026, 8, 24, 7, 2, 11, 0, time.UTC),
			FetchedAt:       time.Date(2026, 8, 24, 19, 0, 0, 0, time.UTC),
		},
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
		fmt.Fprint(w, economicResponse("2026-08-24T07:02:11", "24 Aug"))
	}))
	defer server.Close()

	originalFetchedAt := time.Date(2026, 8, 24, 19, 0, 0, 0, time.UTC)
	now := time.Date(2026, 8, 25, 7, 0, 0, 0, time.UTC)
	service := &Service{
		selection: selection.State{ShowEconomicIndicators: true},
		indicators: &economic.Indicators{
			EffectiveDate:   "2026-08-24",
			SourceUpdatedAt: time.Date(2026, 8, 24, 7, 2, 11, 0, time.UTC),
			FetchedAt:       originalFetchedAt,
		},
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
		fmt.Fprint(w, economicResponse("2026-08-25T07:05:00", "25 Aug"))
	}))
	defer server.Close()

	now := time.Date(2026, 8, 25, 7, 30, 0, 0, time.UTC)
	service := &Service{
		selection: selection.State{ShowEconomicIndicators: true},
		indicators: &economic.Indicators{
			EffectiveDate:   "2026-08-24",
			SourceUpdatedAt: time.Date(2026, 8, 24, 7, 2, 11, 0, time.UTC),
		},
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
	if got := service.EconomicIndicators().SourceUpdatedAt; !got.Equal(time.Date(2026, 8, 25, 7, 5, 0, 0, time.UTC)) {
		t.Fatalf("SourceUpdatedAt = %v", got)
	}
}

func TestEconomicIndicatorsHiddenWhenDisabled(t *testing.T) {
	t.Parallel()
	service := &Service{indicators: &economic.Indicators{Source: economic.SourceName, SourceURL: "https://example.test", EffectiveDate: "2026-08-24", Nisaab: 1, Krugerrand: 2}}
	if got := service.EconomicIndicators(); got != nil {
		t.Fatalf("EconomicIndicators() = %+v, want nil", got)
	}
}
