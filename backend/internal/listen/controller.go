package listen

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
)

const (
	reevaluateInterval          = time.Second
	MaxSourceVolume             = 150
	MinRadioResumeDelayMinutes  = 1
	MaxRadioResumeDelayMinutes  = 30
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

const (
	ActiveNone   ActiveSource = "none"
	ActiveMasjid ActiveSource = "masjid"
	ActiveRadio  ActiveSource = "radio"
)

type Status struct {
	Listening               bool         `json:"listening"`
	ActiveSource            ActiveSource `json:"active_source"`
	ActiveStreamID          string       `json:"active_stream_id,omitempty"`
	MasjidID                string       `json:"masjid_id,omitempty"`
	MasjidOnline            bool         `json:"masjid_online"`
	RadioID                 string       `json:"radio_id,omitempty"`
	MasjidVolume            int          `json:"masjid_volume"`
	RadioVolume             int          `json:"radio_volume"`
	RadioResumeDelayMinutes int          `json:"radio_resume_delay_minutes"`
	RadioResumePending      bool         `json:"radio_resume_pending"`
	RadioResumeAt           string       `json:"radio_resume_at,omitempty"`
	Error                   string       `json:"error,omitempty"`
}

type Controller struct {
	mu                 sync.Mutex
	startOnce          sync.Once
	availability       Availability
	output             Output
	selectedMasjid     *stream.Stream
	selectedRadio      *stream.Stream
	listening          bool
	masjidVolume       int
	radioVolume        int
	radioResumeDelay   time.Duration
	radioResumeAt      time.Time
	activeSource       ActiveSource
	activeStreamID     string
	masjidOnline       bool
	lastError          string
	wake               chan struct{}
}

func New(availability Availability, output Output) *Controller {
	return &Controller{
		availability:     availability,
		output:           output,
		masjidVolume:     100,
		radioVolume:      70,
		radioResumeDelay: DefaultRadioResumeDelayMinutes * time.Minute,
		activeSource:     ActiveNone,
		wake:             make(chan struct{}, 1),
	}
}

func (c *Controller) Start(ctx context.Context) {
	c.startOnce.Do(func() { go c.run(ctx) })
}

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

func (c *Controller) SetMasjidVolume(volume int) error {
	return c.setSourceVolume(stream.KindMasjid, volume)
}

func (c *Controller) SetRadioVolume(volume int) error {
	return c.setSourceVolume(stream.KindRadio, volume)
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
	c.listening = true
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
	status := Status{
		Listening:               c.listening,
		ActiveSource:            c.activeSource,
		ActiveStreamID:          c.activeStreamID,
		MasjidOnline:            c.masjidOnline,
		MasjidVolume:            c.masjidVolume,
		RadioVolume:             c.radioVolume,
		RadioResumeDelayMinutes: int(c.radioResumeDelay / time.Minute),
		RadioResumePending:      !c.radioResumeAt.IsZero(),
		Error:                   c.lastError,
	}
	if !c.radioResumeAt.IsZero() {
		status.RadioResumeAt = c.radioResumeAt.UTC().Format(time.RFC3339)
	}
	if c.selectedMasjid != nil {
		status.MasjidID = c.selectedMasjid.ID
	}
	if c.selectedRadio != nil {
		status.RadioID = c.selectedRadio.ID
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
	masjid := cloneStream(c.selectedMasjid)
	radio := cloneStream(c.selectedRadio)
	masjidVolume := c.masjidVolume
	radioVolume := c.radioVolume
	availability := c.availability
	activeSource := c.activeSource
	resumeAt := c.radioResumeAt
	resumeDelay := c.radioResumeDelay
	c.mu.Unlock()

	masjidOnline := false
	if masjid != nil && availability != nil {
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

	// Masjid priority is immediate. If it returns during the post-broadcast
	// radio hold, cancel the hold and restore the masjid without delay.
	if masjidOnline && masjid != nil {
		c.clearRadioResumeDelay()
		c.activate(ActiveMasjid, *masjid, masjidVolume)
		return
	}

	// A delay is started only when a masjid that was actively playing goes
	// offline. Starting MasjidPi with an already-offline masjid does not delay
	// normal radio playback.
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

func (c *Controller) notify() {
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
