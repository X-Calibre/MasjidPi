package updates

import (
	"testing"
	"time"
)

func TestControllerApproveAvailableRelease(t *testing.T) {
	now := time.Date(2026, time.September, 16, 9, 0, 0, 0, time.UTC)
	controller, _ := newTestController(t, "v1.6.0", &sequenceReleaseSource{
		results: []sourceResult{{release: testRelease("v1.7.0"), found: true}},
	}, now)
	if _, err := controller.Check(t.Context()); err != nil {
		t.Fatal(err)
	}

	state, err := controller.Approve()
	if err != nil {
		t.Fatal(err)
	}
	if state.ApprovedAt == nil || !state.ApprovedAt.Equal(now) {
		t.Fatalf("approved at = %v, want %v", state.ApprovedAt, now)
	}
	if state.PostponedUntil != nil {
		t.Fatalf("postponed until = %v, want nil", state.PostponedUntil)
	}
}

func TestControllerPostponeHonorsOriginalDeadline(t *testing.T) {
	now := time.Date(2026, time.September, 16, 9, 0, 0, 0, time.UTC)
	controller, _ := newTestController(t, "v1.6.0", &sequenceReleaseSource{
		results: []sourceResult{{release: testRelease("v1.7.0"), found: true}},
	}, now)
	if _, err := controller.Check(t.Context()); err != nil {
		t.Fatal(err)
	}

	until := now.Add(7 * 24 * time.Hour)
	state, err := controller.Postpone(until)
	if err != nil {
		t.Fatal(err)
	}
	if state.PostponedUntil == nil || !state.PostponedUntil.Equal(until) {
		t.Fatalf("postponed until = %v, want %v", state.PostponedUntil, until)
	}
	if _, err := controller.Postpone(now.Add(ApprovalPeriod + time.Second)); err == nil {
		t.Fatal("postponement beyond deadline succeeded")
	}
}

func TestControllerDecisionRequiresAvailableRelease(t *testing.T) {
	now := time.Date(2026, time.September, 16, 9, 0, 0, 0, time.UTC)
	controller, _ := newTestController(t, "v1.6.0", &sequenceReleaseSource{}, now)
	if _, err := controller.Approve(); err == nil {
		t.Fatal("approval without a release succeeded")
	}
	if _, err := controller.Postpone(now.Add(time.Hour)); err == nil {
		t.Fatal("postponement without a release succeeded")
	}
}
