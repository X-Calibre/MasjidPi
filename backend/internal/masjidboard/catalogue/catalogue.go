package catalogue

import (
	"fmt"
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

// Catalogue is the complete last-known-good discovery snapshot for the user's
// configured geographic scope, plus records retained by reconciliation when
// they temporarily disappear upstream. It is intentionally not a worldwide
// mirror of the provider catalogue.
type Catalogue struct {
	RetrievedAt time.Time `json:"retrieved_at"`
	ValidatedAt time.Time `json:"validated_at"`
	Records     []Record  `json:"records"`
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
