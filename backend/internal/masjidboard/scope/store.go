package scope

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"

	"github.com/X-Calibre/MasjidPi/backend/internal/atomicfile"
)

// Store persists the configured MasjidBoard discovery scope. The state is
// tiny and loaded only when configuration/catalogue maintenance needs it.
type Store struct {
	path string
	mu   sync.Mutex
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

// Load returns the persisted scope. A missing file means MasjidBoard discovery
// has not yet been configured and returns the zero value without error.
func (s *Store) Load() (State, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return State{}, nil
	}
	if err != nil {
		return State{}, fmt.Errorf("masjidboard scope: read scope: %w", err)
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return State{}, fmt.Errorf("masjidboard scope: decode scope: %w", err)
	}
	state = state.Normalized()
	if err := state.Validate(); err != nil {
		return State{}, fmt.Errorf("masjidboard scope: invalid persisted scope: %w", err)
	}
	return state, nil
}

// Save validates and atomically persists a configured scope. Saving an
// unchanged scope is a no-op. An unconfigured zero value is not persisted.
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
		return fmt.Errorf("masjidboard scope: read current scope: %w", err)
	}

	if err := atomicfile.WriteJSON(s.path, state, 0600); err != nil {
		return fmt.Errorf("masjidboard scope: persist scope: %w", err)
	}
	return nil
}
