package printer

import (
	"bytes"
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func rawTextElements(t *testing.T, input string) []string {
	t.Helper()
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	var values []string
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
			var s strings.Builder
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				s.WriteString(c.Data)
			}
			values = append(values, s.String())
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return values
}
func TestPrettyHTMLPreservesRawText(t *testing.T) {
	for _, input := range []string{"<script>const value = `\n  first\n\n second  \n`;</script>", "<style>\n  a::after { content: 'a & b'; }\n\n</style>", `<div><script>if (a < b && c > d) run();</script></div>`} {
		t.Run(input, func(t *testing.T) {
			var out bytes.Buffer
			if _, err := PrettyHTML(&out, strings.NewReader(input), false); err != nil {
				t.Fatal(err)
			}
			want, got := rawTextElements(t, input), rawTextElements(t, out.String())
			if len(want) != len(got) {
				t.Fatalf("elements: %v vs %v", got, want)
			}
			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("raw text changed: got %q, want %q", got[i], want[i])
				}
			}
		})
	}
}
