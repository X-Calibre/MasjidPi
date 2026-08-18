package selection

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/catalogue"
)

func selected(id, name string, offset int64) Board {
	return Board{
		CatalogueID:      "masjidboardlive:" + id,
		Provider:         "masjidboardlive",
		ExternalID:       id,
		Name:             name,
		TimeZoneOffsetMS: offset,
	}
}

func TestValidateAllowsUpToThreeBoards(t *testing.T) {
	state := State{Boards: []Board{
		selected("brits-jamia", "Brits Jamia Masjid", 7200000),
		selected("brits-taqwa", "Masjid Taqwa", 7200000),
		selected("brits-darul-uloom", "Jamiah Yusuf Darul Uloom Brits", 7200000),
	}}
	if err := Validate(state); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestValidateRejectsFourthBoard(t *testing.T) {
	state := State{Boards: []Board{
		selected("a", "A", 0),
		selected("b", "B", 0),
		selected("c", "C", 0),
		selected("d", "D", 0),
	}}
	if err := Validate(state); err == nil {
		t.Fatal("Validate() expected maximum-board error")
	}
}

func TestValidateRejectsDuplicateBoard(t *testing.T) {
	board := selected("brits-jamia", "Brits Jamia Masjid", 7200000)
	if err := Validate(State{Boards: []Board{board, board}}); err == nil {
		t.Fatal("Validate() expected duplicate-board error")
	}
}

func TestValidateRejectsMismatchedCatalogueID(t *testing.T) {
	board := selected("brits-jamia", "Brits Jamia Masjid", 7200000)
	board.CatalogueID = "wrong"
	if err := Validate(State{Boards: []Board{board}}); err == nil {
		t.Fatal("Validate() expected catalogue-ID error")
	}
}

func TestFromCatalogueRecord(t *testing.T) {
	record := catalogue.Record{
		ID:               "masjidboardlive:brits-jamia",
		Provider:         "masjidboardlive",
		ExternalID:       "brits-jamia",
		Name:             "Brits Jamia Masjid",
		TimeZoneOffsetMS: 7200000,
		Status:           catalogue.StatusActive,
	}
	board, err := FromCatalogueRecord(record)
	if err != nil {
		t.Fatalf("FromCatalogueRecord() error = %v", err)
	}
	if board.CatalogueID != record.ID || board.ExternalID != record.ExternalID || board.TimeZoneOffsetMS != record.TimeZoneOffsetMS {
		t.Fatalf("board = %+v", board)
	}
}

func TestStoreMissingSelection(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "selection.json"))
	state, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(state.Boards) != 0 {
		t.Fatalf("Boards = %d, want 0", len(state.Boards))
	}
}

func TestStoreSaveLoadPreservesOrder(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "selection.json")
	store := NewStore(path)
	want := State{Boards: []Board{
		selected("brits-taqwa", "Masjid Taqwa", 7200000),
		selected("brits-jamia", "Brits Jamia Masjid", 7200000),
	}}
	if err := store.Save(want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	got, err := NewStore(path).Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(got.Boards) != 2 || got.Boards[0].ExternalID != "brits-taqwa" || got.Boards[1].ExternalID != "brits-jamia" {
		t.Fatalf("order not preserved: %+v", got.Boards)
	}
}

func TestStoreUnchangedSaveIsNoOp(t *testing.T) {
	path := filepath.Join(t.TempDir(), "selection.json")
	store := NewStore(path)
	state := State{Boards: []Board{selected("brits-jamia", "Brits Jamia Masjid", 7200000)}}
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}
	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(20 * time.Millisecond)
	if err := store.Save(state); err != nil {
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

func TestStoreCachesSelectionForRuntime(t *testing.T) {
	path := filepath.Join(t.TempDir(), "selection.json")
	store := NewStore(path)
	state := State{Boards: []Board{selected("brits-jamia", "Brits Jamia Masjid", 7200000)}}
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("broken"), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatalf("cached Load() error = %v", err)
	}
	if len(got.Boards) != 1 || got.Boards[0].ExternalID != "brits-jamia" {
		t.Fatalf("cached state = %+v", got)
	}
}

func TestStoreRejectsInvalidPersistedState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "selection.json")
	if err := os.WriteFile(path, []byte(`{"boards":[{"catalogue_id":"wrong","provider":"masjidboardlive","external_id":"brits-jamia","name":"Brits Jamia Masjid","time_zone_offset_ms":7200000}]}`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := NewStore(path).Load(); err == nil {
		t.Fatal("Load() expected validation error")
	}
}

func TestStoreReturnsDefensiveCopies(t *testing.T) {
	path := filepath.Join(t.TempDir(), "selection.json")
	store := NewStore(path)
	state := State{Boards: []Board{selected("brits-jamia", "Brits Jamia Masjid", 7200000)}}
	if err := store.Save(state); err != nil {
		t.Fatal(err)
	}
	first, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	first.Boards[0].Name = "Mutated"
	second, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if second.Boards[0].Name != "Brits Jamia Masjid" {
		t.Fatalf("cached state mutated through returned copy: %+v", second)
	}
}
