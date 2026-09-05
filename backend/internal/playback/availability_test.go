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
	updatedAt time.Time
	events    chan string
}

func (f *fakeAvailability) Status(string) (bool, bool, time.Time) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.available, f.known, f.updatedAt
}

func (f *fakeAvailability) Events() <-chan string { return f.events }

func (f *fakeAvailability) set(available bool) {
	f.mu.Lock()
	f.available = available
	f.known = true
	f.updatedAt = time.Now().Add(-DefaultMountStartupDelay)
	f.mu.Unlock()
	f.events <- "one"
}

func (f *fakeAvailability) setStartedNow() {
	f.mu.Lock()
	f.available = true
	f.known = true
	f.updatedAt = time.Now()
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

	if got := fake.stopCount(); got != 1 {
		t.Fatalf("stop calls = %d, want 1", got)
	}

	plays := fake.playCount()
	time.Sleep(50 * time.Millisecond)
	if fake.playCount() != plays {
		t.Fatalf("play calls increased while availability was unknown: %d -> %d", plays, fake.playCount())
	}
}

func TestManagerWaitsForMountStartupDelay(t *testing.T) {
	fake := &fakePlayer{}
	availability := &fakeAvailability{events: make(chan string, 4)}
	manager := New(fake, Config{
		MountStartupDelay:   40 * time.Millisecond,
		StatusCheckInterval: 10 * time.Millisecond,
	})
	manager.SetAvailability(availability)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	manager.Start(ctx)
	manager.Play(stream.Stream{ID: "one", Name: "Masjid One", URL: "relay://one"})
	availability.setStartedNow()

	time.Sleep(20 * time.Millisecond)
	if got := fake.playCount(); got != 0 {
		t.Fatalf("play calls during mount startup delay = %d, want 0", got)
	}
	// The race detector can delay this goroutine substantially on a busy CI
	// runner. Keep the production startup delay small while allowing enough
	// wall-clock time for the asynchronous assertion to observe it.
	waitFor(t, 3*time.Second, func() bool { return fake.playCount() == 1 })
}

func TestManagerCancelsMountStartupWhenMountStops(t *testing.T) {
	fake := &fakePlayer{}
	availability := &fakeAvailability{events: make(chan string, 4)}
	manager := New(fake, Config{
		MountStartupDelay:   80 * time.Millisecond,
		StatusCheckInterval: 10 * time.Millisecond,
	})
	manager.SetAvailability(availability)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	manager.Start(ctx)
	manager.Play(stream.Stream{ID: "one", Name: "Masjid One", URL: "relay://one"})
	availability.setStartedNow()
	time.Sleep(20 * time.Millisecond)
	availability.set(false)
	time.Sleep(100 * time.Millisecond)

	if got := fake.playCount(); got != 0 {
		t.Fatalf("play calls after mount stopped during startup delay = %d, want 0", got)
	}
}
