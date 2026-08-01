package app

import (
	"time"
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

	// Temporary workaround while MPV creates its IPC socket.
	// We'll replace this with WaitForSocket() in the next milestone.
	time.Sleep(1000 * time.Millisecond)

	ipc := player.NewIPC(cfg.Player.Socket)

	if err := ipc.Connect(); err != nil {
		return err
	}
	defer ipc.Close()

	cmd := player.Command{
		Command: []any{
			"get_property",
			"mpv-version",
		},
	}

	if err := ipc.Send(cmd); err != nil {
		return err
	}

	var resp player.Response

	if err := ipc.Receive(&resp); err != nil {
		return err
	}

	log.Info(
	"Connected to MPV",
	"version", resp.Data,
	"status", resp.Error,
)

	// Start the HTTP server.
	server := api.New(cfg.HTTP.Address, log)

	return server.Start()
}
