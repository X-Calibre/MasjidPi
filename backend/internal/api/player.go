package api

import (
	"encoding/json"
	"net/http"
)

type PlayRequest struct {
	URL string `json:"url"`
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

	if req.URL == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"missing URL",
		)
		return
	}

	if err := s.player.Play(req.URL); err != nil {
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
