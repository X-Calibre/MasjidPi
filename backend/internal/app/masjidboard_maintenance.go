package app

import (
	"context"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/config"
	masjidboardcatalogue "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/catalogue"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/hierarchy"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/maintenance"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/provider/masjidboardlive"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/scope"
)

const masjidBoardMaintenanceCheckInterval = 24 * time.Hour

func startMasjidBoardMaintenance(ctx context.Context, paths config.Paths, log interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
}) *maintenance.Service {
	discovery := masjidboardlive.DiscoveryClient{}
	service := &maintenance.Service{
		Hierarchy: &hierarchy.Updater{
			Source: masjidboardlive.HierarchySource{Client: discovery},
			Store:  hierarchy.NewStore(paths.MasjidBoardHierarchy),
		},
		Catalogue: &masjidboardcatalogue.ScopedUpdater{
			Scope: scope.NewStore(paths.MasjidBoardScope),
			Updater: masjidboardcatalogue.Updater{
				Source: masjidboardlive.CatalogueSource{Client: discovery},
				Store:  masjidboardcatalogue.NewStore(paths.MasjidBoardCatalogue),
			},
		},
	}

	go monitorMasjidBoardMaintenance(ctx, service, log)
	return service
}

func monitorMasjidBoardMaintenance(ctx context.Context, service *maintenance.Service, log interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
}) {
	run := func() {
		result, err := service.RefreshScheduled(ctx, time.Now())
		if err != nil {
			log.Warn("Scheduled MasjidBoard maintenance failed", "error", err)
		} else if result.Hierarchy.Attempted || result.Catalogue.AnyAttempted() {
			log.Info("Scheduled MasjidBoard maintenance completed",
				"hierarchy_attempted", result.Hierarchy.Attempted,
				"hierarchy_updated", result.Hierarchy.Updated,
				"catalogue_attempted", result.Catalogue.AnyAttempted(),
				"catalogue_failed", result.Catalogue.AnyFailed(),
			)
		}
	}

	// Check immediately after startup. Persisted timestamps keep this cheap when
	// the weekly refresh is not due.
	run()

	ticker := time.NewTicker(masjidBoardMaintenanceCheckInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			run()
		}
	}
}
