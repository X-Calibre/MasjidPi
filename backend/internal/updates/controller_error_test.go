package updates

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

func TestControllerDiscoveryErrorPreservesAvailableRelease(
	t *testing.T,
) {
	firstDetected := time.Date(
		2026, time.September, 1,
		8, 0, 0, 0, time.UTC,
	)
	deadline := firstDetected.Add(ApprovalPeriod)
	checkTime := time.Date(
		2026, time.September, 16,
		12, 0, 0, 0, time.UTC,
	)
	available := testRelease("v1.7.0")
	discoveryError := errors.New("temporary GitHub failure")

	source := &sequenceReleaseSource{
		results: []sourceResult{{
			err: discoveryError,
		}},
	}
	controller, store := newTestController(
		t,
		"v1.6.0",
		source,
		checkTime,
	)

	err := store.Save(State{
		SchemaVersion:    StateSchemaVersion,
		CurrentVersion:   "v1.6.0",
		AvailableRelease: &available,
		FirstDetectedAt:  &firstDetected,
		ApprovalDeadline: &deadline,
	})
	if err != nil {
		t.Fatal(err)
	}

	state, err := controller.Check(t.Context())
	if err == nil {
		t.Fatal("Check() error = nil")
	}
	if !strings.Contains(
		err.Error(),
		discoveryError.Error(),
	) {
		t.Fatalf("Check() error = %q", err)
	}

	if state.AvailableRelease == nil ||
		state.AvailableRelease.Version != "v1.7.0" {
		t.Fatalf(
			"available release = %+v",
			state.AvailableRelease,
		)
	}
	requireControllerTime(
		t,
		"first detected",
		state.FirstDetectedAt,
		firstDetected,
	)
	requireControllerTime(
		t,
		"approval deadline",
		state.ApprovalDeadline,
		deadline,
	)
	requireControllerTime(
		t,
		"last checked",
		state.LastCheckedAt,
		checkTime,
	)
	if !strings.Contains(
		state.LastCheckError,
		discoveryError.Error(),
	) {
		t.Fatalf(
			"last check error = %q",
			state.LastCheckError,
		)
	}

	persisted, loadErr := store.Load()
	if loadErr != nil {
		t.Fatal(loadErr)
	}
	if persisted.AvailableRelease == nil ||
		persisted.AvailableRelease.Version != "v1.7.0" {
		t.Fatalf(
			"persisted release = %+v",
			persisted.AvailableRelease,
		)
	}
	if persisted.LastCheckError != state.LastCheckError {
		t.Fatalf(
			"persisted error = %q, want %q",
			persisted.LastCheckError,
			state.LastCheckError,
		)
	}
}

func TestControllerSuccessfulCheckClearsPreviousError(
	t *testing.T,
) {
	firstCheck := time.Date(
		2026, time.September, 16,
		12, 0, 0, 0, time.UTC,
	)
	secondCheck := firstCheck.Add(time.Hour)

	source := &sequenceReleaseSource{
		results: []sourceResult{
			{
				err: errors.New("temporary failure"),
			},
			{
				release: testRelease("v1.7.0"),
				found:   true,
			},
		},
	}
	controller, _ := newTestController(
		t,
		"v1.6.0",
		source,
		firstCheck,
	)

	failed, err := controller.Check(t.Context())
	if err == nil {
		t.Fatal("first Check() error = nil")
	}
	if failed.LastCheckError == "" {
		t.Fatal("first check did not record an error")
	}

	controller.now = func() time.Time {
		return secondCheck
	}

	recovered, err := controller.Check(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if recovered.LastCheckError != "" {
		t.Fatalf(
			"last check error = %q, want empty",
			recovered.LastCheckError,
		)
	}
	if recovered.AvailableRelease == nil ||
		recovered.AvailableRelease.Version != "v1.7.0" {
		t.Fatalf(
			"available release = %+v",
			recovered.AvailableRelease,
		)
	}
	requireControllerTime(
		t,
		"last checked",
		recovered.LastCheckedAt,
		secondCheck,
	)
}

func TestControllerStatusDoesNotDiscoverOrWrite(
	t *testing.T,
) {
	detected := time.Date(
		2026, time.September, 1,
		8, 0, 0, 0, time.UTC,
	)
	deadline := detected.Add(ApprovalPeriod)
	available := testRelease("v1.7.0")

	source := &sequenceReleaseSource{}
	controller, store := newTestController(
		t,
		"v1.6.0-rc.4-image",
		source,
		time.Date(
			2026, time.September, 16,
			13, 0, 0, 0, time.UTC,
		),
	)

	err := store.Save(State{
		SchemaVersion:    StateSchemaVersion,
		CurrentVersion:   "stale-version",
		AvailableRelease: &available,
		FirstDetectedAt:  &detected,
		ApprovalDeadline: &deadline,
	})
	if err != nil {
		t.Fatal(err)
	}

	before, err := os.ReadFile(store.path)
	if err != nil {
		t.Fatal(err)
	}

	state, err := controller.Status()
	if err != nil {
		t.Fatal(err)
	}

	after, err := os.ReadFile(store.path)
	if err != nil {
		t.Fatal(err)
	}

	if source.calls != 0 {
		t.Fatalf(
			"release-source calls = %d, want 0",
			source.calls,
		)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("Status() rewrote persisted state")
	}
	if state.CurrentVersion != "v1.6.0-rc.4-image" {
		t.Fatalf(
			"current version = %q",
			state.CurrentVersion,
		)
	}
	if state.AvailableRelease == nil ||
		state.AvailableRelease.Version != "v1.7.0" {
		t.Fatalf(
			"available release = %+v",
			state.AvailableRelease,
		)
	}

	persisted, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if persisted.CurrentVersion != "stale-version" {
		t.Fatalf(
			"persisted current version = %q",
			persisted.CurrentVersion,
		)
	}
}
