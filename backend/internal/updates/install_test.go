package updates

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"
)

type installerFunc func(context.Context, string) error

func (f installerFunc) Install(ctx context.Context, version string) error {
	return f(ctx, version)
}

type rebooterFunc func(context.Context) error

func (f rebooterFunc) Reboot(ctx context.Context) error {
	return f(ctx)
}

func newInstallController(
	t *testing.T,
	now time.Time,
	installer UpdateInstaller,
	rebooter DeviceRebooter,
) (*Controller, *Store) {
	t.Helper()
	controller, store := newTestController(t, "v1.6.0", nil, now)
	if err := store.Save(installableState(now)); err != nil {
		t.Fatal(err)
	}
	controller.SetInstaller(installer, rebooter)
	return controller, store
}

func immediateInstall(now time.Time) InstallConditions {
	return InstallConditions{Now: now, Immediate: true}
}

func TestInstallPersistsAttemptBeforeStagingAndRebootPendingBeforeReboot(
	t *testing.T,
) {
	now := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	var store *Store
	installer := installerFunc(func(_ context.Context, version string) error {
		state, err := store.Load()
		if err != nil {
			t.Fatal(err)
		}
		if state.Installation == nil ||
			state.Installation.Status != InstallStatusInstalling ||
			state.Installation.Version != version ||
			len(state.Installation.Attempts) != 1 {
			t.Fatalf("state before staging = %+v", state.Installation)
		}
		return nil
	})
	rebooted := false
	rebooter := rebooterFunc(func(context.Context) error {
		state, err := store.Load()
		if err != nil {
			t.Fatal(err)
		}
		if state.Installation == nil ||
			state.Installation.Status != InstallStatusRebootPending ||
			state.Installation.StagedAt == nil {
			t.Fatalf("state before reboot = %+v", state.Installation)
		}
		rebooted = true
		return nil
	})
	controller, stateStore := newInstallController(t, now, installer, rebooter)
	store = stateStore

	state, err := controller.Install(t.Context(), immediateInstall(now))
	if err != nil {
		t.Fatal(err)
	}
	if !rebooted || state.Installation.Status != InstallStatusRebootPending {
		t.Fatalf("rebooted=%t installation=%+v", rebooted, state.Installation)
	}
	attempt := state.Installation.Attempts[0]
	if !attempt.Immediate || !attempt.Succeeded || attempt.FinishedAt == nil {
		t.Fatalf("attempt = %+v", attempt)
	}
}

func TestInstallFailureIsPersistedAndDoesNotReboot(t *testing.T) {
	now := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	wantErr := errors.New("write failed")
	rebootCalls := 0
	controller, store := newInstallController(
		t,
		now,
		installerFunc(func(context.Context, string) error { return wantErr }),
		rebooterFunc(func(context.Context) error { rebootCalls++; return nil }),
	)

	_, err := controller.Install(t.Context(), immediateInstall(now))
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
	state, loadErr := store.Load()
	if loadErr != nil {
		t.Fatal(loadErr)
	}
	if rebootCalls != 0 || state.Installation.Status != InstallStatusFailed ||
		!strings.Contains(state.Installation.LastError, wantErr.Error()) {
		t.Fatalf("reboots=%d installation=%+v", rebootCalls, state.Installation)
	}
}

func TestInstallRebootFailureLeavesStagedState(t *testing.T) {
	now := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	wantErr := errors.New("reboot failed")
	controller, store := newInstallController(
		t,
		now,
		installerFunc(func(context.Context, string) error { return nil }),
		rebooterFunc(func(context.Context) error { return wantErr }),
	)

	_, err := controller.Install(t.Context(), immediateInstall(now))
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
	state, loadErr := store.Load()
	if loadErr != nil {
		t.Fatal(loadErr)
	}
	if state.Installation.Status != InstallStatusRebootPending ||
		!strings.Contains(state.Installation.LastError, wantErr.Error()) {
		t.Fatalf("installation = %+v", state.Installation)
	}
}

func TestInstallRejectsConcurrentAttempt(t *testing.T) {
	now := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	started := make(chan struct{})
	release := make(chan struct{})
	installer := installerFunc(func(context.Context, string) error {
		close(started)
		<-release
		return nil
	})
	controller, _ := newInstallController(
		t, now, installer, rebooterFunc(func(context.Context) error { return nil }),
	)

	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		defer wait.Done()
		_, _ = controller.Install(context.Background(), immediateInstall(now))
	}()
	<-started
	if _, err := controller.Install(t.Context(), immediateInstall(now)); err == nil || !strings.Contains(err.Error(), InstallStatusInstalling) {
		t.Fatalf("concurrent error = %v", err)
	}
	close(release)
	wait.Wait()
}

func TestInstallBlockedByPolicyDoesNotInvokeCommands(t *testing.T) {
	now := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	calls := 0
	controller, store := newInstallController(
		t,
		now,
		installerFunc(func(context.Context, string) error { calls++; return nil }),
		rebooterFunc(func(context.Context) error { calls++; return nil }),
	)
	state, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	state.Download.Status = DownloadStatusFailed
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}

	if _, err := controller.Install(t.Context(), immediateInstall(now)); err == nil {
		t.Fatal("Install() error = nil")
	}
	if calls != 0 {
		t.Fatalf("command calls = %d", calls)
	}
}

func TestInstallDoesNotRestageRebootPendingUpdate(t *testing.T) {
	now := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	calls := 0
	controller, store := newInstallController(
		t,
		now,
		installerFunc(func(context.Context, string) error { calls++; return nil }),
		rebooterFunc(func(context.Context) error { calls++; return nil }),
	)
	state, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	state.Installation = &InstallationState{
		Version:   state.AvailableRelease.Version,
		Status:    InstallStatusRebootPending,
		StartedAt: &now,
		StagedAt:  &now,
		Attempts: []InstallAttempt{{
			StartedAt:  now,
			FinishedAt: &now,
			Succeeded:  true,
		}},
	}
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}

	if _, err := controller.Install(t.Context(), immediateInstall(now)); err == nil || !strings.Contains(err.Error(), InstallStatusRebootPending) {
		t.Fatalf("Install() error = %v", err)
	}
	if calls != 0 {
		t.Fatalf("command calls = %d", calls)
	}
}

func TestInstallDoesNotRestageTrialOrConfirmedUpdate(t *testing.T) {
	now := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	for _, status := range []string{InstallStatusProbation, InstallStatusInstalled} {
		t.Run(status, func(t *testing.T) {
			calls := 0
			controller, store := newInstallController(
				t,
				now,
				installerFunc(func(context.Context, string) error { calls++; return nil }),
				rebooterFunc(func(context.Context) error { calls++; return nil }),
			)
			state, err := store.Load()
			if err != nil {
				t.Fatal(err)
			}
			state.Installation = &InstallationState{
				Version:   state.AvailableRelease.Version,
				Status:    status,
				StartedAt: &now,
				StagedAt:  &now,
				Attempts: []InstallAttempt{{
					StartedAt: now,
				}},
			}
			if status == InstallStatusInstalled {
				state.Installation.ConfirmedAt = &now
			}
			if err := store.Save(state); err != nil {
				t.Fatal(err)
			}

			if _, err := controller.Install(t.Context(), immediateInstall(now)); err == nil || !strings.Contains(err.Error(), status) {
				t.Fatalf("Install() error = %v", err)
			}
			if calls != 0 {
				t.Fatalf("command calls = %d", calls)
			}
		})
	}
}

func TestRetainedAttemptsDropsHistoryOlderThanNinetyDays(t *testing.T) {
	now := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	old := now.Add(-FailureHistoryRetention - time.Second)
	recent := now.Add(-FailureHistoryRetention)
	installation := &InstallationState{Attempts: []InstallAttempt{
		{StartedAt: old},
		{StartedAt: recent},
	}}
	attempts := retainedAttempts(installation, now)
	if len(attempts) != 1 || !attempts[0].StartedAt.Equal(recent) {
		t.Fatalf("attempts = %+v", attempts)
	}
}

func TestCommandInstallerRequiresDownloadDirectory(t *testing.T) {
	err := (CommandInstaller{}).Install(t.Context(), "v1.7.0")
	if err == nil || !strings.Contains(err.Error(), "directory") {
		t.Fatalf("error = %v", err)
	}
}

func TestStartInstallReturnsBeforeStagingCompletes(t *testing.T) {
	now := time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)
	started := make(chan struct{})
	finish := make(chan struct{})
	rebooted := make(chan struct{})
	controller, store := newInstallController(
		t,
		now,
		installerFunc(func(context.Context, string) error {
			close(started)
			<-finish
			return nil
		}),
		rebooterFunc(func(context.Context) error {
			close(rebooted)
			return nil
		}),
	)

	state, err := controller.StartInstall(immediateInstall(now))
	if err != nil {
		t.Fatal(err)
	}
	if state.Installation == nil ||
		state.Installation.Status != InstallStatusInstalling {
		t.Fatalf("installation = %+v", state.Installation)
	}

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("installer did not start")
	}
	persisted, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if persisted.Installation.Status != InstallStatusInstalling {
		t.Fatalf("persisted status = %q", persisted.Installation.Status)
	}

	close(finish)
	select {
	case <-rebooted:
	case <-time.After(time.Second):
		t.Fatal("reboot was not requested")
	}
}
