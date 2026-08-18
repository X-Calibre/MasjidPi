package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

type FavouritesState struct {
	IDs []string `json:"ids"`
}

type Favourites struct {
	path   string
	mu     sync.Mutex
	loaded bool
	ids    []string
}

func NewFavourites(path string) *Favourites {
	return &Favourites{path: path}
}

func (f *Favourites) Load() ([]string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if err := f.loadLocked(); err != nil {
		return nil, err
	}
	return cloneStrings(f.ids), nil
}

func (f *Favourites) Save(ids []string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if err := f.loadLocked(); err != nil {
		return err
	}
	if reflect.DeepEqual(f.ids, ids) {
		return nil
	}

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
	if err := os.Rename(tmpName, f.path); err != nil {
		return err
	}

	f.ids = cloneStrings(ids)
	f.loaded = true
	return nil
}

func (f *Favourites) loadLocked() error {
	if f.loaded {
		return nil
	}

	data, err := os.ReadFile(f.path)
	if errors.Is(err, os.ErrNotExist) {
		f.ids = []string{}
		f.loaded = true
		return nil
	}
	if err != nil {
		return err
	}

	var state FavouritesState
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}
	if state.IDs == nil {
		state.IDs = []string{}
	}
	f.ids = cloneStrings(state.IDs)
	f.loaded = true
	return nil
}

func cloneStrings(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}
	return append([]string(nil), values...)
}
