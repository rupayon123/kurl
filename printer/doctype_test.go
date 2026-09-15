package printer

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func documentTypeAttributes(t *testing.T, text string) []html.Attribute {
	t.Helper()
	doc, err := html.Parse(strings.NewReader(text))
	if err != nil {
		t.Fatal(err)
	}
	for n := doc.FirstChild; n != nil; n = n.NextSibling {
		if n.Type == html.DoctypeNode {
			return n.Attr
		}
	}
	t.Fatal("doctype missing")
	return nil
}

func TestPrettyHTMLPreservesDoctypeIdentifiers(t *testing.T) {
	for _, source := range []string{
		`<!DOCTYPE html><p>Body</p>`,
		`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd"><p>Body</p>`,
		`<!DOCTYPE html SYSTEM "about:legacy-compat"><p>Body</p>`,
	} {
		var output bytes.Buffer
		if _, err := PrettyHTML(&output, strings.NewReader(source), false); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(documentTypeAttributes(t, source), documentTypeAttributes(t, output.String())) {
			t.Fatalf("doctype identifiers lost: %s", output.String())
		}
	}
}
