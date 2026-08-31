package player

import (
	"testing"
	"time"
)

func TestMPVIPCReadyTimeoutAllowsSlowerApplianceStartup(t *testing.T) {
	if mpvIPCReadyTimeout != 5*time.Second {
		t.Fatalf("mpv IPC readiness timeout = %s, want 5s", mpvIPCReadyTimeout)
	}
}
