package player

import "testing"

func TestNormalizeAudioDevices(t *testing.T) {
	items := []map[string]any{
		{"name": "auto", "description": "Autoselect device"},
		{"name": "alsa/bcm2835", "description": "bcm2835 Headphones, bcm2835 Headphones/Hardware device with all software conversions"},
		{"name": "alsa/bcm2835-default", "description": "bcm2835 Headphones, bcm2835 Headphones/Default Audio Device"},
		{"name": "alsa/uac-hw", "description": "UACDemoV1.0, USB Audio/Hardware device with all software conversions"},
		{"name": "alsa/uac-default", "description": "UACDemoV1.0, USB Audio/Default Audio Device"},
		{"name": "alsa/vc4", "description": "vc4-hdmi, MAI PCM i2s-hifi-0/HDMI Audio Output"},
		{"name": "alsa/default", "description": "Default (alsa)"},
		{"name": "jack/default", "description": "Default (jack)"},
		{"name": "sdl/default", "description": "Default (sdl)"},
	}

	got := normalizeAudioDevices(items)
	if len(got) != 4 {
		t.Fatalf("expected 4 user-facing devices, got %d: %#v", len(got), got)
	}

	want := []struct {
		name string
		label string
	}{
		{"auto", "Default audio output"},
		{"alsa/bcm2835", "Headphones"},
		{"alsa/uac-hw", "USB Audio"},
		{"alsa/vc4", "HDMI Audio"},
	}

	for i, expected := range want {
		if got[i].Name != expected.name || got[i].Description != expected.label {
			t.Errorf("device %d: got name=%q description=%q, want name=%q description=%q", i, got[i].Name, got[i].Description, expected.name, expected.label)
		}
	}
}

func TestUserFacingAudioDeviceFiltersNonALSA(t *testing.T) {
	for _, item := range []struct {
		name        string
		description string
	}{
		{"pulse/default", "Default (pulse)"},
		{"sndio/default", "Default (sndio)"},
		{"jack/default", "Default (jack)"},
	} {
		if label, ok := userFacingAudioDevice(item.name, item.description); ok || label != "" {
			t.Errorf("expected %q to be hidden, got label=%q ok=%v", item.name, label, ok)
		}
	}
}
