package updates

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestStoreLoadMissingReturnsDefaultState(t *testing.T) {
	path := filepath.Join(
		t.TempDir(),
		"nested",
		"update_state.json",
	)

	state, err := NewStore(path).Load()
	if err != nil {
		t.Fatal(err)
	}
	if state.SchemaVersion != StateSchemaVersion {
		t.Fatalf(
			"schema version = %d, want %d",
			state.SchemaVersion,
			StateSchemaVersion,
		)
	}
	if state.AvailableRelease != nil {
		t.Fatalf(
			"available release = %+v, want nil",
			state.AvailableRelease,
		)
	}
}

func TestStoreRoundTripAndPermissions(t *testing.T) {
	path := filepath.Join(
		t.TempDir(),
		"update_state.json",
	)
	store := NewStore(path)

	want := validState()
	checked := time.Date(
		2026,
		time.September,
		16,
		8,
		5,
		0,
		0,
		time.UTC,
	)
	want.LastCheckedAt = &checked

	if err := store.Save(want); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0600 {
		t.Fatalf(
			"permissions = %o, want 600",
			got,
		)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf(
			"loaded state = %#v, want %#v",
			got,
			want,
		)
	}
}

func TestStoreRejectedStateDoesNotCreateFile(t *testing.T) {
	path := filepath.Join(
		t.TempDir(),
		"update_state.json",
	)
	state := DefaultState("v1.6.0")
	state.AvailableRelease = &Release{
		Version: "v1.7.0",
	}

	err := NewStore(path).Save(state)
	if err == nil {
		t.Fatal("Save() error = nil, want validation error")
	}
	if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
		t.Fatalf(
			"state file exists after rejected save: %v",
			statErr,
		)
	}
}

func TestStoreRejectsUnsupportedPersistedSchema(t *testing.T) {
	path := filepath.Join(
		t.TempDir(),
		"update_state.json",
	)

	data, err := json.Marshal(State{
		SchemaVersion: 99,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}

	_, err = NewStore(path).Load()
	if err == nil {
		t.Fatal("Load() error = nil, want schema error")
	}
	if !strings.Contains(
		err.Error(),
		"unsupported state schema",
	) {
		t.Fatalf("Load() error = %q", err)
	}
}

func TestStoreRejectsMalformedJSON(t *testing.T) {
	path := filepath.Join(
		t.TempDir(),
		"update_state.json",
	)

	if err := os.WriteFile(
		path,
		[]byte(`{"schema_version":`),
		0600,
	); err != nil {
		t.Fatal(err)
	}

	_, err := NewStore(path).Load()
	if err == nil {
		t.Fatal("Load() error = nil, want decode error")
	}
	if !strings.Contains(err.Error(), "decode state") {
		t.Fatalf("Load() error = %q", err)
	}
}
