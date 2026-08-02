package api

import "net/http"

type StreamResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	City     string `json:"city"`
	Country  string `json:"country"`
	Language string `json:"language"`
}

func (s *Server) streamsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
		return
	}

	streams := s.streams.All()

	response := make([]StreamResponse, 0, len(streams))

	for _, stream := range streams {
		response = append(response, StreamResponse{
			ID:       stream.ID,
			Name:     stream.Name,
			City:     stream.City,
			Country:  stream.Country,
			Language: stream.Language,
		})
	}

	writeJSON(
		w,
		http.StatusOK,
		response,
	)
}
