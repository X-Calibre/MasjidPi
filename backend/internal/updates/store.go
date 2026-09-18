package updates

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/X-Calibre/MasjidPi/backend/internal/atomicfile"
)

// Store durably persists the small update-orchestration state file.
type Store struct {
	path string
	mu   sync.Mutex
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

// Load returns default state when the appliance has never checked for an
// update.
func (s *Store) Load() (State, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return DefaultState(""), nil
	}
	if err != nil {
		return State{}, fmt.Errorf(
			"updates: read state: %w",
			err,
		)
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return State{}, fmt.Errorf(
			"updates: decode state: %w",
			err,
		)
	}

	state = state.normalized()
	if err := state.Validate(); err != nil {
		return State{}, fmt.Errorf(
			"updates: invalid persisted state: %w",
			err,
		)
	}

	return cloneState(state), nil
}

// Save validates and atomically replaces the update state using private file
// permissions because future state may contain local download paths.
func (s *Store) Save(state State) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	state = state.normalized()
	if err := state.Validate(); err != nil {
		return err
	}

	if err := atomicfile.WriteJSON(
		s.path,
		state,
		0600,
	); err != nil {
		return fmt.Errorf(
			"updates: persist state: %w",
			err,
		)
	}

	return nil
}
