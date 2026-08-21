package player

import (
	"encoding/json"
	"net"
	"sync"
	"testing"
	"time"
)

func TestIPCRoundTripSerializesConcurrentCommands(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	responses := make(chan Response)
	ipc := &IPC{
		conn:      client,
		encoder:   json.NewEncoder(client),
		decoder:   json.NewDecoder(client),
		responses: responses,
	}
	go ipc.readLoop(ipc.decoder, responses)

	const requests = 12

	var stateMu sync.Mutex
	outstanding := 0
	maxOutstanding := 0
	serverDone := make(chan struct{})

	go func() {
		defer close(serverDone)
		decoder := json.NewDecoder(server)
		encoder := json.NewEncoder(server)

		for n := 0; n < requests; n++ {
			var cmd Command
			if err := decoder.Decode(&cmd); err != nil {
				return
			}

			stateMu.Lock()
			outstanding++
			if outstanding > maxOutstanding {
				maxOutstanding = outstanding
			}
			stateMu.Unlock()

			// Leave enough time for other callers to reach RoundTrip. Without
			// request serialization, multiple commands become outstanding here.
			time.Sleep(10 * time.Millisecond)

			_ = encoder.Encode(Response{Data: n, Error: "success"})

			stateMu.Lock()
			outstanding--
			stateMu.Unlock()
		}
	}()

	start := make(chan struct{})
	var wg sync.WaitGroup
	errs := make(chan error, requests)

	for n := 0; n < requests; n++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			<-start

			var resp Response
			err := ipc.RoundTrip(Command{Command: []any{"get_property", n}}, &resp)
			errs <- err
		}(n)
	}

	close(start)
	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("RoundTrip failed: %v", err)
		}
	}

	<-serverDone

	stateMu.Lock()
	gotMax := maxOutstanding
	stateMu.Unlock()

	if gotMax != 1 {
		t.Fatalf("maximum outstanding IPC commands = %d, want 1", gotMax)
	}
}

func TestIPCReadLoopClosesOnlyItsOwnResponseChannel(t *testing.T) {
	client1, server1 := net.Pipe()
	client2, server2 := net.Pipe()
	defer server1.Close()
	defer client2.Close()
	defer server2.Close()

	responses1 := make(chan Response)
	responses2 := make(chan Response)

	ipc := &IPC{
		conn:      client1,
		encoder:   json.NewEncoder(client1),
		decoder:   json.NewDecoder(client1),
		responses: responses1,
	}
	go ipc.readLoop(ipc.decoder, responses1)

	// Simulate reconnecting before the old read loop exits.
	ipc.conn = client2
	ipc.encoder = json.NewEncoder(client2)
	ipc.decoder = json.NewDecoder(client2)
	ipc.responses = responses2
	go ipc.readLoop(ipc.decoder, responses2)

	_ = client1.Close()

	select {
	case _, ok := <-responses1:
		if ok {
			t.Fatal("old response channel unexpectedly produced a response")
		}
	case <-time.After(time.Second):
		t.Fatal("old response channel was not closed")
	}

	select {
	case _, ok := <-responses2:
		if !ok {
			t.Fatal("new response channel was closed by old read loop")
		}
	default:
	}
}
