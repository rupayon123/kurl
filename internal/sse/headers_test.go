package sse

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRunSSEHonorsHostHeader(t *testing.T) {
	got := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got <- r.Host
		w.Header().Set("Content-Type", "text/event-stream")
	}))
	defer srv.Close()
	if err := RunSSE(context.Background(), Options{URL: srv.URL, Headers: []string{"hOsT: events.example"}, NoColor: true}); err != nil {
		t.Fatal(err)
	}
	if host := <-got; host != "events.example" {
		t.Fatalf("Host=%q", host)
	}
}
func TestRunSSERejectsMalformedHeaderBeforeRequest(t *testing.T) {
	called := make(chan struct{}, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called <- struct{}{}
		w.Header().Set("Content-Type", "text/event-stream")
	}))
	defer srv.Close()
	if err := RunSSE(context.Background(), Options{URL: srv.URL, Headers: []string{"MissingColon"}, NoColor: true}); err == nil {
		t.Error("malformed header accepted")
	}
	select {
	case <-called:
		t.Error("request sent despite malformed header")
	default:
	}
}
