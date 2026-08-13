package model

import (
	"testing"
	"time"
)

func TestPrayerTimesHaveFiveCorePrayers(t *testing.T) {
	location := time.FixedZone("SAST", 2*60*60)
	date := time.Date(2026, 8, 13, 0, 0, 0, 0, location)
	prayer := func(hour, minute int) *time.Time {
		v := time.Date(2026, 8, 13, hour, minute, 0, 0, location)
		return &v
	}

	board := Board{
		Identity: Identity{MasjidID: "test-masjid", EnglishName: "Test Masjid"},
		DateContext: DateContext{
			GregorianDate: date,
			Timezone:      "Africa/Johannesburg",
		},
		PrayerTimes: PrayerTimes{
			Fajr:    PrayerTime{Adhan: prayer(5, 0)},
			Zuhr:    PrayerTime{Adhan: prayer(12, 30)},
			Asr:     PrayerTime{Adhan: prayer(16, 0)},
			Maghrib: PrayerTime{Adhan: prayer(17, 45)},
			Esha:    PrayerTime{Adhan: prayer(19, 0)},
		},
	}

	if board.Identity.MasjidID == "" {
		t.Fatal("expected masjid identity")
	}
	if board.PrayerTimes.Fajr.Adhan == nil ||
		board.PrayerTimes.Zuhr.Adhan == nil ||
		board.PrayerTimes.Asr.Adhan == nil ||
		board.PrayerTimes.Maghrib.Adhan == nil ||
		board.PrayerTimes.Esha.Adhan == nil {
		t.Fatal("expected all five core prayer times")
	}
}

func TestOptionalContentCanBeAbsent(t *testing.T) {
	board := Board{
		Identity:    Identity{MasjidID: "test-masjid"},
		DateContext: DateContext{GregorianDate: time.Now()},
		PrayerTimes: PrayerTimes{
			Fajr:    PrayerTime{},
			Zuhr:    PrayerTime{},
			Asr:     PrayerTime{},
			Maghrib: PrayerTime{},
			Esha:    PrayerTime{},
		},
	}

	if board.AstronomicalTimes != nil || board.Banking != nil || board.NewMoon != nil {
		t.Fatal("expected optional singular content to be absent")
	}
	if board.JumuahServices != nil || board.Announcements != nil || board.Programmes != nil || board.Media != nil {
		t.Fatal("expected optional collections to be absent")
	}
}
