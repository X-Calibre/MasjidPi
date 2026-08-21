package scope

import (
	"os"
	"path/filepath"
	"testing"
)

func loc(country, region, city string) Location {
	return Location{Country: country, Region: region, City: city}
}

func TestStateConfigured(t *testing.T) {
	if (State{}).Configured() {
		t.Fatal("zero-value state must be unconfigured")
	}
	if !(State{Locations: []Location{loc("South Africa", "North West", "Brits")}}).Configured() {
		t.Fatal("one valid location should be configured")
	}
	if !(State{Locations: []Location{
		loc("South Africa", "North West", "Brits"),
		loc("South Africa", "Gauteng", "Akasia"),
		loc("South Africa", "Gauteng", "Pretoria"),
	}}).Configured() {
		t.Fatal("three valid locations should be configured")
	}
}

func TestValidateAllowsBlankRegion(t *testing.T) {
	state := State{Locations: []Location{loc("South Africa", "", "Brits")}}
	if err := state.Validate(); err != nil {
		t.Fatalf("blank upstream region should be allowed: %v", err)
	}
}

func TestValidateRejectsUnconfigured(t *testing.T) {
	if err := (State{}).Validate(); err == nil {
		t.Fatal("expected unconfigured state to be rejected")
	}
}

func TestValidateRejectsMoreThanThreeLocations(t *testing.T) {
	state := State{Locations: []Location{
		loc("South Africa", "Gauteng", "A"),
		loc("South Africa", "Gauteng", "B"),
		loc("South Africa", "Gauteng", "C"),
		loc("South Africa", "Gauteng", "D"),
	}}
	if err := state.Validate(); err == nil {
		t.Fatal("expected four locations to be rejected")
	}
}

func TestValidateRejectsDuplicateLocationsCaseInsensitively(t *testing.T) {
	state := State{Locations: []Location{
		loc("South Africa", "North West", "Brits"),
		loc(" south africa ", " north west ", " BRITS "),
	}}
	if err := state.Validate(); err == nil {
		t.Fatal("expected duplicate locations to be rejected")
	}
}

func TestValidateAllowsLocationsAcrossRegionsAndCountries(t *testing.T) {
	state := State{Locations: []Location{
		loc("South Africa", "North West", "Brits"),
		loc("South Africa", "Gauteng", "Akasia"),
		loc("Botswana", "", "Gaborone"),
	}}
	if err := state.Validate(); err != nil {
		t.Fatalf("cross-region/country locations should be allowed: %v", err)
	}
}

func TestStoreMissingIsUnconfigured(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "scope.json"))
	state, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if state.Configured() {
		t.Fatalf("missing store should be unconfigured: %#v", state)
	}
}

func TestStoreRoundTripNormalizesAndPreservesOrder(t *testing.T) {
	path := filepath.Join(t.TempDir(), "scope.json")
	store := NewStore(path)

	want := State{Locations: []Location{
		loc(" South Africa ", " North West ", " Brits "),
		loc(" South Africa ", " Gauteng ", " Akasia "),
	}}
	if err := store.Save(want); err != nil {
		t.Fatal(err)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	expected := State{Locations: []Location{
		loc("South Africa", "North West", "Brits"),
		loc("South Africa", "Gauteng", "Akasia"),
	}}
	if len(got.Locations) != len(expected.Locations) {
		t.Fatalf("got %#v, want %#v", got, expected)
	}
	for i := range expected.Locations {
		if got.Locations[i] != expected.Locations[i] {
			t.Fatalf("location %d = %#v, want %#v", i, got.Locations[i], expected.Locations[i])
		}
	}
}

func TestStoreRejectsEmptyConfiguredSave(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "scope.json"))
	if err := store.Save(State{}); err == nil {
		t.Fatal("expected zero-value save to fail")
	}
}

func TestStoreUnchangedSaveDoesNotRewrite(t *testing.T) {
	path := filepath.Join(t.TempDir(), "scope.json")
	store := NewStore(path)
	state := State{Locations: []Location{loc("South Africa", "North West", "Brits")}}
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}

	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}
	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !before.ModTime().Equal(after.ModTime()) {
		t.Fatalf("unchanged save rewrote file: before=%v after=%v", before.ModTime(), after.ModTime())
	}
}
