package sse

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunSSEFiltersEffectiveEventType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: default-data\n\nevent: update\ndata: lower-data\n\nevent: Update\ndata: upper-data\n\nevent:\ndata: empty-type-data\n\n"))
	}))
	defer server.Close()
	for _, tc := range []struct {
		filter       string
		want, absent []string
	}{
		{"update", []string{"lower-data"}, []string{"default-data", "upper-data", "empty-type-data"}},
		{"message", []string{"default-data", "empty-type-data"}, []string{"lower-data", "upper-data"}},
		{"", []string{"default-data", "empty-type-data", "lower-data", "upper-data"}, nil},
	} {
		t.Run(tc.filter, func(t *testing.T) {
			output := filepath.Join(t.TempDir(), "events.log")
			if err := RunSSE(context.Background(), Options{URL: server.URL, FilterType: tc.filter, OutputFile: output, NoColor: true}); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(output)
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range tc.want {
				if !strings.Contains(string(data), want) {
					t.Errorf("missing %q in %s", want, data)
				}
			}
			for _, absent := range tc.absent {
				if strings.Contains(string(data), absent) {
					t.Errorf("unexpected %q in %s", absent, data)
				}
			}
		})
	}
}
