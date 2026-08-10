package playback

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/player"
	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
	"github.com/X-Calibre/MasjidPi/backend/internal/version"
)

const (
	DefaultRetryInterval      = 5 * time.Second
	DefaultReconnectDelay     = 5 * time.Second
	DefaultStartupGracePeriod = 10 * time.Second
	statusCheckInterval       = time.Second
	stateLoopInterval         = 50 * time.Millisecond
	maxRetryDelay             = 5 * time.Minute
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
	AudioDevices() ([]player.AudioDevice, error)
	AudioDevice(name string) error
	Status() (*player.Status, error)
}

type Availability interface {
	IsAvailable(mount string) (available bool, known bool)
	Events() <-chan string
}

type Persistence interface {
	Save(streamID string) error
	Clear() error
}

type Config struct {
	RetryInterval       time.Duration
	ReconnectDelay      time.Duration
	StartupGracePeriod  time.Duration
	StatusCheckInterval time.Duration
	Logger              *slog.Logger
}

type Manager struct {
	mu        sync.Mutex
	startOnce sync.Once

	player       Player
	availability Availability
	persistence  Persistence
	log          *slog.Logger

	retryInterval      time.Duration
	reconnectDelay     time.Duration
	startupGracePeriod time.Duration
	statusInterval     time.Duration

	wake chan struct{}

	selected  *stream.Stream
	listening bool
	state     State
	lastError string
	status    Status
}

type Status struct {
	Version     string `json:"version"`
	State       string `json:"state"`
	Message     string `json:"message"`
	URL         string `json:"url"`
	Volume      int    `json:"volume"`
	Paused      bool   `json:"paused"`
	AudioDevice string `json:"audio_device,omitempty"`
	Listening   bool   `json:"listening"`
	StreamID    string `json:"stream_id,omitempty"`
	StreamName  string `json:"stream_name,omitempty"`
	Error       string `json:"error,omitempty"`
}

func New(player Player, cfg Config) *Manager {
	if cfg.RetryInterval == 0 {
		cfg.RetryInterval = DefaultRetryInterval
	}
	if cfg.ReconnectDelay == 0 {
		cfg.ReconnectDelay = DefaultReconnectDelay
	}
	if cfg.StartupGracePeriod == 0 {
		cfg.StartupGracePeriod = DefaultStartupGracePeriod
	}
	if cfg.StatusCheckInterval == 0 {
		cfg.StatusCheckInterval = statusCheckInterval
	}

	return &Manager{
		player:             player,
		log:                cfg.Logger,
		retryInterval:      cfg.RetryInterval,
		reconnectDelay:     cfg.ReconnectDelay,
		startupGracePeriod: cfg.StartupGracePeriod,
		statusInterval:     cfg.StatusCheckInterval,
		wake:               make(chan struct{}, 1),
		state:              StateIdle,
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

func (m *Manager) SetPersistence(persistence Persistence) {
	m.mu.Lock()
	m.persistence = persistence
	m.mu.Unlock()
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
	persistence := m.persistence
	m.updateStatusLocked(nil)
	m.mu.Unlock()

	if persistence != nil {
		if err := persistence.Save(selected.ID); err != nil {
			m.logPersistenceError("saving last playback stream", err)
		}
	}
	m.notify()
}

func (m *Manager) Stop() {
	m.mu.Lock()
	m.listening = false
	m.state = StateIdle
	m.lastError = ""
	persistence := m.persistence
	m.updateStatusLocked(nil)
	m.mu.Unlock()

	if persistence != nil {
		if err := persistence.Clear(); err != nil {
			m.logPersistenceError("clearing last playback stream", err)
		}
	}
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

func (m *Manager) AudioDevices() ([]player.AudioDevice, error) {
	return m.player.AudioDevices()
}

func (m *Manager) AudioDevice(name string) error {
	if err := m.player.AudioDevice(name); err != nil {
		return err
	}
	m.mu.Lock()
	m.status.AudioDevice = name
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
	activeURL := ""
	playing := false
	attemptStarted := time.Time{}
	nextAttempt := time.Time{}
	retryAttempt := 0
	reconnectAttempt := 0

	for {
		select {
		case <-ctx.Done():
			if active {
				_ = m.player.Stop()
			}
			return
		case <-m.wake:
			m.step(ctx, &active, &activeURL, &playing, &attemptStarted, &nextAttempt, &retryAttempt, &reconnectAttempt)
		case <-loop.C:
			m.step(ctx, &active, &activeURL, &playing, &attemptStarted, &nextAttempt, &retryAttempt, &reconnectAttempt)
		case <-statusTicker.C:
			if active {
				m.checkPlayerStatus(&active, &activeURL, &playing, &attemptStarted, &nextAttempt, &retryAttempt, &reconnectAttempt)
			}
		}
	}
}

func (m *Manager) step(ctx context.Context, active *bool, activeURL *string, playing *bool, attemptStarted *time.Time, nextAttempt *time.Time, retryAttempt, reconnectAttempt *int) {
	selected, listening, availability := m.snapshot()
	if ctx.Err() != nil {
		return
	}
	if !listening || selected == nil {
		if *active {
			_ = m.player.Stop()
			*active = false
			*activeURL = ""
			*playing = false
		}
		*attemptStarted = time.Time{}
		*nextAttempt = time.Time{}
		*retryAttempt = 0
		*reconnectAttempt = 0
		m.setState(StateIdle, "", nil)
		return
	}
	if availability != nil {
		available, known := availability.IsAvailable(selected.ID)
		if !known || !available {
			if *active {
				_ = m.player.Stop()
				*active = false
				*activeURL = ""
				*playing = false
			}
			*attemptStarted = time.Time{}
			*nextAttempt = time.Time{}
			*retryAttempt = 0
			*reconnectAttempt = 0
			m.setState(StateWaiting, "", nil)
			return
		}
	}
	if *active && *activeURL != selected.URL {
		_ = m.player.Stop()
		*active = false
		*activeURL = ""
		*playing = false
		*attemptStarted = time.Time{}
		*nextAttempt = time.Time{}
		*retryAttempt = 0
		*reconnectAttempt = 0
	}
	if *active {
		return
	}
	if !nextAttempt.IsZero() && time.Now().Before(*nextAttempt) {
		m.setState(StateRetrying, m.retryMessage(*nextAttempt), nil)
		return
	}
	m.setState(StateConnecting, "", nil)
	if err := m.player.Play(selected.URL); err != nil {
		delay := backoffDelay(m.retryInterval, *retryAttempt)
		(*retryAttempt)++
		*nextAttempt = time.Now().Add(delay)
		*attemptStarted = time.Time{}
		m.setState(StateRetrying, err.Error(), nil)
		m.logRetry(selected, "relay connection failed", err, delay)
		return
	}
	*active = true
	*activeURL = selected.URL
	*playing = false
	*attemptStarted = time.Now()
	*nextAttempt = time.Time{}
	*reconnectAttempt = 0
	m.setState(StateConnecting, "", nil)
}

func (m *Manager) checkPlayerStatus(active *bool, activeURL *string, playing *bool, attemptStarted *time.Time, nextAttempt *time.Time, retryAttempt, reconnectAttempt *int) {
	if !*active {
		return
	}
	status, err := m.player.Status()
	if err != nil {
		delay := backoffDelay(m.retryInterval, *retryAttempt)
		(*retryAttempt)++
		_ = m.player.Stop()
		*active = false
		*activeURL = ""
		*playing = false
		*attemptStarted = time.Time{}
		*nextAttempt = time.Now().Add(delay)
		m.setState(StateRetrying, err.Error(), nil)
		m.logRetryFromStatus("player status failed", err, delay)
		return
	}
	selected, listening, availability := m.snapshot()
	if !listening || selected == nil {
		return
	}
	if availability != nil {
		available, known := availability.IsAvailable(selected.ID)
		if !known || !available {
			_ = m.player.Stop()
			*active = false
			*activeURL = ""
			*playing = false
			*attemptStarted = time.Time{}
			*nextAttempt = time.Time{}
			*retryAttempt = 0
			*reconnectAttempt = 0
			m.setState(StateWaiting, "", status)
			return
		}
	}
	if status.State == "stopped" {
		if !attemptStarted.IsZero() && time.Since(*attemptStarted) < m.startupGracePeriod {
			m.setState(StateConnecting, "Connecting to stream...", status)
			return
		}
		delay := backoffDelay(m.retryInterval, *retryAttempt)
		(*retryAttempt)++
		_ = m.player.Stop()
		*active = false
		*activeURL = ""
		*playing = false
		*attemptStarted = time.Time{}
		*nextAttempt = time.Now().Add(delay)
		m.setState(StateRetrying, "player stopped unexpectedly", status)
		m.logRetryFromStatus("playback stopped unexpectedly", nil, delay)
		return
	}
	*playing = true
	*attemptStarted = time.Time{}
	*retryAttempt = 0
	*reconnectAttempt = 0
	m.setState(StatePlaying, "", status)
}

func backoffDelay(base time.Duration, attempt int) time.Duration {
	if base <= 0 { base = DefaultRetryInterval }
	if attempt < 0 { attempt = 0 }
	delay := base
	for i := 0; i < attempt && delay < maxRetryDelay; i++ {
		if delay > maxRetryDelay/2 { return maxRetryDelay }
		delay *= 2
	}
	if delay > maxRetryDelay { return maxRetryDelay }
	return delay
}

func (m *Manager) retryMessage(nextAttempt time.Time) string {
	remaining := time.Until(nextAttempt).Round(time.Second)
	if remaining < 0 { remaining = 0 }
	return "Retrying in " + remaining.String()
}

func (m *Manager) logRetry(selected *stream.Stream, reason string, err error, delay time.Duration) {
	if m.log == nil { return }
	args := []any{"stream_id", selected.ID, "stream_name", selected.Name, "retry_in", delay.String(), "reason", reason}
	if err != nil { args = append(args, "error", err) }
	m.log.Warn("Stream playback retry scheduled", args...)
}

func (m *Manager) logRetryFromStatus(reason string, err error, delay time.Duration) {
	if m.log == nil { return }
	args := []any{"retry_in", delay.String(), "reason", reason}
	if err != nil { args = append(args, "error", err) }
	m.log.Warn("Stream playback retry scheduled", args...)
}

func (m *Manager) logPersistenceError(action string, err error) {
	if m.log != nil { m.log.Warn("Playback persistence failed", "action", action, "error", err) }
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
	case StateIdle: status.Message = "Idle"
	case StateWaiting: status.Message = "Waiting for stream"
	case StateConnecting: status.Message = "Connecting..."
	case StatePlaying: status.Message = "Playing"
	case StateRetrying: status.Message = "Reconnecting"
	case StateError: status.Message = "Playback error"
	default: status.Message = ""
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
		status.AudioDevice = playerStatus.AudioDevice
	}
	m.status = status
}

func (m *Manager) notify() {
	select {
	case m.wake <- struct{}{}:
	default:
	}
}
