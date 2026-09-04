package selection

import (
	"fmt"
	"strings"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/catalogue"
)

const (
	MinBoards = 1
	MaxBoards = 3

	DefaultSlideDurationSeconds = 15
	MinSlideDurationSeconds     = 5
	MaxSlideDurationSeconds     = 60

	ThemeEmerald    = "emerald"
	ThemeMidnight   = "midnight"
	ThemeSlate      = "slate"
	ThemeRuby       = "ruby"
	ThemeLight      = "light"
	ThemeBlackWhite = "black-white"
)

var supportedThemes = map[string]struct{}{
	ThemeEmerald: {}, ThemeMidnight: {}, ThemeSlate: {},
	ThemeRuby: {}, ThemeLight: {}, ThemeBlackWhite: {},
}

// Board is the minimal last-known identity required to keep a selected
// MasjidBoard usable without loading or refreshing the full catalogue.
type Board struct {
	CatalogueID        string `json:"catalogue_id"`
	Provider           string `json:"provider"`
	ExternalID         string `json:"external_id"`
	Name               string `json:"name"`
	TimeZoneOffsetMS   int64  `json:"time_zone_offset_ms"`
	ShowDetailedJumuah *bool  `json:"show_detailed_jumuah,omitempty"`
}

func (b Board) ShowDetailedJumuahValue() bool { return enabledByDefault(b.ShowDetailedJumuah) }

// State is the ordered set of boards selected by the user. Order is
// significant and is preserved for display/UI purposes. Display profile is a
// runtime hardware decision and is intentionally not persisted here.
type State struct {
	Boards                 []Board `json:"boards"`
	Theme                  string  `json:"theme,omitempty"`
	SlideDurationSeconds   int     `json:"slide_duration_seconds,omitempty"`
	ShowEconomicIndicators bool    `json:"show_economic_indicators,omitempty"`
	ShowDailyAyah          *bool   `json:"show_daily_ayah,omitempty"`
	ShowDailyHadith        *bool   `json:"show_daily_hadith,omitempty"`
	ShowDailySunnah        *bool   `json:"show_daily_sunnah,omitempty"`
}

func (s State) Configured() bool { return len(s.Boards) > 0 }

func (s State) EffectiveSlideDurationSeconds() int {
	if s.SlideDurationSeconds >= MinSlideDurationSeconds && s.SlideDurationSeconds <= MaxSlideDurationSeconds {
		return s.SlideDurationSeconds
	}
	return DefaultSlideDurationSeconds
}

func (s State) EffectiveTheme() string {
	theme := strings.TrimSpace(strings.ToLower(s.Theme))
	if _, ok := supportedThemes[theme]; ok {
		return theme
	}
	return ThemeEmerald
}

func enabledByDefault(value *bool) bool {
	return value == nil || *value
}

func (s State) ShowDailyAyahValue() bool   { return enabledByDefault(s.ShowDailyAyah) }
func (s State) ShowDailyHadithValue() bool { return enabledByDefault(s.ShowDailyHadith) }
func (s State) ShowDailySunnahValue() bool { return enabledByDefault(s.ShowDailySunnah) }

func (s State) ShowAnyDailyIslamicContent() bool {
	return s.ShowDailyAyahValue() || s.ShowDailyHadithValue() || s.ShowDailySunnahValue()
}

func ThemeSupported(theme string) bool {
	_, ok := supportedThemes[strings.TrimSpace(strings.ToLower(theme))]
	return ok
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

	if state.SlideDurationSeconds != 0 &&
		(state.SlideDurationSeconds < MinSlideDurationSeconds || state.SlideDurationSeconds > MaxSlideDurationSeconds) {
		return fmt.Errorf("masjidboard selection: slide duration must be between %d and %d seconds", MinSlideDurationSeconds, MaxSlideDurationSeconds)
	}
	theme := strings.TrimSpace(strings.ToLower(state.Theme))
	if theme != "" && !ThemeSupported(theme) {
		return fmt.Errorf("masjidboard selection: unsupported display theme %q", state.Theme)
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

func FromCatalogueRecord(record catalogue.Record) (Board, error) {
	if err := catalogue.ValidateRecord(record); err != nil {
		return Board{}, err
	}
	return Board{
		CatalogueID: record.ID, Provider: record.Provider, ExternalID: record.ExternalID,
		Name: record.Name, TimeZoneOffsetMS: record.TimeZoneOffsetMS,
	}, nil
}

func cloneState(state State) State {
	copy := state
	if state.Boards != nil {
		copy.Boards = append([]Board(nil), state.Boards...)
		for index := range copy.Boards {
			copy.Boards[index].ShowDetailedJumuah = cloneBool(state.Boards[index].ShowDetailedJumuah)
		}
	}
	copy.ShowDailyAyah = cloneBool(state.ShowDailyAyah)
	copy.ShowDailyHadith = cloneBool(state.ShowDailyHadith)
	copy.ShowDailySunnah = cloneBool(state.ShowDailySunnah)
	return copy
}

func cloneBool(value *bool) *bool {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
