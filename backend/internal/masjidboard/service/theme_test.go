package service

import (
	"path/filepath"
	"testing"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

func TestSetThemePersistsWithoutChangingBoards(t *testing.T) {
	path := filepath.Join(t.TempDir(), "selection.json")
	store := selection.NewStore(path)
	state := selection.State{
		Boards: []selection.Board{{CatalogueID: "masjidboardlive:test", Provider: "masjidboardlive", ExternalID: "test", Name: "Test Masjid"}},
	}
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}
	service := &Service{selection: state, selectionStore: store}
	if err := service.SetTheme(selection.ThemeRuby); err != nil {
		t.Fatalf("SetTheme() error=%v", err)
	}
	got := service.Selection()
	if got.EffectiveTheme() != selection.ThemeRuby {
		t.Fatalf("state=%+v", got)
	}
	if len(got.Boards) != 1 || got.Boards[0].ExternalID != "test" {
		t.Fatalf("boards=%+v", got.Boards)
	}
	persisted, err := selection.NewStore(path).Load()
	if err != nil {
		t.Fatal(err)
	}
	if persisted.EffectiveTheme() != selection.ThemeRuby {
		t.Fatalf("persisted=%+v", persisted)
	}
}

func TestSetThemeRejectsUnsupportedValue(t *testing.T) {
	service := &Service{selection: selection.State{Boards: []selection.Board{{CatalogueID: "masjidboardlive:test", Provider: "masjidboardlive", ExternalID: "test", Name: "Test Masjid"}}}}
	if err := service.SetTheme("neon"); err == nil {
		t.Fatal("SetTheme() expected unsupported-theme error")
	}
}
