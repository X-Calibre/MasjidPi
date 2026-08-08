package playback

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/player"
	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
)

type fakePlayer struct {
	mu sync.Mutex

	playCalls []string
	stopCalls int
	volume    int
	statuses  []*player.Status
}

func (f *fakePlayer) Play(url string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.playCalls = append(f.playCalls, url)
	return nil
}

func (f *fakePlayer) Stop() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.stopCalls++
	return nil
}

func (f *fakePlayer) Volume(volume int) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.volume = volume
	return nil
}

func (f *fakePlayer) Status() (*player.Status, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.statuses) == 0 {
		return &player.Status{State: "playing", URL: "relay://one", Volume: 70}, nil
	}

	status := f.statuses[0]
	f.statuses = f.statuses[1:]

	return status, nil
}

func (f *fakePlayer) playCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()

	return len(f.playCalls)
}

func TestStopDisablesListeningAndKeepsSelectedStream(t *testing.T) {
	manager := New(&fakePlayer{}, Config{})
	selected := stream.Stream{ID: "one", Name: "Masjid One", URL: "relay://one"}

	manager.Play(selected)
	manager.Stop()

	status := manager.Status()
	if status.State != string(StateIdle) {
		t.Fatalf("state = %q, want %q", status.State, StateIdle)
	}

	if status.Listening {
		t.Fatal("listening = true, want false")
	}

	if status.StreamID != selected.ID {
		t.Fatalf("stream ID = %q, want %q", status.StreamID, selected.ID)
	}
}

func TestManagerRetriesAfterPlaybackStops(t *testing.T) {
	fake := &fakePlayer{
		statuses: []*player.Status{
			{State: "stopped", URL: "relay://one", Volume: 70},
		},
	}

	manager := New(fake, Config{
		RetryInterval:       10 * time.Millisecond,
		ReconnectDelay:      10 * time.Millisecond,
		StatusCheckInterval: 10 * time.Millisecond,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	manager.Start(ctx)
	manager.Play(stream.Stream{ID: "one", Name: "Masjid One", URL: "relay://one"})

	waitFor(t, time.Second, func() bool {
		return fake.playCount() >= 2
	})
}

func TestBackoffDelayDoublesAndCaps(t *testing.T) {
	base := 5 * time.Second
	want := []time.Duration{
		5 * time.Second,
		10 * time.Second,
		20 * time.Second,
		40 * time.Second,
		80 * time.Second,
		160 * time.Second,
		5 * time.Minute,
	}

	for attempt, expected := range want {
		if got := backoffDelay(base, attempt); got != expected {
			t.Fatalf("attempt %d: delay = %s, want %s", attempt, got, expected)
		}
	}
}

func waitFor(t *testing.T, timeout time.Duration, ok func() bool) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if ok() {
			return
		}

		time.Sleep(5 * time.Millisecond)
	}

	t.Fatal("condition was not met before timeout")
}
