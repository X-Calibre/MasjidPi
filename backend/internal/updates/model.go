package updates

import (
	"fmt"
	"time"
)

const StateSchemaVersion = 1

const (
	DownloadStatusDownloading = "downloading"
	DownloadStatusVerified    = "verified"
	DownloadStatusFailed      = "failed"
)

// Release identifies a complete signed appliance update and its detached
// signature.
type Release struct {
	Version      string    `json:"version"`
	PublishedAt  time.Time `json:"published_at"`
	PageURL      string    `json:"page_url"`
	BundleURL    string    `json:"bundle_url"`
	SignatureURL string    `json:"signature_url"`
}

type DownloadState struct {
	Version        string     `json:"version"`
	Status         string     `json:"status"`
	BundleBytes    int64      `json:"bundle_bytes,omitempty"`
	SignatureBytes int64      `json:"signature_bytes,omitempty"`
	StartedAt      *time.Time `json:"started_at,omitempty"`
	VerifiedAt     *time.Time `json:"verified_at,omitempty"`
	LastError      string     `json:"last_error,omitempty"`
}

// State is persisted in shared appliance storage. Later orchestration stages
// can extend this versioned schema with download and installation state.
type State struct {
	SchemaVersion    int            `json:"schema_version"`
	CurrentVersion   string         `json:"current_version,omitempty"`
	AvailableRelease *Release       `json:"available_release,omitempty"`
	FirstDetectedAt  *time.Time     `json:"first_detected_at,omitempty"`
	ApprovalDeadline *time.Time     `json:"approval_deadline,omitempty"`
	ApprovedAt       *time.Time     `json:"approved_at,omitempty"`
	PostponedUntil   *time.Time     `json:"postponed_until,omitempty"`
	Download         *DownloadState `json:"download,omitempty"`
	LastCheckedAt    *time.Time     `json:"last_checked_at,omitempty"`
	LastCheckError   string         `json:"last_check_error,omitempty"`
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
		if s.FirstDetectedAt != nil || s.ApprovalDeadline != nil ||
			s.ApprovedAt != nil || s.PostponedUntil != nil {
			return fmt.Errorf(
				"updates: detection dates require an available release",
			)
		}
		if s.Download != nil {
			return fmt.Errorf("updates: download requires an available release")
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
	if s.ApprovedAt != nil && s.PostponedUntil != nil {
		return fmt.Errorf("updates: release cannot be approved and postponed")
	}
	if s.ApprovedAt != nil && s.ApprovedAt.Before(*s.FirstDetectedAt) {
		return fmt.Errorf("updates: approval precedes detection time")
	}
	if s.PostponedUntil != nil {
		if s.PostponedUntil.Before(*s.FirstDetectedAt) {
			return fmt.Errorf("updates: postponement precedes detection time")
		}
		if s.PostponedUntil.After(*s.ApprovalDeadline) {
			return fmt.Errorf("updates: postponement exceeds approval deadline")
		}
	}
	if s.Download != nil {
		download := s.Download
		if download.Version != release.Version {
			return fmt.Errorf("updates: download version does not match available release")
		}
		switch download.Status {
		case DownloadStatusDownloading, DownloadStatusVerified, DownloadStatusFailed:
		default:
			return fmt.Errorf("updates: invalid download status %q", download.Status)
		}
		if download.Status == DownloadStatusVerified && download.VerifiedAt == nil {
			return fmt.Errorf("updates: verified download requires verification time")
		}
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
	if state.ApprovedAt != nil {
		value := *state.ApprovedAt
		cloned.ApprovedAt = &value
	}
	if state.PostponedUntil != nil {
		value := *state.PostponedUntil
		cloned.PostponedUntil = &value
	}
	if state.Download != nil {
		download := *state.Download
		cloned.Download = &download
		if state.Download.StartedAt != nil {
			value := *state.Download.StartedAt
			cloned.Download.StartedAt = &value
		}
		if state.Download.VerifiedAt != nil {
			value := *state.Download.VerifiedAt
			cloned.Download.VerifiedAt = &value
		}
	}

	return cloned
}
