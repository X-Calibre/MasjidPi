package updates

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeBundlePreparer struct {
	artifact PreparedArtifact
	err      error
	calls    int
}

type blockingBundlePreparer struct {
	started chan struct{}
	release chan struct{}
}

func (p *blockingBundlePreparer) Prepare(
	context.Context,
	Release,
) (PreparedArtifact, error) {
	close(p.started)
	<-p.release
	return PreparedArtifact{BundleBytes: 1, SignatureBytes: 1}, nil
}

func (p *fakeBundlePreparer) Prepare(
	context.Context,
	Release,
) (PreparedArtifact, error) {
	p.calls++
	return p.artifact, p.err
}

func TestControllerPrepareRecordsVerifiedDownload(t *testing.T) {
	now := time.Date(2026, time.September, 16, 9, 0, 0, 0, time.UTC)
	controller, _ := newTestController(t, "v1.6.0", &sequenceReleaseSource{
		results: []sourceResult{{release: testRelease("v1.7.0"), found: true}},
	}, now)
	preparer := &fakeBundlePreparer{artifact: PreparedArtifact{
		BundleBytes: 1234, SignatureBytes: 343,
	}}
	controller.SetBundlePreparer(preparer)
	if _, err := controller.Check(t.Context()); err != nil {
		t.Fatal(err)
	}

	state, err := controller.Prepare(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if state.Download == nil || state.Download.Status != DownloadStatusVerified {
		t.Fatalf("download = %+v", state.Download)
	}
	if state.Download.BundleBytes != 1234 || state.Download.SignatureBytes != 343 {
		t.Fatalf("download sizes = %+v", state.Download)
	}
	if state.Download.VerifiedAt == nil {
		t.Fatal("verified time is nil")
	}
	if preparer.calls != 1 {
		t.Fatalf("prepare calls = %d, want 1", preparer.calls)
	}
	if _, err := controller.Prepare(t.Context()); err != nil {
		t.Fatal(err)
	}
	if preparer.calls != 1 {
		t.Fatalf("verified download prepared again; calls = %d", preparer.calls)
	}
}

func TestControllerPrepareRecordsFailureForRetry(t *testing.T) {
	now := time.Date(2026, time.September, 16, 9, 0, 0, 0, time.UTC)
	controller, _ := newTestController(t, "v1.6.0", &sequenceReleaseSource{
		results: []sourceResult{{release: testRelease("v1.7.0"), found: true}},
	}, now)
	controller.SetBundlePreparer(&fakeBundlePreparer{err: errors.New("network unavailable")})
	if _, err := controller.Check(t.Context()); err != nil {
		t.Fatal(err)
	}

	state, err := controller.Prepare(t.Context())
	if err == nil {
		t.Fatal("prepare error = nil")
	}
	if state.Download == nil || state.Download.Status != DownloadStatusFailed ||
		state.Download.LastError == "" {
		t.Fatalf("download failure = %+v", state.Download)
	}
}

func TestControllerStatusRemainsAvailableDuringPreparation(t *testing.T) {
	now := time.Date(2026, time.September, 16, 9, 0, 0, 0, time.UTC)
	controller, _ := newTestController(t, "v1.6.0", &sequenceReleaseSource{
		results: []sourceResult{{release: testRelease("v1.7.0"), found: true}},
	}, now)
	preparer := &blockingBundlePreparer{
		started: make(chan struct{}),
		release: make(chan struct{}),
	}
	controller.SetBundlePreparer(preparer)
	if _, err := controller.Check(t.Context()); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() {
		_, err := controller.Prepare(t.Context())
		done <- err
	}()
	<-preparer.started

	state, err := controller.Status()
	if err != nil {
		t.Fatal(err)
	}
	if state.Download == nil || state.Download.Status != DownloadStatusDownloading {
		t.Fatalf("download = %+v", state.Download)
	}
	close(preparer.release)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}
