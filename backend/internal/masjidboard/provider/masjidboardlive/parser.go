package masjidboardlive

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
)

// MasjidBoard Live returns 29 positional top-level rows. Keep these indexes
// inside this package so the rest of MasjidBoard never depends on the source
// layout.
const (
	rowJumuah       = 1
	rowClock        = 2
	rowSalah        = 3
	rowAstronomical = 5
	rowMasjid       = 6
)

// Parse normalises the verified core portion of a MasjidBoard Live response.
// Optional board content is only added when its source semantics are verified
// and represented by the domain model.
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

	clock, err := parseClockRow(rows[rowClock])
	if err != nil {
		return model.Board{}, err
	}
	if clock.Timezone == "" {
		return model.Board{}, fmt.Errorf("masjidboardlive: clock row has no timezone")
	}

	loc, err := parseTimezone(clock.Timezone, masjid.OffsetMS)
	if err != nil {
		return model.Board{}, err
	}

	localNow := now.In(loc)
	prayers, err := parseSalahRow(rows[rowSalah], localNow, loc)
	if err != nil {
		return model.Board{}, err
	}

	astronomical, err := parseAstronomicalRow(rows[rowAstronomical], localNow, loc)
	if err != nil {
		return model.Board{}, err
	}

	jumuah, err := parseJumuahRow(rows[rowJumuah], localNow, loc)
	if err != nil {
		return model.Board{}, err
	}

	board := model.Board{
		Identity: model.Identity{
			SourceBoardID: boardID,
			MasjidID:      clock.MasjidID,
			EnglishName:   masjid.Name1,
			ArabicName:    masjid.Name2,
			Location:      clock.Location,
		},
		DateContext: model.DateContext{
			GregorianDate: dateOnly(localNow, loc),
			Timezone:      clock.Timezone,
		},
		PrayerTimes: prayers,
	}

	if astronomical != nil {
		board.AstronomicalTimes = astronomical
	}
	if len(jumuah) > 0 {
		board.JumuahServices = jumuah
	}

	return board, nil
}

type masjidRow struct {
	Name1    string
	Name2    string
	URL      string
	OffsetMS int64
}

func parseMasjidRow(raw json.RawMessage) (masjidRow, error) {
	values, err := rowValues(raw)
	if err != nil {
		return masjidRow{}, fmt.Errorf("masjidboardlive: parse row 6: %w", err)
	}
	if len(values) < 5 {
		return masjidRow{}, fmt.Errorf("masjidboardlive: row 6 has %d fields, need at least 5", len(values))
	}

	offset, err := strconv.ParseInt(stringValue(values, 4), 10, 64)
	if err != nil {
		return masjidRow{}, fmt.Errorf("masjidboardlive: invalid row 6 timezone offset %q: %w", stringValue(values, 4), err)
	}

	return masjidRow{
		Name1:    stringValue(values, 0),
		Name2:    stringValue(values, 1),
		URL:      stringValue(values, 2),
		OffsetMS: offset,
	}, nil
}

type clockRow struct {
	MasjidID string
	Location string
	Timezone string
}

func parseClockRow(raw json.RawMessage) (clockRow, error) {
	values, err := rowValues(raw)
	if err != nil {
		return clockRow{}, fmt.Errorf("masjidboardlive: parse row 2: %w", err)
	}
	if len(values) < 16 {
		return clockRow{}, fmt.Errorf("masjidboardlive: row 2 has %d fields, need at least 16", len(values))
	}

	return clockRow{
		MasjidID: stringValue(values, 12),
		Location: stringValue(values, 14),
		Timezone: stringValue(values, 15),
	}, nil
}

var gmtOffsetRE = regexp.MustCompile(`^GMT([+-])(\d{2})(?::?(\d{2}))?$`)

func parseGMTOffset(value string) (int64, error) {
	match := gmtOffsetRE.FindStringSubmatch(strings.TrimSpace(value))
	if match == nil {
		return 0, fmt.Errorf("masjidboardlive: unsupported timezone %q", value)
	}
	hours, _ := strconv.Atoi(match[2])
	minutes := 0
	if match[3] != "" {
		minutes, _ = strconv.Atoi(match[3])
	}
	seconds := hours*3600 + minutes*60
	if match[1] == "-" {
		seconds = -seconds
	}
	return int64(seconds) * 1000, nil
}

func parseTimezone(label string, offsetMS int64) (*time.Location, error) {
	label = strings.TrimSpace(label)
	if label == "" {
		return nil, fmt.Errorf("masjidboardlive: timezone is empty")
	}
	if offsetMS == 0 {
		parsed, err := parseGMTOffset(label)
		if err != nil {
			return nil, err
		}
		offsetMS = parsed
	}
	return time.FixedZone(label, int(offsetMS/1000)), nil
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

	return model.PrayerTimes{
		Fajr:    model.PrayerTime{Adhan: fields[0], Jamaah: fields[1]},
		Zuhr:    model.PrayerTime{Adhan: fields[2], Jamaah: fields[3]},
		Asr:     model.PrayerTime{Adhan: fields[4], Jamaah: fields[5]},
		Maghrib: model.PrayerTime{Adhan: fields[6], Jamaah: fields[7]},
		Esha:    model.PrayerTime{Adhan: fields[8], Jamaah: fields[9]},
	}, nil
}

// parseJumuahRow handles the portion of row 1 whose semantics are explicitly
// labelled by the upstream response. In the captured response the first
// service is represented as alternating labels and times:
//
//   Lecture -> 12:15
//   Adhan   -> 12:45
//   Khutbah -> 12:55
//
// Later row-1 values are intentionally not interpreted here. The upstream
// response contains additional values, but their semantics have not yet been
// verified well enough to map them into the domain model without guessing.
func parseJumuahRow(raw json.RawMessage, date time.Time, loc *time.Location) ([]model.JumuahService, error) {
	values, err := rowValues(raw)
	if err != nil {
		return nil, fmt.Errorf("masjidboardlive: parse row 1: %w", err)
	}
	if len(values) < 6 {
		return nil, fmt.Errorf("masjidboardlive: row 1 has %d fields, need at least 6", len(values))
	}

	lecture, err := parseOptionalLocalTime(stringValue(values, 1), date, loc)
	if err != nil {
		return nil, fmt.Errorf("masjidboardlive: row 1 lecture time: %w", err)
	}
	adhan, err := parseOptionalLocalTime(stringValue(values, 3), date, loc)
	if err != nil {
		return nil, fmt.Errorf("masjidboardlive: row 1 adhan time: %w", err)
	}
	khutbah, err := parseOptionalLocalTime(stringValue(values, 5), date, loc)
	if err != nil {
		return nil, fmt.Errorf("masjidboardlive: row 1 khutbah time: %w", err)
	}

	if lecture == nil && adhan == nil && khutbah == nil {
		return nil, nil
	}

	return []model.JumuahService{{
		Title:   "Jumu'ah",
		Adhan:   adhan,
		Lecture: lecture,
		Khutbah: khutbah,
	}}, nil
}

func parseOptionalLocalTime(value string, date time.Time, loc *time.Location) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" || value == "-" {
		return nil, nil
	}
	return parseLocalTime(value, date, loc)
}

func parseAstronomicalRow(raw json.RawMessage, date time.Time, loc *time.Location) (*model.AstronomicalTimes, error) {
	values, err := rowValues(raw)
	if err != nil {
		return nil, fmt.Errorf("masjidboardlive: parse row 5: %w", err)
	}
	if len(values) < 9 {
		return nil, fmt.Errorf("masjidboardlive: row 5 has %d fields, need at least 9", len(values))
	}

	parsed := make([]*time.Time, 9)
	for i := range parsed {
		value := stringValue(values, i)
		if value == "" || value == "-" {
			continue
		}
		valueTime, err := parseLocalTime(value, date, loc)
		if err != nil {
			return nil, fmt.Errorf("masjidboardlive: row 5 column %d value %q: %w", i, value, err)
		}
		parsed[i] = valueTime
	}

	for _, value := range parsed {
		if value != nil {
			return &model.AstronomicalTimes{
				Suhur:     parsed[0],
				FajrStart: parsed[1],
				Sunrise:   parsed[2],
				Ishraaq:   parsed[3],
				Duha:      parsed[4],
				AsrShafii: parsed[5],
				AsrHanafi: parsed[6],
				Sunset:    parsed[7],
				EshaStart: parsed[8],
			}, nil
		}
	}

	return nil, nil
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
