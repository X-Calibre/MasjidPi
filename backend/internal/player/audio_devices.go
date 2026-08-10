package player

import "strings"

// normalizeAudioDevices converts mpv's low-level audio-device list into a
// small set of user-facing choices while retaining the real mpv device name
// in AudioDevice.Name for playback and persistence.
func normalizeAudioDevices(items []map[string]any) []AudioDevice {
	devices := make([]AudioDevice, 0, len(items))
	seen := make(map[string]bool)

	for _, item := range items {
		name, _ := item["name"].(string)
		description, _ := item["description"].(string)
		if name == "" {
			continue
		}

		label, ok := userFacingAudioDevice(name, description)
		if !ok || seen[label] {
			continue
		}

		seen[label] = true
		devices = append(devices, AudioDevice{
			Name:        name,
			Description: label,
		})
	}

	return devices
}

func userFacingAudioDevice(name, description string) (string, bool) {
	if name == "auto" {
		return "Default audio output", true
	}

	// MasjidPi uses ALSA. Hide other mpv audio-output backends and their
	// backend-specific defaults from the end user.
	if !strings.HasPrefix(name, "alsa/") {
		return "", false
	}

	lowerName := strings.ToLower(name)
	lowerDescription := strings.ToLower(description)
	if strings.HasSuffix(lowerName, "/default") || strings.HasPrefix(lowerDescription, "default (") {
		return "", false
	}

	label := strings.TrimSpace(description)
	if slash := strings.Index(label, "/"); slash >= 0 {
		label = strings.TrimSpace(label[:slash])
	}

	// Some ALSA profiles describe the same physical device twice, e.g.
	// "bcm2835 Headphones, bcm2835 Headphones". Collapse that duplication.
	parts := strings.Split(label, ",")
	if len(parts) == 2 && strings.EqualFold(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])) {
		label = strings.TrimSpace(parts[0])
	}

	lower := strings.ToLower(label)
	switch {
	case strings.Contains(lower, "usb audio"):
		return "USB Audio", true
	case strings.Contains(lower, "hdmi"):
		return "HDMI Audio", true
	case strings.Contains(lower, "headphones"):
		return "Headphones", true
	case strings.Contains(lower, "analog"):
		return "Analog Audio", true
	}

	if label == "" {
		return "", false
	}
	return label, true
}
