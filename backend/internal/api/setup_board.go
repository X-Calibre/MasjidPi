package api

import (
	"encoding/json"
	"net/http"

	"github.com/X-Calibre/MasjidPi/backend/internal/storage"
)

type boardSetupRequest struct {
	Deferred bool `json:"deferred"`
}

func (s *Server) boardSetup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		state, err := s.preferences.Load()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, boardSetupRequest{Deferred: state.BoardSetupDeferred})
	case http.MethodPut:
		var request boardSetupRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		state, err := s.preferences.Update(func(state *storage.PreferencesState) {
			state.BoardSetupDeferred = request.Deferred
		})
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, boardSetupRequest{Deferred: state.BoardSetupDeferred})
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
