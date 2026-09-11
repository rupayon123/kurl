package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTimeoutCoversWholeRedirectChain(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			return
		case <-time.After(80 * time.Millisecond):
		}
		if len(r.URL.Path) < 4 {
			w.Header().Set("Location", r.URL.Path+"x")
			w.WriteHeader(302)
			return
		}
		w.WriteHeader(200)
	}))
	defer server.Close()
	result, err := Fetch(Options{URL: server.URL, Timeout: 200 * time.Millisecond})
	if result != nil {
		result.Response.Body.Close()
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("redirect chain outlived total timeout: %v", err)
	}
}

func TestTimeoutContextRemainsAliveForResponseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.(http.Flusher).Flush()
		select {
		case <-r.Context().Done():
			return
		case <-time.After(20 * time.Millisecond):
		}
		_, _ = w.Write([]byte("complete"))
	}))
	defer server.Close()
	result, err := Fetch(Options{URL: server.URL, Timeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	defer result.Response.Body.Close()
	data, err := io.ReadAll(result.Response.Body)
	if err != nil || string(data) != "complete" {
		t.Fatalf("body canceled before caller read: %q %v", data, err)
	}
}
