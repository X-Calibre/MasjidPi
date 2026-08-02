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
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var streams []Stream

	if err := json.Unmarshal(data, &streams); err != nil {
		return nil, err
	}

	return &Store{
		streams: streams,
	}, nil
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
