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

	// Start MPV in idle mode.
	proc := player.NewProcess(cfg.Player.Socket)

	if err := proc.Start(); err != nil {
		return err
	}

	log.Info("MPV started")

	// Start the HTTP server.
	server := api.New(cfg.HTTP.Address, log)

	return server.Start()
}
