package catalogue

import (
	"context"
	"fmt"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/scope"
)

// ScopePersistence supplies the persisted 1-3 location discovery scope.
type ScopePersistence interface {
	Load() (scope.State, error)
}

// ScopedUpdater connects persisted discovery-scope configuration to the
// provider-neutral partitioned catalogue updater. It deliberately does not
// retain either scope or catalogue state in memory between operations.
type ScopedUpdater struct {
	Scope   ScopePersistence
	Updater Updater
}

// RefreshScheduled loads the current configured scope and refreshes only
// location partitions that are due according to Updater.RefreshInterval.
// An unconfigured scope is a no-op rather than an application error.
func (u ScopedUpdater) RefreshScheduled(ctx context.Context, now time.Time) (RefreshResult, error) {
	locations, configured, err := u.locations()
	if err != nil {
		return RefreshResult{}, err
	}
	if !configured {
		return RefreshResult{}, nil
	}
	return u.Updater.RefreshScheduled(ctx, locations, now)
}

// RefreshManual loads the current configured scope and immediately attempts
// every configured location. Manual refresh still requires MasjidBoard
// discovery to be configured because there is no meaningful target otherwise.
func (u ScopedUpdater) RefreshManual(ctx context.Context, now time.Time) (RefreshResult, error) {
	locations, configured, err := u.locations()
	if err != nil {
		return RefreshResult{}, err
	}
	if !configured {
		return RefreshResult{}, fmt.Errorf("masjidboard catalogue: discovery scope is not configured")
	}
	return u.Updater.RefreshManual(ctx, locations, now)
}

func (u ScopedUpdater) locations() ([]Location, bool, error) {
	if u.Scope == nil {
		return nil, false, fmt.Errorf("masjidboard catalogue: scope store is required")
	}
	state, err := u.Scope.Load()
	if err != nil {
		return nil, false, fmt.Errorf("masjidboard catalogue: load discovery scope: %w", err)
	}
	if !state.Configured() {
		return nil, false, nil
	}

	locations := make([]Location, 0, len(state.Locations))
	for _, location := range state.Locations {
		locations = append(locations, Location{
			Country: location.Country,
			Region:  location.Region,
			City:    location.City,
		})
	}
	return locations, true, nil
}
