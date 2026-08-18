package masjidboardlive

import (
	"bytes"
	"testing"
	"time"
)

func TestExtractCoreData(t *testing.T) {
	html := []byte(`<html><body><script>
let data = {
 lang : "en",
 mbl_number : "MBL11517PRP",
 customcode : "if (x) { console.log(\"ok\"); }",
 sunday_zuhr_text : "(Sundays & Public Holidays)"
 }
</script><script src="/boards/script.js"></script></body></html>`)

	got, err := ExtractCoreData(html)
	if err != nil {
		t.Fatalf("ExtractCoreData() error = %v", err)
	}

	want := []byte(`{
 lang : "en",
 mbl_number : "MBL11517PRP",
 customcode : "if (x) { console.log(\"ok\"); }",
 sunday_zuhr_text : "(Sundays & Public Holidays)"
 }`)
	if !bytes.Equal(got, want) {
		t.Fatalf("ExtractCoreData() = %q, want %q", got, want)
	}
}

func TestExtractCoreDataFeedsParser(t *testing.T) {
	fixture := loadCoreFixture(t)
	html := append([]byte(`<html><script>let data = `), fixture...)
	html = append(html, []byte(`</script><script src="/boards/script.js"></script></html>`)...)

	object, err := ExtractCoreData(html)
	if err != nil {
		t.Fatalf("ExtractCoreData() error = %v", err)
	}

	now := time.Date(2026, 8, 18, 16, 0, 0, 0, time.UTC)
	result, err := ParseCoreObject(object, coreIdentity(), now)
	if err != nil {
		t.Fatalf("ParseCoreObject() error = %v", err)
	}
	if result.Metadata.MBLNumber != "MBL11517PRP" {
		t.Fatalf("MBLNumber = %q", result.Metadata.MBLNumber)
	}
	assertCoreClock(t, "Fajr Jamaah", result.Board.PrayerTimes.Fajr.Jamaah, 6, 0)
}

func TestExtractCoreDataRejectsMissingAssignment(t *testing.T) {
	if _, err := ExtractCoreData([]byte(`<html><script>let other = {};</script></html>`)); err == nil {
		t.Fatal("ExtractCoreData() expected an error for missing data assignment")
	}
}

func TestExtractCoreDataRejectsUnterminatedObject(t *testing.T) {
	if _, err := ExtractCoreData([]byte(`<script>let data = { lang : "en"`)); err == nil {
		t.Fatal("ExtractCoreData() expected an error for unterminated object")
	}
}
