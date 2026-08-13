package model

import "time"

// Board is the normalised representation of a MasjidBoard.
// Identity, date context, and all five daily prayer times are required
// for a usable board dataset. Other fields are optional enrichment.
type Board struct {
	Identity                Identity
	DateContext             DateContext
	PrayerTimes             PrayerTimes
	AstronomicalTimes       *AstronomicalTimes
	JumuahServices          []JumuahService
	Announcements           []Announcement
	Programmes              []Programme
	Notices                 []Notice
	Media                   []Media
	Banking                 *Banking
	ContributionInformation *ContributionInformation
	NewMoon                 *NewMoon
	DisplayConfiguration    *DisplayConfiguration
}

// Identity identifies the masjid represented by the board.
type Identity struct {
	SourceBoardID string
	MasjidID      string
	EnglishName   string
	ArabicName    string
	Location      string
}

// DateContext describes the date and timezone to which the board applies.
type DateContext struct {
	GregorianDate time.Time
	IslamicDate   string
	Timezone      string
}

// PrayerTimes contains the five daily prayers. These are the core purpose
// of the board and are deliberately represented as fixed semantic fields.
type PrayerTimes struct {
	Fajr    PrayerTime
	Zuhr    PrayerTime
	Asr     PrayerTime
	Maghrib PrayerTime
	Esha    PrayerTime
}

// PrayerTime contains the Adhan and Jamaah times for a prayer.
type PrayerTime struct {
	Adhan  *time.Time
	Jamaah *time.Time
}

// AstronomicalTimes contains optional astronomical/perpetual prayer-related
// times supplied by the upstream source.
type AstronomicalTimes struct {
	Suhur       *time.Time
	FajrStart   *time.Time
	Sunrise     *time.Time
	Ishraaq     *time.Time
	Duha        *time.Time
	SolarNoon   *time.Time
	ZuhrStart   *time.Time
	AsrShafii   *time.Time
	AsrHanafi   *time.Time
	Sunset      *time.Time
	EshaStart   *time.Time
}

// JumuahService represents one Jumuah service. Multiple services are
// supported because a board may have more than one.
type JumuahService struct {
	Title   string
	Adhan   *time.Time
	Lecture *time.Time
	Khutbah *time.Time
	Salah   *time.Time
	Khateeb string
}

// Announcement is optional board content.
type Announcement struct {
	Title       string
	Content     string
	VisibleFrom *time.Time
	VisibleTo   *time.Time
}

// Programme represents an optional scheduled masjid programme.
type Programme struct {
	Title       string
	Content     string
	Start       *time.Time
	End         *time.Time
	VisibleFrom *time.Time
	VisibleTo   *time.Time
}

// Notice represents optional community notice content.
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

// Media represents an optional image/poster or other display asset.
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

// Banking contains optional banking information supplied by the board.
type Banking struct {
	Content string
}

// ContributionInformation contains optional contribution/donation information.
type ContributionInformation struct {
	Content string
}

// NewMoon contains optional lunar/new-moon information.
type NewMoon struct {
	Content string
}

// DisplayConfiguration contains optional board-specific display settings.
// Detailed rendering configuration is intentionally kept outside the core
// domain model until display design is finalised.
type DisplayConfiguration struct {
	Language string
}
