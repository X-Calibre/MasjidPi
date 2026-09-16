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
	status           updates.State
	statusErr        error
	checkResult      updates.State
	checkErr         error
	prepareResult    updates.State
	prepareErr       error
	statusCalls      int
	checkCalls       int
	prepareCalls     int
	installResult    updates.State
	installErr       error
	installAttempted bool
	installCalls     int
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

func (f *fakeScheduledUpdateChecker) Prepare(
	context.Context,
) (updates.State, error) {
	f.prepareCalls++
	return f.prepareResult, f.prepareErr
}

func (f *fakeScheduledUpdateChecker) InstallAutomatic(
	context.Context,
	time.Time,
) (updates.State, bool, error) {
	f.installCalls++
	return f.installResult, f.installAttempted, f.installErr
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

func TestCheckUpdatesIfDuePreparesAvailableRelease(t *testing.T) {
	now := time.Date(2026, time.September, 16, 12, 0, 0, 0, time.UTC)
	release := updates.Release{Version: "v1.7.0"}
	checker := &fakeScheduledUpdateChecker{
		status: updates.DefaultState("v1.6.0"),
		checkResult: updates.State{
			SchemaVersion:    updates.StateSchemaVersion,
			CurrentVersion:   "v1.6.0",
			AvailableRelease: &release,
		},
		prepareResult: updates.State{
			SchemaVersion:    updates.StateSchemaVersion,
			CurrentVersion:   "v1.6.0",
			AvailableRelease: &release,
			Download: &updates.DownloadState{
				Version: "v1.7.0",
				Status:  updates.DownloadStatusVerified,
			},
		},
	}

	checkUpdatesIfDue(t.Context(), checker, now, updateCheckTestLogger())

	if checker.checkCalls != 1 || checker.prepareCalls != 1 {
		t.Fatalf(
			"calls: check=%d prepare=%d, want 1 each",
			checker.checkCalls,
			checker.prepareCalls,
		)
	}
}

func TestCheckUpdatesIfDueStagesEligibleVerifiedRelease(t *testing.T) {
	now := time.Date(2026, 9, 16, 23, 0, 0, 0, time.UTC)
	lastChecked := now
	release := updates.Release{Version: "v1.7.0"}
	checker := &fakeScheduledUpdateChecker{
		status: updates.State{
			SchemaVersion:    updates.StateSchemaVersion,
			AvailableRelease: &release,
			LastCheckedAt:    &lastChecked,
			Download: &updates.DownloadState{
				Version: release.Version,
				Status:  updates.DownloadStatusVerified,
			},
		},
		installResult: updates.State{
			SchemaVersion:    updates.StateSchemaVersion,
			AvailableRelease: &release,
		},
		installAttempted: true,
	}

	checkUpdatesIfDue(t.Context(), checker, now, updateCheckTestLogger())

	if checker.installCalls != 1 {
		t.Fatalf("install calls = %d, want 1", checker.installCalls)
	}
}
