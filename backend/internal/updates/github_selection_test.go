package updates

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func githubTestRelease(
	tag string,
	draft bool,
	prerelease bool,
	includeBundle bool,
	includeSignature bool,
) githubRelease {
	bundleName := "masjidframe-update-" + tag + "-pi3.tar.zst"
	assets := []githubAsset{}

	if includeBundle {
		assets = append(assets, githubAsset{
			Name:  bundleName,
			State: "uploaded",
			BrowserDownloadURL: "https://downloads.example/" +
				bundleName,
		})
	}
	if includeSignature {
		assets = append(assets, githubAsset{
			Name:  bundleName + ".minisig",
			State: "uploaded",
			BrowserDownloadURL: "https://downloads.example/" +
				bundleName + ".minisig",
		})
	}

	return githubRelease{
		TagName:    tag,
		Draft:      draft,
		Prerelease: prerelease,
		HTMLURL: "https://github.example/releases/tag/" +
			tag,
		PublishedAt: time.Date(
			2026,
			time.September,
			16,
			6,
			0,
			0,
			0,
			time.UTC,
		),
		Assets: assets,
	}
}

func githubTestServer(
	t *testing.T,
	releases []githubRelease,
) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(
		func(
			response http.ResponseWriter,
			request *http.Request,
		) {
			if request.Method != http.MethodGet {
				t.Errorf(
					"method = %s, want GET",
					request.Method,
				)
			}
			if got := request.Header.Get("User-Agent"); got !=
				"MasjidFrame Update Checker" {
				t.Errorf(
					"User-Agent = %q",
					got,
				)
			}
			if got := request.Header.Get("Accept"); got !=
				"application/vnd.github+json" {
				t.Errorf(
					"Accept = %q",
					got,
				)
			}

			response.Header().Set(
				"Content-Type",
				"application/json",
			)
			if err := json.NewEncoder(response).Encode(
				releases,
			); err != nil {
				t.Errorf("encode response: %v", err)
			}
		},
	))
}

func TestGitHubClientSelectsHighestCompleteStableRelease(
	t *testing.T,
) {
	releases := []githubRelease{
		githubTestRelease(
			"v1.9.9",
			false,
			false,
			true,
			true,
		),
		githubTestRelease(
			"v2.0.0",
			false,
			false,
			true,
			false,
		),
		githubTestRelease(
			"v9.0.0",
			true,
			false,
			true,
			true,
		),
		githubTestRelease(
			"v8.0.0",
			false,
			true,
			true,
			true,
		),
		githubTestRelease(
			"v7.0.0-rc.1",
			false,
			false,
			true,
			true,
		),
		githubTestRelease(
			"v1.10.0",
			false,
			false,
			true,
			true,
		),
	}

	server := githubTestServer(t, releases)
	defer server.Close()

	client := GitHubClient{
		HTTPClient:  server.Client(),
		ReleasesURL: server.URL,
	}

	release, found, err := client.Latest(
		context.Background(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Fatal("Latest() found = false, want true")
	}
	if release.Version != "v1.10.0" {
		t.Fatalf(
			"version = %q, want v1.10.0",
			release.Version,
		)
	}
	if release.BundleURL !=
		"https://downloads.example/masjidframe-update-v1.10.0-pi3.tar.zst" {
		t.Fatalf(
			"bundle URL = %q",
			release.BundleURL,
		)
	}
	if release.SignatureURL !=
		"https://downloads.example/masjidframe-update-v1.10.0-pi3.tar.zst.minisig" {
		t.Fatalf(
			"signature URL = %q",
			release.SignatureURL,
		)
	}
}

func TestGitHubClientReportsNoEligibleRelease(t *testing.T) {
	releases := []githubRelease{
		githubTestRelease(
			"v1.7.0-rc.1",
			false,
			true,
			true,
			true,
		),
		githubTestRelease(
			"v1.6.0",
			false,
			false,
			true,
			false,
		),
	}

	server := githubTestServer(t, releases)
	defer server.Close()

	release, found, err := (GitHubClient{
		HTTPClient:  server.Client(),
		ReleasesURL: server.URL,
	}).Latest(context.Background())

	if err != nil {
		t.Fatal(err)
	}
	if found {
		t.Fatalf(
			"Latest() = %+v, true; want no release",
			release,
		)
	}
}

func TestGitHubClientRejectsDuplicateRequiredAsset(
	t *testing.T,
) {
	release := githubTestRelease(
		"v1.7.0",
		false,
		false,
		true,
		true,
	)
	release.Assets = append(
		release.Assets,
		release.Assets[0],
	)

	server := githubTestServer(
		t,
		[]githubRelease{release},
	)
	defer server.Close()

	_, found, err := (GitHubClient{
		HTTPClient:  server.Client(),
		ReleasesURL: server.URL,
	}).Latest(context.Background())

	if err != nil {
		t.Fatal(err)
	}
	if found {
		t.Fatal(
			"Latest() accepted duplicate required asset",
		)
	}
}
