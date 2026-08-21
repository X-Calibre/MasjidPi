package masjidboardlive

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestDiscoveryClientSearchMasjids(t *testing.T) {
	var gotForm url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "application/x-www-form-urlencoded") {
			t.Fatalf("Content-Type = %q", got)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm() error = %v", err)
		}
		gotForm = r.PostForm
		_, _ = w.Write([]byte(`[
  "Brits",
  {
    "masjid":"Brits Jamia Masjid",
    "fajr_jamaat":"06:00",
    "zuhr_jamaat":"13:20",
    "asr_jamaat":"17:00",
    "maghrib_adhan":"17:54",
    "esha_jamaat":"19:30",
    "last_updated":"Sun, 22 Mar 2026, 12:47:25",
    "MBL_ID":"MBL11517PRP",
    "city":"Brits",
    "sunset":"17:51",
    "time_zone_milli":"7200000",
    "web_url":"brits-jamia",
    "jumuah_khutbah":"13:00",
    "ramadhaanactive":"Hide",
    "date_adjust":"0",
    "moon_seen":"Y",
    "ladies_facility":"Yes"
  }
]`))
	}))
	defer server.Close()

	client := DiscoveryClient{HTTPClient: server.Client(), Endpoint: server.URL}
	result, err := client.SearchMasjids(context.Background(), MasjidSearch{
		Search:   "Brits",
		Country:  "South Africa",
		Province: "North West",
		City:     "Brits",
	})
	if err != nil {
		t.Fatalf("SearchMasjids() error = %v", err)
	}

	if gotForm.Get("type") != "masjid" || gotForm.Get("search") != "Brits" {
		t.Fatalf("form = %#v", gotForm)
	}
	if gotForm.Get("countryName") != "South Africa" || gotForm.Get("provinceName") != "North West" || gotForm.Get("cityName") != "Brits" {
		t.Fatalf("location form = %#v", gotForm)
	}
	if result.Location != "Brits" {
		t.Fatalf("Location = %q", result.Location)
	}
	if len(result.Entries) != 1 {
		t.Fatalf("Entries = %d, want 1", len(result.Entries))
	}
	entry := result.Entries[0]
	if entry.Name != "Brits Jamia Masjid" || entry.WebURL != "brits-jamia" {
		t.Fatalf("entry = %+v", entry)
	}
	if entry.MBLID != "MBL11517PRP" || entry.TimeZoneOffsetMS != 7200000 {
		t.Fatalf("entry metadata = %+v", entry)
	}
	if entry.JumuahKhutbah != "13:00" || entry.LadiesFacility != "Yes" {
		t.Fatalf("entry optional fields = %+v", entry)
	}
}

func TestParseMasjidDiscoveryRejectsEmptyResponse(t *testing.T) {
	if _, err := parseMasjidDiscovery(nil); err == nil {
		t.Fatal("parseMasjidDiscovery() expected an error for empty response")
	}
}

func TestDiscoveryClientRejectsNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}))
	defer server.Close()

	client := DiscoveryClient{HTTPClient: server.Client(), Endpoint: server.URL}
	if _, err := client.SearchMasjids(context.Background(), MasjidSearch{}); err == nil {
		t.Fatal("SearchMasjids() expected an error for non-2xx response")
	}
}

func TestParseMasjidDiscoveryRejectsMissingIdentity(t *testing.T) {
	rows := []byte(`["Brits",{"masjid":"","web_url":""}]`)
	var raw []json.RawMessage
	if err := json.Unmarshal(rows, &raw); err != nil {
		t.Fatalf("decode test rows: %v", err)
	}
	if _, err := parseMasjidDiscovery(raw); err == nil {
		t.Fatal("parseMasjidDiscovery() expected an error for missing identity")
	}
}

func TestParseMasjidDiscoveryRejectsInvalidOffset(t *testing.T) {
	rows := []byte(`["Brits",{"masjid":"Test Masjid","web_url":"test-masjid","time_zone_milli":"bad"}]`)
	var raw []json.RawMessage
	if err := json.Unmarshal(rows, &raw); err != nil {
		t.Fatalf("decode test rows: %v", err)
	}
	if _, err := parseMasjidDiscovery(raw); err == nil {
		t.Fatal("parseMasjidDiscovery() expected an error for invalid timezone offset")
	}
}
