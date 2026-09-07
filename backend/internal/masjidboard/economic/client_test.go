package economic

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const validAPIResponse = `{
  "id":21226,
  "gregorian_date":"2026-09-07",
  "hijri_day":25,
  "hijri_month":3,
  "hijri_year":1448,
  "hijri_month_name":"Rabi' al-Awwal",
  "usd_zar":"15.9637",
  "gold_24k":"2266.5684",
  "gold_22k":"2077.6877",
  "gold_21k":"1983.2474",
  "gold_18k":"1699.9263",
  "gold_14k":"1322.1649",
  "gold_9k":"849.9632",
  "silver":"33.9823",
  "nisaab":"20809.40",
  "mahr_min":"1040.47",
  "mahr_faatimi":"52023.50",
  "krugerrand":"73485.25",
  "notes":"Weekend — prices copied from 2026-09-04",
  "created_by":10,
  "created_at":"2026-09-07T08:33:23.766Z",
  "updated_at":"2026-09-07T08:33:23.766Z",
  "prices_pulled_at":"2026-09-07T08:33:15.037Z"
}`

func TestClientFetchParsesLatestIndicator(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Errorf("Accept = %q", got)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, validAPIResponse)
	}))
	defer server.Close()

	fetchedAt := time.Date(2026, 9, 7, 9, 30, 0, 0, time.UTC)
	got, err := (Client{APIURL: server.URL, HTTPClient: server.Client(), Now: func() time.Time { return fetchedAt }}).Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if got.EffectiveDate != "2026-09-07" || got.HijriDate != "25 Rabi' al-Awwal 1448" {
		t.Fatalf("dates = %q, %q", got.EffectiveDate, got.HijriDate)
	}
	if got.RandDollar != 15.9637 || got.Gold24Carat != 2266.5684 || got.Gold22Carat != 2077.6877 ||
		got.Gold21Carat != 1983.2474 || got.Gold18Carat != 1699.9263 || got.Gold14Carat != 1322.1649 || got.Gold9Carat != 849.9632 ||
		got.Silver != 33.9823 || got.Nisaab != 20809.40 || got.MinimumMahr != 1040.47 ||
		got.MahrFaatimi != 52023.50 || got.Krugerrand != 73485.25 {
		t.Fatalf("unexpected values: %+v", got)
	}
	wantUpdatedAt := time.Date(2026, 9, 7, 8, 33, 23, 766000000, time.UTC)
	if got.Source != SourceName || got.SourceURL != SourcePageURL || got.FetchedAt != fetchedAt ||
		!got.UpdatedAt.Equal(wantUpdatedAt) || got.Notes != "Weekend — prices copied from 2026-09-04" {
		t.Fatalf("metadata = %+v", got)
	}
	if !got.Valid() || !got.Complete() {
		t.Fatalf("indicator should be valid and complete: %+v", got)
	}
}

func TestClientFetchRejectsHTMLResponse(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		fmt.Fprint(w, "<!DOCTYPE html><html></html>")
	}))
	defer server.Close()

	_, err := (Client{APIURL: server.URL, HTTPClient: server.Client()}).Fetch(context.Background())
	if err == nil || !strings.Contains(err.Error(), `unexpected content type "text/html; charset=UTF-8"`) {
		t.Fatalf("Fetch() error = %v", err)
	}
}

func TestClientFetchRejectsNonSuccessStatus(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, `{"error":"upstream unavailable"}`)
	}))
	defer server.Close()

	_, err := (Client{APIURL: server.URL, HTTPClient: server.Client()}).Fetch(context.Background())
	if err == nil || !strings.Contains(err.Error(), "unexpected HTTP status 502 Bad Gateway") {
		t.Fatalf("Fetch() error = %v", err)
	}
}

func TestClientFetchRejectsMalformedJSON(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{not-json}`)
	}))
	defer server.Close()

	_, err := (Client{APIURL: server.URL, HTTPClient: server.Client()}).Fetch(context.Background())
	if err == nil || !strings.Contains(err.Error(), "decode response") {
		t.Fatalf("Fetch() error = %v", err)
	}
}

func TestClientFetchRejectsInvalidDate(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, strings.Replace(validAPIResponse, `"gregorian_date":"2026-09-07"`, `"gregorian_date":"2026-02-30"`, 1))
	}))
	defer server.Close()

	_, err := (Client{APIURL: server.URL, HTTPClient: server.Client()}).Fetch(context.Background())
	if err == nil || !strings.Contains(err.Error(), "invalid gregorian_date") {
		t.Fatalf("Fetch() error = %v", err)
	}
}

func TestClientFetchRejectsInvalidUpdatedAt(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, strings.Replace(validAPIResponse, `"updated_at":"2026-09-07T08:33:23.766Z"`, `"updated_at":"not-a-time"`, 1))
	}))
	defer server.Close()

	_, err := (Client{APIURL: server.URL, HTTPClient: server.Client()}).Fetch(context.Background())
	if err == nil || !strings.Contains(err.Error(), "invalid updated_at") {
		t.Fatalf("Fetch() error = %v", err)
	}
}

func TestClientFetchRejectsInvalidNumericValue(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, strings.Replace(validAPIResponse, `"nisaab":"20809.40"`, `"nisaab":"not-a-number"`, 1))
	}))
	defer server.Close()

	_, err := (Client{APIURL: server.URL, HTTPClient: server.Client()}).Fetch(context.Background())
	if err == nil || !strings.Contains(err.Error(), "invalid nisaab value") {
		t.Fatalf("Fetch() error = %v", err)
	}
}
