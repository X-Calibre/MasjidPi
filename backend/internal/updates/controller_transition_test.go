package updates

import (
	"testing"
	"time"
)

func TestControllerNewerReleaseResetsDetectionWindow(
	t *testing.T,
) {
	firstCheck := time.Date(
		2026, time.September, 16,
		9, 0, 0, 0, time.UTC,
	)
	secondCheck := firstCheck.Add(7 * 24 * time.Hour)

	source := &sequenceReleaseSource{
		results: []sourceResult{
			{
				release: testRelease("v1.7.0"),
				found:   true,
			},
			{
				release: testRelease("v1.8.0"),
				found:   true,
			},
		},
	}
	controller, _ := newTestController(
		t,
		"v1.6.0",
		source,
		firstCheck,
	)

	first, err := controller.Check(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	controller.now = func() time.Time {
		return secondCheck
	}

	second, err := controller.Check(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if first.AvailableRelease == nil ||
		first.AvailableRelease.Version != "v1.7.0" {
		t.Fatalf(
			"first available release = %+v",
			first.AvailableRelease,
		)
	}
	if second.AvailableRelease == nil ||
		second.AvailableRelease.Version != "v1.8.0" {
		t.Fatalf(
			"second available release = %+v",
			second.AvailableRelease,
		)
	}

	requireControllerTime(
		t,
		"first detected",
		second.FirstDetectedAt,
		secondCheck,
	)
	requireControllerTime(
		t,
		"approval deadline",
		second.ApprovalDeadline,
		secondCheck.Add(ApprovalPeriod),
	)
	requireControllerTime(
		t,
		"last checked",
		second.LastCheckedAt,
		secondCheck,
	)
}

func TestControllerClearsReleaseNotNewerThanInstalled(
	t *testing.T,
) {
	tests := []struct {
		name      string
		installed string
		candidate string
	}{
		{
			name:      "equal stable release",
			installed: "v1.7.0",
			candidate: "v1.7.0",
		},
		{
			name:      "older stable release",
			installed: "v1.7.0",
			candidate: "v1.6.9",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			now := time.Date(
				2026, time.September, 16,
				10, 0, 0, 0, time.UTC,
			)
			detected := now.Add(-24 * time.Hour)
			deadline := detected.Add(ApprovalPeriod)
			previous := testRelease("v1.8.0")

			source := &sequenceReleaseSource{
				results: []sourceResult{{
					release: testRelease(
						test.candidate,
					),
					found: true,
				}},
			}
			controller, store := newTestController(
				t,
				test.installed,
				source,
				now,
			)

			err := store.Save(State{
				SchemaVersion:    StateSchemaVersion,
				CurrentVersion:   test.installed,
				AvailableRelease: &previous,
				FirstDetectedAt:  &detected,
				ApprovalDeadline: &deadline,
			})
			if err != nil {
				t.Fatal(err)
			}

			state, err := controller.Check(t.Context())
			if err != nil {
				t.Fatal(err)
			}

			if state.AvailableRelease != nil {
				t.Fatalf(
					"available release = %+v, want nil",
					state.AvailableRelease,
				)
			}
			if state.FirstDetectedAt != nil {
				t.Fatalf(
					"first detected = %v, want nil",
					state.FirstDetectedAt,
				)
			}
			if state.ApprovalDeadline != nil {
				t.Fatalf(
					"approval deadline = %v, want nil",
					state.ApprovalDeadline,
				)
			}
			requireControllerTime(
				t,
				"last checked",
				state.LastCheckedAt,
				now,
			)
		})
	}
}

func TestControllerFinalStableSupersedesInstalledRC(
	t *testing.T,
) {
	now := time.Date(
		2026, time.September, 16,
		11, 0, 0, 0, time.UTC,
	)
	source := &sequenceReleaseSource{
		results: []sourceResult{{
			release: testRelease("v1.6.0"),
			found:   true,
		}},
	}
	controller, _ := newTestController(
		t,
		"v1.6.0-rc.4-image",
		source,
		now,
	)

	state, err := controller.Check(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if state.AvailableRelease == nil ||
		state.AvailableRelease.Version != "v1.6.0" {
		t.Fatalf(
			"available release = %+v, want v1.6.0",
			state.AvailableRelease,
		)
	}
	requireControllerTime(
		t,
		"first detected",
		state.FirstDetectedAt,
		now,
	)
	requireControllerTime(
		t,
		"approval deadline",
		state.ApprovalDeadline,
		now.Add(ApprovalPeriod),
	)
}

func requireControllerTime(
	t *testing.T,
	name string,
	actual *time.Time,
	expected time.Time,
) {
	t.Helper()

	if actual == nil {
		t.Fatalf("%s = nil, want %v", name, expected)
	}
	if !actual.Equal(expected) {
		t.Fatalf(
			"%s = %v, want %v",
			name,
			*actual,
			expected,
		)
	}
}
