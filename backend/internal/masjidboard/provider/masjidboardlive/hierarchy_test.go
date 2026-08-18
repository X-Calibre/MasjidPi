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

func TestHierarchyCitiesUsesPrimaryRows(t *testing.T) {
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
