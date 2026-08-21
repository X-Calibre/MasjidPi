package hierarchy

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

// Store persists the lightweight global discovery hierarchy. State is read on
// demand and changed state is replaced atomically.
type Store struct {
	path string
	mu   sync.Mutex
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

// Load returns the last-known-good hierarchy. A missing file is treated as an
// empty/uninitialised hierarchy.
func (s *Store) Load() (State, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return State{}, nil
	}
	if err != nil {
		return State{}, fmt.Errorf("masjidboard hierarchy: read hierarchy: %w", err)
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return State{}, fmt.Errorf("masjidboard hierarchy: decode hierarchy: %w", err)
	}
	state = state.Normalized()
	if err := state.Validate(); err != nil {
		return State{}, fmt.Errorf("masjidboard hierarchy: invalid persisted hierarchy: %w", err)
	}
	return state, nil
}

// Save validates and atomically persists hierarchy state. Identical state is a
// no-op so weekly freshness checks do not cause needless SD-card writes.
func (s *Store) Save(state State) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	state = state.Normalized()
	if err := state.Validate(); err != nil {
		return err
	}

	if data, err := os.ReadFile(s.path); err == nil {
		var current State
		if json.Unmarshal(data, &current) == nil {
			current = current.Normalized()
			if reflect.DeepEqual(current, state) {
				return nil
			}
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("masjidboard hierarchy: read current hierarchy: %w", err)
	}

	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("masjidboard hierarchy: create storage directory: %w", err)
	}

	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("masjidboard hierarchy: encode hierarchy: %w", err)
	}

	tmp, err := os.CreateTemp(dir, ".masjidboard-hierarchy-*.tmp")
	if err != nil {
		return fmt.Errorf("masjidboard hierarchy: create temporary file: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if err := tmp.Chmod(0600); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("masjidboard hierarchy: set temporary file mode: %w", err)
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("masjidboard hierarchy: write temporary file: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("masjidboard hierarchy: sync temporary file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("masjidboard hierarchy: close temporary file: %w", err)
	}
	if err := os.Rename(tmpName, s.path); err != nil {
		return fmt.Errorf("masjidboard hierarchy: replace hierarchy: %w", err)
	}
	return nil
}
