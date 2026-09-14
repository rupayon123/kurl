package printer

import (
	"bytes"
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func preText(t *testing.T, input string) string {
	t.Helper()
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	var result string
	var visit func(*html.Node, bool)
	visit = func(n *html.Node, inPre bool) {
		inPre = inPre || (n.Type == html.ElementNode && n.Data == "pre")
		if inPre && n.Type == html.TextNode {
			result += n.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			visit(c, inPre)
		}
	}
	visit(doc, false)
	return result
}
func TestPrettyHTMLPreservesPreformattedBlockChildren(t *testing.T) {
	input := "<pre>  first\n<div>  second\n\n third </div>  last\n</pre>"
	var out bytes.Buffer
	if _, err := PrettyHTML(&out, strings.NewReader(input), false); err != nil {
		t.Fatal(err)
	}
	if got, want := preText(t, out.String()), preText(t, input); got != want {
		t.Fatalf("pre text changed: got %q want %q", got, want)
	}
}

func TestPrettyHTMLPreservesLeadingPreformattedNewline(t *testing.T) {
	input := "<pre>\n\ntext</pre>"
	var out bytes.Buffer
	if _, err := PrettyHTML(&out, strings.NewReader(input), false); err != nil {
		t.Fatal(err)
	}
	if preText(t, out.String()) != preText(t, input) {
		t.Fatalf("leading newline lost: %q", out.String())
	}
}
