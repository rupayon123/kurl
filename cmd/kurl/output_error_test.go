package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestRequestReportsBufferedOutputFailure(t *testing.T) {
	if url := os.Getenv("KURL_TEST_CLOSED_OUTPUT_URL"); url != "" {
		if err := os.Stdout.Close(); err != nil {
			t.Fatal(err)
		}
		runRequest(cliOptions{url: url, method: "GET", timeout: time.Second})
		os.Exit(0)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()
	cmd := exec.Command(os.Args[0], "-test.run=^TestRequestReportsBufferedOutputFailure$")
	cmd.Env = append(os.Environ(), "KURL_TEST_CLOSED_OUTPUT_URL="+srv.URL)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("request reported success despite a closed output destination")
	}
	if !strings.Contains(string(out), "file already closed") {
		t.Fatalf("unexpected failure: %s", out)
	}
}

func TestGraphQLReportsBufferedOutputFailure(t *testing.T) {
	if url := os.Getenv("KURL_TEST_GRAPHQL_CLOSED_OUTPUT_URL"); url != "" {
		if err := os.Stdout.Close(); err != nil {
			t.Fatal(err)
		}
		handleGraphQLCommand([]string{url, "--query", "{ id }"})
		os.Exit(0)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"id":1}}`))
	}))
	defer srv.Close()
	cmd := exec.Command(os.Args[0], "-test.run=^TestGraphQLReportsBufferedOutputFailure$")
	cmd.Env = append(os.Environ(), "KURL_TEST_GRAPHQL_CLOSED_OUTPUT_URL="+srv.URL)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("GraphQL succeeded despite closed output")
	}
	if !strings.Contains(string(out), "file already closed") {
		t.Fatalf("unexpected error: %s", out)
	}
}
