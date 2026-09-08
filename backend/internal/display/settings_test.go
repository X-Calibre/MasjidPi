package display

import (
	"os"
	"path/filepath"
	"testing"
)

func TestControllerUpdatesAndRestoresSettings(t *testing.T) {
	root := t.TempDir()
	device := filepath.Join(root, "rpi_backlight")
	if err := os.Mkdir(device, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(device, "max_brightness"), []byte("31\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(device, "brightness"), []byte("31\n"), 0644); err != nil {
		t.Fatal(err)
	}
	statePath := filepath.Join(t.TempDir(), "display.json")
	c := NewController(statePath, root)
	brightness, temperature := 50, "medium"
	settings, err := c.Update(&brightness, &temperature)
	if err != nil {
		t.Fatal(err)
	}
	if !settings.BrightnessAvailable || settings.BrightnessPercent != 50 || settings.ColorTemperature != "medium" {
		t.Fatalf("unexpected settings: %+v", settings)
	}
	data, _ := os.ReadFile(filepath.Join(device, "brightness"))
	if string(data) != "16\n" {
		t.Fatalf("brightness = %q, want 16", data)
	}
	if err := os.WriteFile(filepath.Join(device, "brightness"), []byte("31\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := c.Restore(); err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(filepath.Join(device, "brightness"))
	if string(data) != "16\n" {
		t.Fatalf("restored brightness = %q, want 16", data)
	}
}

func TestControllerRejectsInvalidValues(t *testing.T) {
	c := NewController(filepath.Join(t.TempDir(), "display.json"), t.TempDir())
	brightness := 9
	if _, err := c.Update(&brightness, nil); err == nil {
		t.Fatal("expected brightness validation error")
	}
	temperature := "orange"
	if _, err := c.Update(nil, &temperature); err == nil {
		t.Fatal("expected temperature validation error")
	}
}
