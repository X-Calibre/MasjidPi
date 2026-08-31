package api

import (
	"encoding/json"
	"net/http"

	"github.com/X-Calibre/MasjidPi/backend/internal/config"
	"github.com/X-Calibre/MasjidPi/backend/internal/player"
	"github.com/X-Calibre/MasjidPi/backend/internal/storage"
)

type VolumeRequest struct {
	Volume      *int   `json:"volume,omitempty"`
	AudioDevice string `json:"audio_device,omitempty"`
	Persist     *bool  `json:"persist,omitempty"`
}

func (s *Server) volume(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Query().Get("devices") == "1" {
		devices, err := s.playback.AudioDevices()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		state := s.audioDeviceState
		if state == nil {
			if paths, pathErr := config.RuntimePaths(); pathErr == nil {
				state = storage.NewAudioDeviceState(paths.AudioDeviceState)
			}
		}
		if state != nil {
			if saved, ok, loadErr := state.Load(); loadErr == nil && ok && !containsAudioDevice(devices, saved) {
				devices = append(devices, player.AudioDevice{
					Name:        saved,
					Description: player.AudioDeviceDescription(saved),
					Unavailable: true,
				})
			}
		}
		writeJSON(w, http.StatusOK, devices)
		return
	}

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req VolumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if req.AudioDevice != "" {
		if err := s.playback.AudioDevice(req.AudioDevice); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		state := s.audioDeviceState
		if state == nil {
			paths, err := config.RuntimePaths()
			if err != nil {
				writeError(w, http.StatusInternalServerError, err.Error())
				return
			}
			state = storage.NewAudioDeviceState(paths.AudioDeviceState)
		}
		if err := state.Save(req.AudioDevice); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if req.Volume != nil {
		persist := true
		if req.Persist != nil {
			persist = *req.Persist
		}

		if persist {
			if err := s.playback.Volume(*req.Volume); err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
		} else if err := s.playback.SetVolumeTransient(*req.Volume); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	writeJSON(w, http.StatusOK, s.playback.Status())
}

func containsAudioDevice(devices []player.AudioDevice, name string) bool {
	for _, device := range devices {
		if device.Name == name {
			return true
		}
	}
	return false
}
