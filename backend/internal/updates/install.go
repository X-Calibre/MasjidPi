package updates

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

type CommandInstaller struct {
	Directory string
	Command   string
}

func (i CommandInstaller) Install(ctx context.Context, version string) error {
	if _, err := ParseStableVersion(version); err != nil {
		return err
	}
	if i.Directory == "" {
		return fmt.Errorf("updates: download directory is required")
	}
	command := i.Command
	if command == "" {
		command = "/usr/local/sbin/masjidpi-update"
	}
	bundle := filepath.Join(
		i.Directory,
		fmt.Sprintf("masjidpi-update-%s-pi3.tar.zst", version),
	)
	output, err := exec.CommandContext(
		ctx,
		command,
		"install",
		bundle,
		bundle+".minisig",
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("updates: install bundle: %w: %s", err, string(output))
	}
	return nil
}

type CommandRebooter struct {
	Command string
}

func (r CommandRebooter) Reboot(ctx context.Context) error {
	command := r.Command
	if command == "" {
		command = "/usr/bin/systemctl"
	}
	if output, err := exec.CommandContext(ctx, command, "reboot").CombinedOutput(); err != nil {
		return fmt.Errorf("updates: request reboot: %w: %s", err, string(output))
	}
	return nil
}

func (c *Controller) Install(
	ctx context.Context,
	conditions InstallConditions,
) (State, error) {
	c.mu.Lock()
	state, err := c.store.Load()
	if err != nil {
		c.mu.Unlock()
		return State{}, err
	}
	state.CurrentVersion = c.currentVersion
	if state.Installation != nil &&
		(state.Installation.Status == InstallStatusInstalling ||
			state.Installation.Status == InstallStatusRebootPending ||
			state.Installation.Status == InstallStatusProbation ||
			state.Installation.Status == InstallStatusInstalled) {
		c.mu.Unlock()
		return cloneState(state), fmt.Errorf(
			"updates: installation state is %s",
			state.Installation.Status,
		)
	}
	decision := EvaluateInstall(state, conditions)
	if !decision.Allowed {
		c.mu.Unlock()
		return cloneState(state), fmt.Errorf("updates: installation blocked: %s", decision.Reason)
	}
	if c.installing {
		c.mu.Unlock()
		return cloneState(state), fmt.Errorf("updates: installation is already running")
	}
	if c.installer == nil || c.rebooter == nil {
		c.mu.Unlock()
		return cloneState(state), fmt.Errorf("updates: installer is unavailable")
	}

	now := c.now().UTC()
	version := state.AvailableRelease.Version
	attempts := retainedAttempts(state.Installation, now)
	attempts = append(attempts, InstallAttempt{
		StartedAt: now,
		Immediate: conditions.Immediate,
	})
	state.Installation = &InstallationState{
		Version:   version,
		Status:    InstallStatusInstalling,
		StartedAt: &now,
		Attempts:  attempts,
	}
	if _, err := c.saveState(state); err != nil {
		c.mu.Unlock()
		return cloneState(state), err
	}
	installer := c.installer
	rebooter := c.rebooter
	c.installing = true
	c.mu.Unlock()

	installErr := installer.Install(ctx, version)

	c.mu.Lock()
	c.installing = false
	state, loadErr := c.store.Load()
	if loadErr != nil {
		c.mu.Unlock()
		return State{}, loadErr
	}
	state.CurrentVersion = c.currentVersion
	finished := c.now().UTC()
	if state.Installation == nil || state.Installation.Version != version {
		c.mu.Unlock()
		return cloneState(state), fmt.Errorf("updates: installation state changed while staging")
	}
	if len(state.Installation.Attempts) == 0 {
		c.mu.Unlock()
		return cloneState(state), fmt.Errorf("updates: installation attempt disappeared while staging")
	}
	attempt := &state.Installation.Attempts[len(state.Installation.Attempts)-1]
	attempt.FinishedAt = &finished
	if installErr != nil {
		attempt.Error = installErr.Error()
		state.Installation.Status = InstallStatusFailed
		state.Installation.LastError = installErr.Error()
		saved, saveErr := c.saveState(state)
		c.mu.Unlock()
		if saveErr != nil {
			return saved, errors.Join(installErr, saveErr)
		}
		return saved, installErr
	}

	attempt.Succeeded = true
	state.Installation.Status = InstallStatusRebootPending
	state.Installation.StagedAt = &finished
	state.Installation.LastError = ""
	saved, saveErr := c.saveState(state)
	c.mu.Unlock()
	if saveErr != nil {
		return saved, saveErr
	}

	if rebootErr := rebooter.Reboot(ctx); rebootErr != nil {
		c.mu.Lock()
		latest, loadErr := c.store.Load()
		if loadErr == nil && latest.Installation != nil &&
			latest.Installation.Version == version {
			latest.Installation.LastError = rebootErr.Error()
			latest, loadErr = c.saveState(latest)
		}
		c.mu.Unlock()
		if loadErr != nil {
			return latest, errors.Join(rebootErr, loadErr)
		}
		return latest, rebootErr
	}
	return saved, nil
}

func retainedAttempts(
	installation *InstallationState,
	now time.Time,
) []InstallAttempt {
	if installation == nil {
		return nil
	}
	cutoff := now.Add(-FailureHistoryRetention)
	attempts := make([]InstallAttempt, 0, len(installation.Attempts))
	for _, attempt := range installation.Attempts {
		if !attempt.StartedAt.Before(cutoff) {
			attempts = append(attempts, attempt)
		}
	}
	return attempts
}
