package network

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type runnerResponse struct {
	out []byte
	err error
}

type fakeRunner struct {
	responses []runnerResponse
	calls     [][]string
	stdin     []string
}

func (f *fakeRunner) Run(_ context.Context, args []string, stdin string) ([]byte, error) {
	f.calls = append(f.calls, append([]string(nil), args...))
	f.stdin = append(f.stdin, stdin)
	if len(f.responses) == 0 {
		return nil, errors.New("unexpected call")
	}
	response := f.responses[0]
	f.responses = f.responses[1:]
	return response.out, response.err
}

func TestScanDeduplicatesAndSortsNetworks(t *testing.T) {
	runner := &fakeRunner{responses: []runnerResponse{{out: []byte(
		":Weak:30:WPA2\n*:Masjid\\:Office:78:WPA2\n:Weak:55:WPA2\n:Guest:64:--\n",
	)}}}
	manager := newNetworkManager(runner)

	got, err := manager.Scan(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := []WiFiNetwork{
		{SSID: "Masjid:Office", Signal: 78, Security: "WPA2", Active: true},
		{SSID: "Guest", Signal: 64, Security: "Open"},
		{SSID: "Weak", Signal: 55, Security: "WPA2"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("networks = %#v, want %#v", got, want)
	}
}

func TestStatusUsesActiveSSIDRatherThanConnectionProfileName(t *testing.T) {
	runner := &fakeRunner{responses: []runnerResponse{
		{out: []byte("802-11-wireless\nloopback\n")},
		{out: []byte("*:It Hertz When IP\n:The Elder WAN\n")},
	}}
	manager := newNetworkManager(runner)

	got, err := manager.Status(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := WiFiStatus{Supported: true, Configured: true, Connected: true, SSID: "It Hertz When IP"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("status = %#v, want %#v", got, want)
	}
}

func TestConnectPassesPasswordOnlyOnStandardInput(t *testing.T) {
	runner := &fakeRunner{responses: []runnerResponse{{}}}
	manager := newNetworkManager(runner)

	if err := manager.Connect(context.Background(), "Home WiFi", "very-secret"); err != nil {
		t.Fatal(err)
	}
	if len(runner.calls) != 1 || len(runner.stdin) != 1 {
		t.Fatalf("unexpected calls: %#v", runner.calls)
	}
	if runner.stdin[0] != "very-secret\n" {
		t.Fatalf("stdin = %q", runner.stdin[0])
	}
	for _, arg := range runner.calls[0] {
		if arg == "very-secret" {
			t.Fatal("password was exposed in command arguments")
		}
	}
}

func TestConnectReturnsSafeError(t *testing.T) {
	runner := &fakeRunner{responses: []runnerResponse{{out: []byte("secret diagnostic"), err: errors.New("exit 10")}}}
	manager := newNetworkManager(runner)

	err := manager.Connect(context.Background(), "Home", "password")
	if err == nil || err.Error() != "could not connect; check the Wi-Fi password and try again" {
		t.Fatalf("error = %v", err)
	}
}

func TestDeviceAccessUsesDHCPHostnameDomainAndAddress(t *testing.T) {
	runner := &fakeRunner{responses: []runnerResponse{
		{out: []byte("wlan0:wifi:connected\nlo:loopback:connected (externally)\n")},
		{out: []byte("10.78.63.4/24\n")},
		{out: []byte("domain_name = internal.cassim.net.za | host_name = zc-masjidpi-test | ip_address = 10.78.63.4\n")},
	}}
	manager := newNetworkManager(runner)

	got, err := manager.DeviceAccess(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := DeviceAccess{IPAddress: "10.78.63.4", FQDN: "zc-masjidpi-test.internal.cassim.net.za"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("access = %#v, want %#v", got, want)
	}
}

func TestDHCPFQDNRequiresNetworkSuppliedDomain(t *testing.T) {
	if got := dhcpFQDN(map[string]string{"host_name": "masjidframe"}); got != "" {
		t.Fatalf("dhcpFQDN() = %q, want empty", got)
	}
	if got := dhcpFQDN(map[string]string{"fqdn": "Frame-01.Masjid.Example."}); got != "frame-01.masjid.example" {
		t.Fatalf("dhcpFQDN() = %q", got)
	}
}
