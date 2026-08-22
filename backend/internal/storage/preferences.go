package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/X-Calibre/MasjidPi/backend/internal/atomicfile"
)

type PreferencesState struct {
	LastStreamID string `json:"last_stream_id,omitempty"`
	Autoplay     bool   `json:"autoplay"`
}

type Preferences struct {
	path   string
	mu     sync.Mutex
	loaded bool
	exists bool
	state  PreferencesState
}

func NewPreferences(path string) *Preferences {
	return &Preferences{path: path}
}

func (p *Preferences) Load() (PreferencesState, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if err := p.loadLocked(); err != nil {
		return PreferencesState{}, err
	}
	return p.state, nil
}

func (p *Preferences) Save(state PreferencesState) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if err := p.loadLocked(); err != nil {
		return err
	}
	if p.exists && p.state == state {
		return nil
	}

	if err := atomicfile.WriteJSON(p.path, state, 0600); err != nil {
		return err
	}

	p.state = state
	p.exists = true
	return nil
}

func (p *Preferences) loadLocked() error {
	if p.loaded {
		return nil
	}

	data, err := os.ReadFile(p.path)
	if errors.Is(err, os.ErrNotExist) {
		p.state = PreferencesState{}
		p.exists = false
		p.loaded = true
		return nil
	}
	if err != nil {
		return err
	}

	var state PreferencesState
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	p.state = state
	p.exists = true
	p.loaded = true
	return nil
}
