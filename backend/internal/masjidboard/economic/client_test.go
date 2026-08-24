package economic

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientFetchParsesLatestIndicatorRowByHeading(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{
            "date":"2026-08-14T09:02:11",
            "link":"https://www.jamiatsa.org/rabi-al-awwal-1448-2/",
            "title":{"rendered":"Rabi &#8216;al Awwal 1448"},
            "content":{"rendered":"<table><thead><tr><th>Hijri</th><th>Date</th><th>Rand-Dollar</th><th>24 Carat</th><th>22 Carat</th><th>18 Carat</th><th>Silver</th><th>Nisaab</th><th>Min Mahr</th><th>Mahr Faatimi</th><th>Krugerrand</th></tr></thead><tbody><tr><td>11</td><td>24 Aug</td><td>R16.01</td><td>R2385.85</td><td>R2187.03</td><td>R1789.39</td><td>R35.45</td><td>R21708.16</td><td>R1085.40</td><td>R54270.41</td><td>R77626.36</td></tr><tr><td>8</td><td>21 Aug</td><td>R16.06</td><td>R2356.62</td><td>R2160.24</td><td>R1767.47</td><td>R35.60</td><td>R21800.02</td><td>R1090.00</td><td>R54500.04</td><td>R76538.96</td></tr></tbody></table>"}
        }]`)
	}))
	defer server.Close()

	fetchedAt := time.Date(2026, 8, 24, 19, 0, 0, 0, time.UTC)
	got, err := (Client{APIURL: server.URL, HTTPClient: server.Client(), Now: func() time.Time { return fetchedAt }}).Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if got.EffectiveDate != "2026-08-24" || got.HijriDate != "11 Rabi ‘al Awwal 1448" {
		t.Fatalf("dates = %q, %q", got.EffectiveDate, got.HijriDate)
	}
	if got.Nisaab != 21708.16 || got.Krugerrand != 77626.36 || got.Gold24Carat != 2385.85 || got.Silver != 35.45 {
		t.Fatalf("unexpected values: %+v", got)
	}
	if got.Source != SourceName || got.FetchedAt != fetchedAt {
		t.Fatalf("metadata = %+v", got)
	}
}

func TestClientFetchRejectsMissingRequiredColumn(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"date":"2026-08-14T09:02:11","link":"https://example.test","title":{"rendered":"Rabi 1448"},"content":{"rendered":"<table><thead><tr><th>Hijri</th><th>Date</th></tr></thead><tbody><tr><td>11</td><td>24 Aug</td></tr></tbody></table>"}}]`)
	}))
	defer server.Close()
	if _, err := (Client{APIURL: server.URL, HTTPClient: server.Client()}).Fetch(context.Background()); err == nil {
		t.Fatal("Fetch() expected missing-column error")
	}
}
