package response

import "testing"

func TestParseRawPreservesBodyBytes(t *testing.T) {
	for _, body := range []string{"  first\r\n\r\nlast \r\n", "\n\n", "\t ", "{\"message\":\"ok\"}\r\n"} {
		for _, ending := range []string{"\r\n", "\n"} {
			input := "HTTP/1.1 200 OK" + ending + "Content-Type: text/plain" + ending + ending + body
			got := Parse(input)
			if got.Body != body {
				t.Errorf("body %q became %q with header separator %q", body, got.Body, ending)
			}
		}
	}
}
