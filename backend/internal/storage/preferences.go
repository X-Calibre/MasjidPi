package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type PreferencesState struct {
	LastStreamID string `json:"last_stream_id,omitempty"`
	Autoplay     bool   `json:"autoplay"`
}

type Preferences struct {
	path string
}

func NewPreferences(path string) *Preferences {
	return &Preferences{path: path}
}

func (p *Preferences) Load() (PreferencesState, error) {
	data, err := os.ReadFile(p.path)
	if errors.Is(err, os.ErrNotExist) {
		return PreferencesState{}, nil
	}
	if err != nil {
		return PreferencesState{}, err
	}

	var state PreferencesState
	if err := json.Unmarshal(data, &state); err != nil {
		return PreferencesState{}, err
	}
	return state, nil
}

func (p *Preferences) Save(state PreferencesState) error {
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

	return os.Rename(tmpName, p.path)
}
