package atomicfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteJSONCreatesAndReplacesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "state.json")
	if err := WriteJSON(path, map[string]string{"value": "one"}, 0600); err != nil {
		t.Fatal(err)
	}
	if err := WriteJSON(path, map[string]string{"value": "two"}, 0600); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"value":"two"}` {
		t.Fatalf("data = %s", data)
	}
}

func TestRemoveDeletesFileAndIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	if err := Write(path, []byte("state"), 0600); err != nil {
		t.Fatal(err)
	}
	removed, err := Remove(path)
	if err != nil {
		t.Fatal(err)
	}
	if !removed {
		t.Fatal("Remove() reported that the existing file was not removed")
	}
	removed, err = Remove(path)
	if err != nil {
		t.Fatal(err)
	}
	if removed {
		t.Fatal("second Remove() reported a removal")
	}
}
