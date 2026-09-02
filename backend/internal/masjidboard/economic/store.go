package economic

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/X-Calibre/MasjidPi/backend/internal/atomicfile"
)

type Store struct{ Path string }

func (s Store) Load() (*Indicators, error) {
	data, err := os.ReadFile(s.Path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("economic indicators cache: read: %w", err)
	}
	var indicators Indicators
	if err := json.Unmarshal(data, &indicators); err != nil {
		return nil, fmt.Errorf("economic indicators cache: decode: %w", err)
	}
	if !indicators.Valid() {
		return nil, fmt.Errorf("economic indicators cache: invalid data")
	}
	return &indicators, nil
}

func (s Store) Save(indicators Indicators) error {
	if !indicators.Valid() {
		return fmt.Errorf("economic indicators cache: refusing to save invalid data")
	}
	data, err := json.MarshalIndent(indicators, "", "  ")
	if err != nil {
		return fmt.Errorf("economic indicators cache: encode: %w", err)
	}
	if err := atomicfile.Write(s.Path, data, 0o644); err != nil {
		return fmt.Errorf("economic indicators cache: persist: %w", err)
	}
	return nil
}
