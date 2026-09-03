package dailycontent

import "time"

const (
	SourceName = "MasjidBoard Live"
	SourceURL  = "https://api.masjidboardlive.com/mblfileapi"
)

// Content is the normalized, masjid-independent daily Islamic content
// published by MasjidBoard Live.
type Content struct {
	Ayah        Ayah      `json:"ayah"`
	Hadith      Hadith    `json:"hadith"`
	Sunnah      Sunnah    `json:"sunnah"`
	Language    string    `json:"language"`
	Source      string    `json:"source"`
	SourceURL   string    `json:"source_url"`
	ContentDate string    `json:"content_date,omitempty"`
	FetchedAt   time.Time `json:"fetched_at"`
}

type Ayah struct {
	Surah      string `json:"surah"`
	AyahNumber string `json:"ayah_number"`
	Text       string `json:"text"`
}

type Hadith struct {
	Heading   string `json:"heading"`
	Text      string `json:"text"`
	Reference string `json:"reference,omitempty"`
}

type Sunnah struct {
	Heading   string `json:"heading"`
	Text      string `json:"text"`
	Reference string `json:"reference,omitempty"`
}

// Valid reports whether content is complete enough to replace the
// last-known-good cache. References are optional because some upstream
// languages do not currently supply them.
func (c Content) Valid() bool {
	return c.Source == SourceName && c.SourceURL != "" && c.Language != "" && !c.FetchedAt.IsZero() &&
		c.Ayah.Surah != "" && c.Ayah.AyahNumber != "" && c.Ayah.Text != "" &&
		c.Hadith.Heading != "" && c.Hadith.Text != "" &&
		c.Sunnah.Heading != "" && c.Sunnah.Text != ""
}
