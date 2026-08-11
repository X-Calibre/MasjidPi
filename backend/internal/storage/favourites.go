package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type FavouritesState struct {
	IDs []string `json:"ids"`
}

type Favourites struct {
	path string
}

func NewFavourites(path string) *Favourites {
	return &Favourites{path: path}
}

func (f *Favourites) Load() ([]string, error) {
	data, err := os.ReadFile(f.path)
	if errors.Is(err, os.ErrNotExist) {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}

	var state FavouritesState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return state.IDs, nil
}

func (f *Favourites) Save(ids []string) error {
	if err := os.MkdirAll(filepath.Dir(f.path), 0755); err != nil {
		return err
	}

	data, err := json.Marshal(FavouritesState{IDs: ids})
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(f.path), ".favourites-*.tmp")
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

	return os.Rename(tmpName, f.path)
}
