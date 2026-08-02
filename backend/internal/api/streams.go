package api

import "net/http"

type streamResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location,omitempty"`
	URL      string `json:"url"`
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

	response := make([]streamResponse, 0, len(streams))

	for _, stream := range streams {
		response = append(response, streamResponse{
			ID:       stream.ID,
			Name:     stream.Name,
			Location: stream.Location,
			URL:      stream.URL,
		})
	}

	writeJSON(
		w,
		http.StatusOK,
		response,
	)
}
