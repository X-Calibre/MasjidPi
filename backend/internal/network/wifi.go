package network

import "context"

// WiFiNetwork is a wireless access point visible to the appliance.
type WiFiNetwork struct {
	SSID     string `json:"ssid"`
	Signal   int    `json:"signal"`
	Security string `json:"security"`
	Active   bool   `json:"active"`
}

// WiFiStatus describes whether appliance Wi-Fi setup is available and needed.
type WiFiStatus struct {
	Supported  bool   `json:"supported"`
	Configured bool   `json:"configured"`
	Connected  bool   `json:"connected"`
	SSID       string `json:"ssid,omitempty"`
}

// DeviceAccess contains LAN addressing supplied by the active network. FQDN
// remains empty when DHCP and network DNS do not provide a usable name.
type DeviceAccess struct {
	IPAddress string `json:"ip_address,omitempty"`
	FQDN      string `json:"fqdn,omitempty"`
}

// WiFiManager provides the small subset of networking operations required by
// the first-run appliance interface.
type WiFiManager interface {
	Status(context.Context) (WiFiStatus, error)
	Scan(context.Context) ([]WiFiNetwork, error)
	Connect(context.Context, string, string) error
	DeviceAccess(context.Context) (DeviceAccess, error)
}
