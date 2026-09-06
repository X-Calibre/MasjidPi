package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	masjidnetwork "github.com/X-Calibre/MasjidPi/backend/internal/network"
)

type fakeWiFiManager struct {
	status       masjidnetwork.WiFiStatus
	networks     []masjidnetwork.WiFiNetwork
	connectError error
	connectedTo  string
	password     string
	hidden       bool
	access       masjidnetwork.DeviceAccess
}

func (f *fakeWiFiManager) Status(context.Context) (masjidnetwork.WiFiStatus, error) {
	return f.status, nil
}

func (f *fakeWiFiManager) Scan(context.Context) ([]masjidnetwork.WiFiNetwork, error) {
	return f.networks, nil
}

func (f *fakeWiFiManager) Connect(_ context.Context, ssid, password string, hidden bool) error {
	f.connectedTo = ssid
	f.password = password
	f.hidden = hidden
	return f.connectError
}

func (f *fakeWiFiManager) DeviceAccess(context.Context) (masjidnetwork.DeviceAccess, error) {
	return f.access, nil
}

func TestDeviceAccessReturnsDHCPNetworkDetails(t *testing.T) {
	wifi := &fakeWiFiManager{access: masjidnetwork.DeviceAccess{
		IPAddress: "10.78.63.4",
		FQDN:      "zc-masjidpi-test.internal.cassim.net.za",
	}}
	server := setupTestServer(wifi)
	request := httptest.NewRequest(http.MethodGet, "/api/setup/device-access", nil)
	response := httptest.NewRecorder()

	server.deviceAccess(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if !bytes.Contains(response.Body.Bytes(), []byte(`"fqdn":"zc-masjidpi-test.internal.cassim.net.za"`)) ||
		!bytes.Contains(response.Body.Bytes(), []byte(`"ip_address":"10.78.63.4"`)) {
		t.Fatalf("unexpected body: %s", response.Body.String())
	}
}

func setupTestServer(wifi masjidnetwork.WiFiManager) *Server {
	return &Server{wifi: wifi, logger: slog.New(slog.NewTextHandler(io.Discard, nil))}
}

func TestApplianceEntryRoutesUnconfiguredDeviceToSetup(t *testing.T) {
	server := setupTestServer(&fakeWiFiManager{status: masjidnetwork.WiFiStatus{Supported: true}})
	request := httptest.NewRequest(http.MethodGet, "/appliance", nil)
	response := httptest.NewRecorder()

	server.applianceEntry(response, request)

	if response.Code != http.StatusTemporaryRedirect || response.Header().Get("Location") != "/setup.html?profile=appliance" {
		t.Fatalf("unexpected redirect: %d %q", response.Code, response.Header().Get("Location"))
	}
}

func TestApplianceEntryPreservesBoardWhenWiFiProfileExists(t *testing.T) {
	server := setupTestServer(&fakeWiFiManager{status: masjidnetwork.WiFiStatus{Supported: true, Configured: true}})
	request := httptest.NewRequest(http.MethodGet, "/appliance", nil)
	response := httptest.NewRecorder()

	server.applianceEntry(response, request)

	if response.Code != http.StatusTemporaryRedirect || response.Header().Get("Location") != "/masjidboard.html?profile=appliance" {
		t.Fatalf("unexpected redirect: %d %q", response.Code, response.Header().Get("Location"))
	}
}

func TestWiFiNetworksAllowsOnlyDeviceLoopback(t *testing.T) {
	server := setupTestServer(&fakeWiFiManager{})
	request := httptest.NewRequest(http.MethodGet, "/api/setup/wifi/networks", nil)
	request.RemoteAddr = "192.0.2.10:1234"
	response := httptest.NewRecorder()

	server.wifiNetworks(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", response.Code)
	}
}

func TestWiFiConnectPassesCredentialsWithoutReturningPassword(t *testing.T) {
	wifi := &fakeWiFiManager{}
	server := setupTestServer(wifi)
	request := httptest.NewRequest(http.MethodPost, "/api/setup/wifi/connect", bytes.NewBufferString(`{"ssid":"Home WiFi","password":"private-pass","hidden":true}`))
	request.RemoteAddr = "127.0.0.1:1234"
	response := httptest.NewRecorder()

	server.wifiConnect(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if wifi.connectedTo != "Home WiFi" || wifi.password != "private-pass" || !wifi.hidden {
		t.Fatalf("credentials not passed to manager: ssid=%q password=%q", wifi.connectedTo, wifi.password)
	}
	var payload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(response.Body.Bytes(), []byte("private-pass")) {
		t.Fatal("password was returned by the API")
	}
}

func TestWiFiConnectRejectsOversizedPassword(t *testing.T) {
	server := setupTestServer(&fakeWiFiManager{})
	body := `{"ssid":"Home","password":"` + string(bytes.Repeat([]byte("x"), 65)) + `"}`
	request := httptest.NewRequest(http.MethodPost, "/api/setup/wifi/connect", bytes.NewBufferString(body))
	request.RemoteAddr = "[::1]:1234"
	response := httptest.NewRecorder()

	server.wifiConnect(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", response.Code)
	}
}
