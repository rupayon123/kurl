package printer

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/kavix/kurl/client"
)

func TestRenderXMLMediaTypesAsText(t *testing.T) {
	for _, media := range []string{"application/xml", "application/atom+xml; charset=utf-8", "image/svg+xml"} {
		t.Run(media, func(t *testing.T) {
			body := `<root><value>hello &amp; goodbye</value></root>`
			result := &client.Result{Response: &http.Response{Header: http.Header{"Content-Type": {media}}, Body: io.NopCloser(strings.NewReader(body))}}
			var output bytes.Buffer
			if err := renderBody(&output, result, Options{}); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(output.String(), body) {
				t.Fatalf("XML hidden: %s", output.String())
			}
		})
	}
	for _, media := range []string{"application/pdf", "image/png", "application/octet-stream"} {
		if !isBinary(media) {
			t.Errorf("binary type %s classified as text", media)
		}
	}
}
