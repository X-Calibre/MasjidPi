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
	Gold14Carat   float64   `json:"gold_14_carat_per_gram"`
	Gold9Carat    float64   `json:"gold_9_carat_per_gram"`
	Silver        float64   `json:"silver_per_gram"`
	Nisaab        float64   `json:"nisaab"`
	MinimumMahr   float64   `json:"minimum_mahr"`
	MahrFaatimi   float64   `json:"mahr_faatimi"`
	Krugerrand    float64   `json:"krugerrand"`
	FetchedAt     time.Time `json:"fetched_at"`
}

func (i Indicators) Valid() bool {
	return i.Source != "" && i.SourceURL != "" && i.EffectiveDate != "" &&
		i.Nisaab > 0 && i.Krugerrand > 0
}

// Complete reports whether the cached row includes every value currently
// displayed by MasjidBoard. Older caches remain valid, but are refreshed once
// so newly introduced fields can be backfilled from the same source row.
func (i Indicators) Complete() bool {
	return i.RandDollar > 0 && i.Gold24Carat > 0 && i.Gold22Carat > 0 &&
		i.Gold18Carat > 0 && i.Gold14Carat > 0 && i.Gold9Carat > 0 &&
		i.Silver > 0 && i.Nisaab > 0 && i.MinimumMahr > 0 &&
		i.MahrFaatimi > 0 && i.Krugerrand > 0
}
