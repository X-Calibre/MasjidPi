package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	masjidboardcatalogue "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/catalogue"
	masjidboardruntime "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

type masjidBoardSelectionManager interface {
	Reconfigure(selection.State) error
	Refresh(context.Context) []masjidboardruntime.Result
}

type masjidBoardSelectionRequest struct {
	CatalogueIDs []string `json:"catalogue_ids"`
}

type masjidBoardSelectionResponse struct {
	Configured bool              `json:"configured"`
	Boards     []selection.Board `json:"boards"`
}

func (s *Server) masjidBoardSelection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.writeMasjidBoardSelection(w)
	case http.MethodPut:
		s.updateMasjidBoardSelection(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) writeMasjidBoardSelection(w http.ResponseWriter) {
	if s.masjidBoardService == nil {
		writeJSON(w, http.StatusOK, masjidBoardSelectionResponse{Boards: []selection.Board{}})
		return
	}
	state := s.masjidBoardService.Selection()
	boards := state.Boards
	if boards == nil {
		boards = []selection.Board{}
	}
	writeJSON(w, http.StatusOK, masjidBoardSelectionResponse{Configured: state.Configured(), Boards: boards})
}

func (s *Server) updateMasjidBoardSelection(w http.ResponseWriter, r *http.Request) {
	if s.masjidBoardSelectionManager == nil {
		writeError(w, http.StatusServiceUnavailable, "MasjidBoard selection service is unavailable")
		return
	}
	if strings.TrimSpace(s.masjidBoardCataloguePath) == "" {
		writeError(w, http.StatusServiceUnavailable, "MasjidBoard catalogue is unavailable")
		return
	}

	var request masjidBoardSelectionRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(request.CatalogueIDs) < selection.MinBoards || len(request.CatalogueIDs) > selection.MaxBoards {
		writeError(w, http.StatusBadRequest, "select between 1 and 3 boards")
		return
	}

	state, err := masjidboardcatalogue.NewStore(s.masjidBoardCataloguePath).Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	merged := masjidboardcatalogue.Merge(state)
	byID := make(map[string]masjidboardcatalogue.Record, len(merged.Records))
	for _, record := range merged.Records {
		byID[record.ID] = record
	}

	selected := selection.State{Boards: make([]selection.Board, 0, len(request.CatalogueIDs))}
	seen := make(map[string]struct{}, len(request.CatalogueIDs))
	for _, rawID := range request.CatalogueIDs {
		id := strings.TrimSpace(rawID)
		if id == "" {
			writeError(w, http.StatusBadRequest, "catalogue ID is required")
			return
		}
		if _, exists := seen[id]; exists {
			writeError(w, http.StatusBadRequest, "duplicate catalogue ID")
			return
		}
		seen[id] = struct{}{}

		record, ok := byID[id]
		if !ok || record.Status != masjidboardcatalogue.StatusActive {
			writeError(w, http.StatusBadRequest, "selected board is not active in the local catalogue: "+id)
			return
		}
		board, err := selection.FromCatalogueRecord(record)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		selected.Boards = append(selected.Boards, board)
	}

	if err := s.masjidBoardSelectionManager.Reconfigure(selected); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// A selection change should produce display data immediately instead of
	// requiring an application restart. Bound the live fetch so a slow provider
	// cannot hold the configuration request indefinitely.
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	s.masjidBoardSelectionManager.Refresh(ctx)

	boards := selected.Boards
	if boards == nil {
		boards = []selection.Board{}
	}
	writeJSON(w, http.StatusOK, masjidBoardSelectionResponse{Configured: true, Boards: boards})
}
