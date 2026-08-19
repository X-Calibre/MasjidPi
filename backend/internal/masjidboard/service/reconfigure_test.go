package service

import (
	"context"
	"testing"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/provider"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

func TestServiceReconfigureReplacesRuntimeSelection(t *testing.T) {
	initial := selection.State{Boards: []selection.Board{serviceBoard("one", "One")}}
	service, err := newWithFactory(initial, newMemoryCache(), func(board selection.Board) (provider.Provider, error) {
		return fakeProvider{board: liveBoard(board.ExternalID, board.Name)}, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Unit-constructed services do not have a production selection store, so
	// inject a temporary one to exercise the same persistence path as New().
	service.selectionStore = selection.NewStore(t.TempDir() + "/selection.json")

	next := selection.State{Boards: []selection.Board{
		serviceBoard("two", "Two"),
		serviceBoard("three", "Three"),
	}}
	if err := service.Reconfigure(next); err != nil {
		t.Fatalf("Reconfigure() error = %v", err)
	}
	if got := service.Selection(); len(got.Boards) != 2 || got.Boards[0].ExternalID != "two" || got.Boards[1].ExternalID != "three" {
		t.Fatalf("Selection() = %+v", got)
	}

	results := service.Refresh(context.Background())
	if len(results) != 2 || results[0].Selection.ExternalID != "two" || results[1].Selection.ExternalID != "three" {
		t.Fatalf("Refresh() results = %+v", results)
	}

	persisted, err := service.selectionStore.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(persisted.Boards) != 2 || persisted.Boards[0].ExternalID != "two" {
		t.Fatalf("persisted = %+v", persisted)
	}
}

func TestServiceReconfigureRejectsInvalidSelectionWithoutChangingRuntime(t *testing.T) {
	initial := selection.State{Boards: []selection.Board{serviceBoard("one", "One")}}
	service, err := newWithFactory(initial, newMemoryCache(), func(board selection.Board) (provider.Provider, error) {
		return fakeProvider{board: liveBoard(board.ExternalID, board.Name)}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	service.selectionStore = selection.NewStore(t.TempDir() + "/selection.json")

	if err := service.Reconfigure(selection.State{}); err == nil {
		t.Fatal("Reconfigure() expected validation error")
	}
	if got := service.Selection(); len(got.Boards) != 1 || got.Boards[0].ExternalID != "one" {
		t.Fatalf("selection changed after failed reconfigure: %+v", got)
	}
}
