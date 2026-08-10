package storage

import "testing"

func TestAudioDeviceStateRoundTrip(t *testing.T) {
	path := t.TempDir() + "/audio_device.json"
	state := NewAudioDeviceState(path)

	if name, ok, err := state.Load(); err != nil || ok || name != "" {
		t.Fatalf("initial Load() = %q, %v, %v", name, ok, err)
	}

	if err := state.Save("alsa/plughw:CARD=Headset,DEV=0"); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	name, ok, err := state.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !ok {
		t.Fatal("Load() ok = false, want true")
	}
	if name != "alsa/plughw:CARD=Headset,DEV=0" {
		t.Fatalf("Load() name = %q, want saved device", name)
	}
}
