package display

import (
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
	masjidboardruntime "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
)

// View is the read-only presentation model consumed by a MasjidBoard display.
// It deliberately excludes discovery, configuration, provider metadata and
// diagnostic error strings. Boards remain in the user's selected order.
type View struct {
	Configured bool    `json:"configured"`
	Boards     []Board `json:"boards"`
}

type Board struct {
	CatalogueID          string                    `json:"catalogue_id"`
	Name                 string                    `json:"name"`
	AlternateName        string                    `json:"alternate_name,omitempty"`
	TimeZone             string                    `json:"time_zone,omitempty"`
	Status               masjidboardruntime.Status `json:"status"`
	Stale                bool                      `json:"stale"`
	LastSuccessfulUpdate *time.Time                `json:"last_successful_update,omitempty"`
	Date                 Date                      `json:"date"`
	Prayers              []Prayer                  `json:"prayers,omitempty"`
	Jumuah               []JumuahService           `json:"jumuah,omitempty"`
	Astronomical         *Astronomical              `json:"astronomical,omitempty"`
}

type Date struct {
	Gregorian string `json:"gregorian,omitempty"`
	Islamic   string `json:"islamic,omitempty"`
}

type Prayer struct {
	Key    string           `json:"key"`
	Label  string           `json:"label"`
	Adhan  *model.ClockTime `json:"adhan,omitempty"`
	Jamaah *model.ClockTime `json:"jamaah,omitempty"`
}

type JumuahService struct {
	Adhan           *model.ClockTime `json:"adhan,omitempty"`
	Jamaah          *model.ClockTime `json:"jamaah,omitempty"`
	EffectiveSalaah *model.ClockTime `json:"effective_salaah,omitempty"`
	AlternateAdhan  *model.ClockTime `json:"alternate_adhan,omitempty"`
	AlternateJamaah *model.ClockTime `json:"alternate_jamaah,omitempty"`
	Khateeb         string           `json:"khateeb,omitempty"`
	Events          []JumuahEvent    `json:"events,omitempty"`
}

type JumuahEvent struct {
	Code    string           `json:"code,omitempty"`
	Heading string           `json:"heading"`
	Time    *model.ClockTime `json:"time,omitempty"`
}

type Astronomical struct {
	Suhur         *model.ClockTime `json:"suhur,omitempty"`
	FajrStart     *model.ClockTime `json:"fajr_start,omitempty"`
	Sunrise       *model.ClockTime `json:"sunrise,omitempty"`
	Ishraaq       *model.ClockTime `json:"ishraaq,omitempty"`
	Duha          *model.ClockTime `json:"duha,omitempty"`
	IstiwaCaution *model.ClockTime `json:"istiwa_caution,omitempty"`
	Istiwa        *model.ClockTime `json:"istiwa,omitempty"`
	ZawaalEnd     *model.ClockTime `json:"zawaal_end,omitempty"`
	AsrShafii     *model.ClockTime `json:"asr_shafii,omitempty"`
	AsrHanafi     *model.ClockTime `json:"asr_hanafi,omitempty"`
	Sunset        *model.ClockTime `json:"sunset,omitempty"`
	EshaStart     *model.ClockTime `json:"esha_start,omitempty"`
}
