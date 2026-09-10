package api

import (
	"encoding/json"
	"net/http"
)

type displaySettingsRequest struct {
	BrightnessPercent *int `json:"brightness_percent"`
}

func (s *Server) displaySettingsHandler(w http.ResponseWriter, r *http.Request) {
	if s.displaySettings == nil {
		writeError(w, http.StatusServiceUnavailable, "display settings are unavailable")
		return
	}
	switch r.Method {
	case http.MethodGet:
		settings, err := s.displaySettings.Load()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, settings)
	case http.MethodPut:
		var request displaySettingsRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&request); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if request.BrightnessPercent == nil {
			writeError(w, http.StatusBadRequest, "brightness is required")
			return
		}
		settings, err := s.displaySettings.Update(request.BrightnessPercent)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, settings)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
