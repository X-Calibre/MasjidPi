package updates

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGitHubClientRejectsUnexpectedHTTPStatus(
	t *testing.T,
) {
	server := httptest.NewServer(http.HandlerFunc(
		func(
			response http.ResponseWriter,
			_ *http.Request,
		) {
			http.Error(
				response,
				"unavailable",
				http.StatusServiceUnavailable,
			)
		},
	))
	defer server.Close()

	_, _, err := (GitHubClient{
		HTTPClient:  server.Client(),
		ReleasesURL: server.URL,
	}).Latest(context.Background())

	if err == nil {
		t.Fatal("Latest() error = nil, want HTTP error")
	}
	if !strings.Contains(
		err.Error(),
		"503 Service Unavailable",
	) {
		t.Fatalf("Latest() error = %q", err)
	}
}

func TestGitHubClientRejectsMalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(
			response http.ResponseWriter,
			_ *http.Request,
		) {
			_, _ = response.Write(
				[]byte(`[{"tag_name":`),
			)
		},
	))
	defer server.Close()

	_, _, err := (GitHubClient{
		HTTPClient:  server.Client(),
		ReleasesURL: server.URL,
	}).Latest(context.Background())

	if err == nil {
		t.Fatal("Latest() error = nil, want decode error")
	}
	if !strings.Contains(
		err.Error(),
		"decode GitHub releases",
	) {
		t.Fatalf("Latest() error = %q", err)
	}
}

func TestGitHubClientRejectsOversizedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(
			response http.ResponseWriter,
			_ *http.Request,
		) {
			_, _ = response.Write(
				[]byte(
					strings.Repeat(
						"x",
						githubResponseLimit+1,
					),
				),
			)
		},
	))
	defer server.Close()

	_, _, err := (GitHubClient{
		HTTPClient:  server.Client(),
		ReleasesURL: server.URL,
	}).Latest(context.Background())

	if err == nil {
		t.Fatal("Latest() error = nil, want size error")
	}
	if !strings.Contains(
		err.Error(),
		"response exceeds",
	) {
		t.Fatalf("Latest() error = %q", err)
	}
}

func TestGitHubClientRejectsInsecureAssetURL(
	t *testing.T,
) {
	release := githubTestRelease(
		"v1.7.0",
		false,
		false,
		true,
		true,
	)
	release.Assets[0].BrowserDownloadURL =
		"http://downloads.example/insecure.tar.zst"

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
		t.Fatal("Latest() accepted an insecure asset URL")
	}
}

func TestGitHubClientRejectsInvalidEndpoint(t *testing.T) {
	_, _, err := (GitHubClient{
		ReleasesURL: "://invalid",
	}).Latest(context.Background())

	if err == nil {
		t.Fatal("Latest() error = nil, want URL error")
	}
	if !strings.Contains(
		err.Error(),
		"invalid GitHub releases URL",
	) {
		t.Fatalf("Latest() error = %q", err)
	}
}
