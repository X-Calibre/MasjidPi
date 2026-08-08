package playback

import (
	"testing"
	"time"
)

func TestBackoffDelayProgressesAndCaps(t *testing.T) {
	base := 5 * time.Second
	want := []time.Duration{
		5 * time.Second,
		10 * time.Second,
		20 * time.Second,
		40 * time.Second,
		80 * time.Second,
		160 * time.Second,
		5 * time.Minute,
		5 * time.Minute,
	}

	for attempt, expected := range want {
		if got := backoffDelay(base, attempt); got != expected {
			t.Fatalf("attempt %d: backoffDelay() = %s, want %s", attempt, got, expected)
		}
	}
}
