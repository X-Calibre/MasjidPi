package masjidboardlive

import (
	"embed"
	"strings"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
)

//go:embed testdata/core-brits-jamia.js
var coreFixture embed.FS

func loadCoreFixture(t *testing.T) []byte {
	t.Helper()
	data, err := coreFixture.ReadFile("testdata/core-brits-jamia.js")
	if err != nil {
		t.Fatalf("read Core fixture: %v", err)
	}
	return data
}

func coreIdentity() model.BoardIdentity {
	return model.BoardIdentity{
		ID:       "brits-jamia",
		Name:     "Brits Jamia Masjid",
		TimeZone: "GMT+02",
	}
}

func assertCoreClock(t *testing.T, name string, got *model.ClockTime, hour, minute int) {
	t.Helper()
	if got == nil {
		t.Fatalf("%s is nil", name)
	}
	if got.Hour != hour || got.Minute != minute {
		t.Fatalf("%s = %02d:%02d, want %02d:%02d", name, got.Hour, got.Minute, hour, minute)
	}
}

func TestParseCoreObjectBritsJamia(t *testing.T) {
	now := time.Date(2026, 8, 18, 16, 0, 0, 0, time.UTC)
	result, err := ParseCoreObject(loadCoreFixture(t), coreIdentity(), now)
	if err != nil {
		t.Fatalf("ParseCoreObject() error = %v", err)
	}

	board := result.Board
	if board.Identity.ID != "brits-jamia" || board.Identity.Name != "Brits Jamia Masjid" {
		t.Fatalf("Identity = %+v", board.Identity)
	}
	if got := board.DateContext.GregorianDate.Format("2006-01-02"); got != "2026-08-18" {
		t.Fatalf("GregorianDate = %q", got)
	}

	assertCoreClock(t, "Fajr Adhan", board.PrayerTimes.Fajr.Adhan, 5, 40)
	assertCoreClock(t, "Fajr Jamaah", board.PrayerTimes.Fajr.Jamaah, 6, 0)
	assertCoreClock(t, "Dhuhr Adhan", board.PrayerTimes.Dhuhr.Adhan, 13, 0)
	assertCoreClock(t, "Dhuhr Jamaah", board.PrayerTimes.Dhuhr.Jamaah, 13, 20)
	assertCoreClock(t, "Asr Adhan", board.PrayerTimes.Asr.Adhan, 16, 40)
	assertCoreClock(t, "Asr Jamaah", board.PrayerTimes.Asr.Jamaah, 17, 0)
	assertCoreClock(t, "Maghrib Adhan", board.PrayerTimes.Maghrib.Adhan, 17, 54)
	if board.PrayerTimes.Maghrib.Jamaah != nil {
		t.Fatalf("Maghrib Jamaah = %+v, want nil", board.PrayerTimes.Maghrib.Jamaah)
	}
	assertCoreClock(t, "Esha Adhan", board.PrayerTimes.Esha.Adhan, 19, 15)
	assertCoreClock(t, "Esha Jamaah", board.PrayerTimes.Esha.Jamaah, 19, 30)

	if board.Astronomical == nil {
		t.Fatal("Astronomical is nil")
	}
	assertCoreClock(t, "Suhur", board.Astronomical.Suhur, 5, 17)
	assertCoreClock(t, "FajrStart", board.Astronomical.FajrStart, 5, 17)
	assertCoreClock(t, "Sunrise", board.Astronomical.Sunrise, 6, 34)
	assertCoreClock(t, "Ishraaq", board.Astronomical.Ishraaq, 6, 49)
	assertCoreClock(t, "Duha", board.Astronomical.Duha, 9, 23)
	assertCoreClock(t, "IstiwaCaution", board.Astronomical.IstiwaCaution, 12, 10)
	assertCoreClock(t, "Istiwa", board.Astronomical.Istiwa, 12, 13)
	assertCoreClock(t, "ZawaalEnd", board.Astronomical.ZawaalEnd, 12, 16)
	assertCoreClock(t, "AsrShafii", board.Astronomical.AsrShafii, 15, 27)
	assertCoreClock(t, "AsrHanafi", board.Astronomical.AsrHanafi, 16, 15)
	assertCoreClock(t, "Sunset", board.Astronomical.Sunset, 17, 51)
	assertCoreClock(t, "EshaStart", board.Astronomical.EshaStart, 19, 8)

	if len(board.PrayerTimes.Jumuah) != 1 {
		t.Fatalf("Jumuah services = %d, want 1", len(board.PrayerTimes.Jumuah))
	}
	jumuah := board.PrayerTimes.Jumuah[0]
	if len(jumuah.Events) != 3 {
		t.Fatalf("Jumuah events = %d, want 3", len(jumuah.Events))
	}
	if jumuah.Events[0].Code != "0" || jumuah.Events[0].Heading != "Adhan" {
		t.Fatalf("Jumuah event 1 = %+v", jumuah.Events[0])
	}
	if jumuah.Events[1].Code != "1" || jumuah.Events[1].Heading != "Lecture" {
		t.Fatalf("Jumuah event 2 = %+v", jumuah.Events[1])
	}
	if jumuah.Events[2].Code != "6" || jumuah.Events[2].Heading != "Khutbah" {
		t.Fatalf("Jumuah event 3 = %+v", jumuah.Events[2])
	}
	assertCoreClock(t, "Jumuah Adhan", jumuah.Adhan, 12, 25)
	if jumuah.Jamaah != nil {
		t.Fatalf("Jumuah Jamaah = %+v, want nil because Core exposes Khutbah rather than a dedicated Jamaah field", jumuah.Jamaah)
	}
	assertCoreClock(t, "Jumuah effective Salaah", jumuah.EffectiveSalaah(), 13, 0)

	if result.Metadata.MBLNumber != "MBL11517PRP" {
		t.Fatalf("MBLNumber = %q", result.Metadata.MBLNumber)
	}
	if result.Metadata.LastUpdated != "Sun, 22 Mar 2026, 12:47:25" {
		t.Fatalf("LastUpdated = %q", result.Metadata.LastUpdated)
	}
	if result.Metadata.LiveStreamProvider != "SmartBilal" {
		t.Fatalf("LiveStreamProvider = %q", result.Metadata.LiveStreamProvider)
	}
	if result.Metadata.LiveStreamURL != "https://media.smartbilal.com/masjid/britsj" {
		t.Fatalf("LiveStreamURL = %q", result.Metadata.LiveStreamURL)
	}
}

func TestParseCoreObjectAllowsJumuahPlaceholders(t *testing.T) {
	original := string(loadCoreFixture(t))
	raw := original
	raw = strings.Replace(raw, "jumuahTime1 : \"12:25\"", "jumuahTime1 : \"~~~~\"", 1)
	raw = strings.Replace(raw, "jumuahTime2 : \"12:40\"", "jumuahTime2 : \"\"", 1)
	raw = strings.Replace(raw, "jumuahTime3 : \"13:00\"", "jumuahTime3 : \"~~~~\"", 1)
	raw = strings.Replace(raw, "jumuahHeadings : \"0,1,6\"", "jumuahHeadings : \"0,,6\"", 1)

	if raw == original {
		t.Fatal("test fixture replacements did not apply")
	}

	result, err := ParseCoreObject([]byte(raw), coreIdentity(), time.Now())
	if err != nil {
		t.Fatalf("ParseCoreObject() error = %v", err)
	}
	if len(result.Board.PrayerTimes.Jumuah) != 1 {
		t.Fatalf("Jumuah services = %d, want 1", len(result.Board.PrayerTimes.Jumuah))
	}
	service := result.Board.PrayerTimes.Jumuah[0]
	if service.Adhan != nil || service.Jamaah != nil {
		t.Fatalf("Jumuah dedicated times = %+v, want nil placeholders", service)
	}
}

func TestParseCoreObjectRejectsMissingCorePrayer(t *testing.T) {
	original := string(loadCoreFixture(t))
	raw := original
	raw = strings.Replace(raw, "fajrAthan : \"05:40\"", "fajrAthan : \"\"", 1)
	raw = strings.Replace(raw, "fajrJamaah : \"06:00\"", "fajrJamaah : \"~~~~\"", 1)

	if raw == original {
		t.Fatal("test fixture replacements did not apply")
	}

	if _, err := ParseCoreObject([]byte(raw), coreIdentity(), time.Now()); err == nil {
		t.Fatal("ParseCoreObject() expected an error for missing Fajr")
	}
}

func TestParseCoreObjectRejectsMissingIdentity(t *testing.T) {
	identity := coreIdentity()
	identity.ID = ""
	if _, err := ParseCoreObject(loadCoreFixture(t), identity, time.Now()); err == nil {
		t.Fatal("ParseCoreObject() expected an error for missing identity ID")
	}
}
