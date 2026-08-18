package stream

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Store struct {
	mu      sync.RWMutex
	streams []Stream
}

func New(filename string) (*Store, error) {
	store := &Store{}

	if err := store.Reload(filename); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Store) Reload(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var streams []Stream
	if err := json.Unmarshal(data, &streams); err != nil {
		return err
	}

	s.Replace(streams)
	return nil
}

func (s *Store) Replace(streams []Stream) {
	copyOfStreams := append([]Stream(nil), streams...)
	s.mu.Lock()
	s.streams = copyOfStreams
	s.mu.Unlock()
}

func (s *Store) All() []Stream {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Stream(nil), s.streams...)
}

func (s *Store) FindByID(id string) (*Stream, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, item := range s.streams {
		if item.ID == id {
			stream := item
			return &stream, nil
		}
	}

	return nil, fmt.Errorf("stream %q not found", id)
}
