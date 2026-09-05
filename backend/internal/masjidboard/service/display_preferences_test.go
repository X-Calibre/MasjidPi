package service

import (
	"path/filepath"
	"testing"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

func TestSetSlideDurationPersists(t *testing.T) {
	path := filepath.Join(t.TempDir(), "selection.json")
	store := selection.NewStore(path)
	state := selection.State{Boards: []selection.Board{{
		CatalogueID: "masjidboardlive:test",
		Provider:    "masjidboardlive",
		ExternalID:  "test",
		Name:        "Test Masjid",
	}}}
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}

	service := &Service{selection: state, selectionStore: store}
	if err := service.SetSlideDurationSeconds(30); err != nil {
		t.Fatalf("SetSlideDurationSeconds() error = %v", err)
	}
	if got := service.Selection().EffectiveSlideDurationSeconds(); got != 30 {
		t.Fatalf("runtime slide duration=%d", got)
	}
	persisted, err := selection.NewStore(path).Load()
	if err != nil {
		t.Fatal(err)
	}
	if got := persisted.EffectiveSlideDurationSeconds(); got != 30 {
		t.Fatalf("persisted slide duration=%d", got)
	}
	if err := service.SetSlideDurationSeconds(61); err == nil {
		t.Fatal("SetSlideDurationSeconds() expected range error")
	}
}

func TestSetDailyIslamicContentPreferencesPersistsExplicitValues(t *testing.T) {
	path := filepath.Join(t.TempDir(), "selection.json")
	store := selection.NewStore(path)
	state := selection.State{Boards: []selection.Board{{
		CatalogueID: "masjidboardlive:test", Provider: "masjidboardlive", ExternalID: "test", Name: "Test Masjid",
	}}}
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}
	service := &Service{selection: state, selectionStore: store}
	if err := service.SetDailyIslamicContentPreferences(false, true, false); err != nil {
		t.Fatalf("SetDailyIslamicContentPreferences() error = %v", err)
	}
	persisted, err := selection.NewStore(path).Load()
	if err != nil {
		t.Fatal(err)
	}
	if persisted.ShowDailyAyah == nil || persisted.ShowDailyHadith == nil || persisted.ShowDailySunnah == nil {
		t.Fatalf("preferences were not persisted explicitly: %+v", persisted)
	}
	if persisted.ShowDailyAyahValue() || !persisted.ShowDailyHadithValue() || persisted.ShowDailySunnahValue() {
		t.Fatalf("persisted preferences = %+v", persisted)
	}
}

func TestSetShowDuaAfterAdhanPersists(t *testing.T) {
	path := filepath.Join(t.TempDir(), "selection.json")
	store := selection.NewStore(path)
	state := selection.State{Boards: []selection.Board{{CatalogueID: "masjidboardlive:test", Provider: "masjidboardlive", ExternalID: "test", Name: "Test Masjid"}}}
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}
	service := &Service{selection: state, selectionStore: store}
	if err := service.SetShowDuaAfterAdhan(true); err != nil {
		t.Fatalf("SetShowDuaAfterAdhan() error = %v", err)
	}
	persisted, err := selection.NewStore(path).Load()
	if err != nil {
		t.Fatal(err)
	}
	if !persisted.ShowDuaAfterAdhanValue() {
		t.Fatal("Dua after Adhan preference was not persisted")
	}
}
