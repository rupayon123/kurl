package sse

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRunSSEValidatesContentTypeBeforeOpeningLog(t *testing.T) {
	for _, tc := range []struct {
		name, contentType string
		valid             bool
	}{
		{"event stream", "text/event-stream", true},
		{"charset parameter", "text/event-stream; charset=utf-8", true},
		{"JSON", "application/json", false},
		{"HTML", "text/html", false},
		{"missing", "", false},
		{"malformed", "text/event-stream; broken", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header()["Content-Type"] = []string{tc.contentType}
				_, _ = w.Write([]byte("data: hello\n\n"))
			}))
			defer server.Close()
			output := filepath.Join(t.TempDir(), "events.log")
			if err := os.WriteFile(output, []byte("existing log"), 0600); err != nil {
				t.Fatal(err)
			}
			err := RunSSE(context.Background(), Options{URL: server.URL, OutputFile: output, NoColor: true})
			if tc.valid {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			if err == nil {
				t.Error("expected invalid content type error")
			}
			data, readErr := os.ReadFile(output)
			if readErr != nil {
				t.Fatal(readErr)
			}
			if string(data) != "existing log" {
				t.Errorf("invalid response overwrote existing log: %q", data)
			}
		})
	}
}
