package hierarchy

import (
	"testing"
	"time"
)

func TestNormalizedMergesDuplicateLocations(t *testing.T) {
	now := time.Date(2026, 8, 18, 20, 0, 0, 0, time.UTC)
	state := State{
		RetrievedAt: now,
		ValidatedAt: now,
		Countries: []Country{
			{
				Name:  " South Africa ",
				Count: 614,
				Regions: []Region{
					{Name: "Limpopo", Count: 25, Cities: []Location{{Name: "Polokwane", Count: 4}}},
					{Name: "", Count: 1, Cities: []Location{{Name: "Unknown Town", Count: 1}}},
				},
			},
			{
				Name:  "South Africa",
				Count: 1,
				Regions: []Region{
					{Name: " Limpopo ", Count: 1, Cities: []Location{{Name: "Polokwane", Count: 1}}},
				},
			},
		},
	}

	got := state.Normalized()
	if len(got.Countries) != 1 || got.Countries[0].Count != 615 {
		t.Fatalf("countries = %+v", got.Countries)
	}
	if len(got.Countries[0].Regions) != 2 {
		t.Fatalf("regions = %+v", got.Countries[0].Regions)
	}
	var limpopo *Region
	for i := range got.Countries[0].Regions {
		if got.Countries[0].Regions[i].Name == "Limpopo" {
			limpopo = &got.Countries[0].Regions[i]
		}
	}
	if limpopo == nil || limpopo.Count != 26 || len(limpopo.Cities) != 1 || limpopo.Cities[0].Count != 5 {
		t.Fatalf("Limpopo = %+v", limpopo)
	}
}

func TestValidateAllowsBlankRegion(t *testing.T) {
	now := time.Now().UTC()
	state := State{
		RetrievedAt: now,
		ValidatedAt: now,
		Countries: []Country{{
			Name:    "South Africa",
			Count:   1,
			Regions: []Region{{Name: "", Count: 1, Cities: []Location{{Name: "Unknown Town", Count: 1}}}},
		}},
	}
	if err := state.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestValidateRejectsMissingCountry(t *testing.T) {
	now := time.Now().UTC()
	state := State{RetrievedAt: now, ValidatedAt: now}
	if err := state.Validate(); err == nil {
		t.Fatal("Validate() expected an error")
	}
}
