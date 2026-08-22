package api

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/X-Calibre/MasjidPi/backend/internal/components"
	masjidboardruntime "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

type fakeBoardRefreshRuntime struct {
	refreshCalls int
	state        selection.State
}

func (f *fakeBoardRefreshRuntime) Configured() bool                     { return f.state.Configured() }
func (f *fakeBoardRefreshRuntime) Selection() selection.State           { return f.state }
func (f *fakeBoardRefreshRuntime) Results() []masjidboardruntime.Result { return nil }
func (f *fakeBoardRefreshRuntime) Reconfigure(state selection.State) error {
	f.state = state
	return nil
}
func (f *fakeBoardRefreshRuntime) Refresh(context.Context) []masjidboardruntime.Result {
	f.refreshCalls++
	return []masjidboardruntime.Result{{}}
}

func TestMasjidBoardBoardsRefresh(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	root := t.TempDir()
	server := New(Config{Address: ":0", Frontend: root, PreferencesPath: root + "/preferences.json", Installed: components.Installed{Listen: true, Board: true}}, Dependencies{Logger: logger})
	runtime := &fakeBoardRefreshRuntime{}
	server.SetMasjidBoardService(runtime)

	req := httptest.NewRequest(http.MethodPost, "/api/masjidboard/boards/refresh", nil)
	rec := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if runtime.refreshCalls != 1 {
		t.Fatalf("refreshCalls=%d, want 1", runtime.refreshCalls)
	}
}
