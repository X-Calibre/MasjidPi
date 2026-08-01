package player

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type IPC struct {
	socket string
	conn   net.Conn
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

	return nil
}

func (i *IPC) Close() error {
	if i.conn == nil {
		return nil
	}

	return i.conn.Close()
}

func (i *IPC) Send(command any) error {
	encoder := json.NewEncoder(i.conn)
	return encoder.Encode(command)
}

func (i *IPC) Receive(v any) error {
	decoder := json.NewDecoder(i.conn)
	return decoder.Decode(v)
}
