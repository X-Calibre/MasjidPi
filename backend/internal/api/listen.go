package api

import (
	"encoding/json"
	"net/http"

	"github.com/X-Calibre/MasjidPi/backend/internal/listen"
	"github.com/X-Calibre/MasjidPi/backend/internal/storage"
)

type listenStatusResponse struct {
	listen.Status
	MasterVolume          int    `json:"master_volume"`
	MasterVolumeSupported bool   `json:"master_volume_supported"`
	PlaybackState         string `json:"playback_state"`
	PlaybackMessage       string `json:"playback_message,omitempty"`
	PlaybackURL           string `json:"playback_url,omitempty"`
	PlaybackEndpoint      string `json:"playback_endpoint,omitempty"`
	PlaybackFallbackUsed  bool   `json:"playback_fallback_used,omitempty"`
}

type listenSelectionRequest struct {
	MasjidID *string `json:"masjid_id"`
	RadioID  *string `json:"radio_id"`
}

type listenVolumeRequest struct {
	Source string `json:"source"`
	Volume int    `json:"volume"`
}

type listenRadioDelayRequest struct {
	Minutes int `json:"minutes"`
}

func (s *Server) listenStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if s.listen == nil || s.playback == nil {
		writeError(w, http.StatusServiceUnavailable, "listen controller is not configured")
		return
	}

	controllerStatus := s.listen.Status()
	playbackStatus := s.playback.Status()
	writeJSON(w, http.StatusOK, listenStatusResponse{
		Status:                controllerStatus,
		MasterVolume:          playbackStatus.Volume,
		MasterVolumeSupported: playbackStatus.VolumeSupported,
		PlaybackState:         playbackStatus.State,
		PlaybackMessage:       playbackStatus.Message,
		PlaybackURL:           playbackStatus.URL,
		PlaybackEndpoint:      playbackStatus.Endpoint,
		PlaybackFallbackUsed:  playbackStatus.FallbackUsed,
	})
}

func (s *Server) listenSelection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if s.listen == nil || s.streams == nil || s.preferences == nil {
		writeError(w, http.StatusServiceUnavailable, "listen controller is not configured")
		return
	}

	var req listenSelectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.MasjidID == nil && req.RadioID == nil {
		writeError(w, http.StatusBadRequest, "masjid_id or radio_id is required")
		return
	}

	prefs, err := s.preferences.Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if req.MasjidID != nil {
		if *req.MasjidID == "" {
			if err := s.listen.SelectMasjid(nil); err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			prefs.SelectedMasjidID = ""
			prefs.LastStreamID = ""
		} else {
			selected, err := s.streams.FindByID(*req.MasjidID)
			if err != nil {
				writeError(w, http.StatusNotFound, err.Error())
				return
			}
			if err := s.listen.SelectMasjid(selected); err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			prefs.SelectedMasjidID = selected.ID
			prefs.LastStreamID = selected.ID
		}
	}

	if req.RadioID != nil {
		if *req.RadioID == "" {
			if err := s.listen.SelectRadio(nil); err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			prefs.SelectedRadioID = ""
		} else {
			selected, err := s.streams.FindByID(*req.RadioID)
			if err != nil {
				writeError(w, http.StatusNotFound, err.Error())
				return
			}
			if err := s.listen.SelectRadio(selected); err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			prefs.SelectedRadioID = selected.ID
		}
	}

	if err := s.preferences.Save(prefs); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, s.listen.Status())
}

func (s *Server) listenVolume(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if s.listen == nil || s.preferences == nil {
		writeError(w, http.StatusServiceUnavailable, "listen controller is not configured")
		return
	}

	var req listenVolumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.Volume < 0 || req.Volume > listen.MaxSourceVolume {
		writeError(w, http.StatusBadRequest, "source volume must be between 0 and 150")
		return
	}

	prefs, err := s.preferences.Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	switch req.Source {
	case "masjid":
		if err := s.listen.SetMasjidVolume(req.Volume); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		prefs.MasjidVolume = req.Volume
	case "radio":
		if err := s.listen.SetRadioVolume(req.Volume); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		prefs.RadioVolume = req.Volume
	default:
		writeError(w, http.StatusBadRequest, "source must be masjid or radio")
		return
	}
	prefs.SourceVolumesSet = true

	if err := s.preferences.Save(prefs); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, s.listen.Status())
}

func (s *Server) listenRadioDelay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if s.listen == nil || s.preferences == nil {
		writeError(w, http.StatusServiceUnavailable, "listen controller is not configured")
		return
	}

	var req listenRadioDelayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := s.listen.SetRadioResumeDelayMinutes(req.Minutes); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	prefs, err := s.preferences.Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	prefs.RadioResumeDelayMinutes = req.Minutes
	if err := s.preferences.Save(prefs); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, s.listen.Status())
}

func (s *Server) listenStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if s.listen == nil || s.preferences == nil {
		writeError(w, http.StatusServiceUnavailable, "listen controller is not configured")
		return
	}

	s.listen.Listen()
	if err := s.persistResumeListening(true); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, s.listen.Status())
}

func (s *Server) listenStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if s.listen == nil || s.preferences == nil {
		writeError(w, http.StatusServiceUnavailable, "listen controller is not configured")
		return
	}

	s.listen.Stop()
	if err := s.persistResumeListening(false); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, s.listen.Status())
}

func (s *Server) persistResumeListening(enabled bool) error {
	prefs, err := s.preferences.Load()
	if err != nil {
		return err
	}
	prefs.ResumeListening = enabled
	prefs.Autoplay = enabled
	return s.preferences.Save(storage.PreferencesState(prefs))
}
