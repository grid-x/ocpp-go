package ws

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

// The error channel is created lazily and closed by Stop, while the connection
// goroutines report errors on it. Run with -race: sending while Stop closes the
// channel used to be a data race and could panic with "send on closed channel".
func TestServerErrorsWhileStopping(t *testing.T) {
	server := NewServer()
	errs := server.Errors()
	go func() {
		for range errs {
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				server.error(errors.New("write failed"))
			}
		}()
	}
	server.Stop()
	wg.Wait()
}

func TestClientErrorsWhileStopping(t *testing.T) {
	client := NewClient()
	client.SetRequestedSubProtocol("ocpp1.6")
	errs := client.Errors()
	go func() {
		for range errs {
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				client.error(errors.New("write failed"))
			}
		}()
	}
	client.Stop()
	wg.Wait()
}

// error() must not block when nobody drains the channel, otherwise Stop, which
// needs the same lock to close it, would deadlock.
func TestServerErrorDoesNotBlockStop(t *testing.T) {
	server := NewServer()
	_ = server.Errors() // never drained
	for i := 0; i < 10; i++ {
		server.error(errors.New("write failed"))
	}
	server.Stop()
}

func TestClientErrorDoesNotBlockStop(t *testing.T) {
	client := NewClient()
	client.SetRequestedSubProtocol("ocpp1.6")
	_ = client.Errors() // never drained
	for i := 0; i < 10; i++ {
		client.error(errors.New("write failed"))
	}
	client.Stop()
}

// startTestServer starts a server on an OS-chosen port and returns its URL.
// Addr() is polled through its accessor, which is what callers do.
func startTestServer(t *testing.T) (*Server, string) {
	t.Helper()
	server := NewServer()
	server.SetMessageHandler(func(ws Channel, data []byte) error { return nil })
	go server.Start(0, "/{ws}")
	deadline := time.Now().Add(5 * time.Second)
	for server.Addr() == nil {
		if time.Now().After(deadline) {
			t.Fatal("server did not start")
		}
		time.Sleep(time.Millisecond)
	}
	return server, fmt.Sprintf("ws://localhost:%d/cp1", server.Addr().Port)
}

// Write must not panic or race when the connection is torn down while a send is
// in flight: the client used to close outQueue in cleanup while Write sent on it.
func TestClientWriteWhileDisconnecting(t *testing.T) {
	server, url := startTestServer(t)
	defer server.Stop()

	client := NewClient()
	client.SetRequestedSubProtocol("ocpp1.6")
	require.NoError(t, client.Start(url))

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 500; j++ {
				_ = client.Write([]byte("x"))
			}
		}()
	}
	// The server drops the connection, so the client's write pump runs cleanup
	// on its own, with no Stop involved.
	require.Eventually(t, func() bool {
		return server.StopConnection("cp1", websocket.CloseError{Code: websocket.CloseNormalClosure}) == nil
	}, 5*time.Second, time.Millisecond, "server never saw the client connection")
	wg.Wait()
	client.Stop()
}

// The handlers are read by the connection goroutines and swapped by callers
// (ocppj.Client.Stop replaces the disconnect handler).
func TestClientHandlerSwapWhileConnected(t *testing.T) {
	server, url := startTestServer(t)
	defer server.Stop()

	client := NewClient()
	client.SetRequestedSubProtocol("ocpp1.6")
	client.SetDisconnectedHandler(func(err error) {})
	client.SetReconnectedHandler(func() {})
	require.NoError(t, client.Start(url))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 200; j++ {
			client.SetDisconnectedHandler(func(err error) {})
			client.SetReconnectedHandler(func() {})
			client.SetMessageHandler(func(data []byte) error { return nil })
		}
	}()
	require.Eventually(t, func() bool {
		return server.StopConnection("cp1", websocket.CloseError{Code: websocket.CloseNormalClosure}) == nil
	}, 5*time.Second, time.Millisecond, "server never saw the client connection")
	wg.Wait()
	client.Stop()
}
