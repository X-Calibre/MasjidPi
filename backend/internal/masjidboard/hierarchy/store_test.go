package hierarchy

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func testHierarchyState() State {
	now := time.Date(2026, 8, 18, 20, 0, 0, 0, time.UTC)
	return State{
		RetrievedAt: now,
		ValidatedAt: now,
		Countries: []Country{{
			Name:  "South Africa",
			Count: 615,
			Regions: []Region{{
				Name:  "North West",
				Count: 23,
				Cities: []Location{{Name: "Brits", Count: 3}},
			}},
		}},
	}
}

func TestStoreMissingReturnsZeroState(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "hierarchy.json"))
	state, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(state.Countries) != 0 || !state.RetrievedAt.IsZero() {
		t.Fatalf("state = %+v", state)
	}
}

func TestStoreSaveLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hierarchy.json")
	store := NewStore(path)
	want := testHierarchyState()
	if err := store.Save(want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(got.Countries) != 1 || got.Countries[0].Name != "South Africa" || got.Countries[0].Regions[0].Cities[0].Name != "Brits" {
		t.Fatalf("got = %+v", got)
	}
}

func TestStoreRejectsInvalidCandidateWithoutReplacingGoodState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hierarchy.json")
	store := NewStore(path)
	good := testHierarchyState()
	if err := store.Save(good); err != nil {
		t.Fatalf("initial Save() error = %v", err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	bad := good
	bad.Countries = nil
	if err := store.Save(bad); err == nil {
		t.Fatal("Save() expected validation error")
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(before) != string(after) {
		t.Fatal("invalid candidate replaced last-known-good hierarchy")
	}
}
