package stream

import (
	"net/url"
	"strings"
)

const (
	liveMasjidRelayHost  = "relay.livemasjid.com:8443"
	liveMasjidIcecastURL = "https://icecast.livemasjid.com/"
)

type Stream struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Location     string   `json:"location,omitempty"`
	URL          string   `json:"url"`
	FallbackURLs []string `json:"fallback_urls,omitempty"`
}

// PlaybackURLs returns the primary stream URL followed by its fallbacks. Old
// LiveMasjid catalogue files only contain the relay URL, so derive the Icecast
// equivalent from the mount path when no explicit fallback has been cached.
func (s Stream) PlaybackURLs() []string {
	candidates := make([]string, 0, 1+len(s.FallbackURLs))
	seen := make(map[string]struct{}, 1+len(s.FallbackURLs))
	appendUnique := func(candidate string) {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" {
			return
		}
		if _, ok := seen[candidate]; ok {
			return
		}
		seen[candidate] = struct{}{}
		candidates = append(candidates, candidate)
	}

	appendUnique(s.URL)
	for _, fallback := range s.FallbackURLs {
		appendUnique(fallback)
	}

	if len(s.FallbackURLs) == 0 {
		if primary, err := url.Parse(s.URL); err == nil &&
			primary.Scheme == "https" && primary.Host == liveMasjidRelayHost {
			mount := strings.TrimPrefix(primary.EscapedPath(), "/")
			if mount != "" && !strings.Contains(mount, "/") {
				appendUnique(liveMasjidIcecastURL + mount)
			}
		}
	}

	return candidates
}
