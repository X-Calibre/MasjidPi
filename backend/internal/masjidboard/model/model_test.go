package model

import "testing"

func TestPrayerTimesUseLocalClockValues(t *testing.T) {
	board := Board{
		Identity: BoardIdentity{
			ID:            "test-masjid",
			Name:          "Test Masjid",
			AlternateName: "Test Alt",
			TimeZone:      "GMT+02",
		},
		PrayerTimes: PrayerTimes{
			Fajr:    PrayerTime{Adhan: &ClockTime{Hour: 5, Minute: 0}},
			Dhuhr:   PrayerTime{Adhan: &ClockTime{Hour: 12, Minute: 30}},
			Asr:     PrayerTime{Adhan: &ClockTime{Hour: 16, Minute: 0}},
			Maghrib: PrayerTime{Adhan: &ClockTime{Hour: 17, Minute: 45}},
			Esha:    PrayerTime{Adhan: &ClockTime{Hour: 19, Minute: 0}},
		},
	}

	if board.Identity.ID == "" || board.Identity.Name == "" || board.Identity.TimeZone == "" {
		t.Fatal("expected fundamental board identity")
	}
	if board.PrayerTimes.Fajr.Adhan.Hour != 5 || board.PrayerTimes.Fajr.Adhan.Minute != 0 {
		t.Fatal("expected local wall-clock Fajr time")
	}
	if board.PrayerTimes.Dhuhr.Adhan.Hour != 12 || board.PrayerTimes.Dhuhr.Adhan.Minute != 30 {
		t.Fatal("expected local wall-clock Dhuhr time")
	}
}

func TestJumuahBelongsToPrayerTimes(t *testing.T) {
	board := Board{PrayerTimes: PrayerTimes{
		Jumuah: []JumuahService{{
			Adhan:  &ClockTime{Hour: 12, Minute: 45},
			Jamaah: &ClockTime{Hour: 13, Minute: 0},
		}},
	}}

	if len(board.PrayerTimes.Jumuah) != 1 {
		t.Fatalf("Jumuah services = %d, want 1", len(board.PrayerTimes.Jumuah))
	}
	if board.PrayerTimes.Jumuah[0].Jamaah.Hour != 13 || board.PrayerTimes.Jumuah[0].Jamaah.Minute != 0 {
		t.Fatal("expected Jumuah Jamaah time")
	}
}

func TestOptionalContentCanBeAbsent(t *testing.T) {
	board := Board{}
	if board.Astronomical != nil || board.Banking != nil || board.NewMoon != nil {
		t.Fatal("expected optional singular content to be absent")
	}
	if board.Announcements != nil || board.Programmes != nil || board.Media != nil {
		t.Fatal("expected optional collections to be absent")
	}
}
