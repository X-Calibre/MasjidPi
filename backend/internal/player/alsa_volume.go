package player

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var ErrHardwareVolumeUnsupported = errors.New("audio device does not expose a controllable hardware volume")

var mixerVolumePattern = regexp.MustCompile(`\[(\d+)%\]`)

var preferredMixerControls = []string{"Master", "Speaker", "Headphone", "PCM"}

type ALSAVolume struct {
	run func(name string, args ...string) ([]byte, error)
}

func NewALSAVolume() *ALSAVolume {
	return &ALSAVolume{run: func(name string, args ...string) ([]byte, error) {
		return exec.Command(name, args...).CombinedOutput()
	}}
}

func (a *ALSAVolume) Get(device string) (int, bool, error) {
	mixer, ok := mixerDevice(device)
	if !ok {
		return 100, false, nil
	}

	control, output, err := a.readControl(mixer)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return 100, false, nil
		}
		return 100, false, err
	}
	_ = control

	match := mixerVolumePattern.FindStringSubmatch(string(output))
	if len(match) != 2 {
		return 100, false, fmt.Errorf("unable to parse ALSA volume from amixer output")
	}
	volume, err := strconv.Atoi(match[1])
	if err != nil {
		return 100, false, fmt.Errorf("parse ALSA volume: %w", err)
	}
	return volume, true, nil
}

func (a *ALSAVolume) Set(device string, volume int) (bool, error) {
	if volume < 0 || volume > 100 {
		return false, fmt.Errorf("volume must be between 0 and 100")
	}

	mixer, ok := mixerDevice(device)
	if !ok {
		return false, nil
	}

	control, _, err := a.readControl(mixer)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return false, nil
		}
		return false, err
	}

	output, err := a.run("amixer", "-D", mixer, "sset", control, fmt.Sprintf("%d%%", volume))
	if err != nil {
		return false, fmt.Errorf("set ALSA volume: %s", strings.TrimSpace(string(output)))
	}
	return true, nil
}

func (a *ALSAVolume) readControl(mixer string) (string, []byte, error) {
	for _, control := range preferredMixerControls {
		output, err := a.run("amixer", "-D", mixer, "sget", control)
		if err == nil && mixerVolumePattern.Match(output) {
			return control, output, nil
		}
		if errors.Is(err, exec.ErrNotFound) {
			return "", nil, err
		}
	}
	return "", nil, ErrHardwareVolumeUnsupported
}

func mixerDevice(device string) (string, bool) {
	device = strings.TrimPrefix(device, "alsa/")
	if strings.HasPrefix(device, "hw:CARD=") {
		return strings.SplitN(device, ",", 2)[0], true
	}
	if strings.HasPrefix(device, "plughw:CARD=") {
		return "hw:" + strings.TrimPrefix(strings.SplitN(device, ",", 2)[0], "plughw:"), true
	}
	if strings.HasPrefix(device, "hw:") {
		return strings.SplitN(device, ",", 2)[0], true
	}
	if strings.HasPrefix(device, "plughw:") {
		return "hw:" + strings.TrimPrefix(strings.SplitN(device, ",", 2)[0], "plughw:"), true
	}
	return "", false
}
