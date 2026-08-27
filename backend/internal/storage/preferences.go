package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/X-Calibre/MasjidPi/backend/internal/atomicfile"
)

const (
	DefaultMasjidVolume       = 100
	DefaultRadioVolume        = 70
	DefaultRadioResumeDelay   = 5
	MinRadioResumeDelay       = 1
	MaxRadioResumeDelay       = 30
	DefaultRadioScheduleStart = "06:00"
	DefaultRadioScheduleStop  = "22:00"
	DefaultRadioMode          = "schedule"
)

type PreferencesState struct {
	LastStreamID string `json:"last_stream_id,omitempty"`
	Autoplay     bool   `json:"autoplay"`

	SelectedMasjidID         string `json:"selected_masjid_id,omitempty"`
	SelectedRadioID          string `json:"selected_radio_id,omitempty"`
	ResumeListening          bool   `json:"resume_listening"`
	MasjidEnabled            *bool  `json:"masjid_enabled,omitempty"`
	RadioEnabled             *bool  `json:"radio_enabled,omitempty"`
	MasjidVolume             int    `json:"masjid_volume"`
	RadioVolume              int    `json:"radio_volume"`
	SourceVolumesSet         bool   `json:"source_volumes_set,omitempty"`
	RadioResumeDelayMinutes  int    `json:"radio_resume_delay_minutes,omitempty"`
	RadioScheduleEnabled     bool   `json:"radio_schedule_enabled,omitempty"`
	RadioScheduleStart       string `json:"radio_schedule_start,omitempty"`
	RadioScheduleStop        string `json:"radio_schedule_stop,omitempty"`
	RadioMode                string `json:"radio_mode,omitempty"`
}

func boolPointer(value bool) *bool { return &value }

func (s PreferencesState) Normalized() PreferencesState {
	if s.SelectedMasjidID == "" && s.LastStreamID != "" {
		s.SelectedMasjidID = s.LastStreamID
	}
	if !s.ResumeListening && s.Autoplay {
		s.ResumeListening = true
	}
	if s.MasjidEnabled == nil {
		s.MasjidEnabled = boolPointer(true)
	}
	if s.RadioEnabled == nil {
		s.RadioEnabled = boolPointer(true)
	}
	if !*s.MasjidEnabled {
		s.RadioEnabled = boolPointer(false)
	}
	if !s.SourceVolumesSet {
		s.MasjidVolume = DefaultMasjidVolume
		s.RadioVolume = DefaultRadioVolume
	}
	if s.RadioResumeDelayMinutes < MinRadioResumeDelay || s.RadioResumeDelayMinutes > MaxRadioResumeDelay {
		s.RadioResumeDelayMinutes = DefaultRadioResumeDelay
	}
	if s.RadioScheduleStart == "" {
		s.RadioScheduleStart = DefaultRadioScheduleStart
	}
	if s.RadioScheduleStop == "" {
		s.RadioScheduleStop = DefaultRadioScheduleStop
	}
	if s.RadioMode != "schedule" && s.RadioMode != "stopped" {
		s.RadioMode = DefaultRadioMode
	}
	return s
}

type Preferences struct {
	path   string
	mu     sync.Mutex
	loaded bool
	exists bool
	state  PreferencesState
}

func NewPreferences(path string) *Preferences { return &Preferences{path: path} }

func (p *Preferences) Load() (PreferencesState, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if err := p.loadLocked(); err != nil {
		return PreferencesState{}, err
	}
	return p.state.Normalized(), nil
}

func (p *Preferences) Save(state PreferencesState) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if err := p.loadLocked(); err != nil {
		return err
	}
	state = state.Normalized()
	state.SourceVolumesSet = true
	if p.exists && p.state.Normalized() == state {
		return nil
	}
	if err := atomicfile.WriteJSON(p.path, state, 0600); err != nil {
		return err
	}
	p.state = state
	p.exists = true
	return nil
}

func (p *Preferences) loadLocked() error {
	if p.loaded {
		return nil
	}
	data, err := os.ReadFile(p.path)
	if errors.Is(err, os.ErrNotExist) {
		p.state = PreferencesState{}
		p.exists = false
		p.loaded = true
		return nil
	}
	if err != nil {
		return err
	}
	var state PreferencesState
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}
	p.state = state
	p.exists = true
	p.loaded = true
	return nil
}
