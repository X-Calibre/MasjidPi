package storage

import (
	"os"
	"testing"
)

func TestVolumeSaveLoadPerDevice(t *testing.T) {
	path := t.TempDir() + "/volume.json"
	state := NewVolume(path)

	if err := state.Save("alsa/plughw:CARD=USB,DEV=0", 85); err != nil {
		t.Fatalf("save volume: %v", err)
	}
	if err := state.Save("alsa/plughw:CARD=HDMI,DEV=0", 40); err != nil {
		t.Fatalf("save second volume: %v", err)
	}

	volume, ok, err := state.Load("alsa/plughw:CARD=USB,DEV=0")
	if err != nil || !ok || volume != 85 {
		t.Fatalf("expected USB volume 85, got volume=%d ok=%v err=%v", volume, ok, err)
	}

	volume, ok, err = state.Load("alsa/plughw:CARD=HDMI,DEV=0")
	if err != nil || !ok || volume != 40 {
		t.Fatalf("expected HDMI volume 40, got volume=%d ok=%v err=%v", volume, ok, err)
	}
}

func TestVolumeLoadMissing(t *testing.T) {
	state := NewVolume(t.TempDir() + "/volume.json")

	volume, ok, err := state.Load("alsa/plughw:CARD=USB,DEV=0")
	if err != nil {
		t.Fatalf("load missing volume: %v", err)
	}
	if ok || volume != 0 {
		t.Fatalf("expected no saved volume, got volume=%d ok=%v", volume, ok)
	}
}

func TestVolumeMigratesLegacyState(t *testing.T) {
	path := t.TempDir() + "/volume.json"
	if err := os.WriteFile(path, []byte(`{"volume":30}`), 0600); err != nil {
		t.Fatalf("write legacy state: %v", err)
	}

	state := NewVolume(path)
	volume, ok, err := state.Load("alsa/plughw:CARD=USB,DEV=0")
	if err != nil || !ok || volume != 30 {
		t.Fatalf("expected legacy volume 30, got volume=%d ok=%v err=%v", volume, ok, err)
	}
}

func TestVolumeRejectsOutOfRange(t *testing.T) {
	state := NewVolume(t.TempDir() + "/volume.json")
	if err := state.Save("alsa/plughw:CARD=USB,DEV=0", 101); err == nil {
		t.Fatal("expected out-of-range volume to fail")
	}
}

func TestVolumeSaveSkipsUnchangedValue(t *testing.T) {
	path := t.TempDir() + "/volume.json"
	state := NewVolume(path)
	device := "alsa/plughw:CARD=USB,DEV=0"

	if err := state.Save(device, 85); err != nil {
		t.Fatalf("initial save: %v", err)
	}

	before, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat before unchanged save: %v", err)
	}

	if err := state.Save(device, 85); err != nil {
		t.Fatalf("unchanged save: %v", err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat after unchanged save: %v", err)
	}
	if !os.SameFile(before, after) {
		t.Fatal("unchanged volume replaced the state file")
	}
}
