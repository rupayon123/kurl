package sse

import (
	"strings"
	"testing"
	"testing/iotest"
)

func TestEventStreamLineEndingsAndBOM(t *testing.T) {
	for _, ending := range []string{"\n", "\r\n", "\r"} {
		for _, bom := range []string{"", "\ufeff"} {
			t.Run(bom+ending, func(t *testing.T) {
				// One-byte reads cover CRLF and UTF-8 BOM split across network chunks.
				input := bom + strings.Join([]string{"data: first", "", "data: second", "", ""}, ending)
				var got []Event
				err := ParseStream(iotest.OneByteReader(strings.NewReader(input)), func(e Event) { got = append(got, e) })
				if err != nil || len(got) != 2 || got[0].Data != "first" || got[1].Data != "second" {
					t.Fatalf("got %+v, %v", got, err)
				}
			})
		}
	}
}
