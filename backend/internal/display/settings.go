package display

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/X-Calibre/MasjidPi/backend/internal/atomicfile"
)

const (
	DefaultBrightness = 100
	MinBrightness     = 10
	MaxBrightness     = 100
)

type State struct {
	BrightnessPercent int `json:"brightness_percent"`
}

type Settings struct {
	State
	BrightnessAvailable bool `json:"brightness_available"`
}

type Controller struct {
	mu            sync.Mutex
	statePath     string
	backlightRoot string
}

func NewController(statePath, backlightRoot string) *Controller {
	if backlightRoot == "" {
		backlightRoot = "/sys/class/backlight"
	}
	return &Controller{statePath: statePath, backlightRoot: backlightRoot}
}

func normalize(state State) State {
	if state.BrightnessPercent < MinBrightness || state.BrightnessPercent > MaxBrightness {
		state.BrightnessPercent = DefaultBrightness
	}
	return state
}

func (c *Controller) Load() (Settings, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.loadLocked()
}

func (c *Controller) Update(brightness *int) (Settings, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	settings, err := c.loadLocked()
	if err != nil {
		return Settings{}, err
	}
	state := settings.State
	if brightness != nil {
		if *brightness < MinBrightness || *brightness > MaxBrightness {
			return Settings{}, fmt.Errorf("brightness must be between %d and %d percent", MinBrightness, MaxBrightness)
		}
		state.BrightnessPercent = *brightness
		if err := c.writeBrightness(state.BrightnessPercent); err != nil {
			return Settings{}, err
		}
	}
	if err := atomicfile.WriteJSON(c.statePath, state, 0600); err != nil {
		return Settings{}, err
	}
	return Settings{State: state, BrightnessAvailable: c.backlightPath() != ""}, nil
}

func (c *Controller) Restore() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	settings, err := c.loadLocked()
	if err != nil {
		return err
	}
	if !settings.BrightnessAvailable {
		return nil
	}
	return c.writeBrightness(settings.BrightnessPercent)
}

func (c *Controller) loadLocked() (Settings, error) {
	state := State{BrightnessPercent: DefaultBrightness}
	data, err := os.ReadFile(c.statePath)
	if err == nil {
		if err := json.Unmarshal(data, &state); err != nil {
			return Settings{}, err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return Settings{}, err
	}
	state = normalize(state)
	path := c.backlightPath()
	if errors.Is(err, os.ErrNotExist) && path != "" {
		if current, readErr := brightnessPercent(path); readErr == nil {
			state.BrightnessPercent = current
		}
	}
	return Settings{State: state, BrightnessAvailable: path != ""}, nil
}

func (c *Controller) backlightPath() string {
	entries, err := os.ReadDir(c.backlightRoot)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		path := filepath.Join(c.backlightRoot, entry.Name())
		if _, err := os.Stat(filepath.Join(path, "brightness")); err != nil {
			continue
		}
		if _, err := os.Stat(filepath.Join(path, "max_brightness")); err == nil {
			return path
		}
	}
	return ""
}

func brightnessPercent(path string) (int, error) {
	current, err := readInt(filepath.Join(path, "brightness"))
	if err != nil {
		return 0, err
	}
	maximum, err := readInt(filepath.Join(path, "max_brightness"))
	if err != nil || maximum <= 0 {
		return 0, errors.New("invalid backlight maximum brightness")
	}
	return int(math.Round(float64(current) * 100 / float64(maximum))), nil
}

func (c *Controller) writeBrightness(percent int) error {
	path := c.backlightPath()
	if path == "" {
		return errors.New("display backlight control is unavailable")
	}
	maximum, err := readInt(filepath.Join(path, "max_brightness"))
	if err != nil || maximum <= 0 {
		return errors.New("invalid backlight maximum brightness")
	}
	value := int(math.Round(float64(maximum) * float64(percent) / 100))
	if value < 1 {
		value = 1
	}
	return os.WriteFile(filepath.Join(path, "brightness"), []byte(strconv.Itoa(value)+"\n"), 0644)
}

func readInt(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}
