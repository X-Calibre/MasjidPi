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

	server := api.New(
		cfg.HTTP.Address,
		log,
		mpv,
	)

	return server.Start()
}
