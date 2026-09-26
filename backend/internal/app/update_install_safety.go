package app

import (
	"time"

	"github.com/X-Calibre/MasjidFrame/backend/internal/masjidboard/model"
	masjidboardruntime "github.com/X-Calibre/MasjidFrame/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidFrame/backend/internal/playback"
	"github.com/X-Calibre/MasjidFrame/backend/internal/updates"
)

type updatePlaybackProvider interface {
	Status() playback.Status
	Stop()
}

type updateBoardProvider interface {
	Results() []masjidboardruntime.Result
}

func currentInstallConditions(
	now time.Time,
	immediate bool,
	interruptPlayback bool,
	playbackProvider updatePlaybackProvider,
	boardProvider updateBoardProvider,
) updates.InstallConditions {
	conditions := updates.InstallConditions{
		Now:               now,
		Immediate:         immediate,
		InterruptPlayback: interruptPlayback,
	}
	if playbackProvider != nil {
		conditions.PlaybackActive = playbackProvider.Status().Listening
	}
	if boardProvider != nil {
		conditions.AdhanTimes = primaryAdhanTimes(now, boardProvider.Results())
	}
	return conditions
}

func primaryAdhanTimes(
	now time.Time,
	results []masjidboardruntime.Result,
) []time.Time {
	if len(results) == 0 || results[0].Board == nil {
		return nil
	}
	board := results[0].Board
	prayers := board.PrayerTimes
	clocks := []*model.ClockTime{
		prayers.Fajr.Adhan,
		prayers.Dhuhr.Adhan,
		prayers.Asr.Adhan,
		prayers.Maghrib.Adhan,
		prayers.Esha.Adhan,
	}
	if now.Weekday() == time.Friday {
		for _, service := range prayers.Jumuah {
			clocks = append(clocks, service.Adhan)
		}
	}

	times := make([]time.Time, 0, len(clocks))
	year, month, day := now.Date()
	for _, clock := range clocks {
		if clock == nil {
			continue
		}
		times = append(times, time.Date(
			year,
			month,
			day,
			clock.Hour,
			clock.Minute,
			0,
			0,
			now.Location(),
		))
	}
	return times
}
