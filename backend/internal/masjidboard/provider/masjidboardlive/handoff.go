package masjidboardlive

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
)

const millisecondsPerMinute int64 = 60 * 1000

// BoardIdentity converts a discovery catalogue entry into the provider-
// independent identity needed by the individual Core board provider.
func (e CatalogueEntry) BoardIdentity() (model.BoardIdentity, error) {
	name := strings.TrimSpace(e.Name)
	webURL := strings.TrimSpace(e.WebURL)
	if name == "" {
		return model.BoardIdentity{}, fmt.Errorf("masjidboardlive: catalogue entry name is required")
	}
	if webURL == "" {
		return model.BoardIdentity{}, fmt.Errorf("masjidboardlive: catalogue entry web_url is required")
	}

	timezone, err := formatGMTOffset(e.TimeZoneOffsetMS)
	if err != nil {
		return model.BoardIdentity{}, err
	}

	return model.BoardIdentity{
		ID:       webURL,
		Name:     name,
		TimeZone: timezone,
	}, nil
}

// NewCoreClientFromCatalogue creates an individual-board Core provider from a
// FindMasjid catalogue entry. The public web_url remains the external board ID.
func NewCoreClientFromCatalogue(entry CatalogueEntry) (CoreClient, error) {
	return NewCoreClientFromCatalogueWithHTTPClient(entry, nil)
}

// NewCoreClientFromCatalogueWithHTTPClient is the testable form of
// NewCoreClientFromCatalogue and allows callers to inject an HTTP client.
func NewCoreClientFromCatalogueWithHTTPClient(entry CatalogueEntry, client *http.Client) (CoreClient, error) {
	identity, err := entry.BoardIdentity()
	if err != nil {
		return CoreClient{}, err
	}
	return CoreClient{
		HTTPClient: client,
		WebURL:     strings.TrimSpace(entry.WebURL),
		Identity:   identity,
	}, nil
}

func formatGMTOffset(offsetMS int64) (string, error) {
	if offsetMS%millisecondsPerMinute != 0 {
		return "", fmt.Errorf("masjidboardlive: timezone offset %dms is not minute-aligned", offsetMS)
	}

	totalMinutes := offsetMS / millisecondsPerMinute
	if totalMinutes < -24*60 || totalMinutes > 24*60 {
		return "", fmt.Errorf("masjidboardlive: timezone offset %dms is outside supported range", offsetMS)
	}

	sign := "+"
	if totalMinutes < 0 {
		sign = "-"
		totalMinutes = -totalMinutes
	}

	hours := totalMinutes / 60
	minutes := totalMinutes % 60
	return fmt.Sprintf("GMT%s%02d:%02d", sign, hours, minutes), nil
}
