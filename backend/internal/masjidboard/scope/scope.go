package scope

import (
	"fmt"
	"strings"
)

const (
	MinLocations = 1
	MaxLocations = 3
)

// Location is one geographic discovery location used to build and refresh the
// local MasjidBoard catalogue. Region may be blank because the upstream
// FindMasjid hierarchy contains entries without a province/region in some
// countries.
type Location struct {
	Country string `json:"country"`
	Region  string `json:"region,omitempty"`
	City    string `json:"city"`
}

// State is the persisted geographic discovery scope. A configured MasjidBoard
// installation has between one and three ordered locations. The catalogue is
// the union of boards discovered from those locations rather than a worldwide
// mirror of the upstream directory.
type State struct {
	Locations []Location `json:"locations"`
}

// Configured reports whether the state represents a complete configured
// discovery scope. The zero value is the unconfigured state used before
// initial setup.
func (s State) Configured() bool {
	if len(s.Locations) < MinLocations || len(s.Locations) > MaxLocations {
		return false
	}
	return s.Validate() == nil
}

// Validate verifies a configured discovery scope. One to three unique
// locations are required. Country and city are required for every location;
// region may be blank.
func (s State) Validate() error {
	if len(s.Locations) < MinLocations {
		return fmt.Errorf("masjidboard scope: at least %d location is required", MinLocations)
	}
	if len(s.Locations) > MaxLocations {
		return fmt.Errorf("masjidboard scope: at most %d locations are allowed", MaxLocations)
	}

	seen := make(map[string]struct{}, len(s.Locations))
	for i, location := range s.Locations {
		location = location.Normalized()
		if err := location.Validate(); err != nil {
			return fmt.Errorf("masjidboard scope: location %d: %w", i+1, err)
		}

		key := location.key()
		if _, exists := seen[key]; exists {
			return fmt.Errorf("masjidboard scope: duplicate location %q", location.DisplayName())
		}
		seen[key] = struct{}{}
	}
	return nil
}

// Normalized returns a trimmed copy suitable for persistence and provider
// requests. User order is retained because the configuration UI may use it as
// the preferred browsing order.
func (s State) Normalized() State {
	out := State{Locations: make([]Location, len(s.Locations))}
	for i, location := range s.Locations {
		out.Locations[i] = location.Normalized()
	}
	return out
}

func (l Location) Validate() error {
	if strings.TrimSpace(l.Country) == "" {
		return fmt.Errorf("country is required")
	}
	if strings.TrimSpace(l.City) == "" {
		return fmt.Errorf("city is required")
	}
	return nil
}

func (l Location) Normalized() Location {
	return Location{
		Country: strings.TrimSpace(l.Country),
		Region:  strings.TrimSpace(l.Region),
		City:    strings.TrimSpace(l.City),
	}
}

// DisplayName returns a human-readable location label for diagnostics.
func (l Location) DisplayName() string {
	l = l.Normalized()
	parts := []string{l.Country}
	if l.Region != "" {
		parts = append(parts, l.Region)
	}
	parts = append(parts, l.City)
	return strings.Join(parts, " / ")
}

func (l Location) key() string {
	l = l.Normalized()
	return strings.ToLower(l.Country) + "\x00" + strings.ToLower(l.Region) + "\x00" + strings.ToLower(l.City)
}
