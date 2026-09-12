package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnvironmentPreservesExplicitMixedCaseSchemes(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	if err := os.Mkdir(filepath.Join(home, ".kurl"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".kurl", "environments.json"), []byte(`{"test":{"base_url":"https://example.invalid/v1"}}`), 0600); err != nil {
		t.Fatal(err)
	}
	for _, target := range []string{"HTTP://localhost/path", "HTTPS://localhost/path", "Ws://localhost/socket", "WSS://localhost/socket"} {
		t.Run(target, func(t *testing.T) {
			opts := cliOptions{env: "test", url: target}
			if err := applyEnvironment(&opts); err != nil {
				t.Fatal(err)
			}
			if opts.url != target {
				t.Fatalf("URL = %q, want explicit URL %q", opts.url, target)
			}
		})
	}
}
