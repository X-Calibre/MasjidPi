package updates

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"
)

type sourceResult struct {
	release Release
	found   bool
	err     error
}

type sequenceReleaseSource struct {
	results []sourceResult
	calls   int
}

func (s *sequenceReleaseSource) Latest(
	context.Context,
) (Release, bool, error) {
	if s.calls >= len(s.results) {
		return Release{}, false, errors.New(
			"unexpected release-source call",
		)
	}

	result := s.results[s.calls]
	s.calls++
	return result.release, result.found, result.err
}

func testRelease(version string) Release {
	return Release{
		Version: version,
		PublishedAt: time.Date(
			2026,
			time.September,
			15,
			12,
			0,
			0,
			0,
			time.UTC,
		),
		PageURL: "https://github.example/releases/tag/" +
			version,
		BundleURL: "https://downloads.example/" +
			"masjidpi-update-" + version +
			"-pi3.tar.zst",
		SignatureURL: "https://downloads.example/" +
			"masjidpi-update-" + version +
			"-pi3.tar.zst.minisig",
	}
}

func newTestController(
	t *testing.T,
	currentVersion string,
	source ReleaseSource,
	now time.Time,
) (*Controller, *Store) {
	t.Helper()

	store := NewStore(
		filepath.Join(
			t.TempDir(),
			"update_state.json",
		),
	)
	controller := NewController(
		store,
		source,
		currentVersion,
	)
	controller.now = func() time.Time {
		return now
	}

	return controller, store
}
