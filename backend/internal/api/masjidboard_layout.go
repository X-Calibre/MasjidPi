package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

type masjidBoardLayoutResponse struct {
	Layout                 string `json:"layout"`
	Theme                  string `json:"theme"`
	SlideDurationSeconds   int    `json:"slide_duration_seconds"`
	ShowEconomicIndicators bool   `json:"show_economic_indicators"`
}

type masjidBoardLayoutRequest struct {
	Layout                 string `json:"layout"`
	Theme                  string `json:"theme"`
	SlideDurationSeconds   int    `json:"slide_duration_seconds"`
	ShowEconomicIndicators *bool  `json:"show_economic_indicators"`
}

type masjidBoardLayoutSetter interface{ SetLayout(string) error }
type masjidBoardThemeSetter interface{ SetTheme(string) error }
type masjidBoardSlideDurationSetter interface{ SetSlideDurationSeconds(int) error }
type masjidBoardEconomicSetter interface{ SetShowEconomicIndicators(bool) error }
type masjidBoardEconomicRefresher interface{ RefreshEconomicIndicators(context.Context) error }

func (s *Server) masjidBoardLayout(w http.ResponseWriter, r *http.Request) {
	current := func() masjidBoardLayoutResponse {
		response := masjidBoardLayoutResponse{Layout: selection.LayoutLandscape, Theme: selection.ThemeEmerald, SlideDurationSeconds: selection.DefaultSlideDurationSeconds}
		if s.masjidBoardService != nil {
			state := s.masjidBoardService.Selection()
			response.Layout = state.EffectiveLayout()
			response.Theme = state.EffectiveTheme()
			response.SlideDurationSeconds = state.EffectiveSlideDurationSeconds()
			response.ShowEconomicIndicators = state.ShowEconomicIndicators
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
		if layout == "" && theme == "" && request.SlideDurationSeconds == 0 && request.ShowEconomicIndicators == nil {
			writeError(w, http.StatusBadRequest, "layout, theme or slide duration is required")
			return
		}
		if layout != "" && layout != selection.LayoutLandscape && layout != selection.LayoutPortrait {
			writeError(w, http.StatusBadRequest, "layout must be landscape or portrait")
			return
		}
		if theme != "" && !selection.ThemeSupported(theme) {
			writeError(w, http.StatusBadRequest, "unsupported MasjidBoard theme")
			return
		}
		if request.SlideDurationSeconds != 0 &&
			(request.SlideDurationSeconds < selection.MinSlideDurationSeconds || request.SlideDurationSeconds > selection.MaxSlideDurationSeconds) {
			writeError(w, http.StatusBadRequest, "slide duration must be between 5 and 60 seconds")
			return
		}

		if layout != "" {
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
		if request.SlideDurationSeconds != 0 {
			setter, ok := s.masjidBoardService.(masjidBoardSlideDurationSetter)
			if !ok || setter == nil {
				writeError(w, http.StatusServiceUnavailable, "MasjidBoard slide duration service is unavailable")
				return
			}
			if err := setter.SetSlideDurationSeconds(request.SlideDurationSeconds); err != nil {
				writeError(w, http.StatusConflict, err.Error())
				return
			}
		}
		if request.ShowEconomicIndicators != nil {
			setter, ok := s.masjidBoardService.(masjidBoardEconomicSetter)
			if !ok || setter == nil {
				writeError(w, http.StatusServiceUnavailable, "MasjidBoard economic indicators service is unavailable")
				return
			}
			if err := setter.SetShowEconomicIndicators(*request.ShowEconomicIndicators); err != nil {
				writeError(w, http.StatusConflict, err.Error())
				return
			}
			if *request.ShowEconomicIndicators {
				if refresher, ok := s.masjidBoardService.(masjidBoardEconomicRefresher); ok {
					go func() {
						ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
						defer cancel()
						_ = refresher.RefreshEconomicIndicators(ctx)
					}()
				}
			}
		}
		writeJSON(w, http.StatusOK, current())
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
