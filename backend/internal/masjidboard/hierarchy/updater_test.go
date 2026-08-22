package hierarchy

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeHierarchySource struct {
	countries []Location
	regions   map[string][]Location
	cities    map[string][]Location
	err       error
	calls     int
}

func (f *fakeHierarchySource) Countries(context.Context) ([]Location, error) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	return f.countries, nil
}

func (f *fakeHierarchySource) Regions(_ context.Context, country string) ([]Location, error) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	return f.regions[country], nil
}

func (f *fakeHierarchySource) Cities(_ context.Context, country, region string) ([]Location, error) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	return f.cities[country+"|"+region], nil
}

type fakeHierarchyStore struct {
	state State
	saves int
}

func (s *fakeHierarchyStore) Load() (State, error) { return s.state, nil }
func (s *fakeHierarchyStore) Save(state State) error {
	s.saves++
	s.state = state
	return nil
}

func TestUpdaterScheduledSkipsFreshHierarchy(t *testing.T) {
	now := time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC)
	store := &fakeHierarchyStore{state: State{
		RetrievedAt: now.Add(-24 * time.Hour),
		ValidatedAt: now.Add(-24 * time.Hour),
		Countries:   []Country{{Name: "South Africa", Count: 1, Regions: []Region{{Name: "North West", Count: 1, Cities: []Location{{Name: "Brits", Count: 1}}}}}},
	}}
	source := &fakeHierarchySource{}

	result, err := (Updater{Source: source, Store: store}).RefreshScheduled(context.Background(), now)
	if err != nil {
		t.Fatalf("RefreshScheduled() error = %v", err)
	}
	if result.Attempted {
		t.Fatal("RefreshScheduled() attempted fresh hierarchy")
	}
	if source.calls != 0 || store.saves != 0 {
		t.Fatalf("calls=%d saves=%d, want 0/0", source.calls, store.saves)
	}
}

func TestUpdaterManualBuildsCompleteHierarchy(t *testing.T) {
	now := time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC)
	source := &fakeHierarchySource{
		countries: []Location{{Name: "South Africa", Count: 3}},
		regions: map[string][]Location{
			"South Africa": {{Name: "North West", Count: 3}},
		},
		cities: map[string][]Location{
			"South Africa|North West": {{Name: "Brits", Count: 2}, {Name: "Rustenburg", Count: 1}},
		},
	}
	store := &fakeHierarchyStore{}

	result, err := (Updater{Source: source, Store: store}).RefreshManual(context.Background(), now)
	if err != nil {
		t.Fatalf("RefreshManual() error = %v", err)
	}
	if !result.Attempted || !result.Updated {
		t.Fatalf("result = %+v", result)
	}
	if store.saves != 1 {
		t.Fatalf("saves = %d, want 1", store.saves)
	}
	if len(result.State.Countries) != 1 || len(result.State.Countries[0].Regions) != 1 {
		t.Fatalf("state = %+v", result.State)
	}
	cities := result.State.Countries[0].Regions[0].Cities
	if len(cities) != 2 || cities[0].Name != "Brits" || cities[1].Name != "Rustenburg" {
		t.Fatalf("cities = %+v", cities)
	}
	if result.State.Countries[0].Regions[0].UnresolvedCount != 0 {
		t.Fatalf("unresolved = %d, want 0", result.State.Countries[0].Regions[0].UnresolvedCount)
	}
}

func TestUpdaterMergesDuplicateRegionLabelsBeforeCityFetch(t *testing.T) {
	now := time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC)
	source := &fakeHierarchySource{
		countries: []Location{{Name: "South Africa", Count: 26}},
		regions: map[string][]Location{
			"South Africa": {{Name: "Limpopo", Count: 25}, {Name: "Limpopo", Count: 1}},
		},
		cities: map[string][]Location{
			"South Africa|Limpopo": {{Name: "Polokwane", Count: 25}},
		},
	}
	store := &fakeHierarchyStore{}

	result, err := (Updater{Source: source, Store: store}).RefreshManual(context.Background(), now)
	if err != nil {
		t.Fatalf("RefreshManual() error = %v", err)
	}
	regions := result.State.Countries[0].Regions
	if len(regions) != 1 || regions[0].Name != "Limpopo" || regions[0].Count != 26 {
		t.Fatalf("regions = %+v", regions)
	}
	if regions[0].UnresolvedCount != 1 {
		t.Fatalf("unresolved = %d, want 1", regions[0].UnresolvedCount)
	}
	if source.calls != 3 {
		t.Fatalf("source calls = %d, want 3 (countries, regions, cities once)", source.calls)
	}
}

func TestUpdaterPreservesMixedBlankRegionAsUnresolved(t *testing.T) {
	now := time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC)
	source := &fakeHierarchySource{
		countries: []Location{{Name: "South Africa", Count: 3}},
		regions: map[string][]Location{
			"South Africa": {{Name: "", Count: 1}, {Name: "Gauteng", Count: 2}},
		},
		cities: map[string][]Location{
			"South Africa|Gauteng": {{Name: "Pretoria", Count: 2}},
		},
	}
	store := &fakeHierarchyStore{}

	result, err := (Updater{Source: source, Store: store}).RefreshManual(context.Background(), now)
	if err != nil {
		t.Fatalf("RefreshManual() error = %v", err)
	}
	regions := result.State.Countries[0].Regions
	if len(regions) != 2 {
		t.Fatalf("regions = %+v", regions)
	}
	if regions[0].Name != "" || regions[0].UnresolvedCount != 1 || len(regions[0].Cities) != 0 {
		t.Fatalf("blank region = %+v", regions[0])
	}
	if source.calls != 3 {
		t.Fatalf("source calls = %d, want 3 (countries, regions, one city request)", source.calls)
	}
}

func TestUpdaterRejectsCityCountExceedingRegionAndPreservesLastKnownGood(t *testing.T) {
	now := time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC)
	old := State{
		RetrievedAt: now.Add(-8 * 24 * time.Hour),
		ValidatedAt: now.Add(-8 * 24 * time.Hour),
		Countries:   []Country{{Name: "South Africa", Count: 1, Regions: []Region{{Name: "North West", Count: 1, Cities: []Location{{Name: "Brits", Count: 1}}}}}},
	}
	store := &fakeHierarchyStore{state: old}
	source := &fakeHierarchySource{
		countries: []Location{{Name: "South Africa", Count: 2}},
		regions:   map[string][]Location{"South Africa": {{Name: "North West", Count: 2}}},
		cities:    map[string][]Location{"South Africa|North West": {{Name: "Brits", Count: 3}}},
	}

	result, err := (Updater{Source: source, Store: store}).RefreshManual(context.Background(), now)
	if err == nil {
		t.Fatal("RefreshManual() expected count validation error")
	}
	if result.State.ValidatedAt != old.ValidatedAt {
		t.Fatalf("fallback state changed: %+v", result.State)
	}
	if store.saves != 0 || store.state.ValidatedAt != old.ValidatedAt {
		t.Fatalf("last-known-good was replaced: saves=%d state=%+v", store.saves, store.state)
	}
}

func TestUpdaterPreservesLastKnownGoodOnSourceFailure(t *testing.T) {
	now := time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC)
	old := State{RetrievedAt: now.Add(-8 * 24 * time.Hour), ValidatedAt: now.Add(-8 * 24 * time.Hour), Countries: []Country{{Name: "South Africa", Count: 0}}}
	store := &fakeHierarchyStore{state: old}
	source := &fakeHierarchySource{err: errors.New("offline")}

	result, err := (Updater{Source: source, Store: store}).RefreshScheduled(context.Background(), now)
	if err == nil {
		t.Fatal("RefreshScheduled() expected source error")
	}
	if !result.Attempted || result.State.ValidatedAt != old.ValidatedAt {
		t.Fatalf("result = %+v", result)
	}
	if store.saves != 0 {
		t.Fatalf("saves = %d, want 0", store.saves)
	}
}
