package app

import (
	"context"
	"sync"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/updates"
)

// applianceUpdates binds the platform-neutral update controller to live
// appliance safety inputs without introducing playback or Board dependencies
// into the updates package.
type applianceUpdates struct {
	*updates.Controller
	mu                sync.RWMutex
	playback          updatePlaybackProvider
	board             updateBoardProvider
	trial             updates.TrialStatusSource
	releaseRecordPath string
	now               func() time.Time
}

func (a *applianceUpdates) SetTrialStatusSource(source updates.TrialStatusSource) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.trial = source
}

func (a *applianceUpdates) SetReleaseRecordPath(path string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.releaseRecordPath = path
}

func (a *applianceUpdates) Status() (updates.State, error) {
	a.mu.RLock()
	source := a.trial
	releaseRecordPath := a.releaseRecordPath
	a.mu.RUnlock()
	if source == nil {
		return a.Controller.Status()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	status, err := source.Status(ctx)
	if err != nil {
		return a.Controller.Status()
	}
	runningReleaseVersion := ""
	if releaseRecordPath != "" {
		version, readErr := updates.ReadInstalledReleaseVersion(
			releaseRecordPath,
		)
		if readErr != nil {
			return a.Controller.Status()
		}
		runningReleaseVersion = version
	}
	return a.Controller.ReconcileTrial(status, runningReleaseVersion)
}

func (a *applianceUpdates) Check(ctx context.Context) (updates.State, error) {
	state, err := a.Status()
	if err != nil {
		return state, err
	}
	if state.Installation != nil {
		switch state.Installation.Status {
		case updates.InstallStatusInstalling,
			updates.InstallStatusRebootPending,
			updates.InstallStatusProbation:
			return state, nil
		}
	}
	return a.Controller.Check(ctx)
}

func newApplianceUpdates(controller *updates.Controller) *applianceUpdates {
	return &applianceUpdates{Controller: controller, now: time.Now}
}

func (a *applianceUpdates) SetSafetyProviders(
	playback updatePlaybackProvider,
	board updateBoardProvider,
) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.playback = playback
	a.board = board
}

func (a *applianceUpdates) Install(
	ctx context.Context,
	immediate bool,
	interruptPlayback bool,
) (updates.State, error) {
	a.mu.RLock()
	playbackProvider := a.playback
	boardProvider := a.board
	now := a.now()
	a.mu.RUnlock()
	conditions := currentInstallConditions(
		now,
		immediate,
		interruptPlayback,
		playbackProvider,
		boardProvider,
	)
	state, err := a.Status()
	if err != nil {
		return state, err
	}
	if !updates.EvaluateInstall(state, conditions).Allowed {
		return a.Controller.Install(ctx, conditions)
	}
	if conditions.PlaybackActive && immediate && interruptPlayback {
		playbackProvider.Stop()
	}
	return a.Controller.Install(ctx, conditions)
}

func (a *applianceUpdates) InstallAutomatic(
	ctx context.Context,
	now time.Time,
) (updates.State, bool, error) {
	state, err := a.Status()
	if err != nil || !updates.AutomaticInstallDue(state, now) {
		return state, false, err
	}

	a.mu.RLock()
	playbackProvider := a.playback
	boardProvider := a.board
	a.mu.RUnlock()
	conditions := currentInstallConditions(
		now,
		false,
		false,
		playbackProvider,
		boardProvider,
	)
	if !updates.EvaluateInstall(state, conditions).Allowed {
		return state, false, nil
	}

	state, err = a.Controller.Install(ctx, conditions)
	return state, true, err
}
