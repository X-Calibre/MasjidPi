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

	if cfg.Playback.RetryInterval == "" {
		cfg.Playback.RetryInterval = "5s"
	}

	if cfg.Playback.ReconnectDelay == "" {
		cfg.Playback.ReconnectDelay = "5s"
	}
}
