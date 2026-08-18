package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFavouritesSaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "favourites.json")
	favourites := NewFavourites(path)

	want := []string{"one", "two", "three"}
	if err := favourites.Save(want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := favourites.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("Load() returned %d IDs, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Load()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestFavouritesLoadMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	favourites := NewFavourites(path)

	got, err := favourites.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got == nil {
		t.Fatal("Load() returned nil slice")
	}
	if len(got) != 0 {
		t.Fatalf("Load() returned %d IDs, want 0", len(got))
	}
}

func TestFavouritesLoadCorruptFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "favourites.json")
	if err := os.WriteFile(path, []byte(`{"ids":[`), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	favourites := NewFavourites(path)
	if _, err := favourites.Load(); err == nil {
		t.Fatal("Load() error = nil, want JSON error")
	}
}

func TestFavouritesSaveCreatesParentDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "state", "favourites.json")
	favourites := NewFavourites(path)

	if err := favourites.Save([]string{"one"}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("saved favourites file not found: %v", err)
	}
}

func TestFavouritesLoadUsesCachedState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "favourites.json")
	if err := os.WriteFile(path, []byte(`{"ids":["one"]}`), 0600); err != nil {
		t.Fatalf("write initial state: %v", err)
	}

	favourites := NewFavourites(path)
	got, err := favourites.Load()
	if err != nil || len(got) != 1 || got[0] != "one" {
		t.Fatalf("initial Load() = %v, %v", got, err)
	}

	if err := os.WriteFile(path, []byte(`{"ids":["external"]}`), 0600); err != nil {
		t.Fatalf("replace backing file: %v", err)
	}

	got, err = favourites.Load()
	if err != nil {
		t.Fatalf("cached Load() error = %v", err)
	}
	if len(got) != 1 || got[0] != "one" {
		t.Fatalf("cached Load() = %v, want [one]", got)
	}
}

func TestFavouritesLoadReturnsCopy(t *testing.T) {
	path := filepath.Join(t.TempDir(), "favourites.json")
	favourites := NewFavourites(path)
	if err := favourites.Save([]string{"one", "two"}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := favourites.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	got[0] = "mutated"

	again, err := favourites.Load()
	if err != nil {
		t.Fatalf("second Load() error = %v", err)
	}
	if again[0] != "one" {
		t.Fatalf("cached state was mutated through returned slice: %v", again)
	}
}
