package playback

import "github.com/X-Calibre/MasjidPi/backend/internal/player"

func (f *fakePlayer) AudioDevices() ([]player.AudioDevice, error) {
	return []player.AudioDevice{{Name: "auto", Description: "Autoselect device"}}, nil
}

func (f *fakePlayer) AudioDevice(name string) error {
	return nil
}
