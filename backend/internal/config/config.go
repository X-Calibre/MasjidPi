package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config contains all application configuration.
type Config struct {
	HTTP     HTTPConfig     `yaml:"http"`
	Player   PlayerConfig   `yaml:"player"`
	Streams  StreamConfig   `yaml:"streams"`
	Playback PlaybackConfig `yaml:"playback"`
}

// HTTPConfig contains HTTP server settings.
type HTTPConfig struct {
	Address string `yaml:"address"`
}

// PlayerConfig contains player settings.
type PlayerConfig struct {
	Socket string `yaml:"socket"`
	Volume int    `yaml:"volume"`
}

// StreamConfig contains stream-related settings.
type StreamConfig struct {
	RefreshInterval string `yaml:"refresh_interval"`
}

// PlaybackConfig contains playback retry settings.
type PlaybackConfig struct {
	RetryInterval  string `yaml:"retry_interval"`
	ReconnectDelay string `yaml:"reconnect_delay"`
}

// Load reads a YAML configuration file from disk.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := &Config{}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	applyDefaults(cfg)

	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return cfg, nil
}
