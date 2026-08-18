package playback

import "errors"

// SetVolumeTransient changes the live player volume without writing persistent
// state. Call PersistVolume when the user has finished making adjustments.
func (m *Manager) SetVolumeTransient(volume int) error {
	if volume < 0 || volume > 100 {
		return errors.New("volume must be between 0 and 100")
	}

	status, err := m.player.Status()
	if err != nil {
		return err
	}
	if err := m.player.Volume(volume); err != nil {
		return err
	}

	m.mu.Lock()
	m.volume = volume
	m.volumeSet = true
	m.volumeDevice = status.AudioDevice
	m.volumeSupported = status.VolumeSupported
	m.status.Volume = volume
	m.status.VolumeSupported = status.VolumeSupported
	m.status.AudioDevice = status.AudioDevice
	m.mu.Unlock()
	return nil
}

// PersistVolume writes the currently active volume once. It is intentionally
// separate from SetVolumeTransient so interactive controls do not write to
// flash storage for every slider movement.
func (m *Manager) PersistVolume() error {
	m.mu.Lock()
	volume := m.volume
	device := m.volumeDevice
	volumeSet := m.volumeSet
	supported := m.volumeSupported
	volumeStore := m.volumeStore
	m.mu.Unlock()

	if !volumeSet || !supported || device == "" || volumeStore == nil {
		return nil
	}
	return volumeStore.Save(device, volume)
}
