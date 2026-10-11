package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/X-Calibre/MasjidFrame/backend/internal/player"
	"github.com/X-Calibre/MasjidFrame/backend/internal/updates"
)

// Use the same connected 720x1280 DSI check as the local display launcher.
func detectedDisplayProfile(drmPath string) string {
	connectors, err := filepath.Glob(filepath.Join(drmPath, "card*-DSI-*"))
	if err != nil {
		return "standard"
	}
	for _, connector := range connectors {
		status, statusErr := os.ReadFile(filepath.Join(connector, "status"))
		modes, modesErr := os.ReadFile(filepath.Join(connector, "modes"))
		if statusErr != nil || modesErr != nil || strings.TrimSpace(string(status)) != "connected" {
			continue
		}
		for _, mode := range strings.Fields(string(modes)) {
			if mode == "720x1280" {
				return "appliance-720"
			}
		}
	}
	return "standard"
}

func applianceAudioMode(drmPath, releasePath string) bool {
	if _, err := os.Stat(releasePath); err == nil {
		return true
	}
	return detectedDisplayProfile(drmPath) == "appliance-720"
}

func restrictApplianceAudioDevices(devices []player.AudioDevice) []player.AudioDevice {
	if applianceAudioMode("/sys/class/drm", updates.DefaultReleaseRecordPath) {
		return player.ApplianceAudioDevices(devices)
	}
	return devices
}

func (s *Server) displayProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"profile": detectedDisplayProfile("/sys/class/drm")})
}
