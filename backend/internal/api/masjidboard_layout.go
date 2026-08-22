package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

type masjidBoardLayoutResponse struct {
	Layout string `json:"layout"`
}

type masjidBoardLayoutRequest struct {
	Layout string `json:"layout"`
}

func (s *Server) masjidBoardLayout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		layout := selection.LayoutStandard
		if s.masjidBoardService != nil {
			layout = s.masjidBoardService.Selection().EffectiveLayout()
		}
		writeJSON(w, http.StatusOK, masjidBoardLayoutResponse{Layout: layout})
	case http.MethodPut:
		if s.masjidBoardService == nil || s.masjidBoardSelectionManager == nil {
			writeError(w, http.StatusServiceUnavailable, "MasjidBoard selection service is unavailable")
			return
		}

		var request masjidBoardLayoutRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&request); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		layout := strings.TrimSpace(strings.ToLower(request.Layout))
		if layout != selection.LayoutStandard && layout != selection.LayoutDetailed {
			writeError(w, http.StatusBadRequest, "layout must be standard or detailed")
			return
		}

		state := s.masjidBoardService.Selection()
		if !state.Configured() {
			writeError(w, http.StatusConflict, "select at least one MasjidBoard before choosing a display layout")
			return
		}
		state.Layout = layout
		if err := s.masjidBoardSelectionManager.Reconfigure(state); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, masjidBoardLayoutResponse{Layout: state.EffectiveLayout()})
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
