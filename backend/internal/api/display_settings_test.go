package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/X-Calibre/MasjidPi/backend/internal/display"
)

func TestDisplaySettingsHandlerUpdatesBrightness(t *testing.T) {
	root := t.TempDir()
	device := filepath.Join(root, "touch-display-2")
	if err := os.Mkdir(device, 0755); err != nil { t.Fatal(err) }
	for name, value := range map[string]string{"brightness":"31\n","max_brightness":"31\n"} {
		if err := os.WriteFile(filepath.Join(device, name), []byte(value), 0644); err != nil { t.Fatal(err) }
	}
	server := &Server{displaySettings: display.NewController(filepath.Join(t.TempDir(), "display.json"), root)}
	request := httptest.NewRequest(http.MethodPut, "/api/display/settings", strings.NewReader(`{"brightness_percent":50}`))
	response := httptest.NewRecorder()
	server.displaySettingsHandler(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"brightness_percent":50`) {
		t.Fatalf("unexpected response: %s", response.Body.String())
	}
}

func TestDisplaySettingsHandlerRejectsRemovedColorTemperature(t *testing.T) {
	server := &Server{displaySettings: display.NewController(filepath.Join(t.TempDir(), "display.json"), t.TempDir())}
	request := httptest.NewRequest(http.MethodPut, "/api/display/settings", strings.NewReader(`{"color_temperature":"mild"}`))
	response := httptest.NewRecorder()
	server.displaySettingsHandler(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}
