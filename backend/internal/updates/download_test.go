package updates

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

type recordingVerifier struct {
	err       error
	bundle    string
	signature string
}

func (v *recordingVerifier) Verify(
	_ context.Context,
	bundle string,
	signature string,
) error {
	v.bundle = bundle
	v.signature = signature
	return v.err
}

func TestDownloaderPreparesVerifiedArtifact(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/bundle":
				_, _ = w.Write([]byte("signed bundle"))
			case "/signature":
				_, _ = w.Write([]byte("signature"))
			default:
				http.NotFound(w, r)
			}
		},
	))
	defer server.Close()

	directory := t.TempDir()
	verifier := &recordingVerifier{}
	downloader := Downloader{
		Directory:  directory,
		HTTPClient: server.Client(),
		Verifier:   verifier,
	}
	release := testRelease("v1.7.0")
	release.BundleURL = server.URL + "/bundle"
	release.SignatureURL = server.URL + "/signature"

	artifact, err := downloader.Prepare(t.Context(), release)
	if err != nil {
		t.Fatal(err)
	}
	if artifact.BundleBytes != int64(len("signed bundle")) ||
		artifact.SignatureBytes != int64(len("signature")) {
		t.Fatalf("artifact sizes = %+v", artifact)
	}
	if verifier.bundle != artifact.BundlePath ||
		verifier.signature != artifact.SignaturePath {
		t.Fatalf("verified paths = %q, %q", verifier.bundle, verifier.signature)
	}
	for _, path := range []string{artifact.BundlePath, artifact.SignaturePath} {
		info, statErr := os.Stat(path)
		if statErr != nil {
			t.Fatal(statErr)
		}
		if info.Mode().Perm() != 0600 {
			t.Fatalf("%s permissions = %o, want 600", path, info.Mode().Perm())
		}
	}
	if matches, _ := filepath.Glob(filepath.Join(directory, ".download-*")); len(matches) != 0 {
		t.Fatalf("temporary downloads remain: %v", matches)
	}
}

func TestDownloaderRejectsInsecureURLs(t *testing.T) {
	release := testRelease("v1.7.0")
	release.BundleURL = "http://downloads.example/bundle"
	_, err := (Downloader{
		Directory: t.TempDir(),
		Verifier:  &recordingVerifier{},
	}).Prepare(t.Context(), release)
	if err == nil {
		t.Fatal("insecure download URL succeeded")
	}
}

func TestDownloaderRemovesFilesAfterVerificationFailure(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("payload"))
		},
	))
	defer server.Close()

	directory := t.TempDir()
	release := testRelease("v1.7.0")
	release.BundleURL = server.URL + "/bundle"
	release.SignatureURL = server.URL + "/signature"
	_, err := (Downloader{
		Directory:  directory,
		HTTPClient: server.Client(),
		Verifier: &recordingVerifier{
			err: errors.New("untrusted bundle"),
		},
	}).Prepare(t.Context(), release)
	if err == nil {
		t.Fatal("verification failure returned nil")
	}
	entries, readErr := os.ReadDir(directory)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if len(entries) != 0 {
		t.Fatalf("files remain after failure: %v", entries)
	}
}

func TestDownloaderCleanupRemovesOnlyRequestedRelease(t *testing.T) {
	directory := t.TempDir()
	target := filepath.Join(directory, "masjidframe-update-v1.7.0-pi3.tar.zst")
	other := filepath.Join(directory, "masjidframe-update-v1.8.0-pi3.tar.zst")
	for _, path := range []string{target, target + ".minisig", other} {
		if err := os.WriteFile(path, []byte("test"), 0600); err != nil {
			t.Fatal(err)
		}
	}

	downloader := Downloader{Directory: directory}
	if err := downloader.Cleanup("v1.7.0"); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{target, target + ".minisig"} {
		if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("%s still exists: %v", path, err)
		}
	}
	if _, err := os.Stat(other); err != nil {
		t.Fatalf("unrelated release removed: %v", err)
	}
	if err := downloader.Cleanup("v1.7.0"); err != nil {
		t.Fatalf("repeated cleanup: %v", err)
	}
}
