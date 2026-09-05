package model

import "time"

// Board is the normalised representation of a MasjidBoard.
// Identity and the five daily prayer times are the required core.
// Everything else is optional enrichment.
type Board struct {
	Identity      BoardIdentity
	DateContext   DateContext
	PrayerTimes   PrayerTimes
	Astronomical  *AstronomicalTimes
	Announcements []Announcement
	Programmes    []Programme
	Notices       []Notice
	Media         []Media
	Banking       *Banking
	NewMoon       *NewMoon
}

// BoardIdentity identifies the masjid represented by the board.
type BoardIdentity struct {
	ID            string
	Name          string
	AlternateName string
	TimeZone      string
}

// DateContext describes the date and timezone to which the board applies.
// IslamicDateAdjustment and ForceIslamicDate30 preserve the upstream
// moon-sighting context so a cached board can keep its Islamic date current
// without needing another network refresh at sunset.
type DateContext struct {
	GregorianDate         time.Time
	IslamicDate           string
	IslamicDateAdjustment int
	ForceIslamicDate30    bool
}

// PrayerTimes contains the five daily prayers and optional Friday Jumu'ah
// services. On Friday, Jumu'ah replaces Dhuhr as the congregational prayer.
type PrayerTimes struct {
	Fajr         PrayerTime
	Dhuhr        PrayerTime
	Asr          PrayerTime
	Maghrib      PrayerTime
	Esha         PrayerTime
	SpecialDhuhr *SpecialPrayerTime
	Jumuah       []JumuahService
}

// ClockTime is a local wall-clock time. The board timezone is stored once on
// BoardIdentity rather than duplicated in every prayer time.
type ClockTime struct {
	Hour   int
	Minute int
}

// PrayerTime contains the optional Adhan and Jamaah times for one prayer.
// The prayer itself is part of the required five-prayer schedule; either
// value may be absent when the upstream board does not supply it.
type PrayerTime struct {
	Adhan  *ClockTime
	Jamaah *ClockTime
}

// SpecialPrayerTime is an additional provider-supplied prayer time whose
// applicability is described by the provider label, for example
// "(Sundays & Public Holidays)" or "(Everyday)".
type SpecialPrayerTime struct {
	Time  *ClockTime
	Label string
}

// JumuahService represents one Friday congregational service.
//
// MasjidBoard Live supplies three configurable heading/time pairs in addition
// to dedicated Jumu'ah Adhan, Jamaah, Islamic-time and Khateeb fields.
// The event heading and source code are preserved so presentation-specific
// translation can be applied later without losing the upstream value.
type JumuahService struct {
	Adhan         *ClockTime
	Jamaah        *ClockTime
	IslamicAdhan  *ClockTime
	IslamicJamaah *ClockTime
	// AlternateAdhan and AlternateJamaah are retained only so older cached
	// boards can be decoded. The Premium source columns are Islamic-time
	// representations, not additional Jumu'ah services.
	AlternateAdhan  *ClockTime
	AlternateJamaah *ClockTime
	Khateeb         string
	Events          []JumuahEvent
}

// JumuahEvent is one of the three detailed Friday heading/time pairs supplied
// by MasjidBoard Live. Code is the source jumuahHeadingsArray value; Heading
// is the human-readable heading present in the board data.
type JumuahEvent struct {
	Code    string
	Heading string
	Time    *ClockTime
}

// EffectiveSalaah returns only an explicitly supplied Jumu'ah Jamaah/Salaah
// time. A Khutbah event is not treated as Salaah because the source may expose
// Khutbah without publishing the actual prayer time.
func (s JumuahService) EffectiveSalaah() *ClockTime {
	return s.Jamaah
}

type AstronomicalTimes struct {
	Suhur         *ClockTime
	FajrStart     *ClockTime
	Sunrise       *ClockTime
	Ishraaq       *ClockTime
	Duha          *ClockTime
	IstiwaCaution *ClockTime
	Istiwa        *ClockTime
	ZawaalEnd     *ClockTime
	AsrShafii     *ClockTime
	AsrHanafi     *ClockTime
	Sunset        *ClockTime
	EshaStart     *ClockTime
}

type Announcement struct {
	Category    string
	Title       string
	Content     string
	VisibleFrom *time.Time
	VisibleTo   *time.Time
}

type Programme struct {
	Title       string
	Content     string
	Start       *time.Time
	End         *time.Time
	VisibleFrom *time.Time
	VisibleTo   *time.Time
}

type Notice struct {
	Type        NoticeType
	Title       string
	Content     string
	Fields      map[string]string
	VisibleFrom *time.Time
	VisibleTo   *time.Time
}

type NoticeType string

const (
	NoticeTypeGeneral      NoticeType = "general"
	NoticeTypeNikah        NoticeType = "nikah"
	NoticeTypeFuneral      NoticeType = "funeral"
	NoticeTypeWellWish     NoticeType = "well_wishes"
	NoticeTypeEid          NoticeType = "eid"
	NoticeTypeSalaahChange NoticeType = "salaah_change"
	NoticeTypeDawah        NoticeType = "dawah"
	NoticeTypeThreeDay     NoticeType = "three_day_jamaat"
)

type Media struct {
	SourceID       string
	SourceURL      string
	LocalPath      string
	MediaType      string
	ContentHash    string
	VisibleFrom    *time.Time
	VisibleTo      *time.Time
	DisplaySeconds int
}

type Banking struct {
	Content string
	Fields  map[string]string
}

type NewMoon struct {
	Fields map[string]string
}
