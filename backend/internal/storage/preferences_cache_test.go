package storage

import (
	"os"
	"sync"
	"testing"
)

func TestPreferencesUpdatePreservesConcurrentFields(t *testing.T) {
	state := NewPreferences(t.TempDir() + "/preferences.json")
	start := make(chan struct{})
	var workers sync.WaitGroup
	workers.Add(2)
	go func() {
		defer workers.Done()
		<-start
		_, _ = state.Update(func(value *PreferencesState) { value.SelectedMasjidID = "masjid-1" })
	}()
	go func() {
		defer workers.Done()
		<-start
		_, _ = state.Update(func(value *PreferencesState) { value.SelectedRadioID = "radio-1" })
	}()
	close(start)
	workers.Wait()

	got, err := state.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.SelectedMasjidID != "masjid-1" || got.SelectedRadioID != "radio-1" {
		t.Fatalf("concurrent updates lost: %+v", got)
	}
}

func TestPreferencesLoadUsesCachedState(t *testing.T) {
	path := t.TempDir() + "/preferences.json"
	state := NewPreferences(path)
	initial := PreferencesState{LastStreamID: "stream-1", Autoplay: true}

	if err := state.Save(initial); err != nil {
		t.Fatalf("initial save: %v", err)
	}

	if err := os.WriteFile(path, []byte(`{"last_stream_id":"stream-2","autoplay":false}`), 0600); err != nil {
		t.Fatalf("replace backing file: %v", err)
	}

	got, err := state.Load()
	if err != nil {
		t.Fatalf("load cached preferences: %v", err)
	}
	want := initial.Normalized()
	want.SourceVolumesSet = true
	if got != want {
		t.Fatalf("Load() = %+v, want cached %+v", got, want)
	}
}

func TestPreferencesLegacyStateNormalizesForPriorityListening(t *testing.T) {
	path := t.TempDir() + "/preferences.json"
	if err := os.WriteFile(path, []byte(`{"last_stream_id":"masjid-1","autoplay":true}`), 0600); err != nil {
		t.Fatalf("write legacy preferences: %v", err)
	}

	state, err := NewPreferences(path).Load()
	if err != nil {
		t.Fatalf("load legacy preferences: %v", err)
	}
	if state.SelectedMasjidID != "masjid-1" {
		t.Fatalf("SelectedMasjidID = %q, want masjid-1", state.SelectedMasjidID)
	}
	if !state.ResumeListening {
		t.Fatal("ResumeListening = false, want true")
	}
	if state.MasjidVolume != DefaultMasjidVolume {
		t.Fatalf("MasjidVolume = %d, want %d", state.MasjidVolume, DefaultMasjidVolume)
	}
	if state.RadioVolume != DefaultRadioVolume {
		t.Fatalf("RadioVolume = %d, want %d", state.RadioVolume, DefaultRadioVolume)
	}
	if state.RadioResumeDelayMinutes != DefaultRadioResumeDelay {
		t.Fatalf("RadioResumeDelayMinutes = %d, want %d", state.RadioResumeDelayMinutes, DefaultRadioResumeDelay)
	}
}
