package player

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

const (
	// MasjidPi plays live audio rather than seekable video. Keep enough forward
	// cache for prolonged network jitter without retaining mpv's much larger
	// generic media-player cache, and retain only a small non-seekable back
	// buffer. These limits are understood by every supported mpv release.
	demuxerForwardCacheBytes = 32 * 1024 * 1024
	demuxerBackCacheBytes    = 4 * 1024 * 1024
)

type Process struct {
	mu     sync.Mutex
	socket string
	cmd    *exec.Cmd
	done   chan struct{}
}

func NewProcess(socket string) *Process {
	return &Process{socket: socket}
}

func (p *Process) Start() error {
	p.mu.Lock()
	if p.cmd != nil {
		select {
		case <-p.done:
			p.cmd = nil
			p.done = nil
		default:
			p.mu.Unlock()
			return nil
		}
	}

	// Remove any stale socket from a previous crash.
	_ = os.Remove(p.socket)

	cmd := exec.Command("mpv", mpvCommandArgs(p.socket)...)

	// MPV runs as a background service. Do not forward its terminal/progress
	// output to MasjidPi's stdout/stderr and therefore into journald.
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		p.mu.Unlock()
		return err
	}

	done := make(chan struct{})
	p.cmd = cmd
	p.done = done
	p.mu.Unlock()

	go func() {
		_ = cmd.Wait()
		close(done)
	}()

	return nil
}

func mpvCommandArgs(socket string) []string {
	return []string{
		"--idle=yes",
		"--no-video",
		"--no-ytdl",
		"--really-quiet",
		"--terminal=no",
		"--ao=alsa",
		"--demuxer-max-bytes=" + fmt.Sprint(demuxerForwardCacheBytes),
		"--demuxer-max-back-bytes=" + fmt.Sprint(demuxerBackCacheBytes),
		"--input-ipc-server=" + socket,
	}
}

func (p *Process) Stop() error {
	p.mu.Lock()
	cmd := p.cmd
	done := p.done
	p.mu.Unlock()

	if cmd == nil || cmd.Process == nil {
		return nil
	}

	err := cmd.Process.Kill()
	if err != nil && !errors.Is(err, os.ErrProcessDone) {
		return err
	}

	<-done
	_ = os.Remove(p.socket)

	p.mu.Lock()
	if p.cmd == cmd {
		p.cmd = nil
		p.done = nil
	}
	p.mu.Unlock()

	return nil
}

func (p *Process) Restart() error {
	if err := p.Stop(); err != nil {
		return err
	}
	return p.Start()
}

func (p *Process) Exited() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cmd == nil {
		return true
	}

	select {
	case <-p.done:
		return true
	default:
		return false
	}
}
