package main

import (
	"golang.org/x/net/websocket"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestWebSocketSendsInputAboveScannerDefault(t *testing.T) {
	if url := os.Getenv("KURL_WS_INPUT_TEST"); url != "" {
		runWebSocket(cliOptions{url: url, noColor: true, timeout: time.Second})
		os.Exit(0)
	}
	received := make(chan string, 1)
	server := httptest.NewServer(websocket.Handler(func(ws *websocket.Conn) {
		var message string
		_ = websocket.Message.Receive(ws, &message)
		received <- message
	}))
	defer server.Close()
	payload := strings.Repeat("x", 70*1024)
	cmd := exec.Command(os.Args[0], "-test.run=^TestWebSocketSendsInputAboveScannerDefault$")
	cmd.Env = append(os.Environ(), "KURL_WS_INPUT_TEST=WS"+strings.TrimPrefix(server.URL, "http"))
	cmd.Stdin = strings.NewReader(payload + "\n")
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	select {
	case got := <-received:
		if got != payload {
			t.Fatalf("received %d bytes, want %d", len(got), len(payload))
		}
	case <-time.After(time.Second):
		t.Fatal("server did not receive input")
	}
}

func TestWebSocketSchemeCase(t *testing.T) {
	for _, raw := range []string{"ws://example.test", "WS://example.test", "wSs://example.test"} {
		if !isWebSocketURL(raw) {
			t.Errorf("missed WebSocket URL %q", raw)
		}
	}
	for _, raw := range []string{"http://example.test", "ws-not-a-scheme://example.test"} {
		if isWebSocketURL(raw) {
			t.Errorf("misclassified %q", raw)
		}
	}
}
