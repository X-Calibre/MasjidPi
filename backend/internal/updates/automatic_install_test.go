package updates

import (
	"testing"
	"time"
)

func TestAutomaticInstallDueLimitsRetriesByNight(t *testing.T) {
	location := time.FixedZone("SAST", 2*60*60)
	now := time.Date(2026, 9, 18, 23, 30, 0, 0, location)
	failed := func(daysAgo int) InstallAttempt {
		return InstallAttempt{
			StartedAt: now.AddDate(0, 0, -daysAgo),
			Error:     "staging failed",
		}
	}

	tests := []struct {
		name         string
		installation *InstallationState
		want         bool
	}{
		{name: "never attempted", want: true},
		{
			name: "failed earlier night",
			installation: &InstallationState{
				Status:   InstallStatusFailed,
				Attempts: []InstallAttempt{failed(1)},
			},
			want: true,
		},
		{
			name: "failed tonight",
			installation: &InstallationState{
				Status:   InstallStatusFailed,
				Attempts: []InstallAttempt{failed(0)},
			},
		},
		{
			name: "three failed nights",
			installation: &InstallationState{
				Status:   InstallStatusFailed,
				Attempts: []InstallAttempt{failed(1), failed(2), failed(3)},
			},
		},
		{
			name:         "reboot pending",
			installation: &InstallationState{Status: InstallStatusRebootPending},
		},
		{
			name:         "probation",
			installation: &InstallationState{Status: InstallStatusProbation},
		},
		{
			name:         "installed",
			installation: &InstallationState{Status: InstallStatusInstalled},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := State{Installation: test.installation}
			if got := AutomaticInstallDue(state, now); got != test.want {
				t.Fatalf("AutomaticInstallDue() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestAutomaticInstallDueCountsDistinctLocalDates(t *testing.T) {
	location := time.FixedZone("SAST", 2*60*60)
	now := time.Date(2026, 9, 18, 23, 0, 0, 0, location)
	state := State{Installation: &InstallationState{
		Status: InstallStatusFailed,
		Attempts: []InstallAttempt{
			{StartedAt: time.Date(2026, 9, 17, 21, 30, 0, 0, time.UTC), Error: "one"},
			{StartedAt: time.Date(2026, 9, 17, 22, 30, 0, 0, time.UTC), Error: "two"},
		},
	}}
	if !AutomaticInstallDue(state, now) {
		t.Fatal("two attempts on one earlier local night exhausted retry budget")
	}
}
