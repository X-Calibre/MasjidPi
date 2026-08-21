package selection

import (
	"fmt"
	"strings"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/catalogue"
)

const (
	MinBoards = 1
	MaxBoards = 3
)

// Board is the minimal last-known identity required to keep a selected
// MasjidBoard usable without loading or refreshing the full catalogue.
type Board struct {
	CatalogueID      string `json:"catalogue_id"`
	Provider         string `json:"provider"`
	ExternalID       string `json:"external_id"`
	Name             string `json:"name"`
	TimeZoneOffsetMS int64  `json:"time_zone_offset_ms"`
}

// State is the ordered set of boards selected by the user. Order is
// significant and is preserved for display/UI purposes.
//
// The zero value represents an unconfigured MasjidBoard installation and is
// returned when no persisted selection exists. A configured State must contain
// between MinBoards and MaxBoards boards.
type State struct {
	Boards []Board `json:"boards"`
}

// Configured reports whether the state contains a user configuration.
func (s State) Configured() bool {
	return len(s.Boards) > 0
}

// Validate verifies a configured selection. An empty State is the internal
// unconfigured state and is not a valid persisted user selection.
func Validate(state State) error {
	if len(state.Boards) < MinBoards {
		return fmt.Errorf("masjidboard selection: at least %d board must be selected", MinBoards)
	}
	if len(state.Boards) > MaxBoards {
		return fmt.Errorf("masjidboard selection: %d boards selected; maximum is %d", len(state.Boards), MaxBoards)
	}

	seen := make(map[string]struct{}, len(state.Boards))
	for i, board := range state.Boards {
		provider := strings.TrimSpace(board.Provider)
		externalID := strings.TrimSpace(board.ExternalID)
		name := strings.TrimSpace(board.Name)
		if provider == "" {
			return fmt.Errorf("masjidboard selection: board %d provider is required", i+1)
		}
		if externalID == "" {
			return fmt.Errorf("masjidboard selection: board %d external ID is required", i+1)
		}
		if name == "" {
			return fmt.Errorf("masjidboard selection: board %d name is required", i+1)
		}

		wantID, err := catalogue.ID(provider, externalID)
		if err != nil {
			return fmt.Errorf("masjidboard selection: board %d: %w", i+1, err)
		}
		if strings.TrimSpace(board.CatalogueID) != wantID {
			return fmt.Errorf("masjidboard selection: board %d catalogue ID %q does not match %q", i+1, board.CatalogueID, wantID)
		}
		if _, exists := seen[wantID]; exists {
			return fmt.Errorf("masjidboard selection: duplicate board %q", wantID)
		}
		seen[wantID] = struct{}{}
	}
	return nil
}

// FromCatalogueRecord creates a persisted selection entry from a catalogue
// record while retaining only the fields required for runtime operation.
func FromCatalogueRecord(record catalogue.Record) (Board, error) {
	if err := catalogue.ValidateRecord(record); err != nil {
		return Board{}, err
	}
	return Board{
		CatalogueID:      record.ID,
		Provider:         record.Provider,
		ExternalID:       record.ExternalID,
		Name:             record.Name,
		TimeZoneOffsetMS: record.TimeZoneOffsetMS,
	}, nil
}

func cloneState(state State) State {
	copy := state
	if state.Boards != nil {
		copy.Boards = append([]Board(nil), state.Boards...)
	}
	return copy
}
