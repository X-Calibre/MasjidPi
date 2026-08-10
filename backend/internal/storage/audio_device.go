package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type AudioDevice struct {
	Name string `json:"name"`
}

type AudioDeviceState struct {
	path string
}

func NewAudioDeviceState(path string) *AudioDeviceState {
	return &AudioDeviceState{path: path}
}

func (s *AudioDeviceState) Load() (string, bool, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	var state AudioDevice
	if err := json.Unmarshal(data, &state); err != nil {
		return "", false, err
	}
	if state.Name == "" {
		return "", false, nil
	}
	return state.Name, true, nil
}

func (s *AudioDeviceState) Save(name string) error {
	if name == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0755); err != nil {
		return err
	}

	data, err := json.Marshal(AudioDevice{Name: name})
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}
