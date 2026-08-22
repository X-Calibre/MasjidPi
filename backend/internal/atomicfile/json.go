package atomicfile

import (
	"encoding/json"
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
	return os.Rename(tmpName, path)
}
