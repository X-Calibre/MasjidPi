package app

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
	masjidboardruntime "github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/runtime"
	"github.com/X-Calibre/MasjidPi/backend/internal/playback"
	"github.com/X-Calibre/MasjidPi/backend/internal/updates"
)

type recordingUpdatePlayback struct {
	status playback.Status
	stops  int
}

func (p *recordingUpdatePlayback) Status() playback.Status { return p.status }
func (p *recordingUpdatePlayback) Stop()                   { p.stops++ }

type recordingInstaller struct{ calls int }

func (i *recordingInstaller) Install(context.Context, string) error {
	i.calls++
	return nil
}

type recordingRebooter struct{ calls int }

func (r *recordingRebooter) Reboot(context.Context) error {
	r.calls++
	return nil
}

func boardAtAdhan(now time.Time) fakeUpdateBoardProvider {
	return fakeUpdateBoardProvider{results: []masjidboardruntime.Result{{
		Board: &model.Board{PrayerTimes: model.PrayerTimes{
			Dhuhr: model.PrayerTime{Adhan: &model.ClockTime{
				Hour:   now.Hour(),
				Minute: now.Minute(),
			}},
		}},
	}}}
}

func installReadyController(
	t *testing.T,
	now time.Time,
) (*applianceUpdates, *recordingInstaller, *recordingRebooter) {
	t.Helper()
	store := updates.NewStore(filepath.Join(t.TempDir(), "update_state.json"))
	release := updates.Release{
		Version:      "v1.7.0",
		PublishedAt:  now.Add(-48 * time.Hour),
		PageURL:      "https://example.test/v1.7.0",
		BundleURL:    "https://example.test/update.tar.zst",
		SignatureURL: "https://example.test/update.tar.zst.minisig",
	}
	detected := now.Add(-24 * time.Hour)
	deadline := detected.Add(updates.ApprovalPeriod)
	verified := now.Add(-time.Hour)
	if err := store.Save(updates.State{
		SchemaVersion:    updates.StateSchemaVersion,
		CurrentVersion:   "v1.6.0",
		AvailableRelease: &release,
		FirstDetectedAt:  &detected,
		ApprovalDeadline: &deadline,
		Download: &updates.DownloadState{
			Version:    release.Version,
			Status:     updates.DownloadStatusVerified,
			VerifiedAt: &verified,
		},
	}); err != nil {
		t.Fatal(err)
	}
	controller := updates.NewController(store, nil, "v1.6.0")
	installer := &recordingInstaller{}
	rebooter := &recordingRebooter{}
	controller.SetInstaller(installer, rebooter)
	appliance := newApplianceUpdates(controller)
	appliance.now = func() time.Time { return now }
	return appliance, installer, rebooter
}

func TestImmediateInstallStopsPlaybackAfterSafetyChecks(t *testing.T) {
	now := time.Date(2026, 9, 16, 12, 0, 0, 0, time.Local)
	appliance, installer, rebooter := installReadyController(t, now)
	playbackProvider := &recordingUpdatePlayback{
		status: playback.Status{Listening: true},
	}
	appliance.SetSafetyProviders(playbackProvider, nil)

	if _, err := appliance.Install(t.Context(), true, true); err != nil {
		t.Fatal(err)
	}
	if playbackProvider.stops != 1 || installer.calls != 1 || rebooter.calls != 1 {
		t.Fatalf(
			"stops=%d installs=%d reboots=%d",
			playbackProvider.stops,
			installer.calls,
			rebooter.calls,
		)
	}
}

func TestImmediateInstallDoesNotStopPlaybackWhenAdhanGuardBlocks(t *testing.T) {
	now := time.Date(2026, 9, 16, 12, 0, 0, 0, time.Local)
	appliance, installer, rebooter := installReadyController(t, now)
	playbackProvider := &recordingUpdatePlayback{
		status: playback.Status{Listening: true},
	}
	boardProvider := fakeUpdateBoardProvider{results: nil}
	appliance.SetSafetyProviders(playbackProvider, boardProvider)

	conditions := currentInstallConditions(now, true, true, playbackProvider, nil)
	conditions.AdhanTimes = []time.Time{now}
	state, err := appliance.Status()
	if err != nil {
		t.Fatal(err)
	}
	if updates.EvaluateInstall(state, conditions).Allowed {
		t.Fatal("test setup should be blocked by Adhan guard")
	}

	// A Board provider with no results cannot fabricate an Adhan, so verify the
	// adapter's ordering through a provider containing an exact current Adhan.
	boardProvider = boardAtAdhan(now)
	appliance.SetSafetyProviders(playbackProvider, boardProvider)
	if _, err := appliance.Install(t.Context(), true, true); err == nil {
		t.Fatal("Install() error = nil")
	}
	if playbackProvider.stops != 0 || installer.calls != 0 || rebooter.calls != 0 {
		t.Fatalf(
			"stops=%d installs=%d reboots=%d",
			playbackProvider.stops,
			installer.calls,
			rebooter.calls,
		)
	}
}
