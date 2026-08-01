package player

import (
	"os"
	"os/exec"
)

type Process struct {
	cmd *exec.Cmd
}

func NewProcess(socket string) *Process {
	return &Process{
		cmd: exec.Command(
			"mpv",
			"--idle=yes",
			"--no-video",
			"--input-ipc-server="+socket,
		),
	}
}

func (p *Process) Start() error {
	p.cmd.Stdout = os.Stdout
	p.cmd.Stderr = os.Stderr

	return p.cmd.Start()
}

func (p *Process) Stop() error {
	if p.cmd == nil || p.cmd.Process == nil {
		return nil
	}

	return p.cmd.Process.Kill()
}
