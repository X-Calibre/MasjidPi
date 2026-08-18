package catalogue

import "github.com/X-Calibre/MasjidPi/backend/internal/stream"

const LiveMasjidPageURL = "https://www.livemasjid.com"

// Update downloads and parses the latest LiveMasjid catalogue in memory, then
// persists the generated catalogue only when its content has changed.
func Update(catalogueFile string) ([]stream.Stream, error) {
	client := NewClient()

	resp, err := client.Get(LiveMasjidPageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	streams, err := ParseHTML(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := WriteCatalogue(catalogueFile, streams); err != nil {
		return nil, err
	}

	return streams, nil
}
