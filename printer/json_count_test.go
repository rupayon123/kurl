package printer

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type failingJSONWriter struct {
	buffer    bytes.Buffer
	remaining int
	failure   error
}

func (w *failingJSONWriter) Write(p []byte) (int, error) {
	if len(p) > w.remaining {
		n, _ := w.buffer.Write(p[:w.remaining])
		w.remaining = 0
		return n, w.failure
	}
	n, _ := w.buffer.Write(p)
	w.remaining -= n
	return n, nil
}
func TestPrettyJSONReportsActualBytes(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	for _, enabled := range []bool{false, true} {
		var out bytes.Buffer
		n, err := PrettyJSON(&out, strings.NewReader(`{"雪":["quoted\"",true,null,42]}`), enabled)
		if err != nil {
			t.Fatal(err)
		}
		if n != int64(out.Len()) {
			t.Fatalf("color=%v reported %d bytes; wrote %d", enabled, n, out.Len())
		}
	}
}
func TestPrettyJSONCountsPartialFailedWrites(t *testing.T) {
	failure := errors.New("destination full")
	out := &failingJSONWriter{remaining: 7, failure: failure}
	n, err := PrettyJSON(out, strings.NewReader(`{"hello":"world"}`), false)
	if !errors.Is(err, failure) {
		t.Fatalf("error=%v", err)
	}
	if n != int64(out.buffer.Len()) {
		t.Fatalf("reported %d bytes; wrote %d", n, out.buffer.Len())
	}
}
