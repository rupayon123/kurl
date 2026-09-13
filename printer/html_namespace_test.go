package printer

import (
	"bytes"
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func TestPrettyHTMLPreservesNamespacedAttributes(t *testing.T) {
	input := `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"><use xlink:href="#shape" xml:lang="en"></use></svg>`
	var out bytes.Buffer
	if _, err := PrettyHTML(&out, strings.NewReader(input), false); err != nil {
		t.Fatal(err)
	}
	for _, attribute := range []string{`xmlns:xlink="http://www.w3.org/1999/xlink"`, `xlink:href="#shape"`, `xml:lang="en"`} {
		if !strings.Contains(out.String(), attribute) {
			t.Errorf("missing qualified attribute %s in %s", attribute, out.String())
		}
	}
}
func TestRenderVoidTagPreservesAttributeNamespace(t *testing.T) {
	node := &html.Node{Type: html.ElementNode, Data: "br", Attr: []html.Attribute{{Namespace: "xml", Key: "lang", Val: "en"}}}
	if got := renderVoidTag(node, false); got != `<br xml:lang="en" />` {
		t.Fatalf("void tag=%q", got)
	}
}
