package radio

import "github.com/X-Calibre/MasjidPi/backend/internal/stream"

// Catalogue returns the built-in South African Islamic radio stations that
// have been validated with MasjidPi's mpv stack. Radio entries are maintained
// by MasjidPi rather than scraped at runtime so endpoint changes can be tested
// before being shipped to appliances.
func Catalogue() []stream.Stream {
	return []stream.Stream{
		{
			ID:       "radio-cii",
			Kind:     stream.KindRadio,
			Name:     "Channel Islam International",
			Location: "Gauteng",
			URL:      "https://edge.iono.fm/xice/109_medium.aac",
		},
		{
			ID:       "radio-radio-islam",
			Kind:     stream.KindRadio,
			Name:     "Radio Islam International",
			Location: "Gauteng",
			URL:      "https://cast1.my-control-panel.com/proxy/netmoham/radioislam.mp3",
		},
		{
			ID:       "radio-786",
			Kind:     stream.KindRadio,
			Name:     "Radio 786",
			Location: "Western Cape",
			URL:      "https://stream.krypton.co.za/proxy/radio786?mp=/stream",
		},
		{
			ID:       "radio-voc",
			Kind:     stream.KindRadio,
			Name:     "Voice of the Cape",
			Location: "Western Cape",
			URL:      "https://streaming.fabrik.fm/voc/echocast/audio/low/index.m3u8",
		},
		{
			ID:       "radio-al-ansaar",
			Kind:     stream.KindRadio,
			Name:     "Radio Al Ansaar",
			Location: "KwaZulu-Natal",
			URL:      "https://edge.iono.fm/xice/467_medium.aac",
		},
		{
			ID:       "radio-markaz-sahaba",
			Kind:     stream.KindRadio,
			Name:     "Markaz Sahaba Online Radio",
			Location: "KwaZulu-Natal",
			URL:      "http://zas4.ndx.co.za:9088/stream",
		},
		{
			ID:       "radio-sirius-fm",
			Kind:     stream.KindRadio,
			Name:     "Sirius FM 105.7",
			Location: "Gauteng",
			URL:      "http://s8.voscast.com:7112/;",
		},
		{
			ID:       "radio-salaamedia",
			Kind:     stream.KindRadio,
			Name:     "Salaamedia",
			Location: "Gauteng",
			URL:      "http://capeant.antfarm.co.za:1935/salaam/salaam.stream/playlist.m3u8",
		},
	}
}

// Merge appends the built-in radio catalogue to a set of masjid streams.
// Copies are returned so callers can safely replace their runtime store.
func Merge(masjids []stream.Stream) []stream.Stream {
	radios := Catalogue()
	merged := make([]stream.Stream, 0, len(masjids)+len(radios))
	merged = append(merged, masjids...)
	merged = append(merged, radios...)
	return merged
}
