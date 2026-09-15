package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"testing"
)

func TestSubcommandsRejectUnexpectedArguments(t *testing.T) {
	if os.Getenv("KURL_ARGUMENT_TEST") == "1" {
		for i, arg := range os.Args {
			if arg == "--" {
				if os.Args[i+1] == "graphql" {
					handleGraphQLCommand(os.Args[i+2:])
				} else {
					handleSSECommand(os.Args[i+2:])
				}
				return
			}
		}
		t.Fatal("missing child arguments")
	}
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{}}`))
	}))
	defer server.Close()
	for _, command := range []string{"graphql", "sse"} {
		for _, unexpected := range []string{"--misspelled-option", "extra-endpoint"} {
			t.Run(command+"/"+unexpected, func(t *testing.T) {
				args := []string{"-test.run=^TestSubcommandsRejectUnexpectedArguments$", "--", command, server.URL}
				if command == "graphql" {
					args = append(args, "--query", "{ id }")
				}
				args = append(args, unexpected)
				cmd := exec.Command(os.Args[0], args...)
				cmd.Env = append(os.Environ(), "KURL_ARGUMENT_TEST=1")
				output, err := cmd.CombinedOutput()
				if err == nil {
					t.Fatalf("unexpected argument succeeded: %s", output)
				}
				if !strings.Contains(string(output), unexpected) {
					t.Fatalf("error omitted argument: %s", output)
				}
			})
		}
	}
	if got := requests.Load(); got != 0 {
		t.Errorf("invalid commands sent %d requests", got)
	}
}
