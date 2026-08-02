package catalogue

import (
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
)

func ParseHTML(filename string) ([]stream.Stream, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, err
	}

	var streams []stream.Stream

	doc.Find(".bs-component").Each(func(i int, card *goquery.Selection) {

		link := card.Find(".masjidname")

		name := clean(link.Text())

		href, _ := link.Attr("href")

		id := strings.TrimSpace(strings.TrimPrefix(href, "/"))

		location := clean(card.Find(".location").Text())

		streams = append(streams, stream.Stream{
			ID:       id,
			Name:     name,
			Location: location,
			URL:      RelayBaseURL + id,
		})

	})

	return streams, nil
}

func clean(text string) string {
	return strings.Join(strings.Fields(text), " ")
}
