package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/cache"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/provider"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/provider/masjidboardlive"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

// Config contains the persistent paths required to start the MasjidBoard
// runtime. The full discovery catalogue is deliberately not required here.
type Config struct {
	SelectionPath string
	CacheDir      string
	HTTPClient    *http.Client
}

// Service owns the small runtime-critical MasjidBoard state. It loads the
// configured 1-3 board selection once at startup, constructs independent
// providers and exposes the latest per-board runtime results.
type Service struct {
	selection selection.State
	runtime   *runtime.Coordinator

	mu      sync.RWMutex
	results []runtime.Result
}

// New constructs the production MasjidBoard startup service.
func New(config Config) (*Service, error) {
	selectionStore := selection.NewStore(config.SelectionPath)
	state, err := selectionStore.Load()
	if err != nil {
		return nil, fmt.Errorf("masjidboard service: load selection: %w", err)
	}

	cacheStore := cache.NewStore(config.CacheDir)
	return newWithFactory(state, cacheStore, func(board selection.Board) (provider.Provider, error) {
		client, err := masjidboardlive.NewCoreClientFromSelectionWithHTTPClient(board, config.HTTPClient)
		if err != nil {
			return nil, err
		}
		return client, nil
	})
}

type providerFactory func(selection.Board) (provider.Provider, error)

func newWithFactory(state selection.State, cacheStore runtime.CacheStore, factory providerFactory) (*Service, error) {
	if !state.Configured() {
		coordinator, err := runtime.New(nil, cacheStore)
		if err != nil {
			return nil, err
		}
		return &Service{selection: state, runtime: coordinator}, nil
	}
	if err := selection.Validate(state); err != nil {
		return nil, fmt.Errorf("masjidboard service: invalid selection: %w", err)
	}
	if cacheStore == nil {
		return nil, fmt.Errorf("masjidboard service: cache store is required")
	}
	if factory == nil {
		return nil, fmt.Errorf("masjidboard service: provider factory is required")
	}

	items := make([]runtime.Item, 0, len(state.Boards))
	for i, board := range state.Boards {
		p, err := factory(board)
		if err != nil {
			return nil, fmt.Errorf("masjidboard service: construct provider %d: %w", i+1, err)
		}
		items = append(items, runtime.Item{Selection: board, Provider: p})
	}

	coordinator, err := runtime.New(items, cacheStore)
	if err != nil {
		return nil, fmt.Errorf("masjidboard service: create runtime: %w", err)
	}
	return &Service{selection: selection.State{Boards: append([]selection.Board(nil), state.Boards...)}, runtime: coordinator}, nil
}

// Configured reports whether MasjidBoard has a persisted user selection.
func (s *Service) Configured() bool {
	return s != nil && s.selection.Configured()
}

// Selection returns the startup selection in user-defined display order.
func (s *Service) Selection() selection.State {
	if s == nil {
		return selection.State{}
	}
	return selection.State{Boards: append([]selection.Board(nil), s.selection.Boards...)}
}

// Refresh updates all configured boards independently and publishes one result
// per selected board. A failure on one board never suppresses the others.
func (s *Service) Refresh(ctx context.Context) []runtime.Result {
	if s == nil || s.runtime == nil {
		return nil
	}
	results := s.runtime.FetchAll(ctx)

	s.mu.Lock()
	s.results = append([]runtime.Result(nil), results...)
	s.mu.Unlock()

	return append([]runtime.Result(nil), results...)
}

// Results returns the latest published runtime results. Before the first
// refresh it returns an empty slice.
func (s *Service) Results() []runtime.Result {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]runtime.Result(nil), s.results...)
}
