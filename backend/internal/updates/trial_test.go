package updates

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCommandTrialStatus(t *testing.T) {
	command := filepath.Join(t.TempDir(), "masjidframe-ab")
	if err := os.WriteFile(command, []byte(`#!/bin/sh
test "$1" = machine-status || exit 2
printf '%s\n' running_slot=b active_slot=b rollback_slot=a upgrade_available=1
`), 0755); err != nil {
		t.Fatal(err)
	}
	status, err := (CommandTrialStatus{Command: command}).Status(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if status.RunningSlot != "b" || !status.UpgradeAvailable {
		t.Fatalf("status = %+v", status)
	}
}

func TestParseTrialStatus(t *testing.T) {
	status, err := parseTrialStatus(strings.Join([]string{
		"running_slot=b",
		"active_slot=b",
		"rollback_slot=a",
		"upgrade_available=1",
		"bootcount=1",
	}, "\n"))
	if err != nil {
		t.Fatal(err)
	}
	if status.RunningSlot != "b" || status.ActiveSlot != "b" ||
		status.RollbackSlot != "a" || !status.UpgradeAvailable {
		t.Fatalf("status = %+v", status)
	}
}

func TestParseTrialStatusRejectsIncompleteState(t *testing.T) {
	if _, err := parseTrialStatus("running_slot=a\nupgrade_available=0\n"); err == nil {
		t.Fatal("parseTrialStatus() error = nil")
	}
}

func TestReconcileTrialLifecycle(t *testing.T) {
	now := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	controller, store := trialTestController(t, now, "v1.7.0")

	probation, err := controller.ReconcileTrial(TrialStatus{
		RunningSlot:      "b",
		ActiveSlot:       "b",
		RollbackSlot:     "a",
		UpgradeAvailable: true,
	}, "v1.7.0")
	if err != nil {
		t.Fatal(err)
	}
	if probation.Installation.Status != InstallStatusProbation {
		t.Fatalf("status = %q", probation.Installation.Status)
	}

	confirmed, err := controller.ReconcileTrial(TrialStatus{
		RunningSlot:  "b",
		ActiveSlot:   "b",
		RollbackSlot: "b",
	}, "v1.7.0")
	if err != nil {
		t.Fatal(err)
	}
	if confirmed.Installation.Status != InstallStatusInstalled ||
		confirmed.Installation.ConfirmedAt == nil {
		t.Fatalf("installation = %+v", confirmed.Installation)
	}

	persisted, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if persisted.Installation.Status != InstallStatusInstalled {
		t.Fatalf("persisted status = %q", persisted.Installation.Status)
	}
}

func TestReconcileTrialUsesSignedInstalledReleaseIdentity(t *testing.T) {
	now := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	controller, store := trialTestController(t, now, "v1.6.0-rc.5-image")
	installedRelease := "v1.6.0-lab.65a1fc6.3"
	persisted, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	release := testRelease(installedRelease)
	persisted.AvailableRelease = &release
	persisted.Installation.Version = installedRelease
	if err := store.Save(persisted); err != nil {
		t.Fatal(err)
	}

	state, err := controller.ReconcileTrial(TrialStatus{
		RunningSlot:      "b",
		ActiveSlot:       "b",
		RollbackSlot:     "a",
		UpgradeAvailable: true,
	}, installedRelease)
	if err != nil {
		t.Fatal(err)
	}
	if state.Installation.Status != InstallStatusProbation {
		t.Fatalf("status = %q", state.Installation.Status)
	}
	if state.CurrentVersion != "v1.6.0-rc.5-image" {
		t.Fatalf("current version = %q", state.CurrentVersion)
	}

	state, err = controller.ReconcileTrial(TrialStatus{
		RunningSlot:  "b",
		ActiveSlot:   "b",
		RollbackSlot: "b",
	}, installedRelease)
	if err != nil {
		t.Fatal(err)
	}
	if state.Installation.Status != InstallStatusInstalled {
		t.Fatalf("status = %q", state.Installation.Status)
	}
}

func TestReadInstalledReleaseVersion(t *testing.T) {
	path := filepath.Join(t.TempDir(), "release.json")
	if err := os.WriteFile(path, []byte(`{"release_version":"v1.7.0"}`), 0644); err != nil {
		t.Fatal(err)
	}
	version, err := ReadInstalledReleaseVersion(path)
	if err != nil {
		t.Fatal(err)
	}
	if version != "v1.7.0" {
		t.Fatalf("version = %q", version)
	}
}

func TestReadInstalledReleaseVersionRejectsInvalidRecords(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{name: "malformed", content: `{`},
		{name: "missing version", content: `{}`},
		{name: "blank version", content: `{"release_version":" "}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "release.json")
			if err := os.WriteFile(path, []byte(test.content), 0644); err != nil {
				t.Fatal(err)
			}
			if _, err := ReadInstalledReleaseVersion(path); err == nil {
				t.Fatal("ReadInstalledReleaseVersion() error = nil")
			}
		})
	}
}

func TestReconcileTrialRecordsRollbackFailure(t *testing.T) {
	now := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	controller, _ := trialTestController(t, now, "v1.6.0")
	state, err := controller.ReconcileTrial(TrialStatus{
		RunningSlot:  "a",
		ActiveSlot:   "a",
		RollbackSlot: "a",
	}, "v1.6.0")
	if err != nil {
		t.Fatal(err)
	}
	installation := state.Installation
	if installation.Status != InstallStatusRolledBack ||
		installation.RolledBackAt == nil ||
		installation.Attempts[0].Error == "" {
		t.Fatalf("installation = %+v", installation)
	}
	if !AutomaticInstallDue(state, now.Add(36*time.Hour)) {
		t.Fatal("retry should be allowed on the next maintenance night")
	}
}

func trialTestController(
	t *testing.T,
	now time.Time,
	currentVersion string,
) (*Controller, *Store) {
	t.Helper()
	store := NewStore(filepath.Join(t.TempDir(), "update_state.json"))
	release := testRelease("v1.7.0")
	detected := now.Add(-time.Hour)
	deadline := detected.Add(ApprovalPeriod)
	staged := now.Add(-time.Minute)
	finished := staged
	if err := store.Save(State{
		SchemaVersion:    StateSchemaVersion,
		CurrentVersion:   "v1.6.0",
		AvailableRelease: &release,
		FirstDetectedAt:  &detected,
		ApprovalDeadline: &deadline,
		Installation: &InstallationState{
			Version:   release.Version,
			Status:    InstallStatusRebootPending,
			StartedAt: &staged,
			StagedAt:  &staged,
			Attempts: []InstallAttempt{{
				StartedAt:  staged,
				FinishedAt: &finished,
				Succeeded:  true,
			}},
		},
	}); err != nil {
		t.Fatal(err)
	}
	controller := NewController(store, nil, currentVersion)
	controller.now = func() time.Time { return now }
	return controller, store
}

type artifactCleanerFunc func(string) error

func (f artifactCleanerFunc) Cleanup(version string) error {
	return f(version)
}

func TestReconcileTrialCleansArtifactsOnlyAfterConfirmation(t *testing.T) {
	now := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	controller, store := trialTestController(t, now, "v1.7.0")
	state, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	state.Download = &DownloadState{
		Version: state.AvailableRelease.Version,
		Status:  DownloadStatusVerified,
	}
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}

	calls := 0
	controller.artifactCleaner = artifactCleanerFunc(func(version string) error {
		calls++
		if version != "v1.7.0" {
			t.Fatalf("cleanup version = %q", version)
		}
		return nil
	})

	probation, err := controller.ReconcileTrial(TrialStatus{
		RunningSlot:      "b",
		ActiveSlot:       "b",
		RollbackSlot:     "a",
		UpgradeAvailable: true,
	}, "v1.7.0")
	if err != nil {
		t.Fatal(err)
	}
	if probation.Download == nil || calls != 0 {
		t.Fatalf("probation download=%+v cleanup calls=%d", probation.Download, calls)
	}

	confirmed, err := controller.ReconcileTrial(TrialStatus{
		RunningSlot:  "b",
		ActiveSlot:   "b",
		RollbackSlot: "b",
	}, "v1.7.0")
	if err != nil {
		t.Fatal(err)
	}
	if confirmed.Download != nil || calls != 1 {
		t.Fatalf("confirmed download=%+v cleanup calls=%d", confirmed.Download, calls)
	}
}

func TestReconcileTrialRetainsArtifactsAfterRollback(t *testing.T) {
	now := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	controller, store := trialTestController(t, now, "v1.6.0")
	state, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	state.Download = &DownloadState{
		Version: state.AvailableRelease.Version,
		Status:  DownloadStatusVerified,
	}
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}
	calls := 0
	controller.artifactCleaner = artifactCleanerFunc(func(string) error {
		calls++
		return nil
	})

	rolledBack, err := controller.ReconcileTrial(TrialStatus{
		RunningSlot:  "a",
		ActiveSlot:   "a",
		RollbackSlot: "a",
	}, "v1.6.0")
	if err != nil {
		t.Fatal(err)
	}
	if rolledBack.Download == nil || calls != 0 {
		t.Fatalf("rollback download=%+v cleanup calls=%d", rolledBack.Download, calls)
	}
}
