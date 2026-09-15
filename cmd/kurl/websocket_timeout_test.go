package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestWebSocketHandshakeHonorsTimeout(t *testing.T) {
	if url := os.Getenv("KURL_WS_TIMEOUT_TEST"); url != "" {
		runWebSocket(cliOptions{url: url, timeout: 40 * time.Millisecond})
		os.Exit(0)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestWebSocketHandshakeHonorsTimeout$")
	cmd.Env = append(os.Environ(), "KURL_WS_TIMEOUT_TEST=ws"+strings.TrimPrefix(server.URL, "http"))
	out, err := cmd.CombinedOutput()
	if err == nil || ctx.Err() != nil || !strings.Contains(string(out), "deadline exceeded") {
		t.Fatalf("configured handshake timeout was not reported: err=%v context=%v output=%s", err, ctx.Err(), out)
	}
}
