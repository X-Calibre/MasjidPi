package masjidboardlive

import (
	"context"
	"fmt"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/provider"
)

// EnrichedClient keeps the public Core board as the authoritative timetable
// provider and adds optional community content from the Premium board when it
// is available. A Premium failure never turns a successful Core refresh into
// a stale or unavailable board.
type EnrichedClient struct {
	Core    provider.Provider
	Premium provider.Provider
}

func (c EnrichedClient) Fetch(ctx context.Context) (model.Board, error) {
	if c.Core == nil {
		return model.Board{}, fmt.Errorf("masjidboardlive: Core provider is required")
	}
	core, err := c.Core.Fetch(ctx)
	if err != nil {
		return model.Board{}, err
	}
	if c.Premium == nil {
		return core, nil
	}

	premium, err := c.Premium.Fetch(ctx)
	if err != nil {
		return core, nil
	}
	mergePremiumEnrichment(&core, premium)
	return core, nil
}

func mergePremiumEnrichment(core *model.Board, premium model.Board) {
	if core == nil {
		return
	}
	core.Announcements = premium.Announcements
	core.Programmes = premium.Programmes
	core.Notices = premium.Notices
	core.Media = premium.Media
	core.Banking = premium.Banking
	core.NewMoon = premium.NewMoon
	core.PrayerTimes.SpecialDhuhr = premium.PrayerTimes.SpecialDhuhr
}
