package app

import (
	"fmt"

	"github.com/X-Calibre/MasjidPi/backend/internal/api"
	"github.com/X-Calibre/MasjidPi/backend/internal/version"
)

// Run starts the application.
func Run() error {
	fmt.Printf("Starting %s %s...\n", version.AppName, version.Version)

	return api.Start(":8080")
}
