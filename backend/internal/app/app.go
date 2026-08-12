package app

import (
	"context"
	"fmt"
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
	"github.com/X-Calibre/MasjidPi/backend/internal/storage"
	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
	"github.com/X-Calibre/MasjidPi/backend/internal/version"
)

const audioDeviceCheckInterval = 2 * time.Second

func Run() error {
	log := logger.New()
	log.Info("Starting application", "name", version.AppName, "version", version.Version)

	paths, err := config.RuntimePaths()
	if err != nil {
		return fmt.Errorf("resolve runtime paths: %w", err)
	}
	cfg, err := config.Load(paths.Config)
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}
	log.Info("Configuration loaded")

	streamStore, err := stream.New(paths.Catalogue)
	if err != nil {
		return fmt.Errorf("load stream catalogue: %w", err)
	}
	log.Info("Loaded stream catalogue", "streams", len(streamStore.All()))

	mpv := player.New(cfg.Player.Socket)
	if err := mpv.Start(); err != nil {
		return fmt.Errorf("start MPV: %w", err)
	}
	defer func() {
		log.Info("Stopping MPV")
		_ = mpv.Close()
	}()

	audioDeviceState := storage.NewAudioDeviceState(paths.AudioDeviceState)
	if name, ok, err := audioDeviceState.Load(); err != nil {
		log.Warn("Could not load saved audio device", "error", err)
	} else if ok {
		if err := mpv.AudioDevice(name); err != nil {
			log.Warn("Could not restore saved audio device", "audio_device", name, "error", err)
		} else {
			log.Info("Restored audio device", "audio_device", name)
		}
	}

	playbackConfig, err := newPlaybackConfig(cfg)
	if err != nil {
		return fmt.Errorf("create playback configuration: %w", err)
	}
	playbackConfig.Logger = log

	playbackManager := playback.New(mpv, playbackConfig)
	volumeState := storage.NewVolume(paths.VolumeState)
	playbackManager.SetVolumePersistence(volumeState)
	if err := playbackManager.InitializeVolume(); err != nil {
		return fmt.Errorf("initialize hardware volume: %w", err)
	}

	playbackState := storage.NewPlayback(paths.PlaybackState)
	playbackManager.SetPersistence(playbackState)

	if streamID, ok, err := playbackState.Load(); err != nil {
		log.Warn("Could not load last playback stream", "error", err)
	} else if ok {
		selected, err := streamStore.FindByID(streamID)
		if err != nil {
			log.Warn("Last playback stream is no longer in the catalogue", "stream_id", streamID)
			if clearErr := playbackState.Clear(); clearErr != nil {
				log.Warn("Could not clear invalid playback state", "error", clearErr)
			}
		} else {
			log.Info("Resuming last playback stream", "stream_id", selected.ID, "stream_name", selected.Name)
			playbackManager.Play(*selected)
		}
	}

	mpvVersion, err := mpv.Version()
	if err != nil {
		return fmt.Errorf("get MPV version: %w", err)
	}
	status, err := mpv.Status()
	if err != nil {
		return fmt.Errorf("get MPV status: %w", err)
	}
	log.Info("Player status", "status", status)
	log.Info("Connected to MPV", "version", mpvVersion)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	liveStatus := livestatus.New("livemasjid.com", 1883, log)
	liveStatus.Start(ctx)
	defer liveStatus.Close()
	playbackManager.SetAvailability(liveStatus)
	log.Info("LiveMasjid live-status monitor started")

	playbackManager.Start(ctx)
	go monitorAudioDevice(ctx, playbackManager, mpv, audioDeviceState, log)

	favourites := storage.NewFavourites(paths.FavouritesState)
	server := api.New(cfg.HTTP.Address, log, playbackManager, streamStore, favourites, paths.Frontend, paths.Catalogue, paths.DataRoot)

	go func() {
		<-ctx.Done()
		log.Info("Shutdown requested")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error("HTTP shutdown failed", "error", err)
		}
	}()

	if err := server.Start(); err != nil {
		return fmt.Errorf("HTTP server stopped: %w", err)
	}
	return nil
}

func monitorAudioDevice(ctx context.Context, manager *playback.Manager, mpv *player.MPV, state *storage.AudioDeviceState, log interface {
	Warn(msg string, args ...any)
}) {
	ticker := time.NewTicker(audioDeviceCheckInterval)
	defer ticker.Stop()

	lastMode := ""
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			name, ok, err := state.Load()
			if err != nil || !ok || name == "" {
				continue
			}

			devices, err := mpv.AudioDevices()
			if err != nil {
				continue
			}
			available := false
			for _, device := range devices {
				if device.Name == name {
					available = true
					break
				}
			}

			if available {
				current, err := mpv.GetProperty("audio-device")
				if err != nil {
					continue
				}
				currentName, _ := current.(string)
				if currentName != name {
					if err := manager.AudioDevice(name); err != nil {
						continue
					}
					if lastMode != "restored" {
						log.Warn("Restored audio device after it became available", "audio_device", name)
						lastMode = "restored"
					}
				}
				continue
			}

			current, err := mpv.GetProperty("audio-device")
			if err != nil {
				continue
			}
			currentName, _ := current.(string)
			if currentName == name {
				if err := manager.AudioDevice("auto"); err != nil {
					continue
				}
				if lastMode != "fallback" {
					log.Warn("Audio device unavailable, falling back to automatic output", "audio_device", name)
					lastMode = "fallback"
				}
			}
		}
	}
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
	return playback.Config{RetryInterval: retryInterval, ReconnectDelay: reconnectDelay}, nil
}
