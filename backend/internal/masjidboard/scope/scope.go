package scope

import (
	"fmt"
	"strings"
)

// State is the persisted geographic discovery scope used to build and refresh
// the local MasjidBoard catalogue. The catalogue is intentionally limited to
// the configured town/city rather than mirroring the complete upstream index.
type State struct {
	Country string `json:"country"`
	Region  string `json:"region,omitempty"`
	City    string `json:"city"`
}

// Configured reports whether the state represents a configured discovery
// scope. The zero value is the unconfigured state used before initial setup.
func (s State) Configured() bool {
	return strings.TrimSpace(s.Country) != "" && strings.TrimSpace(s.City) != ""
}

// Validate verifies a configured discovery scope. Region is allowed to be
// empty because the upstream FindMasjid hierarchy contains entries without a
// province/region in some countries.
func (s State) Validate() error {
	country := strings.TrimSpace(s.Country)
	city := strings.TrimSpace(s.City)

	if country == "" && city == "" && strings.TrimSpace(s.Region) == "" {
		return fmt.Errorf("masjidboard scope: unconfigured state is not a valid configured scope")
	}
	if country == "" {
		return fmt.Errorf("masjidboard scope: country is required")
	}
	if city == "" {
		return fmt.Errorf("masjidboard scope: city is required")
	}
	return nil
}

// Normalized returns a trimmed copy suitable for persistence and provider
// requests.
func (s State) Normalized() State {
	return State{
		Country: strings.TrimSpace(s.Country),
		Region:  strings.TrimSpace(s.Region),
		City:    strings.TrimSpace(s.City),
	}
}
