package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"
)

const unchangedCheckpointInterval = 24 * time.Hour

// Store persists one independent last-known-good cache entry per selected
// board. Cache files are addressed by a hash of the stable catalogue ID so the
// storage layout remains safe across supported operating systems.
type Store struct {
	dir string
	mu  sync.Mutex
}

func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

// Load returns the persisted last-known-good entry for catalogueID. found is
// false when that board has never been cached successfully.
func (s *Store) Load(catalogueID string) (entry Entry, found bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	path, err := s.pathFor(catalogueID)
	if err != nil {
		return Entry{}, false, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Entry{}, false, nil
	}
	if err != nil {
		return Entry{}, false, fmt.Errorf("masjidboard cache: read entry: %w", err)
	}
	if err := json.Unmarshal(data, &entry); err != nil {
		return Entry{}, false, fmt.Errorf("masjidboard cache: decode entry: %w", err)
	}
	if err := Validate(entry); err != nil {
		return Entry{}, false, fmt.Errorf("masjidboard cache: invalid persisted entry: %w", err)
	}
	if entry.CatalogueID != strings.TrimSpace(catalogueID) {
		return Entry{}, false, fmt.Errorf("masjidboard cache: persisted catalogue ID %q does not match requested %q", entry.CatalogueID, strings.TrimSpace(catalogueID))
	}
	return entry, true, nil
}

// Save atomically replaces the last-known-good entry after a successful board
// refresh. If the timetable itself is unchanged, repeated successful refreshes
// are not written every time. A daily checkpoint still advances SuccessfulAt
// so persisted cache freshness remains reasonably representative across
// restarts without creating an SD-card write every refresh interval.
func (s *Store) Save(entry Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := Validate(entry); err != nil {
		return err
	}
	path, err := s.pathFor(entry.CatalogueID)
	if err != nil {
		return err
	}

	if data, err := os.ReadFile(path); err == nil {
		var existing Entry
		if json.Unmarshal(data, &existing) == nil {
			if reflect.DeepEqual(existing, entry) {
				return nil
			}
			if existing.CatalogueID == entry.CatalogueID &&
				reflect.DeepEqual(existing.Board, entry.Board) &&
				entry.SuccessfulAt.Sub(existing.SuccessfulAt) < unchangedCheckpointInterval {
				return nil
			}
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("masjidboard cache: read existing entry: %w", err)
	}

	if err := os.MkdirAll(s.dir, 0755); err != nil {
		return fmt.Errorf("masjidboard cache: create storage directory: %w", err)
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("masjidboard cache: encode entry: %w", err)
	}

	tmp, err := os.CreateTemp(s.dir, ".masjidboard-cache-*.tmp")
	if err != nil {
		return fmt.Errorf("masjidboard cache: create temporary file: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if err := tmp.Chmod(0600); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("masjidboard cache: set temporary file mode: %w", err)
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("masjidboard cache: write temporary file: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("masjidboard cache: sync temporary file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("masjidboard cache: close temporary file: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("masjidboard cache: replace entry: %w", err)
	}
	return nil
}

func (s *Store) pathFor(catalogueID string) (string, error) {
	catalogueID = strings.TrimSpace(catalogueID)
	if catalogueID == "" {
		return "", fmt.Errorf("masjidboard cache: catalogue ID is required")
	}
	sum := sha256.Sum256([]byte(catalogueID))
	name := hex.EncodeToString(sum[:]) + ".json"
	return filepath.Join(s.dir, name), nil
}
