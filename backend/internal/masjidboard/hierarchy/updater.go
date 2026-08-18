package hierarchy

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

const DefaultRefreshInterval = 7 * 24 * time.Hour

// Source supplies the lightweight MasjidBoard discovery hierarchy. Provider
// implementations should preserve upstream counts and blank region names.
type Source interface {
	Countries(context.Context) ([]Location, error)
	Regions(context.Context, string) ([]Location, error)
	Cities(context.Context, string, string) ([]Location, error)
}

// Persistence is the last-known-good hierarchy store used by Updater.
type Persistence interface {
	Load() (State, error)
	Save(State) error
}

// Updater builds, validates and transactionally persists the complete
// country -> region -> city hierarchy. Scheduled refreshes are due weekly;
// manual refreshes bypass the age check.
type Updater struct {
	Source          Source
	Store           Persistence
	RefreshInterval time.Duration
}

type RefreshResult struct {
	Attempted bool
	Updated   bool
	State     State
}

// RefreshScheduled refreshes only when no valid hierarchy exists or the last
// successful validation is at least RefreshInterval old.
func (u Updater) RefreshScheduled(ctx context.Context, now time.Time) (RefreshResult, error) {
	current, err := u.current()
	if err != nil {
		return RefreshResult{}, err
	}
	interval := u.RefreshInterval
	if interval <= 0 {
		interval = DefaultRefreshInterval
	}
	if !current.ValidatedAt.IsZero() && now.Before(current.ValidatedAt.Add(interval)) {
		return RefreshResult{Attempted: false, Updated: false, State: current}, nil
	}
	return u.refresh(ctx, now, current)
}

// RefreshManual always attempts an upstream refresh regardless of age.
func (u Updater) RefreshManual(ctx context.Context, now time.Time) (RefreshResult, error) {
	current, err := u.current()
	if err != nil {
		return RefreshResult{}, err
	}
	return u.refresh(ctx, now, current)
}

func (u Updater) current() (State, error) {
	if u.Store == nil {
		return State{}, fmt.Errorf("masjidboard hierarchy: store is required")
	}
	return u.Store.Load()
}

func (u Updater) refresh(ctx context.Context, now time.Time, current State) (RefreshResult, error) {
	if u.Source == nil {
		return RefreshResult{}, fmt.Errorf("masjidboard hierarchy: source is required")
	}
	candidate, err := u.build(ctx, now)
	if err != nil {
		return RefreshResult{Attempted: true, State: current}, err
	}
	if err := candidate.Validate(); err != nil {
		return RefreshResult{Attempted: true, State: current}, err
	}
	if err := validateCounts(candidate); err != nil {
		return RefreshResult{Attempted: true, State: current}, err
	}
	if err := u.Store.Save(candidate); err != nil {
		return RefreshResult{Attempted: true, State: current}, err
	}
	return RefreshResult{Attempted: true, Updated: !sameHierarchy(current, candidate), State: candidate}, nil
}

func (u Updater) build(ctx context.Context, now time.Time) (State, error) {
	countries, err := u.Source.Countries(ctx)
	if err != nil {
		return State{}, fmt.Errorf("masjidboard hierarchy: fetch countries: %w", err)
	}
	countries = mergeLocations(countries, false)

	state := State{RetrievedAt: now, ValidatedAt: now}
	for _, countryEntry := range countries {
		country := Country{Name: countryEntry.Name, Count: countryEntry.Count}

		regions, err := u.Source.Regions(ctx, country.Name)
		if err != nil {
			return State{}, fmt.Errorf("masjidboard hierarchy: fetch regions for %q: %w", country.Name, err)
		}
		regions = mergeLocations(regions, true)

		for _, regionEntry := range regions {
			region := Region{Name: regionEntry.Name, Count: regionEntry.Count}

			// Some countries have no province layer at all; those are represented
			// by a single blank region and are still resolvable through Cities(). A
			// blank bucket alongside named regions is different: live FindMasjid
			// has been observed to advertise such a bucket but return HTTP 500 when
			// asked for its cities. Preserve its board count as unresolved rather
			// than discarding the entire otherwise-valid global hierarchy.
			if region.Name == "" && len(regions) > 1 {
				region.UnresolvedCount = region.Count
				country.Regions = append(country.Regions, region)
				continue
			}

			cities, err := u.Source.Cities(ctx, country.Name, region.Name)
			if err != nil {
				return State{}, fmt.Errorf("masjidboard hierarchy: fetch cities for %q / %q: %w", country.Name, region.Name, err)
			}
			region.Cities = mergeLocations(cities, false)

			cityTotal := 0
			for _, city := range region.Cities {
				cityTotal += city.Count
			}
			if cityTotal > region.Count {
				return State{}, fmt.Errorf("masjidboard hierarchy: city counts for %q / %q total %d, exceed expected %d", country.Name, region.Name, cityTotal, region.Count)
			}
			region.UnresolvedCount = region.Count - cityTotal
			country.Regions = append(country.Regions, region)
		}
		state.Countries = append(state.Countries, country)
	}
	return state.Normalized(), nil
}

func mergeLocations(entries []Location, allowBlank bool) []Location {
	counts := make(map[string]int)
	for _, entry := range entries {
		name := strings.TrimSpace(entry.Name)
		if name == "" && !allowBlank {
			continue
		}
		counts[name] += entry.Count
	}
	out := make([]Location, 0, len(counts))
	for name, count := range counts {
		out = append(out, Location{Name: name, Count: count})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func validateCounts(state State) error {
	for _, country := range state.Countries {
		var regionTotal int
		for _, region := range country.Regions {
			regionTotal += region.Count
			var cityTotal int
			for _, city := range region.Cities {
				cityTotal += city.Count
			}
			if cityTotal+region.UnresolvedCount != region.Count {
				return fmt.Errorf("masjidboard hierarchy: city counts for %q / %q total %d plus unresolved %d, expected %d", country.Name, region.Name, cityTotal, region.UnresolvedCount, region.Count)
			}
		}
		if regionTotal != country.Count {
			return fmt.Errorf("masjidboard hierarchy: region counts for %q total %d, expected %d", country.Name, regionTotal, country.Count)
		}
	}
	return nil
}

func sameHierarchy(a, b State) bool {
	if len(a.Countries) != len(b.Countries) {
		return false
	}
	for i := range a.Countries {
		ac, bc := a.Countries[i], b.Countries[i]
		if ac.Name != bc.Name || ac.Count != bc.Count || len(ac.Regions) != len(bc.Regions) {
			return false
		}
		for j := range ac.Regions {
			ar, br := ac.Regions[j], bc.Regions[j]
			if ar.Name != br.Name || ar.Count != br.Count || ar.UnresolvedCount != br.UnresolvedCount || len(ar.Cities) != len(br.Cities) {
				return false
			}
			for k := range ar.Cities {
				if ar.Cities[k] != br.Cities[k] {
					return false
				}
		}
	}
	return true
}
