package display

import (
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
	masjidboardruntime "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

// Build creates the read-only display view from selected-board runtime state.
// Selection order is authoritative. A selected board with no usable live or
// cached timetable is still returned with status=unavailable so the display
// can represent the slot without inventing data.
func Build(configured bool, selected selection.State, results []masjidboardruntime.Result) View {
	return BuildAt(configured, selected, results, time.Now())
}

// BuildAt is the deterministic form of Build used by tests. Recalculating the
// Islamic date while building the view means cached timetable data still rolls
// over at the same sunset boundary as MasjidBoard Live, without waiting for a
// network refresh.
func BuildAt(configured bool, selected selection.State, results []masjidboardruntime.Result, now time.Time) View {
	byID := make(map[string]masjidboardruntime.Result, len(results))
	for _, result := range results {
		byID[result.Selection.CatalogueID] = result
	}

	view := View{Configured: configured, Boards: make([]Board, 0, len(selected.Boards))}
	for _, selectedBoard := range selected.Boards {
		item := Board{
			CatalogueID: selectedBoard.CatalogueID,
			Name:        selectedBoard.Name,
			Status:      masjidboardruntime.StatusUnavailable,
		}

		result, ok := byID[selectedBoard.CatalogueID]
		if !ok {
			view.Boards = append(view.Boards, item)
			continue
		}

		item.Status = result.Status
		item.Stale = result.Status == masjidboardruntime.StatusStale
		if !result.LastSuccessfulUpdate.IsZero() {
			updated := result.LastSuccessfulUpdate
			item.LastSuccessfulUpdate = &updated
		}
		if result.Board != nil {
			populateBoard(&item, *result.Board, now)
		}
		view.Boards = append(view.Boards, item)
	}
	return view
}

func populateBoard(out *Board, board model.Board, now time.Time) {
	out.Name = board.Identity.Name
	out.AlternateName = board.Identity.AlternateName
	out.TimeZone = board.Identity.TimeZone
	if !board.DateContext.GregorianDate.IsZero() {
		out.Date.Gregorian = board.DateContext.GregorianDate.Format("2006-01-02")
	}
	out.Date.Islamic = board.DateContext.IslamicDate

	// A non-empty stored Islamic date marks a board whose provider supplied the
	// calculation context. Recalculate from the current clock so cached data is
	// not stuck on the pre-sunset date until the next provider refresh.
	if board.DateContext.IslamicDate != "" {
		loc := now.Location()
		if !board.DateContext.GregorianDate.IsZero() {
			loc = board.DateContext.GregorianDate.Location()
		}
		var sunset *model.ClockTime
		if board.Astronomical != nil {
			sunset = board.Astronomical.Sunset
		}
		out.Date.Islamic = model.CalculateIslamicDate(
			now.In(loc),
			sunset,
			board.DateContext.IslamicDateAdjustment,
			board.DateContext.ForceIslamicDate30,
		).String()
	}

	out.Prayers = []Prayer{
		{Key: "fajr", Label: "Fajr", Adhan: displayTime(board.PrayerTimes.Fajr.Adhan), Jamaah: displayTime(board.PrayerTimes.Fajr.Jamaah)},
		{Key: "dhuhr", Label: "Dhuhr", Adhan: displayTime(board.PrayerTimes.Dhuhr.Adhan), Jamaah: displayTime(board.PrayerTimes.Dhuhr.Jamaah)},
		{Key: "asr", Label: "Asr", Adhan: displayTime(board.PrayerTimes.Asr.Adhan), Jamaah: displayTime(board.PrayerTimes.Asr.Jamaah)},
		{Key: "maghrib", Label: "Maghrib", Adhan: displayTime(board.PrayerTimes.Maghrib.Adhan), Jamaah: displayTime(board.PrayerTimes.Maghrib.Jamaah)},
		{Key: "esha", Label: "Esha", Adhan: displayTime(board.PrayerTimes.Esha.Adhan), Jamaah: displayTime(board.PrayerTimes.Esha.Jamaah)},
	}

	if len(board.PrayerTimes.Jumuah) > 0 {
		out.Jumuah = make([]JumuahService, 0, len(board.PrayerTimes.Jumuah))
		for _, service := range board.PrayerTimes.Jumuah {
			presented := JumuahService{
				Adhan:           displayTime(service.Adhan),
				Jamaah:          displayTime(service.Jamaah),
				EffectiveSalaah: displayTime(service.EffectiveSalaah()),
				AlternateAdhan:  displayTime(service.AlternateAdhan),
				AlternateJamaah: displayTime(service.AlternateJamaah),
				Khateeb:         service.Khateeb,
			}
			if len(service.Events) > 0 {
				presented.Events = make([]JumuahEvent, 0, len(service.Events))
				for _, event := range service.Events {
					presented.Events = append(presented.Events, JumuahEvent{
						Code: event.Code, Heading: event.Heading, Time: displayTime(event.Time),
					})
				}
			}
			out.Jumuah = append(out.Jumuah, presented)
		}
	}

	if board.Astronomical != nil {
		a := board.Astronomical
		out.Astronomical = &Astronomical{
			Suhur: displayTime(a.Suhur), FajrStart: displayTime(a.FajrStart), Sunrise: displayTime(a.Sunrise),
			Ishraaq: displayTime(a.Ishraaq), Duha: displayTime(a.Duha), IstiwaCaution: displayTime(a.IstiwaCaution),
			Istiwa: displayTime(a.Istiwa), ZawaalEnd: displayTime(a.ZawaalEnd), AsrShafii: displayTime(a.AsrShafii),
			AsrHanafi: displayTime(a.AsrHanafi), Sunset: displayTime(a.Sunset), EshaStart: displayTime(a.EshaStart),
		}
	}

	if len(board.Announcements) > 0 {
		out.Announcements = make([]Announcement, 0, len(board.Announcements))
		for _, announcement := range board.Announcements {
			out.Announcements = append(out.Announcements, Announcement{
				Title: announcement.Title, Content: announcement.Content,
			})
		}
	}
	if len(board.Programmes) > 0 {
		out.Programmes = make([]Programme, 0, len(board.Programmes))
		for _, programme := range board.Programmes {
			out.Programmes = append(out.Programmes, Programme{Title: programme.Title, Content: programme.Content})
		}
	}
	if len(board.Notices) > 0 {
		out.Notices = make([]Notice, 0, len(board.Notices))
		for _, notice := range board.Notices {
			fields := make(map[string]string, len(notice.Fields))
			for key, value := range notice.Fields {
				fields[key] = value
			}
			out.Notices = append(out.Notices, Notice{
				Type: string(notice.Type), Title: notice.Title, Content: notice.Content, Fields: fields,
			})
		}
	}
	if board.NewMoon != nil {
		fields := make(map[string]string, len(board.NewMoon.Fields))
		for key, value := range board.NewMoon.Fields {
			fields[key] = value
		}
		out.NewMoon = &NewMoon{Fields: fields}
	}
}

func displayTime(value *model.ClockTime) *ClockTime {
	if value == nil {
		return nil
	}
	return &ClockTime{Hour: value.Hour, Minute: value.Minute}
}
