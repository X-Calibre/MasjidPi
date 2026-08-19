package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/hierarchy"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/scope"
)

type masjidBoardScopeResponse struct {
	Configured bool             `json:"configured"`
	Locations  []scope.Location `json:"locations"`
}

func (s *Server) masjidBoardHierarchy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if strings.TrimSpace(s.masjidBoardHierarchyPath) == "" {
		writeError(w, http.StatusServiceUnavailable, "MasjidBoard hierarchy is not configured")
		return
	}
	state, err := hierarchy.NewStore(s.masjidBoardHierarchyPath).Load()
	if err != nil {
		s.logger.Warn("Could not load MasjidBoard hierarchy", "error", err)
		writeError(w, http.StatusInternalServerError, "could not load MasjidBoard hierarchy")
		return
	}
	writeJSON(w, http.StatusOK, state)
}

func (s *Server) masjidBoardScope(w http.ResponseWriter, r *http.Request) {
	if strings.TrimSpace(s.masjidBoardScopePath) == "" {
		writeError(w, http.StatusServiceUnavailable, "MasjidBoard scope is not configured")
		return
	}
	store := scope.NewStore(s.masjidBoardScopePath)

	switch r.Method {
	case http.MethodGet:
		state, err := store.Load()
		if err != nil {
			s.logger.Warn("Could not load MasjidBoard scope", "error", err)
			writeError(w, http.StatusInternalServerError, "could not load MasjidBoard scope")
			return
		}
		writeJSON(w, http.StatusOK, masjidBoardScopeResponse{Configured: state.Configured(), Locations: state.Locations})

	case http.MethodPut:
		var state scope.State
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&state); err != nil {
			writeError(w, http.StatusBadRequest, "invalid MasjidBoard scope")
			return
		}
		state = state.Normalized()
		if err := state.Validate(); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := store.Save(state); err != nil {
			s.logger.Warn("Could not save MasjidBoard scope", "error", err)
			writeError(w, http.StatusInternalServerError, "could not save MasjidBoard scope")
			return
		}
		writeJSON(w, http.StatusOK, masjidBoardScopeResponse{Configured: true, Locations: state.Locations})

	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
