package api

import (
	"net/http"
)

func (s *Server) updateStatus(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		writeError(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
		return
	}
	if s.updateController == nil {
		writeError(
			w,
			http.StatusServiceUnavailable,
			"update service is unavailable",
		)
		return
	}

	state, err := s.updateController.Status()
	if err != nil {
		s.logger.Error(
			"Could not load update status",
			"error",
			err,
		)
		writeError(
			w,
			http.StatusInternalServerError,
			"could not load update status",
		)
		return
	}

	writeJSON(w, http.StatusOK, state)
}

func (s *Server) updateCheck(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		writeError(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
		return
	}
	if s.updateController == nil {
		writeError(
			w,
			http.StatusServiceUnavailable,
			"update service is unavailable",
		)
		return
	}

	state, err := s.updateController.Check(r.Context())
	if err != nil {
		s.logger.Warn(
			"Update check failed",
			"error",
			err,
		)
		writeJSON(
			w,
			http.StatusBadGateway,
			state,
		)
		return
	}

	writeJSON(w, http.StatusOK, state)
}
