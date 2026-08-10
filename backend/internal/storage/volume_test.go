package storage

import "testing"

func TestVolumeSaveLoad(t *testing.T) {
	path := t.TempDir() + "/volume.json"
	state := NewVolume(path)

	if err := state.Save(85); err != nil {
		t.Fatalf("save volume: %v", err)
	}

	volume, ok, err := state.Load()
	if err != nil {
		t.Fatalf("load volume: %v", err)
	}
	if !ok {
		t.Fatal("expected saved volume to be present")
	}
	if volume != 85 {
		t.Fatalf("expected volume 85, got %d", volume)
	}
}

func TestVolumeLoadMissing(t *testing.T) {
	state := NewVolume(t.TempDir() + "/volume.json")

	volume, ok, err := state.Load()
	if err != nil {
		t.Fatalf("load missing volume: %v", err)
	}
	if ok || volume != 0 {
		t.Fatalf("expected no saved volume, got volume=%d ok=%v", volume, ok)
	}
}
