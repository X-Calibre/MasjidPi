package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/X-Calibre/MasjidPi/backend/internal/playback"
	"github.com/X-Calibre/MasjidPi/backend/internal/player"
	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
)

type apiTestPlayer struct {
	volume int
}

func (p *apiTestPlayer) Play(string) error { return nil }
func (p *apiTestPlayer) Stop() error       { return nil }
func (p *apiTestPlayer) Volume(volume int) error {
	p.volume = volume
	return nil
}
func (p *apiTestPlayer) AudioDevices() ([]player.AudioDevice, error) {
	return nil, nil
}
func (p *apiTestPlayer) AudioDevice(string) error { return nil }
func (p *apiTestPlayer) Status() (*player.Status, error) {
	return &player.Status{Volume: p.volume}, nil
}

func newAPITestServer(t *testing.T) *Server {
	t.Helper()

	player := &apiTestPlayer{volume: 100}
	manager := playback.New(player, playback.Config{})

	catalogue := filepath.Join(t.TempDir(), "streams.json")
	if err := os.WriteFile(catalogue, []byte(`[
		{"id":"activetakbeer","name":"Takbeer Audio","url":"https://relay.livemasjid.com:8443/activetakbeer"}
	]`), 0600); err != nil {
		t.Fatal(err)
	}

	streams, err := stream.New(catalogue)
	if err != nil {
		t.Fatal(err)
	}

	return &Server{playback: manager, streams: streams}
}

func decodeStatus(t *testing.T, response *httptest.ResponseRecorder) playback.Status {
	t.Helper()

	var status playback.Status
	if err := json.NewDecoder(response.Body).Decode(&status); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return status
}

func TestPlayReturnsCanonicalStatus(t *testing.T) {
	s := newAPITestServer(t)
	req := httptest.NewRequest(http.MethodPost, "/api/player/play", bytes.NewBufferString(`{"id":"activetakbeer"}`))
	res := httptest.NewRecorder()

	s.play(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}

	status := decodeStatus(t, res)
	if status.State != string(playback.StateWaiting) {
		t.Fatalf("expected waiting state, got %q", status.State)
	}
	if status.StreamID != "activetakbeer" {
		t.Fatalf("expected stream id activetakbeer, got %q", status.StreamID)
	}
	if status.StreamName != "Takbeer Audio" {
		t.Fatalf("expected stream name Takbeer Audio, got %q", status.StreamName)
	}
}

func TestPlayRejectsMissingID(t *testing.T) {
	s := newAPITestServer(t)
	req := httptest.NewRequest(http.MethodPost, "/api/player/play", bytes.NewBufferString(`{}`))
	res := httptest.NewRecorder()

	s.play(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.Code)
	}
}

func TestPlayRejectsInvalidJSON(t *testing.T) {
	s := newAPITestServer(t)
	req := httptest.NewRequest(http.MethodPost, "/api/player/play", bytes.NewBufferString(`{"id":`))
	res := httptest.NewRecorder()

	s.play(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.Code)
	}
}

func TestPlayReturnsNotFoundForUnknownStream(t *testing.T) {
	s := newAPITestServer(t)
	req := httptest.NewRequest(http.MethodPost, "/api/player/play", bytes.NewBufferString(`{"id":"unknown"}`))
	res := httptest.NewRecorder()

	s.play(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", res.Code)
	}
}

func TestPlayRejectsWrongMethod(t *testing.T) {
	s := newAPITestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/player/play", nil)
	res := httptest.NewRecorder()

	s.play(res, req)

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", res.Code)
	}
}

func TestStopReturnsCanonicalStatus(t *testing.T) {
	s := newAPITestServer(t)
	s.playback.Play(stream.Stream{
		ID: "activetakbeer", Name: "Takbeer Audio", URL: "https://relay.livemasjid.com:8443/activetakbeer",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/player/stop", nil)
	res := httptest.NewRecorder()

	s.stop(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}

	status := decodeStatus(t, res)
	if status.State != string(playback.StateIdle) {
		t.Fatalf("expected idle state, got %q", status.State)
	}
	if status.Listening {
		t.Fatal("expected listening=false after stop")
	}
}

func TestStopRejectsWrongMethod(t *testing.T) {
	s := newAPITestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/player/stop", nil)
	res := httptest.NewRecorder()

	s.stop(res, req)

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", res.Code)
	}
}

func TestVolumeReturnsCanonicalStatus(t *testing.T) {
	s := newAPITestServer(t)
	req := httptest.NewRequest(http.MethodPost, "/api/player/volume", bytes.NewBufferString(`{"volume":42}`))
	res := httptest.NewRecorder()

	s.volume(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}

	status := decodeStatus(t, res)
	if status.Volume != 42 {
		t.Fatalf("expected volume 42, got %d", status.Volume)
	}
}

func TestVolumeRejectsWrongMethod(t *testing.T) {
	s := newAPITestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/player/volume", nil)
	res := httptest.NewRecorder()

	s.volume(res, req)

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", res.Code)
	}
}

func TestStatusReturnsCanonicalStatus(t *testing.T) {
	s := newAPITestServer(t)
	s.playback.Play(stream.Stream{
		ID: "activetakbeer", Name: "Takbeer Audio", URL: "https://relay.livemasjid.com:8443/activetakbeer",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/player/status", nil)
	res := httptest.NewRecorder()

	s.status(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}

	status := decodeStatus(t, res)
	if status.StreamID != "activetakbeer" {
		t.Fatalf("expected stream id activetakbeer, got %q", status.StreamID)
	}
	if status.Listening != true {
		t.Fatal("expected listening=true")
	}
}

func TestStatusRejectsWrongMethod(t *testing.T) {
	s := newAPITestServer(t)
	req := httptest.NewRequest(http.MethodPost, "/api/player/status", nil)
	res := httptest.NewRecorder()

	s.status(res, req)

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", res.Code)
	}
}
