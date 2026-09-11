package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchUsesCustomHostHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Header().Set("X-Received-Host", r.Host) }))
	defer server.Close()
	result, err := Fetch(Options{URL: server.URL, Headers: []string{"Host: virtual.example"}, Timeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	defer result.Response.Body.Close()
	if got := result.Response.Header.Get("X-Received-Host"); got != "virtual.example" {
		t.Fatalf("got Host %q, want virtual.example", got)
	}
}
