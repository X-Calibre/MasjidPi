package player

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestProcessRestartsAfterUnexpectedExit(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("MasjidPi process tests require Linux")
	}

	dir := t.TempDir()
	mpvPath := filepath.Join(dir, "mpv")
	script := "#!/bin/sh\nsleep 0.1\n"
	if err := os.WriteFile(mpvPath, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	oldPath := os.Getenv("PATH")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+oldPath)

	p := NewProcess(filepath.Join(t.TempDir(), "masjidpi.sock"))
	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	deadline := time.Now().Add(2 * time.Second)
	for !p.Exited() && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if !p.Exited() {
		t.Fatal("process did not exit")
	}

	if err := p.Restart(); err != nil {
		t.Fatal(err)
	}
	if p.Exited() {
		t.Fatal("restarted process is not running")
	}

	if err := p.Stop(); err != nil {
		t.Fatal(err)
	}
}
