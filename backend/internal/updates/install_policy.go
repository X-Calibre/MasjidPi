package updates

import "time"

const (
	InstallWindowStartHour = 23
	InstallWindowEndHour   = 3
	InstallOverrun         = 30 * time.Minute
	AdhanGuardPeriod       = 15 * time.Minute
)

type InstallEvaluation struct {
	Allowed               bool       `json:"allowed"`
	Reason                string     `json:"reason,omitempty"`
	AutomaticallyApproved bool       `json:"automatically_approved,omitempty"`
	WindowEndsAt          *time.Time `json:"window_ends_at,omitempty"`
}

type InstallConditions struct {
	Now               time.Time
	Immediate         bool
	PlaybackActive    bool
	InterruptPlayback bool
	AdhanTimes        []time.Time
}

// EvaluateInstall applies the appliance installation policy without changing
// update state or invoking the installer.
func EvaluateInstall(
	state State,
	conditions InstallConditions,
) InstallEvaluation {
	now := conditions.Now
	if now.IsZero() {
		now = time.Now()
	}

	if state.AvailableRelease == nil {
		return blockedInstall("no update is available")
	}
	if state.Download == nil ||
		state.Download.Version != state.AvailableRelease.Version ||
		state.Download.Status != DownloadStatusVerified {
		return blockedInstall("the update is not downloaded and verified")
	}

	automaticApproval := state.ApprovalDeadline != nil &&
		!now.Before(*state.ApprovalDeadline)
	approved := state.ApprovedAt != nil || automaticApproval
	if !approved && !conditions.Immediate {
		if state.PostponedUntil != nil && now.Before(*state.PostponedUntil) {
			return blockedInstall("the update is postponed")
		}
		return blockedInstall("the update is awaiting approval")
	}

	if conditions.PlaybackActive &&
		(!conditions.Immediate || !conditions.InterruptPlayback) {
		return blockedInstall("audio playback is active")
	}
	if withinAdhanGuard(now, conditions.AdhanTimes) {
		return blockedInstall("an Adhan is within 15 minutes")
	}

	decision := InstallEvaluation{
		Allowed:               true,
		AutomaticallyApproved: automaticApproval,
	}
	if conditions.Immediate {
		return decision
	}

	windowEnd, withinWindow := installationWindow(now)
	if !withinWindow {
		return blockedInstall("outside the 23:00–03:00 installation window")
	}
	finishBy := windowEnd.Add(InstallOverrun)
	decision.WindowEndsAt = &finishBy
	return decision
}

func installationWindow(now time.Time) (time.Time, bool) {
	location := now.Location()
	year, month, day := now.Date()
	start := time.Date(year, month, day, InstallWindowStartHour, 0, 0, 0, location)
	end := time.Date(year, month, day+1, InstallWindowEndHour, 0, 0, 0, location)
	if now.Hour() < InstallWindowEndHour {
		start = time.Date(year, month, day-1, InstallWindowStartHour, 0, 0, 0, location)
		end = time.Date(year, month, day, InstallWindowEndHour, 0, 0, 0, location)
	}
	return end, !now.Before(start) && now.Before(end)
}

func withinAdhanGuard(now time.Time, adhanTimes []time.Time) bool {
	for _, adhan := range adhanTimes {
		difference := now.Sub(adhan)
		if difference < 0 {
			difference = -difference
		}
		if difference <= AdhanGuardPeriod {
			return true
		}
	}
	return false
}

func blockedInstall(reason string) InstallEvaluation {
	return InstallEvaluation{Reason: reason}
}
