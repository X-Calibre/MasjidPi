package api

import (
	"net/http"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
	masjidboardruntime "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
)

type masjidBoardStatusResponse struct {
	Configured bool                    `json:"configured"`
	Boards     []masjidBoardStatusItem `json:"boards"`
}

type masjidBoardStatusItem struct {
	CatalogueID          string                   `json:"catalogue_id"`
	Provider             string                   `json:"provider"`
	ExternalID           string                   `json:"external_id"`
	Name                 string                   `json:"name"`
	TimeZoneOffsetMS     int64                    `json:"time_zone_offset_ms"`
	Status               masjidboardruntime.Status `json:"status"`
	UsingCachedData      bool                     `json:"using_cached_data"`
	UpdateFailed         bool                     `json:"update_failed"`
	LastAttempt          *time.Time               `json:"last_attempt,omitempty"`
	LastSuccessfulUpdate *time.Time               `json:"last_successful_update,omitempty"`
	UpdateError          string                   `json:"update_error,omitempty"`
	PersistenceError     string                   `json:"persistence_error,omitempty"`
	Board                *model.Board             `json:"board,omitempty"`
}

func (s *Server) masjidBoardStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if s.masjidBoardService == nil {
		writeJSON(w, http.StatusOK, masjidBoardStatusResponse{
			Configured: false,
			Boards:     []masjidBoardStatusItem{},
		})
		return
	}

	configured := s.masjidBoardService.Configured()
	selection := s.masjidBoardService.Selection()
	results := s.masjidBoardService.Results()

	byCatalogueID := make(map[string]masjidboardruntime.Result, len(results))
	for _, result := range results {
		byCatalogueID[result.Selection.CatalogueID] = result
	}

	boards := make([]masjidBoardStatusItem, 0, len(selection.Boards))
	for _, selected := range selection.Boards {
		item := masjidBoardStatusItem{
			CatalogueID:      selected.CatalogueID,
			Provider:         selected.Provider,
			ExternalID:       selected.ExternalID,
			Name:             selected.Name,
			TimeZoneOffsetMS: selected.TimeZoneOffsetMS,
			Status:           masjidboardruntime.StatusUnavailable,
		}

		if result, ok := byCatalogueID[selected.CatalogueID]; ok {
			item.Status = result.Status
			item.UsingCachedData = result.Status == masjidboardruntime.StatusStale
			item.UpdateFailed = result.UpdateError != nil
			item.Board = result.Board
			if !result.LastAttempt.IsZero() {
				attempt := result.LastAttempt
				item.LastAttempt = &attempt
			}
			if !result.LastSuccessfulUpdate.IsZero() {
				success := result.LastSuccessfulUpdate
				item.LastSuccessfulUpdate = &success
			}
			if result.UpdateError != nil {
				item.UpdateError = result.UpdateError.Error()
			}
			if result.PersistenceError != nil {
				item.PersistenceError = result.PersistenceError.Error()
			}
		}

		boards = append(boards, item)
	}

	writeJSON(w, http.StatusOK, masjidBoardStatusResponse{
		Configured: configured,
		Boards:     boards,
	})
}
