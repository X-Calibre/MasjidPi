package service

import (
	"path/filepath"
	"testing"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

func TestSetLayoutPersistsWithoutChangingBoards(t *testing.T) {
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
	if err := service.SetLayout(selection.LayoutDetailed); err != nil {
		t.Fatalf("SetLayout() error = %v", err)
	}

	got := service.Selection()
	if got.EffectiveLayout() != selection.LayoutDetailed {
		t.Fatalf("runtime layout=%q", got.EffectiveLayout())
	}
	if len(got.Boards) != 1 || got.Boards[0].ExternalID != "test" {
		t.Fatalf("boards changed: %+v", got.Boards)
	}

	persisted, err := selection.NewStore(path).Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if persisted.EffectiveLayout() != selection.LayoutDetailed {
		t.Fatalf("persisted layout=%q", persisted.EffectiveLayout())
	}
	if len(persisted.Boards) != 1 || persisted.Boards[0].ExternalID != "test" {
		t.Fatalf("persisted boards changed: %+v", persisted.Boards)
	}
}

func TestSetLayoutRejectsUnsupportedValue(t *testing.T) {
	service := &Service{selection: selection.State{Boards: []selection.Board{{
		CatalogueID: "masjidboardlive:test",
		Provider:    "masjidboardlive",
		ExternalID:  "test",
		Name:        "Test Masjid",
	}}}}
	if err := service.SetLayout("wide"); err == nil {
		t.Fatal("SetLayout() expected unsupported-layout error")
	}
}
