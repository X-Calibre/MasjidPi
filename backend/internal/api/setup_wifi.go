package api

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

type wifiConnectRequest struct {
	SSID     string `json:"ssid"`
	Password string `json:"password"`
}

func (s *Server) wifiStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if s.wifi == nil {
		writeJSON(w, http.StatusOK, map[string]bool{"supported": false})
		return
	}
	status, err := s.wifi.Status(r.Context())
	if err != nil {
		s.logger.Warn("Could not read Wi-Fi status", "error", err)
		writeError(w, http.StatusServiceUnavailable, "Wi-Fi status is temporarily unavailable")
		return
	}
	writeJSON(w, http.StatusOK, status)
}

func (s *Server) wifiNetworks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !isLoopbackRequest(r) {
		writeError(w, http.StatusForbidden, "Wi-Fi setup is available only on the device")
		return
	}
	if s.wifi == nil {
		writeError(w, http.StatusServiceUnavailable, "Wi-Fi setup is unavailable")
		return
	}
	networks, err := s.wifi.Scan(r.Context())
	if err != nil {
		s.logger.Warn("Could not scan Wi-Fi networks", "error", err)
		writeError(w, http.StatusServiceUnavailable, "Could not scan for Wi-Fi networks")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"networks": networks})
}

func (s *Server) wifiConnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !isLoopbackRequest(r) {
		writeError(w, http.StatusForbidden, "Wi-Fi setup is available only on the device")
		return
	}
	if s.wifi == nil {
		writeError(w, http.StatusServiceUnavailable, "Wi-Fi setup is unavailable")
		return
	}

	var request wifiConnectRequest
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid Wi-Fi connection request")
		return
	}
	request.SSID = strings.TrimSpace(request.SSID)
	if request.SSID == "" || len(request.SSID) > 32 || len(request.Password) > 64 {
		writeError(w, http.StatusBadRequest, "invalid Wi-Fi network name or password")
		return
	}
	if err := s.wifi.Connect(r.Context(), request.SSID, request.Password); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"connected": true})
}

func isLoopbackRequest(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
