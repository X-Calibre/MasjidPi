package api

import (
	"net/http"
	"os"
	"strings"
)

type installedComponents struct {
	Listen bool `json:"listen"`
	Board  bool `json:"board"`
}

func currentInstalledComponents() installedComponents {
	// Backward compatibility: installations created before component profiles
	// existed expose both capabilities until the installer writes an explicit
	// MASJIDPI_COMPONENTS value.
	value := strings.TrimSpace(os.Getenv("MASJIDPI_COMPONENTS"))
	if value == "" {
		return installedComponents{Listen: true, Board: true}
	}

	components := installedComponents{}
	for _, item := range strings.Split(value, ",") {
		switch strings.TrimSpace(strings.ToLower(item)) {
		case "listen":
			components.Listen = true
		case "board":
			components.Board = true
		case "both":
			components.Listen = true
			components.Board = true
		}
	}
	return components
}

func (s *Server) components(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, currentInstalledComponents())
}
