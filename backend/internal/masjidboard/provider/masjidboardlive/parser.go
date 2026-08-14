package masjidboardlive

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
)

// MasjidBoard Live returns 29 positional top-level rows. Keep these indexes
// inside this package so the rest of MasjidBoard never depends on the source
// layout.
const (
	rowJumuah = 1
	rowSalah  = 3
	rowMasjid = 6
)

// Parse normalises the verified core portion of a MasjidBoard Live response.
// The provider deliberately fails when the required identity, timezone, or
// five-prayer data cannot be established. Optional rows are parsed only when
// their values are required by the current model and can be identified
// without guessing at undocumented semantics.
func Parse(rows []json.RawMessage, boardID string, now time.Time) (model.Board, error) {
	if boardID == "" {
		return model.Board{}, fmt.Errorf("masjidboardlive: board ID is required")
	}
	if len(rows) < 29 {
		return model.Board{}, fmt.Errorf("masjidboardlive: expected 29 rows, got %d", len(rows))
	}

	masjid, err := parseMasjidRow(rows[rowMasjid])
	if err != nil {
		return model.Board{}, err
	}
	if masjid.Timezone == "" {
		return model.Board{}, fmt.Errorf("masjidboardlive: masjid row has no timezone")
	}

	loc, err := time.LoadLocation(masjid.Timezone)
	if err != nil {
		return model.Board{}, fmt.Errorf("masjidboardlive: load timezone %q: %w", masjid.Timezone, err)
	}

	localNow := now.In(loc)
	prayers, err := parseSalahRow(rows[rowSalah], localNow, loc)
	if err != nil {
		return model.Board{}, err
	}

	return model.Board{
		Identity: model.Identity{
			SourceBoardID: boardID,
			MasjidID:      boardID,
			EnglishName:   masjid.Name1,
			ArabicName:    masjid.Name2,
		},
		DateContext: model.DateContext{
			GregorianDate: dateOnly(localNow, loc),
			Timezone:      masjid.Timezone,
		},
		PrayerTimes: prayers,
	}, nil
}

type masjidRow struct {
	Name1    string
	Name2    string
	URL      string
	Timezone string
}

func parseMasjidRow(raw json.RawMessage) (masjidRow, error) {
	values, err := rowValues(raw)
	if err != nil {
		return masjidRow{}, fmt.Errorf("masjidboardlive: parse row 6: %w", err)
	}
	if len(values) < 4 {
		return masjidRow{}, fmt.Errorf("masjidboardlive: row 6 has %d fields, need at least 4", len(values))
	}

	return masjidRow{
		Name1:    stringValue(values, 0),
		Name2:    stringValue(values, 1),
		URL:      stringValue(values, 2),
		Timezone: stringValue(values, 3),
	}, nil
}

func parseSalahRow(raw json.RawMessage, date time.Time, loc *time.Location) (model.PrayerTimes, error) {
	values, err := rowValues(raw)
	if err != nil {
		return model.PrayerTimes{}, fmt.Errorf("masjidboardlive: parse row 3: %w", err)
	}
	if len(values) < 10 {
		return model.PrayerTimes{}, fmt.Errorf("masjidboardlive: row 3 has %d fields, need at least 10", len(values))
	}

	fields := make([]*time.Time, 10)
	for i := range fields {
		value := stringValue(values, i)
		if value == "" || value == "-" {
			continue
		}
		parsed, err := parseLocalTime(value, date, loc)
		if err != nil {
			return model.PrayerTimes{}, fmt.Errorf("masjidboardlive: row 3 column %d value %q: %w", i, value, err)
		}
		fields[i] = parsed
	}

	// Row 3 columns 0..9 are verified by functions_uo_latest.js:
	// Fajr Adhan/Jamaah, Zuhr Adhan/Jamaah, Asr Adhan/Jamaah,
	// Maghrib/Iftar, Maghrib Jamaah, Esha Adhan/Jamaah.
	return model.PrayerTimes{
		Fajr:    model.PrayerTime{Adhan: fields[0], Jamaah: fields[1]},
		Zuhr:    model.PrayerTime{Adhan: fields[2], Jamaah: fields[3]},
		Asr:     model.PrayerTime{Adhan: fields[4], Jamaah: fields[5]},
		Maghrib: model.PrayerTime{Adhan: fields[6], Jamaah: fields[7]},
		Esha:    model.PrayerTime{Adhan: fields[8], Jamaah: fields[9]},
	}, nil
}

func rowValues(raw json.RawMessage) ([]json.RawMessage, error) {
	var values []json.RawMessage
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil, err
	}
	return values, nil
}

func stringValue(values []json.RawMessage, index int) string {
	if index < 0 || index >= len(values) {
		return ""
	}

	var value string
	if err := json.Unmarshal(values[index], &value); err != nil {
		return ""
	}
	return strings.TrimSpace(value)
}

func parseLocalTime(value string, date time.Time, loc *time.Location) (*time.Time, error) {
	value = strings.TrimSpace(value)
	layouts := []string{
		"15:04",
		"15:04:05",
		"3:04 PM",
		"3:04PM",
		"3:04 pm",
		"3:04pm",
	}

	for _, layout := range layouts {
		parsed, err := time.ParseInLocation(layout, value, loc)
		if err == nil {
			result := time.Date(date.Year(), date.Month(), date.Day(), parsed.Hour(), parsed.Minute(), parsed.Second(), 0, loc)
			return &result, nil
		}
	}

	return nil, fmt.Errorf("unsupported time format")
}

func dateOnly(t time.Time, loc *time.Location) time.Time {
	local := t.In(loc)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
}

// Keep the Jumuah row index named here even while its full optional mapping is
// deferred. This documents the verified 29-row source layout without making
// unverified assumptions about optional display semantics.
var _ = rowJumuah
