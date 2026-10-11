package player

import "strings"

// ApplianceAudioDevices exposes only the two physical appliance outputs.
// Keep the discovered ALSA names so selection does not depend on card numbers.
func ApplianceAudioDevices(devices []AudioDevice) []AudioDevice {
	outputs := make([]AudioDevice, 0, 2)
	for _, label := range []string{"MasjidFrame Built-In Audio", "External Stereo Output"} {
		for _, device := range devices {
			if !strings.HasPrefix(device.Name, "alsa/") || device.Unavailable {
				continue
			}
			matches := device.Description == "MasjidFrame Built-In Audio"
			if label == "External Stereo Output" {
				// The Pi's analogue jack is the bcm2835 Headphones card;
				// do not treat arbitrary USB headphone outputs as the Pi jack.
				matches = strings.Contains(strings.ToLower(device.Name), "card=headphones,") ||
					(strings.Contains(strings.ToLower(device.Name), "bcm2835") && device.Description == "Headphones")
			}
			if matches {
				device.Description = label
				outputs = append(outputs, device)
				break
			}
		}
	}
	return outputs
}
