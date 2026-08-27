package playback

import (
	"errors"

	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
)

type softwareVolumePlayer interface {
	SoftwareVolume(volume int) error
}

// ListenOutput adapts the existing endpoint/retry playback manager to the
// higher-level Listen priority controller. The controller decides whether a
// masjid or radio source should be active; Manager continues to own endpoint
// fallback, retries and MPV lifecycle handling.
type ListenOutput struct {
	manager *Manager
}

func NewListenOutput(manager *Manager) *ListenOutput {
	return &ListenOutput{manager: manager}
}

func (o *ListenOutput) Activate(selected stream.Stream, softwareVolume int) error {
	if o == nil || o.manager == nil {
		return errors.New("playback manager is not configured")
	}
	if err := o.SetSoftwareVolume(softwareVolume); err != nil {
		return err
	}
	o.manager.Play(selected)
	return nil
}

func (o *ListenOutput) SetSoftwareVolume(volume int) error {
	if o == nil || o.manager == nil {
		return errors.New("playback manager is not configured")
	}
	if volume < 0 || volume > 100 {
		return errors.New("volume must be between 0 and 100")
	}
	player, ok := o.manager.player.(softwareVolumePlayer)
	if !ok {
		return errors.New("player does not support software volume")
	}
	return player.SoftwareVolume(volume)
}

func (o *ListenOutput) Stop() error {
	if o == nil || o.manager == nil {
		return errors.New("playback manager is not configured")
	}
	o.manager.Stop()
	return nil
}
