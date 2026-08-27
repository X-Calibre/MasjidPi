package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/api"
	"github.com/X-Calibre/MasjidPi/backend/internal/catalogue"
	"github.com/X-Calibre/MasjidPi/backend/internal/components"
	"github.com/X-Calibre/MasjidPi/backend/internal/config"
	"github.com/X-Calibre/MasjidPi/backend/internal/listen"
	"github.com/X-Calibre/MasjidPi/backend/internal/livestatus"
	"github.com/X-Calibre/MasjidPi/backend/internal/logger"
	"github.com/X-Calibre/MasjidPi/backend/internal/playback"
	"github.com/X-Calibre/MasjidPi/backend/internal/player"
	"github.com/X-Calibre/MasjidPi/backend/internal/radio"
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

	installed := components.Current()
	log.Info("Installed component profile", "listen", installed.Listen, "board", installed.Board)
	if !installed.Listen && !installed.Board {
		return fmt.Errorf("no MasjidPi components are installed")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if !installed.Listen {
		masjidBoardService, masjidBoardMaintenance := startMasjidBoard(ctx, paths, log)
		server := newAPIServer(cfg, paths, installed, api.Dependencies{
			Logger:                 log,
			MasjidBoardService:     masjidBoardService,
			MasjidBoardMaintenance: masjidBoardMaintenance,
		})
		return runHTTPServer(ctx, server, log)
	}

	streamStore, err := stream.New(paths.Catalogue)
	if err != nil {
		return fmt.Errorf("load stream catalogue: %w", err)
	}
	streamStore.Replace(radio.Merge(streamStore.All()))
	log.Info("Loaded stream catalogue", "streams", len(streamStore.All()), "radio_stations", len(radio.Catalogue()))

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

	liveStatus := livestatus.New("livemasjid.com", 1883, log)
	liveStatus.Start(ctx)
	defer liveStatus.Close()
	log.Info("LiveMasjid live-status monitor started")

	playbackManager.Start(ctx)
	listenController := listen.New(liveStatus, playback.NewListenOutput(playbackManager))
	preferences := storage.NewPreferences(paths.PreferencesState)
	prefs, err := preferences.Load()
	if err != nil {
		return fmt.Errorf("load Listen preferences: %w", err)
	}

	listenController.SetMasjidEnabled(*prefs.MasjidEnabled)
	if err := listenController.SetRadioEnabled(*prefs.RadioEnabled); err != nil {
		return fmt.Errorf("restore radio power state: %w", err)
	}
	if err := listenController.SetMasjidVolume(prefs.MasjidVolume); err != nil {
		return fmt.Errorf("restore masjid volume: %w", err)
	}
	if err := listenController.SetRadioVolume(prefs.RadioVolume); err != nil {
		return fmt.Errorf("restore radio volume: %w", err)
	}
	if err := listenController.SetRadioResumeDelayMinutes(prefs.RadioResumeDelayMinutes); err != nil {
		return fmt.Errorf("restore radio resume delay: %w", err)
	}
	if err := listenController.SetRadioSchedule(prefs.RadioScheduleEnabled, prefs.RadioScheduleStart, prefs.RadioScheduleStop); err != nil {
		return fmt.Errorf("restore radio schedule: %w", err)
	}

	if prefs.SelectedMasjidID != "" {
		selected, findErr := streamStore.FindByID(prefs.SelectedMasjidID)
		if findErr != nil {
			log.Warn("Selected masjid is no longer in the catalogue", "stream_id", prefs.SelectedMasjidID)
		} else if selectErr := listenController.SelectMasjid(selected); selectErr != nil {
			log.Warn("Could not restore selected masjid", "stream_id", selected.ID, "error", selectErr)
		} else {
			log.Info("Restored selected masjid", "stream_id", selected.ID, "stream_name", selected.Name)
		}
	}
	if prefs.SelectedRadioID != "" {
		selected, findErr := streamStore.FindByID(prefs.SelectedRadioID)
		if findErr != nil {
			log.Warn("Selected radio station is no longer in the catalogue", "stream_id", prefs.SelectedRadioID)
		} else if selectErr := listenController.SelectRadio(selected); selectErr != nil {
			log.Warn("Could not restore selected radio station", "stream_id", selected.ID, "error", selectErr)
		} else {
			log.Info("Restored selected radio station", "stream_id", selected.ID, "stream_name", selected.Name)
		}
	}
	if *prefs.RadioEnabled {
		if err := listenController.SetRadioMode(listen.RadioMode(prefs.RadioMode)); err != nil {
			return fmt.Errorf("restore radio mode: %w", err)
		}
	}
	if prefs.ResumeListening {
		listenController.Listen()
		log.Info("Restoring Listen active state")
	}
	listenController.Start(ctx)

	go monitorAudioDevice(ctx, playbackManager, mpv, audioDeviceState, log)

	favourites := storage.NewFavourites(paths.FavouritesState)
	dependencies := api.Dependencies{
		Logger:           log,
		Playback:         playbackManager,
		Listen:           listenController,
		Streams:          streamStore,
		Favourites:       favourites,
		Preferences:      preferences,
		AudioDeviceState: audioDeviceState,
	}
	if installed.Board {
		masjidBoardService, masjidBoardMaintenance := startMasjidBoard(ctx, paths, log)
		dependencies.MasjidBoardService = masjidBoardService
		dependencies.MasjidBoardMaintenance = masjidBoardMaintenance
	}
	server := newAPIServer(cfg, paths, installed, dependencies)

	catalogueRefreshInterval, err := time.ParseDuration(cfg.Streams.RefreshInterval)
	if err != nil {
		return fmt.Errorf("parse catalogue refresh interval: %w", err)
	}
	if catalogueRefreshInterval <= 0 {
		return fmt.Errorf("catalogue refresh interval must be greater than zero")
	}
	go monitorCatalogueRefresh(ctx, catalogueRefreshInterval, paths.Catalogue, streamStore, log)
	return runHTTPServer(ctx, server, log)
}

func newAPIServer(cfg *config.Config, paths config.Paths, installed components.Installed, dependencies api.Dependencies) *api.Server {
	return api.New(api.Config{
		Address:                  cfg.HTTP.Address,
		Frontend:                 paths.Frontend,
		CatalogueFile:            paths.Catalogue,
		CatalogueDataRoot:        paths.DataRoot,
		PreferencesPath:          paths.PreferencesState,
		Installed:                installed,
		MasjidBoardHierarchyPath: paths.MasjidBoardHierarchy,
		MasjidBoardCataloguePath: paths.MasjidBoardCatalogue,
		MasjidBoardScopePath:     paths.MasjidBoardScope,
	}, dependencies)
}

func runHTTPServer(ctx context.Context, server *api.Server, log interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}) error {
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

func monitorCatalogueRefresh(ctx context.Context, interval time.Duration, catalogueFile string, streams *stream.Store, log interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			updated, err := catalogue.Update(catalogueFile)
			if err != nil {
				log.Warn("Scheduled catalogue update failed", "error", err)
				continue
			}
			streams.Replace(radio.Merge(updated))
			log.Info("Scheduled catalogue refresh completed", "streams", len(streams.All()), "radio_stations", len(radio.Catalogue()))
		}
	}
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

func newPlaybackConfig(cfg *config.Config) (playback.Config, error) {
	retryInterval, err := time.ParseDuration(cfg.Playback.RetryInterval)
	if err != nil {
		return playback.Config{}, err
	}
	return playback.Config{RetryInterval: retryInterval}, nil
}
