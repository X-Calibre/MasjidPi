package radio

import (
	"testing"

	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
)

func TestCatalogueContainsOnlyRadioStreams(t *testing.T) {
	stations := Catalogue()
	if len(stations) != 8 {
		t.Fatalf("station count = %d, want 8", len(stations))
	}

	seen := make(map[string]bool, len(stations))
	for _, station := range stations {
		if station.SourceKind() != stream.KindRadio {
			t.Fatalf("station %q kind = %q, want radio", station.ID, station.SourceKind())
		}
		if station.ID == "" || station.Name == "" || station.URL == "" {
			t.Fatalf("incomplete station: %#v", station)
		}
		if seen[station.ID] {
			t.Fatalf("duplicate station ID %q", station.ID)
		}
		seen[station.ID] = true
	}
}

func TestMergePreservesMasjidsAndAppendsRadio(t *testing.T) {
	masjids := []stream.Stream{{ID: "masjid-a", Name: "Masjid A", URL: "https://example.test/a"}}
	merged := Merge(masjids)

	if len(merged) != 9 {
		t.Fatalf("merged count = %d, want 9", len(merged))
	}
	if merged[0].ID != "masjid-a" || merged[0].SourceKind() != stream.KindMasjid {
		t.Fatalf("first stream = %#v, want original masjid", merged[0])
	}
	for _, station := range merged[1:] {
		if station.SourceKind() != stream.KindRadio {
			t.Fatalf("merged station %q kind = %q, want radio", station.ID, station.SourceKind())
		}
	}
}
