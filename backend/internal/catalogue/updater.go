package catalogue

const (
	PageURL       = "https://www.livemasjid.com"
	PageFile      = "data/page.html"
	CatalogueFile = "data/catalogue.json"
)

// Update downloads the latest LiveMasjid catalogue,
// parses it and writes catalogue.json.

func Update() error {

	client := NewClient()

	if err := client.Download(PageURL, PageFile); err != nil {
		return err
	}

	streams, err := ParseHTML(PageFile)
	if err != nil {
		return err
	}

	if err := WriteCatalogue(CatalogueFile, streams); err != nil {
		return err
	}

	return nil
}
