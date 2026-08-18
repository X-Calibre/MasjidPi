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
	if got != initial {
		t.Fatalf("Load() = %+v, want cached %+v", got, initial)
	}
}
