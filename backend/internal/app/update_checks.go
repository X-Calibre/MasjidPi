package app

import (
	"context"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/updates"
)

const updateCheckPollInterval = time.Hour

type updateChecker interface {
	Status() (updates.State, error)
	Check(context.Context) (updates.State, error)
	Prepare(context.Context) (updates.State, error)
	InstallAutomatic(context.Context, time.Time) (updates.State, bool, error)
}

type updateCheckLogger interface {
	Info(string, ...any)
	Warn(string, ...any)
}

func monitorUpdateChecks(
	ctx context.Context,
	checker updateChecker,
	log updateCheckLogger,
) {
	checkUpdatesIfDue(
		ctx,
		checker,
		time.Now(),
		log,
	)

	ticker := time.NewTicker(
		updateCheckPollInterval,
	)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			checkUpdatesIfDue(
				ctx,
				checker,
				now,
				log,
			)
		}
	}
}

func checkUpdatesIfDue(
	ctx context.Context,
	checker updateChecker,
	now time.Time,
	log updateCheckLogger,
) {
	state, err := checker.Status()
	if err != nil {
		log.Warn(
			"Could not load automatic update state",
			"error",
			err,
		)
		return
	}

	if updates.CheckDue(state, now) {
		state, err = checker.Check(ctx)
		if err != nil {
			log.Warn(
				"Automatic update check failed",
				"error",
				err,
			)
			return
		}
	}

	if state.AvailableRelease != nil {
		log.Info(
			"Stable appliance update available",
			"version",
			state.AvailableRelease.Version,
			"approval_deadline",
			state.ApprovalDeadline,
		)
		if state.Download == nil ||
			state.Download.Status != updates.DownloadStatusVerified {
			state, err = checker.Prepare(ctx)
			if err != nil {
				log.Warn(
					"Automatic update preparation failed",
					"error",
					err,
				)
				return
			}
			log.Info(
				"Signed appliance update is ready",
				"version",
				state.AvailableRelease.Version,
			)
		}
		if state.Download != nil &&
			state.Download.Status == updates.DownloadStatusVerified {
			state, attempted, installErr := checker.InstallAutomatic(ctx, now)
			if installErr != nil {
				log.Warn(
					"Automatic update installation failed",
					"error",
					installErr,
				)
				return
			}
			if attempted {
				log.Info(
					"Automatic update installation staged",
					"version",
					state.AvailableRelease.Version,
				)
			}
		}
		return
	}

	log.Info(
		"Automatic update check completed",
		"current_version",
		state.CurrentVersion,
	)
}
