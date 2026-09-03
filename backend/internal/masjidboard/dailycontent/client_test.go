package dailycontent

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const validResponse = `let translations = {
  "ayahSurah":{"en":" Surah 59 – Al-Hashr "},
  "AyahNo":{"en":"Āyah 7"},
  "ayah":{"en":"First  line &amp; meaning"},
  "hadithHeading":{"en":"Hadīth"},
  "hadith":{"en":"A <b>short</b> narration."},
  "hadithRef":{"en":"Bukhāri"},
  "sunnahHeading":{"en":"Sunnah of Travelling – 3"},
  "sunnah":{"en":"First action.<br><br><br>Second action."},
  "sunnahRef":{"en":"Muslim"},
  "Thu Sep 03 2026 00:00:00 GMT+0200 (South Africa Standard Time)":{"en":""}
};`

func TestClientFetchesAndNormalizesContent(t *testing.T) {
	t.Parallel()
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		if got := r.Header.Get("User-Agent"); got != "MasjidPi Daily Islamic Content" {
			t.Errorf("User-Agent = %q", got)
		}
		fmt.Fprint(w, validResponse)
	}))
	defer server.Close()
	fetchedAt := time.Date(2026, 9, 3, 18, 0, 0, 0, time.UTC)
	content, err := (Client{APIURL: server.URL, HTTPClient: server.Client(), Now: func() time.Time { return fetchedAt }}).Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if requests.Load() != 1 || content.ContentDate != "2026-09-03" || content.FetchedAt != fetchedAt {
		t.Fatalf("metadata = %+v; requests = %d", content, requests.Load())
	}
	if content.Ayah.Text != "First line & meaning" || content.Hadith.Text != "A short narration." {
		t.Fatalf("normalized content = %+v", content)
	}
	if content.Sunnah.Text != "First action.\n\nSecond action." {
		t.Fatalf("Sunnah text = %q", content.Sunnah.Text)
	}
}

func TestClientRejectsMalformedOrIncompleteResponses(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		body string
	}{
		{"not assignment", `{}`},
		{"invalid JSON", `let translations = {`},
		{"missing field", strings.Replace(validResponse, `"ayah":`, `"missingAyah":`, 1)},
		{"missing language value", strings.Replace(validResponse, `"ayah":{"en":"First  line &amp; meaning"}`, `"ayah":{"ar":"نص"}`, 1)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, test.body) }))
			defer server.Close()
			if _, err := (Client{APIURL: server.URL, HTTPClient: server.Client()}).Fetch(context.Background()); err == nil {
				t.Fatal("Fetch() expected an error")
			}
		})
	}
}

func TestClientRejectsUnexpectedStatusAndOversizedResponse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status int
		body   string
	}{
		{"status", http.StatusServiceUnavailable, "unavailable"},
		{"oversized", http.StatusOK, strings.Repeat("x", maxResponseSize+1)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(test.status)
				fmt.Fprint(w, test.body)
			}))
			defer server.Close()
			if _, err := (Client{APIURL: server.URL, HTTPClient: server.Client()}).Fetch(context.Background()); err == nil {
				t.Fatal("Fetch() expected an error")
			}
		})
	}
}
