package main

import (
	"testing"
	"time"
)

func TestSavedRequestPreservesFilteringAndProtocol(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	original := cliOptions{method: "GET", url: "https://example.com", timeout: time.Second, filterQuery: ".items", filterKeys: "id,name", filterFlatten: true, http3: true}
	if err := saveRequestLocally("filtered", original); err != nil {
		t.Fatal(err)
	}
	got, err := loadRequestLocally("filtered")
	if err != nil {
		t.Fatal(err)
	}
	if got.filterQuery != original.filterQuery || got.filterKeys != original.filterKeys || !got.filterFlatten || !got.http3 {
		t.Fatalf("saved options lost: query=%q keys=%q flatten=%v http3=%v", got.filterQuery, got.filterKeys, got.filterFlatten, got.http3)
	}
}
