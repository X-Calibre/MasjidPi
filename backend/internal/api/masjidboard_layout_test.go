package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	masjidboardruntime "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

type fakeDisplaySettingsService struct{ state selection.State }

func (f *fakeDisplaySettingsService) Configured() bool                     { return f.state.Configured() }
func (f *fakeDisplaySettingsService) Selection() selection.State           { return f.state }
func (f *fakeDisplaySettingsService) Results() []masjidboardruntime.Result { return nil }
func (f *fakeDisplaySettingsService) SetTheme(theme string) error          { f.state.Theme = theme; return nil }
func (f *fakeDisplaySettingsService) SetSlideDurationSeconds(seconds int) error {
	f.state.SlideDurationSeconds = seconds
	return nil
}
func (f *fakeDisplaySettingsService) SetShowEconomicIndicators(show bool) error {
	f.state.ShowEconomicIndicators = show
	return nil
}
func (f *fakeDisplaySettingsService) SetDailyIslamicContentPreferences(ayah, hadith, sunnah bool) error {
	f.state.ShowDailyAyah = selectionBoolPointer(ayah)
	f.state.ShowDailyHadith = selectionBoolPointer(hadith)
	f.state.ShowDailySunnah = selectionBoolPointer(sunnah)
	return nil
}

func testDisplaySettingsState() selection.State {
	return selection.State{Boards: []selection.Board{{CatalogueID: "masjidboardlive:test", Provider: "masjidboardlive", ExternalID: "test", Name: "Test Masjid"}}}
}

func TestMasjidBoardDisplaySettingsDefaults(t *testing.T) {
	service := &fakeDisplaySettingsService{state: testDisplaySettingsState()}
	server := &Server{masjidBoardService: service}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/masjidboard/layout", nil)
	server.masjidBoardLayout(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"theme":"emerald","slide_duration_seconds":15,"show_economic_indicators":false,"show_daily_ayah":true,"show_daily_hadith":true,"show_daily_sunnah":true}` {
		t.Fatalf("body=%s", got)
	}
}

func TestMasjidBoardDisplaySettingsPUTUpdatesPreferences(t *testing.T) {
	service := &fakeDisplaySettingsService{state: testDisplaySettingsState()}
	server := &Server{masjidBoardService: service}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/masjidboard/layout", strings.NewReader(`{"theme":"ruby","slide_duration_seconds":30,"show_economic_indicators":true,"show_daily_ayah":false,"show_daily_hadith":true,"show_daily_sunnah":false}`))
	req.Header.Set("Content-Type", "application/json")
	server.masjidBoardLayout(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if service.state.Theme != selection.ThemeRuby || service.state.SlideDurationSeconds != 30 || !service.state.ShowEconomicIndicators {
		t.Fatalf("state=%+v", service.state)
	}
	if service.state.ShowDailyAyahValue() || !service.state.ShowDailyHadithValue() || service.state.ShowDailySunnahValue() {
		t.Fatalf("daily content state=%+v", service.state)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `{"theme":"ruby","slide_duration_seconds":30,"show_economic_indicators":true,"show_daily_ayah":false,"show_daily_hadith":true,"show_daily_sunnah":false}` {
		t.Fatalf("body=%s", got)
	}
}

func TestMasjidBoardDisplaySettingsPUTPreservesUnspecifiedDailyPreferences(t *testing.T) {
	falseValue := false
	service := &fakeDisplaySettingsService{state: testDisplaySettingsState()}
	service.state.ShowDailyAyah = &falseValue
	server := &Server{masjidBoardService: service}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/masjidboard/layout", strings.NewReader(`{"show_daily_hadith":false}`))
	server.masjidBoardLayout(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if service.state.ShowDailyAyahValue() || service.state.ShowDailyHadithValue() || !service.state.ShowDailySunnahValue() {
		t.Fatalf("state=%+v", service.state)
	}
}

func TestMasjidBoardDisplaySettingsPUTRejectsInvalidSlideDuration(t *testing.T) {
	server := &Server{masjidBoardService: &fakeDisplaySettingsService{state: testDisplaySettingsState()}}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/masjidboard/layout", strings.NewReader(`{"slide_duration_seconds":61}`))
	server.masjidBoardLayout(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestMasjidBoardDisplaySettingsPUTRejectsRemovedLayoutField(t *testing.T) {
	server := &Server{masjidBoardService: &fakeDisplaySettingsService{state: testDisplaySettingsState()}}
	for _, layout := range []string{"portrait", "landscape", "standard", "detailed"} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/api/masjidboard/layout", strings.NewReader(`{"layout":"`+layout+`"}`))
		server.masjidBoardLayout(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("layout=%q status=%d body=%s", layout, rec.Code, rec.Body.String())
		}
	}
}

func TestMasjidBoardDisplaySettingsPUTRejectsInvalidTheme(t *testing.T) {
	server := &Server{masjidBoardService: &fakeDisplaySettingsService{state: testDisplaySettingsState()}}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/masjidboard/layout", strings.NewReader(`{"theme":"neon"}`))
	server.masjidBoardLayout(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}
