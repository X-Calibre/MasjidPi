package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/X-Calibre/MasjidPi/backend/internal/atomicfile"
)

type VolumeState struct {
	Volumes map[string]int `json:"volumes,omitempty"`
	// Volume is retained for migration from the pre-Phase-6 single-volume format.
	Volume *int `json:"volume,omitempty"`
}

type Volume struct {
	path   string
	mu     sync.Mutex
	loaded bool
	state  VolumeState
}

func NewVolume(path string) *Volume {
	return &Volume{path: path}
}

func (v *Volume) Load(device string) (int, bool, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if err := v.loadLocked(); err != nil {
		return 0, false, err
	}
	return volumeFromState(v.state, device)
}

func (v *Volume) Save(device string, volume int) error {
	if device == "" {
		return errors.New("audio device cannot be empty")
	}
	if volume < 0 || volume > 100 {
		return errors.New("volume must be between 0 and 100")
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	if err := v.loadLocked(); err != nil {
		return err
	}
	if existingVolume, ok := v.state.Volumes[device]; ok && existingVolume == volume {
		return nil
	}

	if v.state.Volumes == nil {
		v.state.Volumes = make(map[string]int)
	}
	v.state.Volumes[device] = volume
	v.state.Volume = nil

	if err := atomicfile.WriteJSON(v.path, v.state, 0600); err != nil {
		return err
	}

	return nil
}

func (v *Volume) loadLocked() error {
	if v.loaded {
		return nil
	}

	state := VolumeState{Volumes: make(map[string]int)}
	data, err := os.ReadFile(v.path)
	if errors.Is(err, os.ErrNotExist) {
		v.state = state
		v.loaded = true
		return nil
	}
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}
	if state.Volumes == nil {
		state.Volumes = make(map[string]int)
	}
	v.state = state
	v.loaded = true
	return nil
}

func volumeFromState(state VolumeState, device string) (int, bool, error) {
	if state.Volumes != nil {
		if volume, ok := state.Volumes[device]; ok && volume >= 0 && volume <= 100 {
			return volume, true, nil
		}
	}
	if state.Volume != nil && *state.Volume >= 0 && *state.Volume <= 100 {
		return *state.Volume, true, nil
	}
	return 0, false, nil
}
