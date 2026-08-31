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
