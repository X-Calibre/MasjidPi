package masjidboardlive

import "testing"

func TestCatalogueEntryBoardIdentity(t *testing.T) {
	entry := CatalogueEntry{
		Name:             "Brits Jamia Masjid",
		WebURL:           "brits-jamia",
		TimeZoneOffsetMS: 7200000,
	}

	identity, err := entry.BoardIdentity()
	if err != nil {
		t.Fatalf("BoardIdentity() error = %v", err)
	}
	if identity.ID != "brits-jamia" || identity.Name != "Brits Jamia Masjid" || identity.TimeZone != "GMT+02:00" {
		t.Fatalf("identity = %+v", identity)
	}
}

func TestCatalogueEntryBoardIdentityFractionalPositiveOffset(t *testing.T) {
	entry := CatalogueEntry{
		Name:             "Test Masjid",
		WebURL:           "test-masjid",
		TimeZoneOffsetMS: (5*60 + 30) * 60 * 1000,
	}

	identity, err := entry.BoardIdentity()
	if err != nil {
		t.Fatalf("BoardIdentity() error = %v", err)
	}
	if identity.TimeZone != "GMT+05:30" {
		t.Fatalf("TimeZone = %q, want GMT+05:30", identity.TimeZone)
	}
}

func TestCatalogueEntryBoardIdentityFractionalNegativeOffset(t *testing.T) {
	entry := CatalogueEntry{
		Name:             "Test Masjid",
		WebURL:           "test-masjid",
		TimeZoneOffsetMS: -(3*60 + 30) * 60 * 1000,
	}

	identity, err := entry.BoardIdentity()
	if err != nil {
		t.Fatalf("BoardIdentity() error = %v", err)
	}
	if identity.TimeZone != "GMT-03:30" {
		t.Fatalf("TimeZone = %q, want GMT-03:30", identity.TimeZone)
	}
}

func TestCatalogueEntryBoardIdentityRejectsMissingFields(t *testing.T) {
	for _, entry := range []CatalogueEntry{
		{WebURL: "test"},
		{Name: "Test Masjid"},
	} {
		if _, err := entry.BoardIdentity(); err == nil {
			t.Fatalf("BoardIdentity() expected an error for %+v", entry)
		}
	}
}

func TestFormatGMTOffsetRejectsNonMinuteOffset(t *testing.T) {
	if _, err := formatGMTOffset(1001); err == nil {
		t.Fatal("formatGMTOffset() expected an error for non-minute offset")
	}
}

func TestNewCoreClientFromCatalogue(t *testing.T) {
	entry := CatalogueEntry{
		Name:             "Brits Jamia Masjid",
		WebURL:           "brits-jamia",
		TimeZoneOffsetMS: 7200000,
	}

	client, err := NewCoreClientFromCatalogue(entry)
	if err != nil {
		t.Fatalf("NewCoreClientFromCatalogue() error = %v", err)
	}
	if client.WebURL != "brits-jamia" {
		t.Fatalf("WebURL = %q", client.WebURL)
	}
	if client.Identity.TimeZone != "GMT+02:00" {
		t.Fatalf("TimeZone = %q", client.Identity.TimeZone)
	}
}
