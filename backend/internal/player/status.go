package player

type Status struct {
	Version     string `json:"version"`
	State       string `json:"state"`
	URL         string `json:"url"`
	Volume      int    `json:"volume"`
	Paused      bool   `json:"paused"`
	AudioDevice string `json:"audio_device"`
}

type AudioDevice struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
