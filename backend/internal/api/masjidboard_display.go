package api

import (
	"net/http"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/display"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

// masjidBoardDisplay exposes only presentation-oriented MasjidBoard state. It
// intentionally omits provider discovery details, configuration controls and
// diagnostic error strings that belong to the administrative/status APIs.
func (s *Server) masjidBoardDisplay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if s.masjidBoardService == nil {
		writeJSON(w, http.StatusOK, display.Build(false, selection.State{}, nil))
		return
	}

	view := display.Build(
		s.masjidBoardService.Configured(),
		s.masjidBoardService.Selection(),
		s.masjidBoardService.Results(),
	)
	if provider, ok := s.masjidBoardService.(masjidBoardEconomicProvider); ok {
		view.EconomicIndicators = provider.EconomicIndicators()
	}
	writeJSON(w, http.StatusOK, view)
}
