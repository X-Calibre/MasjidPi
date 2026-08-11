package config

import (
	"fmt"
	"time"
)

func validate(cfg *Config) error {
	if cfg.HTTP.Address == "" {
		return fmt.Errorf("http address must not be empty")
	}

	if cfg.Player.Socket == "" {
		return fmt.Errorf("player socket must not be empty")
	}

	if cfg.Player.Volume < 0 || cfg.Player.Volume > 100 {
		return fmt.Errorf("player volume must be between 0 and 100")
	}

	if cfg.Streams.RefreshInterval != "" {
		if _, err := time.ParseDuration(cfg.Streams.RefreshInterval); err != nil {
			return fmt.Errorf("streams refresh_interval: %w", err)
		}
	}

	if _, err := time.ParseDuration(cfg.Playback.RetryInterval); err != nil {
		return fmt.Errorf("playback retry_interval: %w", err)
	}

	if _, err := time.ParseDuration(cfg.Playback.ReconnectDelay); err != nil {
		return fmt.Errorf("playback reconnect_delay: %w", err)
	}

	return nil
}
