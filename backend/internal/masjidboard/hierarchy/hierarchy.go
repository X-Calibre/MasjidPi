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
	out := State{RetrievedAt: s.RetrievedAt, ValidatedAt: s.ValidatedAt}
	countries := make(map[string]*Country)

	for _, sourceCountry := range s.Countries {
		countryName := strings.TrimSpace(sourceCountry.Name)
		if countryName == "" {
			continue
		}
		country := countries[countryName]
		if country == nil {
			country = &Country{Name: countryName}
			countries[countryName] = country
		}
		country.Count += sourceCountry.Count

		regionMap := make(map[string]*Region)
		for i := range country.Regions {
			regionMap[country.Regions[i].Name] = &country.Regions[i]
		}
		for _, sourceRegion := range sourceCountry.Regions {
			regionName := strings.TrimSpace(sourceRegion.Name)
			region := regionMap[regionName]
			if region == nil {
				country.Regions = append(country.Regions, Region{Name: regionName})
				region = &country.Regions[len(country.Regions)-1]
				regionMap[regionName] = region
			}
			region.Count += sourceRegion.Count

			cityCounts := make(map[string]int)
			for _, city := range region.Cities {
				cityCounts[strings.TrimSpace(city.Name)] += city.Count
			}
			for _, city := range sourceRegion.Cities {
				name := strings.TrimSpace(city.Name)
				if name != "" {
					cityCounts[name] += city.Count
				}
			}
			region.Cities = region.Cities[:0]
			for name, count := range cityCounts {
				if name != "" {
					region.Cities = append(region.Cities, Location{Name: name, Count: count})
				}
			}
			sort.Slice(region.Cities, func(i, j int) bool { return region.Cities[i].Name < region.Cities[j].Name })
		}
		sort.Slice(country.Regions, func(i, j int) bool { return country.Regions[i].Name < country.Regions[j].Name })
	}

	for _, country := range countries {
		out.Countries = append(out.Countries, *country)
	}
	sort.Slice(out.Countries, func(i, j int) bool { return out.Countries[i].Name < out.Countries[j].Name })
	return out
}
