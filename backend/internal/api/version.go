package api

import (
	"encoding/json"
	"net/http"

	"github.com/X-Calibre/MasjidPi/backend/internal/version"
)

func (s *Server) version(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(map[string]string{
		"name":    version.AppName,
		"version": version.Version,
	})
}
