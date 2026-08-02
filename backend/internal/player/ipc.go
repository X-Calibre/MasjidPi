package player

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type IPC struct {
	socket string

	conn    net.Conn
	encoder *json.Encoder
	decoder *json.Decoder

	responses chan Response
}

func NewIPC(socket string) *IPC {
	return &IPC{
		socket: socket,
	}
}

func (i *IPC) Connect() error {
	conn, err := net.DialTimeout("unix", i.socket, 3*time.Second)
	if err != nil {
		return fmt.Errorf("connect to mpv: %w", err)
	}

	i.conn = conn
	i.encoder = json.NewEncoder(conn)
	i.decoder = json.NewDecoder(conn)

	i.responses = make(chan Response)

	go i.readLoop()

	return nil
}

func (i *IPC) Close() error {
	if i.conn == nil {
		return nil
	}

	return i.conn.Close()
}

func (i *IPC) Send(command any) error {
	return i.encoder.Encode(command)
}

func (i *IPC) readLoop() {
	for {
		var resp Response

		if err := i.decoder.Decode(&resp); err != nil {
			close(i.responses)
			return
		}

		// Ignore MPV event messages for now.
		if resp.Event != "" {
			continue
		}

		i.responses <- resp
	}
}

func (i *IPC) Receive(v *Response) error {
	resp, ok := <-i.responses
	if !ok {
		return fmt.Errorf("ipc connection closed")
	}

	*v = resp
	return nil
}
