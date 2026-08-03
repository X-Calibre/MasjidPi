package api

import "net/http"

func (s *Server) updateCatalogue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
		return
	}

	if err := s.refreshCatalogue(); err != nil {
		s.logger.Error(
			"Failed to update catalogue",
			"error", err,
		)

		writeError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	s.logger.Info("Catalogue updated successfully")

	writeJSON(
		w,
		http.StatusOK,
		map[string]string{
			"message": "Catalogue updated successfully",
		},
	)
}
