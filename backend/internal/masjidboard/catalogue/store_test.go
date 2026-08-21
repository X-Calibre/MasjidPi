package catalogue

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func testPartition(country, region, city, externalID, name string, when time.Time) Partition {
	id, _ := ID("masjidboardlive", externalID)
	return Partition{
		Location:    Location{Country: country, Region: region, City: city},
		RetrievedAt: when,
		ValidatedAt: when,
		Records: []Record{{
			ID: id, Provider: "masjidboardlive", ExternalID: externalID, Name: name,
			Country: country, Region: region, City: city,
			DiscoveredAt: when, LastSeenAt: when, Status: StatusActive,
		}},
	}
}

func testState() State {
	now := time.Date(2026, 8, 18, 18, 30, 0, 0, time.UTC)
	return State{Partitions: []Partition{
		testPartition("South Africa", "North West", "Brits", "brits-jamia", "Brits Jamia Masjid", now),
		testPartition("South Africa", "Gauteng", "Akasia", "akasia-masjid", "Akasia Masjid", now.Add(time.Hour)),
	}}
}

func TestStoreLoadMissingFile(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "catalogue.json"))
	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(got.Partitions) != 0 {
		t.Fatalf("Load() = %+v, want empty state", got)
	}
}

func TestStoreSaveAndLoadPartitions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "catalogue.json")
	store := NewStore(path)
	want := testState()
	if err := store.Save(want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("mode = %o, want 600", info.Mode().Perm())
	}
	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Partitions) != 2 {
		t.Fatalf("partitions = %d, want 2", len(got.Partitions))
	}
	merged := Merge(got)
	if len(merged.Records) != 2 {
		t.Fatalf("merged records = %d, want 2", len(merged.Records))
	}
}

func TestStoreSavePartitionPreservesOtherLocations(t *testing.T) {
	path := filepath.Join(t.TempDir(), "catalogue.json")
	store := NewStore(path)
	if err := store.Save(testState()); err != nil {
		t.Fatal(err)
	}

	now := time.Date(2026, 8, 19, 10, 0, 0, 0, time.UTC)
	updated := testPartition("South Africa", "North West", "Brits", "brits-taqwa", "Masjid Taqwa", now)
	if err := store.SavePartition(updated); err != nil {
		t.Fatal(err)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Partitions) != 2 {
		t.Fatalf("partitions = %d, want 2", len(got.Partitions))
	}
	merged := Merge(got)
	if len(merged.Records) != 2 {
		t.Fatalf("merged records = %+v", merged.Records)
	}
	foundAkasia := false
	foundTaqwa := false
	for _, record := range merged.Records {
		foundAkasia = foundAkasia || record.ExternalID == "akasia-masjid"
		foundTaqwa = foundTaqwa || record.ExternalID == "brits-taqwa"
	}
	if !foundAkasia || !foundTaqwa {
		t.Fatalf("SavePartition lost data: %+v", merged.Records)
	}
}

func TestMergeDeduplicatesBoardAcrossLocations(t *testing.T) {
	now := time.Date(2026, 8, 19, 10, 0, 0, 0, time.UTC)
	first := testPartition("South Africa", "North West", "Town A", "border-masjid", "Border Masjid", now)
	second := testPartition("South Africa", "Gauteng", "Town B", "border-masjid", "Border Masjid", now.Add(time.Hour))
	first.Records[0].Status = StatusMissing
	second.Records[0].Status = StatusActive

	merged := Merge(State{Partitions: []Partition{first, second}})
	if len(merged.Records) != 1 {
		t.Fatalf("merged records = %d, want 1", len(merged.Records))
	}
	if merged.Records[0].Status != StatusActive {
		t.Fatalf("merged status = %q, want active", merged.Records[0].Status)
	}
	if !merged.ValidatedAt.Equal(now) {
		t.Fatalf("merged validated_at = %v, want oldest %v", merged.ValidatedAt, now)
	}
}

func TestStoreRejectsDuplicateLocationPartitions(t *testing.T) {
	now := time.Date(2026, 8, 19, 10, 0, 0, 0, time.UTC)
	state := State{Partitions: []Partition{
		testPartition("South Africa", "North West", "Brits", "one", "One", now),
		testPartition(" south africa ", " north west ", " BRITS ", "two", "Two", now),
	}}
	if err := NewStore(filepath.Join(t.TempDir(), "catalogue.json")).Save(state); err == nil {
		t.Fatal("expected duplicate location validation error")
	}
}

func TestStoreRejectsMalformedJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "catalogue.json")
	if err := os.WriteFile(path, []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := NewStore(path).Load(); err == nil {
		t.Fatal("Load() expected malformed JSON error")
	}
}

func TestStoreUnchangedSaveIsNoOp(t *testing.T) {
	path := filepath.Join(t.TempDir(), "catalogue.json")
	store := NewStore(path)
	state := testState()
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}
	before, _ := os.Stat(path)
	time.Sleep(20 * time.Millisecond)
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}
	after, _ := os.Stat(path)
	if !after.ModTime().Equal(before.ModTime()) {
		t.Fatalf("unchanged Save() rewrote file")
	}
}

func TestStoreLoadReadsDiskOnDemand(t *testing.T) {
	path := filepath.Join(t.TempDir(), "catalogue.json")
	store := NewStore(path)
	state := testState()
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}
	state.Partitions[0].Records[0].Name = "Changed On Disk"
	state = normalizeState(state)
	data, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	merged := Merge(got)
	found := false
	for _, record := range merged.Records {
		if record.Name == "Changed On Disk" {
			found = true
		}
	}
	if !found {
		t.Fatal("Load() did not reread disk")
	}
}

func TestStoreReturnsDefensiveCopies(t *testing.T) {
	path := filepath.Join(t.TempDir(), "catalogue.json")
	store := NewStore(path)
	if err := store.Save(testState()); err != nil {
		t.Fatal(err)
	}
	first, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	first.Partitions[0].Records[0].Name = "Mutated"
	second, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if second.Partitions[0].Records[0].Name == "Mutated" {
		t.Fatal("persisted state was mutated through returned copy")
	}
}

func TestStoreFailedPartitionSaveKeepsLastKnownGood(t *testing.T) {
	path := filepath.Join(t.TempDir(), "catalogue.json")
	store := NewStore(path)
	if err := store.Save(testState()); err != nil {
		t.Fatal(err)
	}
	bad := testState().Partitions[0]
	bad.Records[0].ID = "wrong"
	if err := store.SavePartition(bad); err == nil {
		t.Fatal("SavePartition() expected validation error")
	}
	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Partitions) != 2 {
		t.Fatalf("last-known-good changed: %+v", got)
	}
}
