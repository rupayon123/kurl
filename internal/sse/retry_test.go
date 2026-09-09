package sse

import (
	"strings"
	"testing"
	"time"
)

func TestRetryRequiresNonnegativeIntegerMilliseconds(t *testing.T) {
	for _, value := range []string{"-1", "+1", "1.5", "1s", " 10", "", "999999999999999999999999"} {
		t.Run(value, func(t *testing.T) {
			var got []Event
			err := ParseStream(strings.NewReader("retry: 1000\nretry: "+value+"\ndata: test\n\n"), func(e Event) { got = append(got, e) })
			if err != nil || len(got) != 1 || got[0].Retry != time.Second {
				t.Fatalf("got %+v, %v", got, err)
			}
		})
	}
}

func TestRetryPersistsAcrossEventBoundaries(t *testing.T) {
	var got []Event
	err := ParseStream(strings.NewReader("retry: 1000\n\ndata: first\n\ndata: second\n\nretry: 0\ndata: third\n\n"), func(e Event) { got = append(got, e) })
	if err != nil || len(got) != 3 {
		t.Fatalf("got %+v, %v", got, err)
	}
	if got[0].Retry != time.Second || got[1].Retry != time.Second || got[2].Retry != 0 {
		t.Fatalf("retry state lost: %+v", got)
	}
}
