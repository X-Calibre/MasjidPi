package catalogue

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeCatalogueSource struct {
	results map[string]Catalogue
	errors  map[string]error
	calls   []string
}

func (f *fakeCatalogueSource) Fetch(_ context.Context, location Location, _ time.Time) (Catalogue, error) {
	key := location.key()
	f.calls = append(f.calls, key)
	if err := f.errors[key]; err != nil {
		return Catalogue{}, err
	}
	return f.results[key], nil
}

type fakePartitionStore struct {
	state State
	saves []Partition
}

func (s *fakePartitionStore) Load() (State, error) { return s.state, nil }
func (s *fakePartitionStore) SavePartition(partition Partition) error {
	s.saves = append(s.saves, partition)
	key := partition.Location.key()
	for i := range s.state.Partitions {
		if s.state.Partitions[i].Location.key() == key {
			s.state.Partitions[i] = partition
			return nil
		}
	}
	s.state.Partitions = append(s.state.Partitions, partition)
	return nil
}
func (s *fakePartitionStore) RetainLocations(locations []Location) error {
	allowed := map[string]struct{}{}
	for _, location := range locations {
		allowed[location.key()] = struct{}{}
	}
	kept := s.state.Partitions[:0]
	for _, partition := range s.state.Partitions {
		if _, ok := allowed[partition.Location.key()]; ok {
			kept = append(kept, partition)
		}
	}
	s.state.Partitions = kept
	return nil
}

func candidate(location Location, id, name string, when time.Time) Catalogue {
	return Catalogue{RetrievedAt: when, ValidatedAt: when, Records: []Record{{
		ID: id, Provider: "masjidboardlive", ExternalID: id[len("masjidboardlive:"):], Name: name,
		Country: location.Country, Region: location.Region, City: location.City, Status: StatusActive,
	}}}
}

func TestCatalogueUpdaterScheduledRefreshesOnlyDueLocations(t *testing.T) {
	now := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)
	brits := Location{Country: "South Africa", Region: "North West", City: "Brits"}
	fawkner := Location{Country: "Australia", Region: "Victoria", City: "Fawkner"}
	store := &fakePartitionStore{state: State{Partitions: []Partition{{Location: brits, RetrievedAt: now.Add(-24 * time.Hour), ValidatedAt: now.Add(-24 * time.Hour), Records: candidate(brits, "masjidboardlive:brits-jamia", "Brits Jamia", now.Add(-24*time.Hour)).Records}}}}
	source := &fakeCatalogueSource{results: map[string]Catalogue{fawkner.key(): candidate(fawkner, "masjidboardlive:fawkner-masjid", "Fawkner Masjid", now)}, errors: map[string]error{}}
	result, err := (Updater{Source: source, Store: store}).RefreshScheduled(context.Background(), []Location{brits, fawkner}, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(source.calls) != 1 || source.calls[0] != fawkner.key() {
		t.Fatalf("calls=%v", source.calls)
	}
	if result.Locations[0].Attempted || !result.Locations[1].Attempted || !result.Locations[1].Updated {
		t.Fatalf("result=%+v", result)
	}
}

func TestCatalogueUpdaterFailureIsolation(t *testing.T) {
	now := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)
	a := Location{Country: "South Africa", Region: "North West", City: "Brits"}
	b := Location{Country: "South Africa", Region: "Gauteng", City: "Akasia"}
	oldB := Partition{Location: b, RetrievedAt: now.Add(-8 * 24 * time.Hour), ValidatedAt: now.Add(-8 * 24 * time.Hour), Records: candidate(b, "masjidboardlive:old-b", "Old B", now.Add(-8*24*time.Hour)).Records}
	store := &fakePartitionStore{state: State{Partitions: []Partition{oldB}}}
	source := &fakeCatalogueSource{results: map[string]Catalogue{a.key(): candidate(a, "masjidboardlive:brits-jamia", "Brits Jamia", now)}, errors: map[string]error{b.key(): errors.New("offline")}}
	result, err := (Updater{Source: source, Store: store}).RefreshManual(context.Background(), []Location{a, b}, now)
	if err != nil {
		t.Fatal(err)
	}
	if !result.AnyFailed() || len(store.saves) != 1 || store.saves[0].Location.key() != a.key() {
		t.Fatalf("result=%+v saves=%+v", result, store.saves)
	}
	if store.state.Partitions[0].Location.key() == b.key() && store.state.Partitions[0].Records[0].Name != "Old B" {
		t.Fatal("failed location lost last-known-good partition")
	}
}

func TestCatalogueUpdaterReconcilesMissingPerLocation(t *testing.T) {
	now := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)
	loc := Location{Country: "South Africa", Region: "North West", City: "Brits"}
	old := candidate(loc, "masjidboardlive:brits-taqwa", "Masjid Taqwa", now.Add(-8*24*time.Hour))
	old.Records[0].DiscoveredAt = now.Add(-30 * 24 * time.Hour)
	store := &fakePartitionStore{state: State{Partitions: []Partition{{Location: loc, RetrievedAt: old.RetrievedAt, ValidatedAt: old.ValidatedAt, Records: old.Records}}}}
	source := &fakeCatalogueSource{results: map[string]Catalogue{loc.key(): {RetrievedAt: now, ValidatedAt: now}}, errors: map[string]error{}}
	_, err := (Updater{Source: source, Store: store}).RefreshManual(context.Background(), []Location{loc}, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(store.saves) != 1 || store.saves[0].Records[0].Status != StatusMissing {
		t.Fatalf("saved=%+v", store.saves)
	}
}

func TestCatalogueUpdaterPrunesRemovedLocations(t *testing.T) {
	now := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)
	keep := Location{Country: "South Africa", Region: "North West", City: "Brits"}
	remove := Location{Country: "South Africa", Region: "Gauteng", City: "Akasia"}
	store := &fakePartitionStore{state: State{Partitions: []Partition{{Location: keep, RetrievedAt: now, ValidatedAt: now}, {Location: remove, RetrievedAt: now, ValidatedAt: now}}}}
	source := &fakeCatalogueSource{results: map[string]Catalogue{}, errors: map[string]error{}}
	if _, err := (Updater{Source: source, Store: store}).RefreshScheduled(context.Background(), []Location{keep}, now); err != nil {
		t.Fatal(err)
	}
	if len(store.state.Partitions) != 1 || store.state.Partitions[0].Location.key() != keep.key() {
		t.Fatalf("partitions=%+v", store.state.Partitions)
	}
}

func TestCatalogueUpdaterRejectsInvalidLocationSet(t *testing.T) {
	updater := Updater{Source: &fakeCatalogueSource{}, Store: &fakePartitionStore{}}
	if _, err := updater.RefreshManual(context.Background(), nil, time.Now()); err == nil {
		t.Fatal("expected empty location error")
	}
	loc := Location{Country: "South Africa", City: "Brits"}
	if _, err := updater.RefreshManual(context.Background(), []Location{loc, loc}, time.Now()); err == nil {
		t.Fatal("expected duplicate location error")
	}
}
