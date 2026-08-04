package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func defaultSocketPath() string {
	return filepath.Join(
		os.TempDir(),
		fmt.Sprintf("masjidpi-%d.sock", os.Getuid()),
	)
}

func applyDefaults(cfg *Config) {

	if cfg.Player.Socket == "" {
		cfg.Player.Socket = defaultSocketPath()
	}

	if cfg.Player.Volume == 0 {
		cfg.Player.Volume = 100
	}
}
