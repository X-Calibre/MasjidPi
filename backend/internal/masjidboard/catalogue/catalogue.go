package catalogue

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Status describes a catalogue record's discovery/availability state.
type Status string

const (
	StatusActive      Status = "active"
	StatusMissing     Status = "missing"
	StatusUnavailable Status = "unavailable"
)

// Record is the provider-neutral identity and discovery metadata for one
// MasjidBoard catalogue entry. Timetable data belongs to the board provider,
// not to the catalogue.
type Record struct {
	ID                 string            `json:"id"`
	Provider           string            `json:"provider"`
	ExternalID         string            `json:"external_id"`
	Name               string            `json:"name"`
	City               string            `json:"city,omitempty"`
	Region             string            `json:"region,omitempty"`
	Country            string            `json:"country,omitempty"`
	TimeZoneOffsetMS   int64             `json:"time_zone_offset_ms"`
	ProviderMetadata   map[string]string `json:"provider_metadata,omitempty"`
	DiscoveredAt       time.Time         `json:"discovered_at"`
	LastSeenAt         time.Time         `json:"last_seen_at"`
	Status             Status            `json:"status"`
}

// Catalogue is one last-known-good discovery snapshot. It is used for
// reconciliation within a single configured location and as the merged view
// presented to callers.
type Catalogue struct {
	RetrievedAt time.Time `json:"retrieved_at"`
	ValidatedAt time.Time `json:"validated_at"`
	Records     []Record  `json:"records"`
}

// Location identifies one configured catalogue source location. Region may be
// blank because MasjidBoard Live supports countries without a province layer.
type Location struct {
	Country string `json:"country"`
	Region  string `json:"region,omitempty"`
	City    string `json:"city"`
}

func (l Location) Normalized() Location {
	return Location{
		Country: strings.TrimSpace(l.Country),
		Region:  strings.TrimSpace(l.Region),
		City:    strings.TrimSpace(l.City),
	}
}

func (l Location) Validate() error {
	l = l.Normalized()
	if l.Country == "" {
		return fmt.Errorf("masjidboard catalogue: location country is required")
	}
	if l.City == "" {
		return fmt.Errorf("masjidboard catalogue: location city is required")
	}
	return nil
}

func (l Location) key() string {
	l = l.Normalized()
	return strings.ToLower(l.Country) + "\x00" + strings.ToLower(l.Region) + "\x00" + strings.ToLower(l.City)
}

// Partition is the independently refreshable last-known-good catalogue for
// one configured discovery location.
type Partition struct {
	Location    Location  `json:"location"`
	RetrievedAt time.Time `json:"retrieved_at"`
	ValidatedAt time.Time `json:"validated_at"`
	Records     []Record  `json:"records"`
}

// State is the disk-first catalogue persistence model. Each configured
// location is retained independently so one failed location refresh cannot
// invalidate successful or previously cached data from another location.
type State struct {
	Partitions []Partition `json:"partitions"`
}

// Merge returns the deduplicated catalogue view across all persisted location
// partitions. A record active in any partition is active in the merged view.
// For duplicate records, the most recently seen copy supplies the metadata.
// The merged timestamps are the oldest non-zero partition timestamps, making
// them conservative freshness indicators for the combined view.
func Merge(state State) Catalogue {
	byID := make(map[string]Record)
	var retrievedAt time.Time
	var validatedAt time.Time

	for _, partition := range state.Partitions {
		if retrievedAt.IsZero() || (!partition.RetrievedAt.IsZero() && partition.RetrievedAt.Before(retrievedAt)) {
			retrievedAt = partition.RetrievedAt
		}
		if validatedAt.IsZero() || (!partition.ValidatedAt.IsZero() && partition.ValidatedAt.Before(validatedAt)) {
			validatedAt = partition.ValidatedAt
		}

		for _, incoming := range partition.Records {
			existing, ok := byID[incoming.ID]
			if !ok {
				byID[incoming.ID] = cloneRecord(incoming)
				continue
			}

			chosen := existing
			if incoming.LastSeenAt.After(existing.LastSeenAt) {
				chosen = cloneRecord(incoming)
			}
			if existing.Status == StatusActive || incoming.Status == StatusActive {
				chosen.Status = StatusActive
			}
			if chosen.DiscoveredAt.IsZero() || (!existing.DiscoveredAt.IsZero() && existing.DiscoveredAt.Before(chosen.DiscoveredAt)) {
				chosen.DiscoveredAt = existing.DiscoveredAt
			}
			if !incoming.DiscoveredAt.IsZero() && (chosen.DiscoveredAt.IsZero() || incoming.DiscoveredAt.Before(chosen.DiscoveredAt)) {
				chosen.DiscoveredAt = incoming.DiscoveredAt
			}
			byID[incoming.ID] = chosen
		}
	}

	records := make([]Record, 0, len(byID))
	for _, record := range byID {
		records = append(records, record)
	}
	sort.Slice(records, func(i, j int) bool { return records[i].ID < records[j].ID })

	return Catalogue{RetrievedAt: retrievedAt, ValidatedAt: validatedAt, Records: records}
}

// ID returns the stable provider-neutral catalogue identity.
func ID(provider, externalID string) (string, error) {
	provider = strings.TrimSpace(provider)
	externalID = strings.TrimSpace(externalID)
	if provider == "" {
		return "", fmt.Errorf("masjidboard catalogue: provider is required")
	}
	if externalID == "" {
		return "", fmt.Errorf("masjidboard catalogue: external ID is required")
	}
	return provider + ":" + externalID, nil
}

// ValidateRecord verifies the minimum identity required for safe reconciliation.
func ValidateRecord(record Record) error {
	wantID, err := ID(record.Provider, record.ExternalID)
	if err != nil {
		return err
	}
	if strings.TrimSpace(record.ID) != wantID {
		return fmt.Errorf("masjidboard catalogue: record ID %q does not match %q", record.ID, wantID)
	}
	if strings.TrimSpace(record.Name) == "" {
		return fmt.Errorf("masjidboard catalogue: name is required for %s", wantID)
	}
	if record.Status != "" && record.Status != StatusActive && record.Status != StatusMissing && record.Status != StatusUnavailable {
		return fmt.Errorf("masjidboard catalogue: invalid status %q for %s", record.Status, wantID)
	}
	return nil
}
