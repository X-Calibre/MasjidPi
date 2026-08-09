package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/api"
	"github.com/X-Calibre/MasjidPi/backend/internal/config"
	"github.com/X-Calibre/MasjidPi/backend/internal/livestatus"
	"github.com/X-Calibre/MasjidPi/backend/internal/logger"
	"github.com/X-Calibre/MasjidPi/backend/internal/playback"
	"github.com/X-Calibre/MasjidPi/backend/internal/player"
	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
	"github.com/X-Calibre/MasjidPi/backend/internal/version"
)

func Run() error {
	log := logger.New()

	log.Info(
		"Starting application",
		"name", version.AppName,
		"version", version.Version,
	)

	paths, err := config.RuntimePaths()
	if err != nil {
		return err
	}

	cfg, err := config.Load(paths.Config)
	if err != nil {
		return err
	}

	log.Info("Configuration loaded")

	streamStore, err := stream.New(paths.Catalogue)
	if err != nil {
		return err
	}

	log.Info(
		"Loaded stream catalogue",
		"streams", len(streamStore.All()),
	)

	mpv := player.New(cfg.Player.Socket)
	if err := mpv.Start(); err != nil {
		return err
	}
	defer func() {
		log.Info("Stopping MPV")
		_ = mpv.Close()
	}()

	playbackConfig, err := newPlaybackConfig(cfg)
	if err != nil {
		return err
	}
	playbackConfig.Logger = log

	playbackManager := playback.New(mpv, playbackConfig)
	if err := playbackManager.Volume(cfg.Player.Volume); err != nil {
		return err
	}

	version, err := mpv.Version()
	if err != nil {
		return err
	}

	status, err := mpv.Status()
	if err != nil {
		return err
	}

	log.Info("Player status", "status", status)
	log.Info("Connected to MPV", "version", version)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// LiveMasjid publishes mount start/stop events over MQTT. MasjidPi keeps
	// this connection open instead of repeatedly polling the LiveMasjid site.
	liveStatus := livestatus.New("livemasjid.com", 1883, log)
	liveStatus.Start(ctx)
	defer liveStatus.Close()
	playbackManager.SetAvailability(liveStatus)
	log.Info("LiveMasjid live-status monitor started")

	playbackManager.Start(ctx)

	server := api.New(
		cfg.HTTP.Address,
		log,
		playbackManager,
		streamStore,
		paths.Frontend,
	)

	go func() {
		<-ctx.Done()
		log.Info("Shutdown requested")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error("HTTP shutdown failed", "error", err)
		}
	}()

	return server.Start()
}

func newPlaybackConfig(cfg *config.Config) (playback.Config, error) {
	retryInterval, err := time.ParseDuration(cfg.Playback.RetryInterval)
	if err != nil {
		return playback.Config{}, err
	}

	reconnectDelay, err := time.ParseDuration(cfg.Playback.ReconnectDelay)
	if err != nil {
		return playback.Config{}, err
	}

	return playback.Config{
		RetryInterval:  retryInterval,
		ReconnectDelay: reconnectDelay,
	}, nil
}
