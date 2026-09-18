package updates

import (
	"testing"
	"time"
)

func TestControllerFirstDetectionStartsApprovalPeriod(
	t *testing.T,
) {
	now := time.Date(
		2026,
		time.September,
		16,
		9,
		0,
		0,
		0,
		time.UTC,
	)
	source := &sequenceReleaseSource{
		results: []sourceResult{{
			release: testRelease("v1.7.0"),
			found:   true,
		}},
	}
	controller, store := newTestController(
		t,
		"v1.6.0",
		source,
		now,
	)

	state, err := controller.Check(
		t.Context(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if state.CurrentVersion != "v1.6.0" {
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
	if state.FirstDetectedAt == nil ||
		!state.FirstDetectedAt.Equal(now) {
		t.Fatalf(
			"first detected = %v, want %v",
			state.FirstDetectedAt,
			now,
		)
	}

	wantDeadline := now.Add(ApprovalPeriod)
	if state.ApprovalDeadline == nil ||
		!state.ApprovalDeadline.Equal(wantDeadline) {
		t.Fatalf(
			"approval deadline = %v, want %v",
			state.ApprovalDeadline,
			wantDeadline,
		)
	}
	if state.LastCheckedAt == nil ||
		!state.LastCheckedAt.Equal(now) {
		t.Fatalf(
			"last checked = %v, want %v",
			state.LastCheckedAt,
			now,
		)
	}
	if state.LastCheckError != "" {
		t.Fatalf(
			"last check error = %q",
			state.LastCheckError,
		)
	}

	persisted, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if persisted.AvailableRelease == nil ||
		persisted.AvailableRelease.Version != "v1.7.0" {
		t.Fatalf(
			"persisted release = %+v",
			persisted.AvailableRelease,
		)
	}
}

func TestControllerRepeatedDetectionKeepsOriginalDeadline(
	t *testing.T,
) {
	firstCheck := time.Date(
		2026,
		time.September,
		16,
		9,
		0,
		0,
		0,
		time.UTC,
	)
	secondCheck := firstCheck.Add(7 * 24 * time.Hour)

	source := &sequenceReleaseSource{
		results: []sourceResult{
			{
				release: testRelease("v1.7.0"),
				found:   true,
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

	first, err := controller.Check(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	controller.now = func() time.Time {
		return secondCheck
	}
	second, err := controller.Check(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if second.FirstDetectedAt == nil ||
		first.FirstDetectedAt == nil ||
		!second.FirstDetectedAt.Equal(
			*first.FirstDetectedAt,
		) {
		t.Fatalf(
			"first detection changed from %v to %v",
			first.FirstDetectedAt,
			second.FirstDetectedAt,
		)
	}
	if second.ApprovalDeadline == nil ||
		first.ApprovalDeadline == nil ||
		!second.ApprovalDeadline.Equal(
			*first.ApprovalDeadline,
		) {
		t.Fatalf(
			"deadline changed from %v to %v",
			first.ApprovalDeadline,
			second.ApprovalDeadline,
		)
	}
	if second.LastCheckedAt == nil ||
		!second.LastCheckedAt.Equal(secondCheck) {
		t.Fatalf(
			"last checked = %v, want %v",
			second.LastCheckedAt,
			secondCheck,
		)
	}
}
