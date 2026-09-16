package updates

import (
	"strings"
	"testing"
	"time"
)

func installableState(now time.Time) State {
	release := testRelease("v1.7.0")
	detected := now.Add(-24 * time.Hour)
	deadline := detected.Add(ApprovalPeriod)
	verified := now.Add(-time.Hour)
	approved := now.Add(-30 * time.Minute)
	return State{
		SchemaVersion:    StateSchemaVersion,
		CurrentVersion:   "v1.6.0",
		AvailableRelease: &release,
		FirstDetectedAt:  &detected,
		ApprovalDeadline: &deadline,
		ApprovedAt:       &approved,
		Download: &DownloadState{
			Version:    release.Version,
			Status:     DownloadStatusVerified,
			VerifiedAt: &verified,
		},
	}
}

func TestEvaluateInstallAllowsApprovedVerifiedUpdateInWindow(t *testing.T) {
	now := time.Date(2026, time.September, 16, 23, 15, 0, 0, time.Local)
	decision := EvaluateInstall(installableState(now), InstallConditions{Now: now})
	if !decision.Allowed {
		t.Fatalf("decision = %+v", decision)
	}
	wantEnd := time.Date(2026, time.September, 17, 3, 30, 0, 0, time.Local)
	if decision.WindowEndsAt == nil || !decision.WindowEndsAt.Equal(wantEnd) {
		t.Fatalf("window end = %v, want %v", decision.WindowEndsAt, wantEnd)
	}
}

func TestEvaluateInstallWindowBoundaries(t *testing.T) {
	location := time.FixedZone("SAST", 2*60*60)
	tests := []struct {
		name    string
		now     time.Time
		allowed bool
	}{
		{"before window", time.Date(2026, 9, 16, 22, 59, 59, 0, location), false},
		{"window opens", time.Date(2026, 9, 16, 23, 0, 0, 0, location), true},
		{"after midnight", time.Date(2026, 9, 17, 2, 59, 59, 0, location), true},
		{"window closes", time.Date(2026, 9, 17, 3, 0, 0, 0, location), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decision := EvaluateInstall(installableState(test.now), InstallConditions{Now: test.now})
			if decision.Allowed != test.allowed {
				t.Fatalf("decision = %+v, want allowed=%t", decision, test.allowed)
			}
		})
	}
}

func TestEvaluateInstallAutomaticallyApprovesAtDeadline(t *testing.T) {
	now := time.Date(2026, 9, 16, 23, 30, 0, 0, time.UTC)
	state := installableState(now)
	state.ApprovedAt = nil
	deadline := now
	state.ApprovalDeadline = &deadline
	decision := EvaluateInstall(state, InstallConditions{Now: now})
	if !decision.Allowed || !decision.AutomaticallyApproved {
		t.Fatalf("decision = %+v", decision)
	}
}

func TestEvaluateInstallHonorsPostponement(t *testing.T) {
	now := time.Date(2026, 9, 16, 23, 30, 0, 0, time.UTC)
	state := installableState(now)
	state.ApprovedAt = nil
	postponed := now.Add(24 * time.Hour)
	state.PostponedUntil = &postponed
	decision := EvaluateInstall(state, InstallConditions{Now: now})
	if decision.Allowed || !strings.Contains(decision.Reason, "postponed") {
		t.Fatalf("decision = %+v", decision)
	}
}

func TestEvaluateInstallPlaybackPolicy(t *testing.T) {
	now := time.Date(2026, 9, 16, 23, 30, 0, 0, time.UTC)
	state := installableState(now)
	if decision := EvaluateInstall(state, InstallConditions{
		Now: now, PlaybackActive: true,
	}); decision.Allowed {
		t.Fatalf("automatic decision = %+v", decision)
	}
	if decision := EvaluateInstall(state, InstallConditions{
		Now: now, Immediate: true, PlaybackActive: true,
	}); decision.Allowed {
		t.Fatalf("unconfirmed immediate decision = %+v", decision)
	}
	if decision := EvaluateInstall(state, InstallConditions{
		Now: now, Immediate: true, PlaybackActive: true, InterruptPlayback: true,
	}); !decision.Allowed {
		t.Fatalf("confirmed immediate decision = %+v", decision)
	}
}

func TestEvaluateInstallBlocksWithinAdhanGuard(t *testing.T) {
	now := time.Date(2026, 9, 16, 23, 30, 0, 0, time.UTC)
	state := installableState(now)
	for _, offset := range []time.Duration{-15 * time.Minute, 15 * time.Minute} {
		decision := EvaluateInstall(state, InstallConditions{
			Now: now, AdhanTimes: []time.Time{now.Add(offset)},
		})
		if decision.Allowed || !strings.Contains(decision.Reason, "Adhan") {
			t.Fatalf("offset %v decision = %+v", offset, decision)
		}
	}
}

func TestEvaluateInstallRequiresVerifiedBundle(t *testing.T) {
	now := time.Date(2026, 9, 16, 23, 30, 0, 0, time.UTC)
	state := installableState(now)
	state.Download.Status = DownloadStatusFailed
	decision := EvaluateInstall(state, InstallConditions{Now: now, Immediate: true})
	if decision.Allowed || !strings.Contains(decision.Reason, "verified") {
		t.Fatalf("decision = %+v", decision)
	}
}
