package stream

import (
	"encoding/json"
	"fmt"
	"os"
)

type Store struct {
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

	s.streams = streams

	return nil
}

func (s *Store) All() []Stream {
	return s.streams
}

func (s *Store) FindByID(id string) (*Stream, error) {
	for _, stream := range s.streams {
		if stream.ID == id {
			return &stream, nil
		}
	}

	return nil, fmt.Errorf("stream %q not found", id)
}
