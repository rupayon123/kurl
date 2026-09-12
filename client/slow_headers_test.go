package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchAllowsSlowHeadersWithinConfiguredTimeout(t *testing.T) {
	for _, timeout := range []time.Duration{0, 30 * time.Second} {
		t.Run(timeout.String(), func(t *testing.T) {
			t.Parallel()
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				select {
				case <-time.After(11 * time.Second):
					w.WriteHeader(http.StatusOK)
				case <-r.Context().Done():
				}
			}))
			defer srv.Close()
			result, err := Fetch(Options{URL: srv.URL, Timeout: timeout})
			if err != nil {
				t.Fatalf("headers within configured timeout %s rejected: %v", timeout, err)
			}
			defer result.Response.Body.Close()
			if result.Response.StatusCode != http.StatusOK {
				t.Fatalf("status = %d", result.Response.StatusCode)
			}
		})
	}
}
