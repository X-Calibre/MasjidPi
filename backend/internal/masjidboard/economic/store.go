package economic

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o755); err != nil {
		return fmt.Errorf("economic indicators cache: create directory: %w", err)
	}
	data, err := json.MarshalIndent(indicators, "", "  ")
	if err != nil {
		return fmt.Errorf("economic indicators cache: encode: %w", err)
	}
	temporary, err := os.CreateTemp(filepath.Dir(s.Path), ".economic-indicators-*.tmp")
	if err != nil {
		return fmt.Errorf("economic indicators cache: create temporary file: %w", err)
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	if err := temporary.Chmod(0o644); err != nil {
		temporary.Close()
		return fmt.Errorf("economic indicators cache: set permissions: %w", err)
	}
	if _, err := temporary.Write(data); err != nil {
		temporary.Close()
		return fmt.Errorf("economic indicators cache: write: %w", err)
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return fmt.Errorf("economic indicators cache: sync: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("economic indicators cache: close: %w", err)
	}
	if err := os.Rename(temporaryPath, s.Path); err != nil {
		return fmt.Errorf("economic indicators cache: replace: %w", err)
	}
	return nil
}
