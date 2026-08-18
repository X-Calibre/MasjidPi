package catalogue

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func testCatalogue() Catalogue {
	now := time.Date(2026, 8, 18, 18, 30, 0, 0, time.UTC)
	return Catalogue{
		RetrievedAt: now,
		ValidatedAt: now,
		Records: []Record{
			{
				ID:               "masjidboardlive:brits-jamia",
				Provider:         "masjidboardlive",
				ExternalID:       "brits-jamia",
				Name:             "Brits Jamia Masjid",
				City:             "Brits",
				Region:           "North West",
				Country:          "South Africa",
				TimeZoneOffsetMS: 7200000,
				ProviderMetadata: map[string]string{"mbl_id": "MBL11517PRP"},
				DiscoveredAt:     now,
				LastSeenAt:       now,
				Status:           StatusActive,
			},
		},
	}
}

func TestStoreLoadMissingFile(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "catalogue.json"))
	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !got.RetrievedAt.IsZero() || !got.ValidatedAt.IsZero() || len(got.Records) != 0 {
		t.Fatalf("Load() = %+v, want empty catalogue", got)
	}
}

func TestStoreSaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "catalogue.json")
	store := NewStore(path)
	want := testCatalogue()

	if err := store.Save(want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("mode = %o, want 600", info.Mode().Perm())
	}

	freshStore := NewStore(path)
	got, err := freshStore.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !cataloguesEqual(got, want) {
		t.Fatalf("Load() = %+v, want %+v", got, want)
	}
}

func TestStoreRejectsMalformedJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "catalogue.json")
	if err := os.WriteFile(path, []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}
	store := NewStore(path)
	if _, err := store.Load(); err == nil {
		t.Fatal("Load() expected malformed JSON error")
	}
}

func TestStoreRejectsInvalidPersistedCatalogue(t *testing.T) {
	path := filepath.Join(t.TempDir(), "catalogue.json")
	if err := os.WriteFile(path, []byte(`{"retrieved_at":"2026-08-18T18:30:00Z","validated_at":"2026-08-18T18:30:00Z","records":[{"id":"wrong","provider":"masjidboardlive","external_id":"brits-jamia","name":"Brits Jamia Masjid"}]}`), 0600); err != nil {
		t.Fatal(err)
	}
	store := NewStore(path)
	if _, err := store.Load(); err == nil {
		t.Fatal("Load() expected validation error")
	}
}

func TestStoreUnchangedSaveIsNoOp(t *testing.T) {
	path := filepath.Join(t.TempDir(), "catalogue.json")
	store := NewStore(path)
	state := testCatalogue()
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}

	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(20 * time.Millisecond)
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}
	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !after.ModTime().Equal(before.ModTime()) {
		t.Fatalf("unchanged Save() rewrote file: before=%v after=%v", before.ModTime(), after.ModTime())
	}
}

func TestStoreChangedSaveReplacesState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "catalogue.json")
	store := NewStore(path)
	first := testCatalogue()
	if err := store.Save(first); err != nil {
		t.Fatal(err)
	}

	second := testCatalogue()
	second.Records[0].Name = "Brits Jamia"
	second.ValidatedAt = second.ValidatedAt.Add(time.Hour)
	second.RetrievedAt = second.RetrievedAt.Add(time.Hour)
	if err := store.Save(second); err != nil {
		t.Fatal(err)
	}

	got, err := NewStore(path).Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.Records[0].Name != "Brits Jamia" || !got.ValidatedAt.Equal(second.ValidatedAt) {
		t.Fatalf("changed state not persisted: %+v", got)
	}
}

func TestStoreLoadCachesState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "catalogue.json")
	state := testCatalogue()
	store := NewStore(path)
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path, []byte("broken"), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatalf("cached Load() error = %v", err)
	}
	if got.Records[0].Name != state.Records[0].Name {
		t.Fatalf("cached Load() = %+v", got)
	}
}

func TestStoreReturnsDefensiveCopies(t *testing.T) {
	path := filepath.Join(t.TempDir(), "catalogue.json")
	store := NewStore(path)
	state := testCatalogue()
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}

	first, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	first.Records[0].Name = "Mutated"
	first.Records[0].ProviderMetadata["mbl_id"] = "changed"

	second, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if second.Records[0].Name != "Brits Jamia Masjid" || second.Records[0].ProviderMetadata["mbl_id"] != "MBL11517PRP" {
		t.Fatalf("store state was mutated through returned copy: %+v", second.Records[0])
	}
}

func TestStoreFailedSaveKeepsLastKnownGood(t *testing.T) {
	path := filepath.Join(t.TempDir(), "catalogue.json")
	store := NewStore(path)
	state := testCatalogue()
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}

	bad := testCatalogue()
	bad.Records[0].ID = "wrong"
	if err := store.Save(bad); err == nil {
		t.Fatal("Save() expected validation error")
	}

	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.Records[0].ID != state.Records[0].ID {
		t.Fatalf("last-known-good state changed after failed Save(): %+v", got)
	}
}

func cataloguesEqual(a, b Catalogue) bool {
	if !a.RetrievedAt.Equal(b.RetrievedAt) || !a.ValidatedAt.Equal(b.ValidatedAt) || len(a.Records) != len(b.Records) {
		return false
	}
	for i := range a.Records {
		x, y := a.Records[i], b.Records[i]
		if x.ID != y.ID || x.Provider != y.Provider || x.ExternalID != y.ExternalID || x.Name != y.Name || x.City != y.City || x.Region != y.Region || x.Country != y.Country || x.TimeZoneOffsetMS != y.TimeZoneOffsetMS || x.Status != y.Status || !x.DiscoveredAt.Equal(y.DiscoveredAt) || !x.LastSeenAt.Equal(y.LastSeenAt) {
			return false
		}
		if len(x.ProviderMetadata) != len(y.ProviderMetadata) {
			return false
		}
		for key, value := range x.ProviderMetadata {
			if y.ProviderMetadata[key] != value {
				return false
			}
		}
	}
	return true
}
