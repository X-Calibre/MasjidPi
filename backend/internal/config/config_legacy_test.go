package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMigratesLegacyMasjidPiSocket(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.yaml")
	data := []byte("http:\n  address: \":8080\"\nplayer:\n  socket: \"/run/masjidpi/mpv.sock\"\n")

	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Player.Socket != "/run/masjidframe/mpv.sock" {
		t.Fatalf("Player.Socket = %q", cfg.Player.Socket)
	}
}

func TestLoadPreservesCustomSocket(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.yaml")
	data := []byte("http:\n  address: \":8080\"\nplayer:\n  socket: \"/srv/custom/mpv.sock\"\n")

	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Player.Socket != "/srv/custom/mpv.sock" {
		t.Fatalf("Player.Socket = %q", cfg.Player.Socket)
	}
}
