package catalogue

const LiveMasjidPageURL = "https://www.livemasjid.com"

// Update downloads the latest LiveMasjid catalogue, parses it and writes the
// generated catalogue to the supplied runtime paths.
func Update(pageFile, catalogueFile string) error {
	client := NewClient()

	if err := client.Download(LiveMasjidPageURL, pageFile); err != nil {
		return err
	}

	streams, err := ParseHTML(pageFile)
	if err != nil {
		return err
	}

	if err := WriteCatalogue(catalogueFile, streams); err != nil {
		return err
	}

	return nil
}
