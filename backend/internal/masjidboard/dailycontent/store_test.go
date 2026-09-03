package dailycontent

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func validContent() Content {
	return Content{
		Ayah:     Ayah{Surah: "Surah 1", AyahNumber: "Ayah 1", Text: "Text"},
		Hadith:   Hadith{Heading: "Hadith", Text: "Text", Reference: "Reference"},
		Sunnah:   Sunnah{Heading: "Sunnah", Text: "Text", Reference: "Reference"},
		Language: "en", Source: SourceName, SourceURL: SourceURL, ContentDate: "2026-09-03", FetchedAt: time.Now().UTC(),
	}
}

func TestStoreRoundTrip(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "cache", "daily.json")
	store := Store{Path: path}
	want := validContent()
	if err := store.Save(want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	got, err := store.Load()
	if err != nil || got == nil || got.Ayah.Text != want.Ayah.Text || got.ContentDate != want.ContentDate {
		t.Fatalf("Load() = %+v, %v", got, err)
	}
}

func TestStoreMissingAndInvalidCache(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "daily.json")
	store := Store{Path: path}
	if got, err := store.Load(); err != nil || got != nil {
		t.Fatalf("missing Load() = %+v, %v", got, err)
	}
	if err := store.Save(Content{}); err == nil {
		t.Fatal("Save() expected invalid-content error")
	}
	if err := os.WriteFile(path, []byte(`{"source":"broken"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Load(); err == nil {
		t.Fatal("Load() expected invalid-cache error")
	}
}
