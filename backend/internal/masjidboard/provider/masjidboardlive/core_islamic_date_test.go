package masjidboardlive

import (
	"strings"
	"testing"
	"time"
)

func TestParseCoreObjectPopulatesIslamicDateAndMetadata(t *testing.T) {
	// Brits Jamia publishes sunset at 17:51 and has adjustment 0 / forceDate30 Y.
	// Before sunset rollover, 16 July 2026 is 1 Safar 1448 in the exact
	// MasjidBoard Live Kuwaiti-calendar algorithm.
	loc := time.FixedZone("GMT+02", 2*60*60)
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, loc)

	result, err := ParseCoreObject(loadCoreFixture(t), coreIdentity(), now)
	if err != nil {
		t.Fatalf("ParseCoreObject() error = %v", err)
	}
	if result.Board.DateContext.IslamicDate != "1 Safar 1448" {
		t.Fatalf("IslamicDate = %q, want %q", result.Board.DateContext.IslamicDate, "1 Safar 1448")
	}
	if result.Metadata.IslamicDateAdjust != 0 {
		t.Fatalf("IslamicDateAdjust = %d, want 0", result.Metadata.IslamicDateAdjust)
	}
	if result.Metadata.ForceDate30 != "Y" {
		t.Fatalf("ForceDate30 = %q, want Y", result.Metadata.ForceDate30)
	}
}

func TestParseCoreObjectIslamicDateRollsOverAfterSunset(t *testing.T) {
	loc := time.FixedZone("GMT+02", 2*60*60)
	// Upstream refreshes the Islamic date 185 seconds after its published
	// sunset. Brits Jamia fixture sunset is 17:51.
	now := time.Date(2026, 7, 16, 17, 54, 6, 0, loc)

	result, err := ParseCoreObject(loadCoreFixture(t), coreIdentity(), now)
	if err != nil {
		t.Fatalf("ParseCoreObject() error = %v", err)
	}
	if result.Board.DateContext.IslamicDate != "2 Safar 1448" {
		t.Fatalf("IslamicDate after sunset = %q, want %q", result.Board.DateContext.IslamicDate, "2 Safar 1448")
	}
}

func TestParseCoreObjectHonoursIslamicDateAdjust(t *testing.T) {
	raw := strings.Replace(
		string(loadCoreFixture(t)),
		`islamicDateAdjust : "0"`,
		`islamicDateAdjust : "-1"`,
		1,
	)
	loc := time.FixedZone("GMT+02", 2*60*60)
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, loc)

	result, err := ParseCoreObject([]byte(raw), coreIdentity(), now)
	if err != nil {
		t.Fatalf("ParseCoreObject() error = %v", err)
	}
	if result.Board.DateContext.IslamicDate != "30 Muharram 1448" {
		t.Fatalf("adjusted IslamicDate = %q, want %q", result.Board.DateContext.IslamicDate, "30 Muharram 1448")
	}
	if result.Metadata.IslamicDateAdjust != -1 {
		t.Fatalf("IslamicDateAdjust = %d, want -1", result.Metadata.IslamicDateAdjust)
	}
}

func TestParseCoreObjectHonoursForceDate30(t *testing.T) {
	raw := strings.Replace(
		string(loadCoreFixture(t)),
		`forceDate30 : "Y"`,
		`forceDate30 : "N"`,
		1,
	)
	loc := time.FixedZone("GMT+02", 2*60*60)
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, loc)

	result, err := ParseCoreObject([]byte(raw), coreIdentity(), now)
	if err != nil {
		t.Fatalf("ParseCoreObject() error = %v", err)
	}
	if result.Board.DateContext.IslamicDate != "30 Muharram 1448" {
		t.Fatalf("forced IslamicDate = %q, want %q", result.Board.DateContext.IslamicDate, "30 Muharram 1448")
	}
}

func TestParseCoreObjectRejectsInvalidIslamicDateAdjust(t *testing.T) {
	raw := strings.Replace(
		string(loadCoreFixture(t)),
		`islamicDateAdjust : "0"`,
		`islamicDateAdjust : "tomorrow"`,
		1,
	)

	if _, err := ParseCoreObject([]byte(raw), coreIdentity(), time.Now()); err == nil {
		t.Fatal("ParseCoreObject() expected an error for invalid islamicDateAdjust")
	}
}
