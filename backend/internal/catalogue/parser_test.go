package catalogue

import (
	"strings"
	"testing"
)

func TestParseHTMLCachesRelayAndIcecastURLs(t *testing.T) {
	html := `<div class="bs-component"><a class="masjidname" href="/annoor-relay">An Noor</a><span class="location">South Africa</span></div>`
	streams, err := ParseHTML(strings.NewReader(html))
	if err != nil {
		t.Fatalf("ParseHTML: %v", err)
	}
	if len(streams) != 1 {
		t.Fatalf("stream count = %d, want 1", len(streams))
	}
	if got, want := streams[0].URL, RelayBaseURL+"annoor-relay"; got != want {
		t.Fatalf("primary URL = %q, want %q", got, want)
	}
	if len(streams[0].FallbackURLs) != 1 || streams[0].FallbackURLs[0] != IcecastBaseURL+"annoor-relay" {
		t.Fatalf("fallback URLs = %#v", streams[0].FallbackURLs)
	}
}
