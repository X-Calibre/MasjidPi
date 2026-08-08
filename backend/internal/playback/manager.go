package playback

import (
	"context"
	"sync"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/player"
	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
	"github.com/X-Calibre/MasjidPi/backend/internal/version"
)

const (
	DefaultRetryInterval  = 30 * time.Second
	DefaultReconnectDelay = 5 * time.Second
	statusCheckInterval   = time.Second
	stateLoopInterval     = 50 * time.Millisecond
)

type State string

const (
	StateIdle       State = "idle"
	StateWaiting    State = "waiting"
	StateConnecting State = "connecting"
	StatePlaying    State = "playing"
	StateRetrying   State = "retrying"
	StateError      State = "error"
)

type Player interface {
	Play(url string) error
	Stop() error
	Volume(volume int) error
	Status() (*player.Status, error)
}

// Availability supplies push-based stream availability information.
// Implementations maintain the availability cache from LiveMasjid events;
// the playback manager only reads that local cache and never polls LiveMasjid.
type Availability interface {
	IsAvailable(mount string) (available bool, known bool)
	Events() <-chan string
}

type Config struct {
	RetryInterval       time.Duration
	ReconnectDelay      time.Duration
	StatusCheckInterval time.Duration
}

type Manager struct {
	mu        sync.Mutex
	startOnce sync.Once

	player       Player
	availability Availability

	retryInterval  time.Duration
	reconnectDelay time.Duration
	statusInterval time.Duration

	wake chan struct{}

	selected  *stream.Stream
	listening bool
	state     State
	lastError string
	status    Status
}

type Status struct {
	Version    string `json:"version"`
	State      string `json:"state"`
	Message    string `json:"message"`
	URL        string `json:"url"`
	Volume     int    `json:"volume"`
	Paused     bool   `json:"paused"`
	Listening  bool   `json:"listening"`
	StreamID   string `json:"stream_id,omitempty"`
	StreamName string `json:"stream_name,omitempty"`
	Error      string `json:"error,omitempty"`
}

func New(player Player, cfg Config) *Manager {
	if cfg.RetryInterval == 0 {
		cfg.RetryInterval = DefaultRetryInterval
	}
	if cfg.ReconnectDelay == 0 {
		cfg.ReconnectDelay = DefaultReconnectDelay
	}
	if cfg.StatusCheckInterval == 0 {
		cfg.StatusCheckInterval = statusCheckInterval
	}

	return &Manager{
		player:         player,
		retryInterval:  cfg.RetryInterval,
		reconnectDelay: cfg.ReconnectDelay,
		statusInterval: cfg.StatusCheckInterval,
		wake:            make(chan struct{}, 1),
		state:           StateIdle,
		status: Status{
			Version: version.Version,
			State:   string(StateIdle),
		},
	}
}

func (m *Manager) SetAvailability(availability Availability) {
	m.mu.Lock()
	m.availability = availability
	m.mu.Unlock()
	m.notify()
}

func (m *Manager) Start(ctx context.Context) {
	m.startOnce.Do(func() { go m.run(ctx) })
}

func (m *Manager) Play(selected stream.Stream) {
	m.mu.Lock()
	m.selected = &selected
	m.listening = true
	m.state = StateWaiting
	m.lastError = ""
	m.updateStatusLocked(nil)
	m.mu.Unlock()
	m.notify()
}

func (m *Manager) Stop() {
	m.mu.Lock()
	m.listening = false
	m.state = StateIdle
	m.lastError = ""
	m.updateStatusLocked(nil)
	m.mu.Unlock()
	m.notify()
}

func (m *Manager) Volume(volume int) error {
	if err := m.player.Volume(volume); err != nil {
		return err
	}
	m.mu.Lock()
	m.status.Volume = volume
	m.mu.Unlock()
	return nil
}

func (m *Manager) Status() Status {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.status
}

func (m *Manager) run(ctx context.Context) {
	loop := time.NewTicker(stateLoopInterval)
	defer loop.Stop()

	statusTicker := time.NewTicker(m.statusInterval)
	defer statusTicker.Stop()

	active := false
	playing := false
	nextAttempt := time.Time{}

	for {
		select {
		case <-ctx.Done():
			if active {
				_ = m.player.Stop()
			}
			return
		case <-m.wake:
			m.step(ctx, &active, &playing, &nextAttempt, false)
		case <-loop.C:
			m.step(ctx, &active, &playing, &nextAttempt, false)
		case <-statusTicker.C:
			if active {
				m.checkPlayerStatus(&active, &playing, &nextAttempt)
			}
		}
	}
}

func (m *Manager) step(ctx context.Context, active, playing *bool, nextAttempt *time.Time, force bool) {
	selected, listening, availability := m.snapshot()
	if ctx.Err() != nil {
		return
	}

	if !listening || selected == nil {
		if *active {
			_ = m.player.Stop()
			*active = false
			*playing = false
		}
		m.setState(StateIdle, "", nil)
		return
	}

	available := true
	known := false
	if availability != nil {
		available, known = availability.IsAvailable(selected.ID)
		if !known || !available {
			if *active {
				_ = m.player.Stop()
				*active = false
				*playing = false
			}
			m.setState(StateWaiting, "", nil)
			return
		}
	}

	if *active {
		return
	}
	if !force && !nextAttempt.IsZero() && time.Now().Before(*nextAttempt) {
		m.setState(StateRetrying, "", nil)
		return
	}

	m.setState(StateConnecting, "", nil)
	if err := m.player.Play(selected.URL); err != nil {
		m.setState(StateRetrying, err.Error(), nil)
		*nextAttempt = time.Now().Add(m.retryInterval)
		return
	}

	*active = true
	*playing = false
	*nextAttempt = time.Time{}
	m.setState(StatePlaying, "", nil)

	_ = available
	_ = known
}

func (m *Manager) checkPlayerStatus(active, playing *bool, nextAttempt *time.Time) {
	if !*active {
		return
	}

	status, err := m.player.Status()
	if err != nil {
		m.setState(StateRetrying, err.Error(), nil)
		_ = m.player.Stop()
		*active = false
		*playing = false
		*nextAttempt = time.Now().Add(m.reconnectDelay)
		return
	}

	selected, listening, availability := m.snapshot()
	if !listening || selected == nil {
		return
	}
	if availability != nil {
		if available, known := availability.IsAvailable(selected.ID); known && !available {
			_ = m.player.Stop()
			*active = false
			*playing = false
			*nextAttempt = time.Time{}
			m.setState(StateWaiting, "", status)
			return
		}
	}

	if status.State == "stopped" {
		_ = m.player.Stop()
		*active = false
		*playing = false
		*nextAttempt = time.Now().Add(m.reconnectDelay)
		m.setState(StateRetrying, "", status)
		return
	}

	*playing = true
	m.setState(StatePlaying, "", status)
}

func (m *Manager) snapshot() (*stream.Stream, bool, Availability) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.selected, m.listening, m.availability
}

func (m *Manager) setState(state State, message string, playerStatus *player.Status) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state = state
	m.lastError = message
	m.updateStatusLocked(playerStatus)
}

func (m *Manager) updateStatusLocked(playerStatus *player.Status) {
	status := m.status
	status.State = string(m.state)
	switch m.state {
	case StateIdle:
		status.Message = "Idle"
	case StateWaiting:
		status.Message = "Waiting for stream"
	case StateConnecting:
		status.Message = "Connecting..."
	case StatePlaying:
		status.Message = "Playing"
	case StateRetrying:
		status.Message = "Reconnecting"
	case StateError:
		status.Message = "Playback error"
	default:
		status.Message = ""
	}
	status.Listening = m.listening
	status.Error = m.lastError
	if m.selected == nil {
		status.StreamID = ""
		status.StreamName = ""
		status.URL = ""
	} else {
		status.StreamID = m.selected.ID
		status.StreamName = m.selected.Name
		status.URL = m.selected.URL
	}
	if playerStatus != nil {
		status.URL = playerStatus.URL
		status.Volume = playerStatus.Volume
		status.Paused = playerStatus.Paused
	}
	m.status = status
}

func (m *Manager) notify() {
	select {
	case m.wake <- struct{}{}:
	default:
	}
}
