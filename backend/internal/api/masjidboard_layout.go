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
	Theme                  string `json:"theme"`
	SlideDurationSeconds   int    `json:"slide_duration_seconds"`
	ShowEconomicIndicators bool   `json:"show_economic_indicators"`
	ShowDailyAyah          bool   `json:"show_daily_ayah"`
	ShowDailyHadith        bool   `json:"show_daily_hadith"`
	ShowDailySunnah        bool   `json:"show_daily_sunnah"`
}

type masjidBoardLayoutRequest struct {
	Theme                  string `json:"theme"`
	SlideDurationSeconds   int    `json:"slide_duration_seconds"`
	ShowEconomicIndicators *bool  `json:"show_economic_indicators"`
	ShowDailyAyah          *bool  `json:"show_daily_ayah"`
	ShowDailyHadith        *bool  `json:"show_daily_hadith"`
	ShowDailySunnah        *bool  `json:"show_daily_sunnah"`
}

type masjidBoardThemeSetter interface{ SetTheme(string) error }
type masjidBoardSlideDurationSetter interface{ SetSlideDurationSeconds(int) error }
type masjidBoardEconomicSetter interface{ SetShowEconomicIndicators(bool) error }
type masjidBoardEconomicRefresher interface{ RefreshEconomicIndicators(context.Context) error }
type masjidBoardDailyContentSetter interface {
	SetDailyIslamicContentPreferences(bool, bool, bool) error
}
type masjidBoardDailyContentRefresher interface{ RefreshDailyIslamicContent(context.Context) error }

func (s *Server) masjidBoardLayout(w http.ResponseWriter, r *http.Request) {
	current := func() masjidBoardLayoutResponse {
		response := masjidBoardLayoutResponse{
			Theme: selection.ThemeEmerald, SlideDurationSeconds: selection.DefaultSlideDurationSeconds,
			ShowDailyAyah: true, ShowDailyHadith: true, ShowDailySunnah: true,
		}
		if s.masjidBoardService != nil {
			state := s.masjidBoardService.Selection()
			response.Theme = state.EffectiveTheme()
			response.SlideDurationSeconds = state.EffectiveSlideDurationSeconds()
			response.ShowEconomicIndicators = state.ShowEconomicIndicators
			response.ShowDailyAyah = state.ShowDailyAyahValue()
			response.ShowDailyHadith = state.ShowDailyHadithValue()
			response.ShowDailySunnah = state.ShowDailySunnahValue()
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

		theme := strings.TrimSpace(strings.ToLower(request.Theme))
		if theme == "" && request.SlideDurationSeconds == 0 && request.ShowEconomicIndicators == nil &&
			request.ShowDailyAyah == nil && request.ShowDailyHadith == nil && request.ShowDailySunnah == nil {
			writeError(w, http.StatusBadRequest, "theme, slide duration or additional information setting is required")
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
		if request.ShowDailyAyah != nil || request.ShowDailyHadith != nil || request.ShowDailySunnah != nil {
			currentState := current()
			showAyah, showHadith, showSunnah := currentState.ShowDailyAyah, currentState.ShowDailyHadith, currentState.ShowDailySunnah
			if request.ShowDailyAyah != nil {
				showAyah = *request.ShowDailyAyah
			}
			if request.ShowDailyHadith != nil {
				showHadith = *request.ShowDailyHadith
			}
			if request.ShowDailySunnah != nil {
				showSunnah = *request.ShowDailySunnah
			}
			setter, ok := s.masjidBoardService.(masjidBoardDailyContentSetter)
			if !ok || setter == nil {
				writeError(w, http.StatusServiceUnavailable, "MasjidBoard daily content service is unavailable")
				return
			}
			if err := setter.SetDailyIslamicContentPreferences(showAyah, showHadith, showSunnah); err != nil {
				writeError(w, http.StatusConflict, err.Error())
				return
			}
			if showAyah || showHadith || showSunnah {
				if refresher, ok := s.masjidBoardService.(masjidBoardDailyContentRefresher); ok {
					go func() {
						ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
						defer cancel()
						_ = refresher.RefreshDailyIslamicContent(ctx)
					}()
				}
			}
		}
		writeJSON(w, http.StatusOK, current())
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
