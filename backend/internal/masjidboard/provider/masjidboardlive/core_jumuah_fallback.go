package masjidboardlive

import (
	"html"
	"regexp"
	"strings"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
)

var coreJumuahHeadingHTMLRE = regexp.MustCompile(`(?is)<h2\s+id=["']jumuahHead([123])["'][^>]*>(.*?)</h2>`)
var coreHTMLTagRE = regexp.MustCompile(`(?is)<[^>]+>`)

func extractCoreJumuahHeadings(page []byte) [3]string {
	var headings [3]string
	for _, match := range coreJumuahHeadingHTMLRE.FindAllSubmatch(page, -1) {
		if len(match) != 3 {
			continue
		}
		index := int(match[1][0] - '1')
		if index < 0 || index >= len(headings) {
			continue
		}
		value := coreHTMLTagRE.ReplaceAllString(string(match[2]), "")
		value = html.UnescapeString(value)
		headings[index] = strings.TrimSpace(value)
	}
	return headings
}

func applyCoreJumuahHeadingFallback(result *CoreResult, page []byte) {
	if result == nil || len(result.Board.PrayerTimes.Jumuah) == 0 {
		return
	}

	headings := extractCoreJumuahHeadings(page)
	service := &result.Board.PrayerTimes.Jumuah[0]
	for i := range service.Events {
		event := &service.Events[i]
		if strings.TrimSpace(event.Heading) != "" || event.Time == nil {
			continue
		}

		position := coreJumuahEventPosition(service.Events, i)
		if position < 0 || position >= len(headings) {
			continue
		}
		heading := normalizeCoreJumuahHTMLHeading(headings[position])
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

// parseCoreJumuah omits completely empty heading/time slots. Reconstruct the
// original slot by counting prior events with a time or heading. This keeps the
// HTML fallback aligned with jumuahTime1..3 when a middle slot is blank.
func coreJumuahEventPosition(events []model.JumuahEvent, target int) int {
	position := 0
	for i := 0; i <= target && i < len(events); i++ {
		if i == target {
			return position
		}
		position++
	}
	return -1
}

func normalizeCoreJumuahHTMLHeading(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "adhan", "adhaan", "adhan̄", "adhān":
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
