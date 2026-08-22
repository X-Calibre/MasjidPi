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

	playCalls   []string
	stopCalls   int
	volume      int
	volumeCalls []int
	statuses    []*player.Status
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
	f.volumeCalls = append(f.volumeCalls, volume)
	return nil
}

func (f *fakePlayer) Status() (*player.Status, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.statuses) == 0 {
		return &player.Status{State: "playing", URL: "relay://one", Volume: 70, VolumeSupported: true, AudioDevice: "auto"}, nil
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

func (f *fakePlayer) stopCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.stopCalls
}

func (f *fakePlayer) volumeCallsSnapshot() []int {
	f.mu.Lock()
	defer f.mu.Unlock()

	return append([]int(nil), f.volumeCalls...)
}

type fakePersistence struct {
	mu      sync.Mutex
	saved   []string
	cleared int
}

func (f *fakePersistence) Save(streamID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.saved = append(f.saved, streamID)
	return nil
}

func (f *fakePersistence) Clear() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.cleared++
	return nil
}

func (f *fakePersistence) lastSaved() string {
	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.saved) == 0 {
		return ""
	}
	return f.saved[len(f.saved)-1]
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

func TestManagerPersistsPlayAndClearsStop(t *testing.T) {
	persistence := &fakePersistence{}
	manager := New(&fakePlayer{}, Config{})
	manager.SetPersistence(persistence)
	selected := stream.Stream{ID: "one", Name: "Masjid One", URL: "relay://one"}

	manager.Play(selected)
	if got := persistence.lastSaved(); got != selected.ID {
		t.Fatalf("saved stream ID = %q, want %q", got, selected.ID)
	}

	manager.Stop()
	if persistence.cleared != 1 {
		t.Fatalf("clear calls = %d, want 1", persistence.cleared)
	}
}

func TestManagerSwitchesStreamsDuringPlayback(t *testing.T) {
	fake := &fakePlayer{}
	manager := New(fake, Config{
		StartupGracePeriod:  200 * time.Millisecond,
		StatusCheckInterval: 10 * time.Millisecond,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	manager.Start(ctx)
	manager.Play(stream.Stream{ID: "one", Name: "Masjid One", URL: "relay://one"})

	waitFor(t, time.Second, func() bool {
		return fake.playCount() == 1
	})

	manager.Play(stream.Stream{ID: "two", Name: "Masjid Two", URL: "relay://two"})

	waitFor(t, time.Second, func() bool {
		return fake.playCount() == 2
	})

	if got := fake.stopCount(); got != 1 {
		t.Fatalf("stop calls = %d, want 1", got)
	}

	status := manager.Status()
	if status.StreamID != "two" {
		t.Fatalf("selected stream ID = %q, want %q", status.StreamID, "two")
	}
	if status.URL != "relay://two" {
		t.Fatalf("selected stream URL = %q, want %q", status.URL, "relay://two")
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
		StartupGracePeriod:  1 * time.Millisecond,
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

func TestManagerRestoresVolumeAfterPlayerRecovery(t *testing.T) {
	fake := &fakePlayer{
		statuses: []*player.Status{
			{State: "stopped", URL: "relay://one", Volume: 38, VolumeSupported: true, AudioDevice: "auto"},
			{State: "playing", URL: "relay://one", Volume: 100, VolumeSupported: true, AudioDevice: "auto"},
		},
	}

	manager := New(fake, Config{
		RetryInterval:       10 * time.Millisecond,
		StartupGracePeriod:  1 * time.Millisecond,
		StatusCheckInterval: 10 * time.Millisecond,
	})

	if err := manager.Volume(38); err != nil {
		t.Fatalf("set volume: %v", err)
	}
	fake.mu.Lock()
	fake.volumeCalls = nil
	fake.volume = 0
	fake.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	manager.Start(ctx)
	manager.Play(stream.Stream{ID: "one", Name: "Masjid One", URL: "relay://one"})

	waitFor(t, time.Second, func() bool {
		calls := fake.volumeCallsSnapshot()
		return len(calls) > 0 && calls[len(calls)-1] == 38
	})

	status := manager.Status()
	if status.Volume != 38 {
		t.Fatalf("manager volume = %d, want 38", status.Volume)
	}
}

func TestManagerWaitsForLiveStatusBeforePlaying(t *testing.T) {
	fake := &fakePlayer{}
	manager := New(fake, Config{
		StartupGracePeriod:  200 * time.Millisecond,
		StatusCheckInterval: 10 * time.Millisecond,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	manager.Start(ctx)
	manager.Play(stream.Stream{ID: "one", Name: "Masjid One", URL: "relay://one"})

	waitFor(t, time.Second, func() bool {
		return fake.playCount() == 1
	})

	if status := manager.Status(); status.State != string(StateConnecting) {
		t.Fatalf("state = %q, want %q", status.State, StateConnecting)
	}

	waitFor(t, time.Second, func() bool {
		return manager.Status().State == string(StatePlaying)
	})

	if got := fake.playCount(); got != 1 {
		t.Fatalf("play calls = %d, want 1", got)
	}
}

func TestManagerDoesNotRetryDuringStartupGracePeriod(t *testing.T) {
	fake := &fakePlayer{}
	for i := 0; i < 20; i++ {
		fake.statuses = append(fake.statuses, &player.Status{
			State:  "stopped",
			URL:    "relay://one",
			Volume: 70,
		})
	}

	manager := New(fake, Config{
		RetryInterval:       10 * time.Millisecond,
		StartupGracePeriod:  100 * time.Millisecond,
		StatusCheckInterval: 10 * time.Millisecond,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	manager.Start(ctx)
	manager.Play(stream.Stream{ID: "one", Name: "Masjid One", URL: "relay://one"})

	waitFor(t, 50*time.Millisecond, func() bool {
		return manager.Status().State == string(StateConnecting)
	})

	if got := fake.playCount(); got != 1 {
		t.Fatalf("play calls during grace period = %d, want 1", got)
	}

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
