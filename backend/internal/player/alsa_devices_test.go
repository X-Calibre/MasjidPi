package player

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

type fallbackAudioDevices struct {
	devices []AudioDevice
}

func (f fallbackAudioDevices) AudioDevices() ([]AudioDevice, error) { return f.devices, nil }

func TestALSAAudioDevicesReflectHotPlugChanges(t *testing.T) {
	root := t.TempDir()
	createALSATestCard(t, root, 0, "Headphones", "", 0)

	provider := NewALSAAudioDevices(root, nil)
	devices, err := provider.AudioDevices()
	if err != nil {
		t.Fatal(err)
	}
	assertAudioDevice(t, devices, "alsa/plughw:CARD=Headphones,DEV=0", "Headphones")

	createALSATestCard(t, root, 1, "UACDemoV10", "DRIVER=snd-usb-audio\n", 0)
	devices, err = provider.AudioDevices()
	if err != nil {
		t.Fatal(err)
	}
	assertAudioDevice(t, devices, "alsa/plughw:CARD=UACDemoV10,DEV=0", "USB Audio")

	if err := os.Remove(filepath.Join(root, "pcmC1D0p")); err != nil {
		t.Fatal(err)
	}
	devices, err = provider.AudioDevices()
	if err != nil {
		t.Fatal(err)
	}
	for _, device := range devices {
		if device.Name == "alsa/plughw:CARD=UACDemoV10,DEV=0" {
			t.Fatal("disconnected USB device remained in fresh ALSA inventory")
		}
	}
}

func TestALSAAudioDevicesFallsBackWhenSysfsUnavailable(t *testing.T) {
	want := []AudioDevice{{Name: "auto", Description: "Default audio output"}}
	provider := NewALSAAudioDevices(filepath.Join(t.TempDir(), "missing"), fallbackAudioDevices{devices: want})
	got, err := provider.AudioDevices()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != want[0] {
		t.Fatalf("fallback devices = %#v, want %#v", got, want)
	}
}

func TestALSAAudioDevicesDoesNotReturnStaleFallbackForEmptySysfs(t *testing.T) {
	stale := []AudioDevice{{Name: "alsa/plughw:CARD=OLD,DEV=0", Description: "Old device"}}
	provider := NewALSAAudioDevices(t.TempDir(), fallbackAudioDevices{devices: stale})
	got, err := provider.AudioDevices()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != "auto" {
		t.Fatalf("devices = %#v, want only automatic output", got)
	}
}

func createALSATestCard(t *testing.T, root string, card int, id, uevent string, devices ...int) {
	t.Helper()
	cardPath := filepath.Join(root, "card"+strconv.Itoa(card))
	if err := os.MkdirAll(filepath.Join(cardPath, "device"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cardPath, "id"), []byte(id+"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cardPath, "device", "uevent"), []byte(uevent), 0644); err != nil {
		t.Fatal(err)
	}
	for _, device := range devices {
		name := "pcmC" + strconv.Itoa(card) + "D" + strconv.Itoa(device) + "p"
		if err := os.WriteFile(filepath.Join(root, name), nil, 0644); err != nil {
			t.Fatal(err)
		}
	}
}

func assertAudioDevice(t *testing.T, devices []AudioDevice, name, description string) {
	t.Helper()
	for _, device := range devices {
		if device.Name == name && device.Description == description {
			return
		}
	}
	t.Fatalf("device %q (%q) not found in %#v", name, description, devices)
}
