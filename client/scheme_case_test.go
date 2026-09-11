package client

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFetchAcceptsCaseInsensitiveScheme(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }))
	defer server.Close()
	for _, scheme := range []string{"HTTP://", "HtTp://"} {
		target := scheme + strings.TrimPrefix(server.URL, "http://")
		result, err := Fetch(Options{URL: target, Timeout: time.Second})
		if err != nil {
			t.Errorf("%s: %v", scheme, err)
			continue
		}
		result.Response.Body.Close()
	}
}
