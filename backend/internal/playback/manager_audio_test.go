package playback

import (
	"testing"

	"github.com/X-Calibre/MasjidPi/backend/internal/player"
)

type fakeAudioDeviceProvider struct {
	devices []player.AudioDevice
}

func (f fakeAudioDeviceProvider) AudioDevices() ([]player.AudioDevice, error) { return f.devices, nil }

func (f *fakePlayer) AudioDevices() ([]player.AudioDevice, error) {
	return []player.AudioDevice{{Name: "auto", Description: "Autoselect device"}}, nil
}

func (f *fakePlayer) AudioDevice(name string) error {
	return nil
}

func TestManagerUsesIndependentAudioDeviceProvider(t *testing.T) {
	manager := New(&fakePlayer{}, Config{})
	want := []player.AudioDevice{{Name: "alsa/plughw:CARD=USB,DEV=0", Description: "USB Audio"}}
	manager.SetAudioDeviceProvider(fakeAudioDeviceProvider{devices: want})

	got, err := manager.AudioDevices()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != want[0] {
		t.Fatalf("AudioDevices() = %#v, want %#v", got, want)
	}
}
