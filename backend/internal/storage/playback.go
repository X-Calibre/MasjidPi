package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type PlaybackState struct {
	StreamID string `json:"stream_id"`
}

type Playback struct {
	path string
}

func NewPlayback(path string) *Playback {
	return &Playback{path: path}
}

func (p *Playback) Load() (string, bool, error) {
	data, err := os.ReadFile(p.path)
	if errors.Is(err, os.ErrNotExist) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	var state PlaybackState
	if err := json.Unmarshal(data, &state); err != nil {
		return "", false, err
	}
	if state.StreamID == "" {
		return "", false, nil
	}

	return state.StreamID, true, nil
}

func (p *Playback) Save(streamID string) error {
	if streamID == "" {
		return p.Clear()
	}

	if err := os.MkdirAll(filepath.Dir(p.path), 0755); err != nil {
		return err
	}

	data, err := json.Marshal(PlaybackState{StreamID: streamID})
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(p.path), ".playback-*.tmp")
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

func (p *Playback) Clear() error {
	if err := os.Remove(p.path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
