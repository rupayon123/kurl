package response

import (
	"net/http"
	"testing"
)

func TestFromHTTPResponsePreservesWhitespaceOnlyBody(t *testing.T) {
	for _, body := range []string{" ", "\t\r\n", "\n\n"} {
		result := FromHTTPResponse("HTTP/1.1 200 OK", http.Header{"Content-Type": []string{"text/plain"}}, []byte(body))
		if result.Body != body {
			t.Errorf("body %q became %q", body, result.Body)
		}
		if len(result.Warnings) != 0 {
			t.Errorf("nonempty body reported empty: %v", result.Warnings)
		}
	}
	empty := FromHTTPResponse("HTTP/1.1 204 No Content", nil, nil)
	if empty.Body != "" || len(empty.Warnings) == 0 {
		t.Fatal("truly empty body should retain its warning")
	}
}
