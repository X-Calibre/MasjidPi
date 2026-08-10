package api

import (
	"encoding/json"
	"net/http"

	"github.com/X-Calibre/MasjidPi/backend/internal/storage"
)

type PreferencesRequest struct {
	LastStreamID string `json:"last_stream_id,omitempty"`
	Autoplay     bool   `json:"autoplay"`
}

func (s *Server) preferencesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		state, err := s.preferences.Load()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, PreferencesRequest{LastStreamID: state.LastStreamID, Autoplay: state.Autoplay})
	case http.MethodPut:
		var req PreferencesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		if err := s.preferences.Save(storage.PreferencesState{LastStreamID: req.LastStreamID, Autoplay: req.Autoplay}); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, req)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
