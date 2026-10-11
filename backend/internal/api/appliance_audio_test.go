package api

import (
	"os"
	"path/filepath"
	"testing"
)

func TestApplianceAudioMode(t *testing.T) {
	root := t.TempDir()
	release := filepath.Join(root, "release.json")
	if applianceAudioMode(root, release) {
		t.Fatal("standard installation detected as appliance")
	}
	if err := os.WriteFile(release, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}
	if !applianceAudioMode(root, release) {
		t.Fatal("appliance image without attached display was not detected")
	}
	connector := filepath.Join(root, "card1-DSI-1")
	if err := os.Mkdir(connector, 0700); err != nil {
		t.Fatal(err)
	}
	for name, value := range map[string]string{"status": "connected", "modes": "720x1280"} {
		if err := os.WriteFile(filepath.Join(connector, name), []byte(value), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if !applianceAudioMode(root, filepath.Join(root, "missing")) {
		t.Fatal("Touch Display 2 appliance was not detected")
	}
}
