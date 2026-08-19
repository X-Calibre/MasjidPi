package api

import (
	"context"
	"net/http"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/catalogue"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/hierarchy"
)

type masjidBoardMaintenance interface {
	RefreshHierarchy(context.Context, time.Time) (hierarchy.RefreshResult, error)
	RefreshCatalogue(context.Context, time.Time) (catalogue.RefreshResult, error)
}

func (s *Server) masjidBoardHierarchyRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if s.masjidBoardMaintenance == nil {
		writeError(w, http.StatusServiceUnavailable, "MasjidBoard maintenance is not available")
		return
	}
	result, err := s.masjidBoardMaintenance.RefreshHierarchy(r.Context(), time.Now())
	if err != nil {
		s.logger.Warn("Manual MasjidBoard hierarchy refresh failed", "error", err)
		writeError(w, http.StatusBadGateway, "MasjidBoard hierarchy refresh failed")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) masjidBoardCatalogueRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if s.masjidBoardMaintenance == nil {
		writeError(w, http.StatusServiceUnavailable, "MasjidBoard maintenance is not available")
		return
	}
	result, err := s.masjidBoardMaintenance.RefreshCatalogue(r.Context(), time.Now())
	if err != nil {
		s.logger.Warn("Manual MasjidBoard catalogue refresh failed", "error", err)
		writeError(w, http.StatusBadGateway, "MasjidBoard catalogue refresh failed")
		return
	}
	writeJSON(w, http.StatusOK, result)
}
