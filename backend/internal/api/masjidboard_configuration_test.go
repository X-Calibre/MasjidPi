package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/hierarchy"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/scope"
)

func TestMasjidBoardHierarchyReturnsPersistedHierarchy(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hierarchy.json")
	when := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)
	state := hierarchy.State{RetrievedAt: when, ValidatedAt: when, Countries: []hierarchy.Country{{Name: "South Africa", Count: 2, Regions: []hierarchy.Region{{Name: "North West", Count: 2, Cities: []hierarchy.Location{{Name: "Brits", Count: 2}}}}}}}
	if err := hierarchy.NewStore(path).Save(state); err != nil { t.Fatal(err) }

	s := &Server{masjidBoardHierarchyPath: path}
	res := httptest.NewRecorder()
	s.masjidBoardHierarchy(res, httptest.NewRequest(http.MethodGet, "/api/masjidboard/hierarchy", nil))
	if res.Code != http.StatusOK { t.Fatalf("status=%d", res.Code) }
	var got hierarchy.State
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil { t.Fatal(err) }
	if len(got.Countries) != 1 || got.Countries[0].Regions[0].Cities[0].Name != "Brits" { t.Fatalf("got=%+v", got) }
}

func TestMasjidBoardScopeRoundTripMultipleLocations(t *testing.T) {
	path := filepath.Join(t.TempDir(), "scope.json")
	s := &Server{masjidBoardScopePath: path}
	body := []byte(`{"locations":[{"country":"South Africa","region":"North West","city":"Brits"},{"country":"South Africa","region":"Gauteng","city":"Akasia"}]}`)
	put := httptest.NewRecorder()
	s.masjidBoardScope(put, httptest.NewRequest(http.MethodPut, "/api/masjidboard/scope", bytes.NewReader(body)))
	if put.Code != http.StatusOK { t.Fatalf("PUT status=%d body=%s", put.Code, put.Body.String()) }

	get := httptest.NewRecorder()
	s.masjidBoardScope(get, httptest.NewRequest(http.MethodGet, "/api/masjidboard/scope", nil))
	if get.Code != http.StatusOK { t.Fatalf("GET status=%d", get.Code) }
	var response masjidBoardScopeResponse
	if err := json.NewDecoder(get.Body).Decode(&response); err != nil { t.Fatal(err) }
	if !response.Configured || len(response.Locations) != 2 || response.Locations[1].City != "Akasia" { t.Fatalf("response=%+v", response) }
}

func TestMasjidBoardScopeRejectsMoreThanThreeLocations(t *testing.T) {
	s := &Server{masjidBoardScopePath: filepath.Join(t.TempDir(), "scope.json")}
	state := scope.State{Locations: []scope.Location{{Country:"A",City:"1"},{Country:"A",City:"2"},{Country:"A",City:"3"},{Country:"A",City:"4"}}}
	data, _ := json.Marshal(state)
	res := httptest.NewRecorder()
	s.masjidBoardScope(res, httptest.NewRequest(http.MethodPut, "/api/masjidboard/scope", bytes.NewReader(data)))
	if res.Code != http.StatusBadRequest { t.Fatalf("status=%d", res.Code) }
}
