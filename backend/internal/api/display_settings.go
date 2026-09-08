package api

import (
	"encoding/json"
	"net/http"
)

type displaySettingsRequest struct {
	BrightnessPercent *int    `json:"brightness_percent"`
	ColorTemperature  *string `json:"color_temperature"`
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
		if request.BrightnessPercent == nil && request.ColorTemperature == nil {
			writeError(w, http.StatusBadRequest, "brightness or color temperature is required")
			return
		}
		settings, err := s.displaySettings.Update(request.BrightnessPercent, request.ColorTemperature)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, settings)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
