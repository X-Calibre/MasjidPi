package catalogue

import (
	"testing"
	"time"
)

func testTime(day int) time.Time {
	return time.Date(2026, 8, day, 12, 0, 0, 0, time.UTC)
}

func testRecord(slug, name string) Record {
	return Record{
		ID:               "masjidboardlive:" + slug,
		Provider:         "masjidboardlive",
		ExternalID:       slug,
		Name:             name,
		City:             "Brits",
		Region:           "North West",
		Country:          "South Africa",
		TimeZoneOffsetMS: 7200000,
		ProviderMetadata: map[string]string{
			"mbl_id": "MBL11517PRP",
		},
	}
}

func catalogueAt(day int, records ...Record) Catalogue {
	return Catalogue{
		RetrievedAt: testTime(day),
		ValidatedAt: testTime(day),
		Records:     records,
	}
}

func TestID(t *testing.T) {
	got, err := ID("masjidboardlive", "brits-jamia")
	if err != nil {
		t.Fatalf("ID() error = %v", err)
	}
	if got != "masjidboardlive:brits-jamia" {
		t.Fatalf("ID() = %q", got)
	}
}

func TestReconcileAddsNewRecord(t *testing.T) {
	candidate := catalogueAt(18, testRecord("brits-jamia", "Brits Jamia Masjid"))

	got, err := Reconcile(Catalogue{}, candidate)
	if err != nil {
		t.Fatalf("Reconcile() error = %v", err)
	}
	if len(got.Records) != 1 {
		t.Fatalf("records = %d, want 1", len(got.Records))
	}
	record := got.Records[0]
	if record.Status != StatusActive {
		t.Fatalf("status = %q", record.Status)
	}
	if !record.DiscoveredAt.Equal(testTime(18)) || !record.LastSeenAt.Equal(testTime(18)) {
		t.Fatalf("timestamps = discovered %v last_seen %v", record.DiscoveredAt, record.LastSeenAt)
	}
}

func TestReconcilePreservesIdentityAcrossRename(t *testing.T) {
	existing := testRecord("brits-jamia", "Old Name")
	existing.Status = StatusActive
	existing.DiscoveredAt = testTime(10)
	existing.LastSeenAt = testTime(10)
	current := catalogueAt(10, existing)

	incoming := testRecord("brits-jamia", "Brits Jamia Masjid")
	incoming.City = "Brits Central"
	incoming.ProviderMetadata["mbl_id"] = "MBL-UPDATED"
	candidate := catalogueAt(18, incoming)

	got, err := Reconcile(current, candidate)
	if err != nil {
		t.Fatalf("Reconcile() error = %v", err)
	}
	record := got.Records[0]
	if record.ID != existing.ID {
		t.Fatalf("ID changed to %q", record.ID)
	}
	if record.Name != "Brits Jamia Masjid" || record.City != "Brits Central" {
		t.Fatalf("mutable metadata not updated: %+v", record)
	}
	if record.ProviderMetadata["mbl_id"] != "MBL-UPDATED" {
		t.Fatalf("provider metadata not updated: %+v", record.ProviderMetadata)
	}
	if !record.DiscoveredAt.Equal(testTime(10)) {
		t.Fatalf("DiscoveredAt = %v, want original", record.DiscoveredAt)
	}
	if !record.LastSeenAt.Equal(testTime(18)) {
		t.Fatalf("LastSeenAt = %v", record.LastSeenAt)
	}
}

func TestReconcileRetainsMissingRecord(t *testing.T) {
	existing := testRecord("brits-jamia", "Brits Jamia Masjid")
	existing.Status = StatusActive
	existing.DiscoveredAt = testTime(10)
	existing.LastSeenAt = testTime(17)

	got, err := Reconcile(catalogueAt(17, existing), catalogueAt(18))
	if err != nil {
		t.Fatalf("Reconcile() error = %v", err)
	}
	if len(got.Records) != 1 {
		t.Fatalf("records = %d, want retained record", len(got.Records))
	}
	record := got.Records[0]
	if record.Status != StatusMissing {
		t.Fatalf("status = %q, want missing", record.Status)
	}
	if !record.LastSeenAt.Equal(testTime(17)) {
		t.Fatalf("LastSeenAt changed to %v", record.LastSeenAt)
	}
}

func TestReconcileTreatsChangedSlugAsNewIdentity(t *testing.T) {
	existing := testRecord("old-slug", "Same Masjid")
	existing.Status = StatusActive
	existing.DiscoveredAt = testTime(10)
	existing.LastSeenAt = testTime(17)

	incoming := testRecord("new-slug", "Same Masjid")
	got, err := Reconcile(catalogueAt(17, existing), catalogueAt(18, incoming))
	if err != nil {
		t.Fatalf("Reconcile() error = %v", err)
	}
	if len(got.Records) != 2 {
		t.Fatalf("records = %d, want 2", len(got.Records))
	}

	byID := map[string]Record{}
	for _, record := range got.Records {
		byID[record.ID] = record
	}
	if byID["masjidboardlive:old-slug"].Status != StatusMissing {
		t.Fatalf("old slug was not retained as missing")
	}
	if byID["masjidboardlive:new-slug"].Status != StatusActive {
		t.Fatalf("new slug was not added as active")
	}
}

func TestReconcileReactivatesMissingRecord(t *testing.T) {
	existing := testRecord("brits-jamia", "Brits Jamia Masjid")
	existing.Status = StatusMissing
	existing.DiscoveredAt = testTime(10)
	existing.LastSeenAt = testTime(15)

	got, err := Reconcile(catalogueAt(17, existing), catalogueAt(18, testRecord("brits-jamia", "Brits Jamia Masjid")))
	if err != nil {
		t.Fatalf("Reconcile() error = %v", err)
	}
	if got.Records[0].Status != StatusActive {
		t.Fatalf("status = %q, want active", got.Records[0].Status)
	}
	if !got.Records[0].DiscoveredAt.Equal(testTime(10)) {
		t.Fatalf("DiscoveredAt was not preserved")
	}
}

func TestReconcileRejectsDuplicateCandidate(t *testing.T) {
	record := testRecord("brits-jamia", "Brits Jamia Masjid")
	if _, err := Reconcile(Catalogue{}, catalogueAt(18, record, record)); err == nil {
		t.Fatal("Reconcile() expected duplicate error")
	}
}

func TestReconcileRejectsInvalidCandidate(t *testing.T) {
	record := testRecord("brits-jamia", "Brits Jamia Masjid")
	record.ID = "wrong:id"
	if _, err := Reconcile(Catalogue{}, catalogueAt(18, record)); err == nil {
		t.Fatal("Reconcile() expected invalid identity error")
	}
}

func TestReconcileRejectsMissingCandidateFreshness(t *testing.T) {
	record := testRecord("brits-jamia", "Brits Jamia Masjid")
	if _, err := Reconcile(Catalogue{}, Catalogue{Records: []Record{record}}); err == nil {
		t.Fatal("Reconcile() expected missing freshness error")
	}
}

func TestReconcileDoesNotAliasProviderMetadata(t *testing.T) {
	incoming := testRecord("brits-jamia", "Brits Jamia Masjid")
	candidate := catalogueAt(18, incoming)
	got, err := Reconcile(Catalogue{}, candidate)
	if err != nil {
		t.Fatalf("Reconcile() error = %v", err)
	}

	incoming.ProviderMetadata["mbl_id"] = "CHANGED"
	if got.Records[0].ProviderMetadata["mbl_id"] != "MBL11517PRP" {
		t.Fatalf("reconciled metadata aliases candidate map")
	}
}
