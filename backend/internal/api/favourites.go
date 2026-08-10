package api

import (
	"encoding/json"
	"net/http"
)

type FavouritesRequest struct {
	IDs []string `json:"ids"`
}

func (s *Server) favouritesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ids, err := s.favourites.Load()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, FavouritesRequest{IDs: ids})

	case http.MethodPut:
		var req FavouritesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		if req.IDs == nil {
			req.IDs = []string{}
		}
		if err := s.favourites.Save(req.IDs); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, req)

	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
