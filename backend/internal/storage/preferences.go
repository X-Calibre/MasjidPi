package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
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

	if err := os.MkdirAll(filepath.Dir(p.path), 0755); err != nil {
		return err
	}

	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(p.path), ".preferences-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if err := tmp.Chmod(0600); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpName, p.path); err != nil {
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
