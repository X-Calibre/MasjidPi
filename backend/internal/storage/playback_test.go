package storage

import (
	"path/filepath"
	"testing"
)

func TestPlaybackSaveLoadAndClear(t *testing.T) {
	path := filepath.Join(t.TempDir(), "playback.json")
	state := NewPlayback(path)

	if _, ok, err := state.Load(); err != nil || ok {
		t.Fatalf("initial Load() = (%v, %v), want empty state", err, ok)
	}

	if err := state.Save("activetakbeer"); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, ok, err := state.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !ok || got != "activetakbeer" {
		t.Fatalf("Load() = (%q, %v), want (%q, true)", got, ok, "activetakbeer")
	}

	if err := state.Clear(); err != nil {
		t.Fatalf("Clear() error = %v", err)
	}

	if _, ok, err := state.Load(); err != nil || ok {
		t.Fatalf("Load() after Clear() = (%v, %v), want empty state", err, ok)
	}
}
