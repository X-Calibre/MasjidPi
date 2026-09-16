package updates

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

const ApprovalPeriod = 30 * 24 * time.Hour

type ReleaseSource interface {
	Latest(context.Context) (Release, bool, error)
}

type StateStore interface {
	Load() (State, error)
	Save(State) error
}

type Controller struct {
	mu             sync.Mutex
	store          StateStore
	source         ReleaseSource
	currentVersion string
	now            func() time.Time
}

func NewController(
	store StateStore,
	source ReleaseSource,
	currentVersion string,
) *Controller {
	return &Controller{
		store:          store,
		source:         source,
		currentVersion: currentVersion,
		now:            time.Now,
	}
}

// Status returns the persisted state with the running version refreshed from
// the current binary. It performs no network request and no disk write.
func (c *Controller) Status() (State, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	state, err := c.store.Load()
	if err != nil {
		return State{}, err
	}
	state.CurrentVersion = c.currentVersion
	return cloneState(state), nil
}

// Check discovers the latest complete stable release and atomically records
// the result. The first detection time and approval deadline remain unchanged
// while the same release stays available.
func (c *Controller) Check(
	ctx context.Context,
) (State, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	state, err := c.store.Load()
	if err != nil {
		return State{}, err
	}

	now := c.now().UTC()
	state.CurrentVersion = c.currentVersion
	state.LastCheckedAt = timePointer(now)

	release, found, checkErr := c.source.Latest(ctx)
	if checkErr != nil {
		return c.recordCheckError(
			state,
			fmt.Errorf(
				"updates: discover latest release: %w",
				checkErr,
			),
		)
	}

	if !found {
		clearAvailableRelease(&state)
		state.LastCheckError = ""
		return c.saveState(state)
	}

	candidateVersion, err := ParseStableVersion(
		release.Version,
	)
	if err != nil {
		return c.recordCheckError(state, err)
	}

	newer, err := IsNewerThanInstalled(
		candidateVersion,
		c.currentVersion,
	)
	if err != nil {
		return c.recordCheckError(state, err)
	}

	if !newer {
		clearAvailableRelease(&state)
		state.LastCheckError = ""
		return c.saveState(state)
	}

	if state.AvailableRelease == nil ||
		state.AvailableRelease.Version != release.Version {
		detected := now
		deadline := now.Add(ApprovalPeriod)
		state.FirstDetectedAt = &detected
		state.ApprovalDeadline = &deadline
		state.ApprovedAt = nil
		state.PostponedUntil = nil
	}

	releaseCopy := release
	state.AvailableRelease = &releaseCopy
	state.LastCheckError = ""

	return c.saveState(state)
}

// Approve records consent to install the currently available release.
func (c *Controller) Approve() (State, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	state, err := c.store.Load()
	if err != nil {
		return State{}, err
	}
	state.CurrentVersion = c.currentVersion
	if state.AvailableRelease == nil {
		return cloneState(state), fmt.Errorf("updates: no release is available for approval")
	}

	now := c.now().UTC()
	state.ApprovedAt = &now
	state.PostponedUntil = nil
	return c.saveState(state)
}

// Postpone defers approval until a time no later than the original deadline.
func (c *Controller) Postpone(until time.Time) (State, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	state, err := c.store.Load()
	if err != nil {
		return State{}, err
	}
	state.CurrentVersion = c.currentVersion
	if state.AvailableRelease == nil || state.ApprovalDeadline == nil {
		return cloneState(state), fmt.Errorf("updates: no release is available to postpone")
	}

	now := c.now().UTC()
	until = until.UTC()
	if !until.After(now) {
		return cloneState(state), fmt.Errorf("updates: postponement must be in the future")
	}
	if until.After(*state.ApprovalDeadline) {
		return cloneState(state), fmt.Errorf("updates: postponement exceeds approval deadline")
	}

	state.ApprovedAt = nil
	state.PostponedUntil = &until
	return c.saveState(state)
}

func (c *Controller) recordCheckError(
	state State,
	checkErr error,
) (State, error) {
	state.LastCheckError = checkErr.Error()

	saved, saveErr := c.saveState(state)
	if saveErr != nil {
		return saved, errors.Join(checkErr, saveErr)
	}
	return saved, checkErr
}

func (c *Controller) saveState(
	state State,
) (State, error) {
	if err := c.store.Save(state); err != nil {
		return cloneState(state), err
	}
	return cloneState(state), nil
}

func clearAvailableRelease(state *State) {
	state.AvailableRelease = nil
	state.FirstDetectedAt = nil
	state.ApprovalDeadline = nil
	state.ApprovedAt = nil
	state.PostponedUntil = nil
}

func timePointer(value time.Time) *time.Time {
	return &value
}
