package playback

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
)

type fakeAvailability struct {
	mu        sync.Mutex
	available bool
	known     bool
	events    chan string
}

func (f *fakeAvailability) IsAvailable(string) (bool, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.available, f.known
}

func (f *fakeAvailability) Events() <-chan string { return f.events }

func (f *fakeAvailability) set(available bool) {
	f.mu.Lock()
	f.available = available
	f.known = true
	f.mu.Unlock()
	f.events <- "one"
}

func (f *fakeAvailability) setUnknown() {
	f.mu.Lock()
	f.known = false
	f.mu.Unlock()
	f.events <- "one"
}

func TestManagerRespondsToAvailabilityEvents(t *testing.T) {
	fake := &fakePlayer{}
	availability := &fakeAvailability{events: make(chan string, 4)}
	manager := New(fake, Config{StatusCheckInterval: 10 * time.Millisecond})
	manager.SetAvailability(availability)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	manager.Start(ctx)
	manager.Play(stream.Stream{ID: "one", Name: "Masjid One", URL: "relay://one"})

	time.Sleep(50 * time.Millisecond)
	if fake.playCount() != 0 {
		t.Fatalf("play calls = %d, want 0 while stream is unknown/offline", fake.playCount())
	}

	availability.set(true)
	waitFor(t, time.Second, func() bool { return fake.playCount() == 1 })
}

func TestManagerStopsRetryingWhenMountGoesOffline(t *testing.T) {
	fake := &fakePlayer{}
	availability := &fakeAvailability{available: true, known: true, events: make(chan string, 4)}
	manager := New(fake, Config{
		RetryInterval:       20 * time.Millisecond,
		ReconnectDelay:      20 * time.Millisecond,
		StatusCheckInterval: 10 * time.Millisecond,
	})
	manager.SetAvailability(availability)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	manager.Start(ctx)
	manager.Play(stream.Stream{ID: "one", Name: "Masjid One", URL: "relay://one"})

	waitFor(t, time.Second, func() bool { return fake.playCount() >= 1 })
	availability.set(false)

	time.Sleep(80 * time.Millisecond)
	plays := fake.playCount()
	time.Sleep(80 * time.Millisecond)

	if fake.playCount() != plays {
		t.Fatalf("play calls increased while mount was offline: %d -> %d", plays, fake.playCount())
	}
}

func TestManagerStopsPlaybackWhenAvailabilityBecomesUnknown(t *testing.T) {
	fake := &fakePlayer{}
	availability := &fakeAvailability{available: true, known: true, events: make(chan string, 4)}
	manager := New(fake, Config{
		StatusCheckInterval: 10 * time.Millisecond,
	})
	manager.SetAvailability(availability)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	manager.Start(ctx)
	manager.Play(stream.Stream{ID: "one", Name: "Masjid One", URL: "relay://one"})

	waitFor(t, time.Second, func() bool {
		return manager.Status().State == string(StatePlaying)
	})

	availability.setUnknown()

	waitFor(t, time.Second, func() bool {
		return manager.Status().State == string(StateWaiting)
	})

	if fake.stopCalls != 1 {
		t.Fatalf("stop calls = %d, want 1", fake.stopCalls)
	}

	plays := fake.playCount()
	time.Sleep(50 * time.Millisecond)
	if fake.playCount() != plays {
		t.Fatalf("play calls increased while availability was unknown: %d -> %d", plays, fake.playCount())
	}
}
