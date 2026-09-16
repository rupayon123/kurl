package sse

import (
	"strings"
	"testing"
)

func TestParseStreamReplacesInvalidUTF8(t *testing.T) {
	var events []Event
	err := ParseStream(strings.NewReader("data: valid ☃ invalid \xff\xfe\n\n"), func(event Event) { events = append(events, event) })
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 || events[0].Data != "valid ☃ invalid \ufffd\ufffd" {
		t.Fatalf("invalid UTF-8 was not decoded with replacement: %#v", events)
	}
}
