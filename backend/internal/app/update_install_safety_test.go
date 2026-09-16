package app

import (
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
	masjidboardruntime "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/playback"
)

type fakeUpdatePlaybackProvider struct {
	status playback.Status
}

func (p fakeUpdatePlaybackProvider) Status() playback.Status {
	return p.status
}

func (p fakeUpdatePlaybackProvider) Stop() {}

type fakeUpdateBoardProvider struct {
	results []masjidboardruntime.Result
}

func (p fakeUpdateBoardProvider) Results() []masjidboardruntime.Result {
	return p.results
}

func TestCurrentInstallConditionsUsesPlaybackAndPrimaryBoard(t *testing.T) {
	now := time.Date(2026, time.September, 18, 23, 0, 0, 0, time.FixedZone("SAST", 7200))
	primary := &model.Board{PrayerTimes: model.PrayerTimes{
		Fajr: model.PrayerTime{Adhan: &model.ClockTime{Hour: 5, Minute: 30}},
		Esha: model.PrayerTime{Adhan: &model.ClockTime{Hour: 19, Minute: 15}},
		Jumuah: []model.JumuahService{{
			Adhan: &model.ClockTime{Hour: 12, Minute: 45},
		}},
	}}
	secondary := &model.Board{PrayerTimes: model.PrayerTimes{
		Fajr: model.PrayerTime{Adhan: &model.ClockTime{Hour: 6, Minute: 0}},
	}}
	conditions := currentInstallConditions(
		now,
		true,
		true,
		fakeUpdatePlaybackProvider{status: playback.Status{Listening: true}},
		fakeUpdateBoardProvider{results: []masjidboardruntime.Result{
			{Board: primary},
			{Board: secondary},
		}},
	)

	if !conditions.PlaybackActive || !conditions.Immediate || !conditions.InterruptPlayback {
		t.Fatalf("conditions = %+v", conditions)
	}
	if len(conditions.AdhanTimes) != 3 {
		t.Fatalf("Adhan times = %v, want three primary-board times", conditions.AdhanTimes)
	}
	for _, adhan := range conditions.AdhanTimes {
		if adhan.Location() != now.Location() || adhan.Day() != now.Day() {
			t.Fatalf("Adhan time = %v, want current local date", adhan)
		}
	}
}

func TestCurrentInstallConditionsSupportsMissingComponents(t *testing.T) {
	conditions := currentInstallConditions(time.Now(), false, false, nil, nil)
	if conditions.PlaybackActive || len(conditions.AdhanTimes) != 0 {
		t.Fatalf("conditions = %+v", conditions)
	}
}
