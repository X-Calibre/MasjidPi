package masjidboardlive

import (
	"html"
	"regexp"
	"strconv"
	"strings"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
)

var coreJumuahHeadingHTMLRE = regexp.MustCompile(`(?is)<h2\s+id=["']jumuahHead([123])["'][^>]*>(.*?)</h2>`)
var coreJumuahTimeHTMLRE = regexp.MustCompile(`(?is)<h3\s+id=["']jumuahTime([123])["'][^>]*>(.*?)</h3>`)
var coreHTMLTagRE = regexp.MustCompile(`(?is)<[^>]+>`)

type coreJumuahHTMLSlot struct {
	heading string
	time    *model.ClockTime
}

func extractCoreJumuahHTMLSlots(page []byte) [3]coreJumuahHTMLSlot {
	var slots [3]coreJumuahHTMLSlot

	for _, match := range coreJumuahHeadingHTMLRE.FindAllSubmatch(page, -1) {
		if len(match) != 3 {
			continue
		}
		index := int(match[1][0] - '1')
		if index < 0 || index >= len(slots) {
			continue
		}
		slots[index].heading = cleanCoreHTMLText(match[2])
	}

	for _, match := range coreJumuahTimeHTMLRE.FindAllSubmatch(page, -1) {
		if len(match) != 3 {
			continue
		}
		index := int(match[1][0] - '1')
		if index < 0 || index >= len(slots) {
			continue
		}
		value := cleanCoreHTMLText(match[2])
		if value == "" {
			continue
		}
		if parsed, ok := parseCoreHTMLClock(value); ok {
			slots[index].time = parsed
		}
	}

	return slots
}

func cleanCoreHTMLText(raw []byte) string {
	value := coreHTMLTagRE.ReplaceAllString(string(raw), "")
	value = html.UnescapeString(value)
	value = strings.ReplaceAll(value, "\u00a0", " ")
	return strings.TrimSpace(value)
}

func parseCoreHTMLClock(value string) (*model.ClockTime, bool) {
	parts := strings.Split(strings.TrimSpace(value), ":")
	if len(parts) != 2 {
		return nil, false
	}
	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, false
	}
	minute, err := strconv.Atoi(parts[1])
	if err != nil || hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return nil, false
	}
	return &model.ClockTime{Hour: hour, Minute: minute}, true
}

func applyCoreJumuahHeadingFallback(result *CoreResult, page []byte) {
	if result == nil || len(result.Board.PrayerTimes.Jumuah) == 0 {
		return
	}

	slots := extractCoreJumuahHTMLSlots(page)
	service := &result.Board.PrayerTimes.Jumuah[0]
	for i := range service.Events {
		event := &service.Events[i]
		if strings.TrimSpace(event.Heading) != "" || event.Time == nil {
			continue
		}

		slot := matchingCoreJumuahHTMLSlot(slots, event.Time)
		if slot == nil {
			continue
		}
		heading := normalizeCoreJumuahHTMLHeading(slot.heading)
		if heading == "" {
			continue
		}
		event.Heading = heading
		event.Code = coreJumuahCodeForHeading(heading)
		if heading == "Adhan" && service.Adhan == nil {
			service.Adhan = event.Time
		}
	}
}

func matchingCoreJumuahHTMLSlot(slots [3]coreJumuahHTMLSlot, wanted *model.ClockTime) *coreJumuahHTMLSlot {
	if wanted == nil {
		return nil
	}
	for i := range slots {
		if slots[i].time == nil {
			continue
		}
		if slots[i].time.Hour == wanted.Hour && slots[i].time.Minute == wanted.Minute {
			return &slots[i]
		}
	}
	return nil
}

func normalizeCoreJumuahHTMLHeading(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "adhan", "adhaan", "adhān":
		return "Adhan"
	case "lecture":
		return "Lecture"
	case "sunan", "sunnan":
		return "Sunan"
	case "khutbah", "khutba":
		return "Khutbah"
	default:
		return strings.TrimSpace(value)
	}
}

func coreJumuahCodeForHeading(heading string) string {
	switch heading {
	case "Adhan":
		return "0"
	case "Lecture":
		return "1"
	case "Sunan":
		return "3"
	case "Khutbah":
		return "6"
	default:
		return ""
	}
}
