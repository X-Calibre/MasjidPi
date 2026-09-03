package masjidboardlive

import (
	"regexp"
	"strings"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
)

var coreUpcomingSalaahHTMLRE = regexp.MustCompile(`(?is)<(?:h[1-6]|div|span)\b[^>]*\bid=["'](fajrNextDate|fajrNextTime|asrNextDate|asrNextTime|eshaNextDate|eshaNextTime)["'][^>]*>(.*?)</(?:h[1-6]|div|span)\s*>`)

// applyCoreUpcomingSalaahChangeFallback extracts the rendered next-change
// values used by public/Core-only boards. These fields are not present in the
// Core page's embedded data object, so Premium parsing remains the preferred
// structured source when it is available.
func applyCoreUpcomingSalaahChangeFallback(result *CoreResult, page []byte) {
	if result == nil || hasSalaahChangeNotice(result.Board.Notices) {
		return
	}

	values := extractCoreUpcomingSalaahHTML(page)
	for _, prayer := range []struct {
		name    string
		dateKey string
		timeKey string
	}{
		{name: "Fajr", dateKey: "fajrNextDate", timeKey: "fajrNextTime"},
		{name: "Asr", dateKey: "asrNextDate", timeKey: "asrNextTime"},
		{name: "Esha", dateKey: "eshaNextDate", timeKey: "eshaNextTime"},
	} {
		dateValue := strings.TrimSpace(values[prayer.dateKey])
		timeValue := strings.TrimSpace(values[prayer.timeKey])
		if isAbsent(dateValue) || isAbsent(timeValue) || dateValue == "00:00" {
			continue
		}

		fields := map[string]string{
			"prayer":         prayer.name,
			"effective_date": dateValue,
			"new_time":       timeValue,
		}
		result.Board.Notices = append(result.Board.Notices, model.Notice{
			Type:    model.NoticeTypeSalaahChange,
			Title:   prayer.name + " Time Change",
			Content: joinNoticeValues(fields, "effective_date", "new_time"),
			Fields:  fields,
		})
	}
}

func extractCoreUpcomingSalaahHTML(page []byte) map[string]string {
	values := make(map[string]string, 6)
	for _, match := range coreUpcomingSalaahHTMLRE.FindAllSubmatch(page, -1) {
		if len(match) != 3 {
			continue
		}
		values[string(match[1])] = cleanCoreHTMLText(match[2])
	}
	return values
}

func hasSalaahChangeNotice(notices []model.Notice) bool {
	for _, notice := range notices {
		if notice.Type == model.NoticeTypeSalaahChange {
			return true
		}
	}
	return false
}
