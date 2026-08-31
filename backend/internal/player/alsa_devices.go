package player

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var alsaPlaybackPCM = regexp.MustCompile(`^pcmC([0-9]+)D([0-9]+)p$`)

type AudioDeviceProvider interface {
	AudioDevices() ([]AudioDevice, error)
}

// ALSAAudioDevices discovers playback devices from Linux sysfs on every call.
// MPV's audio-device-list is a startup snapshot and does not reliably reflect
// USB audio hot-plug events.
type ALSAAudioDevices struct {
	soundClassPath string
	fallback       AudioDeviceProvider
}

func NewALSAAudioDevices(soundClassPath string, fallback AudioDeviceProvider) *ALSAAudioDevices {
	return &ALSAAudioDevices{soundClassPath: soundClassPath, fallback: fallback}
}

func (d *ALSAAudioDevices) AudioDevices() ([]AudioDevice, error) {
	devices, err := discoverALSAAudioDevices(d.soundClassPath)
	if err == nil {
		return devices, nil
	}
	if d.fallback != nil {
		return d.fallback.AudioDevices()
	}
	if err != nil {
		return nil, err
	}
	return devices, nil
}

type alsaPCM struct {
	card   int
	device int
	cardID string
	label  string
}

func discoverALSAAudioDevices(soundClassPath string) ([]AudioDevice, error) {
	entries, err := os.ReadDir(soundClassPath)
	if err != nil {
		return nil, fmt.Errorf("read ALSA sound devices: %w", err)
	}

	pcms := make([]alsaPCM, 0)
	for _, entry := range entries {
		matches := alsaPlaybackPCM.FindStringSubmatch(entry.Name())
		if matches == nil {
			continue
		}
		card, _ := strconv.Atoi(matches[1])
		device, _ := strconv.Atoi(matches[2])
		cardPath := filepath.Join(soundClassPath, fmt.Sprintf("card%d", card))
		idBytes, readErr := os.ReadFile(filepath.Join(cardPath, "id"))
		if readErr != nil {
			continue
		}
		cardID := strings.TrimSpace(string(idBytes))
		if cardID == "" {
			continue
		}
		pcms = append(pcms, alsaPCM{
			card:   card,
			device: device,
			cardID: cardID,
			label:  alsaCardLabel(cardPath, cardID),
		})
	}

	sort.Slice(pcms, func(i, j int) bool {
		if pcms[i].card == pcms[j].card {
			return pcms[i].device < pcms[j].device
		}
		return pcms[i].card < pcms[j].card
	})

	devices := []AudioDevice{{Name: "auto", Description: "Default audio output"}}
	labelCounts := make(map[string]int)
	for _, pcm := range pcms {
		label := pcm.label
		labelCounts[label]++
		if labelCounts[label] > 1 {
			label = fmt.Sprintf("%s (%s, device %d)", label, pcm.cardID, pcm.device)
		}
		devices = append(devices, AudioDevice{
			Name:        fmt.Sprintf("alsa/plughw:CARD=%s,DEV=%d", pcm.cardID, pcm.device),
			Description: label,
		})
	}
	return devices, nil
}

func alsaCardLabel(cardPath, cardID string) string {
	lowerID := strings.ToLower(cardID)
	switch {
	case strings.Contains(lowerID, "hdmi"):
		return "HDMI Audio"
	case strings.Contains(lowerID, "headphone"):
		return "Headphones"
	case strings.Contains(lowerID, "analog"):
		return "Analog Audio"
	}

	uevent, _ := os.ReadFile(filepath.Join(cardPath, "device", "uevent"))
	if strings.Contains(strings.ToLower(string(uevent)), "snd-usb-audio") {
		return "USB Audio"
	}
	return cardID
}

func AudioDeviceDescription(name string) string {
	const marker = "CARD="
	start := strings.Index(name, marker)
	if start < 0 {
		return name
	}
	card := name[start+len(marker):]
	if comma := strings.IndexByte(card, ','); comma >= 0 {
		card = card[:comma]
	}
	lower := strings.ToLower(card)
	switch {
	case strings.Contains(lower, "hdmi"):
		return "HDMI Audio"
	case strings.Contains(lower, "headphone"):
		return "Headphones"
	case strings.Contains(lower, "analog"):
		return "Analog Audio"
	case strings.Contains(lower, "usb"), strings.Contains(lower, "uac"):
		return "USB Audio"
	}
	return card
}
