package catalogue

import (
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/X-Calibre/MasjidPi/backend/internal/stream"
)

func ParseHTML(r io.Reader) ([]stream.Stream, error) {
	doc, err := goquery.NewDocumentFromReader(r)
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
