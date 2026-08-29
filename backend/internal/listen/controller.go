package listen

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
)

const (
	reevaluateInterval             = time.Second
	MaxSourceVolume                = 150
	MinRadioResumeDelayMinutes     = 1
	MaxRadioResumeDelayMinutes     = 30
	DefaultRadioResumeDelayMinutes = 5
)

type Availability interface {
	Status(mount string) (available bool, known bool, updatedAt time.Time)
	Events() <-chan string
}

type Output interface {
	Activate(selected stream.Stream, softwareVolume int) error
	SetSoftwareVolume(volume int) error
	Stop() error
}

type ActiveSource string
type RadioMode string

const (
	ActiveNone   ActiveSource = "none"
	ActiveMasjid ActiveSource = "masjid"
	ActiveRadio  ActiveSource = "radio"

	RadioModeSchedule RadioMode = "schedule"
	RadioModePlayNow  RadioMode = "play_now"
	RadioModeStopped  RadioMode = "stopped"
)

type Status struct {
	Listening               bool         `json:"listening"`
	MasjidEnabled           bool         `json:"masjid_enabled"`
	RadioEnabled            bool         `json:"radio_enabled"`
	ActiveSource            ActiveSource `json:"active_source"`
	ActiveStreamID          string       `json:"active_stream_id,omitempty"`
	ActiveStreamName        string       `json:"active_stream_name,omitempty"`
	MasjidID                string       `json:"masjid_id,omitempty"`
	MasjidName              string       `json:"masjid_name,omitempty"`
	MasjidOnline            bool         `json:"masjid_online"`
	RadioID                 string       `json:"radio_id,omitempty"`
	RadioName               string       `json:"radio_name,omitempty"`
	MasjidVolume            int          `json:"masjid_volume"`
	RadioVolume             int          `json:"radio_volume"`
	RadioMode               RadioMode    `json:"radio_mode"`
	RadioResumeDelayMinutes int          `json:"radio_resume_delay_minutes"`
	RadioResumePending      bool         `json:"radio_resume_pending"`
	RadioResumeAt           string       `json:"radio_resume_at,omitempty"`
	RadioScheduleEnabled    bool         `json:"radio_schedule_enabled"`
	RadioScheduleStart      string       `json:"radio_schedule_start,omitempty"`
	RadioScheduleStop       string       `json:"radio_schedule_stop,omitempty"`
	RadioScheduleAllowsNow  bool         `json:"radio_schedule_allows_now"`
	Error                   string       `json:"error,omitempty"`
}

type Controller struct {
	mu                   sync.Mutex
	startOnce            sync.Once
	availability         Availability
	output               Output
	selectedMasjid       *stream.Stream
	selectedRadio        *stream.Stream
	listening            bool
	masjidEnabled        bool
	radioEnabled         bool
	masjidVolume         int
	radioVolume          int
	radioMode            RadioMode
	playNowScheduleState bool
	radioResumeDelay     time.Duration
	radioResumeAt        time.Time
	radioScheduleEnabled bool
	radioScheduleStart   string
	radioScheduleStop    string
	activeSource         ActiveSource
	activeStreamID       string
	masjidOnline         bool
	lastError            string
	wake                 chan struct{}
}

func New(availability Availability, output Output) *Controller {
	return &Controller{
		availability:     availability,
		output:           output,
		masjidEnabled:    true,
		radioEnabled:     true,
		masjidVolume:     100,
		radioVolume:      70,
		radioMode:        RadioModeSchedule,
		radioResumeDelay: DefaultRadioResumeDelayMinutes * time.Minute,
		activeSource:     ActiveNone,
		wake:             make(chan struct{}, 1),
	}
}

func (c *Controller) Start(ctx context.Context) { c.startOnce.Do(func() { go c.run(ctx) }) }

func (c *Controller) SelectMasjid(selected *stream.Stream) error {
	if selected != nil && selected.SourceKind() != stream.KindMasjid {
		return errors.New("selected primary stream must be a masjid")
	}
	c.mu.Lock()
	c.selectedMasjid = cloneStream(selected)
	c.mu.Unlock()
	c.notify()
	return nil
}

func (c *Controller) SelectRadio(selected *stream.Stream) error {
	if selected != nil && selected.SourceKind() != stream.KindRadio {
		return errors.New("selected secondary stream must be a radio station")
	}
	c.mu.Lock()
	c.selectedRadio = cloneStream(selected)
	c.mu.Unlock()
	c.notify()
	return nil
}

func (c *Controller) SetMasjidEnabled(enabled bool) {
	c.mu.Lock()
	c.masjidEnabled = enabled
	if !enabled {
		c.radioEnabled = false
		c.listening = false
		c.radioResumeAt = time.Time{}
		c.masjidOnline = false
	}
	c.mu.Unlock()
	c.notify()
}

func (c *Controller) SetRadioEnabled(enabled bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if enabled && !c.masjidEnabled {
		return errors.New("masjid must be powered on before radio can be powered on")
	}
	c.radioEnabled = enabled
	if !enabled {
		c.radioResumeAt = time.Time{}
	}
	c.notifyLocked()
	return nil
}

func (c *Controller) SetMasjidVolume(volume int) error {
	return c.setSourceVolume(stream.KindMasjid, volume)
}
func (c *Controller) SetRadioVolume(volume int) error {
	return c.setSourceVolume(stream.KindRadio, volume)
}

func (c *Controller) SetRadioMode(mode RadioMode) error {
	if mode != RadioModeSchedule && mode != RadioModePlayNow && mode != RadioModeStopped {
		return errors.New("radio mode must be schedule, play_now or stopped")
	}
	c.mu.Lock()
	if !c.radioEnabled {
		c.mu.Unlock()
		return errors.New("radio is powered off")
	}
	if mode == RadioModePlayNow {
		if c.masjidOnline {
			c.mu.Unlock()
			return errors.New("masjid is currently online")
		}
		if c.selectedRadio == nil {
			c.mu.Unlock()
			return errors.New("no radio station is selected")
		}
		c.playNowScheduleState = !c.radioScheduleEnabled || radioWindowAllows(time.Now(), c.radioScheduleStart, c.radioScheduleStop)
	}
	c.radioMode = mode
	c.radioResumeAt = time.Time{}
	c.mu.Unlock()
	c.notify()
	return nil
}

func (c *Controller) SetRadioResumeDelayMinutes(minutes int) error {
	if minutes < MinRadioResumeDelayMinutes || minutes > MaxRadioResumeDelayMinutes {
		return errors.New("radio resume delay must be between 1 and 30 minutes")
	}
	c.mu.Lock()
	c.radioResumeDelay = time.Duration(minutes) * time.Minute
	if !c.radioResumeAt.IsZero() {
		c.radioResumeAt = time.Now().Add(c.radioResumeDelay)
	}
	c.mu.Unlock()
	c.notify()
	return nil
}

func (c *Controller) SetRadioSchedule(enabled bool, start, stop string) error {
	if enabled {
		if _, err := parseClock(start); err != nil {
			return fmt.Errorf("invalid radio schedule start: %w", err)
		}
		if _, err := parseClock(stop); err != nil {
			return fmt.Errorf("invalid radio schedule stop: %w", err)
		}
		if start == stop {
			return errors.New("radio schedule start and stop times must differ")
		}
	}
	c.mu.Lock()
	c.radioScheduleEnabled = enabled
	c.radioScheduleStart = start
	c.radioScheduleStop = stop
	c.mu.Unlock()
	c.notify()
	return nil
}

func (c *Controller) setSourceVolume(kind stream.Kind, volume int) error {
	if volume < 0 || volume > MaxSourceVolume {
		return errors.New("source volume must be between 0 and 150")
	}
	c.mu.Lock()
	if kind == stream.KindMasjid {
		c.masjidVolume = volume
	} else {
		c.radioVolume = volume
	}
	active := (kind == stream.KindMasjid && c.activeSource == ActiveMasjid) ||
		(kind == stream.KindRadio && c.activeSource == ActiveRadio)
	output := c.output
	c.mu.Unlock()
	if active && output != nil {
		if err := output.SetSoftwareVolume(volume); err != nil {
			c.mu.Lock()
			c.lastError = err.Error()
			c.mu.Unlock()
			return err
		}
	}
	return nil
}

func (c *Controller) Listen() {
	c.mu.Lock()
	if c.masjidEnabled {
		c.listening = true
	}
	c.mu.Unlock()
	c.notify()
}

func (c *Controller) Stop() {
	c.mu.Lock()
	c.listening = false
	c.radioResumeAt = time.Time{}
	c.mu.Unlock()
	c.notify()
}

func (c *Controller) Status() Status {
	c.mu.Lock()
	defer c.mu.Unlock()
	allowsNow := !c.radioScheduleEnabled || radioWindowAllows(time.Now(), c.radioScheduleStart, c.radioScheduleStop)
	status := Status{
		Listening:               c.listening,
		MasjidEnabled:           c.masjidEnabled,
		RadioEnabled:            c.radioEnabled,
		ActiveSource:            c.activeSource,
		ActiveStreamID:          c.activeStreamID,
		MasjidOnline:            c.masjidOnline,
		MasjidVolume:            c.masjidVolume,
		RadioVolume:             c.radioVolume,
		RadioMode:               c.radioMode,
		RadioResumeDelayMinutes: int(c.radioResumeDelay / time.Minute),
		RadioResumePending:      !c.radioResumeAt.IsZero(),
		RadioScheduleEnabled:    c.radioScheduleEnabled,
		RadioScheduleStart:      c.radioScheduleStart,
		RadioScheduleStop:       c.radioScheduleStop,
		RadioScheduleAllowsNow:  allowsNow,
		Error:                   c.lastError,
	}
	if !c.radioResumeAt.IsZero() {
		status.RadioResumeAt = c.radioResumeAt.UTC().Format(time.RFC3339)
	}
	if c.selectedMasjid != nil {
		status.MasjidID = c.selectedMasjid.ID
		status.MasjidName = c.selectedMasjid.Name
	}
	if c.selectedRadio != nil {
		status.RadioID = c.selectedRadio.ID
		status.RadioName = c.selectedRadio.Name
	}
	switch c.activeSource {
	case ActiveMasjid:
		status.ActiveStreamName = status.MasjidName
	case ActiveRadio:
		status.ActiveStreamName = status.RadioName
	}
	return status
}

func (c *Controller) run(ctx context.Context) {
	ticker := time.NewTicker(reevaluateInterval)
	defer ticker.Stop()
	var events <-chan string
	if c.availability != nil {
		events = c.availability.Events()
	}
	for {
		select {
		case <-ctx.Done():
			c.deactivate()
			return
		case <-c.wake:
			c.step()
		case <-events:
			c.step()
		case <-ticker.C:
			c.step()
		}
	}
}

func (c *Controller) step() {
	c.mu.Lock()
	listening := c.listening
	masjidEnabled := c.masjidEnabled
	radioEnabled := c.radioEnabled
	masjid := cloneStream(c.selectedMasjid)
	radio := cloneStream(c.selectedRadio)
	masjidVolume := c.masjidVolume
	radioVolume := c.radioVolume
	availability := c.availability
	activeSource := c.activeSource
	resumeAt := c.radioResumeAt
	resumeDelay := c.radioResumeDelay
	scheduleEnabled := c.radioScheduleEnabled
	scheduleStart := c.radioScheduleStart
	scheduleStop := c.radioScheduleStop
	mode := c.radioMode
	playNowScheduleState := c.playNowScheduleState
	c.mu.Unlock()

	masjidOnline := false
	if masjidEnabled && masjid != nil && availability != nil {
		available, known, _ := availability.Status(masjid.ID)
		masjidOnline = known && available
	}
	c.mu.Lock()
	c.masjidOnline = masjidOnline
	c.mu.Unlock()

	if !listening {
		c.clearRadioResumeDelay()
		c.deactivate()
		return
	}

	if masjidEnabled && masjidOnline && masjid != nil {
		c.clearRadioResumeDelay()
		if mode == RadioModePlayNow {
			c.mu.Lock()
			c.radioMode = RadioModeSchedule
			c.mu.Unlock()
		}
		c.activate(ActiveMasjid, *masjid, masjidVolume)
		return
	}

	if !radioEnabled {
		c.clearRadioResumeDelay()
		c.deactivate()
		return
	}

	allowedNow := !scheduleEnabled || radioWindowAllows(time.Now(), scheduleStart, scheduleStop)
	if mode == RadioModePlayNow && scheduleEnabled && allowedNow != playNowScheduleState {
		c.mu.Lock()
		c.radioMode = RadioModeSchedule
		c.mu.Unlock()
		mode = RadioModeSchedule
	}
	if mode == RadioModeStopped {
		c.clearRadioResumeDelay()
		c.deactivate()
		return
	}
	if mode == RadioModePlayNow {
		c.clearRadioResumeDelay()
		if radio != nil {
			c.activate(ActiveRadio, *radio, radioVolume)
		} else {
			c.deactivate()
		}
		return
	}
	if !allowedNow {
		c.clearRadioResumeDelay()
		c.deactivate()
		return
	}
	if activeSource == ActiveMasjid && resumeAt.IsZero() {
		c.mu.Lock()
		c.radioResumeAt = time.Now().Add(resumeDelay)
		c.mu.Unlock()
		c.deactivate()
		return
	}
	c.mu.Lock()
	resumeAt = c.radioResumeAt
	c.mu.Unlock()
	if !resumeAt.IsZero() {
		if time.Now().Before(resumeAt) {
			c.deactivate()
			return
		}
		c.clearRadioResumeDelay()
	}
	if radio != nil {
		c.activate(ActiveRadio, *radio, radioVolume)
		return
	}
	c.deactivate()
}

func parseClock(value string) (int, error) {
	parsed, err := time.Parse("15:04", value)
	if err != nil {
		return 0, errors.New("time must use HH:MM format")
	}
	return parsed.Hour()*60 + parsed.Minute(), nil
}

func radioWindowAllows(now time.Time, start, stop string) bool {
	startMinutes, startErr := parseClock(start)
	stopMinutes, stopErr := parseClock(stop)
	if startErr != nil || stopErr != nil || startMinutes == stopMinutes {
		return false
	}
	nowMinutes := now.Hour()*60 + now.Minute()
	if startMinutes < stopMinutes {
		return nowMinutes >= startMinutes && nowMinutes < stopMinutes
	}
	return nowMinutes >= startMinutes || nowMinutes < stopMinutes
}

func (c *Controller) clearRadioResumeDelay() {
	c.mu.Lock()
	c.radioResumeAt = time.Time{}
	c.mu.Unlock()
}

func (c *Controller) activate(source ActiveSource, selected stream.Stream, volume int) {
	c.mu.Lock()
	if c.activeSource == source && c.activeStreamID == selected.ID {
		c.lastError = ""
		c.mu.Unlock()
		return
	}
	output := c.output
	c.mu.Unlock()
	if output == nil {
		return
	}
	if err := output.Activate(selected, volume); err != nil {
		c.mu.Lock()
		c.lastError = err.Error()
		c.mu.Unlock()
		return
	}
	c.mu.Lock()
	c.activeSource = source
	c.activeStreamID = selected.ID
	c.lastError = ""
	c.mu.Unlock()
}

func (c *Controller) deactivate() {
	c.mu.Lock()
	if c.activeSource == ActiveNone {
		c.mu.Unlock()
		return
	}
	output := c.output
	c.mu.Unlock()
	if output != nil {
		if err := output.Stop(); err != nil {
			c.mu.Lock()
			c.lastError = err.Error()
			c.mu.Unlock()
			return
		}
	}
	c.mu.Lock()
	c.activeSource = ActiveNone
	c.activeStreamID = ""
	c.lastError = ""
	c.mu.Unlock()
}

func (c *Controller) notify() { c.mu.Lock(); c.notifyLocked(); c.mu.Unlock() }

func (c *Controller) notifyLocked() {
	select {
	case c.wake <- struct{}{}:
	default:
	}
}

func cloneStream(selected *stream.Stream) *stream.Stream {
	if selected == nil {
		return nil
	}
	copy := *selected
	copy.FallbackURLs = append([]string(nil), selected.FallbackURLs...)
	return &copy
}
