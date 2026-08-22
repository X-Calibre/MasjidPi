package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/X-Calibre/MasjidPi/backend/internal/atomicfile"
)

type PlaybackState struct {
	StreamID string `json:"stream_id"`
}

type Playback struct {
	path string

	mu       sync.Mutex
	loaded   bool
	streamID string
	hasState bool
}

func NewPlayback(path string) *Playback {
	return &Playback{path: path}
}

func (p *Playback) Load() (string, bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.loaded {
		return p.streamID, p.hasState, nil
	}

	data, err := os.ReadFile(p.path)
	if errors.Is(err, os.ErrNotExist) {
		p.loaded = true
		p.streamID = ""
		p.hasState = false
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	var state PlaybackState
	if err := json.Unmarshal(data, &state); err != nil {
		return "", false, err
	}

	p.loaded = true
	p.streamID = state.StreamID
	p.hasState = state.StreamID != ""
	return p.streamID, p.hasState, nil
}

func (p *Playback) Save(streamID string) error {
	if streamID == "" {
		return p.Clear()
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.loaded {
		data, err := os.ReadFile(p.path)
		if errors.Is(err, os.ErrNotExist) {
			p.loaded = true
			p.streamID = ""
			p.hasState = false
		} else if err != nil {
			return err
		} else {
			var existing PlaybackState
			if err := json.Unmarshal(data, &existing); err != nil {
				return err
			}
			p.loaded = true
			p.streamID = existing.StreamID
			p.hasState = existing.StreamID != ""
		}
	}

	if p.hasState && p.streamID == streamID {
		return nil
	}

	if err := atomicfile.WriteJSON(p.path, PlaybackState{StreamID: streamID}, 0600); err != nil {
		return err
	}

	p.loaded = true
	p.streamID = streamID
	p.hasState = true
	return nil
}

func (p *Playback) Clear() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.loaded && !p.hasState {
		return nil
	}

	if err := os.Remove(p.path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	p.loaded = true
	p.streamID = ""
	p.hasState = false
	return nil
}
