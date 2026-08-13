package masjidboardlive

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const defaultEndpoint = "https://api.masjidboardlive.com/mblapi"

// Client retrieves raw MasjidBoard Live responses. Parsing and normalisation
// are intentionally kept separate from transport concerns.
type Client struct {
	HTTPClient *http.Client
	Endpoint   string
	BoardID    string
}

func (c Client) Fetch(ctx context.Context) ([]json.RawMessage, error) {
	if c.BoardID == "" {
		return nil, fmt.Errorf("masjidboardlive: board ID is required")
	}

	endpoint := c.Endpoint
	if endpoint == "" {
		endpoint = defaultEndpoint
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("masjidboardlive: parse endpoint: %w", err)
	}
	q := u.Query()
	q.Set("id", c.BoardID)
	u.RawQuery = q.Encode()

	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("masjidboardlive: create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("masjidboardlive: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("masjidboardlive: unexpected HTTP status %s", resp.Status)
	}

	var rows []json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&rows); err != nil {
		return nil, fmt.Errorf("masjidboardlive: decode response: %w", err)
	}

	return rows, nil
}
