package masjidboardlive

import (
	"context"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/hierarchy"
)

// HierarchySource adapts the MasjidBoard Live discovery client to the
// provider-neutral hierarchy refresh service.
type HierarchySource struct {
	Client DiscoveryClient
}

func (s HierarchySource) Countries(ctx context.Context) ([]hierarchy.Location, error) {
	entries, err := s.Client.Countries(ctx)
	if err != nil {
		return nil, err
	}
	return hierarchyLocations(entries), nil
}

func (s HierarchySource) Regions(ctx context.Context, country string) ([]hierarchy.Location, error) {
	entries, err := s.Client.Regions(ctx, country)
	if err != nil {
		return nil, err
	}
	return hierarchyLocations(entries), nil
}

func (s HierarchySource) Cities(ctx context.Context, country, region string) ([]hierarchy.Location, error) {
	entries, err := s.Client.Cities(ctx, country, region)
	if err != nil {
		return nil, err
	}
	return hierarchyLocations(entries), nil
}

func hierarchyLocations(entries []HierarchyEntry) []hierarchy.Location {
	out := make([]hierarchy.Location, len(entries))
	for i, entry := range entries {
		out[i] = hierarchy.Location{Name: entry.Name, Count: entry.Count}
	}
	return out
}
