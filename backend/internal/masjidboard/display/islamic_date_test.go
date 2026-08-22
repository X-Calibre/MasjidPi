package display

import (
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
	masjidboardruntime "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

func TestBuildAtRecalculatesIslamicDateForCachedBoard(t *testing.T) {
	selectedBoard := selected("brits-jamia", "Brits Jamia Masjid")
	loc := time.FixedZone("GMT+02", 2*60*60)
	board := model.Board{
		Identity: model.BoardIdentity{ID: "brits-jamia", Name: "Brits Jamia Masjid", TimeZone: "GMT+02"},
		DateContext: model.DateContext{
			GregorianDate:         time.Date(2026, 7, 16, 0, 0, 0, 0, loc),
			IslamicDate:           "1 Safar 1448",
			IslamicDateAdjustment: 0,
		},
		Astronomical: &model.AstronomicalTimes{Sunset: ct(17, 51)},
	}
	results := []masjidboardruntime.Result{{
		Selection: selectedBoard,
		Board:     &board,
		Status:    masjidboardruntime.StatusStale,
	}}

	before := BuildAt(
		true,
		selection.State{Boards: []selection.Board{selectedBoard}},
		results,
		time.Date(2026, 7, 16, 17, 54, 5, 0, loc),
	)
	if got := before.Boards[0].Date.Islamic; got != "1 Safar 1448" {
		t.Fatalf("before sunset rollover Islamic date = %q", got)
	}

	after := BuildAt(
		true,
		selection.State{Boards: []selection.Board{selectedBoard}},
		results,
		time.Date(2026, 7, 16, 17, 54, 6, 0, loc),
	)
	if got := after.Boards[0].Date.Islamic; got != "2 Safar 1448" {
		t.Fatalf("after sunset rollover Islamic date = %q", got)
	}
}
