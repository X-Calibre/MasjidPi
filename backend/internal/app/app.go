package app

import (
	"github.com/X-Calibre/MasjidPi/backend/internal/api"
	"github.com/X-Calibre/MasjidPi/backend/internal/config"
	"github.com/X-Calibre/MasjidPi/backend/internal/logger"
	"github.com/X-Calibre/MasjidPi/backend/internal/player"
	"github.com/X-Calibre/MasjidPi/backend/internal/version"
)

func Run() error {
	log := logger.New()

	log.Info(
		"Starting application",
		"name", version.AppName,
		"version", version.Version,
	)

	cfg, err := config.Load("configs/default.yaml")
	if err != nil {
		return err
	}

	log.Info("Configuration loaded")

	mpv := player.New(cfg.Player.Socket)

	if err := mpv.Start(); err != nil {
		return err
	}

	version, err := mpv.Version()
	if err != nil {
		return err
	}

	log.Info(
		"Connected to MPV",
		"version", version,
	)

	const testStream = "https://relay.livemasjid.com:8443/activetakbeer"

	log.Info(
		"Playing test stream",
		"url", testStream,
	)

	if err := mpv.Play(testStream); err != nil {
		return err
	}

	log.Info("Playback started")

	server := api.New(cfg.HTTP.Address, log)

	return server.Start()
}
