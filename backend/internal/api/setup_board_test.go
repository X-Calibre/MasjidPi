package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	masjidnetwork "github.com/X-Calibre/MasjidFrame/backend/internal/network"
	"github.com/X-Calibre/MasjidFrame/backend/internal/storage"
)

func TestBoardSetupDeferralPersistsAndRoutesToBoard(t *testing.T) {
	preferences := storage.NewPreferences(filepath.Join(t.TempDir(), "preferences.json"))
	server := setupTestServer(&fakeWiFiManager{status: masjidnetwork.WiFiStatus{
		Supported: true, Configured: true, Connected: true,
	}})
	server.preferences = preferences
	server.masjidBoardService = fakeMasjidBoardStatusProvider{configured: false}

	request := httptest.NewRequest(http.MethodPut, "/api/setup/board", bytes.NewBufferString(`{"deferred":true}`))
	response := httptest.NewRecorder()
	server.boardSetup(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}

	state, err := preferences.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !state.BoardSetupDeferred {
		t.Fatal("BoardSetupDeferred = false, want true")
	}

	request = httptest.NewRequest(http.MethodGet, "/appliance", nil)
	response = httptest.NewRecorder()
	server.applianceEntry(response, request)
	if response.Code != http.StatusTemporaryRedirect ||
		response.Header().Get("Location") != "/masjidboard.html?profile=appliance-720" {
		t.Fatalf("unexpected redirect: %d %q", response.Code, response.Header().Get("Location"))
	}
}

func TestBoardSetupResumeClearsDeferral(t *testing.T) {
	preferences := storage.NewPreferences(filepath.Join(t.TempDir(), "preferences.json"))
	if _, err := preferences.Update(func(state *storage.PreferencesState) {
		state.BoardSetupDeferred = true
	}); err != nil {
		t.Fatal(err)
	}
	server := setupTestServer(nil)
	server.preferences = preferences

	request := httptest.NewRequest(http.MethodPut, "/api/setup/board", bytes.NewBufferString(`{"deferred":false}`))
	response := httptest.NewRecorder()
	server.boardSetup(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	state, err := preferences.Load()
	if err != nil {
		t.Fatal(err)
	}
	if state.BoardSetupDeferred {
		t.Fatal("BoardSetupDeferred = true, want false")
	}
}

func TestBoardSetupRejectsInvalidRequest(t *testing.T) {
	server := setupTestServer(nil)
	server.preferences = storage.NewPreferences(filepath.Join(t.TempDir(), "preferences.json"))
	request := httptest.NewRequest(http.MethodPut, "/api/setup/board", bytes.NewBufferString(`{"deferred":`))
	response := httptest.NewRecorder()
	server.boardSetup(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", response.Code)
	}
}
