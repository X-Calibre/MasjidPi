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
            "content":{"rendered":"<table><thead><tr><th>Hijri</th><th>Date</th><th>Rand-Dollar</th><th>24 Carat</th><th>22 Carat</th><th>18 Carat</th><th>14 Carat</th><th>9 Carat</th><th>Silver</th><th>Nisaab</th><th>Min Mahr</th><th>Mahr Faatimi</th><th>Krugerrand</th></tr></thead><tbody><tr><td>11</td><td>24 Aug</td><td>R16.01</td><td>R2385.85</td><td>R2187.03</td><td>R1789.39</td><td>R1391.75</td><td>R894.69</td><td>R35.45</td><td>R21708.16</td><td>R1085.40</td><td>R54270.41</td><td>R77626.36</td></tr><tr><td>8</td><td>21 Aug</td><td>R16.06</td><td>R2356.62</td><td>R2160.24</td><td>R1767.47</td><td>R1374.70</td><td>R883.73</td><td>R35.60</td><td>R21800.02</td><td>R1090.00</td><td>R54500.04</td><td>R76538.96</td></tr></tbody></table>"}
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
	if got.RandDollar != 16.01 || got.Gold14Carat != 1391.75 || got.Gold9Carat != 894.69 ||
		got.Nisaab != 21708.16 || got.Krugerrand != 77626.36 || got.Gold24Carat != 2385.85 || got.Silver != 35.45 {
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

func TestParseEffectiveDateAcceptsJamiatFormatsAndEveryMonth(t *testing.T) {
	t.Parallel()
	tests := []struct {
		value    string
		postYear int
		want     string
	}{
		{"01 Jan", 2026, "2026-01-01"},
		{"02 Feb", 2026, "2026-02-02"},
		{"03 Mar", 2026, "2026-03-03"},
		{"04 Apr", 2026, "2026-04-04"},
		{"05 May", 2026, "2026-05-05"},
		{"06 Jun", 2026, "2026-06-06"},
		{"07 Jul", 2026, "2026-07-07"},
		{"08 Aug", 2026, "2026-08-08"},
		{"09 Sep", 2026, "2026-09-09"},
		{"10 Oct", 2026, "2026-10-10"},
		{"11 Nov", 2026, "2026-11-11"},
		{"12 Dec", 2026, "2026-12-12"},
		{"02 Sept", 2026, "2026-09-02"},
		{"16-Jun", 2026, "2026-06-16"},
		{"30–Apr", 2026, "2026-04-30"},
		{"13 Jul 18", 2018, "2018-07-13"},
		{"12 Mar '21", 2021, "2021-03-12"},
		{"2 September 2026", 2026, "2026-09-02"},
	}
	for _, test := range tests {
		test := test
		t.Run(test.value, func(t *testing.T) {
			t.Parallel()
			got, err := parseEffectiveDate(test.value, test.postYear)
			if err != nil {
				t.Fatalf("parseEffectiveDate() error = %v", err)
			}
			if formatted := got.Format("2006-01-02"); formatted != test.want {
				t.Fatalf("parseEffectiveDate() = %s, want %s", formatted, test.want)
			}
		})
	}
}

func TestParseEffectiveDateRejectsInvalidValues(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"", "Sept", "31 Feb", "2 Smarch", "2 Sep 2", "2 Sep twenty"} {
		if _, err := parseEffectiveDate(value, 2026); err == nil {
			t.Errorf("parseEffectiveDate(%q) expected error", value)
		}
	}
}
