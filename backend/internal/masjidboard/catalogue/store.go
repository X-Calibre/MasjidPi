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

// Store persists the last-known-good MasjidBoard catalogue. It loads once,
// caches state in memory, avoids rewriting identical state, and replaces
// changed state atomically.
type Store struct {
	path   string
	mu     sync.Mutex
	loaded bool
	exists bool
	state  Catalogue
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

// Load returns the cached catalogue after loading it from disk at most once.
// A missing file is treated as an empty catalogue.
func (s *Store) Load() (Catalogue, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.loadLocked(); err != nil {
		return Catalogue{}, err
	}
	return cloneCatalogue(s.state), nil
}

// Save validates and atomically persists a catalogue. Identical state is a
// no-op. The in-memory last-known-good state is changed only after the atomic
// replacement succeeds.
func (s *Store) Save(state Catalogue) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := validateCatalogue(state); err != nil {
		return err
	}
	if err := s.loadLocked(); err != nil {
		return err
	}

	state = cloneCatalogue(state)
	if s.exists && reflect.DeepEqual(s.state, state) {
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

	s.state = state
	s.exists = true
	return nil
}

func (s *Store) loadLocked() error {
	if s.loaded {
		return nil
	}

	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		s.state = Catalogue{}
		s.exists = false
		s.loaded = true
		return nil
	}
	if err != nil {
		return fmt.Errorf("masjidboard catalogue: read catalogue: %w", err)
	}

	var state Catalogue
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("masjidboard catalogue: decode catalogue: %w", err)
	}
	if err := validateCatalogue(state); err != nil {
		return fmt.Errorf("masjidboard catalogue: invalid persisted catalogue: %w", err)
	}

	s.state = cloneCatalogue(state)
	s.exists = true
	s.loaded = true
	return nil
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
