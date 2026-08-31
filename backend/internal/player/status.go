package player

type Status struct {
	Version         string        `json:"version"`
	State           string        `json:"state"`
	URL             string        `json:"url"`
	Volume          int           `json:"volume"`
	VolumeSupported bool          `json:"volume_supported"`
	Paused          bool          `json:"paused"`
	AudioDevice     string        `json:"audio_device"`
	AudioDevices    []AudioDevice `json:"audio_devices,omitempty"`
}

type AudioDevice struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Unavailable bool   `json:"unavailable,omitempty"`
}
