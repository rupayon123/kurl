package sse

import (
	"reflect"
	"strings"
	"testing"
)

func TestDataDispatchAndPersistentEventID(t *testing.T) {
	for _, tc := range []struct {
		name, input string
		want        []Event
	}{
		{"empty data", "data:\n\n", []Event{{Data: ""}}},
		{"leading empty data", "data:\ndata: second\n\n", []Event{{Data: "\nsecond"}}},
		{"id-only block", "id: 7\n\ndata: hello\n\ndata: again\n\n", []Event{{ID: "7", Data: "hello"}, {ID: "7", Data: "again"}}},
		{"event-only block", "event: stale\n\ndata: hello\n\n", []Event{{Data: "hello"}}},
		{"null id ignored", "id: 7\ndata: first\n\nid: bad\x00id\ndata: next\n\n", []Event{{ID: "7", Data: "first"}, {ID: "7", Data: "next"}}},
		{"empty id resets", "id: 7\ndata: first\n\nid:\ndata: next\n\n", []Event{{ID: "7", Data: "first"}, {Data: "next"}}},
		{"no dispatch at eof", "data: unfinished", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got []Event
			err := ParseStream(strings.NewReader(tc.input), func(e Event) { got = append(got, e) })
			if err != nil || !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %#v, %v; want %#v", got, err, tc.want)
			}
		})
	}
}
