package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/X-Calibre/MasjidPi/backend/internal/atomicfile"
)

type AudioDevice struct {
	Name string `json:"name"`
}

type AudioDeviceState struct {
	path string

	mu     sync.RWMutex
	loaded bool
	name   string
	ok     bool
}

func NewAudioDeviceState(path string) *AudioDeviceState {
	return &AudioDeviceState{path: path}
}

func (s *AudioDeviceState) Load() (string, bool, error) {
	s.mu.RLock()
	if s.loaded {
		name, ok := s.name, s.ok
		s.mu.RUnlock()
		return name, ok, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.loaded {
		return s.name, s.ok, nil
	}

	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		s.loaded = true
		s.name = ""
		s.ok = false
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	var state AudioDevice
	if err := json.Unmarshal(data, &state); err != nil {
		return "", false, err
	}

	s.loaded = true
	s.name = state.Name
	s.ok = state.Name != ""
	return s.name, s.ok, nil
}

func (s *AudioDeviceState) Save(name string) error {
	if name == "" {
		return nil
	}

	current, ok, err := s.Load()
	if err != nil {
		return err
	}
	if ok && current == name {
		return nil
	}

	if err := atomicfile.WriteJSON(s.path, AudioDevice{Name: name}, 0644); err != nil {
		return err
	}

	s.mu.Lock()
	s.loaded = true
	s.name = name
	s.ok = true
	s.mu.Unlock()

	return nil
}
