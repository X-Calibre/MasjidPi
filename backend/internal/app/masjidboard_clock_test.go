package app

import (
	"testing"
	"time"
)

func TestMasjidBoardCalendarDateChangedAfterNTPCorrection(t *testing.T) {
	originalLocation := time.Local
	time.Local = time.FixedZone("Africa/Johannesburg", 2*60*60)
	t.Cleanup(func() { time.Local = originalLocation })

	startup := time.Date(2026, time.September, 18, 19, 0, 42, 0, time.Local)
	synchronized := time.Date(2026, time.September, 21, 19, 41, 0, 0, time.Local)

	previousDate := masjidBoardLocalDate(startup)
	if previousDate != "2026-09-18" {
		t.Fatalf("startup date = %q, want 2026-09-18", previousDate)
	}
	if !masjidBoardCalendarDateChanged(previousDate, synchronized) {
		t.Fatal("calendar date change after NTP correction was not detected")
	}
}

func TestMasjidBoardCalendarDateChangedDoesNotRepeatForCurrentDate(t *testing.T) {
	originalLocation := time.Local
	time.Local = time.FixedZone("Africa/Johannesburg", 2*60*60)
	t.Cleanup(func() { time.Local = originalLocation })

	now := time.Date(2026, time.September, 21, 19, 41, 0, 0, time.Local)
	currentDate := masjidBoardLocalDate(now)

	if masjidBoardCalendarDateChanged(currentDate, now.Add(30*time.Second)) {
		t.Fatal("unchanged local calendar date requested another refresh")
	}
}

func TestMasjidBoardCalendarDateChangedRequiresPreviousRefresh(t *testing.T) {
	now := time.Date(2026, time.September, 21, 17, 41, 0, 0, time.UTC)
	if masjidBoardCalendarDateChanged("", now) {
		t.Fatal("unconfigured service requested a date-change refresh")
	}
}
