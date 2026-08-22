package config

import "testing"

func validConfig() Config {
	cfg := Config{}
	applyDefaults(&cfg)
	cfg.HTTP.Address = ":8080"
	return cfg
}

func TestValidateAcceptsValidConfig(t *testing.T) {
	cfg := validConfig()
	if err := validate(&cfg); err != nil {
		t.Fatalf("validate() error = %v", err)
	}
}

func TestValidateRejectsEmptyHTTPAddress(t *testing.T) {
	cfg := validConfig()
	cfg.HTTP.Address = ""
	if err := validate(&cfg); err == nil {
		t.Fatal("validate() error = nil, want error")
	}
}

func TestValidateRejectsInvalidVolume(t *testing.T) {
	for _, volume := range []int{-1, 101} {
		cfg := validConfig()
		cfg.Player.Volume = volume
		if err := validate(&cfg); err == nil {
			t.Fatalf("validate() error = nil for volume %d, want error", volume)
		}
	}
}

func TestValidateRejectsInvalidDurations(t *testing.T) {
	tests := []struct {
		name string
		set  func(*Config)
	}{
		{"refresh interval", func(cfg *Config) { cfg.Streams.RefreshInterval = "not-a-duration" }},
		{"retry interval", func(cfg *Config) { cfg.Playback.RetryInterval = "not-a-duration" }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validConfig()
			tt.set(&cfg)
			if err := validate(&cfg); err == nil {
				t.Fatal("validate() error = nil, want error")
			}
		})
	}
}
