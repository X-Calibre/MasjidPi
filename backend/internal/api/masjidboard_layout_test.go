package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	masjidboardruntime "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

type fakeLayoutService struct {
	state selection.State
}

func (f *fakeLayoutService) Configured() bool { return f.state.Configured() }
func (f *fakeLayoutService) Selection() selection.State { return f.state }
func (f *fakeLayoutService) Results() []masjidboardruntime.Result { return nil }
func (f *fakeLayoutService) SetLayout(layout string) error {
	f.state.Layout = layout
	return nil
}

func TestMasjidBoardLayoutDefaultsToStandard(t *testing.T) {
	service := &fakeLayoutService{state: selection.State{Boards: []selection.Board{{
		CatalogueID: "masjidboardlive:test",
		Provider: "masjidboardlive",
		ExternalID: "test",
		Name: "Test Masjid",
	}}}}
	server := &Server{masjidBoardService: service}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/masjidboard/layout", nil)

	server.masjidBoardLayout(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"layout":"standard"}` {
		t.Fatalf("body=%s", got)
	}
}

func TestMasjidBoardLayoutPUTUpdatesPreference(t *testing.T) {
	service := &fakeLayoutService{state: selection.State{Boards: []selection.Board{{
		CatalogueID: "masjidboardlive:test",
		Provider: "masjidboardlive",
		ExternalID: "test",
		Name: "Test Masjid",
	}}}}
	server := &Server{masjidBoardService: service}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/masjidboard/layout", strings.NewReader(`{"layout":"detailed"}`))
	req.Header.Set("Content-Type", "application/json")

	server.masjidBoardLayout(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if service.state.Layout != selection.LayoutDetailed {
		t.Fatalf("layout=%q", service.state.Layout)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"layout":"detailed"}` {
		t.Fatalf("body=%s", got)
	}
}

func TestMasjidBoardLayoutPUTRejectsInvalidLayout(t *testing.T) {
	service := &fakeLayoutService{state: selection.State{Boards: []selection.Board{{
		CatalogueID: "masjidboardlive:test",
		Provider: "masjidboardlive",
		ExternalID: "test",
		Name: "Test Masjid",
	}}}}
	server := &Server{masjidBoardService: service}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/masjidboard/layout", strings.NewReader(`{"layout":"wide"}`))

	server.masjidBoardLayout(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}
