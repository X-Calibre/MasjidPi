package catalogue

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/scope"
)

type fakeScopeStore struct {
	state scope.State
	err   error
}

func (s fakeScopeStore) Load() (scope.State, error) {
	if s.err != nil {
		return scope.State{}, s.err
	}
	return s.state, nil
}

func TestScopedUpdaterScheduledUsesAllConfiguredLocations(t *testing.T) {
	now := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)
	store := &fakePartitionStore{}
	source := &fakeCatalogueSource{records: map[string][]Record{
		"South Africa|North West|Brits": {{ID: "masjidboardlive:brits", Provider: "masjidboardlive", ExternalID: "brits", Name: "Brits Masjid"}},
		"South Africa|Gauteng|Akasia":   {{ID: "masjidboardlive:akasia", Provider: "masjidboardlive", ExternalID: "akasia", Name: "Akasia Masjid"}},
	}}

	u := ScopedUpdater{
		Scope: fakeScopeStore{state: scope.State{Locations: []scope.Location{
			{Country: "South Africa", Region: "North West", City: "Brits"},
			{Country: "South Africa", Region: "Gauteng", City: "Akasia"},
		}}},
		Updater: Updater{Source: source, Store: store},
	}

	result, err := u.RefreshScheduled(context.Background(), now)
	if err != nil {
		t.Fatalf("RefreshScheduled() error = %v", err)
	}
	if len(result.Locations) != 2 {
		t.Fatalf("locations = %d, want 2", len(result.Locations))
	}
	for i, item := range result.Locations {
		if !item.Attempted || item.Error != nil {
			t.Fatalf("location %d result = %+v", i, item)
		}
	}
}

func TestScopedUpdaterScheduledUnconfiguredIsNoOp(t *testing.T) {
	u := ScopedUpdater{Scope: fakeScopeStore{}, Updater: Updater{}}
	result, err := u.RefreshScheduled(context.Background(), time.Now())
	if err != nil {
		t.Fatalf("RefreshScheduled() error = %v", err)
	}
	if result.AnyAttempted() || len(result.Locations) != 0 {
		t.Fatalf("result = %+v, want no-op", result)
	}
}

func TestScopedUpdaterManualRejectsUnconfiguredScope(t *testing.T) {
	u := ScopedUpdater{Scope: fakeScopeStore{}, Updater: Updater{}}
	if _, err := u.RefreshManual(context.Background(), time.Now()); err == nil {
		t.Fatal("RefreshManual() expected unconfigured error")
	}
}

func TestScopedUpdaterPropagatesScopeLoadFailure(t *testing.T) {
	u := ScopedUpdater{Scope: fakeScopeStore{err: errors.New("broken scope")}}
	if _, err := u.RefreshScheduled(context.Background(), time.Now()); err == nil {
		t.Fatal("RefreshScheduled() expected scope load error")
	}
}
