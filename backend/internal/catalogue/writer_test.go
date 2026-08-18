package catalogue

import (
	"os"
	"testing"

	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
)

func TestWriteCatalogueSkipsUnchangedContent(t *testing.T) {
	path := t.TempDir() + "/catalogue.json"
	streams := []stream.Stream{{ID: "one", Name: "One", URL: "https://example.com/one"}}

	if err := WriteCatalogue(path, streams); err != nil {
		t.Fatalf("initial write: %v", err)
	}

	before, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat before unchanged write: %v", err)
	}

	if err := WriteCatalogue(path, streams); err != nil {
		t.Fatalf("unchanged write: %v", err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat after unchanged write: %v", err)
	}
	if !os.SameFile(before, after) {
		t.Fatal("unchanged catalogue replaced the state file")
	}
}

func TestWriteCatalogueReplacesChangedContentAtomically(t *testing.T) {
	path := t.TempDir() + "/catalogue.json"
	first := []stream.Stream{{ID: "one", Name: "One", URL: "https://example.com/one"}}
	second := []stream.Stream{{ID: "two", Name: "Two", URL: "https://example.com/two"}}

	if err := WriteCatalogue(path, first); err != nil {
		t.Fatalf("initial write: %v", err)
	}
	if err := WriteCatalogue(path, second); err != nil {
		t.Fatalf("changed write: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read changed catalogue: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("changed catalogue is empty")
	}
}
