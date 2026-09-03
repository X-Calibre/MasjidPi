package display

import (
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/dailycontent"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

func displayDailyContent() *dailycontent.Content {
	return &dailycontent.Content{
		Ayah:     dailycontent.Ayah{Surah: "Surah 1", AyahNumber: "Ayah 1", Text: "Ayah"},
		Hadith:   dailycontent.Hadith{Heading: "Hadith", Text: "Hadith", Reference: "Reference"},
		Sunnah:   dailycontent.Sunnah{Heading: "Sunnah", Text: "Sunnah", Reference: "Reference"},
		Language: "en", Source: dailycontent.SourceName, SourceURL: dailycontent.SourceURL,
		ContentDate: "2026-09-03", FetchedAt: time.Date(2026, 9, 3, 12, 0, 0, 0, time.UTC),
	}
}

func TestPresentDailyIslamicContentDefaultsAllCategoriesEnabled(t *testing.T) {
	presented := PresentDailyIslamicContent(displayDailyContent(), selection.State{})
	if presented == nil || presented.Ayah == nil || presented.Hadith == nil || presented.Sunnah == nil {
		t.Fatalf("presented=%+v", presented)
	}
}

func TestPresentDailyIslamicContentFiltersDisabledCategories(t *testing.T) {
	falseValue := false
	state := selection.State{ShowDailyHadith: &falseValue, ShowDailySunnah: &falseValue}
	presented := PresentDailyIslamicContent(displayDailyContent(), state)
	if presented == nil || presented.Ayah == nil || presented.Hadith != nil || presented.Sunnah != nil {
		t.Fatalf("presented=%+v", presented)
	}
	state.ShowDailyAyah = &falseValue
	if got := PresentDailyIslamicContent(displayDailyContent(), state); got != nil {
		t.Fatalf("all-disabled content=%+v", got)
	}
}
