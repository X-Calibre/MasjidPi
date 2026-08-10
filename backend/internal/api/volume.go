package api

import (
	"encoding/json"
	"net/http"
)

type VolumeRequest struct {
	Volume      int    `json:"volume"`
	AudioDevice string `json:"audio_device,omitempty"`
}

func (s *Server) volume(w http.ResponseWriter, r *http.Request) {
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
	}

	if err := s.playback.Volume(req.Volume); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, s.playback.Status())
}
