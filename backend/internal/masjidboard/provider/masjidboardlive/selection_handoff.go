package masjidboardlive

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/selection"
)

const providerName = "masjidboardlive"

// NewCoreClientFromSelection constructs one Core provider from persisted
// runtime-critical selection state. It does not require the full catalogue.
func NewCoreClientFromSelection(board selection.Board) (CoreClient, error) {
	return NewCoreClientFromSelectionWithHTTPClient(board, nil)
}

// NewCoreClientFromSelectionWithHTTPClient is the testable form of
// NewCoreClientFromSelection and allows callers to inject an HTTP client.
func NewCoreClientFromSelectionWithHTTPClient(board selection.Board, client *http.Client) (CoreClient, error) {
	if strings.TrimSpace(board.Provider) != providerName {
		return CoreClient{}, fmt.Errorf("masjidboardlive: unsupported selected-board provider %q", board.Provider)
	}
	if strings.TrimSpace(board.ExternalID) == "" {
		return CoreClient{}, fmt.Errorf("masjidboardlive: selected board external ID is required")
	}
	if strings.TrimSpace(board.Name) == "" {
		return CoreClient{}, fmt.Errorf("masjidboardlive: selected board name is required")
	}

	timezone, err := formatGMTOffset(board.TimeZoneOffsetMS)
	if err != nil {
		return CoreClient{}, err
	}

	return CoreClient{
		HTTPClient: client,
		WebURL:     strings.TrimSpace(board.ExternalID),
		Identity: model.BoardIdentity{
			ID:       strings.TrimSpace(board.ExternalID),
			Name:     strings.TrimSpace(board.Name),
			TimeZone: timezone,
		},
	}, nil
}

// NewCoreClientsFromSelection constructs one independent Core provider for
// each configured board while preserving the user's selection order.
// An unconfigured zero-value selection produces no providers and no error.
func NewCoreClientsFromSelection(state selection.State) ([]CoreClient, error) {
	return NewCoreClientsFromSelectionWithHTTPClient(state, nil)
}

// NewCoreClientsFromSelectionWithHTTPClient is the testable form of
// NewCoreClientsFromSelection and applies the same injected HTTP client to all
// constructed providers.
func NewCoreClientsFromSelectionWithHTTPClient(state selection.State, client *http.Client) ([]CoreClient, error) {
	if !state.Configured() {
		return nil, nil
	}
	if err := selection.Validate(state); err != nil {
		return nil, err
	}

	clients := make([]CoreClient, 0, len(state.Boards))
	for i, board := range state.Boards {
		core, err := NewCoreClientFromSelectionWithHTTPClient(board, client)
		if err != nil {
			return nil, fmt.Errorf("masjidboardlive: selected board %d: %w", i+1, err)
		}
		clients = append(clients, core)
	}
	return clients, nil
}
