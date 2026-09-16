package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/updates"
)

type fakeUpdateController struct {
	statusState updates.State
	statusErr   error
	checkState  updates.State
	checkErr    error
	statusCalls int
	checkCalls  int
}

func (f *fakeUpdateController) Status() (
	updates.State,
	error,
) {
	f.statusCalls++
	return f.statusState, f.statusErr
}

func (f *fakeUpdateController) Check(
	context.Context,
) (updates.State, error) {
	f.checkCalls++
	return f.checkState, f.checkErr
}

func newUpdateTestServer(
	t *testing.T,
	controller updateController,
) *Server {
	t.Helper()

	root := t.TempDir()
	logger := slog.New(
		slog.NewTextHandler(io.Discard, nil),
	)

	return New(
		Config{
			Address:         ":0",
			Frontend:        root,
			PreferencesPath: root + "/preferences.json",
		},
		Dependencies{
			Logger:  logger,
			Updates: controller,
		},
	)
}

func performUpdateRequest(
	t *testing.T,
	server *Server,
	method string,
	path string,
) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(method, path, nil)
	recorder := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(
		recorder,
		request,
	)

	return recorder
}

func decodeUpdateState(
	t *testing.T,
	recorder *httptest.ResponseRecorder,
) updates.State {
	t.Helper()

	var state updates.State
	if err := json.Unmarshal(
		recorder.Body.Bytes(),
		&state,
	); err != nil {
		t.Fatalf(
			"decode response %q: %v",
			recorder.Body.String(),
			err,
		)
	}
	return state
}

func TestUpdateStatus(t *testing.T) {
	detected := time.Date(
		2026, time.September, 16,
		9, 0, 0, 0, time.UTC,
	)
	release := updates.Release{
		Version: "v1.7.0",
	}
	controller := &fakeUpdateController{
		statusState: updates.State{
			SchemaVersion:    updates.StateSchemaVersion,
			CurrentVersion:   "v1.6.0",
			AvailableRelease: &release,
			FirstDetectedAt:  &detected,
		},
	}
	server := newUpdateTestServer(t, controller)

	recorder := performUpdateRequest(
		t,
		server,
		http.MethodGet,
		"/api/update/status",
	)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"status = %d, body = %s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	state := decodeUpdateState(t, recorder)
	if state.CurrentVersion != "v1.6.0" {
		t.Fatalf(
			"current version = %q",
			state.CurrentVersion,
		)
	}
	if state.AvailableRelease == nil ||
		state.AvailableRelease.Version != "v1.7.0" {
		t.Fatalf(
			"available release = %+v",
			state.AvailableRelease,
		)
	}
	if controller.statusCalls != 1 {
		t.Fatalf(
			"status calls = %d, want 1",
			controller.statusCalls,
		)
	}
	if controller.checkCalls != 0 {
		t.Fatalf(
			"check calls = %d, want 0",
			controller.checkCalls,
		)
	}
}

func TestUpdateCheck(t *testing.T) {
	release := updates.Release{
		Version: "v1.7.0",
	}
	controller := &fakeUpdateController{
		checkState: updates.State{
			SchemaVersion:    updates.StateSchemaVersion,
			CurrentVersion:   "v1.6.0",
			AvailableRelease: &release,
		},
	}
	server := newUpdateTestServer(t, controller)

	recorder := performUpdateRequest(
		t,
		server,
		http.MethodPost,
		"/api/update/check",
	)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"status = %d, body = %s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	state := decodeUpdateState(t, recorder)
	if state.AvailableRelease == nil ||
		state.AvailableRelease.Version != "v1.7.0" {
		t.Fatalf(
			"available release = %+v",
			state.AvailableRelease,
		)
	}
	if controller.checkCalls != 1 {
		t.Fatalf(
			"check calls = %d, want 1",
			controller.checkCalls,
		)
	}
	if controller.statusCalls != 0 {
		t.Fatalf(
			"status calls = %d, want 0",
			controller.statusCalls,
		)
	}
}

func TestUpdateCheckReturnsRecordedStateOnDiscoveryFailure(
	t *testing.T,
) {
	controller := &fakeUpdateController{
		checkState: updates.State{
			SchemaVersion:  updates.StateSchemaVersion,
			CurrentVersion: "v1.6.0",
			LastCheckError: "temporary GitHub failure",
		},
		checkErr: errors.New("temporary GitHub failure"),
	}
	server := newUpdateTestServer(t, controller)

	recorder := performUpdateRequest(
		t,
		server,
		http.MethodPost,
		"/api/update/check",
	)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf(
			"status = %d, body = %s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	state := decodeUpdateState(t, recorder)
	if state.LastCheckError != "temporary GitHub failure" {
		t.Fatalf(
			"last check error = %q",
			state.LastCheckError,
		)
	}
}

func TestUpdateEndpointsRejectWrongMethods(t *testing.T) {
	controller := &fakeUpdateController{}
	server := newUpdateTestServer(t, controller)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "status post",
			method: http.MethodPost,
			path:   "/api/update/status",
		},
		{
			name:   "check get",
			method: http.MethodGet,
			path:   "/api/update/check",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := performUpdateRequest(
				t,
				server,
				test.method,
				test.path,
			)

			if recorder.Code != http.StatusMethodNotAllowed {
				t.Fatalf(
					"status = %d, body = %s",
					recorder.Code,
					recorder.Body.String(),
				)
			}
		})
	}

	if controller.statusCalls != 0 ||
		controller.checkCalls != 0 {
		t.Fatalf(
			"unexpected calls: status=%d check=%d",
			controller.statusCalls,
			controller.checkCalls,
		)
	}
}

func TestUpdateEndpointsRequireController(t *testing.T) {
	server := newUpdateTestServer(t, nil)

	tests := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/api/update/status",
		},
		{
			method: http.MethodPost,
			path:   "/api/update/check",
		},
	}

	for _, test := range tests {
		recorder := performUpdateRequest(
			t,
			server,
			test.method,
			test.path,
		)

		if recorder.Code != http.StatusServiceUnavailable {
			t.Fatalf(
				"%s %s status = %d, body = %s",
				test.method,
				test.path,
				recorder.Code,
				recorder.Body.String(),
			)
		}
	}
}

func TestUpdateStatusLoadFailure(t *testing.T) {
	controller := &fakeUpdateController{
		statusErr: errors.New("state unavailable"),
	}
	server := newUpdateTestServer(t, controller)

	recorder := performUpdateRequest(
		t,
		server,
		http.MethodGet,
		"/api/update/status",
	)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"status = %d, body = %s",
			recorder.Code,
			recorder.Body.String(),
		)
	}
}
