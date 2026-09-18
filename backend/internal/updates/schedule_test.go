package updates

import (
	"testing"
	"time"
)

func TestCheckDue(t *testing.T) {
	now := time.Date(
		2026, time.September, 16,
		12, 0, 0, 0, time.UTC,
	)

	tests := []struct {
		name        string
		lastChecked *time.Time
		want        bool
	}{
		{
			name: "never checked",
			want: true,
		},
		{
			name: "one second before interval",
			lastChecked: scheduleTimePointer(
				now.Add(-CheckInterval + time.Second),
			),
			want: false,
		},
		{
			name: "exactly at interval",
			lastChecked: scheduleTimePointer(
				now.Add(-CheckInterval),
			),
			want: true,
		},
		{
			name: "past interval",
			lastChecked: scheduleTimePointer(
				now.Add(-CheckInterval - time.Hour),
			),
			want: true,
		},
		{
			name: "recent check",
			lastChecked: scheduleTimePointer(
				now.Add(-24 * time.Hour),
			),
			want: false,
		},
		{
			name: "future timestamp",
			lastChecked: scheduleTimePointer(
				now.Add(time.Hour),
			),
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := State{
				SchemaVersion: StateSchemaVersion,
				LastCheckedAt: test.lastChecked,
			}

			got := CheckDue(state, now)
			if got != test.want {
				t.Fatalf(
					"CheckDue() = %t, want %t",
					got,
					test.want,
				)
			}
		})
	}
}

func TestCheckDueNormalizesTimeZones(t *testing.T) {
	now := time.Date(
		2026, time.September, 16,
		12, 0, 0, 0, time.UTC,
	)
	location := time.FixedZone("SAST", 2*60*60)
	lastChecked := now.
		Add(-CheckInterval).
		In(location)

	state := State{
		SchemaVersion: StateSchemaVersion,
		LastCheckedAt: &lastChecked,
	}

	if !CheckDue(state, now) {
		t.Fatal("CheckDue() = false, want true")
	}
}

func scheduleTimePointer(value time.Time) *time.Time {
	return &value
}
