package api

import (
	"net/http"

	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
)

type streamResponse struct {
	ID       string      `json:"id"`
	Kind     stream.Kind `json:"kind"`
	Name     string      `json:"name"`
	Location string      `json:"location,omitempty"`
	URL      string      `json:"url"`
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

	kind := stream.Kind(r.URL.Query().Get("kind"))
	if kind != "" && kind != stream.KindMasjid && kind != stream.KindRadio {
		writeError(w, http.StatusBadRequest, "invalid stream kind")
		return
	}

	streams := s.streams.All()
	response := make([]streamResponse, 0, len(streams))

	for _, item := range streams {
		itemKind := item.SourceKind()
		if kind != "" && itemKind != kind {
			continue
		}
		response = append(response, streamResponse{
			ID:       item.ID,
			Kind:     itemKind,
			Name:     item.Name,
			Location: item.Location,
			URL:      item.URL,
		})
	}

	writeJSON(w, http.StatusOK, response)
}
