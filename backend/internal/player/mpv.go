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

	deadline := time.Now().Add(3 * time.Second)

	for {
		if err := m.ipc.Connect(); err == nil {
			return nil
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("wait for mpv ipc: timed out")
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func (m *MPV) Stop() error {
	_, err := m.execute("stop")
	return err
}

func (m *MPV) Close() error {
	_ = m.ipc.Close()

	return m.process.Stop()
}

func (m *MPV) execute(command ...any) (*Response, error) {
	cmd := Command{Command: command}

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

func (m *MPV) GetProperty(name string) (any, error) {
	resp, err := m.execute("get_property", name)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (m *MPV) SetProperty(name string, value any) error {
	_, err := m.execute("set_property", name, value)
	return err
}

func (m *MPV) Version() (string, error) {
	value, err := m.GetProperty("mpv-version")
	if err != nil {
		return "", err
	}

	version, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("mpv-version has type %T with value %#v", value, value)
	}
	return version, nil
}

func (m *MPV) Play(url string) error {
	_, err := m.execute("loadfile", url)
	return err
}

func (m *MPV) Volume(volume int) error {
	if volume < 0 || volume > 125 {
		return fmt.Errorf("volume must be between 0 and 125")
	}
	return m.SetProperty("volume", volume)
}

func (m *MPV) AudioDevices() ([]AudioDevice, error) {
	value, err := m.GetProperty("audio-device-list")
	if err != nil {
		return nil, err
	}

	items, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("audio-device-list is %T (%v)", value, value)
	}

	devices := make([]AudioDevice, 0, len(items))
	for _, item := range items {
		entry, ok := item.(map[string]any)
		if !ok {
			continue
		}

		name, _ := entry["name"].(string)
		description, _ := entry["description"].(string)
		if name == "" {
			continue
		}

		devices = append(devices, AudioDevice{
			Name:        name,
			Description: description,
		})
	}

	return devices, nil
}

func (m *MPV) AudioDevice(name string) error {
	if name == "" {
		return fmt.Errorf("audio device cannot be empty")
	}
	return m.SetProperty("audio-device", name)
}

func (m *MPV) Status() (*Status, error) {
	version, err := m.Version()
	if err != nil {
		return nil, err
	}

	pausedValue, err := m.GetProperty("pause")
	if err != nil {
		return nil, err
	}
	paused, ok := pausedValue.(bool)
	if !ok {
		return nil, fmt.Errorf("pause is %T (%v)", pausedValue, pausedValue)
	}

	volumeValue, err := m.GetProperty("volume")
	if err != nil {
		return nil, err
	}
	volumeFloat, ok := volumeValue.(float64)
	if !ok {
		return nil, fmt.Errorf("volume is %T (%v)", volumeValue, volumeValue)
	}

	idleValue, err := m.GetProperty("core-idle")
	if err != nil {
		return nil, err
	}
	idle, ok := idleValue.(bool)
	if !ok {
		return nil, fmt.Errorf("core-idle is %T (%v)", idleValue, idleValue)
	}

	state := "playing"
	if idle {
		state = "stopped"
	}

	path := ""
	pathValue, err := m.GetProperty("path")
	if err == nil {
		path, ok = pathValue.(string)
		if !ok {
			return nil, fmt.Errorf("path is %T (%v)", pathValue, pathValue)
		}
	}

	audioDevice := ""
	if value, err := m.GetProperty("audio-device"); err == nil {
		audioDevice, _ = value.(string)
	}

	return &Status{
		Version:     version,
		State:       state,
		URL:         path,
		Volume:      int(volumeFloat),
		Paused:      paused,
		AudioDevice: audioDevice,
	}, nil
}
