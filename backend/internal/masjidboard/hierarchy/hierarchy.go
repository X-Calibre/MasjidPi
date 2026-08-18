package hierarchy

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Location is one node in the persisted discovery hierarchy. Count is the
// upstream number of boards beneath that node and is retained as validation
// metadata for later refreshes.
type Location struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// Region is one province/region and the towns/cities currently available
// beneath it.
type Region struct {
	Name   string     `json:"name"`
	Count  int        `json:"count"`
	Cities []Location `json:"cities"`
}

// Country is one active MasjidBoard Live country and its region hierarchy.
type Country struct {
	Name    string   `json:"name"`
	Count   int      `json:"count"`
	Regions []Region `json:"regions"`
}

// State is the complete lightweight location hierarchy. It deliberately does
// not contain board records; those belong to the configured-scope catalogue.
type State struct {
	RetrievedAt time.Time `json:"retrieved_at"`
	ValidatedAt time.Time `json:"validated_at"`
	Countries   []Country `json:"countries"`
}

// Validate checks the minimum invariants required before hierarchy state can
// replace the last-known-good persisted copy.
func (s State) Validate() error {
	if s.RetrievedAt.IsZero() {
		return fmt.Errorf("masjidboard hierarchy: retrieved_at is required")
	}
	if s.ValidatedAt.IsZero() {
		return fmt.Errorf("masjidboard hierarchy: validated_at is required")
	}
	if len(s.Countries) == 0 {
		return fmt.Errorf("masjidboard hierarchy: at least one country is required")
	}
	for _, country := range s.Countries {
		if strings.TrimSpace(country.Name) == "" {
			return fmt.Errorf("masjidboard hierarchy: country name is required")
		}
		if country.Count < 0 {
			return fmt.Errorf("masjidboard hierarchy: invalid count for country %q", country.Name)
		}
		for _, region := range country.Regions {
			if region.Count < 0 {
				return fmt.Errorf("masjidboard hierarchy: invalid count for region %q", region.Name)
			}
			for _, city := range region.Cities {
				if strings.TrimSpace(city.Name) == "" {
					return fmt.Errorf("masjidboard hierarchy: city name is required in %q", region.Name)
				}
				if city.Count < 0 {
					return fmt.Errorf("masjidboard hierarchy: invalid count for city %q", city.Name)
				}
			}
		}
	}
	return nil
}

// Normalized returns deterministic hierarchy state suitable for comparison and
// persistence. Duplicate upstream labels at the same hierarchy level are
// merged by name and their counts are summed. Blank region names are retained
// because MasjidBoard Live has been observed to expose legitimate blank region
// buckets.
func (s State) Normalized() State {
	type regionAccumulator struct {
		count  int
		cities map[string]int
	}
	type countryAccumulator struct {
		count   int
		regions map[string]*regionAccumulator
	}

	countries := make(map[string]*countryAccumulator)
	for _, sourceCountry := range s.Countries {
		countryName := strings.TrimSpace(sourceCountry.Name)
		if countryName == "" {
			continue
		}
		country := countries[countryName]
		if country == nil {
			country = &countryAccumulator{regions: make(map[string]*regionAccumulator)}
			countries[countryName] = country
		}
		country.count += sourceCountry.Count

		for _, sourceRegion := range sourceCountry.Regions {
			regionName := strings.TrimSpace(sourceRegion.Name)
			region := country.regions[regionName]
			if region == nil {
				region = &regionAccumulator{cities: make(map[string]int)}
				country.regions[regionName] = region
			}
			region.count += sourceRegion.Count
			for _, sourceCity := range sourceRegion.Cities {
				cityName := strings.TrimSpace(sourceCity.Name)
				if cityName != "" {
					region.cities[cityName] += sourceCity.Count
				}
			}
		}
	}

	out := State{RetrievedAt: s.RetrievedAt, ValidatedAt: s.ValidatedAt}
	for countryName, sourceCountry := range countries {
		country := Country{Name: countryName, Count: sourceCountry.count}
		for regionName, sourceRegion := range sourceCountry.regions {
			region := Region{Name: regionName, Count: sourceRegion.count}
			for cityName, count := range sourceRegion.cities {
				region.Cities = append(region.Cities, Location{Name: cityName, Count: count})
			}
			sort.Slice(region.Cities, func(i, j int) bool { return region.Cities[i].Name < region.Cities[j].Name })
			country.Regions = append(country.Regions, region)
		}
		sort.Slice(country.Regions, func(i, j int) bool { return country.Regions[i].Name < country.Regions[j].Name })
		out.Countries = append(out.Countries, country)
	}
	sort.Slice(out.Countries, func(i, j int) bool { return out.Countries[i].Name < out.Countries[j].Name })
	return out
}
