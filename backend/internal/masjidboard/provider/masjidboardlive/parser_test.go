package masjidboardlive

import (
	"embed"
	"encoding/json"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
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
	if len(rows) < 29 {
		t.Fatalf("fixture rows = %d, want at least 29", len(rows))
	}

	now := time.Date(2026, 9, 11, 9, 0, 0, 0, time.UTC)
	board, err := Parse(rows, "1Zpg5LKfd_ZoEQsA0rsyWNBrUgY6QVaHnGdPfuKHF24A", now)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if board.Identity.ID != "azaadville-darul-uloom" {
		t.Fatalf("Identity.ID = %q", board.Identity.ID)
	}
	if board.Identity.Name != "Azaadville" {
		t.Fatalf("Identity.Name = %q", board.Identity.Name)
	}
	if board.Identity.AlternateName != "Madrasah Arabia Islamia" {
		t.Fatalf("Identity.AlternateName = %q", board.Identity.AlternateName)
	}
	if board.Identity.TimeZone != "GMT+02" {
		t.Fatalf("Identity.TimeZone = %q", board.Identity.TimeZone)
	}
	if got := board.DateContext.GregorianDate.Format("2006-01-02"); got != "2026-09-11" {
		t.Fatalf("GregorianDate = %q", got)
	}

	assertClock := func(name string, got *model.ClockTime, hour, minute int) {
		t.Helper()
		if got == nil {
			t.Fatalf("%s is nil", name)
		}
		if got.Hour != hour || got.Minute != minute {
			t.Fatalf("%s = %02d:%02d, want %02d:%02d", name, got.Hour, got.Minute, hour, minute)
		}
	}

	assertClock("Fajr Adhan", board.PrayerTimes.Fajr.Adhan, 5, 30)
	assertClock("Fajr Jamaah", board.PrayerTimes.Fajr.Jamaah, 5, 50)
	assertClock("Dhuhr Adhan", board.PrayerTimes.Dhuhr.Adhan, 12, 30)
	assertClock("Dhuhr Jamaah", board.PrayerTimes.Dhuhr.Jamaah, 12, 45)
	assertClock("Asr Adhan", board.PrayerTimes.Asr.Adhan, 16, 15)
	assertClock("Asr Jamaah", board.PrayerTimes.Asr.Jamaah, 16, 30)
	assertClock("Maghrib Adhan", board.PrayerTimes.Maghrib.Adhan, 17, 52)
	if board.PrayerTimes.Maghrib.Jamaah != nil {
		t.Fatalf("Maghrib Jamaah = %+v, want nil", board.PrayerTimes.Maghrib.Jamaah)
	}
	assertClock("Esha Adhan", board.PrayerTimes.Esha.Adhan, 19, 15)
	assertClock("Esha Jamaah", board.PrayerTimes.Esha.Jamaah, 19, 30)

	if len(board.PrayerTimes.Jumuah) != 1 {
		t.Fatalf("Jumuah services = %d, want 1", len(board.PrayerTimes.Jumuah))
	}
	jumuah := board.PrayerTimes.Jumuah[0]
	assertClock("Jumuah Adhan", jumuah.Adhan, 12, 45)
	assertClock("Jumuah Jamaah", jumuah.Jamaah, 12, 55)
	assertClock("Jumuah Alternate Adhan", jumuah.AlternateAdhan, 18, 56)
	assertClock("Jumuah Alternate Jamaah", jumuah.AlternateJamaah, 19, 6)
	if jumuah.Khateeb != "Sunnats after Adhān" {
		t.Fatalf("Jumuah Khateeb = %q", jumuah.Khateeb)
	}
	if len(jumuah.Events) != 3 {
		t.Fatalf("Jumuah events = %d, want 3", len(jumuah.Events))
	}
	if jumuah.Events[0].Code != "1" || jumuah.Events[0].Heading != "Lecture" {
		t.Fatalf("Jumuah event 1 = %+v", jumuah.Events[0])
	}
	assertClock("Jumuah event 1", jumuah.Events[0].Time, 12, 15)
	if jumuah.Events[1].Code != "0" || jumuah.Events[1].Heading != "Adhān" {
		t.Fatalf("Jumuah event 2 = %+v", jumuah.Events[1])
	}
	assertClock("Jumuah event 2", jumuah.Events[1].Time, 12, 45)
	if jumuah.Events[2].Code != "6" || jumuah.Events[2].Heading != "Khutbah" {
		t.Fatalf("Jumuah event 3 = %+v", jumuah.Events[2])
	}
	assertClock("Jumuah event 3", jumuah.Events[2].Time, 12, 55)
	assertClock("Jumuah effective Salaah", jumuah.EffectiveSalaah(), 12, 55)

	if board.Astronomical == nil {
		t.Fatal("Astronomical is nil")
	}
	assertClock("Suhur", board.Astronomical.Suhur, 5, 19)
	assertClock("FajrStart", board.Astronomical.FajrStart, 5, 19)
	assertClock("Sunrise", board.Astronomical.Sunrise, 6, 37)
	assertClock("Ishraaq", board.Astronomical.Ishraaq, 6, 52)
	assertClock("Duha", board.Astronomical.Duha, 9, 25)
	assertClock("AsrShafii", board.Astronomical.AsrShafii, 15, 25)
	assertClock("AsrHanafi", board.Astronomical.AsrHanafi, 16, 13)
	assertClock("Sunset", board.Astronomical.Sunset, 17, 49)
	assertClock("EshaStart", board.Astronomical.EshaStart, 19, 7)
}

func TestParseJumuahDoesNotInferSalaahFromKhutbah(t *testing.T) {
	row := json.RawMessage(`["Adhān","12:25","Sunan","12:55","Khutbah","13:00","Ml M Bhamjee","12:25","","18:35","19:10","0,3,6"]`)
	service, err := parseJumuahRow(row)
	if err != nil {
		t.Fatalf("parseJumuahRow() error = %v", err)
	}
	if service.Jamaah != nil {
		t.Fatal("Jumuah Jamaah should be absent")
	}
	if got := service.EffectiveSalaah(); got != nil {
		t.Fatalf("EffectiveSalaah() = %+v, want nil without explicit Jamaah", got)
	}
	if len(service.Events) != 3 || service.Events[2].Heading != "Khutbah" || service.Events[2].Time == nil || service.Events[2].Time.Hour != 13 {
		t.Fatalf("Jumuah events = %+v", service.Events)
	}
}

func TestParseJumuahAllowsEmptyRow(t *testing.T) {
	row := json.RawMessage(`["","","","","","","","","","","","#N/A"]`)
	service, err := parseJumuahRow(row)
	if err != nil {
		t.Fatalf("parseJumuahRow() error = %v", err)
	}
	if service != nil {
		t.Fatalf("service = %+v, want nil", service)
	}
}

func TestCorePlaceholderIsAbsent(t *testing.T) {
	if !isAbsent("~~~~") {
		t.Fatal("expected ~~~~ to be treated as absent")
	}
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

func TestParseRejectsMissingCorePrayer(t *testing.T) {
	rows := loadCapturedRows(t)
	rows[rowSalah] = json.RawMessage(`["","", "12:30","12:45", "16:15","16:30", "17:52","", "19:15","19:30"]`)
	if _, err := Parse(rows, "board-id", time.Now()); err == nil {
		t.Fatal("Parse() expected an error when a core prayer has no usable time")
	}
}

func TestParseAllowsMissingOptionalPrayerValue(t *testing.T) {
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
