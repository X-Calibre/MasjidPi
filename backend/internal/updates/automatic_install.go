package updates

import "time"

const MaximumFailedInstallNights = 3

// AutomaticInstallDue prevents repeated staging attempts during one night and
// stops automatic retries after three distinct failed nights. A newly
// discovered release receives a fresh Installation state and retry budget.
func AutomaticInstallDue(state State, now time.Time) bool {
	if state.Installation == nil {
		return true
	}
	if state.Installation.Status == InstallStatusInstalling ||
		state.Installation.Status == InstallStatusRebootPending ||
		state.Installation.Status == InstallStatusProbation ||
		state.Installation.Status == InstallStatusInstalled {
		return false
	}

	failedNights := make(map[string]struct{})
	for _, attempt := range state.Installation.Attempts {
		if attempt.Error == "" {
			continue
		}
		night := installNight(attempt.StartedAt.In(now.Location()))
		failedNights[night] = struct{}{}
	}
	if len(failedNights) >= MaximumFailedInstallNights {
		return false
	}
	_, alreadyFailedTonight := failedNights[installNight(now)]
	return !alreadyFailedTonight
}

func installNight(value time.Time) string {
	if value.Hour() < InstallWindowEndHour {
		value = value.AddDate(0, 0, -1)
	}
	return value.Format("2006-01-02")
}
