package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
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
// providers and exposes the latest per-board runtime results. Selection may be
// reconfigured live through the API without restarting the application.
type Service struct {
	mu sync.RWMutex

	selection selection.State
	runtime   *runtime.Coordinator
	results   []runtime.Result

	selectionStore *selection.Store
	cacheStore     runtime.CacheStore
	factory        providerFactory
}

// New constructs the production MasjidBoard startup service.
func New(config Config) (*Service, error) {
	selectionStore := selection.NewStore(config.SelectionPath)
	state, err := selectionStore.Load()
	if err != nil {
		return nil, fmt.Errorf("masjidboard service: load selection: %w", err)
	}

	cacheStore := cache.NewStore(config.CacheDir)
	factory := func(board selection.Board) (provider.Provider, error) {
		client, err := masjidboardlive.NewCoreClientFromSelectionWithHTTPClient(board, config.HTTPClient)
		if err != nil {
			return nil, err
		}
		return client, nil
	}

	service, err := newWithFactory(state, cacheStore, factory)
	if err != nil {
		return nil, err
	}
	service.selectionStore = selectionStore
	return service, nil
}

type providerFactory func(selection.Board) (provider.Provider, error)

func newWithFactory(state selection.State, cacheStore runtime.CacheStore, factory providerFactory) (*Service, error) {
	coordinator, err := buildCoordinator(state, cacheStore, factory)
	if err != nil {
		return nil, err
	}
	return &Service{
		selection:  cloneSelection(state),
		runtime:    coordinator,
		cacheStore: cacheStore,
		factory:    factory,
	}, nil
}

func buildCoordinator(state selection.State, cacheStore runtime.CacheStore, factory providerFactory) (*runtime.Coordinator, error) {
	if !state.Configured() {
		coordinator, err := runtime.New(nil, cacheStore)
		if err != nil {
			return nil, err
		}
		return coordinator, nil
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
	for i, board := range state.Boards) {
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
	return coordinator, nil
}

// Reconfigure validates and persists a new ordered 1-3 board selection, then
// atomically replaces the active runtime coordinator. Existing cached board
// data is retained on disk and may be reused by the new coordinator.
func (s *Service) Reconfigure(state selection.State) error {
	if s == nil {
		return fmt.Errorf("masjidboard service: service is unavailable")
	}
	if err := selection.Validate(state); err != nil {
		return err
	}

	s.mu.RLock()
	cacheStore := s.cacheStore
	factory := s.factory
	selectionStore := s.selectionStore
	s.mu.RUnlock()

	coordinator, err := buildCoordinator(state, cacheStore, factory)
	if err != nil {
		return err
	}
	if selectionStore == nil {
		return fmt.Errorf("masjidboard service: selection store is unavailable")
	}
	if err := selectionStore.Save(state); err != nil {
		return fmt.Errorf("masjidboard service: persist selection: %w", err)
	}

	s.mu.Lock()
	s.selection = cloneSelection(state)
	s.runtime = coordinator
	s.results = nil
	s.mu.Unlock()
	return nil
}

// SetLayout persists only the HDMI/display presentation preference. It leaves
// the active runtime coordinator and current timetable results untouched.
func (s *Service) SetLayout(layout string) error {
	if s == nil {
		return fmt.Errorf("masjidboard service: service is unavailable")
	}
	layout = strings.TrimSpace(strings.ToLower(layout))
	if layout != selection.LayoutStandard && layout != selection.LayoutDetailed {
		return fmt.Errorf("masjidboard service: unsupported display layout %q", layout)
	}

	s.mu.RLock()
	state := cloneSelection(s.selection)
	selectionStore := s.selectionStore
	s.mu.RUnlock()
	if !state.Configured() {
		return fmt.Errorf("masjidboard service: select at least one board before choosing a display layout")
	}
	if selectionStore == nil {
		return fmt.Errorf("masjidboard service: selection store is unavailable")
	}
	state.Layout = layout
	if err := selection.Validate(state); err != nil {
		return err
	}
	if err := selectionStore.Save(state); err != nil {
		return fmt.Errorf("masjidboard service: persist layout: %w", err)
	}

	s.mu.Lock()
	s.selection.Layout = layout
	s.mu.Unlock()
	return nil
}

// Configured reports whether MasjidBoard has a persisted user selection.
func (s *Service) Configured() bool {
	if s == nil {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.selection.Configured()
}

// Selection returns the current selection in user-defined display order.
func (s *Service) Selection() selection.State {
	if s == nil {
		return selection.State{}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneSelection(s.selection)
}

// Refresh updates all configured boards independently and publishes one result
// per selected board. A failure on one board never suppresses the others.
func (s *Service) Refresh(ctx context.Context) []runtime.Result {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	coordinator := s.runtime
	s.mu.RUnlock()
	if coordinator == nil {
		return nil
	}

	results := coordinator.FetchAll(ctx)

	s.mu.Lock()
	// Do not publish results from a coordinator that was replaced while this
	// network refresh was in flight.
	if s.runtime == coordinator {
		s.results = append([]runtime.Result(nil), results...)
	}
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

func cloneSelection(state selection.State) selection.State {
	return selection.State{
		Boards: append([]selection.Board(nil), state.Boards...),
		Layout: state.Layout,
	}
}
