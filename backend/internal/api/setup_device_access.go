package api

import "net/http"

func (s *Server) deviceAccess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if s.wifi == nil {
		writeError(w, http.StatusServiceUnavailable, "device network details are unavailable")
		return
	}
	access, err := s.wifi.DeviceAccess(r.Context())
	if err != nil {
		s.logger.Warn("Could not read device network access details", "error", err)
		writeError(w, http.StatusServiceUnavailable, "device network details are unavailable")
		return
	}
	writeJSON(w, http.StatusOK, access)
}
