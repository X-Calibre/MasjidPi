package masjidboardlive

import (
	"context"
	"time"

	masjidboardcatalogue "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/catalogue"
)

const catalogueProvider = "masjidboardlive"

// CatalogueSource adapts FindMasjid city-level discovery to the provider-neutral
// partitioned catalogue updater.
type CatalogueSource struct {
	Client DiscoveryClient
}

func (s CatalogueSource) Fetch(ctx context.Context, location masjidboardcatalogue.Location, now time.Time) (masjidboardcatalogue.Catalogue, error) {
	result, err := s.Client.SearchMasjids(ctx, MasjidSearch{
		Search:   location.City,
		Country:  location.Country,
		Province: location.Region,
		City:     location.City,
	})
	if err != nil {
		return masjidboardcatalogue.Catalogue{}, err
	}

	records := make([]masjidboardcatalogue.Record, 0, len(result.Entries))
	for _, entry := range result.Entries {
		id, err := masjidboardcatalogue.ID(catalogueProvider, entry.WebURL)
		if err != nil {
			return masjidboardcatalogue.Catalogue{}, err
		}
		metadata := map[string]string{}
		if entry.MBLID != "" {
			metadata["mbl_id"] = entry.MBLID
		}
		if entry.LastUpdated != "" {
			metadata["last_updated"] = entry.LastUpdated
		}
		records = append(records, masjidboardcatalogue.Record{
			ID:               id,
			Provider:         catalogueProvider,
			ExternalID:       entry.WebURL,
			Name:             entry.Name,
			City:             entry.City,
			Region:           location.Region,
			Country:          location.Country,
			TimeZoneOffsetMS: entry.TimeZoneOffsetMS,
			ProviderMetadata: metadata,
			Status:           masjidboardcatalogue.StatusActive,
		})
	}

	return masjidboardcatalogue.Catalogue{
		RetrievedAt: now,
		ValidatedAt: now,
		Records:     records,
	}, nil
}
