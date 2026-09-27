package api

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectedDisplayProfile(t *testing.T) {
	root := t.TempDir()
	connector := filepath.Join(root, "card1-DSI-1")
	if err := os.MkdirAll(connector, 0755); err != nil {
		t.Fatal(err)
	}
	write := func(file, value string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(connector, file), []byte(value), 0644); err != nil {
			t.Fatal(err)
		}
	}
	write("status", "disconnected\n")
	write("modes", "720x1280\n")
	if got := detectedDisplayProfile(root); got != "standard" {
		t.Fatalf("disconnected DSI profile = %q", got)
	}
	write("status", "connected\n")
	if got := detectedDisplayProfile(root); got != "appliance-720" {
		t.Fatalf("connected Touch Display 2 profile = %q", got)
	}
	write("modes", "1920x1080\n")
	if got := detectedDisplayProfile(root); got != "standard" {
		t.Fatalf("different DSI mode profile = %q", got)
	}
}
