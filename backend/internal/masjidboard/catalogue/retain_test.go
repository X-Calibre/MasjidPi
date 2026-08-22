package catalogue

import (
	"path/filepath"
	"testing"
	"time"
)

func TestStoreRetainLocationsRemovesObsoletePartitions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "catalogue.json")
	store := NewStore(path)
	now := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)
	brits := Location{Country: "South Africa", Region: "North West", City: "Brits"}
	akasia := Location{Country: "South Africa", Region: "Gauteng", City: "Akasia"}
	if err := store.Save(State{Partitions: []Partition{{Location: brits, RetrievedAt: now, ValidatedAt: now}, {Location: akasia, RetrievedAt: now, ValidatedAt: now}}}); err != nil {
		t.Fatal(err)
	}
	if err := store.RetainLocations([]Location{brits}); err != nil {
		t.Fatal(err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Partitions) != 1 || got.Partitions[0].Location.key() != brits.key() {
		t.Fatalf("partitions=%+v", got.Partitions)
	}
}
