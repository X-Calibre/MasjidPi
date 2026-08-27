package api

import (
	"encoding/json"
	"net/http"

	"github.com/X-Calibre/MasjidPi/backend/internal/storage"
)

type PreferencesRequest struct {
	LastStreamID string `json:"last_stream_id,omitempty"`
	Autoplay     bool   `json:"autoplay"`

	SelectedMasjidID string `json:"selected_masjid_id,omitempty"`
	SelectedRadioID  string `json:"selected_radio_id,omitempty"`
	ResumeListening  bool   `json:"resume_listening"`
	MasjidVolume     int    `json:"masjid_volume"`
	RadioVolume      int    `json:"radio_volume"`
	SourceVolumesSet bool   `json:"source_volumes_set,omitempty"`
}

func preferencesResponse(state storage.PreferencesState) PreferencesRequest {
	state = state.Normalized()
	return PreferencesRequest{
		LastStreamID:     state.LastStreamID,
		Autoplay:         state.Autoplay,
		SelectedMasjidID: state.SelectedMasjidID,
		SelectedRadioID:  state.SelectedRadioID,
		ResumeListening:  state.ResumeListening,
		MasjidVolume:     state.MasjidVolume,
		RadioVolume:      state.RadioVolume,
		SourceVolumesSet: state.SourceVolumesSet,
	}
}

func (s *Server) preferencesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		state, err := s.preferences.Load()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, preferencesResponse(state))
	case http.MethodPut:
		var req PreferencesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		state := storage.PreferencesState{
			LastStreamID:     req.LastStreamID,
			Autoplay:         req.Autoplay,
			SelectedMasjidID: req.SelectedMasjidID,
			SelectedRadioID:  req.SelectedRadioID,
			ResumeListening:  req.ResumeListening,
			MasjidVolume:     req.MasjidVolume,
			RadioVolume:      req.RadioVolume,
			SourceVolumesSet: req.SourceVolumesSet,
		}
		if err := s.preferences.Save(state); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		saved, err := s.preferences.Load()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, preferencesResponse(saved))
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
