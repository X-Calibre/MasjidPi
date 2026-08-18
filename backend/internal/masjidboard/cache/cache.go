package cache

import (
	"fmt"
	"strings"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
)

// Entry is the last successfully retrieved and validated timetable for one
// selected MasjidBoard. Failed refresh attempts must never replace this data.
type Entry struct {
	CatalogueID string      `json:"catalogue_id"`
	SuccessfulAt time.Time  `json:"successful_at"`
	Board       model.Board `json:"board"`
}

// Validate verifies that a cache entry is safe to use as last-known-good data.
func Validate(entry Entry) error {
	if strings.TrimSpace(entry.CatalogueID) == "" {
		return fmt.Errorf("masjidboard cache: catalogue ID is required")
	}
	if entry.SuccessfulAt.IsZero() {
		return fmt.Errorf("masjidboard cache: successful_at is required")
	}
	if strings.TrimSpace(entry.Board.Identity.ID) == "" {
		return fmt.Errorf("masjidboard cache: board identity ID is required")
	}
	if strings.TrimSpace(entry.Board.Identity.Name) == "" {
		return fmt.Errorf("masjidboard cache: board identity name is required")
	}
	return nil
}
