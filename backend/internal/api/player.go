package api

import (
	"encoding/json"
	"net/http"
)

type PlayRequest struct {
	ID string `json:"id"`
}

func (s *Server) play(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
		return
	}

	var req PlayRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid JSON",
		)
		return
	}

	if req.ID == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"missing ID",
		)
		return
	}

	stream, err := s.streams.FindByID(req.ID)
	if err != nil {
		writeError(
			w,
			http.StatusNotFound,
			err.Error(),
		)
		return
	}

	s.playback.Play(*stream)

	writeJSON(
		w,
		http.StatusOK,
		map[string]string{
			"status": "playing",
		},
	)
}

func (s *Server) stop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
		return
	}

	s.playback.Stop()

	writeJSON(
		w,
		http.StatusOK,
		map[string]string{
			"status": "stopped",
		},
	)
}
