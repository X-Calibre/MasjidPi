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

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/dailycontent"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

const dailyResponse = `let translations = {
"ayahSurah":{"en":"Surah 1"},"AyahNo":{"en":"Ayah 1"},"ayah":{"en":"Ayah text"},
"hadithHeading":{"en":"Hadith"},"hadith":{"en":"Hadith text"},"hadithRef":{"en":"Bukhari"},
"sunnahHeading":{"en":"Sunnah"},"sunnah":{"en":"Sunnah text"},"sunnahRef":{"en":"Muslim"},
"Thu Sep 03 2026 00:00:00 GMT+0200 (South Africa Standard Time)":{"en":""}}`

func TestRefreshDailyIslamicContentFetchesOncePerJohannesburgDay(t *testing.T) {
	t.Parallel()
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		fmt.Fprint(w, dailyResponse)
	}))
	defer server.Close()
	now := time.Date(2026, 9, 3, 22, 30, 0, 0, time.UTC) // 4 September in Johannesburg.
	path := filepath.Join(t.TempDir(), "daily.json")
	service := &Service{
		selection:          selection.State{},
		dailyContentClient: dailycontent.Client{APIURL: server.URL, HTTPClient: server.Client(), Now: func() time.Time { return now }},
		dailyContentStore:  dailycontent.Store{Path: path},
	}
	if err := service.RefreshDailyIslamicContent(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := service.RefreshDailyIslamicContent(context.Background()); err != nil {
		t.Fatal(err)
	}
	if requests.Load() != 1 || service.DailyIslamicContent() == nil {
		t.Fatalf("requests = %d, content = %+v", requests.Load(), service.DailyIslamicContent())
	}
	if cached, err := (dailycontent.Store{Path: path}).Load(); err != nil || cached == nil {
		t.Fatalf("cached content = %+v, %v", cached, err)
	}
}

func TestRefreshDailyIslamicContentRetriesNextDay(t *testing.T) {
	t.Parallel()
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		fmt.Fprint(w, dailyResponse)
	}))
	defer server.Close()
	now := time.Date(2026, 9, 3, 10, 0, 0, 0, time.UTC)
	service := &Service{
		selection:          selection.State{},
		dailyContentClient: dailycontent.Client{APIURL: server.URL, HTTPClient: server.Client(), Now: func() time.Time { return now }},
		dailyContentStore:  dailycontent.Store{Path: filepath.Join(t.TempDir(), "daily.json")},
	}
	if err := service.RefreshDailyIslamicContent(context.Background()); err != nil {
		t.Fatal(err)
	}
	now = now.Add(24 * time.Hour)
	if err := service.RefreshDailyIslamicContent(context.Background()); err != nil {
		t.Fatal(err)
	}
	if requests.Load() != 2 {
		t.Fatalf("requests = %d, want 2", requests.Load())
	}
}

func TestRefreshDailyIslamicContentDisabledAndFailureFallback(t *testing.T) {
	t.Parallel()
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()
	falseValue := false
	cached := dailycontent.Content{
		Ayah:     dailycontent.Ayah{Surah: "Surah 1", AyahNumber: "Ayah 1", Text: "Cached"},
		Hadith:   dailycontent.Hadith{Heading: "Hadith", Text: "Cached"},
		Sunnah:   dailycontent.Sunnah{Heading: "Sunnah", Text: "Cached"},
		Language: "en", Source: dailycontent.SourceName, SourceURL: dailycontent.SourceURL,
		FetchedAt: time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC),
	}
	service := &Service{
		selection:    selection.State{ShowDailyAyah: &falseValue, ShowDailyHadith: &falseValue, ShowDailySunnah: &falseValue},
		dailyContent: &cached,
		dailyContentClient: dailycontent.Client{APIURL: server.URL, HTTPClient: server.Client(), Now: func() time.Time {
			return time.Date(2026, 9, 3, 10, 0, 0, 0, time.UTC)
		}},
	}
	if err := service.RefreshDailyIslamicContent(context.Background()); err != nil || requests.Load() != 0 {
		t.Fatalf("disabled refresh = %v; requests = %d", err, requests.Load())
	}
	service.selection.ShowDailyAyah = nil // Default-enabled migration behavior.
	if err := service.RefreshDailyIslamicContent(context.Background()); err == nil {
		t.Fatal("enabled refresh expected upstream error")
	}
	if requests.Load() != 1 || service.DailyIslamicContent() == nil || service.DailyIslamicContent().Ayah.Text != "Cached" {
		t.Fatalf("fallback content = %+v; requests = %d", service.DailyIslamicContent(), requests.Load())
	}
}
