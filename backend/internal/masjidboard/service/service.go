package service

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/cache"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/economic"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/provider"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/provider/masjidboardlive"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

type Config struct {
	SelectionPath string
	CacheDir      string
	HTTPClient    *http.Client
}

type Service struct {
	mu             sync.RWMutex
	selection      selection.State
	runtime        *runtime.Coordinator
	results        []runtime.Result
	selectionStore *selection.Store
	cacheStore     runtime.CacheStore
	factory        providerFactory
	economicClient economic.Client
	economicStore  economic.Store
	indicators     *economic.Indicators
}

func New(config Config) (*Service, error) {
	selectionStore := selection.NewStore(config.SelectionPath)
	state, err := selectionStore.Load()
	if err != nil {
		return nil, fmt.Errorf("masjidboard service: load selection: %w", err)
	}
	cacheStore := cache.NewStore(config.CacheDir)
	economicStore := economic.Store{Path: filepath.Join(config.CacheDir, "islamic_economic_indicators.json")}
	indicators, err := economicStore.Load()
	if err != nil {
		// A corrupt optional cache must not prevent MasjidBoard or Listen from
		// starting. The next successful fetch will replace it atomically.
		indicators = nil
	}
	factory := func(board selection.Board) (provider.Provider, error) {
		core, err := masjidboardlive.NewCoreClientFromSelectionWithHTTPClient(board, config.HTTPClient)
		if err != nil {
			return nil, err
		}
		return masjidboardlive.EnrichedClient{
			Core: core,
			Premium: masjidboardlive.PremiumClient{
				HTTPClient: config.HTTPClient,
				Mid:        strings.TrimSpace(board.ExternalID),
			},
		}, nil
	}
	service, err := newWithFactory(state, cacheStore, factory)
	if err != nil {
		return nil, err
	}
	service.selectionStore = selectionStore
	service.economicClient = economic.Client{HTTPClient: config.HTTPClient}
	service.economicStore = economicStore
	service.indicators = indicators
	return service, nil
}

type providerFactory func(selection.Board) (provider.Provider, error)

func newWithFactory(state selection.State, cacheStore runtime.CacheStore, factory providerFactory) (*Service, error) {
	coordinator, err := buildCoordinator(state, cacheStore, factory)
	if err != nil {
		return nil, err
	}
	return &Service{selection: cloneSelection(state), runtime: coordinator, cacheStore: cacheStore, factory: factory}, nil
}

func buildCoordinator(state selection.State, cacheStore runtime.CacheStore, factory providerFactory) (*runtime.Coordinator, error) {
	if !state.Configured() {
		return runtime.New(nil, cacheStore)
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
	return coordinator, nil
}

func (s *Service) Reconfigure(state selection.State) error {
	if s == nil {
		return fmt.Errorf("masjidboard service: service is unavailable")
	}
	if err := selection.Validate(state); err != nil {
		return err
	}
	s.mu.RLock()
	cacheStore, factory, selectionStore := s.cacheStore, s.factory, s.selectionStore
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
	s.selection, s.runtime, s.results = cloneSelection(state), coordinator, nil
	s.mu.Unlock()
	return nil
}

func (s *Service) SetSlideDurationSeconds(seconds int) error {
	if seconds < selection.MinSlideDurationSeconds || seconds > selection.MaxSlideDurationSeconds {
		return fmt.Errorf("masjidboard service: slide duration must be between %d and %d seconds", selection.MinSlideDurationSeconds, selection.MaxSlideDurationSeconds)
	}
	return s.updateDisplayPreference(func(state *selection.State) { state.SlideDurationSeconds = seconds }, "slide duration")
}

func (s *Service) SetTheme(theme string) error {
	theme = strings.TrimSpace(strings.ToLower(theme))
	if !selection.ThemeSupported(theme) {
		return fmt.Errorf("masjidboard service: unsupported display theme %q", theme)
	}
	return s.updateDisplayPreference(func(state *selection.State) { state.Theme = theme }, "theme")
}

func (s *Service) SetShowEconomicIndicators(show bool) error {
	return s.updateDisplayPreference(func(state *selection.State) { state.ShowEconomicIndicators = show }, "economic indicators")
}

func (s *Service) updateDisplayPreference(update func(*selection.State), label string) error {
	if s == nil {
		return fmt.Errorf("masjidboard service: service is unavailable")
	}
	s.mu.RLock()
	state, store := cloneSelection(s.selection), s.selectionStore
	s.mu.RUnlock()
	if !state.Configured() {
		return fmt.Errorf("masjidboard service: select at least one board before choosing a display %s", label)
	}
	if store == nil {
		return fmt.Errorf("masjidboard service: selection store is unavailable")
	}
	update(&state)
	if err := selection.Validate(state); err != nil {
		return err
	}
	if err := store.Save(state); err != nil {
		return fmt.Errorf("masjidboard service: persist %s: %w", label, err)
	}
	s.mu.Lock()
	s.selection = cloneSelection(state)
	s.mu.Unlock()
	return nil
}

func (s *Service) Configured() bool {
	if s == nil {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.selection.Configured()
}

func (s *Service) Selection() selection.State {
	if s == nil {
		return selection.State{}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneSelection(s.selection)
}

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
	_ = s.RefreshEconomicIndicators(ctx)
	s.mu.Lock()
	if s.runtime == coordinator {
		s.results = append([]runtime.Result(nil), results...)
	}
	s.mu.Unlock()
	return append([]runtime.Result(nil), results...)
}

const economicRefreshHour = 9

var economicRefreshLocation = time.FixedZone("Africa/Johannesburg", 2*60*60)

func economicRefreshDue(current *economic.Indicators, now time.Time) bool {
	if current == nil {
		return true
	}
	if !current.Complete() {
		return true
	}
	localNow := now.In(economicRefreshLocation)
	if localNow.Hour() < economicRefreshHour {
		return false
	}
	if localNow.Weekday() == time.Saturday || localNow.Weekday() == time.Sunday {
		weekendCutoff := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 10, 30, 0, 0, economicRefreshLocation)
		if !localNow.Before(weekendCutoff) {
			return false
		}
	}
	effectiveDate, err := time.Parse("2006-01-02", current.EffectiveDate)
	if err != nil {
		return true
	}
	return effectiveDate.Year() != localNow.Year() || effectiveDate.YearDay() != localNow.YearDay()
}

func (s *Service) RefreshEconomicIndicators(ctx context.Context) error {
	s.mu.RLock()
	enabled := s.selection.ShowEconomicIndicators
	current := s.indicators
	client, store := s.economicClient, s.economicStore
	s.mu.RUnlock()
	now := time.Now
	if client.Now != nil {
		now = client.Now
	}
	if !enabled || !economicRefreshDue(current, now()) {
		return nil
	}
	indicators, err := client.Fetch(ctx)
	if err != nil {
		return err
	}
	if current != nil {
		currentDate, currentErr := time.Parse("2006-01-02", current.EffectiveDate)
		fetchedDate, fetchedErr := time.Parse("2006-01-02", indicators.EffectiveDate)
		if currentErr == nil && fetchedErr == nil && !fetchedDate.After(currentDate) && current.Complete() {
			return nil
		}
	}
	if err := store.Save(indicators); err != nil {
		return err
	}
	s.mu.Lock()
	s.indicators = &indicators
	s.mu.Unlock()
	return nil
}

func (s *Service) EconomicIndicators() *economic.Indicators {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.selection.ShowEconomicIndicators || s.indicators == nil {
		return nil
	}
	copy := *s.indicators
	return &copy
}

func (s *Service) Results() []runtime.Result {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]runtime.Result(nil), s.results...)
}

func cloneSelection(state selection.State) selection.State {
	copy := state
	copy.Boards = append([]selection.Board(nil), state.Boards...)
	return copy
}
