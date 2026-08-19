package api

import (
	"context"
	"net/http"
	"time"
)

type masjidBoardBoardsRefreshResponse struct {
	Refreshed int `json:"refreshed"`
}

func (s *Server) masjidBoardBoardsRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if s.masjidBoardSelectionManager == nil {
		writeError(w, http.StatusServiceUnavailable, "MasjidBoard selection service is unavailable")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	results := s.masjidBoardSelectionManager.Refresh(ctx)

	writeJSON(w, http.StatusOK, masjidBoardBoardsRefreshResponse{Refreshed: len(results)})
}
