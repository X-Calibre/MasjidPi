package masjidboardlive

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/provider"
)

var _ provider.Provider = CoreClient{}

func TestCoreClientFetchAt(t *testing.T) {
	fixture := string(loadCoreFixture(t))
	upcoming := `<h3 id="fajrNextDate">15 Sep</h3><h3 id="fajrNextTime">05:15</h3>`
	var gotQuery string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<html><script>let data = " + fixture + "</script>" + upcoming + "</html>"))
	}))
	defer server.Close()

	client := CoreClient{
		HTTPClient: server.Client(),
		Endpoint:   server.URL,
		WebURL:     "brits-jamia",
		Identity:   coreIdentity(),
	}

	now := time.Date(2026, 8, 18, 16, 0, 0, 0, time.UTC)
	result, err := client.FetchAt(context.Background(), now)
	if err != nil {
		t.Fatalf("FetchAt() error = %v", err)
	}
	if gotQuery != "brits-jamia" {
		t.Fatalf("query = %q, want %q", gotQuery, "brits-jamia")
	}
	if result.Metadata.MBLNumber != "MBL11517PRP" {
		t.Fatalf("MBLNumber = %q", result.Metadata.MBLNumber)
	}
	assertCoreClock(t, "Fajr Jamaah", result.Board.PrayerTimes.Fajr.Jamaah, 6, 0)
	if len(result.Board.Notices) != 1 || result.Board.Notices[0].Fields["new_time"] != "05:15" {
		t.Fatalf("upcoming Salaah changes = %+v", result.Board.Notices)
	}
}

func TestCoreClientFetchReturnsNormalisedBoard(t *testing.T) {
	fixture := string(loadCoreFixture(t))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("<html><script>let data = " + fixture + "</script></html>"))
	}))
	defer server.Close()

	client := CoreClient{
		HTTPClient: server.Client(),
		Endpoint:   server.URL,
		WebURL:     "brits-jamia",
		Identity:   coreIdentity(),
	}

	board, err := client.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if board.Identity.ID != "brits-jamia" {
		t.Fatalf("Identity.ID = %q", board.Identity.ID)
	}
	assertCoreClock(t, "Fajr Jamaah", board.PrayerTimes.Fajr.Jamaah, 6, 0)
}

func TestCoreClientRejectsNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	client := CoreClient{
		HTTPClient: server.Client(),
		Endpoint:   server.URL,
		WebURL:     "missing-board",
		Identity:   coreIdentity(),
	}

	if _, err := client.FetchAt(context.Background(), time.Now()); err == nil {
		t.Fatal("FetchAt() expected an error for non-2xx response")
	}
}

func TestCoreClientRejectsPageWithoutCoreData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html><script>let other = {};</script></html>`))
	}))
	defer server.Close()

	client := CoreClient{
		HTTPClient: server.Client(),
		Endpoint:   server.URL,
		WebURL:     "brits-jamia",
		Identity:   coreIdentity(),
	}

	if _, err := client.FetchAt(context.Background(), time.Now()); err == nil {
		t.Fatal("FetchAt() expected an error for page without Core data")
	}
}

func TestCoreClientRejectsInvalidCoreData(t *testing.T) {
	fixture := string(loadCoreFixture(t))
	fixture = strings.Replace(fixture, `fajrAthan : "05:40"`, `fajrAthan : ""`, 1)
	fixture = strings.Replace(fixture, `fajrJamaah : "06:00"`, `fajrJamaah : "~~~~"`, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("<html><script>let data = " + fixture + "</script></html>"))
	}))
	defer server.Close()

	client := CoreClient{
		HTTPClient: server.Client(),
		Endpoint:   server.URL,
		WebURL:     "brits-jamia",
		Identity:   coreIdentity(),
	}

	if _, err := client.FetchAt(context.Background(), time.Now()); err == nil {
		t.Fatal("FetchAt() expected an error for invalid Core data")
	}
}

func TestCoreClientRejectsMissingWebURL(t *testing.T) {
	client := CoreClient{Identity: coreIdentity()}
	if _, err := client.FetchAt(context.Background(), time.Now()); err == nil {
		t.Fatal("FetchAt() expected an error for missing web URL")
	}
}
