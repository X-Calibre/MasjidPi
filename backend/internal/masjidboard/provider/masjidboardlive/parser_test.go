package masjidboardlive

import (
	"encoding/json"
	"testing"
	"time"
)

func capturedRows() []json.RawMessage {
	rows := make([]json.RawMessage, 29)
	for i := range rows {
		rows[i] = json.RawMessage(`[]`)
	}

	// Representative fixture shaped exactly like the captured 29-row
	// mblapi response. The positional indexes are taken from the verified
	// functions_uo_latest.js mapping documented in MASJIDBOARD-LIVE.md.
	rows[rowSalah] = json.RawMessage(`[
		"05:12", "05:30",
		"12:31", "13:00",
		"16:05", "16:30",
		"17:48", "18:00",
		"19:05", "19:30",
		"12:45", "en", "12:45", "", "", "", ""
	]`)
	rows[rowMasjid] = json.RawMessage(`[
		"Test Masjid", "مسجد الاختبار", "https://example.invalid", "Africa/Johannesburg",
		"7200000", "", "", "", "0", "0", "", "", "", "", "", ""
	]`)

	// Preserve the known Jumuah row shape even though optional Jumuah parsing
	// is intentionally not part of this core parser yet.
	rows[rowJumuah] = json.RawMessage(`[
		"Jumuah 1", "13:15", "Jumuah 2", "14:00", "Jumuah 3", "15:00",
		"Test Khateeb", "13:00", "13:15", "", "", ""
	]`)

	return rows
}

func TestParseCaptured29RowShape(t *testing.T) {
	now := time.Date(2026, 8, 14, 9, 0, 0, 0, time.UTC)

	board, err := Parse(capturedRows(), "captured-board-id", now)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if board.Identity.SourceBoardID != "captured-board-id" {
		t.Fatalf("SourceBoardID = %q", board.Identity.SourceBoardID)
	}
	if board.Identity.EnglishName != "Test Masjid" {
		t.Fatalf("EnglishName = %q", board.Identity.EnglishName)
	}
	if board.Identity.ArabicName != "مسجد الاختبار" {
		t.Fatalf("ArabicName = %q", board.Identity.ArabicName)
	}
	if board.DateContext.Timezone != "Africa/Johannesburg" {
		t.Fatalf("Timezone = %q", board.DateContext.Timezone)
	}

	assertTime := func(name string, got *time.Time, hour, minute int) {
		t.Helper()
		if got == nil {
			t.Fatalf("%s is nil", name)
		}
		if got.Hour() != hour || got.Minute() != minute {
			t.Fatalf("%s = %s, want %02d:%02d", name, got.Format("15:04"), hour, minute)
		}
		if got.Location().String() != "Africa/Johannesburg" {
			t.Fatalf("%s location = %s", name, got.Location())
		}
	}

	assertTime("Fajr Adhan", board.PrayerTimes.Fajr.Adhan, 5, 12)
	assertTime("Fajr Jamaah", board.PrayerTimes.Fajr.Jamaah, 5, 30)
	assertTime("Zuhr Adhan", board.PrayerTimes.Zuhr.Adhan, 12, 31)
	assertTime("Zuhr Jamaah", board.PrayerTimes.Zuhr.Jamaah, 13, 0)
	assertTime("Asr Adhan", board.PrayerTimes.Asr.Adhan, 16, 5)
	assertTime("Asr Jamaah", board.PrayerTimes.Asr.Jamaah, 16, 30)
	assertTime("Maghrib Adhan", board.PrayerTimes.Maghrib.Adhan, 17, 48)
	assertTime("Maghrib Jamaah", board.PrayerTimes.Maghrib.Jamaah, 18, 0)
	assertTime("Esha Adhan", board.PrayerTimes.Esha.Adhan, 19, 5)
	assertTime("Esha Jamaah", board.PrayerTimes.Esha.Jamaah, 19, 30)
}

func TestParseRejectsIncomplete29RowResponse(t *testing.T) {
	rows := make([]json.RawMessage, 28)
	if _, err := Parse(rows, "board-id", time.Now()); err == nil {
		t.Fatal("Parse() expected an error for fewer than 29 rows")
	}
}

func TestParseRejectsMissingTimezone(t *testing.T) {
	rows := capturedRows()
	rows[rowMasjid] = json.RawMessage(`[
		"Test Masjid", "", "", ""
	]`)

	if _, err := Parse(rows, "board-id", time.Now()); err == nil {
		t.Fatal("Parse() expected an error for missing timezone")
	}
}

func TestParseAllowsMissingOptionalPrayerValues(t *testing.T) {
	rows := capturedRows()
	rows[rowSalah] = json.RawMessage(`[
		"05:12", "", "12:31", "13:00", "", "16:30", "17:48", "18:00", "19:05", ""
	]`)

	board, err := Parse(rows, "board-id", time.Date(2026, 8, 14, 9, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if board.PrayerTimes.Fajr.Adhan == nil || board.PrayerTimes.Fajr.Jamaah != nil {
		t.Fatal("expected Fajr Adhan only")
	}
	if board.PrayerTimes.Asr.Adhan != nil || board.PrayerTimes.Asr.Jamaah == nil {
		t.Fatal("expected Asr Jamaah only")
	}
}
