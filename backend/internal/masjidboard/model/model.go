package model

import "time"

// Board is the normalised representation of a MasjidBoard.
// Identity and the five daily prayer times are the required core.
// Everything else is optional enrichment.
type Board struct {
	Identity       BoardIdentity
	DateContext    DateContext
	PrayerTimes    PrayerTimes
	Astronomical   *AstronomicalTimes
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
type DateContext struct {
	GregorianDate time.Time
	IslamicDate   string
}

// PrayerTimes contains the five daily prayers and optional Friday Jumu'ah
// services. On Friday, Jumu'ah replaces Dhuhr as the congregational prayer.
type PrayerTimes struct {
	Fajr    PrayerTime
	Dhuhr   PrayerTime
	Asr     PrayerTime
	Maghrib PrayerTime
	Esha    PrayerTime
	Jumuah  []JumuahService
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

// JumuahService represents one Friday congregational service. Multiple
// services are supported. Individual fields may be absent.
type JumuahService struct {
	Label   string
	Adhan   *ClockTime
	Jamaah  *ClockTime
	Khateeb string
}

type AstronomicalTimes struct {
	Suhur     *ClockTime
	FajrStart *ClockTime
	Sunrise   *ClockTime
	Ishraaq   *ClockTime
	Duha      *ClockTime
	AsrShafii *ClockTime
	AsrHanafi *ClockTime
	Sunset    *ClockTime
	EshaStart *ClockTime
}

type Announcement struct {
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
	VisibleFrom *time.Time
	VisibleTo   *time.Time
}

type NoticeType string

const (
	NoticeTypeGeneral  NoticeType = "general"
	NoticeTypeNikah    NoticeType = "nikah"
	NoticeTypeFuneral  NoticeType = "funeral"
	NoticeTypeWellWish NoticeType = "well_wishes"
	NoticeTypeEid      NoticeType = "eid"
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
}

type NewMoon struct {
	Content string
}
