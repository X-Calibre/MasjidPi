package provider

import (
	"context"

	"github.com/X-Calibre/MasjidFrame/backend/internal/masjidboard/model"
)

// Provider retrieves and normalises board data from an external source.
// Implementations must return the semantic model and must not expose
// source-specific response structures to the rest of the application.
type Provider interface {
	Fetch(ctx context.Context) (model.Board, error)
}
