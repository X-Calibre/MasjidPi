package maintenance

import (
	"context"
	"fmt"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/catalogue"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/hierarchy"
)

// Service coordinates the lightweight global hierarchy refresh and the
// configured multi-location catalogue refresh. Timetable runtime is separate.
type Service struct {
	Hierarchy *hierarchy.Updater
	Catalogue *catalogue.ScopedUpdater
}

type ScheduledResult struct {
	Hierarchy hierarchy.RefreshResult
	Catalogue catalogue.RefreshResult
}

func (s *Service) RefreshScheduled(ctx context.Context, now time.Time) (ScheduledResult, error) {
	if s == nil || s.Hierarchy == nil || s.Catalogue == nil {
		return ScheduledResult{}, fmt.Errorf("masjidboard maintenance: service is not configured")
	}

	hierarchyResult, hierarchyErr := s.Hierarchy.RefreshScheduled(ctx, now)
	catalogueResult, catalogueErr := s.Catalogue.RefreshScheduled(ctx, now)

	if hierarchyErr != nil && catalogueErr != nil {
		return ScheduledResult{Hierarchy: hierarchyResult, Catalogue: catalogueResult}, fmt.Errorf("masjidboard maintenance: hierarchy refresh: %v; catalogue refresh: %v", hierarchyErr, catalogueErr)
	}
	if hierarchyErr != nil {
		return ScheduledResult{Hierarchy: hierarchyResult, Catalogue: catalogueResult}, fmt.Errorf("masjidboard maintenance: hierarchy refresh: %w", hierarchyErr)
	}
	if catalogueErr != nil {
		return ScheduledResult{Hierarchy: hierarchyResult, Catalogue: catalogueResult}, fmt.Errorf("masjidboard maintenance: catalogue refresh: %w", catalogueErr)
	}
	return ScheduledResult{Hierarchy: hierarchyResult, Catalogue: catalogueResult}, nil
}

func (s *Service) RefreshHierarchy(ctx context.Context, now time.Time) (hierarchy.RefreshResult, error) {
	if s == nil || s.Hierarchy == nil {
		return hierarchy.RefreshResult{}, fmt.Errorf("masjidboard maintenance: hierarchy updater is not configured")
	}
	return s.Hierarchy.RefreshManual(ctx, now)
}

func (s *Service) RefreshCatalogue(ctx context.Context, now time.Time) (catalogue.RefreshResult, error) {
	if s == nil || s.Catalogue == nil {
		return catalogue.RefreshResult{}, fmt.Errorf("masjidboard maintenance: catalogue updater is not configured")
	}
	return s.Catalogue.RefreshManual(ctx, now)
}
