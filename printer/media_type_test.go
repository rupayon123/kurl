package printer

import (
	"bytes"
	"github.com/kavix/kurl/client"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestRenderStructuredMediaTypes(t *testing.T) {
	for _, tc := range []struct{ media, body, want string }{
		{"application/problem+json", `{"message":"missing"}`, `"message"`},
		{"application/vnd.api+json; charset=utf-8", `{"message":"missing"}`, `"message"`},
		{"application/xhtml+xml", "<p>hello</p>", "<p>hello</p>"},
	} {
		t.Run(tc.media, func(t *testing.T) {
			var out bytes.Buffer
			r := &client.Result{Response: &http.Response{Header: http.Header{"Content-Type": {tc.media}}, Body: io.NopCloser(strings.NewReader(tc.body))}}
			if err := renderBody(&out, r, Options{}); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(out.String(), tc.want) {
				t.Fatalf("body hidden: %q", out.String())
			}
		})
	}
}
func TestMediaTypeParametersDoNotSelectFormatter(t *testing.T) {
	if isJSON(`text/plain; note="application/json"`, 0) {
		t.Fatal("parameter selected JSON formatter")
	}
	if isHTML(`text/plain; note="text/html"`) {
		t.Fatal("parameter selected HTML formatter")
	}
}
