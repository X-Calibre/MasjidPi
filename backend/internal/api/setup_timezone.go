package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

type timezoneRequest struct {
	Name string `json:"name"`
}

func (s *Server) timezones(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !isLoopbackRequest(r) {
		writeError(w, http.StatusForbidden, "timezone setup is available only on the device")
		return
	}
	if s.timezoneController == nil {
		writeError(w, http.StatusServiceUnavailable, "timezone setup is unavailable")
		return
	}

	country := strings.TrimSpace(r.URL.Query().Get("country"))
	if country == "" {
		writeError(w, http.StatusBadRequest, "country is required")
		return
	}
	zones, err := s.timezoneController.Zones(country)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	current, configured, err := s.timezoneController.Current()
	if err != nil {
		s.logger.Warn("Could not read saved timezone", "error", err)
		current = ""
		configured = false
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"country":    country,
		"zones":      zones,
		"current":    current,
		"configured": configured,
	})
}

func (s *Server) timezone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !isLoopbackRequest(r) {
		writeError(w, http.StatusForbidden, "timezone setup is available only on the device")
		return
	}
	if s.timezoneController == nil {
		writeError(w, http.StatusServiceUnavailable, "timezone setup is unavailable")
		return
	}

	var request timezoneRequest
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid timezone request")
		return
	}
	request.Name = strings.TrimSpace(request.Name)
	if request.Name == "" || len(request.Name) > 128 {
		writeError(w, http.StatusBadRequest, "invalid timezone")
		return
	}
	if err := s.timezoneController.Set(r.Context(), request.Name); err != nil {
		s.logger.Warn("Could not set timezone", "timezone", request.Name, "error", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"name":    request.Name,
		"updated": true,
	})
}
