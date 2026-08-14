package masjidboardlive

import (
	"embed"
	"encoding/json"
	"testing"
	"time"
)

//go:embed testdata/azaadville-darul-uloom-core.json
var capturedFixture embed.FS

func loadCapturedRows(t *testing.T) []json.RawMessage {
	t.Helper()
	data, err := capturedFixture.ReadFile("testdata/azaadville-darul-uloom-core.json")
	if err != nil {
		t.Fatalf("read captured fixture: %v", err)
	}

	var rows []json.RawMessage
	if err := json.Unmarshal(data, &rows); err != nil {
		t.Fatalf("decode captured fixture: %v", err)
	}
	return rows
}

func TestParseCapturedMasjidBoardLiveCore(t *testing.T) {
	rows := loadCapturedRows(t)
	if len(rows) != 29 {
		t.Fatalf("fixture rows = %d, want 29", len(rows))
	}

	now := time.Date(2026, 9, 11, 9, 0, 0, 0, time.UTC)
	board, err := Parse(rows, "1Zpg5LKfd_ZoEQsA0rsyWNBrUgY6QVaHnGdPfuKHF24A", now)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if board.Identity.SourceBoardID != "1Zpg5LKfd_ZoEQsA0rsyWNBrUgY6QVaHnGdPfuKHF24A" {
		t.Fatalf("SourceBoardID = %q", board.Identity.SourceBoardID)
	}
	if board.Identity.MasjidID != "azaadville-darul-uloom" {
		t.Fatalf("MasjidID = %q", board.Identity.MasjidID)
	}
	if board.Identity.EnglishName != "Azaadville" {
		t.Fatalf("EnglishName = %q", board.Identity.EnglishName)
	}
	if board.Identity.Location != "Johannesburg" {
		t.Fatalf("Location = %q", board.Identity.Location)
	}
	if board.DateContext.Timezone != "GMT+02" {
		t.Fatalf("Timezone = %q", board.DateContext.Timezone)
	}
	if got := board.DateContext.GregorianDate.Format("2006-01-02"); got != "2026-09-11" {
		t.Fatalf("GregorianDate = %q", got)
	}

	assertTime := func(name string, got *time.Time, hour, minute int) {
		t.Helper()
		if got == nil {
			t.Fatalf("%s is nil", name)
		}
		if got.Hour() != hour || got.Minute() != minute {
			t.Fatalf("%s = %s, want %02d:%02d", name, got.Format("15:04"), hour, minute)
		}
		_, offset := got.Zone()
		if offset != 2*60*60 {
			t.Fatalf("%s UTC offset = %d, want %d", name, offset, 2*60*60)
		}
	}

	assertTime("Fajr Adhan", board.PrayerTimes.Fajr.Adhan, 5, 30)
	assertTime("Fajr Jamaah", board.PrayerTimes.Fajr.Jamaah, 5, 50)
	assertTime("Zuhr Adhan", board.PrayerTimes.Zuhr.Adhan, 12, 30)
	assertTime("Zuhr Jamaah", board.PrayerTimes.Zuhr.Jamaah, 12, 45)
	assertTime("Asr Adhan", board.PrayerTimes.Asr.Adhan, 16, 15)
	assertTime("Asr Jamaah", board.PrayerTimes.Asr.Jamaah, 16, 30)
	assertTime("Maghrib Adhan", board.PrayerTimes.Maghrib.Adhan, 17, 52)
	if board.PrayerTimes.Maghrib.Jamaah != nil {
		t.Fatalf("Maghrib Jamaah = %v, want nil from captured response", board.PrayerTimes.Maghrib.Jamaah)
	}
	assertTime("Esha Adhan", board.PrayerTimes.Esha.Adhan, 19, 15)
	assertTime("Esha Jamaah", board.PrayerTimes.Esha.Jamaah, 19, 30)
}

func TestParseRejectsIncomplete29RowResponse(t *testing.T) {
	rows := loadCapturedRows(t)[:28]
	if _, err := Parse(rows, "board-id", time.Now()); err == nil {
		t.Fatal("Parse() expected an error for fewer than 29 rows")
	}
}

func TestParseRejectsMissingTimezone(t *testing.T) {
	rows := loadCapturedRows(t)
	rows[rowClock] = json.RawMessage(`[]`)

	if _, err := Parse(rows, "board-id", time.Now()); err == nil {
		t.Fatal("Parse() expected an error for missing timezone")
	}
}

func TestParseAllowsMissingPrayerValue(t *testing.T) {
	rows := loadCapturedRows(t)
	rows[rowSalah] = json.RawMessage(`["05:30","", "12:30","12:45", "16:15","16:30", "17:52","", "19:15","19:30"]`)

	board, err := Parse(rows, "board-id", time.Date(2026, 9, 11, 9, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if board.PrayerTimes.Fajr.Adhan == nil || board.PrayerTimes.Fajr.Jamaah != nil {
		t.Fatal("expected Fajr Adhan only")
	}
	if board.PrayerTimes.Maghrib.Adhan == nil || board.PrayerTimes.Maghrib.Jamaah != nil {
		t.Fatal("expected Maghrib Adhan only")
	}
}
