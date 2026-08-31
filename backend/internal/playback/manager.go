package playback

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"sync"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/player"
	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
	"github.com/X-Calibre/MasjidPi/backend/internal/version"
)

const (
	DefaultRetryInterval      = 5 * time.Second
	DefaultStartupGracePeriod = 10 * time.Second
	DefaultMountStartupDelay  = 2 * time.Second
	statusCheckInterval       = time.Second
	maxRetryDelay             = 5 * time.Minute
	preferredEndpointTTL      = 6 * time.Hour
)

type State string

const (
	StateIdle       State = "idle"
	StateWaiting    State = "waiting"
	StateConnecting State = "connecting"
	StatePlaying    State = "playing"
	StateRetrying   State = "retrying"
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
	Status(mount string) (available bool, known bool, updatedAt time.Time)
	Events() <-chan string
}

type Persistence interface {
	Save(streamID string) error
	Clear() error
}
type VolumePersistence interface {
	Load(device string) (int, bool, error)
	Save(device string, volume int) error
}

type Config struct {
	RetryInterval       time.Duration
	StartupGracePeriod  time.Duration
	MountStartupDelay   time.Duration
	StatusCheckInterval time.Duration
	Logger              *slog.Logger
}

type Manager struct {
	mu                 sync.Mutex
	startOnce          sync.Once
	player             Player
	audioDevices       player.AudioDeviceProvider
	availability       Availability
	persistence        Persistence
	volumeStore        VolumePersistence
	log                *slog.Logger
	retryInterval      time.Duration
	startupGracePeriod time.Duration
	mountStartupDelay  time.Duration
	statusInterval     time.Duration
	wake               chan struct{}
	selected           *stream.Stream
	currentURL         string
	listening          bool
	state              State
	lastError          string
	status             Status
	volume             int
	volumeSet          bool
	volumeDevice       string
	volumeSupported    bool
	preferredEndpoints map[string]preferredEndpoint
}

type preferredEndpoint struct {
	url       string
	expiresAt time.Time
}

type Status struct {
	Version         string `json:"version"`
	State           string `json:"state"`
	Message         string `json:"message"`
	URL             string `json:"url"`
	Volume          int    `json:"volume"`
	VolumeSupported bool   `json:"volume_supported"`
	Paused          bool   `json:"paused"`
	AudioDevice     string `json:"audio_device,omitempty"`
	Listening       bool   `json:"listening"`
	StreamID        string `json:"stream_id,omitempty"`
	StreamName      string `json:"stream_name,omitempty"`
	Endpoint        string `json:"endpoint,omitempty"`
	FallbackUsed    bool   `json:"fallback_used,omitempty"`
	Error           string `json:"error,omitempty"`
}

func New(player Player, cfg Config) *Manager {
	if cfg.RetryInterval == 0 {
		cfg.RetryInterval = DefaultRetryInterval
	}
	if cfg.StartupGracePeriod == 0 {
		cfg.StartupGracePeriod = DefaultStartupGracePeriod
	}
	if cfg.MountStartupDelay == 0 {
		cfg.MountStartupDelay = DefaultMountStartupDelay
	}
	if cfg.StatusCheckInterval == 0 {
		cfg.StatusCheckInterval = statusCheckInterval
	}
	return &Manager{
		player:             player,
		audioDevices:       player,
		log:                cfg.Logger,
		retryInterval:      cfg.RetryInterval,
		startupGracePeriod: cfg.StartupGracePeriod,
		mountStartupDelay:  cfg.MountStartupDelay,
		statusInterval:     cfg.StatusCheckInterval,
		wake:               make(chan struct{}, 1),
		state:              StateIdle,
		status:             Status{Version: version.Version, State: string(StateIdle), Volume: 100},
		preferredEndpoints: make(map[string]preferredEndpoint),
	}
}

func (m *Manager) SetAudioDeviceProvider(provider player.AudioDeviceProvider) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if provider == nil {
		m.audioDevices = m.player
		return
	}
	m.audioDevices = provider
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
func (m *Manager) SetVolumePersistence(persistence VolumePersistence) {
	m.mu.Lock()
	m.volumeStore = persistence
	m.mu.Unlock()
}
func (m *Manager) Start(ctx context.Context) { m.startOnce.Do(func() { go m.run(ctx) }) }

func (m *Manager) InitializeVolume() error {
	status, err := m.player.Status()
	if err != nil {
		return err
	}
	volume := status.Volume
	if !status.VolumeSupported {
		volume = 100
	}

	m.mu.Lock()
	m.volume = volume
	m.volumeSet = true
	m.volumeDevice = status.AudioDevice
	m.volumeSupported = status.VolumeSupported
	volumeStore := m.volumeStore
	m.mu.Unlock()

	if volumeStore == nil || status.AudioDevice == "" || !status.VolumeSupported {
		return nil
	}
	if saved, ok, err := volumeStore.Load(status.AudioDevice); err != nil {
		m.logPersistenceError("loading volume", err)
		return nil
	} else if ok {
		if err := m.player.Volume(saved); err != nil {
			return err
		}
		m.mu.Lock()
		m.volume = saved
		m.status.Volume = saved
		m.mu.Unlock()
	}
	return nil
}

func (m *Manager) Play(selected stream.Stream) {
	m.mu.Lock()
	m.selected = &selected
	m.currentURL = selected.URL
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
	m.currentURL = ""
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
	if volume < 0 || volume > 100 {
		return errors.New("volume must be between 0 and 100")
	}
	status, err := m.player.Status()
	if err != nil {
		return err
	}
	if err := m.player.Volume(volume); err != nil {
		return err
	}

	m.mu.Lock()
	m.volume = volume
	m.volumeSet = true
	m.volumeDevice = status.AudioDevice
	m.volumeSupported = true
	m.status.Volume = volume
	m.status.VolumeSupported = true
	m.status.AudioDevice = status.AudioDevice
	volumeStore := m.volumeStore
	m.mu.Unlock()
	if volumeStore != nil && status.AudioDevice != "" {
		if err := volumeStore.Save(status.AudioDevice, volume); err != nil {
			m.logPersistenceError("saving volume", err)
			return err
		}
	}
	return nil
}

func (m *Manager) AudioDevices() ([]player.AudioDevice, error) {
	m.mu.Lock()
	provider := m.audioDevices
	m.mu.Unlock()
	return provider.AudioDevices()
}

func (m *Manager) AudioDevice(name string) error {
	if err := m.player.AudioDevice(name); err != nil {
		return err
	}
	status, err := m.player.Status()
	if err != nil {
		return err
	}
	return m.applyAudioDeviceVolume(status)
}

func (m *Manager) applyAudioDeviceVolume(status *player.Status) error {
	volume := status.Volume
	volumeSupported := status.VolumeSupported
	if !volumeSupported {
		volume = 100
	}

	m.mu.Lock()
	volumeStore := m.volumeStore
	m.mu.Unlock()
	if volumeStore != nil && status.AudioDevice != "" && volumeSupported {
		if saved, ok, err := volumeStore.Load(status.AudioDevice); err != nil {
			m.logPersistenceError("loading volume", err)
		} else if ok {
			if err := m.player.Volume(saved); err != nil {
				return err
			}
			volume = saved
		}
	}

	m.mu.Lock()
	m.volume = volume
	m.volumeSet = true
	m.volumeDevice = status.AudioDevice
	m.volumeSupported = volumeSupported
	m.status.Volume = volume
	m.status.VolumeSupported = volumeSupported
	m.status.AudioDevice = status.AudioDevice
	m.mu.Unlock()
	return nil
}

func (m *Manager) Status() Status { m.mu.Lock(); defer m.mu.Unlock(); return m.status }

type runtimeState struct {
	active         bool
	activeURL      string
	selectedID     string
	candidates     []string
	candidateIndex int
	attemptStarted time.Time
	nextAttempt    time.Time
	retryAttempt   int
}

func (state *runtimeState) reset() { *state = runtimeState{} }

func (m *Manager) stopActive(state *runtimeState) {
	if state.active {
		_ = m.player.Stop()
	}
	state.reset()
}

func (m *Manager) scheduleRetry(state *runtimeState) time.Duration {
	delay := backoffDelay(m.retryInterval, state.retryAttempt)
	state.retryAttempt++
	if state.active {
		_ = m.player.Stop()
	}
	state.active = false
	state.activeURL = ""
	state.attemptStarted = time.Time{}
	state.candidateIndex = 0
	state.nextAttempt = time.Now().Add(delay)
	return delay
}

func (m *Manager) prepareCandidates(state *runtimeState, selected *stream.Stream) {
	if state.selectedID == selected.ID && len(state.candidates) > 0 {
		return
	}
	state.reset()
	state.selectedID = selected.ID
	state.candidates = selected.PlaybackURLs()

	m.mu.Lock()
	preferred, ok := m.preferredEndpoints[selected.ID]
	if ok && time.Now().After(preferred.expiresAt) {
		delete(m.preferredEndpoints, selected.ID)
		ok = false
	}
	m.mu.Unlock()
	if !ok || len(state.candidates) < 2 || state.candidates[0] == preferred.url {
		return
	}
	for index, candidate := range state.candidates {
		if candidate == preferred.url {
			state.candidates[0], state.candidates[index] = state.candidates[index], state.candidates[0]
			return
		}
	}
}

func (m *Manager) tryNextCandidate(state *runtimeState) bool {
	if state.active {
		_ = m.player.Stop()
	}
	state.active = false
	state.activeURL = ""
	state.attemptStarted = time.Time{}
	if state.candidateIndex+1 >= len(state.candidates) {
		return false
	}
	state.candidateIndex++
	state.nextAttempt = time.Time{}
	m.notify()
	return true
}

func (m *Manager) run(ctx context.Context) {
	statusTicker := time.NewTicker(m.statusInterval)
	defer statusTicker.Stop()
	stepTimer := time.NewTimer(time.Hour)
	if !stepTimer.Stop() {
		<-stepTimer.C
	}
	defer stepTimer.Stop()
	var stepDeadline <-chan time.Time
	m.mu.Lock()
	availability := m.availability
	m.mu.Unlock()
	var availabilityEvents <-chan string
	if availability != nil {
		availabilityEvents = availability.Events()
	}
	state := runtimeState{}
	resetStepTimer := func() {
		if !stepTimer.Stop() {
			select {
			case <-stepTimer.C:
			default:
			}
		}
		delay, ok := m.nextStepDelay(&state)
		if !ok {
			stepDeadline = nil
			return
		}
		stepTimer.Reset(delay)
		stepDeadline = stepTimer.C
	}
	for {
		select {
		case <-ctx.Done():
			m.stopActive(&state)
			return
		case <-m.wake:
			m.step(ctx, &state)
			resetStepTimer()
		case <-availabilityEvents:
			m.step(ctx, &state)
			resetStepTimer()
		case <-stepDeadline:
			m.step(ctx, &state)
			resetStepTimer()
		case <-statusTicker.C:
			if state.active {
				m.checkPlayerStatus(&state)
			}
			resetStepTimer()
		}
	}
}

func (m *Manager) nextStepDelay(state *runtimeState) (time.Duration, bool) {
	now := time.Now()
	deadline := state.nextAttempt
	selected, listening, availability := m.snapshot()
	if listening && selected != nil && availability != nil && !state.active {
		if available, known, updatedAt := availability.Status(selected.ID); available && known && !updatedAt.IsZero() {
			startupDeadline := updatedAt.Add(m.mountStartupDelay)
			if startupDeadline.After(now) && (deadline.IsZero() || startupDeadline.Before(deadline)) {
				deadline = startupDeadline
			}
		}
	}
	if deadline.IsZero() {
		return 0, false
	}
	delay := deadline.Sub(now)
	if delay < time.Millisecond {
		delay = time.Millisecond
	}
	return delay, true
}

func (m *Manager) step(ctx context.Context, state *runtimeState) {
	selected, listening, availability := m.snapshot()
	if ctx.Err() != nil {
		return
	}
	if !listening || selected == nil {
		m.stopActive(state)
		m.setState(StateIdle, "", nil)
		return
	}
	if availability != nil {
		available, known, updatedAt := availability.Status(selected.ID)
		if !known || !available {
			m.stopActive(state)
			m.setState(StateWaiting, "", nil)
			return
		}
		if !state.active && !updatedAt.IsZero() && time.Now().Before(updatedAt.Add(m.mountStartupDelay)) {
			m.stopActive(state)
			m.setState(StateWaiting, "", nil)
			return
		}
	}
	if state.selectedID != "" && state.selectedID != selected.ID {
		m.stopActive(state)
	}
	m.prepareCandidates(state, selected)
	if len(state.candidates) == 0 {
		delay := m.scheduleRetry(state)
		m.setState(StateRetrying, "stream has no playback URL", nil)
		m.logRetry(selected, "stream has no playback URL", nil, delay)
		return
	}
	candidateURL := state.candidates[state.candidateIndex]
	if state.active && state.activeURL != candidateURL {
		m.stopActive(state)
		m.prepareCandidates(state, selected)
		candidateURL = state.candidates[state.candidateIndex]
	}
	if state.active {
		return
	}
	if !state.nextAttempt.IsZero() && time.Now().Before(state.nextAttempt) {
		m.setState(StateRetrying, m.retryMessage(state.nextAttempt), nil)
		return
	}
	m.setCurrentEndpoint(candidateURL)
	m.setState(StateConnecting, "", nil)
	if err := m.player.Play(candidateURL); err != nil {
		if m.tryNextCandidate(state) {
			m.setState(StateConnecting, "", nil)
			m.logEndpointFailure(selected, candidateURL, err)
			return
		}
		delay := m.scheduleRetry(state)
		m.setState(StateRetrying, err.Error(), nil)
		m.logRetry(selected, "all stream endpoints failed", err, delay)
		return
	}
	state.active = true
	state.activeURL = candidateURL
	state.attemptStarted = time.Now()
	state.nextAttempt = time.Time{}
	m.setState(StateConnecting, "", nil)
}

func (m *Manager) checkPlayerStatus(state *runtimeState) {
	if !state.active {
		return
	}
	status, err := m.player.Status()
	if err != nil {
		delay := m.scheduleRetry(state)
		m.setState(StateRetrying, err.Error(), nil)
		m.logRetryFromStatus("player status failed", err, delay)
		return
	}
	selected, listening, availability := m.snapshot()
	if !listening || selected == nil {
		return
	}
	if availability != nil {
		available, known, _ := availability.Status(selected.ID)
		if !known || !available {
			m.stopActive(state)
			m.setState(StateWaiting, "", status)
			return
		}
	}
	if status.AudioDevice != "" {
		m.mu.Lock()
		deviceChanged := m.volumeDevice != status.AudioDevice
		m.mu.Unlock()
		if deviceChanged {
			if err := m.applyAudioDeviceVolume(status); err != nil {
				delay := m.scheduleRetry(state)
				m.setState(StateRetrying, err.Error(), status)
				m.logRetryFromStatus("restoring volume for audio device", err, delay)
				return
			}
		}
	}
	if status.State == "stopped" {
		if !state.attemptStarted.IsZero() && time.Since(state.attemptStarted) < m.startupGracePeriod {
			m.setState(StateConnecting, "Connecting to stream...", status)
			return
		}
		failedURL := state.activeURL
		if m.tryNextCandidate(state) {
			m.setState(StateConnecting, "", status)
			m.logEndpointFailure(selected, failedURL, nil)
			return
		}
		delay := m.scheduleRetry(state)
		m.setState(StateRetrying, "player stopped unexpectedly", status)
		m.logRetryFromStatus("all stream endpoints stopped unexpectedly", nil, delay)
		return
	}
	newlyPlaying := !state.attemptStarted.IsZero()
	state.attemptStarted = time.Time{}
	state.retryAttempt = 0
	m.rememberEndpoint(selected.ID, state.activeURL)
	if newlyPlaying {
		m.logEndpointSuccess(selected, state.activeURL)
	}
	if err := m.restoreVolume(&status.Volume); err != nil {
		delay := m.scheduleRetry(state)
		m.setState(StateRetrying, err.Error(), status)
		m.logRetryFromStatus("restoring volume after player recovery", err, delay)
		return
	}
	m.setState(StatePlaying, "", status)
}

func (m *Manager) restoreVolume(playerVolume *int) error {
	m.mu.Lock()
	volume := m.volume
	volumeSet := m.volumeSet
	supported := m.volumeSupported
	m.mu.Unlock()
	if !volumeSet || !supported || *playerVolume == volume {
		return nil
	}
	if err := m.player.Volume(volume); err != nil {
		return err
	}
	*playerVolume = volume
	return nil
}

func backoffDelay(base time.Duration, attempt int) time.Duration {
	if base <= 0 {
		base = DefaultRetryInterval
	}
	if attempt < 0 {
		attempt = 0
	}
	delay := base
	for i := 0; i < attempt && delay < maxRetryDelay; i++ {
		if delay > maxRetryDelay/2 {
			return maxRetryDelay
		}
		delay *= 2
	}
	if delay > maxRetryDelay {
		return maxRetryDelay
	}
	return delay
}
func (m *Manager) retryMessage(nextAttempt time.Time) string {
	remaining := time.Until(nextAttempt).Round(time.Second)
	if remaining < 0 {
		remaining = 0
	}
	return "Retrying in " + remaining.String()
}
func (m *Manager) logRetry(selected *stream.Stream, reason string, err error, delay time.Duration) {
	if m.log == nil {
		return
	}
	args := []any{"stream_id", selected.ID, "stream_name", selected.Name, "retry_in", delay.String(), "reason", reason}
	if err != nil {
		args = append(args, "error", err)
	}
	m.log.Warn("Stream playback retry scheduled", args...)
}
func (m *Manager) logRetryFromStatus(reason string, err error, delay time.Duration) {
	if m.log == nil {
		return
	}
	args := []any{"retry_in", delay.String(), "reason", reason}
	if err != nil {
		args = append(args, "error", err)
	}
	m.log.Warn("Stream playback retry scheduled", args...)
}
func (m *Manager) logEndpointFailure(selected *stream.Stream, endpointURL string, err error) {
	if m.log == nil {
		return
	}
	args := []any{"stream_id", selected.ID, "stream_name", selected.Name, "endpoint", endpointName(endpointURL)}
	if err != nil {
		args = append(args, "error", err)
	}
	m.log.Warn("Stream endpoint failed; trying fallback", args...)
}
func (m *Manager) logEndpointSuccess(selected *stream.Stream, endpointURL string) {
	if m.log == nil {
		return
	}
	m.log.Info(
		"Stream playback established",
		"stream_id", selected.ID,
		"stream_name", selected.Name,
		"endpoint", endpointName(endpointURL),
		"fallback", endpointURL != selected.URL,
	)
}
func (m *Manager) rememberEndpoint(streamID, endpointURL string) {
	if streamID == "" || endpointURL == "" {
		return
	}
	m.mu.Lock()
	m.preferredEndpoints[streamID] = preferredEndpoint{url: endpointURL, expiresAt: time.Now().Add(preferredEndpointTTL)}
	m.mu.Unlock()
}
func (m *Manager) setCurrentEndpoint(endpointURL string) {
	m.mu.Lock()
	m.currentURL = endpointURL
	m.mu.Unlock()
}
func endpointName(endpointURL string) string {
	parsed, err := url.Parse(endpointURL)
	if err != nil {
		return "primary"
	}
	if parsed.Hostname() == "icecast.livemasjid.com" {
		return "icecast"
	}
	if parsed.Hostname() == "relay.livemasjid.com" {
		return "relay"
	}
	return "primary"
}
func (m *Manager) logPersistenceError(action string, err error) {
	if m.log != nil {
		m.log.Warn("Playback persistence failed", "action", action, "error", err)
	}
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
	default:
		status.Message = ""
	}
	status.Listening = m.listening
	status.Error = m.lastError
	if m.selected == nil {
		status.StreamID = ""
		status.StreamName = ""
		status.URL = ""
		status.Endpoint = ""
		status.FallbackUsed = false
	} else {
		status.StreamID = m.selected.ID
		status.StreamName = m.selected.Name
		status.URL = m.currentURL
		if m.currentURL == "" {
			status.Endpoint = ""
			status.FallbackUsed = false
		} else {
			status.Endpoint = endpointName(m.currentURL)
			status.FallbackUsed = m.currentURL != m.selected.URL
		}
	}
	if playerStatus != nil {
		status.Volume = playerStatus.Volume
		status.VolumeSupported = playerStatus.VolumeSupported
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
