package economic

import "time"

const SourceName = "Jamiatul Ulama South Africa"

// Indicators is the normalized latest Islamic economic indicator row.
type Indicators struct {
	Source        string    `json:"source"`
	SourceURL     string    `json:"source_url"`
	EffectiveDate string    `json:"effective_date"`
	HijriDate     string    `json:"hijri_date,omitempty"`
	RandDollar    float64   `json:"rand_dollar"`
	Gold24Carat   float64   `json:"gold_24_carat_per_gram"`
	Gold22Carat   float64   `json:"gold_22_carat_per_gram"`
	Gold18Carat   float64   `json:"gold_18_carat_per_gram"`
	Silver        float64   `json:"silver_per_gram"`
	Nisaab        float64   `json:"nisaab"`
	MinimumMahr   float64   `json:"minimum_mahr"`
	MahrFaatimi   float64   `json:"mahr_faatimi"`
	Krugerrand     float64   `json:"krugerrand"`
	SourceUpdatedAt time.Time `json:"source_updated_at,omitempty"`
	FetchedAt      time.Time `json:"fetched_at"`
}

func (i Indicators) Valid() bool {
	return i.Source != "" && i.SourceURL != "" && i.EffectiveDate != "" &&
		i.Nisaab > 0 && i.Krugerrand > 0
}
