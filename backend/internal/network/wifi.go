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

// WiFiManager provides the small subset of networking operations required by
// the first-run appliance interface.
type WiFiManager interface {
	Status(context.Context) (WiFiStatus, error)
	Scan(context.Context) ([]WiFiNetwork, error)
	Connect(context.Context, string, string) error
}
