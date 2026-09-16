package updates

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const DefaultReleaseRecordPath = "/usr/share/masjidpi/update/release.json"

type TrialStatus struct {
	RunningSlot      string
	ActiveSlot       string
	RollbackSlot     string
	UpgradeAvailable bool
}

type TrialStatusSource interface {
	Status(context.Context) (TrialStatus, error)
}

type CommandTrialStatus struct {
	Command string
}

func (s CommandTrialStatus) Status(ctx context.Context) (TrialStatus, error) {
	command := s.Command
	if command == "" {
		command = "/usr/local/sbin/masjidpi-ab"
	}
	output, err := exec.CommandContext(ctx, command, "machine-status").CombinedOutput()
	if err != nil {
		return TrialStatus{}, fmt.Errorf(
			"updates: read A/B trial status: %w: %s",
			err,
			strings.TrimSpace(string(output)),
		)
	}
	return parseTrialStatus(string(output))
}

func parseTrialStatus(output string) (TrialStatus, error) {
	values := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		key, value, found := strings.Cut(scanner.Text(), "=")
		if found {
			values[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		return TrialStatus{}, fmt.Errorf("updates: scan A/B trial status: %w", err)
	}
	status := TrialStatus{
		RunningSlot:      values["running_slot"],
		ActiveSlot:       values["active_slot"],
		RollbackSlot:     values["rollback_slot"],
		UpgradeAvailable: values["upgrade_available"] == "1",
	}
	if !validSlot(status.RunningSlot) ||
		!validSlot(status.ActiveSlot) ||
		!validSlot(status.RollbackSlot) ||
		(values["upgrade_available"] != "0" && values["upgrade_available"] != "1") {
		return TrialStatus{}, fmt.Errorf("updates: invalid A/B trial status")
	}
	return status, nil
}

func validSlot(slot string) bool { return slot == "a" || slot == "b" }

func ReadInstalledReleaseVersion(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("updates: read installed release: %w", err)
	}
	var record struct {
		ReleaseVersion string `json:"release_version"`
	}
	if err := json.Unmarshal(data, &record); err != nil {
		return "", fmt.Errorf("updates: decode installed release: %w", err)
	}
	version := strings.TrimSpace(record.ReleaseVersion)
	if version == "" {
		return "", fmt.Errorf("updates: installed release version is required")
	}
	return version, nil
}

func (c *Controller) ReconcileTrial(
	status TrialStatus,
	runningReleaseVersion string,
) (State, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	state, err := c.store.Load()
	if err != nil {
		return State{}, err
	}
	state.CurrentVersion = c.currentVersion
	if runningReleaseVersion == "" {
		runningReleaseVersion = c.currentVersion
	}
	installation := state.Installation
	if installation == nil ||
		(installation.Status != InstallStatusRebootPending &&
			installation.Status != InstallStatusProbation) {
		return cloneState(state), nil
	}

	now := c.now().UTC()
	if status.UpgradeAvailable {
		if runningReleaseVersion == installation.Version &&
			status.RunningSlot == status.ActiveSlot {
			installation.Status = InstallStatusProbation
			return c.saveState(state)
		}
		return cloneState(state), nil
	}

	stableSlot := status.RunningSlot == status.ActiveSlot &&
		status.RunningSlot == status.RollbackSlot
	if !stableSlot {
		return cloneState(state), nil
	}
	if runningReleaseVersion == installation.Version {
		installation.Status = InstallStatusInstalled
		installation.ConfirmedAt = &now
		installation.LastError = ""
		return c.saveState(state)
	}

	installation.Status = InstallStatusRolledBack
	installation.RolledBackAt = &now
	installation.LastError = "trial boot rolled back to the previous release"
	if len(installation.Attempts) != 0 {
		attempt := &installation.Attempts[len(installation.Attempts)-1]
		if attempt.Error == "" {
			attempt.Error = installation.LastError
		}
	}
	return c.saveState(state)
}
