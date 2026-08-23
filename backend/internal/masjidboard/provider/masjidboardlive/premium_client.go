package masjidboardlive

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
)

const defaultPremiumEndpoint = "https://premium.masjidboardlive.com/v2/"

// PremiumClient retrieves the richer 29-row board payload exposed by a
// generated MasjidBoard Live Premium page. Mid is the public stable slug. The
// page-supplied opaque board ID is deliberately resolved on every fetch rather
// than persisted as selection state.
type PremiumClient struct {
	HTTPClient *http.Client
	Endpoint   string
	Mid        string
}

// Fetch implements the provider-level board contract.
func (c PremiumClient) Fetch(ctx context.Context) (model.Board, error) {
	return c.FetchAt(ctx, time.Now())
}

// FetchAt retrieves the generated page, extracts its initial theInfo payload,
// and normalises it. The explicit clock keeps parsing tests deterministic.
func (c PremiumClient) FetchAt(ctx context.Context, now time.Time) (model.Board, error) {
	mid := strings.TrimSpace(c.Mid)
	if mid == "" {
		return model.Board{}, fmt.Errorf("masjidboardlive: Premium mid is required")
	}

	endpoint := c.Endpoint
	if endpoint == "" {
		endpoint = defaultPremiumEndpoint
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return model.Board{}, fmt.Errorf("masjidboardlive: parse Premium endpoint: %w", err)
	}
	query := u.Query()
	query.Set("mid", mid)
	u.RawQuery = query.Encode()

	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return model.Board{}, fmt.Errorf("masjidboardlive: create Premium request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return model.Board{}, fmt.Errorf("masjidboardlive: Premium request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return model.Board{}, fmt.Errorf("masjidboardlive: unexpected Premium HTTP status %s", resp.Status)
	}
	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Board{}, fmt.Errorf("masjidboardlive: read Premium response: %w", err)
	}

	boardID, rows, err := ExtractPremiumPageData(html)
	if err != nil {
		return model.Board{}, err
	}
	board, err := Parse(rows, boardID, now)
	if err != nil {
		return model.Board{}, fmt.Errorf("masjidboardlive: parse Premium response: %w", err)
	}
	return board, nil
}
