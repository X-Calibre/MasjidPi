package catalogue

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
)

func WriteCatalogue(filename string, streams []stream.Stream) error {
	data, err := json.MarshalIndent(streams, "", "    ")
	if err != nil {
		return err
	}

	if existing, err := os.ReadFile(filename); err == nil {
		if bytes.Equal(existing, data) {
			return nil
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(filename), ".catalogue-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if err := tmp.Chmod(0644); err != nil {
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

	return os.Rename(tmpName, filename)
}
