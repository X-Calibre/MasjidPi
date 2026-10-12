package player

import (
	"testing"
	"time"
)

func TestMPVStartupTimeout(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  time.Duration
	}{
		{"", 5 * time.Second},
		{"30s", 30 * time.Second},
		{"invalid", 5 * time.Second},
		{"-1s", 5 * time.Second},
		{"0s", 5 * time.Second},
		{"2h", 5 * time.Second},
	} {
		t.Run(tc.value, func(t *testing.T) {
			t.Setenv("MASJIDFRAME_MPV_STARTUP_TIMEOUT", tc.value)
			if got := mpvStartupTimeout(); got != tc.want {
				t.Fatalf("timeout = %s, want %s", got, tc.want)
			}
		})
	}
}
