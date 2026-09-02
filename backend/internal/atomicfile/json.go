package atomicfile

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// WriteJSON replaces path atomically with the JSON representation of value.
func WriteJSON(path string, value any, mode fs.FileMode) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return Write(path, data, mode)
}

// Write replaces path atomically with data.
func Write(path string, data []byte, mode fs.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, ".masjidpi-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if err := tmp.Chmod(mode); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		return err
	}
	return syncDirectory(dir)
}

// Remove deletes path and durably records the directory change. It reports
// whether a file was removed.
func Remove(path string) (bool, error) {
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := syncDirectory(filepath.Dir(path)); err != nil {
		return true, err
	}
	return true, nil
}

func syncDirectory(path string) error {
	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	if err := dir.Sync(); err != nil {
		_ = dir.Close()
		return err
	}
	return dir.Close()
}
