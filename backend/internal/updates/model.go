package updates

import (
	"fmt"
	"time"
)

const StateSchemaVersion = 1

// Release identifies a complete signed appliance update and its detached
// signature.
type Release struct {
	Version      string    `json:"version"`
	PublishedAt  time.Time `json:"published_at"`
	PageURL      string    `json:"page_url"`
	BundleURL    string    `json:"bundle_url"`
	SignatureURL string    `json:"signature_url"`
}

// State is persisted in shared appliance storage. Later orchestration stages
// can extend this versioned schema with download and installation state.
type State struct {
	SchemaVersion    int        `json:"schema_version"`
	CurrentVersion   string     `json:"current_version,omitempty"`
	AvailableRelease *Release   `json:"available_release,omitempty"`
	FirstDetectedAt  *time.Time `json:"first_detected_at,omitempty"`
	ApprovalDeadline *time.Time `json:"approval_deadline,omitempty"`
	LastCheckedAt    *time.Time `json:"last_checked_at,omitempty"`
	LastCheckError   string     `json:"last_check_error,omitempty"`
}

func DefaultState(currentVersion string) State {
	return State{
		SchemaVersion:  StateSchemaVersion,
		CurrentVersion: currentVersion,
	}
}

func (s State) normalized() State {
	if s.SchemaVersion == 0 {
		s.SchemaVersion = StateSchemaVersion
	}
	return cloneState(s)
}

func (s State) Validate() error {
	if s.SchemaVersion != StateSchemaVersion {
		return fmt.Errorf(
			"updates: unsupported state schema version %d",
			s.SchemaVersion,
		)
	}

	if s.AvailableRelease == nil {
		if s.FirstDetectedAt != nil || s.ApprovalDeadline != nil {
			return fmt.Errorf(
				"updates: detection dates require an available release",
			)
		}
		return nil
	}

	release := s.AvailableRelease
	if release.Version == "" {
		return fmt.Errorf(
			"updates: available release version is required",
		)
	}
	if release.PublishedAt.IsZero() {
		return fmt.Errorf(
			"updates: available release publication time is required",
		)
	}
	if release.PageURL == "" ||
		release.BundleURL == "" ||
		release.SignatureURL == "" {
		return fmt.Errorf(
			"updates: available release URLs are required",
		)
	}
	if s.FirstDetectedAt == nil {
		return fmt.Errorf(
			"updates: available release detection time is required",
		)
	}
	if s.ApprovalDeadline == nil {
		return fmt.Errorf(
			"updates: available release approval deadline is required",
		)
	}
	if s.ApprovalDeadline.Before(*s.FirstDetectedAt) {
		return fmt.Errorf(
			"updates: approval deadline precedes detection time",
		)
	}

	return nil
}

func cloneState(state State) State {
	cloned := state

	if state.AvailableRelease != nil {
		release := *state.AvailableRelease
		cloned.AvailableRelease = &release
	}
	if state.FirstDetectedAt != nil {
		value := *state.FirstDetectedAt
		cloned.FirstDetectedAt = &value
	}
	if state.ApprovalDeadline != nil {
		value := *state.ApprovalDeadline
		cloned.ApprovalDeadline = &value
	}
	if state.LastCheckedAt != nil {
		value := *state.LastCheckedAt
		cloned.LastCheckedAt = &value
	}

	return cloned
}
