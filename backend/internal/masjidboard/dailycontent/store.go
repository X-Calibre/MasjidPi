package dailycontent

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/X-Calibre/MasjidPi/backend/internal/atomicfile"
)

type Store struct{ Path string }

func (s Store) Load() (*Content, error) {
	data, err := os.ReadFile(s.Path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("daily Islamic content cache: read: %w", err)
	}
	var content Content
	if err := json.Unmarshal(data, &content); err != nil {
		return nil, fmt.Errorf("daily Islamic content cache: decode: %w", err)
	}
	if !content.Valid() {
		return nil, fmt.Errorf("daily Islamic content cache: invalid data")
	}
	return &content, nil
}

func (s Store) Save(content Content) error {
	if !content.Valid() {
		return fmt.Errorf("daily Islamic content cache: refusing to save invalid data")
	}
	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return fmt.Errorf("daily Islamic content cache: encode: %w", err)
	}
	if err := atomicfile.Write(s.Path, data, 0o644); err != nil {
		return fmt.Errorf("daily Islamic content cache: persist: %w", err)
	}
	return nil
}
