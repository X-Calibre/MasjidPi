package player

import (
	"errors"
	"os"
	"os/exec"
	"sync"
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

	cmd := exec.Command(
		"mpv",
		"--idle=yes",
		"--no-video",
		"--no-ytdl",
		"--really-quiet",
		"--terminal=no",
		"--ao=alsa",
		"--input-ipc-server="+p.socket,
	)

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
