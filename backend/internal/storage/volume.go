package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type VolumeState struct {
	Volumes map[string]int `json:"volumes,omitempty"`
	// Volume is retained for migration from the pre-Phase-6 single-volume format.
	Volume *int `json:"volume,omitempty"`
}

type Volume struct {
	path string
}

func NewVolume(path string) *Volume {
	return &Volume{path: path}
}

func (v *Volume) Load(device string) (int, bool, error) {
	data, err := os.ReadFile(v.path)
	if errors.Is(err, os.ErrNotExist) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}

	var state VolumeState
	if err := json.Unmarshal(data, &state); err != nil {
		return 0, false, err
	}

	if state.Volumes != nil {
		volume, ok := state.Volumes[device]
		if ok && volume >= 0 && volume <= 100 {
			return volume, true, nil
		}
	}

	if state.Volume != nil && *state.Volume >= 0 && *state.Volume <= 100 {
		return *state.Volume, true, nil
	}

	return 0, false, nil
}

func (v *Volume) Save(device string, volume int) error {
	if device == "" {
		return errors.New("audio device cannot be empty")
	}
	if volume < 0 || volume > 100 {
		return errors.New("volume must be between 0 and 100")
	}
	if err := os.MkdirAll(filepath.Dir(v.path), 0755); err != nil {
		return err
	}

	state := VolumeState{Volumes: map[string]int{}}
	if data, err := os.ReadFile(v.path); err == nil {
		var existing VolumeState
		if json.Unmarshal(data, &existing) == nil {
			state.Volumes = existing.Volumes
			if state.Volumes == nil {
				state.Volumes = make(map[string]int)
			}

			// Avoid rewriting the state file when the persisted value is already
			// the requested value. This is especially important for volume changes
			// because the UI may issue repeated updates while a slider is moved.
			if existingVolume, ok := state.Volumes[device]; ok && existingVolume == volume {
				return nil
			}
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	state.Volumes[device] = volume
	state.Volume = nil

	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(v.path), ".volume-*.tmp")
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

	return os.Rename(tmpName, v.path)
}
