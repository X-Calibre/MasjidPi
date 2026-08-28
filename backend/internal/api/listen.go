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
type listenRadioScheduleRequest struct {
	Enabled bool   `json:"enabled"`
	Start   string `json:"start"`
	Stop    string `json:"stop"`
}
type listenRadioModeRequest struct {
	Mode string `json:"mode"`
}
type listenPowerRequest struct {
	Module  string `json:"module"`
	Enabled bool   `json:"enabled"`
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
	var selectedMasjidID, selectedRadioID *string
	if req.MasjidID != nil {
		if *req.MasjidID == "" {
			_ = s.listen.SelectMasjid(nil)
			value := ""
			selectedMasjidID = &value
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
			value := selected.ID
			selectedMasjidID = &value
		}
	}
	if req.RadioID != nil {
		if *req.RadioID == "" {
			_ = s.listen.SelectRadio(nil)
			value := ""
			selectedRadioID = &value
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
			value := selected.ID
			selectedRadioID = &value
		}
	}
	if _, err := s.preferences.Update(func(prefs *storage.PreferencesState) {
		if selectedMasjidID != nil {
			prefs.SelectedMasjidID = *selectedMasjidID
			prefs.LastStreamID = *selectedMasjidID
		}
		if selectedRadioID != nil {
			prefs.SelectedRadioID = *selectedRadioID
		}
	}); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, s.listen.Status())
}

func (s *Server) listenPower(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if s.listen == nil || s.preferences == nil {
		writeError(w, http.StatusServiceUnavailable, "listen controller is not configured")
		return
	}
	var req listenPowerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	var update func(*storage.PreferencesState)
	switch req.Module {
	case "masjid":
		s.listen.SetMasjidEnabled(req.Enabled)
		update = func(prefs *storage.PreferencesState) { prefs.MasjidEnabled = boolPtr(req.Enabled) }
		if !req.Enabled {
			update = func(prefs *storage.PreferencesState) {
				prefs.MasjidEnabled = boolPtr(false)
				prefs.RadioEnabled = boolPtr(false)
				prefs.ResumeListening = false
				prefs.Autoplay = false
			}
		}
	case "radio":
		if req.Enabled && !s.listen.Status().MasjidEnabled {
			s.listen.SetMasjidEnabled(true)
		}
		if err := s.listen.SetRadioEnabled(req.Enabled); err != nil {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		update = func(prefs *storage.PreferencesState) {
			if req.Enabled {
				prefs.MasjidEnabled = boolPtr(true)
			}
			prefs.RadioEnabled = boolPtr(req.Enabled)
		}
	default:
		writeError(w, http.StatusBadRequest, "module must be masjid or radio")
		return
	}
	if _, err := s.preferences.Update(update); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, s.listen.Status())
}

func boolPtr(value bool) *bool { return &value }

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
	var update func(*storage.PreferencesState)
	switch req.Source {
	case "masjid":
		err := s.listen.SetMasjidVolume(req.Volume)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		update = func(prefs *storage.PreferencesState) { prefs.MasjidVolume = req.Volume }
	case "radio":
		err := s.listen.SetRadioVolume(req.Volume)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		update = func(prefs *storage.PreferencesState) { prefs.RadioVolume = req.Volume }
	default:
		writeError(w, http.StatusBadRequest, "source must be masjid or radio")
		return
	}
	if _, err := s.preferences.Update(func(prefs *storage.PreferencesState) {
		update(prefs)
		prefs.SourceVolumesSet = true
	}); err != nil {
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
	var req listenRadioDelayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := s.listen.SetRadioResumeDelayMinutes(req.Minutes); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if _, err := s.preferences.Update(func(prefs *storage.PreferencesState) { prefs.RadioResumeDelayMinutes = req.Minutes }); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, s.listen.Status())
}

func (s *Server) listenRadioSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req listenRadioScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := s.listen.SetRadioSchedule(req.Enabled, req.Start, req.Stop); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if _, err := s.preferences.Update(func(prefs *storage.PreferencesState) {
		prefs.RadioScheduleEnabled = req.Enabled
		prefs.RadioScheduleStart = req.Start
		prefs.RadioScheduleStop = req.Stop
	}); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, s.listen.Status())
}

func (s *Server) listenRadioMode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if s.listen == nil || s.preferences == nil {
		writeError(w, http.StatusServiceUnavailable, "listen controller is not configured")
		return
	}
	var req listenRadioModeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	mode := listen.RadioMode(req.Mode)
	if err := s.listen.SetRadioMode(mode); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	if _, err := s.preferences.Update(func(prefs *storage.PreferencesState) {
		if mode == listen.RadioModeStopped {
			prefs.RadioMode = "stopped"
		} else {
			prefs.RadioMode = "schedule"
		}
	}); err != nil {
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
	s.listen.Listen()
	listening := s.listen.Status().Listening
	if err := s.persistResumeListening(listening); err != nil {
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
	s.listen.Stop()
	if err := s.persistResumeListening(false); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, s.listen.Status())
}

func (s *Server) persistResumeListening(enabled bool) error {
	_, err := s.preferences.Update(func(prefs *storage.PreferencesState) {
		prefs.ResumeListening = enabled
		prefs.Autoplay = enabled
	})
	return err
}
