package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
)

func testEntry() Entry {
	return Entry{
		CatalogueID: "masjidboardlive:brits-jamia",
		SuccessfulAt: time.Date(2026, 8, 18, 18, 42, 0, 0, time.UTC),
		Board: model.Board{
			Identity: model.BoardIdentity{
				ID:       "brits-jamia",
				Name:     "Brits Jamia Masjid",
				TimeZone: "GMT+02:00",
			},
			PrayerTimes: model.PrayerTimes{
				Fajr: model.PrayerTime{Jamaah: &model.ClockTime{Hour: 6, Minute: 0}},
				Asr:  model.PrayerTime{Jamaah: &model.ClockTime{Hour: 17, Minute: 0}},
			},
		},
	}
}

func TestStoreLoadMissingEntry(t *testing.T) {
	store := NewStore(t.TempDir())
	_, found, err := store.Load("masjidboardlive:brits-jamia")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if found {
		t.Fatal("Load() found entry, want no cache")
	}
}

func TestStoreSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)
	want := testEntry()
	if err := store.Save(want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	got, found, err := store.Load(want.CatalogueID)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !found {
		t.Fatal("Load() did not find saved entry")
	}
	if got.CatalogueID != want.CatalogueID || !got.SuccessfulAt.Equal(want.SuccessfulAt) || got.Board.Identity.Name != want.Board.Identity.Name {
		t.Fatalf("Load() = %+v, want %+v", got, want)
	}
	if got.Board.PrayerTimes.Asr.Jamaah == nil || got.Board.PrayerTimes.Asr.Jamaah.Hour != 17 {
		t.Fatalf("Asr cache data not preserved: %+v", got.Board.PrayerTimes.Asr)
	}
}

func TestStoreKeepsBoardsIndependent(t *testing.T) {
	store := NewStore(t.TempDir())
	first := testEntry()
	second := testEntry()
	second.CatalogueID = "masjidboardlive:brits-taqwa"
	second.Board.Identity.ID = "brits-taqwa"
	second.Board.Identity.Name = "Masjid Taqwa"
	second.Board.PrayerTimes.Asr.Jamaah = &model.ClockTime{Hour: 16, Minute: 50}

	if err := store.Save(first); err != nil {
		t.Fatal(err)
	}
	if err := store.Save(second); err != nil {
		t.Fatal(err)
	}
	gotFirst, found, err := store.Load(first.CatalogueID)
	if err != nil || !found {
		t.Fatalf("first Load() found=%v err=%v", found, err)
	}
	gotSecond, found, err := store.Load(second.CatalogueID)
	if err != nil || !found {
		t.Fatalf("second Load() found=%v err=%v", found, err)
	}
	if gotFirst.Board.Identity.ID == gotSecond.Board.Identity.ID {
		t.Fatalf("independent cache entries collapsed: first=%+v second=%+v", gotFirst, gotSecond)
	}
}

func TestStoreUnchangedSaveIsNoOp(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)
	entry := testEntry()
	if err := store.Save(entry); err != nil {
		t.Fatal(err)
	}
	path, err := store.pathFor(entry.CatalogueID)
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(20 * time.Millisecond)
	if err := store.Save(entry); err != nil {
		t.Fatal(err)
	}
	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !after.ModTime().Equal(before.ModTime()) {
		t.Fatalf("unchanged Save() rewrote file: before=%v after=%v", before.ModTime(), after.ModTime())
	}
}

func TestStoreSuccessfulRefreshReplacesCache(t *testing.T) {
	store := NewStore(t.TempDir())
	first := testEntry()
	if err := store.Save(first); err != nil {
		t.Fatal(err)
	}
	second := testEntry()
	second.SuccessfulAt = second.SuccessfulAt.Add(time.Hour)
	second.Board.PrayerTimes.Asr.Jamaah = &model.ClockTime{Hour: 17, Minute: 15}
	if err := store.Save(second); err != nil {
		t.Fatal(err)
	}
	got, found, err := store.Load(second.CatalogueID)
	if err != nil || !found {
		t.Fatalf("Load() found=%v err=%v", found, err)
	}
	if !got.SuccessfulAt.Equal(second.SuccessfulAt) || got.Board.PrayerTimes.Asr.Jamaah.Hour != 17 || got.Board.PrayerTimes.Asr.Jamaah.Minute != 15 {
		t.Fatalf("new successful cache not persisted: %+v", got)
	}
}

func TestStoreFailedCandidateCannotReplaceLastKnownGood(t *testing.T) {
	store := NewStore(t.TempDir())
	good := testEntry()
	if err := store.Save(good); err != nil {
		t.Fatal(err)
	}
	bad := testEntry()
	bad.SuccessfulAt = time.Time{}
	bad.Board.Identity.Name = "Broken candidate"
	if err := store.Save(bad); err == nil {
		t.Fatal("Save() expected validation error")
	}
	got, found, err := store.Load(good.CatalogueID)
	if err != nil || !found {
		t.Fatalf("Load() found=%v err=%v", found, err)
	}
	if got.Board.Identity.Name != good.Board.Identity.Name || !got.SuccessfulAt.Equal(good.SuccessfulAt) {
		t.Fatalf("last-known-good entry was replaced: %+v", got)
	}
}

func TestStoreRejectsMalformedPersistedEntry(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)
	path, err := store.pathFor("masjidboardlive:brits-jamia")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := store.Load("masjidboardlive:brits-jamia"); err == nil {
		t.Fatal("Load() expected malformed JSON error")
	}
}

func TestStoreUsesFilesystemSafeHashedNames(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)
	entry := testEntry()
	if err := store.Save(entry); err != nil {
		t.Fatal(err)
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("files = %d, want 1", len(files))
	}
	if filepath.Ext(files[0].Name()) != ".json" || files[0].Name() == entry.CatalogueID+".json" {
		t.Fatalf("cache filename = %q, want hashed .json name", files[0].Name())
	}
}
