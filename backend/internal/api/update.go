package api

import (
	"encoding/json"
	"net/http"
	"time"
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

func (s *Server) updateApprove(
	w http.ResponseWriter,
	r *http.Request,
) {
	if !s.requireUpdatePost(w, r) {
		return
	}

	state, err := s.updateController.Approve()
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, state)
}

func (s *Server) updatePostpone(
	w http.ResponseWriter,
	r *http.Request,
) {
	if !s.requireUpdatePost(w, r) {
		return
	}

	var request struct {
		Until time.Time `json:"until"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil ||
		request.Until.IsZero() {
		writeError(
			w,
			http.StatusBadRequest,
			"valid postponement time is required",
		)
		return
	}

	state, err := s.updateController.Postpone(request.Until)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, state)
}

func (s *Server) requireUpdatePost(
	w http.ResponseWriter,
	r *http.Request,
) bool {
	if r.Method != http.MethodPost {
		writeError(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
		return false
	}
	if s.updateController == nil {
		writeError(
			w,
			http.StatusServiceUnavailable,
			"update service is unavailable",
		)
		return false
	}
	return true
}
