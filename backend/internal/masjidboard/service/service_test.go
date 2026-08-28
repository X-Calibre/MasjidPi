package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/cache"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/provider"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

type fakeProvider struct {
	board model.Board
	err   error
}

func (p fakeProvider) Fetch(context.Context) (model.Board, error) {
	return p.board, p.err
}

type memoryCache struct {
	mu      sync.RWMutex
	entries map[string]cache.Entry
}

func newMemoryCache() *memoryCache {
	return &memoryCache{entries: make(map[string]cache.Entry)}
}

func (c *memoryCache) Load(id string) (cache.Entry, bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[id]
	return entry, ok, nil
}

func (c *memoryCache) Save(entry cache.Entry) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[entry.CatalogueID] = entry
	return nil
}

func serviceBoard(id, name string) selection.Board {
	return selection.Board{
		CatalogueID:      "masjidboardlive:" + id,
		Provider:         "masjidboardlive",
		ExternalID:       id,
		Name:             name,
		TimeZoneOffsetMS: 7200000,
	}
}

func liveBoard(id, name string) model.Board {
	return model.Board{Identity: model.BoardIdentity{ID: id, Name: name, TimeZone: "GMT+02:00"}}
}

func TestUnconfiguredServiceStartsWithoutProviders(t *testing.T) {
	service, err := newWithFactory(selection.State{}, nil, nil)
	if err != nil {
		t.Fatalf("newWithFactory() error = %v", err)
	}
	if service.Configured() {
		t.Fatal("Configured() = true, want false")
	}
	if got := service.Refresh(context.Background()); len(got) != 0 {
		t.Fatalf("Refresh() returned %d results, want 0", len(got))
	}
}

func TestServiceConstructsProvidersInSelectionOrder(t *testing.T) {
	state := selection.State{Boards: []selection.Board{
		serviceBoard("brits-taqwa", "Masjid Taqwa"),
		serviceBoard("brits-jamia", "Brits Jamia Masjid"),
		serviceBoard("brits-darul-uloom", "Jamiah Yusuf Darul Uloom Brits"),
	}}

	var constructed []string
	service, err := newWithFactory(state, newMemoryCache(), func(board selection.Board) (provider.Provider, error) {
		constructed = append(constructed, board.ExternalID)
		return fakeProvider{board: liveBoard(board.ExternalID, board.Name)}, nil
	})
	if err != nil {
		t.Fatalf("newWithFactory() error = %v", err)
	}
	if len(constructed) != 3 || constructed[0] != "brits-taqwa" || constructed[1] != "brits-jamia" || constructed[2] != "brits-darul-uloom" {
		t.Fatalf("constructed order = %v", constructed)
	}

	results := service.Refresh(context.Background())
	if len(results) != 3 {
		t.Fatalf("Refresh() returned %d results, want 3", len(results))
	}
	for i := range state.Boards {
		if results[i].Selection.ExternalID != state.Boards[i].ExternalID {
			t.Fatalf("result %d selection = %q, want %q", i, results[i].Selection.ExternalID, state.Boards[i].ExternalID)
		}
		if results[i].Status != runtime.StatusCurrent {
			t.Fatalf("result %d status = %q, want current", i, results[i].Status)
		}
	}
}

func TestServiceKeepsBoardFailuresIndependent(t *testing.T) {
	state := selection.State{Boards: []selection.Board{
		serviceBoard("one", "One"),
		serviceBoard("two", "Two"),
		serviceBoard("three", "Three"),
	}}
	cacheStore := newMemoryCache()
	if err := cacheStore.Save(cache.Entry{
		CatalogueID:  state.Boards[1].CatalogueID,
		SuccessfulAt: time.Date(2026, 8, 18, 15, 0, 0, 0, time.UTC),
		Board:        liveBoard("two", "Two cached"),
	}); err != nil {
		t.Fatal(err)
	}

	service, err := newWithFactory(state, cacheStore, func(board selection.Board) (provider.Provider, error) {
		switch board.ExternalID {
		case "one":
			return fakeProvider{board: liveBoard("one", "One")}, nil
		case "two":
			return fakeProvider{err: errors.New("upstream failed")}, nil
		case "three":
			return fakeProvider{err: errors.New("offline")}, nil
		default:
			return nil, errors.New("unexpected board")
		}
	})
	if err != nil {
		t.Fatalf("newWithFactory() error = %v", err)
	}

	results := service.Refresh(context.Background())
	if results[0].Status != runtime.StatusCurrent {
		t.Fatalf("board one status = %q", results[0].Status)
	}
	if results[1].Status != runtime.StatusStale || results[1].Board == nil || results[1].Board.Identity.Name != "Two cached" {
		t.Fatalf("board two result = %+v", results[1])
	}
	if results[2].Status != runtime.StatusUnavailable || results[2].Board != nil {
		t.Fatalf("board three result = %+v", results[2])
	}
}

func TestServicePublishesLatestResults(t *testing.T) {
	state := selection.State{Boards: []selection.Board{serviceBoard("brits-jamia", "Brits Jamia Masjid")}}
	service, err := newWithFactory(state, newMemoryCache(), func(board selection.Board) (provider.Provider, error) {
		return fakeProvider{board: liveBoard(board.ExternalID, board.Name)}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(service.Results()) != 0 {
		t.Fatal("Results() before refresh should be empty")
	}
	service.Refresh(context.Background())
	got := service.Results()
	if len(got) != 1 || got[0].Status != runtime.StatusCurrent {
		t.Fatalf("Results() = %+v", got)
	}
	got[0].Status = runtime.StatusUnavailable
	if service.Results()[0].Status != runtime.StatusCurrent {
		t.Fatal("Results() returned mutable service slice")
	}
}

func TestServiceRejectsProviderConstructionFailure(t *testing.T) {
	state := selection.State{Boards: []selection.Board{serviceBoard("brits-jamia", "Brits Jamia Masjid")}}
	_, err := newWithFactory(state, newMemoryCache(), func(selection.Board) (provider.Provider, error) {
		return nil, errors.New("bad provider")
	})
	if err == nil {
		t.Fatal("newWithFactory() expected provider construction error")
	}
}
