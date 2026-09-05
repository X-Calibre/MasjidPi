package api

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/components"
	masjidboardcatalogue "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/catalogue"
	masjidboardruntime "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

type fakeSelectionRuntime struct {
	state selection.State
}

func (f *fakeSelectionRuntime) Configured() bool                     { return f.state.Configured() }
func (f *fakeSelectionRuntime) Selection() selection.State           { return f.state }
func (f *fakeSelectionRuntime) Results() []masjidboardruntime.Result { return nil }
func (f *fakeSelectionRuntime) Reconfigure(state selection.State) error {
	f.state = state
	return nil
}
func (f *fakeSelectionRuntime) Refresh(context.Context) []masjidboardruntime.Result { return nil }

func TestMasjidBoardSelectionPUTResolvesCatalogueIDs(t *testing.T) {
	path := t.TempDir() + "/catalogue.json"
	now := time.Date(2026, 8, 19, 18, 0, 0, 0, time.UTC)
	store := masjidboardcatalogue.NewStore(path)
	if err := store.Save(masjidboardcatalogue.State{Partitions: []masjidboardcatalogue.Partition{{
		Location:    masjidboardcatalogue.Location{Country: "South Africa", Region: "North West", City: "Brits"},
		RetrievedAt: now,
		ValidatedAt: now,
		Records: []masjidboardcatalogue.Record{{
			ID: "masjidboardlive:brits-jamia", Provider: "masjidboardlive", ExternalID: "brits-jamia",
			Name: "Brits Jamia Masjid", City: "Brits", Region: "North West", Country: "South Africa",
			TimeZoneOffsetMS: 7200000, DiscoveredAt: now, LastSeenAt: now, Status: masjidboardcatalogue.StatusActive,
		}},
	}}}); err != nil {
		t.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	root := t.TempDir()
	server := New(Config{Address: ":0", Frontend: root, PreferencesPath: root + "/preferences.json", Installed: components.Installed{Listen: true, Board: true}}, Dependencies{Logger: logger})
	server.SetMasjidBoardCataloguePath(path)
	runtime := &fakeSelectionRuntime{state: selection.State{
		Theme: selection.ThemeRuby, SlideDurationSeconds: 30, ShowEconomicIndicators: true,
		ShowDailyAyah: selectionBoolPointer(false), ShowDailyHadith: selectionBoolPointer(true), ShowDailySunnah: selectionBoolPointer(false),
		ShowDuaAfterAdhan: selectionBoolPointer(true),
	}}
	server.SetMasjidBoardService(runtime)

	req := httptest.NewRequest(http.MethodPut, "/api/masjidboard/selection", strings.NewReader(`{"catalogue_ids":["masjidboardlive:brits-jamia"]}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if len(runtime.state.Boards) != 1 || runtime.state.Boards[0].ExternalID != "brits-jamia" || runtime.state.Boards[0].TimeZoneOffsetMS != 7200000 {
		t.Fatalf("selection=%+v", runtime.state)
	}
	if runtime.state.Theme != selection.ThemeRuby || runtime.state.SlideDurationSeconds != 30 || !runtime.state.ShowEconomicIndicators {
		t.Fatalf("display preferences were not preserved: %+v", runtime.state)
	}
	if runtime.state.ShowDailyAyahValue() || !runtime.state.ShowDailyHadithValue() || runtime.state.ShowDailySunnahValue() {
		t.Fatalf("daily content preferences were not preserved: %+v", runtime.state)
	}
	if !runtime.state.ShowDuaAfterAdhanValue() {
		t.Fatalf("Dua after Adhan preference was not preserved: %+v", runtime.state)
	}
	if !runtime.state.Boards[0].ShowDetailedJumuahValue() {
		t.Fatalf("new board detailed Jumuah preference must default to enabled: %+v", runtime.state.Boards[0])
	}
}

func TestMasjidBoardSelectionPUTPersistsDetailedJumuahPreference(t *testing.T) {
	path := t.TempDir() + "/catalogue.json"
	now := time.Date(2026, 9, 4, 8, 0, 0, 0, time.UTC)
	record := masjidboardcatalogue.Record{
		ID: "masjidboardlive:test", Provider: "masjidboardlive", ExternalID: "test", Name: "Test Masjid",
		City: "Pretoria", Region: "Gauteng", Country: "South Africa", TimeZoneOffsetMS: 7200000,
		DiscoveredAt: now, LastSeenAt: now, Status: masjidboardcatalogue.StatusActive,
	}
	if err := masjidboardcatalogue.NewStore(path).Save(masjidboardcatalogue.State{Partitions: []masjidboardcatalogue.Partition{{
		Location:    masjidboardcatalogue.Location{Country: "South Africa", Region: "Gauteng", City: "Pretoria"},
		RetrievedAt: now, ValidatedAt: now, Records: []masjidboardcatalogue.Record{record},
	}}}); err != nil {
		t.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	root := t.TempDir()
	server := New(Config{Address: ":0", Frontend: root, PreferencesPath: root + "/preferences.json", Installed: components.Installed{Listen: true, Board: true}}, Dependencies{Logger: logger})
	server.SetMasjidBoardCataloguePath(path)
	runtime := &fakeSelectionRuntime{}
	server.SetMasjidBoardService(runtime)

	req := httptest.NewRequest(http.MethodPut, "/api/masjidboard/selection", strings.NewReader(`{"catalogue_ids":["masjidboardlive:test"],"detailed_jumuah":{"masjidboardlive:test":false}}`))
	rec := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if len(runtime.state.Boards) != 1 || runtime.state.Boards[0].ShowDetailedJumuahValue() {
		t.Fatalf("detailed Jumuah preference was not persisted: %+v", runtime.state.Boards)
	}
}

func TestMasjidBoardSelectionPUTRejectsUnknownBoard(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	root := t.TempDir()
	server := New(Config{Address: ":0", Frontend: root, PreferencesPath: root + "/preferences.json", Installed: components.Installed{Listen: true, Board: true}}, Dependencies{Logger: logger})
	server.SetMasjidBoardCataloguePath(t.TempDir() + "/catalogue.json")
	runtime := &fakeSelectionRuntime{}
	server.SetMasjidBoardService(runtime)

	req := httptest.NewRequest(http.MethodPut, "/api/masjidboard/selection", strings.NewReader(`{"catalogue_ids":["masjidboardlive:missing"]}`))
	rec := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}
