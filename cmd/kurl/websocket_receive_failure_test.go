package main

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestWebSocketReceiveFailureExitsNonzero(t *testing.T) {
	if target := os.Getenv("KURL_WS_BROKEN_FRAME_TEST"); target != "" {
		runWebSocket(cliOptions{url: target, timeout: time.Second, noColor: true})
		os.Exit(0)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, buffer, err := w.(http.Hijacker).Hijack()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		sum := sha1.Sum([]byte(r.Header.Get("Sec-WebSocket-Key") + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
		_, _ = fmt.Fprintf(buffer, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: %s\r\n\r\n", base64.StdEncoding.EncodeToString(sum[:]))
		// Advertise a 64 MiB payload, exceeding the receive limit without allocating it.
		_, _ = buffer.Write([]byte{0x81, 127, 0, 0, 0, 0, 4, 0, 0, 0})
		_ = buffer.Flush()
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestWebSocketReceiveFailureExitsNonzero$")
	cmd.Env = append(os.Environ(), "KURL_WS_BROKEN_FRAME_TEST=ws"+strings.TrimPrefix(server.URL, "http"))
	input, writer, pipeErr := os.Pipe()
	if pipeErr != nil {
		t.Fatal(pipeErr)
	}
	defer input.Close()
	defer writer.Close()
	cmd.Stdin = input
	output, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		t.Fatal("client hung on receive failure")
	}
	if err == nil {
		t.Fatalf("broken frame reported success: %s", output)
	}
	if !strings.Contains(string(output), "frame payload size exceeds limit") {
		t.Fatalf("missing receive error: %s", output)
	}
}
