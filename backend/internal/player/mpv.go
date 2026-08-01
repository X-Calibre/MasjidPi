package player

import (
	"fmt"
	"time"
)

type MPV struct {
	process *Process
	ipc     *IPC
}

func New(socket string) *MPV {
	return &MPV{
		process: NewProcess(socket),
		ipc:     NewIPC(socket),
	}
}

func (m *MPV) Start() error {
	if err := m.process.Start(); err != nil {
		return err
	}

	// Temporary until we implement WaitForReady().
	time.Sleep(500 * time.Millisecond)
	//	if err := m.WaitForReady(3 * time.Second); err != nil {
	//		return err

	if err := m.ipc.Connect(); err != nil {
		return err
	}

	return nil
}

func (m *MPV) Close() error {
	_ = m.ipc.Close()

	return m.process.Stop()
}

func (m *MPV) execute(command ...any) (*Response, error) {
	cmd := Command{
		Command: command,
	}

	if err := m.ipc.Send(cmd); err != nil {
		return nil, err
	}

	var resp Response

	if err := m.ipc.Receive(&resp); err != nil {
		return nil, err
	}

	if resp.Error != "" && resp.Error != "success" {
		return nil, fmt.Errorf("mpv: %s", resp.Error)
	}

	return &resp, nil
}

func (m *MPV) Version() (string, error) {
	resp, err := m.execute(
		"get_property",
		"mpv-version",
	)
	if err != nil {
		return "", err
	}

	version, ok := resp.Data.(string)
	if !ok {
		return "", fmt.Errorf("unexpected response type")
	}

	return version, nil
}
