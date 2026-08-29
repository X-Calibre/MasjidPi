package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/X-Calibre/MasjidPi/backend/internal/listen"
	"github.com/X-Calibre/MasjidPi/backend/internal/playback"
)

func TestListenStatusIncludesActiveAudioDevice(t *testing.T) {
	player := &apiTestPlayer{volume: 100}
	manager := playback.New(player, playback.Config{})
	const audioDevice = "alsa/plughw:CARD=vc4hdmi0,DEV=0"
	if err := manager.AudioDevice(audioDevice); err != nil {
		t.Fatalf("set audio device: %v", err)
	}

	server := &Server{
		listen:   listen.New(nil, nil),
		playback: manager,
	}
	request := httptest.NewRequest(http.MethodGet, "/api/listen/status", nil)
	response := httptest.NewRecorder()

	server.listenStatus(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusOK)
	}
	var status listenStatusResponse
	if err := json.NewDecoder(response.Body).Decode(&status); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if status.AudioDevice != audioDevice {
		t.Fatalf("audio device = %q, want %q", status.AudioDevice, audioDevice)
	}
}
