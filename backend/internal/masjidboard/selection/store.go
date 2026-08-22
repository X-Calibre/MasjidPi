package selection

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"

	"github.com/X-Calibre/MasjidPi/backend/internal/atomicfile"
)

// Store persists the small runtime selection state. Unlike the full catalogue,
// configured selection state is loaded once and retained in memory because
// normal runtime needs it continuously.
type Store struct {
	path   string
	mu     sync.Mutex
	loaded bool
	exists bool
	state  State
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

// Load reads the selection at most once and returns a defensive copy.
// A missing file returns the zero State, which represents an unconfigured
// MasjidBoard installation rather than a valid configured selection.
func (s *Store) Load() (State, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.loadLocked(); err != nil {
		return State{}, err
	}
	return cloneState(s.state), nil
}

// Save validates and atomically persists a configured selection. A configured
// selection must contain one to three boards; an empty selection is rejected.
// Identical state is a no-op. Selection order is significant and therefore
// participates in equality.
func (s *Store) Save(state State) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := Validate(state); err != nil {
		return err
	}
	if err := s.loadLocked(); err != nil {
		return err
	}

	state = cloneState(state)
	if s.exists && reflect.DeepEqual(s.state, state) {
		return nil
	}

	if err := atomicfile.WriteJSON(s.path, state, 0600); err != nil {
		return fmt.Errorf("masjidboard selection: persist state: %w", err)
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
		s.state = State{}
		s.exists = false
		s.loaded = true
		return nil
	}
	if err != nil {
		return fmt.Errorf("masjidboard selection: read state: %w", err)
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("masjidboard selection: decode state: %w", err)
	}
	if err := Validate(state); err != nil {
		return fmt.Errorf("masjidboard selection: invalid persisted state: %w", err)
	}

	s.state = cloneState(state)
	s.exists = true
	s.loaded = true
	return nil
}
