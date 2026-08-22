package selection

import "testing"

func TestEffectiveThemeDefaultsToEmerald(t *testing.T) {
	if got := (State{}).EffectiveTheme(); got != ThemeEmerald {
		t.Fatalf("EffectiveTheme()=%q, want %q", got, ThemeEmerald)
	}
}

func TestValidateAcceptsSupportedThemes(t *testing.T) {
	board := selected("test", "Test Masjid", 0)
	for _, theme := range []string{ThemeEmerald, ThemeMidnight, ThemeSlate, ThemeRuby, ThemeLight, ThemeBlackWhite} {
		if err := Validate(State{Boards: []Board{board}, Theme: theme}); err != nil {
			t.Fatalf("theme %q rejected: %v", theme, err)
		}
	}
}

func TestValidateRejectsUnsupportedTheme(t *testing.T) {
	board := selected("test", "Test Masjid", 0)
	if err := Validate(State{Boards: []Board{board}, Theme: "neon"}); err == nil {
		t.Fatal("Validate() expected unsupported-theme error")
	}
}

func TestStorePreservesTheme(t *testing.T) {
	store := NewStore(t.TempDir() + "/selection.json")
	state := State{Boards: []Board{selected("test", "Test Masjid", 0)}, Theme: ThemeRuby}
	if err := store.Save(state); err != nil { t.Fatal(err) }
	got, err := store.Load(); if err != nil { t.Fatal(err) }
	if got.EffectiveTheme() != ThemeRuby { t.Fatalf("theme=%q", got.EffectiveTheme()) }
}
