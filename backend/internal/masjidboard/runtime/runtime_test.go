package runtime

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/cache"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

type fakeProvider struct {
	board model.Board
	err   error
	calls int
}

func (p *fakeProvider) Fetch(context.Context) (model.Board, error) {
	p.calls++
	return p.board, p.err
}

type fakeCache struct {
	entries map[string]cache.Entry
	loadErr map[string]error
	saveErr error
	saves   []cache.Entry
}

func (c *fakeCache) Load(id string) (cache.Entry, bool, error) {
	if err := c.loadErr[id]; err != nil {
		return cache.Entry{}, false, err
	}
	entry, ok := c.entries[id]
	return entry, ok, nil
}

func (c *fakeCache) Save(entry cache.Entry) error {
	c.saves = append(c.saves, entry)
	if c.saveErr != nil {
		return c.saveErr
	}
	if c.entries == nil {
		c.entries = make(map[string]cache.Entry)
	}
	c.entries[entry.CatalogueID] = entry
	return nil
}

func selectedBoard(id, name string) selection.Board {
	return selection.Board{
		CatalogueID:      "masjidboardlive:" + id,
		Provider:         "masjidboardlive",
		ExternalID:       id,
		Name:             name,
		TimeZoneOffsetMS: 7200000,
	}
}

func board(id, name string, fajrHour int) model.Board {
	return model.Board{
		Identity: model.BoardIdentity{ID: id, Name: name, TimeZone: "GMT+02:00"},
		PrayerTimes: model.PrayerTimes{
			Fajr: model.PrayerTime{Jamaah: &model.ClockTime{Hour: fajrHour, Minute: 0}},
		},
	}
}

func TestFetchAllCurrentPersistsLastKnownGood(t *testing.T) {
	now := time.Date(2026, 8, 18, 19, 0, 0, 0, time.UTC)
	provider := &fakeProvider{board: board("brits-jamia", "Brits Jamia Masjid", 6)}
	store := &fakeCache{entries: make(map[string]cache.Entry)}
	coord, err := newWithClock([]Item{{Selection: selectedBoard("brits-jamia", "Brits Jamia Masjid"), Provider: provider}}, store, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}

	results := coord.FetchAll(context.Background())
	if len(results) != 1 {
		t.Fatalf("results = %d, want 1", len(results))
	}
	got := results[0]
	if got.Status != StatusCurrent || got.Board == nil || got.UpdateError != nil || got.PersistenceError != nil {
		t.Fatalf("result = %+v", got)
	}
	if !got.LastSuccessfulUpdate.Equal(now) || !got.LastAttempt.Equal(now) {
		t.Fatalf("timestamps = attempt %v success %v", got.LastAttempt, got.LastSuccessfulUpdate)
	}
	if len(store.saves) != 1 || store.saves[0].CatalogueID != "masjidboardlive:brits-jamia" || !store.saves[0].SuccessfulAt.Equal(now) {
		t.Fatalf("saved entries = %+v", store.saves)
	}
}

func TestFetchFailureUsesPersistedLastKnownGoodAsStale(t *testing.T) {
	attempt := time.Date(2026, 8, 18, 20, 0, 0, 0, time.UTC)
	success := attempt.Add(-2 * time.Hour)
	id := "masjidboardlive:brits-jamia"
	cachedBoard := board("brits-jamia", "Brits Jamia Masjid", 6)
	provider := &fakeProvider{err: errors.New("upstream unavailable")}
	store := &fakeCache{entries: map[string]cache.Entry{id: {CatalogueID: id, SuccessfulAt: success, Board: cachedBoard}}}
	coord, err := newWithClock([]Item{{Selection: selectedBoard("brits-jamia", "Brits Jamia Masjid"), Provider: provider}}, store, func() time.Time { return attempt })
	if err != nil {
		t.Fatal(err)
	}

	got := coord.FetchAll(context.Background())[0]
	if got.Status != StatusStale || got.Board == nil {
		t.Fatalf("result = %+v", got)
	}
	if got.UpdateError == nil || got.PersistenceError != nil {
		t.Fatalf("errors = update %v persistence %v", got.UpdateError, got.PersistenceError)
	}
	if !got.LastSuccessfulUpdate.Equal(success) || !got.LastAttempt.Equal(attempt) {
		t.Fatalf("timestamps = attempt %v success %v", got.LastAttempt, got.LastSuccessfulUpdate)
	}
	if len(store.saves) != 0 {
		t.Fatalf("failed refresh must not overwrite cache: saves=%d", len(store.saves))
	}
}

func TestFetchFailureWithoutCacheIsUnavailable(t *testing.T) {
	provider := &fakeProvider{err: errors.New("network down")}
	store := &fakeCache{entries: map[string]cache.Entry{}}
	coord, err := newWithClock([]Item{{Selection: selectedBoard("brits-jamia", "Brits Jamia Masjid"), Provider: provider}}, store, time.Now)
	if err != nil {
		t.Fatal(err)
	}
	got := coord.FetchAll(context.Background())[0]
	if got.Status != StatusUnavailable || got.Board != nil || got.UpdateError == nil {
		t.Fatalf("result = %+v", got)
	}
}

func TestFetchAllIsolatesFailuresAndPreservesOrder(t *testing.T) {
	now := time.Date(2026, 8, 18, 21, 0, 0, 0, time.UTC)
	cachedAt := now.Add(-time.Hour)
	first := &fakeProvider{board: board("brits-jamia", "Brits Jamia Masjid", 6)}
	second := &fakeProvider{err: errors.New("timeout")}
	third := &fakeProvider{board: board("brits-darul-uloom", "Jamiah Yusuf Darul Uloom Brits", 5)}
	store := &fakeCache{entries: map[string]cache.Entry{
		"masjidboardlive:brits-taqwa": {
			CatalogueID:  "masjidboardlive:brits-taqwa",
			SuccessfulAt: cachedAt,
			Board:        board("brits-taqwa", "Masjid Taqwa", 6),
		},
	}}
	items := []Item{
		{Selection: selectedBoard("brits-jamia", "Brits Jamia Masjid"), Provider: first},
		{Selection: selectedBoard("brits-taqwa", "Masjid Taqwa"), Provider: second},
		{Selection: selectedBoard("brits-darul-uloom", "Jamiah Yusuf Darul Uloom Brits"), Provider: third},
	}
	coord, err := newWithClock(items, store, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}

	results := coord.FetchAll(context.Background())
	if len(results) != 3 {
		t.Fatalf("results = %d, want 3", len(results))
	}
	wantIDs := []string{"brits-jamia", "brits-taqwa", "brits-darul-uloom"}
	wantStatuses := []Status{StatusCurrent, StatusStale, StatusCurrent}
	for i := range results {
		if results[i].Selection.ExternalID != wantIDs[i] || results[i].Status != wantStatuses[i] {
			t.Fatalf("result %d = id %q status %q", i, results[i].Selection.ExternalID, results[i].Status)
		}
	}
	if first.calls != 1 || second.calls != 1 || third.calls != 1 {
		t.Fatalf("provider calls = %d,%d,%d", first.calls, second.calls, third.calls)
	}
}

func TestSuccessfulFetchWithCacheSaveFailureRemainsCurrent(t *testing.T) {
	provider := &fakeProvider{board: board("brits-jamia", "Brits Jamia Masjid", 6)}
	store := &fakeCache{entries: map[string]cache.Entry{}, saveErr: errors.New("disk full")}
	coord, err := newWithClock([]Item{{Selection: selectedBoard("brits-jamia", "Brits Jamia Masjid"), Provider: provider}}, store, time.Now)
	if err != nil {
		t.Fatal(err)
	}
	got := coord.FetchAll(context.Background())[0]
	if got.Status != StatusCurrent || got.Board == nil || got.UpdateError != nil || got.PersistenceError == nil {
		t.Fatalf("result = %+v", got)
	}
}

func TestFetchFailureWithUnreadableCacheIsUnavailable(t *testing.T) {
	id := "masjidboardlive:brits-jamia"
	provider := &fakeProvider{err: errors.New("upstream unavailable")}
	store := &fakeCache{entries: map[string]cache.Entry{}, loadErr: map[string]error{id: errors.New("corrupt cache")}}
	coord, err := newWithClock([]Item{{Selection: selectedBoard("brits-jamia", "Brits Jamia Masjid"), Provider: provider}}, store, time.Now)
	if err != nil {
		t.Fatal(err)
	}
	got := coord.FetchAll(context.Background())[0]
	if got.Status != StatusUnavailable || got.Board != nil || got.UpdateError == nil || got.PersistenceError == nil {
		t.Fatalf("result = %+v", got)
	}
}

func TestNewAllowsUnconfiguredRuntime(t *testing.T) {
	coord, err := New(nil, nil)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if got := coord.FetchAll(context.Background()); len(got) != 0 {
		t.Fatalf("FetchAll() = %d results, want 0", len(got))
	}
}

func TestNewRejectsFourthBoard(t *testing.T) {
	store := &fakeCache{entries: map[string]cache.Entry{}}
	items := []Item{
		{Selection: selectedBoard("a", "A"), Provider: &fakeProvider{}},
		{Selection: selectedBoard("b", "B"), Provider: &fakeProvider{}},
		{Selection: selectedBoard("c", "C"), Provider: &fakeProvider{}},
		{Selection: selectedBoard("d", "D"), Provider: &fakeProvider{}},
	}
	if _, err := New(items, store); err == nil {
		t.Fatal("New() expected maximum-board error")
	}
}
