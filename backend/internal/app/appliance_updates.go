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
	mu       sync.RWMutex
	playback updatePlaybackProvider
	board    updateBoardProvider
	now      func() time.Time
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
