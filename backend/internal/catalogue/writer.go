package catalogue

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/X-Calibre/MasjidPi/backend/internal/atomicfile"
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

	return atomicfile.Write(filename, data, 0644)
}
