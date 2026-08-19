package app

import (
	"context"

	"github.com/X-Calibre/MasjidPi/backend/internal/config"
	masjidboardservice "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/service"
)

func startMasjidBoard(ctx context.Context, paths config.Paths, log interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
}) *masjidboardservice.Service {
	// Discovery maintenance is independent of selected-board timetable runtime.
	// It starts even before a board selection exists so the hierarchy can be
	// available to the configuration WebUI/API.
	startMasjidBoardMaintenance(ctx, paths, log)

	service, err := masjidboardservice.New(masjidboardservice.Config{
		SelectionPath: paths.MasjidBoardSelection,
		CacheDir:      paths.MasjidBoardCache,
	})
	if err != nil {
		// MasjidBoard is intentionally independent from audio playback. Invalid
		// or unavailable MasjidBoard state must not stop the appliance starting.
		log.Warn("MasjidBoard service could not start", "error", err)
		return nil
	}

	if !service.Configured() {
		log.Info("MasjidBoard is not configured")
		return service
	}

	log.Info("MasjidBoard configured", "boards", len(service.Selection().Boards))

	// Initial board retrieval is asynchronous so an unavailable timetable
	// provider cannot delay the existing audio appliance startup path.
	go func() {
		results := service.Refresh(ctx)
		for _, result := range results {
			args := []any{
				"catalogue_id", result.Selection.CatalogueID,
				"masjid", result.Selection.Name,
				"status", result.Status,
			}

			switch result.Status {
			case "current":
				log.Info("MasjidBoard timetable refreshed", args...)
			case "stale":
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
	}()

	return service
}
