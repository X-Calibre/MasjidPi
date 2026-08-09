package player

import (
	"os"
	"os/exec"
)

type Process struct {
	socket string
	cmd    *exec.Cmd
}

func NewProcess(socket string) *Process {
	return &Process{
		socket: socket,
		cmd: exec.Command(
			"mpv",
			"--idle=yes",
			"--no-video",
			"--no-ytdl",
			"--really-quiet",
			"--terminal=no",
			"--ao=alsa",
			"--input-ipc-server="+socket,
		),
	}
}

func (p *Process) Start() error {
	// Remove any stale socket from a previous crash.
	_ = os.Remove(p.socket)

	// MPV runs as a background service. Do not forward its terminal/progress
	// output to MasjidPi's stdout/stderr and therefore into journald.
	p.cmd.Stdout = nil
	p.cmd.Stderr = nil

	return p.cmd.Start()
}

func (p *Process) Stop() error {
	if p.cmd == nil || p.cmd.Process == nil {
		return nil
	}

	err := p.cmd.Process.Kill()

	_ = p.cmd.Wait()

	_ = os.Remove(p.socket)

	p.cmd = nil

	return err
}
