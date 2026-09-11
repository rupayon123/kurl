package main

import (
	"reflect"
	"testing"
)

func TestMergeHeadersOverridesAllProfileValues(t *testing.T) {
	profile := []string{"Accept: old-one", "X-Other: keep", "accept: old-two"}
	cli := []string{"ACCEPT: new-one", "Accept: new-two"}
	got := mergeHeaders(profile, cli)
	want := []string{"ACCEPT: new-one", "Accept: new-two", "X-Other: keep"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if profile[0] != "Accept: old-one" {
		t.Fatal("mutated profile")
	}
}
