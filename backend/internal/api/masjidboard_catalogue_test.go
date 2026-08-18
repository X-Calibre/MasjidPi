package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	masjidboardcatalogue "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/catalogue"
)

func writeMasjidBoardCatalogueFixture(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "masjidboard_catalogue.json")
	when := time.Date(2026, 8, 18, 18, 0, 0, 0, time.UTC)
	state := masjidboardcatalogue.Catalogue{
		RetrievedAt: when,
		ValidatedAt: when,
		Records: []masjidboardcatalogue.Record{
			{
				ID:               "masjidboardlive:brits-jamia",
				Provider:         "masjidboardlive",
				ExternalID:       "brits-jamia",
				Name:             "Brits Jamia Masjid",
				City:             "Brits",
				Region:           "North West",
				Country:          "South Africa",
				TimeZoneOffsetMS: 7200000,
				DiscoveredAt:     when,
				LastSeenAt:       when,
				Status:           masjidboardcatalogue.StatusActive,
			},
			{
				ID:               "masjidboardlive:brits-taqwa",
				Provider:         "masjidboardlive",
				ExternalID:       "brits-taqwa",
				Name:             "Masjid Taqwa",
				City:             "Brits",
				Region:           "North West",
				Country:          "South Africa",
				TimeZoneOffsetMS: 7200000,
				DiscoveredAt:     when,
				LastSeenAt:       when,
				Status:           masjidboardcatalogue.StatusMissing,
			},
			{
				ID:               "masjidboardlive:fawkner-masjid",
				Provider:         "masjidboardlive",
				ExternalID:       "fawkner-masjid",
				Name:             "Fawkner Masjid",
				City:             "Fawkner",
				Region:           "Victoria",
				Country:          "Australia",
				TimeZoneOffsetMS: 36000000,
				DiscoveredAt:     when,
				LastSeenAt:       when,
				Status:           masjidboardcatalogue.StatusActive,
			},
		},
	}
	if err := masjidboardcatalogue.NewStore(path).Save(state); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	return path
}

func decodeMasjidBoardCatalogueResponse(t *testing.T, response *httptest.ResponseRecorder) masjidBoardCatalogueResponse {
	t.Helper()
	var body masjidBoardCatalogueResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return body
}

func TestMasjidBoardCatalogueReturnsLocalCatalogue(t *testing.T) {
	s := &Server{masjidBoardCataloguePath: writeMasjidBoardCatalogueFixture(t)}
	req := httptest.NewRequest(http.MethodGet, "/api/masjidboard/catalogue", nil)
	res := httptest.NewRecorder()

	s.masjidBoardCatalogue(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}
	body := decodeMasjidBoardCatalogueResponse(t, res)
	if body.Count != 3 || len(body.Records) != 3 {
		t.Fatalf("response count = %d records = %d, want 3", body.Count, len(body.Records))
	}
	if body.RetrievedAt == nil || body.ValidatedAt == nil {
		t.Fatal("expected catalogue timestamps")
	}
}

func TestMasjidBoardCatalogueGlobalSearch(t *testing.T) {
	s := &Server{masjidBoardCataloguePath: writeMasjidBoardCatalogueFixture(t)}
	req := httptest.NewRequest(http.MethodGet, "/api/masjidboard/catalogue?q=victoria", nil)
	res := httptest.NewRecorder()

	s.masjidBoardCatalogue(res, req)

	body := decodeMasjidBoardCatalogueResponse(t, res)
	if body.Count != 1 || body.Records[0].ExternalID != "fawkner-masjid" {
		t.Fatalf("records = %+v", body.Records)
	}
}

func TestMasjidBoardCatalogueCombinesFilters(t *testing.T) {
	s := &Server{masjidBoardCataloguePath: writeMasjidBoardCatalogueFixture(t)}
	req := httptest.NewRequest(http.MethodGet, "/api/masjidboard/catalogue?city=brits&province=north%20west&country=south%20africa&status=missing", nil)
	res := httptest.NewRecorder()

	s.masjidBoardCatalogue(res, req)

	body := decodeMasjidBoardCatalogueResponse(t, res)
	if body.Count != 1 || body.Records[0].ExternalID != "brits-taqwa" {
		t.Fatalf("records = %+v", body.Records)
	}
}

func TestMasjidBoardCatalogueNameFilterIsCaseInsensitive(t *testing.T) {
	s := &Server{masjidBoardCataloguePath: writeMasjidBoardCatalogueFixture(t)}
	req := httptest.NewRequest(http.MethodGet, "/api/masjidboard/catalogue?name=JAMIA", nil)
	res := httptest.NewRecorder()

	s.masjidBoardCatalogue(res, req)

	body := decodeMasjidBoardCatalogueResponse(t, res)
	if body.Count != 1 || body.Records[0].ExternalID != "brits-jamia" {
		t.Fatalf("records = %+v", body.Records)
	}
}

func TestMasjidBoardCatalogueMissingFileReturnsEmptyCatalogue(t *testing.T) {
	s := &Server{masjidBoardCataloguePath: filepath.Join(t.TempDir(), "missing.json")}
	req := httptest.NewRequest(http.MethodGet, "/api/masjidboard/catalogue", nil)
	res := httptest.NewRecorder()

	s.masjidBoardCatalogue(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}
	body := decodeMasjidBoardCatalogueResponse(t, res)
	if body.Count != 0 || len(body.Records) != 0 {
		t.Fatalf("expected empty catalogue, got %+v", body)
	}
	if body.RetrievedAt != nil || body.ValidatedAt != nil {
		t.Fatal("missing catalogue should not report timestamps")
	}
}

func TestMasjidBoardCatalogueRejectsInvalidStatus(t *testing.T) {
	s := &Server{masjidBoardCataloguePath: writeMasjidBoardCatalogueFixture(t)}
	req := httptest.NewRequest(http.MethodGet, "/api/masjidboard/catalogue?status=unknown", nil)
	res := httptest.NewRecorder()

	s.masjidBoardCatalogue(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.Code)
	}
}

func TestMasjidBoardCatalogueRejectsWrongMethod(t *testing.T) {
	s := &Server{masjidBoardCataloguePath: writeMasjidBoardCatalogueFixture(t)}
	req := httptest.NewRequest(http.MethodPost, "/api/masjidboard/catalogue", nil)
	res := httptest.NewRecorder()

	s.masjidBoardCatalogue(res, req)

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", res.Code)
	}
}

func TestMasjidBoardCatalogueRequiresConfiguredPath(t *testing.T) {
	s := &Server{}
	req := httptest.NewRequest(http.MethodGet, "/api/masjidboard/catalogue", nil)
	res := httptest.NewRecorder()

	s.masjidBoardCatalogue(res, req)

	if res.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", res.Code)
	}
}
