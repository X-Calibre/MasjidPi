package economic

import (
	"path/filepath"
	"testing"
	"time"
)

func TestStoreRoundTrip(t *testing.T) {
	t.Parallel()
	store := Store{Path: filepath.Join(t.TempDir(), "nested", "indicators.json")}
	want := Indicators{Source: SourceName, SourceURL: "https://example.test/source", EffectiveDate: "2026-08-24", Nisaab: 21708.16, Krugerrand: 77626.36, FetchedAt: time.Now().UTC().Truncate(time.Second)}
	if err := store.Save(want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got == nil || *got != want {
		t.Fatalf("Load() = %+v, want %+v", got, want)
	}
}

func TestStoreMissingIsEmpty(t *testing.T) {
	t.Parallel()
	got, err := (Store{Path: filepath.Join(t.TempDir(), "missing.json")}).Load()
	if err != nil || got != nil {
		t.Fatalf("Load() = %+v, %v; want nil, nil", got, err)
	}
}
