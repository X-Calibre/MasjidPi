package display

import (
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
	masjidboardruntime "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

func selected(id, name string) selection.Board {
	return selection.Board{
		CatalogueID:      "masjidboardlive:" + id,
		Provider:         "masjidboardlive",
		ExternalID:       id,
		Name:             name,
		TimeZoneOffsetMS: 7200000,
	}
}

func ct(hour, minute int) *model.ClockTime { return &model.ClockTime{Hour: hour, Minute: minute} }

func TestBuildPreservesSelectionOrderAndBuildsTimetable(t *testing.T) {
	one := selected("one", "One")
	two := selected("two", "Two")
	updated := time.Date(2026, 8, 19, 19, 0, 0, 0, time.FixedZone("SAST", 2*60*60))
	board := model.Board{
		Identity:    model.BoardIdentity{ID: "two", Name: "Two Live", TimeZone: "GMT+02:00"},
		DateContext: model.DateContext{GregorianDate: time.Date(2026, 8, 19, 0, 0, 0, 0, updated.Location())},
		PrayerTimes: model.PrayerTimes{
			Fajr:   model.PrayerTime{Adhan: ct(5, 40), Jamaah: ct(6, 0)},
			Asr:    model.PrayerTime{Adhan: ct(16, 40), Jamaah: ct(17, 0)},
			Jumuah: []model.JumuahService{{Events: []model.JumuahEvent{{Code: "6", Heading: "Khutbah", Time: ct(13, 0)}}}},
		},
		Astronomical:  &model.AstronomicalTimes{Sunrise: ct(6, 33), Sunset: ct(17, 51)},
		Announcements: []model.Announcement{{Title: "Community update", Content: "<b>Tonight</b>"}},
		Programmes:    []model.Programme{{Title: "Taleem Programme", Content: "After Esha"}},
		Notices: []model.Notice{{
			Type: model.NoticeTypeFuneral, Title: "Funeral Notice",
			Fields: map[string]string{"name": "Abdullah", "salaah_time": "14:30"},
		}},
		Banking: &model.Banking{Content: "Masjid Contributions", Fields: map[string]string{"bank": "Example Bank", "account_number": "000123456"}},
		NewMoon: &model.NewMoon{Fields: map[string]string{"visibility_date": "12 September 2026"}},
	}

	view := Build(true, selection.State{Layout: selection.LayoutPortrait, Theme: selection.ThemeRuby, SlideDurationSeconds: 30, Boards: []selection.Board{one, two}}, []masjidboardruntime.Result{{
		Selection: two, Board: &board, Status: masjidboardruntime.StatusCurrent, LastSuccessfulUpdate: updated,
	}})

	if !view.Configured || len(view.Boards) != 2 {
		t.Fatalf("view = %+v", view)
	}
	if view.Layout != selection.LayoutPortrait || view.Theme != selection.ThemeRuby || view.SlideDuration != 30 {
		t.Fatalf("display preferences = layout %q theme %q duration %d", view.Layout, view.Theme, view.SlideDuration)
	}
	if view.Boards[0].CatalogueID != one.CatalogueID || view.Boards[0].Status != masjidboardruntime.StatusUnavailable {
		t.Fatalf("first board = %+v", view.Boards[0])
	}
	got := view.Boards[1]
	if got.Name != "Two Live" || got.Status != masjidboardruntime.StatusCurrent || got.Stale {
		t.Fatalf("second board = %+v", got)
	}
	if len(got.Prayers) != 5 || got.Prayers[0].Key != "fajr" || got.Prayers[2].Key != "asr" {
		t.Fatalf("prayers = %+v", got.Prayers)
	}
	if got.Prayers[0].Jamaah == nil || got.Prayers[0].Jamaah.Hour != 6 {
		t.Fatalf("fajr = %+v", got.Prayers[0])
	}
	if len(got.Jumuah) != 1 || got.Jumuah[0].EffectiveSalaah != nil {
		t.Fatalf("jumuah effective salaah = %+v, want nil without explicit Jamaah", got.Jumuah)
	}
	if len(got.Jumuah[0].Events) != 1 || got.Jumuah[0].Events[0].Heading != "Khutbah" || got.Jumuah[0].Events[0].Time == nil || got.Jumuah[0].Events[0].Time.Hour != 13 {
		t.Fatalf("jumuah events = %+v", got.Jumuah[0].Events)
	}
	if got.Astronomical == nil || got.Astronomical.Sunrise == nil || got.Astronomical.Sunset == nil {
		t.Fatalf("astronomical = %+v", got.Astronomical)
	}
	if got.Date.Gregorian != "2026-08-19" {
		t.Fatalf("gregorian = %q", got.Date.Gregorian)
	}
	if len(got.Announcements) != 1 || got.Announcements[0].Content != "<b>Tonight</b>" {
		t.Fatalf("announcements = %+v", got.Announcements)
	}
	if len(got.Programmes) != 1 || got.Programmes[0].Content != "After Esha" {
		t.Fatalf("programmes = %+v", got.Programmes)
	}
	if len(got.Notices) != 1 || got.Notices[0].Type != string(model.NoticeTypeFuneral) || got.Notices[0].Fields["salaah_time"] != "14:30" {
		t.Fatalf("notices = %+v", got.Notices)
	}
	got.Notices[0].Fields["name"] = "changed"
	if board.Notices[0].Fields["name"] != "Abdullah" {
		t.Fatal("display notice fields alias the cached domain model")
	}
	if got.Banking == nil || got.Banking.Title != "Masjid Contributions" || got.Banking.Fields["account_number"] != "000123456" {
		t.Fatalf("banking = %+v", got.Banking)
	}
	got.Banking.Fields["bank"] = "changed"
	if board.Banking.Fields["bank"] != "Example Bank" {
		t.Fatal("display banking fields alias the cached domain model")
	}
	if got.NewMoon == nil || got.NewMoon.Fields["visibility_date"] != "12 September 2026" {
		t.Fatalf("new moon = %+v", got.NewMoon)
	}
	got.NewMoon.Fields["visibility_date"] = "changed"
	if board.NewMoon.Fields["visibility_date"] != "12 September 2026" {
		t.Fatal("display new-moon fields alias the cached domain model")
	}
}

func TestBuildMarksCachedBoardStaleWithoutDiagnostics(t *testing.T) {
	selectedBoard := selected("brits-jamia", "Brits Jamia Masjid")
	updated := time.Date(2026, 8, 19, 18, 0, 0, 0, time.UTC)
	board := model.Board{Identity: model.BoardIdentity{ID: "brits-jamia", Name: "Brits Jamia Masjid", TimeZone: "GMT+02:00"}}

	view := Build(true, selection.State{Boards: []selection.Board{selectedBoard}}, []masjidboardruntime.Result{{
		Selection:            selectedBoard,
		Board:                &board,
		Status:               masjidboardruntime.StatusStale,
		LastSuccessfulUpdate: updated,
	}})

	got := view.Boards[0]
	if got.Status != masjidboardruntime.StatusStale || !got.Stale {
		t.Fatalf("board = %+v", got)
	}
	if got.LastSuccessfulUpdate == nil || !got.LastSuccessfulUpdate.Equal(updated) {
		t.Fatalf("last successful update = %v", got.LastSuccessfulUpdate)
	}
}
