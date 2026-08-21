package components

import (
	"os"
	"strings"
)

// Installed describes the MasjidPi capabilities enabled on this appliance.
type Installed struct {
	Listen bool `json:"listen"`
	Board  bool `json:"board"`
}

// Current returns the installed component profile. Installations created before
// component profiles existed retain the historical Listen + Board behaviour.
func Current() Installed {
	value := strings.TrimSpace(os.Getenv("MASJIDPI_COMPONENTS"))
	if value == "" {
		return Installed{Listen: true, Board: true}
	}

	installed := Installed{}
	for _, item := range strings.Split(value, ",") {
		switch strings.TrimSpace(strings.ToLower(item)) {
		case "listen":
			installed.Listen = true
		case "board":
			installed.Board = true
		case "both":
			installed.Listen = true
			installed.Board = true
		}
	}
	return installed
}
