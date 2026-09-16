package main

import (
	"io"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

func TestWebSocketReportsClosedOutput(t *testing.T) {
	if target := os.Getenv("KURL_WS_CLOSED_OUTPUT_TEST"); target != "" {
		if err := os.Stdout.Close(); err != nil {
			t.Fatal(err)
		}
		runWebSocket(cliOptions{url: target, timeout: time.Second, noColor: true})
		os.Exit(0)
	}
	server := httptest.NewServer(websocket.Handler(func(ws *websocket.Conn) { _, _ = io.Copy(io.Discard, ws) }))
	defer server.Close()
	cmd := exec.Command(os.Args[0], "-test.run=^TestWebSocketReportsClosedOutput$")
	cmd.Env = append(os.Environ(), "KURL_WS_CLOSED_OUTPUT_TEST=ws"+strings.TrimPrefix(server.URL, "http"))
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("closed output reported success: %s", output)
	}
	if !strings.Contains(string(output), "file already closed") {
		t.Fatalf("unexpected failure: %s", output)
	}
}
