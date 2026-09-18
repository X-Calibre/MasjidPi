package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/updates"
)

type fakeUpdateController struct {
	statusState       updates.State
	statusErr         error
	checkState        updates.State
	checkErr          error
	approveState      updates.State
	approveErr        error
	postponeState     updates.State
	postponeErr       error
	installState      updates.State
	installErr        error
	statusCalls       int
	checkCalls        int
	approveCalls      int
	postponeCalls     int
	installCalls      int
	postponedUntil    time.Time
	installImmediate  bool
	interruptPlayback bool
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

func (f *fakeUpdateController) Approve() (updates.State, error) {
	f.approveCalls++
	return f.approveState, f.approveErr
}

func (f *fakeUpdateController) Postpone(
	until time.Time,
) (updates.State, error) {
	f.postponeCalls++
	f.postponedUntil = until
	return f.postponeState, f.postponeErr
}

func (f *fakeUpdateController) Install(
	_ context.Context,
	immediate bool,
	interruptPlayback bool,
) (updates.State, error) {
	f.installCalls++
	f.installImmediate = immediate
	f.interruptPlayback = interruptPlayback
	return f.installState, f.installErr
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

func TestUpdateApprove(t *testing.T) {
	approved := time.Date(2026, time.September, 16, 9, 0, 0, 0, time.UTC)
	controller := &fakeUpdateController{
		approveState: updates.State{
			SchemaVersion: updates.StateSchemaVersion,
			ApprovedAt:    &approved,
		},
	}
	server := newUpdateTestServer(t, controller)
	recorder := performUpdateRequest(t, server, http.MethodPost, "/api/update/approve")

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	state := decodeUpdateState(t, recorder)
	if state.ApprovedAt == nil || !state.ApprovedAt.Equal(approved) {
		t.Fatalf("approved at = %v, want %v", state.ApprovedAt, approved)
	}
	if controller.approveCalls != 1 {
		t.Fatalf("approve calls = %d, want 1", controller.approveCalls)
	}
}

func TestUpdatePostpone(t *testing.T) {
	until := time.Date(2026, time.September, 23, 9, 0, 0, 0, time.UTC)
	controller := &fakeUpdateController{
		postponeState: updates.State{
			SchemaVersion:  updates.StateSchemaVersion,
			PostponedUntil: &until,
		},
	}
	server := newUpdateTestServer(t, controller)
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/update/postpone",
		strings.NewReader(`{"until":"2026-09-23T09:00:00Z"}`),
	)
	recorder := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if !controller.postponedUntil.Equal(until) {
		t.Fatalf("postponed until = %v, want %v", controller.postponedUntil, until)
	}
}

func TestUpdatePostponeRejectsInvalidRequest(t *testing.T) {
	controller := &fakeUpdateController{}
	server := newUpdateTestServer(t, controller)
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/update/postpone",
		strings.NewReader(`{"until":"not-a-date"}`),
	)
	recorder := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if controller.postponeCalls != 0 {
		t.Fatalf("postpone calls = %d, want 0", controller.postponeCalls)
	}
}

func TestUpdateInstallRequiresExplicitPlaybackInterruption(t *testing.T) {
	controller := &fakeUpdateController{installState: updates.State{
		SchemaVersion: updates.StateSchemaVersion,
		Installation: &updates.InstallationState{
			Status: updates.InstallStatusRebootPending,
		},
	}}
	server := newUpdateTestServer(t, controller)
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/update/install",
		strings.NewReader(`{"interrupt_playback":true}`),
	)
	recorder := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if controller.installCalls != 1 || !controller.installImmediate ||
		!controller.interruptPlayback {
		t.Fatalf(
			"install calls=%d immediate=%t interrupt=%t",
			controller.installCalls,
			controller.installImmediate,
			controller.interruptPlayback,
		)
	}
}

func TestUpdateInstallReturnsConflictReason(t *testing.T) {
	controller := &fakeUpdateController{
		installState: updates.State{SchemaVersion: updates.StateSchemaVersion},
		installErr:   errors.New("audio playback is active"),
	}
	server := newUpdateTestServer(t, controller)
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/update/install",
		strings.NewReader(`{"interrupt_playback":false}`),
	)
	recorder := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), "audio playback is active") {
		t.Fatalf("body = %s", recorder.Body.String())
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
		{
			name:   "approve get",
			method: http.MethodGet,
			path:   "/api/update/approve",
		},
		{
			name:   "postpone get",
			method: http.MethodGet,
			path:   "/api/update/postpone",
		},
		{
			name:   "install get",
			method: http.MethodGet,
			path:   "/api/update/install",
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
		{
			method: http.MethodPost,
			path:   "/api/update/approve",
		},
		{
			method: http.MethodPost,
			path:   "/api/update/postpone",
		},
		{
			method: http.MethodPost,
			path:   "/api/update/install",
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
