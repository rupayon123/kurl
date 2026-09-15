package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSavedRequestRejectsInvalidTimeout(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	dir := filepath.Join(home, ".kurl", "requests")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	for _, timeout := range []string{"-1s", "not-a-duration", "9999999999999999999999s"} {
		data, err := json.Marshal(savedRequest{URL: "https://example.test", Timeout: timeout})
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "bad.json"), data, 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := loadRequestLocally("bad"); err == nil {
			t.Errorf("accepted timeout %q", timeout)
		}
	}
	for _, timeout := range []time.Duration{0, time.Second} {
		if err := saveRequestLocally("valid", cliOptions{url: "https://example.test", timeout: timeout}); err != nil {
			t.Fatal(err)
		}
		got, err := loadRequestLocally("valid")
		if err != nil || got.timeout != timeout {
			t.Fatalf("valid timeout %v: %+v, %v", timeout, got, err)
		}
	}
}
