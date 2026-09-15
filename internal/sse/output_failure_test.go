package sse

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type failingEventWriter struct {
	calls   int
	failure error
}

func (w *failingEventWriter) Write(p []byte) (int, error) {
	w.calls++
	if w.calls > 1 {
		return 0, w.failure
	}
	return len(p), nil
}

func TestSSEStopsWhenEventOutputFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: hello\n\n"))
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	failure := errors.New("output unavailable")
	writer := &failingEventWriter{failure: failure}
	if err := runSSE(ctx, Options{URL: server.URL, NoColor: true}, writer); !errors.Is(err, failure) {
		t.Fatalf("got %v, want output failure", err)
	}
	if ctx.Err() != nil {
		t.Fatal("waited for context timeout instead of stopping on write failure")
	}
}
