package client

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClosingResponseReleasesPerRequestTransport(t *testing.T) {
	closed := make(chan struct{}, 1)
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateClosed {
			select {
			case closed <- struct{}{}:
			default:
			}
		}
	}
	server.Start()
	defer server.Close()
	result, err := Fetch(Options{URL: server.URL, Timeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(io.Discard, result.Response.Body); err != nil {
		t.Fatal(err)
	}
	if err := result.Response.Body.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-closed:
	case <-time.After(300 * time.Millisecond):
		t.Fatal("unused per-request transport retained an idle connection")
	}
}
