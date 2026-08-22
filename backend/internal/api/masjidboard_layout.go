package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

type masjidBoardLayoutResponse struct {
	Layout string `json:"layout"`
	Theme  string `json:"theme"`
}

type masjidBoardLayoutRequest struct {
	Layout string `json:"layout"`
	Theme  string `json:"theme"`
}

type masjidBoardLayoutSetter interface { SetLayout(string) error }
type masjidBoardThemeSetter interface { SetTheme(string) error }

func (s *Server) masjidBoardLayout(w http.ResponseWriter, r *http.Request) {
	current := func() masjidBoardLayoutResponse {
		response := masjidBoardLayoutResponse{Layout: selection.LayoutStandard, Theme: selection.ThemeEmerald}
		if s.masjidBoardService != nil {
			state := s.masjidBoardService.Selection()
			response.Layout = state.EffectiveLayout()
			response.Theme = state.EffectiveTheme()
		}
		return response
	}

	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, current())
	case http.MethodPut:
		var request masjidBoardLayoutRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&request); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		layout := strings.TrimSpace(strings.ToLower(request.Layout))
		theme := strings.TrimSpace(strings.ToLower(request.Theme))
		if layout == "" && theme == "" {
			writeError(w, http.StatusBadRequest, "layout or theme is required")
			return
		}
		if layout != "" {
			if layout != selection.LayoutStandard && layout != selection.LayoutDetailed {
				writeError(w, http.StatusBadRequest, "layout must be standard or detailed")
				return
			}
			setter, ok := s.masjidBoardService.(masjidBoardLayoutSetter)
			if !ok || setter == nil {
				writeError(w, http.StatusServiceUnavailable, "MasjidBoard layout service is unavailable")
				return
			}
			if err := setter.SetLayout(layout); err != nil {
				writeError(w, http.StatusConflict, err.Error())
				return
			}
		}
		if theme != "" {
			if !selection.ThemeSupported(theme) {
				writeError(w, http.StatusBadRequest, "unsupported MasjidBoard theme")
				return
			}
			setter, ok := s.masjidBoardService.(masjidBoardThemeSetter)
			if !ok || setter == nil {
				writeError(w, http.StatusServiceUnavailable, "MasjidBoard theme service is unavailable")
				return
			}
			if err := setter.SetTheme(theme); err != nil {
				writeError(w, http.StatusConflict, err.Error())
				return
			}
		}
		writeJSON(w, http.StatusOK, current())
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
