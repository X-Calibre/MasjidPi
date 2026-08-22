package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	masjidboardruntime "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

type fakeLayoutService struct{ state selection.State }

func (f *fakeLayoutService) Configured() bool                     { return f.state.Configured() }
func (f *fakeLayoutService) Selection() selection.State           { return f.state }
func (f *fakeLayoutService) Results() []masjidboardruntime.Result { return nil }
func (f *fakeLayoutService) SetLayout(layout string) error        { f.state.Layout = layout; return nil }
func (f *fakeLayoutService) SetTheme(theme string) error          { f.state.Theme = theme; return nil }

func testLayoutState() selection.State {
	return selection.State{Boards: []selection.Board{{CatalogueID: "masjidboardlive:test", Provider: "masjidboardlive", ExternalID: "test", Name: "Test Masjid"}}}
}

func TestMasjidBoardLayoutDefaultsToStandardEmerald(t *testing.T) {
	service := &fakeLayoutService{state: testLayoutState()}
	server := &Server{masjidBoardService: service}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/masjidboard/layout", nil)
	server.masjidBoardLayout(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"layout":"standard","theme":"emerald"}` {
		t.Fatalf("body=%s", got)
	}
}

func TestMasjidBoardLayoutPUTUpdatesPreferences(t *testing.T) {
	service := &fakeLayoutService{state: testLayoutState()}
	server := &Server{masjidBoardService: service}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/masjidboard/layout", strings.NewReader(`{"layout":"detailed","theme":"ruby"}`))
	req.Header.Set("Content-Type", "application/json")
	server.masjidBoardLayout(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if service.state.Layout != selection.LayoutDetailed || service.state.Theme != selection.ThemeRuby {
		t.Fatalf("state=%+v", service.state)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"layout":"detailed","theme":"ruby"}` {
		t.Fatalf("body=%s", got)
	}
}

func TestMasjidBoardLayoutPUTRejectsInvalidLayout(t *testing.T) {
	server := &Server{masjidBoardService: &fakeLayoutService{state: testLayoutState()}}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/masjidboard/layout", strings.NewReader(`{"layout":"wide"}`))
	server.masjidBoardLayout(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestMasjidBoardLayoutPUTRejectsInvalidTheme(t *testing.T) {
	server := &Server{masjidBoardService: &fakeLayoutService{state: testLayoutState()}}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/masjidboard/layout", strings.NewReader(`{"theme":"neon"}`))
	server.masjidBoardLayout(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}
