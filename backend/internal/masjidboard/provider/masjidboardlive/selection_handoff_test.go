package masjidboardlive

import (
	"net/http"
	"testing"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

func selectedBoard(id, name string, offset int64) selection.Board {
	return selection.Board{
		CatalogueID:      providerName + ":" + id,
		Provider:         providerName,
		ExternalID:       id,
		Name:             name,
		TimeZoneOffsetMS: offset,
	}
}

func TestNewCoreClientFromSelection(t *testing.T) {
	board := selectedBoard("brits-jamia", "Brits Jamia Masjid", 7200000)
	client, err := NewCoreClientFromSelection(board)
	if err != nil {
		t.Fatalf("NewCoreClientFromSelection() error = %v", err)
	}
	if client.WebURL != "brits-jamia" || client.Identity.ID != "brits-jamia" || client.Identity.Name != "Brits Jamia Masjid" {
		t.Fatalf("client = %+v", client)
	}
	if client.Identity.TimeZone != "GMT+02:00" {
		t.Fatalf("timezone = %q, want GMT+02:00", client.Identity.TimeZone)
	}
}

func TestNewCoreClientFromSelectionPreservesFractionalOffset(t *testing.T) {
	board := selectedBoard("test-board", "Test Masjid", 19800000)
	client, err := NewCoreClientFromSelection(board)
	if err != nil {
		t.Fatalf("NewCoreClientFromSelection() error = %v", err)
	}
	if client.Identity.TimeZone != "GMT+05:30" {
		t.Fatalf("timezone = %q, want GMT+05:30", client.Identity.TimeZone)
	}
}

func TestNewCoreClientsFromSelectionPreservesOrder(t *testing.T) {
	state := selection.State{Boards: []selection.Board{
		selectedBoard("brits-taqwa", "Masjid Taqwa", 7200000),
		selectedBoard("brits-jamia", "Brits Jamia Masjid", 7200000),
		selectedBoard("brits-darul-uloom", "Jamiah Yusuf Darul Uloom Brits", 7200000),
	}}
	clients, err := NewCoreClientsFromSelection(state)
	if err != nil {
		t.Fatalf("NewCoreClientsFromSelection() error = %v", err)
	}
	if len(clients) != 3 {
		t.Fatalf("clients = %d, want 3", len(clients))
	}
	want := []string{"brits-taqwa", "brits-jamia", "brits-darul-uloom"}
	for i, id := range want {
		if clients[i].WebURL != id {
			t.Fatalf("clients[%d].WebURL = %q, want %q", i, clients[i].WebURL, id)
		}
	}
}

func TestNewCoreClientsFromSelectionUnconfigured(t *testing.T) {
	clients, err := NewCoreClientsFromSelection(selection.State{})
	if err != nil {
		t.Fatalf("NewCoreClientsFromSelection() error = %v", err)
	}
	if len(clients) != 0 {
		t.Fatalf("clients = %d, want 0", len(clients))
	}
}

func TestNewCoreClientsFromSelectionRejectsInvalidConfiguredState(t *testing.T) {
	state := selection.State{Boards: []selection.Board{
		selectedBoard("brits-jamia", "Brits Jamia Masjid", 7200000),
		selectedBoard("brits-jamia", "Brits Jamia Masjid", 7200000),
	}}
	if _, err := NewCoreClientsFromSelection(state); err == nil {
		t.Fatal("NewCoreClientsFromSelection() expected validation error")
	}
}

func TestNewCoreClientFromSelectionRejectsWrongProvider(t *testing.T) {
	board := selectedBoard("brits-jamia", "Brits Jamia Masjid", 7200000)
	board.Provider = "other"
	if _, err := NewCoreClientFromSelection(board); err == nil {
		t.Fatal("NewCoreClientFromSelection() expected provider error")
	}
}

func TestNewCoreClientsFromSelectionInjectsHTTPClient(t *testing.T) {
	httpClient := &http.Client{}
	state := selection.State{Boards: []selection.Board{
		selectedBoard("brits-jamia", "Brits Jamia Masjid", 7200000),
		selectedBoard("brits-taqwa", "Masjid Taqwa", 7200000),
	}}
	clients, err := NewCoreClientsFromSelectionWithHTTPClient(state, httpClient)
	if err != nil {
		t.Fatalf("NewCoreClientsFromSelectionWithHTTPClient() error = %v", err)
	}
	for i, client := range clients {
		if client.HTTPClient != httpClient {
			t.Fatalf("clients[%d] HTTPClient was not injected", i)
		}
	}
}
