package storage

import (
	"os"
	"testing"
)

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

func TestAudioDeviceStateLoadUsesCachedValue(t *testing.T) {
	path := t.TempDir() + "/audio_device.json"
	if err := os.WriteFile(path, []byte(`{"name":"alsa/device-one"}`), 0644); err != nil {
		t.Fatalf("write initial state: %v", err)
	}

	state := NewAudioDeviceState(path)
	name, ok, err := state.Load()
	if err != nil {
		t.Fatalf("initial Load() error = %v", err)
	}
	if !ok || name != "alsa/device-one" {
		t.Fatalf("initial Load() = %q, %v; want cached device one", name, ok)
	}

	// Change the backing file behind the state object. A second Load must return
	// the in-memory value rather than re-reading persistent storage.
	if err := os.WriteFile(path, []byte(`{"name":"alsa/device-two"}`), 0644); err != nil {
		t.Fatalf("replace backing state: %v", err)
	}

	name, ok, err = state.Load()
	if err != nil {
		t.Fatalf("cached Load() error = %v", err)
	}
	if !ok || name != "alsa/device-one" {
		t.Fatalf("cached Load() = %q, %v; want device one", name, ok)
	}
}
