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

	if err := s.player.Play(stream.URL); err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

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

	if err := s.player.Stop(); err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		map[string]string{
			"status": "stopped",
		},
	)
}
