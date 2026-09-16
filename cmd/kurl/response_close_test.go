package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHeadersOnlyClosesUnreadResponse(t *testing.T) {
	disconnected := make(chan struct{})
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1000")
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		select {
		case <-r.Context().Done():
			close(disconnected)
		case <-release:
		}
	}))
	defer server.Close()
	defer close(release)
	runRequest(cliOptions{url: server.URL, method: "GET", headersOnly: true, noColor: true})
	select {
	case <-disconnected:
	case <-time.After(time.Second):
		t.Fatal("headers-only request returned without closing its unread response")
	}
}
