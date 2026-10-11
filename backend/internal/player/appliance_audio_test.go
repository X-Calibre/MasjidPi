package player

import "testing"

func TestApplianceAudioDevices(t *testing.T) {
	devices := []AudioDevice{
		{Name: "auto", Description: "Default audio output"},
		{Name: "alsa/plughw:CARD=vc4hdmi,DEV=0", Description: "HDMI Audio"},
		{Name: "alsa/plughw:CARD=USB,DEV=0", Description: "Headphones"},
		{Name: "alsa/plughw:CARD=Headphones,DEV=0", Description: "Headphones"},
		{Name: "alsa/plughw:CARD=Device,DEV=0", Description: "MasjidFrame Built-In Audio"},
		{Name: "alsa/plughw:CARD=Device,DEV=1", Description: "MasjidFrame Built-In Audio"},
	}
	got := ApplianceAudioDevices(devices)
	if len(got) != 2 || got[0].Name != devices[4].Name || got[1].Name != devices[3].Name ||
		got[0].Description != "MasjidFrame Built-In Audio" || got[1].Description != "External Stereo Output" {
		t.Fatalf("unexpected appliance outputs: %#v", got)
	}
	if devices[3].Description != "Headphones" {
		t.Fatal("changed original device list")
	}
	devices[4].Unavailable = true
	devices[5].Unavailable = true
	got = ApplianceAudioDevices(devices)
	if len(got) != 1 || got[0].Description != "External Stereo Output" {
		t.Fatalf("unavailable built-in output should not be selectable: %#v", got)
	}
	if got := ApplianceAudioDevices(nil); len(got) != 0 {
		t.Fatalf("invented absent devices: %#v", got)
	}
}
