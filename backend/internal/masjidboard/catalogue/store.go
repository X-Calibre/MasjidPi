package catalogue

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

// Store persists the last-known-good MasjidBoard catalogue. The catalogue is
// disk-first: Load reads it on demand rather than retaining it in memory for
// the lifetime of the process. Save avoids rewriting identical state and
// replaces changed state atomically.
type Store struct {
	path string
	mu   sync.Mutex
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

// Load reads and validates the catalogue from disk on demand. A missing file
// is treated as an empty catalogue.
func (s *Store) Load() (Catalogue, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.loadLocked()
}

// Save validates and atomically persists a catalogue. If the currently
// persisted catalogue is identical, Save is a no-op. A failed save leaves the
// existing on-disk last-known-good catalogue unchanged.
func (s *Store) Save(state Catalogue) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := validateCatalogue(state); err != nil {
		return err
	}

	current, exists, err := s.loadExistingLocked()
	if err != nil {
		return err
	}

	state = cloneCatalogue(state)
	if exists && reflect.DeepEqual(current, state) {
		return nil
	}

	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("masjidboard catalogue: create storage directory: %w", err)
	}

	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("masjidboard catalogue: encode catalogue: %w", err)
	}

	tmp, err := os.CreateTemp(dir, ".masjidboard-catalogue-*.tmp")
	if err != nil {
		return fmt.Errorf("masjidboard catalogue: create temporary file: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if err := tmp.Chmod(0600); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("masjidboard catalogue: set temporary file mode: %w", err)
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("masjidboard catalogue: write temporary file: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("masjidboard catalogue: sync temporary file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("masjidboard catalogue: close temporary file: %w", err)
	}
	if err := os.Rename(tmpName, s.path); err != nil {
		return fmt.Errorf("masjidboard catalogue: replace catalogue: %w", err)
	}

	return nil
}

func (s *Store) loadLocked() (Catalogue, error) {
	state, exists, err := s.loadExistingLocked()
	if err != nil {
		return Catalogue{}, err
	}
	if !exists {
		return Catalogue{}, nil
	}
	return state, nil
}

func (s *Store) loadExistingLocked() (Catalogue, bool, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return Catalogue{}, false, nil
	}
	if err != nil {
		return Catalogue{}, false, fmt.Errorf("masjidboard catalogue: read catalogue: %w", err)
	}

	var state Catalogue
	if err := json.Unmarshal(data, &state); err != nil {
		return Catalogue{}, false, fmt.Errorf("masjidboard catalogue: decode catalogue: %w", err)
	}
	if err := validateCatalogue(state); err != nil {
		return Catalogue{}, false, fmt.Errorf("masjidboard catalogue: invalid persisted catalogue: %w", err)
	}

	return cloneCatalogue(state), true, nil
}

func validateCatalogue(state Catalogue) error {
	if state.RetrievedAt.IsZero() {
		return fmt.Errorf("masjidboard catalogue: retrieved_at is required")
	}
	if state.ValidatedAt.IsZero() {
		return fmt.Errorf("masjidboard catalogue: validated_at is required")
	}

	seen := make(map[string]struct{}, len(state.Records))
	for _, record := range state.Records {
		if err := ValidateRecord(record); err != nil {
			return err
		}
		if _, exists := seen[record.ID]; exists {
			return fmt.Errorf("masjidboard catalogue: duplicate record %q", record.ID)
		}
		seen[record.ID] = struct{}{}
	}
	return nil
}

func cloneCatalogue(state Catalogue) Catalogue {
	copy := state
	if state.Records != nil {
		copy.Records = make([]Record, len(state.Records))
		for i, record := range state.Records {
			copy.Records[i] = cloneRecord(record)
		}
	}
	return copy
}
