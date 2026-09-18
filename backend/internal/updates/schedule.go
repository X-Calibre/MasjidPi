package updates

import "time"

const CheckInterval = 7 * 24 * time.Hour

// CheckDue reports whether automatic release discovery should run. A device
// that has never checked is immediately due. Future timestamps are treated as
// due so a bad clock cannot suppress checks indefinitely.
func CheckDue(
	state State,
	now time.Time,
) bool {
	if state.LastCheckedAt == nil {
		return true
	}

	now = now.UTC()
	lastChecked := state.LastCheckedAt.UTC()

	if lastChecked.After(now) {
		return true
	}

	return !now.Before(
		lastChecked.Add(CheckInterval),
	)
}
