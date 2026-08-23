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

const (
	rowUpcoming       = 0
	rowJumuah        = 1
	rowClock         = 2
	rowSalah         = 3
	rowAstronomical  = 5
	rowMasjid        = 6
	rowTaleem        = 10
	rowAnnouncement  = 11
	rowAnnouncement2 = 12
	rowNikah         = 13
	rowFuneral       = 14
	rowEid           = 17
	rowWellWishes    = 21
)

// Parse normalises the verified core portion of a MasjidBoard Live response.
// The 29-row upstream structure remains confined to this provider package.
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

	prayers, err := parseSalahRow(rows[rowSalah])
	if err != nil {
		return model.Board{}, err
	}

	jumuah, err := parseJumuahRow(rows[rowJumuah])
	if err != nil {
		return model.Board{}, err
	}
	if jumuah != nil {
		prayers.Jumuah = []model.JumuahService{*jumuah}
	}

	astronomical, err := parseAstronomicalRow(rows[rowAstronomical])
	if err != nil {
		return model.Board{}, err
	}
	name := strings.TrimSpace(masjid.Name1)
	alternateName := strings.TrimSpace(masjid.Name2)
	if name == "" {
		name = alternateName
		alternateName = ""
	}
	if name == "" {
		return model.Board{}, fmt.Errorf("masjidboardlive: masjid row has no usable name")
	}
	announcements := parseAnnouncementRows(rows[rowAnnouncement], rows[rowAnnouncement2])
	notices := parseNoticeRows(rows[rowNikah], rows[rowFuneral], rows[rowEid])
	notices = append(parseUpcomingSalaahChanges(rows[rowUpcoming], localNow), notices...)
	notices = append(notices, parseWellWishesRow(rows[rowWellWishes])...)
	programmes := parseTaleemRow(rows[rowTaleem])
	newMoon := parseNewMoonRow(rows[rowClock])

	return model.Board{
		Identity: model.BoardIdentity{
			ID:            clock.MasjidID,
			Name:          name,
			AlternateName: alternateName,
			TimeZone:      clock.Timezone,
		},
		DateContext: model.DateContext{
			GregorianDate: dateOnly(localNow, loc),
			IslamicDate:   stringValueFromRow(rows[rowClock], 5),
		},
		PrayerTimes:   prayers,
		Astronomical:  astronomical,
		Announcements: announcements,
		Programmes:    programmes,
		Notices:       notices,
		NewMoon:       newMoon,
	}, nil
}

func parseUpcomingSalaahChanges(raw json.RawMessage, localNow time.Time) []model.Notice {
	values, err := rowValues(raw)
	if err != nil || len(values) < 19 {
		return nil
	}
	offset := 0
	if isDisplayed(stringValue(values, 9)) {
		offset = 10
	}
	prayers := []string{"Fajr", "Asr", "Esha"}
	notices := make([]model.Notice, 0, len(prayers))
	today := dateOnly(localNow, localNow.Location())
	for index, prayer := range prayers {
		base := offset + index*3
		dateValue := strings.TrimSpace(stringValue(values, base))
		timeValue := strings.TrimSpace(stringValue(values, base+1))
		milliValue := strings.TrimSpace(stringValue(values, base+2))
		if isAbsent(dateValue) || isAbsent(timeValue) {
			continue
		}
		if millis, parseErr := strconv.ParseInt(milliValue, 10, 64); parseErr == nil && millis > 0 {
			effective := time.UnixMilli(millis).In(localNow.Location())
			if dateOnly(effective, localNow.Location()).Before(today) {
				continue
			}
		}
		fields := map[string]string{"prayer": prayer, "effective_date": dateValue, "new_time": timeValue}
		notices = append(notices, model.Notice{
			Type: model.NoticeTypeSalaahChange, Title: prayer + " Time Change",
			Content: joinNoticeValues(fields, "effective_date", "new_time"), Fields: fields,
		})
	}
	return notices
}

func parseTaleemRow(raw json.RawMessage) []model.Programme {
	values, err := rowValues(raw)
	if err != nil {
		return nil
	}
	programmes := make([]model.Programme, 0, 2)
	for start := 0; start+4 < len(values) && start < 10; start += 5 {
		if !isDisplayed(stringValue(values, start+4)) {
			continue
		}
		parts := make([]string, 0, 4)
		for index := start; index < start+4; index++ {
			if value := strings.TrimSpace(stringValue(values, index)); !isAbsent(value) {
				parts = append(parts, value)
			}
		}
		if len(parts) > 0 {
			programmes = append(programmes, model.Programme{Title: "Taleem Programme", Content: strings.Join(parts, "\n")})
		}
	}
	return programmes
}

func parseWellWishesRow(raw json.RawMessage) []model.Notice {
	values, err := rowValues(raw)
	if err != nil || len(values) < 11 || !isDisplayed(stringValue(values, 10)) {
		return nil
	}
	notices := make([]model.Notice, 0, 10)
	for index := 0; index < 10; index++ {
		message := strings.TrimSpace(stringValue(values, index))
		if !isAbsent(message) {
			notices = append(notices, model.Notice{Type: model.NoticeTypeWellWish, Title: "Du'a Requested", Content: message})
		}
	}
	return notices
}

func parseNewMoonRow(raw json.RawMessage) *model.NewMoon {
	values, err := rowValues(raw)
	if err != nil || len(values) < 22 || !isDisplayed(stringValue(values, 21)) {
		return nil
	}
	fields := namedFields(values, map[int]string{
		0: "birth", 1: "first_moonset", 2: "first_age", 3: "first_azimuth", 4: "first_altitude",
		5: "best_visibility", 6: "second_moonset", 7: "second_age", 8: "second_azimuth",
		9: "second_altitude", 10: "birth_date", 11: "visibility_date",
	})
	if len(fields) == 0 {
		return nil
	}
	return &model.NewMoon{Fields: fields}
}

func parseAnnouncementRows(rawRows ...json.RawMessage) []model.Announcement {
	announcements := make([]model.Announcement, 0, 10)
	for _, raw := range rawRows {
		values, err := rowValues(raw)
		if err != nil {
			continue
		}
		for i := 0; i+2 < len(values) && i < 15; i += 3 {
			title := strings.TrimSpace(stringValue(values, i))
			content := strings.TrimSpace(stringValue(values, i+1))
			if !isDisplayed(stringValue(values, i+2)) || (title == "" && content == "") {
				continue
			}
			announcements = append(announcements, model.Announcement{Title: title, Content: content})
		}
	}
	return announcements
}

func parseNoticeRows(nikahRaw, funeralRaw, eidRaw json.RawMessage) []model.Notice {
	notices := make([]model.Notice, 0, 3)
	if notice, ok := parseNikahNotice(nikahRaw); ok {
		notices = append(notices, notice)
	}
	if notice, ok := parseFuneralNotice(funeralRaw); ok {
		notices = append(notices, notice)
	}
	if notice, ok := parseEidNotice(eidRaw); ok {
		notices = append(notices, notice)
	}
	return notices
}

func parseNikahNotice(raw json.RawMessage) (model.Notice, bool) {
	values, err := rowValues(raw)
	if err != nil || len(values) < 9 || !isDisplayed(stringValue(values, 8)) {
		return model.Notice{}, false
	}
	fields := namedFields(values, map[int]string{
		0: "name_one", 1: "groom_relation", 2: "relation_one", 3: "relation_two",
		4: "name_two", 5: "date", 6: "time", 7: "event_timestamp", 10: "bride",
	})
	if len(fields) == 0 {
		return model.Notice{}, false
	}
	return model.Notice{
		Type: model.NoticeTypeNikah, Title: "Nikah Notice",
		Content: joinNoticeValues(fields, "name_one", "groom_relation", "relation_one", "relation_two", "name_two", "bride", "date", "time"),
		Fields:  fields,
	}, true
}

func parseFuneralNotice(raw json.RawMessage) (model.Notice, bool) {
	values, err := rowValues(raw)
	if err != nil || len(values) < 8 || !isDisplayed(stringValue(values, 7)) {
		return model.Notice{}, false
	}
	fields := namedFields(values, map[int]string{
		0: "name", 1: "relation", 2: "address", 3: "pickup",
		4: "cemetery", 5: "salaah_venue", 6: "salaah_time",
	})
	if len(fields) == 0 {
		return model.Notice{}, false
	}
	return model.Notice{
		Type: model.NoticeTypeFuneral, Title: "Funeral Notice",
		Content: joinNoticeValues(fields, "name", "relation", "address", "pickup", "cemetery", "salaah_venue", "salaah_time"),
		Fields:  fields,
	}, true
}

func parseEidNotice(raw json.RawMessage) (model.Notice, bool) {
	values, err := rowValues(raw)
	if err != nil || len(values) < 6 || !isDisplayed(stringValue(values, 5)) {
		return model.Notice{}, false
	}
	fields := namedFields(values, map[int]string{
		0: "date", 1: "venue", 2: "address", 3: "lecture", 4: "salaah",
	})
	if len(fields) == 0 {
		return model.Notice{}, false
	}
	return model.Notice{
		Type: model.NoticeTypeEid, Title: "Eid Salaah Notice",
		Content: joinNoticeValues(fields, "date", "venue", "address", "lecture", "salaah"),
		Fields:  fields,
	}, true
}

func namedFields(values []json.RawMessage, names map[int]string) map[string]string {
	fields := make(map[string]string, len(names))
	for index, name := range names {
		value := strings.TrimSpace(stringValue(values, index))
		if !isAbsent(value) {
			fields[name] = value
		}
	}
	return fields
}

func joinNoticeValues(fields map[string]string, names ...string) string {
	values := make([]string, 0, len(names))
	for _, name := range names {
		if value := strings.TrimSpace(fields[name]); value != "" {
			values = append(values, value)
		}
	}
	return strings.Join(values, "\n")
}

func isDisplayed(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "display", "show", "visible", "true", "yes":
		return true
	default:
		return false
	}
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
	return masjidRow{Name1: stringValue(values, 0), Name2: stringValue(values, 1), URL: stringValue(values, 2), OffsetMS: offset}, nil
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
	return clockRow{MasjidID: stringValue(values, 12), Location: stringValue(values, 14), Timezone: stringValue(values, 15)}, nil
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

func parseSalahRow(raw json.RawMessage) (model.PrayerTimes, error) {
	values, err := rowValues(raw)
	if err != nil {
		return model.PrayerTimes{}, fmt.Errorf("masjidboardlive: parse row 3: %w", err)
	}
	if len(values) < 10 {
		return model.PrayerTimes{}, fmt.Errorf("masjidboardlive: row 3 has %d fields, need at least 10", len(values))
	}

	fields := make([]*model.ClockTime, 10)
	for i := range fields {
		value := stringValue(values, i)
		if isAbsent(value) {
			continue
		}
		parsed, err := parseClockTime(value)
		if err != nil {
			return model.PrayerTimes{}, fmt.Errorf("masjidboardlive: row 3 column %d value %q: %w", i, value, err)
		}
		fields[i] = parsed
	}

	prayers := model.PrayerTimes{
		Fajr:    model.PrayerTime{Adhan: fields[0], Jamaah: fields[1]},
		Dhuhr:   model.PrayerTime{Adhan: fields[2], Jamaah: fields[3]},
		Asr:     model.PrayerTime{Adhan: fields[4], Jamaah: fields[5]},
		Maghrib: model.PrayerTime{Adhan: fields[6], Jamaah: fields[7]},
		Esha:    model.PrayerTime{Adhan: fields[8], Jamaah: fields[9]},
	}

	checks := []struct {
		name  string
		value model.PrayerTime
	}{
		{"Fajr", prayers.Fajr},
		{"Dhuhr", prayers.Dhuhr},
		{"Asr", prayers.Asr},
		{"Maghrib", prayers.Maghrib},
		{"Esha", prayers.Esha},
	}
	for _, check := range checks {
		if check.value.Adhan == nil && check.value.Jamaah == nil {
			return model.PrayerTimes{}, fmt.Errorf("masjidboardlive: missing core prayer time for %s", check.name)
		}
	}
	return prayers, nil
}

// parseJumuahRow normalises row 1. The row contains three detailed
// heading/time pairs followed by dedicated Jumu'ah Adhan, Jamaah, alternate
// language values and the source heading-code configuration.
func parseJumuahRow(raw json.RawMessage) (*model.JumuahService, error) {
	values, err := rowValues(raw)
	if err != nil {
		return nil, fmt.Errorf("masjidboardlive: parse row 1: %w", err)
	}
	if len(values) < 12 {
		return nil, fmt.Errorf("masjidboardlive: row 1 has %d fields, need at least 12", len(values))
	}

	hasData := false
	for i := 0; i < 12; i++ {
		if !isAbsent(stringValue(values, i)) {
			hasData = true
			break
		}
	}
	if !hasData {
		return nil, nil
	}

	codes := strings.Split(stringValue(values, 11), ",")
	events := make([]model.JumuahEvent, 0, 3)
	for i := 0; i < 3; i++ {
		heading := stringValue(values, i*2)
		timeValue := stringValue(values, i*2+1)
		var parsed *model.ClockTime
		if !isAbsent(timeValue) {
			parsed, err = parseClockTime(timeValue)
			if err != nil {
				return nil, fmt.Errorf("masjidboardlive: row 1 Jumuah time %d value %q: %w", i+1, timeValue, err)
			}
		}
		code := ""
		if i < len(codes) {
			code = strings.TrimSpace(codes[i])
		}
		if !isAbsent(heading) || parsed != nil || code != "" {
			events = append(events, model.JumuahEvent{Code: code, Heading: heading, Time: parsed})
		}
	}

	parseOptional := func(index int) (*model.ClockTime, error) {
		value := stringValue(values, index)
		if isAbsent(value) {
			return nil, nil
		}
		parsed, parseErr := parseClockTime(value)
		if parseErr != nil {
			return nil, fmt.Errorf("masjidboardlive: row 1 column %d value %q: %w", index, value, parseErr)
		}
		return parsed, nil
	}

	adhan, err := parseOptional(7)
	if err != nil {
		return nil, err
	}
	jamaah, err := parseOptional(8)
	if err != nil {
		return nil, err
	}
	alternateAdhan, err := parseOptional(9)
	if err != nil {
		return nil, err
	}
	alternateJamaah, err := parseOptional(10)
	if err != nil {
		return nil, err
	}

	return &model.JumuahService{
		Adhan:           adhan,
		Jamaah:          jamaah,
		AlternateAdhan:  alternateAdhan,
		AlternateJamaah: alternateJamaah,
		Khateeb:         stringValue(values, 6),
		Events:          events,
	}, nil
}

func parseAstronomicalRow(raw json.RawMessage) (*model.AstronomicalTimes, error) {
	values, err := rowValues(raw)
	if err != nil {
		return nil, fmt.Errorf("masjidboardlive: parse row 5: %w", err)
	}
	if len(values) < 9 {
		return nil, fmt.Errorf("masjidboardlive: row 5 has %d fields, need at least 9", len(values))
	}
	parsed := make([]*model.ClockTime, 9)
	any := false
	for i := range parsed {
		value := stringValue(values, i)
		if isAbsent(value) {
			continue
		}
		valueTime, err := parseClockTime(value)
		if err != nil {
			return nil, fmt.Errorf("masjidboardlive: row 5 column %d value %q: %w", i, value, err)
		}
		parsed[i] = valueTime
		any = true
	}
	if !any {
		return nil, nil
	}
	return &model.AstronomicalTimes{
		Suhur: parsed[0], FajrStart: parsed[1], Sunrise: parsed[2], Ishraaq: parsed[3], Duha: parsed[4],
		AsrShafii: parsed[5], AsrHanafi: parsed[6], Sunset: parsed[7], EshaStart: parsed[8],
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

func stringValueFromRow(raw json.RawMessage, index int) string {
	values, err := rowValues(raw)
	if err != nil {
		return ""
	}
	return stringValue(values, index)
}

func isAbsent(value string) bool {
	switch strings.TrimSpace(value) {
	case "", "-", "–", "—", "~~~~", "FALSE", "false", "Hide", "hide", "#N/A":
		return true
	default:
		return false
	}
}

func parseClockTime(value string) (*model.ClockTime, error) {
	value = strings.TrimSpace(value)
	layouts := []string{"15:04", "15:04:05", "3:04 PM", "3:04PM", "3:04 pm", "3:04pm"}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return &model.ClockTime{Hour: parsed.Hour(), Minute: parsed.Minute()}, nil
		}
	}
	return nil, fmt.Errorf("unsupported time format")
}

func dateOnly(t time.Time, loc *time.Location) time.Time {
	local := t.In(loc)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
}
