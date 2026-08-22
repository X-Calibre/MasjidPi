package catalogue

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const DefaultRefreshInterval = 7 * 24 * time.Hour

type Source interface {
	Fetch(context.Context, Location, time.Time) (Catalogue, error)
}

type PartitionPersistence interface {
	Load() (State, error)
	SavePartition(Partition) error
	RetainLocations([]Location) error
}

type Updater struct {
	Source          Source
	Store           PartitionPersistence
	RefreshInterval time.Duration
}

type LocationRefreshResult struct {
	Location  Location
	Attempted bool
	Updated   bool
	Error     error
	Partition Partition
}

type RefreshResult struct{ Locations []LocationRefreshResult }

func (r RefreshResult) AnyAttempted() bool {
	for _, x := range r.Locations {
		if x.Attempted {
			return true
		}
	}
	return false
}
func (r RefreshResult) AnyFailed() bool {
	for _, x := range r.Locations {
		if x.Error != nil {
			return true
		}
	}
	return false
}

func (u Updater) RefreshScheduled(ctx context.Context, locations []Location, now time.Time) (RefreshResult, error) {
	return u.refresh(ctx, locations, now, false)
}
func (u Updater) RefreshManual(ctx context.Context, locations []Location, now time.Time) (RefreshResult, error) {
	return u.refresh(ctx, locations, now, true)
}

func (u Updater) refresh(ctx context.Context, locations []Location, now time.Time, force bool) (RefreshResult, error) {
	if u.Source == nil {
		return RefreshResult{}, fmt.Errorf("masjidboard catalogue: source is required")
	}
	if u.Store == nil {
		return RefreshResult{}, fmt.Errorf("masjidboard catalogue: store is required")
	}
	locations, err := normalizeLocations(locations)
	if err != nil {
		return RefreshResult{}, err
	}
	if err := u.Store.RetainLocations(locations); err != nil {
		return RefreshResult{}, fmt.Errorf("masjidboard catalogue: retain configured locations: %w", err)
	}

	state, err := u.Store.Load()
	if err != nil {
		return RefreshResult{}, err
	}
	current := make(map[string]Partition, len(state.Partitions))
	for _, partition := range state.Partitions {
		current[partition.Location.key()] = clonePartition(partition)
	}

	interval := u.RefreshInterval
	if interval <= 0 {
		interval = DefaultRefreshInterval
	}
	result := RefreshResult{Locations: make([]LocationRefreshResult, 0, len(locations))}
	for _, location := range locations {
		previous, exists := current[location.key()]
		item := LocationRefreshResult{Location: location, Partition: previous}
		if !force && exists && !previous.ValidatedAt.IsZero() && now.Before(previous.ValidatedAt.Add(interval)) {
			result.Locations = append(result.Locations, item)
			continue
		}
		item.Attempted = true
		candidate, fetchErr := u.Source.Fetch(ctx, location, now)
		if fetchErr != nil {
			item.Error = fetchErr
			result.Locations = append(result.Locations, item)
			continue
		}
		if candidate.RetrievedAt.IsZero() {
			candidate.RetrievedAt = now
		}
		if candidate.ValidatedAt.IsZero() {
			candidate.ValidatedAt = now
		}
		currentCatalogue := Catalogue{}
		if exists {
			currentCatalogue = Catalogue{RetrievedAt: previous.RetrievedAt, ValidatedAt: previous.ValidatedAt, Records: previous.Records}
		}
		reconciled, reconcileErr := Reconcile(currentCatalogue, candidate)
		if reconcileErr != nil {
			item.Error = reconcileErr
			result.Locations = append(result.Locations, item)
			continue
		}
		partition := Partition{Location: location, RetrievedAt: reconciled.RetrievedAt, ValidatedAt: reconciled.ValidatedAt, Records: reconciled.Records}
		if saveErr := u.Store.SavePartition(partition); saveErr != nil {
			item.Error = saveErr
			result.Locations = append(result.Locations, item)
			continue
		}
		item.Updated = !exists || !samePartition(previous, partition)
		item.Partition = partition
		current[location.key()] = clonePartition(partition)
		result.Locations = append(result.Locations, item)
	}
	return result, nil
}

func normalizeLocations(locations []Location) ([]Location, error) {
	if len(locations) < 1 || len(locations) > 3 {
		return nil, fmt.Errorf("masjidboard catalogue: configured locations must contain 1 to 3 entries")
	}
	seen := make(map[string]struct{}, len(locations))
	out := make([]Location, 0, len(locations))
	for i, location := range locations {
		location = location.Normalized()
		if err := location.Validate(); err != nil {
			return nil, fmt.Errorf("masjidboard catalogue: location %d: %w", i+1, err)
		}
		key := location.key()
		if _, ok := seen[key]; ok {
			return nil, fmt.Errorf("masjidboard catalogue: duplicate location %q / %q / %q", location.Country, location.Region, location.City)
		}
		seen[key] = struct{}{}
		out = append(out, location)
	}
	return out, nil
}

func samePartition(a, b Partition) bool {
	if a.Location.key() != b.Location.key() || !a.RetrievedAt.Equal(b.RetrievedAt) || !a.ValidatedAt.Equal(b.ValidatedAt) || len(a.Records) != len(b.Records) {
		return false
	}
	for i := range a.Records {
		x, y := a.Records[i], b.Records[i]
		if x.ID != y.ID || x.Provider != y.Provider || x.ExternalID != y.ExternalID || x.Name != y.Name || x.City != y.City || x.Region != y.Region || x.Country != y.Country || x.TimeZoneOffsetMS != y.TimeZoneOffsetMS || x.Status != y.Status || !x.DiscoveredAt.Equal(y.DiscoveredAt) || !x.LastSeenAt.Equal(y.LastSeenAt) {
			return false
		}
		if len(x.ProviderMetadata) != len(y.ProviderMetadata) {
			return false
		}
		for key, value := range x.ProviderMetadata {
			if strings.TrimSpace(y.ProviderMetadata[key]) != strings.TrimSpace(value) {
				return false
			}
		}
	}
	return true
}
