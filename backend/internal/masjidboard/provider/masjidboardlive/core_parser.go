package masjidboardlive

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
)

// CoreMetadata contains MasjidBoard Live fields that are useful to the
// provider but do not belong in the provider-independent Board domain model.
type CoreMetadata struct {
	MBLNumber          string
	LastUpdated        string
	Language           string
	TwentyFourHourTime string
	LiveStreamProvider string
	LiveStreamURL      string
	RamadhaanHide      string
	SundayDhuhrText    string
}

// CoreResult is the normalised result of parsing the public Core board data
// object embedded by https://masjidboardlive.com/boards/?<web_url>.
type CoreResult struct {
	Board    model.Board
	Metadata CoreMetadata
}

var coreFieldRE = regexp.MustCompile(`(?m)^\s*([A-Za-z_][A-Za-z0-9_]*)\s*:\s*("(?:\\.|[^"\\])*")\s*,?\s*$`)

// ParseCoreObject normalises the verified 36-field JavaScript data object used
// by public MasjidBoard Live Core boards. Identity and timezone come from the
// discovery/catalogue layer because the embedded Core object does not contain
// enough information to construct the normalised Board identity by itself.
func ParseCoreObject(raw []byte, identity model.BoardIdentity, now time.Time) (CoreResult, error) {
	if identity.ID == "" {
		return CoreResult{}, fmt.Errorf("masjidboardlive: Core board identity ID is required")
	}
	if identity.Name == "" {
		return CoreResult{}, fmt.Errorf("masjidboardlive: Core board identity name is required")
	}
	if identity.TimeZone == "" {
		return CoreResult{}, fmt.Errorf("masjidboardlive: Core board timezone is required")
	}

	fields, err := parseCoreFields(raw)
	if err != nil {
		return CoreResult{}, err
	}
	if strings.TrimSpace(fields["mbl_number"]) == "" {
		return CoreResult{}, fmt.Errorf("masjidboardlive: Core data has no mbl_number")
	}

	loc, err := parseTimezone(identity.TimeZone, 0)
	if err != nil {
		return CoreResult{}, err
	}

	parseOptionalClock := func(name string) (*model.ClockTime, error) {
		value := strings.TrimSpace(fields[name])
		if isAbsent(value) {
			return nil, nil
		}
		parsed, parseErr := parseClockTime(value)
		if parseErr != nil {
			return nil, fmt.Errorf("masjidboardlive: Core field %s value %q: %w", name, value, parseErr)
		}
		return parsed, nil
	}

	prayers, err := parseCorePrayerTimes(parseOptionalClock)
	if err != nil {
		return CoreResult{}, err
	}

	astronomical, err := parseCoreAstronomical(parseOptionalClock)
	if err != nil {
		return CoreResult{}, err
	}

	jumuah, err := parseCoreJumuah(fields, parseOptionalClock)
	if err != nil {
		return CoreResult{}, err
	}
	if jumuah != nil {
		prayers.Jumuah = []model.JumuahService{*jumuah}
	}

	return CoreResult{
		Board: model.Board{
			Identity: identity,
			DateContext: model.DateContext{
				GregorianDate: dateOnly(now, loc),
			},
			PrayerTimes:  prayers,
			Astronomical: astronomical,
		},
		Metadata: CoreMetadata{
			MBLNumber:          fields["mbl_number"],
			LastUpdated:        fields["last_updated"],
			Language:           fields["lang"],
			TwentyFourHourTime: fields["twentyfourhrtime"],
			LiveStreamProvider: fields["liveStreamP"],
			LiveStreamURL:      fields["liveStreamURL"],
			RamadhaanHide:      fields["ramadaanHide"],
			SundayDhuhrText:    fields["sunday_zuhr_text"],
		},
	}, nil
}

func parseCoreFields(raw []byte) (map[string]string, error) {
	matches := coreFieldRE.FindAllSubmatch(raw, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("masjidboardlive: Core data object contains no recognised fields")
	}

	fields := make(map[string]string, len(matches))
	for _, match := range matches {
		value, err := strconv.Unquote(string(match[2]))
		if err != nil {
			return nil, fmt.Errorf("masjidboardlive: decode Core field %s: %w", match[1], err)
		}
		fields[string(match[1])] = strings.TrimSpace(value)
	}
	return fields, nil
}

func parseCorePrayerTimes(parseClock func(string) (*model.ClockTime, error)) (model.PrayerTimes, error) {
	fajrAdhan, err := parseClock("fajrAthan")
	if err != nil {
		return model.PrayerTimes{}, err
	}
	fajrJamaah, err := parseClock("fajrJamaah")
	if err != nil {
		return model.PrayerTimes{}, err
	}
	dhuhrAdhan, err := parseClock("dhuhrAthan")
	if err != nil {
		return model.PrayerTimes{}, err
	}
	dhuhrJamaah, err := parseClock("dhuhrJamaah")
	if err != nil {
		return model.PrayerTimes{}, err
	}
	asrAdhan, err := parseClock("asrAthan")
	if err != nil {
		return model.PrayerTimes{}, err
	}
	asrJamaah, err := parseClock("asrJamaah")
	if err != nil {
		return model.PrayerTimes{}, err
	}
	maghribAdhan, err := parseClock("maghribAthan")
	if err != nil {
		return model.PrayerTimes{}, err
	}
	eshaAdhan, err := parseClock("eshaAthan")
	if err != nil {
		return model.PrayerTimes{}, err
	}
	eshaJamaah, err := parseClock("eshaJamaah")
	if err != nil {
		return model.PrayerTimes{}, err
	}

	prayers := model.PrayerTimes{
		Fajr:    model.PrayerTime{Adhan: fajrAdhan, Jamaah: fajrJamaah},
		Dhuhr:   model.PrayerTime{Adhan: dhuhrAdhan, Jamaah: dhuhrJamaah},
		Asr:     model.PrayerTime{Adhan: asrAdhan, Jamaah: asrJamaah},
		Maghrib: model.PrayerTime{Adhan: maghribAdhan},
		Esha:    model.PrayerTime{Adhan: eshaAdhan, Jamaah: eshaJamaah},
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
			return model.PrayerTimes{}, fmt.Errorf("masjidboardlive: Core data missing prayer time for %s", check.name)
		}
	}

	return prayers, nil
}

func parseCoreAstronomical(parseClock func(string) (*model.ClockTime, error)) (*model.AstronomicalTimes, error) {
	names := []string{
		"sehriEnds",
		"fajrStarts",
		"sunrise",
		"ishraaq",
		"duha",
		"istiwaCaution",
		"istiwa",
		"zawaalEnd",
		"asrShafi",
		"asrHanafi",
		"sunset",
		"eshaStarts",
	}
	values := make([]*model.ClockTime, len(names))
	any := false
	for i, name := range names {
		parsed, err := parseClock(name)
		if err != nil {
			return nil, err
		}
		values[i] = parsed
		if parsed != nil {
			any = true
		}
	}
	if !any {
		return nil, nil
	}

	return &model.AstronomicalTimes{
		Suhur:         values[0],
		FajrStart:     values[1],
		Sunrise:       values[2],
		Ishraaq:       values[3],
		Duha:          values[4],
		IstiwaCaution: values[5],
		Istiwa:        values[6],
		ZawaalEnd:     values[7],
		AsrShafii:     values[8],
		AsrHanafi:     values[9],
		Sunset:        values[10],
		EshaStart:     values[11],
	}, nil
}

func parseCoreJumuah(fields map[string]string, parseClock func(string) (*model.ClockTime, error)) (*model.JumuahService, error) {
	codes := strings.Split(fields["jumuahHeadings"], ",")
	events := make([]model.JumuahEvent, 0, 3)
	for i := 0; i < 3; i++ {
		field := fmt.Sprintf("jumuahTime%d", i+1)
		parsed, err := parseClock(field)
		if err != nil {
			return nil, err
		}
		code := ""
		if i < len(codes) {
			code = strings.TrimSpace(codes[i])
		}
		if parsed == nil && code == "" {
			continue
		}
		events = append(events, model.JumuahEvent{
			Code:    code,
			Heading: coreJumuahHeading(code),
			Time:    parsed,
		})
	}
	if len(events) == 0 {
		return nil, nil
	}

	service := &model.JumuahService{Events: events}
	for _, event := range events {
		if event.Time != nil && event.Code == "0" && service.Adhan == nil {
			service.Adhan = event.Time
		}
	}
	return service, nil
}

func coreJumuahHeading(code string) string {
	switch strings.TrimSpace(code) {
	case "0":
		return "Adhan"
	case "1":
		return "Lecture"
	case "3":
		return "Sunan"
	case "6":
		return "Khutbah"
	default:
		return ""
	}
}
