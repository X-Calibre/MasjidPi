package runtime

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/cache"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/provider"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

// Status describes whether the displayed timetable is current, stale but
// usable, or unavailable because neither a live fetch nor cached data exists.
type Status string

const (
	StatusCurrent     Status = "current"
	StatusStale       Status = "stale"
	StatusUnavailable Status = "unavailable"
)

// Item binds one selected board to its provider. Items are processed in order,
// preserving the user's selected-board display order.
type Item struct {
	Selection selection.Board
	Provider  provider.Provider
}

// Result is the independent runtime state for one selected board.
type Result struct {
	Selection            selection.Board
	Board                *model.Board
	Status               Status
	LastAttempt          time.Time
	LastSuccessfulUpdate time.Time
	UpdateError          error
	PersistenceError     error
}

// CacheStore is the last-known-good persistence contract used by Coordinator.
// cache.Store satisfies this interface.
type CacheStore interface {
	Load(catalogueID string) (cache.Entry, bool, error)
	Save(entry cache.Entry) error
}

// Coordinator refreshes selected boards independently and falls back to each
// board's last-known-good persisted timetable when a live refresh fails.
type Coordinator struct {
	items []Item
	cache CacheStore
	now   func() time.Time
}

func New(items []Item, cacheStore CacheStore) (*Coordinator, error) {
	return newWithClock(items, cacheStore, time.Now)
}

func newWithClock(items []Item, cacheStore CacheStore, now func() time.Time) (*Coordinator, error) {
	if len(items) == 0 {
		return &Coordinator{cache: cacheStore, now: now}, nil
	}
	if len(items) > selection.MaxBoards {
		return nil, fmt.Errorf("masjidboard runtime: %d items configured; maximum is %d", len(items), selection.MaxBoards)
	}
	if cacheStore == nil {
		return nil, fmt.Errorf("masjidboard runtime: cache store is required")
	}
	if now == nil {
		return nil, fmt.Errorf("masjidboard runtime: clock is required")
	}

	state := selection.State{Boards: make([]selection.Board, len(items))}
	for i, item := range items {
		if item.Provider == nil {
			return nil, fmt.Errorf("masjidboard runtime: provider %d is required", i+1)
		}
		state.Boards[i] = item.Selection
	}
	if err := selection.Validate(state); err != nil {
		return nil, fmt.Errorf("masjidboard runtime: invalid selection: %w", err)
	}

	copyItems := append([]Item(nil), items...)
	return &Coordinator{items: copyItems, cache: cacheStore, now: now}, nil
}

// FetchAll refreshes every selected board independently. One board's failure
// never prevents later boards from being attempted or returned.
func (c *Coordinator) FetchAll(ctx context.Context) []Result {
	results := make([]Result, len(c.items))
	for i, item := range c.items {
		results[i] = c.fetchOne(ctx, item)
	}
	return results
}

func (c *Coordinator) fetchOne(ctx context.Context, item Item) Result {
	attemptedAt := c.now()
	result := Result{
		Selection:   item.Selection,
		Status:      StatusUnavailable,
		LastAttempt: attemptedAt,
	}

	board, err := item.Provider.Fetch(ctx)
	if err == nil {
		result.Board = &board
		result.Status = StatusCurrent
		result.LastSuccessfulUpdate = attemptedAt

		entry := cache.Entry{
			CatalogueID:  item.Selection.CatalogueID,
			SuccessfulAt: attemptedAt,
			Board:        board,
		}
		if saveErr := c.cache.Save(entry); saveErr != nil {
			result.PersistenceError = saveErr
		}
		return result
	}

	result.UpdateError = err
	entry, found, cacheErr := c.cache.Load(item.Selection.CatalogueID)
	if cacheErr != nil {
		result.PersistenceError = cacheErr
		result.UpdateError = errors.Join(err, fmt.Errorf("load last-known-good cache: %w", cacheErr))
		return result
	}
	if !found {
		return result
	}

	board = entry.Board
	result.Board = &board
	result.Status = StatusStale
	result.LastSuccessfulUpdate = entry.SuccessfulAt
	return result
}
