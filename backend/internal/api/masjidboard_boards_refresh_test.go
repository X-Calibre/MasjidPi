package api

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	masjidboardruntime "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

type fakeBoardRefreshRuntime struct {
	refreshCalls int
	state        selection.State
}

func (f *fakeBoardRefreshRuntime) Configured() bool { return f.state.Configured() }
func (f *fakeBoardRefreshRuntime) Selection() selection.State { return f.state }
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
	server := New(":0", logger, nil, nil, nil, t.TempDir(), "", t.TempDir())
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
