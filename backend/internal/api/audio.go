package api

import (
	"encoding/json"
	"net/http"
)

type AudioDeviceRequest struct {
	Name string `json:"name"`
}

func (s *Server) audio(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		devices, err := s.playback.AudioDevices()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, devices)

	case http.MethodPost:
		var req AudioDeviceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON")
			return
		}

		if err := s.playback.AudioDevice(req.Name); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, s.playback.Status())

	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
