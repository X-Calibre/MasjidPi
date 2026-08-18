package storage

import (
	"os"
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

func TestPlaybackLoadUsesCachedStateAfterFirstRead(t *testing.T) {
	path := filepath.Join(t.TempDir(), "playback.json")
	if err := os.WriteFile(path, []byte(`{"stream_id":"one"}`), 0600); err != nil {
		t.Fatalf("write initial state: %v", err)
	}

	state := NewPlayback(path)
	got, ok, err := state.Load()
	if err != nil || !ok || got != "one" {
		t.Fatalf("initial Load() = (%q, %v, %v), want (one, true, nil)", got, ok, err)
	}

	if err := os.WriteFile(path, []byte(`{"stream_id":"two"}`), 0600); err != nil {
		t.Fatalf("replace backing state: %v", err)
	}

	got, ok, err = state.Load()
	if err != nil || !ok || got != "one" {
		t.Fatalf("cached Load() = (%q, %v, %v), want cached stream one", got, ok, err)
	}
}
