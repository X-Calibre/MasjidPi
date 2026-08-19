package display

import (
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
	masjidboardruntime "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

// Build creates the read-only display view from selected-board runtime state.
// Selection order is authoritative. A selected board with no usable live or
// cached timetable is still returned with status=unavailable so the display
// can represent the slot without inventing data.
func Build(configured bool, selected selection.State, results []masjidboardruntime.Result) View {
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
			populateBoard(&item, *result.Board)
		}
		view.Boards = append(view.Boards, item)
	}
	return view
}

func populateBoard(out *Board, board model.Board) {
	out.Name = board.Identity.Name
	out.AlternateName = board.Identity.AlternateName
	out.TimeZone = board.Identity.TimeZone
	if !board.DateContext.GregorianDate.IsZero() {
		out.Date.Gregorian = board.DateContext.GregorianDate.Format("2006-01-02")
	}
	out.Date.Islamic = board.DateContext.IslamicDate

	out.Prayers = []Prayer{
		{Key: "fajr", Label: "Fajr", Adhan: cloneTime(board.PrayerTimes.Fajr.Adhan), Jamaah: cloneTime(board.PrayerTimes.Fajr.Jamaah)},
		{Key: "dhuhr", Label: "Dhuhr", Adhan: cloneTime(board.PrayerTimes.Dhuhr.Adhan), Jamaah: cloneTime(board.PrayerTimes.Dhuhr.Jamaah)},
		{Key: "asr", Label: "Asr", Adhan: cloneTime(board.PrayerTimes.Asr.Adhan), Jamaah: cloneTime(board.PrayerTimes.Asr.Jamaah)},
		{Key: "maghrib", Label: "Maghrib", Adhan: cloneTime(board.PrayerTimes.Maghrib.Adhan), Jamaah: cloneTime(board.PrayerTimes.Maghrib.Jamaah)},
		{Key: "esha", Label: "Esha", Adhan: cloneTime(board.PrayerTimes.Esha.Adhan), Jamaah: cloneTime(board.PrayerTimes.Esha.Jamaah)},
	}

	if len(board.PrayerTimes.Jumuah) > 0 {
		out.Jumuah = make([]JumuahService, 0, len(board.PrayerTimes.Jumuah))
		for _, service := range board.PrayerTimes.Jumuah {
			presented := JumuahService{
				Adhan:           cloneTime(service.Adhan),
				Jamaah:          cloneTime(service.Jamaah),
				EffectiveSalaah: cloneTime(service.EffectiveSalaah()),
				AlternateAdhan:  cloneTime(service.AlternateAdhan),
				AlternateJamaah: cloneTime(service.AlternateJamaah),
				Khateeb:         service.Khateeb,
			}
			if len(service.Events) > 0 {
				presented.Events = make([]JumuahEvent, 0, len(service.Events))
				for _, event := range service.Events {
					presented.Events = append(presented.Events, JumuahEvent{
						Code: event.Code, Heading: event.Heading, Time: cloneTime(event.Time),
					})
				}
			}
			out.Jumuah = append(out.Jumuah, presented)
		}
	}

	if board.Astronomical != nil {
		a := board.Astronomical
		out.Astronomical = &Astronomical{
			Suhur: cloneTime(a.Suhur), FajrStart: cloneTime(a.FajrStart), Sunrise: cloneTime(a.Sunrise),
			Ishraaq: cloneTime(a.Ishraaq), Duha: cloneTime(a.Duha), IstiwaCaution: cloneTime(a.IstiwaCaution),
			Istiwa: cloneTime(a.Istiwa), ZawaalEnd: cloneTime(a.ZawaalEnd), AsrShafii: cloneTime(a.AsrShafii),
			AsrHanafi: cloneTime(a.AsrHanafi), Sunset: cloneTime(a.Sunset), EshaStart: cloneTime(a.EshaStart),
		}
	}
}

func cloneTime(value *model.ClockTime) *model.ClockTime {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
