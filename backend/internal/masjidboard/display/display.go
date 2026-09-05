package display

import (
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/dailycontent"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/economic"
	masjidboardruntime "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

// View is the read-only presentation model consumed by a MasjidBoard display.
// It deliberately excludes discovery, configuration, provider metadata and
// diagnostic error strings. Boards remain in the user's selected order.
type View struct {
	Configured          bool                 `json:"configured"`
	Theme               string               `json:"theme"`
	SlideDuration       int                  `json:"slide_duration_seconds"`
	ShowDuaAfterAdhan   bool                 `json:"show_dua_after_adhan"`
	Boards              []Board              `json:"boards"`
	EconomicIndicators  *economic.Indicators `json:"economic_indicators,omitempty"`
	DailyIslamicContent *DailyIslamicContent `json:"daily_islamic_content,omitempty"`
}

type DailyIslamicContent struct {
	Ayah        *dailycontent.Ayah   `json:"ayah,omitempty"`
	Hadith      *dailycontent.Hadith `json:"hadith,omitempty"`
	Sunnah      *dailycontent.Sunnah `json:"sunnah,omitempty"`
	Language    string               `json:"language"`
	Source      string               `json:"source"`
	SourceURL   string               `json:"source_url"`
	ContentDate string               `json:"content_date,omitempty"`
	FetchedAt   time.Time            `json:"fetched_at"`
}

func PresentDailyIslamicContent(content *dailycontent.Content, selected selection.State) *DailyIslamicContent {
	if content == nil || !content.Valid() || !selected.ShowAnyDailyIslamicContent() {
		return nil
	}
	presented := &DailyIslamicContent{
		Language: content.Language, Source: content.Source, SourceURL: content.SourceURL,
		ContentDate: content.ContentDate, FetchedAt: content.FetchedAt,
	}
	if selected.ShowDailyAyahValue() {
		copy := content.Ayah
		presented.Ayah = &copy
	}
	if selected.ShowDailyHadithValue() {
		copy := content.Hadith
		presented.Hadith = &copy
	}
	if selected.ShowDailySunnahValue() {
		copy := content.Sunnah
		presented.Sunnah = &copy
	}
	return presented
}

type Board struct {
	CatalogueID        string                    `json:"catalogue_id"`
	Name               string                    `json:"name"`
	AlternateName      string                    `json:"alternate_name,omitempty"`
	TimeZone           string                    `json:"time_zone,omitempty"`
	Status             masjidboardruntime.Status `json:"status"`
	Stale              bool                      `json:"stale"`
	ShowDetailedJumuah bool                      `json:"show_detailed_jumuah"`
	Date               Date                      `json:"date"`
	Prayers            []Prayer                  `json:"prayers,omitempty"`
	SpecialDhuhr       *SpecialPrayerTime         `json:"special_dhuhr,omitempty"`
	Jumuah             []JumuahService           `json:"jumuah,omitempty"`
	Astronomical       *Astronomical             `json:"astronomical,omitempty"`
	Announcements      []Announcement            `json:"announcements,omitempty"`
	Programmes         []Programme               `json:"programmes,omitempty"`
	Notices            []Notice                  `json:"notices,omitempty"`
	Banking            *Banking                  `json:"banking,omitempty"`
	NewMoon            *NewMoon                  `json:"new_moon,omitempty"`
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

type SpecialPrayerTime struct {
	Time  *ClockTime `json:"time,omitempty"`
	Label string     `json:"label,omitempty"`
}

type JumuahService struct {
	Adhan           *ClockTime    `json:"adhan,omitempty"`
	Jamaah          *ClockTime    `json:"jamaah,omitempty"`
	EffectiveSalaah *ClockTime    `json:"effective_salaah,omitempty"`
	IslamicAdhan    *ClockTime    `json:"islamic_adhan,omitempty"`
	IslamicJamaah   *ClockTime    `json:"islamic_jamaah,omitempty"`
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
	Category string `json:"category,omitempty"`
	Title    string `json:"title,omitempty"`
	Content  string `json:"content,omitempty"`
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

type Banking struct {
	Title  string            `json:"title,omitempty"`
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
