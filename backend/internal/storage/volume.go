package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type VolumeState struct {
	Volume int `json:"volume"`
}

type Volume struct {
	path string
}

func NewVolume(path string) *Volume {
	return &Volume{path: path}
}

func (v *Volume) Load() (int, bool, error) {
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
	if state.Volume < 0 || state.Volume > 125 {
		return 0, false, nil
	}
	return state.Volume, true, nil
}

func (v *Volume) Save(volume int) error {
	if volume < 0 || volume > 125 {
		return errors.New("volume must be between 0 and 125")
	}
	if err := os.MkdirAll(filepath.Dir(v.path), 0755); err != nil {
		return err
	}

	data, err := json.Marshal(VolumeState{Volume: volume})
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
