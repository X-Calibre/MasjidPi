package playback

import "testing"

func TestManagerAudioDevice(t *testing.T) {
	manager := New(&fakePlayer{}, Config{})

	if err := manager.AudioDevice("auto"); err != nil {
		t.Fatalf("AudioDevice() error = %v", err)
	}

	if got := manager.Status().AudioDevice; got != "auto" {
		t.Fatalf("audio device = %q, want %q", got, "auto")
	}
}
