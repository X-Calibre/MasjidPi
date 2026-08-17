package storage

import (
	"os"
	"testing"
)

func assertSameFile(t *testing.T, beforePath string, save func() error) {
	t.Helper()
	before, err := os.Stat(beforePath)
	if err != nil {
		t.Fatalf("stat before save: %v", err)
	}

	if err := save(); err != nil {
		t.Fatalf("unchanged save: %v", err)
	}

	after, err := os.Stat(beforePath)
	if err != nil {
		t.Fatalf("stat after save: %v", err)
	}
	if !os.SameFile(before, after) {
		t.Fatalf("unchanged state replaced %s", beforePath)
	}
}

func TestPreferencesSaveSkipsUnchangedState(t *testing.T) {
	path := t.TempDir() + "/preferences.json"
	state := NewPreferences(path)
	value := PreferencesState{LastStreamID: "stream-1", Autoplay: true}

	if err := state.Save(value); err != nil {
		t.Fatalf("initial save: %v", err)
	}
	assertSameFile(t, path, func() error { return state.Save(value) })
}

func TestAudioDeviceSaveSkipsUnchangedState(t *testing.T) {
	path := t.TempDir() + "/audio-device.json"
	state := NewAudioDeviceState(path)
	device := "alsa/plughw:CARD=USB,DEV=0"

	if err := state.Save(device); err != nil {
		t.Fatalf("initial save: %v", err)
	}
	assertSameFile(t, path, func() error { return state.Save(device) })
}

func TestPlaybackSaveSkipsUnchangedState(t *testing.T) {
	path := t.TempDir() + "/playback.json"
	state := NewPlayback(path)
	streamID := "stream-1"

	if err := state.Save(streamID); err != nil {
		t.Fatalf("initial save: %v", err)
	}
	assertSameFile(t, path, func() error { return state.Save(streamID) })
}

func TestFavouritesSaveSkipsUnchangedState(t *testing.T) {
	path := t.TempDir() + "/favourites.json"
	state := NewFavourites(path)
	ids := []string{"stream-1", "stream-2"}

	if err := state.Save(ids); err != nil {
		t.Fatalf("initial save: %v", err)
	}
	assertSameFile(t, path, func() error { return state.Save(ids) })
}
