package hierarchy

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"

	"github.com/X-Calibre/MasjidPi/backend/internal/atomicfile"
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

	if err := atomicfile.WriteJSON(s.path, state, 0600); err != nil {
		return fmt.Errorf("masjidboard hierarchy: persist hierarchy: %w", err)
	}
	return nil
}
