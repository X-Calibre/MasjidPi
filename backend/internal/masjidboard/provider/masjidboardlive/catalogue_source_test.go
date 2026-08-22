package masjidboardlive

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	masjidboardcatalogue "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/catalogue"
)

func TestCatalogueSourceFetchMapsDiscoveryRecords(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.PostForm.Get("type") != "masjid" || r.PostForm.Get("countryName") != "South Africa" || r.PostForm.Get("provinceName") != "North West" || r.PostForm.Get("cityName") != "Brits" {
			t.Fatalf("form=%#v", r.PostForm)
		}
		_, _ = w.Write([]byte(`["Brits",{"masjid":"Brits Jamia Masjid","MBL_ID":"MBL11517PRP","city":"Brits","time_zone_milli":"7200000","web_url":"brits-jamia","last_updated":"Sun, 22 Mar 2026, 12:47:25"}]`))
	}))
	defer server.Close()

	now := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)
	source := CatalogueSource{Client: DiscoveryClient{HTTPClient: server.Client(), Endpoint: server.URL}}
	got, err := source.Fetch(context.Background(), masjidboardcatalogue.Location{Country: "South Africa", Region: "North West", City: "Brits"}, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Records) != 1 {
		t.Fatalf("records=%+v", got.Records)
	}
	record := got.Records[0]
	if record.ID != "masjidboardlive:brits-jamia" || record.Name != "Brits Jamia Masjid" || record.Region != "North West" || record.Country != "South Africa" || record.TimeZoneOffsetMS != 7200000 {
		t.Fatalf("record=%+v", record)
	}
	if record.ProviderMetadata["mbl_id"] != "MBL11517PRP" || record.ProviderMetadata["last_updated"] == "" {
		t.Fatalf("metadata=%+v", record.ProviderMetadata)
	}
	if !got.RetrievedAt.Equal(now) || !got.ValidatedAt.Equal(now) {
		t.Fatalf("timestamps=%+v", got)
	}
}
