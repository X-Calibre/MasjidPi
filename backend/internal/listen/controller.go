package listen

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
)

const (
	reevaluateInterval = time.Second
	MaxSourceVolume    = 150
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
	Listening      bool         `json:"listening"`
	ActiveSource   ActiveSource `json:"active_source"`
	ActiveStreamID string       `json:"active_stream_id,omitempty"`
	MasjidID       string       `json:"masjid_id,omitempty"`
	MasjidOnline   bool         `json:"masjid_online"`
	RadioID        string       `json:"radio_id,omitempty"`
	MasjidVolume   int          `json:"masjid_volume"`
	RadioVolume    int          `json:"radio_volume"`
	Error          string       `json:"error,omitempty"`
}

type Controller struct {
	mu             sync.Mutex
	startOnce      sync.Once
	availability   Availability
	output         Output
	selectedMasjid *stream.Stream
	selectedRadio  *stream.Stream
	listening      bool
	masjidVolume   int
	radioVolume    int
	activeSource   ActiveSource
	activeStreamID string
	masjidOnline   bool
	lastError      string
	wake           chan struct{}
}

func New(availability Availability, output Output) *Controller {
	return &Controller{
		availability: availability,
		output:       output,
		masjidVolume: 100,
		radioVolume:  70,
		activeSource: ActiveNone,
		wake:         make(chan struct{}, 1),
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
	c.mu.Unlock()
	c.notify()
}

func (c *Controller) Status() Status {
	c.mu.Lock()
	defer c.mu.Unlock()
	status := Status{
		Listening:      c.listening,
		ActiveSource:   c.activeSource,
		ActiveStreamID: c.activeStreamID,
		MasjidOnline:   c.masjidOnline,
		MasjidVolume:   c.masjidVolume,
		RadioVolume:    c.radioVolume,
		Error:          c.lastError,
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
		c.deactivate()
		return
	}
	if masjidOnline && masjid != nil {
		c.activate(ActiveMasjid, *masjid, masjidVolume)
		return
	}
	if radio != nil {
		c.activate(ActiveRadio, *radio, radioVolume)
		return
	}
	c.deactivate()
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
