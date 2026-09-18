package api

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	masjidtimezone "github.com/X-Calibre/MasjidPi/backend/internal/timezone"
)

type fakeTimezoneController struct {
	zones        []masjidtimezone.Zone
	zonesError   error
	current      string
	configured   bool
	currentError error
	setName      string
	setError     error
}

func (f *fakeTimezoneController) Zones(
	string,
) ([]masjidtimezone.Zone, error) {
	return f.zones, f.zonesError
}

func (f *fakeTimezoneController) Current() (string, bool, error) {
	return f.current, f.configured, f.currentError
}

func (f *fakeTimezoneController) Set(
	_ context.Context,
	name string,
) error {
	f.setName = name
	return f.setError
}

func TestTimezonesReturnsCountryOptions(t *testing.T) {
	controller := &fakeTimezoneController{
		zones: []masjidtimezone.Zone{{
			Name: "Africa/Johannesburg",
		}},
	}
	server := setupTestServer(&fakeWiFiManager{})
	server.timezoneController = controller

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/setup/timezones?country=South%20Africa",
		nil,
	)
	request.RemoteAddr = "127.0.0.1:1234"
	response := httptest.NewRecorder()

	server.timezones(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"status = %d, body = %s",
			response.Code,
			response.Body.String(),
		)
	}
	if !bytes.Contains(
		response.Body.Bytes(),
		[]byte(`"name":"Africa/Johannesburg"`),
	) {
		t.Fatalf("unexpected body: %s", response.Body.String())
	}
}

func TestTimezonesRejectsRemoteRequest(t *testing.T) {
	server := setupTestServer(&fakeWiFiManager{})
	server.timezoneController = &fakeTimezoneController{}

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/setup/timezones?country=South%20Africa",
		nil,
	)
	request.RemoteAddr = "192.0.2.10:1234"
	response := httptest.NewRecorder()

	server.timezones(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", response.Code)
	}
}

func TestTimezoneAppliesSelectedZone(t *testing.T) {
	controller := &fakeTimezoneController{}
	server := setupTestServer(&fakeWiFiManager{})
	server.timezoneController = controller

	request := httptest.NewRequest(
		http.MethodPut,
		"/api/setup/timezone",
		bytes.NewBufferString(`{"name":"Africa/Johannesburg"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.RemoteAddr = "[::1]:1234"
	response := httptest.NewRecorder()

	server.timezone(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"status = %d, body = %s",
			response.Code,
			response.Body.String(),
		)
	}
	if controller.setName != "Africa/Johannesburg" {
		t.Fatalf("Set() name = %q", controller.setName)
	}
}

func TestTimezoneDoesNotHideControllerError(t *testing.T) {
	controller := &fakeTimezoneController{
		setError: errors.New("unsupported timezone"),
	}
	server := setupTestServer(&fakeWiFiManager{})
	server.timezoneController = controller

	request := httptest.NewRequest(
		http.MethodPut,
		"/api/setup/timezone",
		bytes.NewBufferString(`{"name":"Invalid/Zone"}`),
	)
	request.RemoteAddr = "127.0.0.1:1234"
	response := httptest.NewRecorder()

	server.timezone(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", response.Code)
	}
}

func TestTimezoneRejectsRemoteRequest(t *testing.T) {
	server := setupTestServer(&fakeWiFiManager{})
	server.timezoneController = &fakeTimezoneController{}

	request := httptest.NewRequest(
		http.MethodPut,
		"/api/setup/timezone",
		bytes.NewBufferString(`{"name":"Africa/Johannesburg"}`),
	)
	request.RemoteAddr = "192.0.2.10:1234"
	response := httptest.NewRecorder()

	server.timezone(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", response.Code)
	}
}
