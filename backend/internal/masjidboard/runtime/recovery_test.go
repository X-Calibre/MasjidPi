package runtime

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestFetchRecoversFromStaleCacheToCurrent(t *testing.T) {
	currentTime := time.Date(2026, 8, 20, 12, 0, 0, 0, time.UTC)
	provider := &fakeProvider{board: board("brits-jamia", "Brits Jamia Masjid", 6)}
	store := &fakeCache{entries: make(map[string]cache.Entry)}
	coord, err := newWithClock(
		[]Item{{Selection: selectedBoard("brits-jamia", "Brits Jamia Masjid"), Provider: provider}},
		store,
		func() time.Time { return currentTime },
	)
	if err != nil {
		t.Fatal(err)
	}

	first := coord.FetchAll(context.Background())[0]
	if first.Status != StatusCurrent || first.Board == nil || first.UpdateError != nil {
		t.Fatalf("initial refresh = %+v, want current", first)
	}
	if !first.LastSuccessfulUpdate.Equal(currentTime) {
		t.Fatalf("initial success timestamp = %v, want %v", first.LastSuccessfulUpdate, currentTime)
	}

	initialSuccess := currentTime
	currentTime = currentTime.Add(30 * time.Minute)
	provider.err = errors.New("provider offline")

	stale := coord.FetchAll(context.Background())[0]
	if stale.Status != StatusStale || stale.Board == nil || stale.UpdateError == nil {
		t.Fatalf("offline refresh = %+v, want stale cached data", stale)
	}
	if !stale.LastSuccessfulUpdate.Equal(initialSuccess) {
		t.Fatalf("stale success timestamp = %v, want cached %v", stale.LastSuccessfulUpdate, initialSuccess)
	}
	if stale.Board.PrayerTimes.Fajr.Jamaah == nil || stale.Board.PrayerTimes.Fajr.Jamaah.Hour != 6 {
		t.Fatalf("stale board did not preserve cached timetable: %+v", stale.Board)
	}

	currentTime = currentTime.Add(30 * time.Minute)
	provider.err = nil
	provider.board = board("brits-jamia", "Brits Jamia Masjid", 7)

	recovered := coord.FetchAll(context.Background())[0]
	if recovered.Status != StatusCurrent || recovered.Board == nil || recovered.UpdateError != nil {
		t.Fatalf("recovery refresh = %+v, want current", recovered)
	}
	if !recovered.LastSuccessfulUpdate.Equal(currentTime) {
		t.Fatalf("recovered success timestamp = %v, want %v", recovered.LastSuccessfulUpdate, currentTime)
	}
	if recovered.Board.PrayerTimes.Fajr.Jamaah == nil || recovered.Board.PrayerTimes.Fajr.Jamaah.Hour != 7 {
		t.Fatalf("recovery did not publish fresh timetable: %+v", recovered.Board)
	}

	cached, found, err := store.Load("masjidboardlive:brits-jamia")
	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Fatal("recovered timetable was not cached")
	}
	if !cached.SuccessfulAt.Equal(currentTime) || cached.Board.PrayerTimes.Fajr.Jamaah == nil || cached.Board.PrayerTimes.Fajr.Jamaah.Hour != 7 {
		t.Fatalf("cache after recovery = %+v", cached)
	}
}
