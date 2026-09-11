package response

import (
	"strings"
	"testing"
)

func TestParseVerboseKeepsFinalExchange(t *testing.T) {
	input := strings.Join([]string{
		"> GET /old HTTP/1.1", "> Host: old.example", "< HTTP/1.1 302 Found", "< Location: https://new.example/", "<", "redirect body",
		"> GET /new HTTP/2", "> Host: new.example", "> X-Protocol: foo HTTP/2", "< HTTP/2 200", "< Content-Type: text/plain", "<", "final body",
	}, "\n")
	got := Parse(input)
	if got.StatusLine != "HTTP/2 200" || got.RequestLine != "GET /new HTTP/2" || got.Body != "final body" {
		t.Errorf("mixed exchanges: %#v", got)
	}
	if len(got.Headers) != 1 || got.Headers[0].Name != "Content-Type" {
		t.Errorf("stale response headers: %v", got.Headers)
	}
	if len(got.RequestHeaders) != 2 || got.RequestHeaders[0].Value != "new.example" {
		t.Errorf("stale request headers: %v", got.RequestHeaders)
	}
}

func TestParseVerboseDiscardsInterimHeaders(t *testing.T) {
	got := Parse("< HTTP/1.1 103 Early Hints\n< Link: </style.css>\n<\n< HTTP/1.1 200 OK\n< Content-Type: text/plain\n<\nbody")
	if got.StatusLine != "HTTP/1.1 200 OK" || len(got.Headers) != 1 || got.Headers[0].Name != "Content-Type" || got.Body != "body" {
		t.Fatalf("mixed interim and final response: %#v", got)
	}
}
