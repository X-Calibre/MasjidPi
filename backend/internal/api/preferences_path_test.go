package api

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/X-Calibre/MasjidPi/backend/internal/components"
)

func TestPreferencesUseConfiguredPersistentPath(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "persistent", "preferences.json")
	t.Setenv("MASJIDPI_HOME", filepath.Join(root, "replaceable-runtime"))
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := New(Config{
		Address:         ":0",
		Frontend:        root,
		PreferencesPath: path,
		Installed:       components.Installed{Listen: true},
	}, Dependencies{Logger: logger})

	req := httptest.NewRequest(http.MethodPut, "/api/preferences", strings.NewReader(`{"last_stream_id":"test","autoplay":true}`))
	rec := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("configured preferences path was not written: %v", err)
	}
	legacy := filepath.Join(os.Getenv("MASJIDPI_HOME"), "backend", "data", "preferences.json")
	if _, err := os.Stat(legacy); !os.IsNotExist(err) {
		t.Fatalf("replaceable runtime path was written: %v", err)
	}
}
