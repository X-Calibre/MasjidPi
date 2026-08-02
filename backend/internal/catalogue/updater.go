package catalogue

func Update() error {

	client := NewClient()

	if err := client.Download(
		"https://www.livemasjid.com",
		"data/page.html",
	); err != nil {
		return err
	}

	streams, err := ParseHTML("data/page.html")
	if err != nil {
		return err
	}

	if err := WriteCatalogue("data/catalogue.json", streams); err != nil {
		return err
	}

	return nil
}
