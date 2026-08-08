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
		wake:           make(chan struct{}, 1),
		state:          StateIdle,
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
	active := false

	for {
		selected, listening, availability := m.snapshot()

		if ctx.Err() != nil {
			if active {
				_ = m.player.Stop()
			}
			return
		}

		if !listening || selected == nil {
			if active {
				_ = m.player.Stop()
				active = false
			}
			m.setState(StateIdle, "", nil)
			if !m.wait(ctx, 0) {
				return
			}
			continue
		}

		if availability != nil {
			available, known := availability.IsAvailable(selected.ID)
			if !known || !available {
				if active {
					_ = m.player.Stop()
					active = false
				}
				m.setState(StateWaiting, "", nil)
				if !m.waitForAvailability(ctx, selected.ID, availability) {
					return
				}
				continue
			}
		}

		m.setState(StateConnecting, "", nil)
		if err := m.player.Play(selected.URL); err != nil {
			m.setState(StateRetrying, err.Error(), nil)
			if !m.wait(ctx, m.retryInterval) {
				return
			}
			continue
		}

		active = true
		delay, ok := m.monitor(ctx, selected.ID, availability)
		if !ok {
			return
		}
		if delay > 0 && !m.wait(ctx, delay) {
			return
		}
	}
}

func (m *Manager) monitor(ctx context.Context, mount string, availability Availability) (time.Duration, bool) {
	ticker := time.NewTicker(m.statusInterval)
	defer ticker.Stop()

	wasPlaying := false
	for {
		select {
		case <-ctx.Done():
			_ = m.player.Stop()
			return 0, false
		case <-m.wake:
			return 0, true
		case <-ticker.C:
			if availability != nil {
				if available, known := availability.IsAvailable(mount); known && !available {
					m.setState(StateWaiting, "", nil)
					return 0, m.waitForAvailability(ctx, mount, availability)
				}
			}

			status, err := m.player.Status()
			if err != nil {
				m.setState(StateRetrying, err.Error(), nil)
				return m.retryDelay(wasPlaying), true
			}
			if status.State == "stopped" {
				if availability != nil {
					if available, known := availability.IsAvailable(mount); known && !available {
						m.setState(StateWaiting, "", status)
						return 0, m.waitForAvailability(ctx, mount, availability)
					}
				}
				m.setState(StateRetrying, "", status)
				return m.retryDelay(wasPlaying), true
			}

			wasPlaying = true
			m.setState(StatePlaying, "", status)
		}
	}
}

func (m *Manager) waitForAvailability(ctx context.Context, mount string, availability Availability) bool {
	events := availability.Events()
	check := func() bool {
		available, known := availability.IsAvailable(mount)
		return known && available
	}

	for {
		if check() {
			return true
		}

		// Events are the normal wake-up mechanism. The short local timer is
		// only a safety net against a lost/coalesced event; it never contacts
		// LiveMasjid or performs any network polling.
		timer := time.NewTimer(100 * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return false
		case eventMount, ok := <-events:
			timer.Stop()
			if !ok {
				return false
			}
			if eventMount == mount && check() {
				return true
			}
		case <-timer.C:
		}
	}
}

func (m *Manager) retryDelay(wasPlaying bool) time.Duration {
	if wasPlaying {
		return m.reconnectDelay
	}
	return m.retryInterval
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

func (m *Manager) wait(ctx context.Context, delay time.Duration) bool {
	if delay <= 0 {
		select {
		case <-ctx.Done():
			return false
		case <-m.wake:
			return true
		}
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-m.wake:
		return true
	case <-timer.C:
		return true
	}
}
