package api

import (
	"net/http"

	"github.com/X-Calibre/MasjidPi/backend/internal/components"
)

type installedComponents = components.Installed

func currentInstalledComponents() installedComponents {
	return components.Current()
}

func (s *Server) components(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, components.Current())
}
