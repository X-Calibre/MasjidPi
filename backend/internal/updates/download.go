package updates

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

const (
	maxBundleBytes    = int64(2 << 30)
	maxSignatureBytes = int64(64 << 10)
	spaceReserveBytes = uint64(512 << 20)
	downloadTimeout   = 6 * time.Hour
)

type PreparedArtifact struct {
	BundlePath     string
	SignaturePath  string
	BundleBytes    int64
	SignatureBytes int64
}

type BundleVerifier interface {
	Verify(context.Context, string, string) error
}

type CommandVerifier struct {
	Command string
}

func (v CommandVerifier) Verify(
	ctx context.Context,
	bundle string,
	signature string,
) error {
	command := v.Command
	if command == "" {
		command = "/usr/local/sbin/masjidframe-update"
	}
	output, err := exec.CommandContext(
		ctx,
		command,
		"verify",
		bundle,
		signature,
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"updates: verify downloaded bundle: %w: %s",
			err,
			string(output),
		)
	}
	return nil
}

type Downloader struct {
	Directory  string
	HTTPClient *http.Client
	Verifier   BundleVerifier
}

func (d Downloader) Prepare(
	ctx context.Context,
	release Release,
) (PreparedArtifact, error) {
	if _, err := ParseStableVersion(release.Version); err != nil {
		return PreparedArtifact{}, err
	}
	if !isHTTPSURL(release.BundleURL) || !isHTTPSURL(release.SignatureURL) {
		return PreparedArtifact{}, fmt.Errorf("updates: download URLs must use HTTPS")
	}
	if d.Directory == "" {
		return PreparedArtifact{}, fmt.Errorf("updates: download directory is required")
	}
	if d.Verifier == nil {
		return PreparedArtifact{}, fmt.Errorf("updates: bundle verifier is required")
	}
	if err := os.MkdirAll(d.Directory, 0700); err != nil {
		return PreparedArtifact{}, fmt.Errorf("updates: create download directory: %w", err)
	}

	bundleName := fmt.Sprintf("masjidframe-update-%s-pi3.tar.zst", release.Version)
	bundlePath := filepath.Join(d.Directory, bundleName)
	signaturePath := bundlePath + ".minisig"
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(bundlePath)
			_ = os.Remove(signaturePath)
		}
	}()

	bundleBytes, err := d.download(ctx, release.BundleURL, bundlePath, maxBundleBytes)
	if err != nil {
		return PreparedArtifact{}, err
	}
	signatureBytes, err := d.download(ctx, release.SignatureURL, signaturePath, maxSignatureBytes)
	if err != nil {
		return PreparedArtifact{}, err
	}
	if err := d.Verifier.Verify(ctx, bundlePath, signaturePath); err != nil {
		return PreparedArtifact{}, err
	}

	cleanup = false
	return PreparedArtifact{
		BundlePath:     bundlePath,
		SignaturePath:  signaturePath,
		BundleBytes:    bundleBytes,
		SignatureBytes: signatureBytes,
	}, nil
}

func (d Downloader) Cleanup(version string) error {
	if _, err := ParseStableVersion(version); err != nil {
		return err
	}
	if d.Directory == "" {
		return fmt.Errorf("updates: download directory is required")
	}
	bundlePath := filepath.Join(
		d.Directory,
		fmt.Sprintf("masjidframe-update-%s-pi3.tar.zst", version),
	)
	var cleanupErr error
	for _, path := range []string{bundlePath, bundlePath + ".minisig"} {
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			cleanupErr = errors.Join(
				cleanupErr,
				fmt.Errorf("updates: remove %s: %w", filepath.Base(path), err),
			)
		}
	}
	return cleanupErr
}

func (d Downloader) download(
	ctx context.Context,
	address string,
	destination string,
	limit int64,
) (int64, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, address, nil)
	if err != nil {
		return 0, fmt.Errorf("updates: create download request: %w", err)
	}
	request.Header.Set("User-Agent", "MasjidFrame Update Downloader")

	client := d.client()
	response, err := client.Do(request)
	if err != nil {
		return 0, fmt.Errorf("updates: download %s: %w", filepath.Base(destination), err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("updates: download %s returned %s", filepath.Base(destination), response.Status)
	}
	if response.ContentLength > limit {
		return 0, fmt.Errorf("updates: download %s exceeds size limit", filepath.Base(destination))
	}
	required := limit
	if response.ContentLength >= 0 {
		required = response.ContentLength
	}
	if err := requireFreeSpace(d.Directory, uint64(required)); err != nil {
		return 0, err
	}

	temporary, err := os.CreateTemp(d.Directory, ".download-*")
	if err != nil {
		return 0, fmt.Errorf("updates: create temporary download: %w", err)
	}
	temporaryPath := temporary.Name()
	keep := false
	defer func() {
		_ = temporary.Close()
		if !keep {
			_ = os.Remove(temporaryPath)
		}
	}()
	if err := temporary.Chmod(0600); err != nil {
		return 0, fmt.Errorf("updates: secure temporary download: %w", err)
	}

	written, err := io.Copy(temporary, io.LimitReader(response.Body, limit+1))
	if err != nil {
		return 0, fmt.Errorf("updates: write temporary download: %w", err)
	}
	if written > limit {
		return 0, fmt.Errorf("updates: download %s exceeds size limit", filepath.Base(destination))
	}
	if err := temporary.Sync(); err != nil {
		return 0, fmt.Errorf("updates: sync temporary download: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return 0, fmt.Errorf("updates: close temporary download: %w", err)
	}
	if err := os.Rename(temporaryPath, destination); err != nil {
		return 0, fmt.Errorf("updates: activate downloaded file: %w", err)
	}
	keep = true
	return written, nil
}

func (d Downloader) client() *http.Client {
	base := d.HTTPClient
	if base == nil {
		base = &http.Client{Timeout: downloadTimeout}
	}
	client := *base
	previous := client.CheckRedirect
	client.CheckRedirect = func(request *http.Request, via []*http.Request) error {
		if request.URL.Scheme != "https" {
			return fmt.Errorf("updates: redirect URL must use HTTPS")
		}
		if previous != nil {
			return previous(request, via)
		}
		if len(via) >= 10 {
			return fmt.Errorf("updates: too many redirects")
		}
		return nil
	}
	return &client
}

func requireFreeSpace(directory string, required uint64) error {
	var info syscall.Statfs_t
	if err := syscall.Statfs(directory, &info); err != nil {
		return fmt.Errorf("updates: inspect available space: %w", err)
	}
	available := uint64(info.Bavail) * uint64(info.Bsize)
	if available < required+spaceReserveBytes {
		return fmt.Errorf("updates: insufficient download space")
	}
	return nil
}
