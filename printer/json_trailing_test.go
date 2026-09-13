package printer

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

type jsonTailErrorReader struct {
	done    bool
	failure error
}

func (r *jsonTailErrorReader) Read(p []byte) (int, error) {
	if !r.done {
		r.done = true
		return copy(p, `{"ok":true}`), nil
	}
	return 0, r.failure
}
func TestPrettyJSONRejectsTrailingContent(t *testing.T) {
	for _, input := range []string{`{} {}`, `[] garbage`, `true false`, `{"a":1}x`} {
		var out bytes.Buffer
		if _, err := PrettyJSON(&out, strings.NewReader(input), false); err == nil {
			t.Fatalf("accepted trailing content: %q", input)
		}
	}
}
func TestPrettyJSONAcceptsTrailingWhitespace(t *testing.T) {
	var out bytes.Buffer
	if _, err := PrettyJSON(&out, strings.NewReader("{} \n\r\t"), false); err != nil {
		t.Fatal(err)
	}
}
func TestPrettyJSONPropagatesTrailingReadError(t *testing.T) {
	failure := errors.New("response truncated")
	if _, err := PrettyJSON(io.Discard, &jsonTailErrorReader{failure: failure}, false); !errors.Is(err, failure) {
		t.Fatalf("error=%v; want %v", err, failure)
	}
}
