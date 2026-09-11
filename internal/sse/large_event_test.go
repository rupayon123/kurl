package sse

import (
	"strings"
	"testing"
)

func TestParseStreamLargeDataLine(t *testing.T) {
	data := strings.Repeat("x", 128*1024)
	var events []Event
	err := ParseStream(strings.NewReader("data: "+data+"\n\n"), func(event Event) { events = append(events, event) })
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("got %d events, want 1", len(events))
	}
	if events[0].Data != data {
		t.Fatal("large event data was not preserved")
	}
}

func TestParseStreamRejectsOversizedLine(t *testing.T) {
	err := ParseStream(strings.NewReader("data: "+strings.Repeat("x", 1024*1024)+"\n\n"), func(Event) { t.Error("oversized event was dispatched") })
	if err == nil {
		t.Fatal("expected bounded scanner to reject oversized line")
	}
}
