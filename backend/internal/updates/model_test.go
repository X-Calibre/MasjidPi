package updates

import (
	"strings"
	"testing"
	"time"
)

func validState() State {
	published := time.Date(
		2026,
		time.September,
		16,
		6,
		0,
		0,
		0,
		time.UTC,
	)
	detected := published.Add(2 * time.Hour)
	deadline := detected.Add(30 * 24 * time.Hour)

	return State{
		SchemaVersion:  StateSchemaVersion,
		CurrentVersion: "v1.6.0",
		AvailableRelease: &Release{
			Version:      "v1.7.0",
			PublishedAt:  published,
			PageURL:      "https://example.test/releases/v1.7.0",
			BundleURL:    "https://example.test/update.tar.zst",
			SignatureURL: "https://example.test/update.tar.zst.minisig",
		},
		FirstDetectedAt:  &detected,
		ApprovalDeadline: &deadline,
	}
}

func TestStateValidateAcceptsCompleteAvailableRelease(t *testing.T) {
	if err := validState().Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestStateValidateRejectsIncompleteAvailableRelease(t *testing.T) {
	state := DefaultState("v1.6.0")
	state.AvailableRelease = &Release{
		Version: "v1.7.0",
	}

	err := state.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want validation error")
	}
	if !strings.Contains(err.Error(), "publication time") {
		t.Fatalf("Validate() error = %q", err)
	}
}

func TestStateValidateRejectsDetectionWithoutRelease(t *testing.T) {
	state := DefaultState("v1.6.0")
	now := time.Now().UTC()
	state.FirstDetectedAt = &now

	err := state.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want validation error")
	}
	if !strings.Contains(
		err.Error(),
		"require an available release",
	) {
		t.Fatalf("Validate() error = %q", err)
	}
}

func TestStateValidateRejectsDeadlineBeforeDetection(t *testing.T) {
	state := validState()
	invalid := state.FirstDetectedAt.Add(-time.Minute)
	state.ApprovalDeadline = &invalid

	err := state.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want deadline error")
	}
	if !strings.Contains(
		err.Error(),
		"deadline precedes detection",
	) {
		t.Fatalf("Validate() error = %q", err)
	}
}

func TestCloneStateDoesNotAliasPointers(t *testing.T) {
	state := validState()
	originalVersion := state.AvailableRelease.Version
	originalDetected := *state.FirstDetectedAt

	cloned := cloneState(state)
	cloned.AvailableRelease.Version = "v2.0.0"
	changed := originalDetected.Add(time.Hour)
	*cloned.FirstDetectedAt = changed

	if state.AvailableRelease.Version != originalVersion {
		t.Fatal("release pointer was aliased")
	}
	if !state.FirstDetectedAt.Equal(originalDetected) {
		t.Fatal("detection-time pointer was aliased")
	}
}
