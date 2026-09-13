package printer

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestPrettyRespectsNoColorEnvironment(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	for _, tc := range []struct {
		name, input string
		format      func(io.Writer, io.Reader, bool) (int64, error)
	}{
		{"json", `{"value":42}`, PrettyJSON},
		{"html", `<div class="content">Hello</div>`, PrettyHTML},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			if _, err := tc.format(&out, strings.NewReader(tc.input), true); err != nil {
				t.Fatal(err)
			}
			if out.Len() == 0 {
				t.Fatal("empty output")
			}
			if strings.Contains(out.String(), "\x1b[") {
				t.Fatalf("NO_COLOR output contains ANSI codes: %q", out.String())
			}
		})
	}
}
