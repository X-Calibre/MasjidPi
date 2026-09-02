package app

import (
	"context"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/config"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/maintenance"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	masjidboardservice "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/service"
)

const masjidBoardTimetableRefreshInterval = 30 * time.Minute

func startMasjidBoard(ctx context.Context, paths config.Paths, log interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
}) (*masjidboardservice.Service, *maintenance.Service) {
	// Discovery maintenance is independent of selected-board timetable runtime.
	// It starts even before a board selection exists so the hierarchy can be
	// available to the configuration WebUI/API.
	maintenanceService := startMasjidBoardMaintenance(ctx, paths, log)

	service, err := masjidboardservice.New(masjidboardservice.Config{
		SelectionPath: paths.MasjidBoardSelection,
		CacheDir:      paths.MasjidBoardCache,
		Log:           log,
	})
	if err != nil {
		// MasjidBoard is intentionally independent from audio playback. Invalid
		// or unavailable MasjidBoard state must not stop the appliance starting.
		log.Warn("MasjidBoard service could not start", "error", err)
		return nil, maintenanceService
	}

	if !service.Configured() {
		log.Info("MasjidBoard is not configured")
	} else {
		log.Info("MasjidBoard configured", "boards", len(service.Selection().Boards))
	}

	// Timetable retrieval is asynchronous so an unavailable provider can never
	// delay the existing audio appliance startup path. The service is checked
	// periodically even when initially unconfigured so a selection saved through
	// the WebUI begins receiving automatic refreshes without an application restart.
	go monitorMasjidBoardTimetables(ctx, service, log)

	return service, maintenanceService
}

func monitorMasjidBoardTimetables(ctx context.Context, service *masjidboardservice.Service, log interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
}) {
	refresh := func() {
		if service == nil || !service.Configured() {
			return
		}
		logMasjidBoardRefreshResults(service.Refresh(ctx), log)
	}

	// Fetch immediately after startup when configured. Persisted last-known-good
	// cache data remains available if the live provider cannot be reached.
	refresh()

	ticker := time.NewTicker(masjidBoardTimetableRefreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			refresh()
		}
	}
}

func logMasjidBoardRefreshResults(results []runtime.Result, log interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
}) {
	for _, result := range results {
		args := []any{
			"catalogue_id", result.Selection.CatalogueID,
			"masjid", result.Selection.Name,
			"status", result.Status,
		}

		switch result.Status {
		case runtime.StatusCurrent:
			log.Info("MasjidBoard timetable refreshed", args...)
		case runtime.StatusStale:
			args = append(args,
				"last_successful_update", result.LastSuccessfulUpdate,
				"error", result.UpdateError,
			)
			log.Warn("MasjidBoard update failed; using last-known-good timetable", args...)
		default:
			args = append(args, "error", result.UpdateError)
			log.Warn("MasjidBoard timetable unavailable", args...)
		}

		if result.PersistenceError != nil {
			log.Warn(
				"MasjidBoard cache persistence failed",
				"catalogue_id", result.Selection.CatalogueID,
				"masjid", result.Selection.Name,
				"error", result.PersistenceError,
			)
		}
	}
}
