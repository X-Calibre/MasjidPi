package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
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
		writeMasjidBoardDisplay(w, r, display.Build(false, selection.State{}, nil))
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
	if provider, ok := s.masjidBoardService.(masjidBoardDailyContentProvider); ok {
		view.DailyIslamicContent = display.PresentDailyIslamicContent(provider.DailyIslamicContent(), s.masjidBoardService.Selection())
	}
	writeMasjidBoardDisplay(w, r, view)
}

func writeMasjidBoardDisplay(w http.ResponseWriter, r *http.Request, view display.View) {
	body, err := json.Marshal(view)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "encode MasjidBoard display")
		return
	}

	digest := sha256.Sum256(body)
	etag := fmt.Sprintf(`"%x"`, digest)
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("ETag", etag)
	if r.Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
