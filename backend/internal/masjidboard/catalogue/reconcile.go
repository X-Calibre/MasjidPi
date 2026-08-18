package catalogue

import (
	"fmt"
	"sort"
)

// Reconcile combines a validated discovery candidate with the current
// last-known-good catalogue. Records present in the candidate become active;
// records absent from the candidate are retained as missing.
func Reconcile(current, candidate Catalogue) (Catalogue, error) {
	if candidate.ValidatedAt.IsZero() {
		return Catalogue{}, fmt.Errorf("masjidboard catalogue: candidate validated_at is required")
	}
	if candidate.RetrievedAt.IsZero() {
		return Catalogue{}, fmt.Errorf("masjidboard catalogue: candidate retrieved_at is required")
	}

	currentByID := make(map[string]Record, len(current.Records))
	for _, record := range current.Records {
		if err := ValidateRecord(record); err != nil {
			return Catalogue{}, fmt.Errorf("current record: %w", err)
		}
		if _, exists := currentByID[record.ID]; exists {
			return Catalogue{}, fmt.Errorf("masjidboard catalogue: duplicate current record %q", record.ID)
		}
		currentByID[record.ID] = cloneRecord(record)
	}

	seen := make(map[string]struct{}, len(candidate.Records))
	reconciled := make([]Record, 0, len(current.Records)+len(candidate.Records))

	for _, incoming := range candidate.Records {
		if err := ValidateRecord(incoming); err != nil {
			return Catalogue{}, fmt.Errorf("candidate record: %w", err)
		}
		if _, duplicate := seen[incoming.ID]; duplicate {
			return Catalogue{}, fmt.Errorf("masjidboard catalogue: duplicate candidate record %q", incoming.ID)
		}
		seen[incoming.ID] = struct{}{}

		record := cloneRecord(incoming)
		record.Status = StatusActive
		record.LastSeenAt = candidate.ValidatedAt

		if existing, ok := currentByID[incoming.ID]; ok {
			record.DiscoveredAt = existing.DiscoveredAt
			if record.DiscoveredAt.IsZero() {
				record.DiscoveredAt = candidate.ValidatedAt
			}
		} else {
			record.DiscoveredAt = candidate.ValidatedAt
		}

		reconciled = append(reconciled, record)
	}

	for id, existing := range currentByID {
		if _, ok := seen[id]; ok {
			continue
		}
		record := cloneRecord(existing)
		record.Status = StatusMissing
		reconciled = append(reconciled, record)
	}

	sort.Slice(reconciled, func(i, j int) bool {
		return reconciled[i].ID < reconciled[j].ID
	})

	return Catalogue{
		RetrievedAt: candidate.RetrievedAt,
		ValidatedAt: candidate.ValidatedAt,
		Records:     reconciled,
	}, nil
}

func cloneRecord(record Record) Record {
	copy := record
	if record.ProviderMetadata != nil {
		copy.ProviderMetadata = make(map[string]string, len(record.ProviderMetadata))
		for key, value := range record.ProviderMetadata {
			copy.ProviderMetadata[key] = value
		}
	}
	return copy
}
