package storage

import (
	"os"
	"testing"
)

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
}
