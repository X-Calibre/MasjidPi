package display

import (
	"time"

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
	Astronomical         *Astronomical             `json:"astronomical,omitempty"`
	Announcements        []Announcement            `json:"announcements,omitempty"`
	Programmes           []Programme               `json:"programmes,omitempty"`
	Notices              []Notice                  `json:"notices,omitempty"`
	NewMoon              *NewMoon                  `json:"new_moon,omitempty"`
}

type Date struct {
	Gregorian string `json:"gregorian,omitempty"`
	Islamic   string `json:"islamic,omitempty"`
}

// ClockTime is a display-facing local wall-clock time. It is deliberately
// distinct from the internal domain type so the API's JSON shape remains a
// stable presentation contract.
type ClockTime struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

type Prayer struct {
	Key    string     `json:"key"`
	Label  string     `json:"label"`
	Adhan  *ClockTime `json:"adhan,omitempty"`
	Jamaah *ClockTime `json:"jamaah,omitempty"`
}

type JumuahService struct {
	Adhan           *ClockTime    `json:"adhan,omitempty"`
	Jamaah          *ClockTime    `json:"jamaah,omitempty"`
	EffectiveSalaah *ClockTime    `json:"effective_salaah,omitempty"`
	AlternateAdhan  *ClockTime    `json:"alternate_adhan,omitempty"`
	AlternateJamaah *ClockTime    `json:"alternate_jamaah,omitempty"`
	Khateeb         string        `json:"khateeb,omitempty"`
	Events          []JumuahEvent `json:"events,omitempty"`
}

type JumuahEvent struct {
	Code    string     `json:"code,omitempty"`
	Heading string     `json:"heading"`
	Time    *ClockTime `json:"time,omitempty"`
}

// Announcement content may contain upstream HTML. The display API transports
// it as data only; presentation code must not inject it into the DOM as trusted
// markup.
type Announcement struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

type Programme struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

type Notice struct {
	Type    string            `json:"type"`
	Title   string            `json:"title,omitempty"`
	Content string            `json:"content,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type NewMoon struct {
	Fields map[string]string `json:"fields,omitempty"`
}

type Astronomical struct {
	Suhur         *ClockTime `json:"suhur,omitempty"`
	FajrStart     *ClockTime `json:"fajr_start,omitempty"`
	Sunrise       *ClockTime `json:"sunrise,omitempty"`
	Ishraaq       *ClockTime `json:"ishraaq,omitempty"`
	Duha          *ClockTime `json:"duha,omitempty"`
	IstiwaCaution *ClockTime `json:"istiwa_caution,omitempty"`
	Istiwa        *ClockTime `json:"istiwa,omitempty"`
	ZawaalEnd     *ClockTime `json:"zawaal_end,omitempty"`
	AsrShafii     *ClockTime `json:"asr_shafii,omitempty"`
	AsrHanafi     *ClockTime `json:"asr_hanafi,omitempty"`
	Sunset        *ClockTime `json:"sunset,omitempty"`
	EshaStart     *ClockTime `json:"esha_start,omitempty"`
}
