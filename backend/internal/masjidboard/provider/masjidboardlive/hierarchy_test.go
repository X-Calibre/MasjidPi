package masjidboardlive

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHierarchyCountries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm() error = %v", err)
		}
		if r.PostForm.Get("type") != "country" {
			t.Fatalf("type = %q", r.PostForm.Get("type"))
		}
		_, _ = w.Write([]byte(`[["South Africa","615"],["Australia","17"]]`))
	}))
	defer server.Close()

	client := DiscoveryClient{HTTPClient: server.Client(), Endpoint: server.URL}
	entries, err := client.Countries(context.Background())
	if err != nil {
		t.Fatalf("Countries() error = %v", err)
	}
	if len(entries) != 2 || entries[0].Name != "South Africa" || entries[0].Count != 615 {
		t.Fatalf("entries = %+v", entries)
	}
}

func TestHierarchyRegionsPreservesBlankBucket(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm() error = %v", err)
		}
		if r.PostForm.Get("type") != "province" || r.PostForm.Get("search") != "South Africa" || r.PostForm.Get("countryName") != "South Africa" {
			t.Fatalf("form = %#v", r.PostForm)
		}
		_, _ = w.Write([]byte(`[["Gauteng","277"],["Limpopo","25"],["","1"],["Limpopo","1"]]`))
	}))
	defer server.Close()

	client := DiscoveryClient{HTTPClient: server.Client(), Endpoint: server.URL}
	entries, err := client.Regions(context.Background(), "South Africa")
	if err != nil {
		t.Fatalf("Regions() error = %v", err)
	}
	if len(entries) != 4 || entries[2].Name != "" || entries[2].Count != 1 {
		t.Fatalf("entries = %+v", entries)
	}
}

func TestHierarchyCitiesUsesPrimaryPairRows(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm() error = %v", err)
		}
		if r.PostForm.Get("type") != "cityProvince" || r.PostForm.Get("search") != "North West" {
			t.Fatalf("form = %#v", r.PostForm)
		}
		if r.PostForm.Get("countryName") != "South Africa" || r.PostForm.Get("provinceName") != "North West" {
			t.Fatalf("location form = %#v", r.PostForm)
		}
		_, _ = w.Write([]byte(`[[["Brits","3"],["Rustenburg","4"]],[["B","1"],["R","1"]]]`))
	}))
	defer server.Close()

	client := DiscoveryClient{HTTPClient: server.Client(), Endpoint: server.URL}
	entries, err := client.Cities(context.Background(), "South Africa", "North West")
	if err != nil {
		t.Fatalf("Cities() error = %v", err)
	}
	if len(entries) != 2 || entries[0].Name != "Brits" || entries[0].Count != 3 {
		t.Fatalf("entries = %+v", entries)
	}
}

func TestHierarchyCitiesSupportsObjectRows(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm() error = %v", err)
		}
		if r.PostForm.Get("type") != "cityProvince" || r.PostForm.Get("search") != "Victoria" {
			t.Fatalf("form = %#v", r.PostForm)
		}
		if r.PostForm.Get("countryName") != "Australia" || r.PostForm.Get("provinceName") != "Victoria" {
			t.Fatalf("location form = %#v", r.PostForm)
		}
		_, _ = w.Write([]byte(`[[{"0":"Coburg North","city":"Coburg North","1":"1","COUNT(*)":"1"},{"0":"Fawkner","city":"Fawkner","1":"3","COUNT(*)":"3"},{"0":"Wallan","city":"Wallan","1":"2","COUNT(*)":"2"}],[{"0":"","LEFT(SUBQUERY.city,1)":"","1":"3","COUNT(LEFT(SUBQUERY.city,1))":"3"}]]`))
	}))
	defer server.Close()

	client := DiscoveryClient{HTTPClient: server.Client(), Endpoint: server.URL}
	entries, err := client.Cities(context.Background(), "Australia", "Victoria")
	if err != nil {
		t.Fatalf("Cities() error = %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("entries = %+v", entries)
	}
	if entries[0].Name != "Coburg North" || entries[0].Count != 1 {
		t.Fatalf("entries[0] = %+v", entries[0])
	}
	if entries[1].Name != "Fawkner" || entries[1].Count != 3 {
		t.Fatalf("entries[1] = %+v", entries[1])
	}
	if entries[2].Name != "Wallan" || entries[2].Count != 2 {
		t.Fatalf("entries[2] = %+v", entries[2])
	}
}

func TestHierarchyRejectsNullResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`null`))
	}))
	defer server.Close()

	client := DiscoveryClient{HTTPClient: server.Client(), Endpoint: server.URL}
	if _, err := client.Countries(context.Background()); err == nil {
		t.Fatal("Countries() expected error for null response")
	}
}

func TestHierarchyRejectsInvalidCount(t *testing.T) {
	if _, err := parseHierarchyPairs([]byte(`[["South Africa","bad"]]`), false); err == nil {
		t.Fatal("parseHierarchyPairs() expected error for invalid count")
	}
}

func TestCityHierarchyRejectsMissingObjectFields(t *testing.T) {
	if _, err := parseCityHierarchyRows([]byte(`[{"city":"Fawkner"}]`)); err == nil {
		t.Fatal("parseCityHierarchyRows() expected error for missing count")
	}
}
