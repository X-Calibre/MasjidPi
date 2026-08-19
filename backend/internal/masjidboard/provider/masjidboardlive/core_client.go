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

const defaultCoreEndpoint = "https://masjidboardlive.com/boards/"

// CoreClient retrieves and normalises a public MasjidBoard Live Core board.
// The public web_url slug is used as the external identifier; discovery is
// responsible for supplying the remaining normalised board identity metadata.
type CoreClient struct {
	HTTPClient *http.Client
	Endpoint   string
	WebURL     string
	Identity   model.BoardIdentity
}

// Fetch implements the provider-level board contract. Provider-specific Core
// metadata remains available through FetchAt for callers that explicitly need
// it, while the rest of MasjidPi sees only the normalised Board model.
func (c CoreClient) Fetch(ctx context.Context) (model.Board, error) {
	result, err := c.FetchAt(ctx, time.Now())
	if err != nil {
		return model.Board{}, err
	}
	return result.Board, nil
}

// FetchAt retrieves the public Core board page, extracts its embedded data
// object and normalises it into the shared MasjidBoard model. The explicit
// clock keeps parsing/date-context tests deterministic and also returns the
// provider-specific metadata captured from the Core object.
func (c CoreClient) FetchAt(ctx context.Context, now time.Time) (CoreResult, error) {
	webURL := strings.TrimSpace(c.WebURL)
	if webURL == "" {
		return CoreResult{}, fmt.Errorf("masjidboardlive: Core web URL is required")
	}

	endpoint := c.Endpoint
	if endpoint == "" {
		endpoint = defaultCoreEndpoint
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return CoreResult{}, fmt.Errorf("masjidboardlive: parse Core endpoint: %w", err)
	}
	// MasjidBoard Live uses the unusual public form /boards/?<web_url>, not
	// /boards/?<web_url>=. Preserve that verified upstream shape explicitly.
	u.RawQuery = url.QueryEscape(webURL)

	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return CoreResult{}, fmt.Errorf("masjidboardlive: create Core request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return CoreResult{}, fmt.Errorf("masjidboardlive: Core request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return CoreResult{}, fmt.Errorf("masjidboardlive: unexpected Core HTTP status %s", resp.Status)
	}

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return CoreResult{}, fmt.Errorf("masjidboardlive: read Core response: %w", err)
	}

	object, err := ExtractCoreData(html)
	if err != nil {
		return CoreResult{}, fmt.Errorf("masjidboardlive: extract Core response: %w", err)
	}

	result, err := ParseCoreObject(object, c.Identity, now)
	if err != nil {
		return CoreResult{}, fmt.Errorf("masjidboardlive: parse Core response: %w", err)
	}
	applyCoreJumuahHeadingFallback(&result, html)
	return result, nil
}
