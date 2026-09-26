package updates

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultGitHubReleasesURL = "https://api.github.com/repos/X-Calibre/MasjidPi/releases?per_page=30"
	githubResponseLimit      = 2 << 20
	githubTimeout            = 20 * time.Second
)

type GitHubClient struct {
	HTTPClient  *http.Client
	ReleasesURL string
}

type githubRelease struct {
	TagName     string        `json:"tag_name"`
	Draft       bool          `json:"draft"`
	Prerelease  bool          `json:"prerelease"`
	HTMLURL     string        `json:"html_url"`
	PublishedAt time.Time     `json:"published_at"`
	Assets      []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	State              string `json:"state"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// Latest returns the highest stable release that contains the exact Pi 3
// bundle and detached-signature asset pair.
func (c GitHubClient) Latest(
	ctx context.Context,
) (Release, bool, error) {
	endpoint := strings.TrimSpace(c.ReleasesURL)
	if endpoint == "" {
		endpoint = DefaultGitHubReleasesURL
	}
	if _, err := url.ParseRequestURI(endpoint); err != nil {
		return Release{}, false, fmt.Errorf(
			"updates: invalid GitHub releases URL: %w",
			err,
		)
	}

	client := c.HTTPClient
	if client == nil {
		client = &http.Client{
			Timeout: githubTimeout,
		}
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return Release{}, false, fmt.Errorf(
			"updates: create GitHub request: %w",
			err,
		)
	}
	request.Header.Set(
		"Accept",
		"application/vnd.github+json",
	)
	request.Header.Set(
		"User-Agent",
		"MasjidFrame Update Checker",
	)
	request.Header.Set(
		"X-GitHub-Api-Version",
		"2022-11-28",
	)

	response, err := client.Do(request)
	if err != nil {
		return Release{}, false, fmt.Errorf(
			"updates: fetch GitHub releases: %w",
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Release{}, false, fmt.Errorf(
			"updates: GitHub releases returned %s",
			response.Status,
		)
	}

	body, err := io.ReadAll(
		io.LimitReader(
			response.Body,
			githubResponseLimit+1,
		),
	)
	if err != nil {
		return Release{}, false, fmt.Errorf(
			"updates: read GitHub releases: %w",
			err,
		)
	}
	if len(body) > githubResponseLimit {
		return Release{}, false, fmt.Errorf(
			"updates: GitHub response exceeds %d bytes",
			githubResponseLimit,
		)
	}

	var releases []githubRelease
	if err := json.Unmarshal(body, &releases); err != nil {
		return Release{}, false, fmt.Errorf(
			"updates: decode GitHub releases: %w",
			err,
		)
	}

	var selected Release
	var selectedVersion StableVersion
	found := false

	for _, candidate := range releases {
		if candidate.Draft ||
			candidate.Prerelease ||
			candidate.PublishedAt.IsZero() ||
			!isHTTPSURL(candidate.HTMLURL) {
			continue
		}

		version, err := ParseStableVersion(candidate.TagName)
		if err != nil {
			continue
		}

		bundleURL := ""
		signatureURL := ""

		// Prefer branded assets, while retaining the legacy pair required by
		// MasjidPi v1.6.0 devices during the MasjidFrame naming transition.
		for _, bundlePrefix := range []string{
			"masjidframe-update-",
			"masjidpi-update-",
		} {
			bundleName := fmt.Sprintf(
				"%s%s-pi3.tar.zst",
				bundlePrefix,
				candidate.TagName,
			)
			candidateBundleURL, bundleCount := uploadedAssetURL(
				candidate.Assets,
				bundleName,
			)
			candidateSignatureURL, signatureCount := uploadedAssetURL(
				candidate.Assets,
				bundleName+".minisig",
			)

			if bundleCount == 1 &&
				signatureCount == 1 &&
				isHTTPSURL(candidateBundleURL) &&
				isHTTPSURL(candidateSignatureURL) {
				bundleURL = candidateBundleURL
				signatureURL = candidateSignatureURL
				break
			}
		}

		if bundleURL == "" || signatureURL == "" {
			continue
		}

		if found && version.Compare(selectedVersion) <= 0 {
			continue
		}

		selected = Release{
			Version:      candidate.TagName,
			PublishedAt:  candidate.PublishedAt.UTC(),
			PageURL:      candidate.HTMLURL,
			BundleURL:    bundleURL,
			SignatureURL: signatureURL,
		}
		selectedVersion = version
		found = true
	}

	return selected, found, nil
}

func uploadedAssetURL(
	assets []githubAsset,
	name string,
) (string, int) {
	found := ""
	count := 0

	for _, asset := range assets {
		if asset.Name != name || asset.State != "uploaded" {
			continue
		}
		found = asset.BrowserDownloadURL
		count++
	}

	return found, count
}

func isHTTPSURL(value string) bool {
	parsed, err := url.ParseRequestURI(value)
	return err == nil &&
		parsed.Scheme == "https" &&
		parsed.Host != ""
}
