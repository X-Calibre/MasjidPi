package stream

import (
	"reflect"
	"testing"
)

func TestPlaybackURLsUsesExplicitFallbacks(t *testing.T) {
	item := Stream{URL: "https://primary.example/one", FallbackURLs: []string{"https://fallback.example/one", "https://fallback.example/one"}}
	want := []string{"https://primary.example/one", "https://fallback.example/one"}
	if got := item.PlaybackURLs(); !reflect.DeepEqual(got, want) {
		t.Fatalf("PlaybackURLs() = %#v, want %#v", got, want)
	}
}

func TestPlaybackURLsDerivesLiveMasjidFallbackForOldCatalogue(t *testing.T) {
	item := Stream{URL: "https://relay.livemasjid.com:8443/annoor-relay"}
	want := []string{"https://relay.livemasjid.com:8443/annoor-relay", "https://icecast.livemasjid.com/annoor-relay"}
	if got := item.PlaybackURLs(); !reflect.DeepEqual(got, want) {
		t.Fatalf("PlaybackURLs() = %#v, want %#v", got, want)
	}
}

func TestPlaybackURLsDoesNotDeriveFallbackForOtherHosts(t *testing.T) {
	item := Stream{URL: "https://example.com/annoor-relay"}
	want := []string{"https://example.com/annoor-relay"}
	if got := item.PlaybackURLs(); !reflect.DeepEqual(got, want) {
		t.Fatalf("PlaybackURLs() = %#v, want %#v", got, want)
	}
}
