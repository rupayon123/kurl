package printer

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrettyHTMLIgnoresTagNamesInCommentsAndRawText(t *testing.T) {
	for _, source := range []string{
		`<!-- <html><head><body> --><p>Fragment</p>`,
		`<script>const template = "<html><head><body>";</script>`,
		`<htmlish>Fragment</htmlish>`,
	} {
		var output bytes.Buffer
		if _, err := PrettyHTML(&output, strings.NewReader(source), false); err != nil {
			t.Fatal(err)
		}
		for _, wrapper := range []string{"</html>", "</head>", "</body>"} {
			if strings.Contains(output.String(), wrapper) {
				t.Fatalf("inserted wrapper %s: %s", wrapper, output.String())
			}
		}
	}
}
