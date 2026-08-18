package scope

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStateConfigured(t *testing.T) {
	if (State{}).Configured() {
		t.Fatal("zero-value state must be unconfigured")
	}
	if !(State{Country: "South Africa", Region: "North West", City: "Brits"}).Configured() {
		t.Fatal("country/city scope should be configured")
	}
}

func TestValidateAllowsBlankRegion(t *testing.T) {
	state := State{Country: "South Africa", City: "Brits"}
	if err := state.Validate(); err != nil {
		t.Fatalf("blank upstream region should be allowed: %v", err)
	}
}

func TestValidateRejectsUnconfigured(t *testing.T) {
	if err := (State{}).Validate(); err == nil {
		t.Fatal("expected unconfigured state to be rejected")
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

func TestStoreRoundTripNormalizes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "scope.json")
	store := NewStore(path)

	want := State{Country: " South Africa ", Region: " North West ", City: " Brits "}
	if err := store.Save(want); err != nil {
		t.Fatal(err)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	expected := State{Country: "South Africa", Region: "North West", City: "Brits"}
	if got != expected {
		t.Fatalf("got %#v, want %#v", got, expected)
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
	state := State{Country: "South Africa", Region: "North West", City: "Brits"}
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
