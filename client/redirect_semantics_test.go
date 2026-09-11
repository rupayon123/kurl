package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchDoesNotFollowNonRedirectStatuses(t *testing.T) {
	for _, code := range []int{300, 304, 305, 306, 399} {
		t.Run(http.StatusText(code), func(t *testing.T) {
			visited := false
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/other" {
					visited = true
					w.WriteHeader(200)
					return
				}
				w.Header().Set("Location", "/other")
				w.WriteHeader(code)
			}))
			defer server.Close()
			result, err := Fetch(Options{URL: server.URL, Timeout: time.Second})
			if err != nil {
				t.Fatal(err)
			}
			defer result.Response.Body.Close()
			if visited || result.Response.StatusCode != code {
				t.Errorf("followed status %d", code)
			}
		})
	}
}

func TestHeadRemainsHeadAcrossSeeOther(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Location", "/other")
			w.WriteHeader(303)
			return
		}
		w.Header().Set("X-Method", r.Method)
	}))
	defer server.Close()
	result, err := Fetch(Options{URL: server.URL, Method: "HEAD", Timeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	defer result.Response.Body.Close()
	if result.Response.Header.Get("X-Method") != "HEAD" {
		t.Fatal("303 changed HEAD to GET")
	}
	_, _ = io.Copy(io.Discard, result.Response.Body)
}
