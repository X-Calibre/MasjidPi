package app

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/updates"
)

type fakeScheduledUpdateChecker struct {
	status      updates.State
	statusErr   error
	checkResult updates.State
	checkErr    error
	statusCalls int
	checkCalls  int
}

func (f *fakeScheduledUpdateChecker) Status() (
	updates.State,
	error,
) {
	f.statusCalls++
	return f.status, f.statusErr
}

func (f *fakeScheduledUpdateChecker) Check(
	context.Context,
) (updates.State, error) {
	f.checkCalls++
	return f.checkResult, f.checkErr
}

func updateCheckTestLogger() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(io.Discard, nil),
	)
}

func TestCheckUpdatesIfDueChecksNeverCheckedDevice(
	t *testing.T,
) {
	now := time.Date(
		2026, time.September, 16,
		12, 0, 0, 0, time.UTC,
	)
	checker := &fakeScheduledUpdateChecker{
		status: updates.DefaultState("v1.6.0"),
		checkResult: updates.State{
			SchemaVersion:  updates.StateSchemaVersion,
			CurrentVersion: "v1.6.0",
			LastCheckedAt:  &now,
		},
	}

	checkUpdatesIfDue(
		t.Context(),
		checker,
		now,
		updateCheckTestLogger(),
	)

	if checker.statusCalls != 1 {
		t.Fatalf(
			"status calls = %d, want 1",
			checker.statusCalls,
		)
	}
	if checker.checkCalls != 1 {
		t.Fatalf(
			"check calls = %d, want 1",
			checker.checkCalls,
		)
	}
}

func TestCheckUpdatesIfDueSkipsRecentCheck(
	t *testing.T,
) {
	now := time.Date(
		2026, time.September, 16,
		12, 0, 0, 0, time.UTC,
	)
	lastChecked := now.Add(-24 * time.Hour)
	checker := &fakeScheduledUpdateChecker{
		status: updates.State{
			SchemaVersion:  updates.StateSchemaVersion,
			CurrentVersion: "v1.6.0",
			LastCheckedAt:  &lastChecked,
		},
	}

	checkUpdatesIfDue(
		t.Context(),
		checker,
		now,
		updateCheckTestLogger(),
	)

	if checker.statusCalls != 1 {
		t.Fatalf(
			"status calls = %d, want 1",
			checker.statusCalls,
		)
	}
	if checker.checkCalls != 0 {
		t.Fatalf(
			"check calls = %d, want 0",
			checker.checkCalls,
		)
	}
}

func TestCheckUpdatesIfDueStopsAfterStatusFailure(
	t *testing.T,
) {
	checker := &fakeScheduledUpdateChecker{
		statusErr: errors.New("state unavailable"),
	}

	checkUpdatesIfDue(
		t.Context(),
		checker,
		time.Now(),
		updateCheckTestLogger(),
	)

	if checker.statusCalls != 1 {
		t.Fatalf(
			"status calls = %d, want 1",
			checker.statusCalls,
		)
	}
	if checker.checkCalls != 0 {
		t.Fatalf(
			"check calls = %d, want 0",
			checker.checkCalls,
		)
	}
}

func TestCheckUpdatesIfDueHandlesDiscoveryFailure(
	t *testing.T,
) {
	checker := &fakeScheduledUpdateChecker{
		status:   updates.DefaultState("v1.6.0"),
		checkErr: errors.New("GitHub unavailable"),
	}

	checkUpdatesIfDue(
		t.Context(),
		checker,
		time.Now(),
		updateCheckTestLogger(),
	)

	if checker.statusCalls != 1 ||
		checker.checkCalls != 1 {
		t.Fatalf(
			"calls: status=%d check=%d, want 1 each",
			checker.statusCalls,
			checker.checkCalls,
		)
	}
}
