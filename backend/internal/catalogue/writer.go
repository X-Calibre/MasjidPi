package catalogue

import (
	"encoding/json"
	"os"

	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
)

func WriteCatalogue(filename string, streams []stream.Stream) error {

	data, err := json.MarshalIndent(streams, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
