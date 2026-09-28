package app

import (
	"testing"

	"github.com/X-Calibre/MasjidFrame/backend/internal/player"
)

func TestPreferredFirstRunAudioDeviceSelectsFirstAvailableUSBDevice(t *testing.T) {
	devices := []player.AudioDevice{
		{Name: "auto", Description: "Default audio output"},
		{Name: "alsa/plughw:CARD=vc4hdmi,DEV=0", Description: "HDMI Audio"},
		{Name: "alsa/plughw:CARD=Device,DEV=0", Description: "USB Audio"},
		{Name: "alsa/plughw:CARD=OtherUSB,DEV=0", Description: "USB Audio"},
	}

	name, ok := preferredFirstRunAudioDevice(devices)
	if !ok {
		t.Fatal("preferredFirstRunAudioDevice() found no device")
	}
	if want := "alsa/plughw:CARD=Device,DEV=0"; name != want {
		t.Fatalf("preferredFirstRunAudioDevice() = %q, want %q", name, want)
	}
}

func TestPreferredFirstRunAudioDevicePrefersBuiltInAudio(t *testing.T) {
	devices := []player.AudioDevice{
		{Name: "auto", Description: "Default audio output"},
		{Name: "alsa/plughw:CARD=OtherUSB,DEV=0", Description: "USB Audio"},
		{Name: "alsa/plughw:CARD=Device,DEV=0", Description: "MasjidFrame Built-In Audio"},
	}

	name, ok := preferredFirstRunAudioDevice(devices)
	if !ok || name != "alsa/plughw:CARD=Device,DEV=0" {
		t.Fatalf("preferredFirstRunAudioDevice() = %q, %t, want built-in audio", name, ok)
	}

	devices[2].Unavailable = true
	name, ok = preferredFirstRunAudioDevice(devices)
	if !ok || name != "alsa/plughw:CARD=OtherUSB,DEV=0" {
		t.Fatalf("unavailable built-in device: got %q, %t, want USB fallback", name, ok)
	}
}

func TestPreferredFirstRunAudioDeviceIgnoresUnavailableUSBDevice(t *testing.T) {
	devices := []player.AudioDevice{
		{
			Name:        "alsa/plughw:CARD=Missing,DEV=0",
			Description: "USB Audio",
			Unavailable: true,
		},
		{Name: "alsa/plughw:CARD=Headphones,DEV=0", Description: "Headphones"},
	}

	if name, ok := preferredFirstRunAudioDevice(devices); ok {
		t.Fatalf("preferredFirstRunAudioDevice() = %q, want no selection", name)
	}
}

func TestPreferredFirstRunAudioDeviceRetainsAutomaticOutputWithoutUSB(t *testing.T) {
	devices := []player.AudioDevice{
		{Name: "auto", Description: "Default audio output"},
		{Name: "alsa/plughw:CARD=vc4hdmi,DEV=0", Description: "HDMI Audio"},
		{Name: "alsa/plughw:CARD=Headphones,DEV=0", Description: "Headphones"},
	}

	if name, ok := preferredFirstRunAudioDevice(devices); ok {
		t.Fatalf("preferredFirstRunAudioDevice() = %q, want automatic output", name)
	}
}
