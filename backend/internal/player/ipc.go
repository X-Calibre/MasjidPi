package player

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

type IPC struct {
	socket string

	conn    net.Conn
	encoder *json.Encoder
	decoder *json.Decoder

	responses chan Response
	requestMu sync.Mutex
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

	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)
	responses := make(chan Response)

	i.conn = conn
	i.encoder = encoder
	i.decoder = decoder
	i.responses = responses

	go i.readLoop(decoder, responses)

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

func (i *IPC) readLoop(decoder *json.Decoder, responses chan Response) {
	defer close(responses)

	for {
		var resp Response

		if err := decoder.Decode(&resp); err != nil {
			return
		}

		// Ignore MPV event messages for now.
		if resp.Event != "" {
			continue
		}

		responses <- resp
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

// RoundTrip keeps a command and its response paired on MPV's shared IPC
// connection. MPV replies are ordered, but multiple goroutines can issue
// commands concurrently, so an unlocked Send followed by Receive can allow
// one caller to consume another caller's response.
func (i *IPC) RoundTrip(command any, v *Response) error {
	i.requestMu.Lock()
	defer i.requestMu.Unlock()

	if err := i.Send(command); err != nil {
		return err
	}
	return i.Receive(v)
}
