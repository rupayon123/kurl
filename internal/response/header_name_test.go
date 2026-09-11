package response

import "testing"

func TestParseDoesNotConsumeJSONAsHeaders(t *testing.T) {
	for _, input := range []string{`{"message":"ok"}`, "{\n  \"message\": \"ok\"\n}\n", "not a header: body\n", ": value\n"} {
		got := Parse(input)
		if len(got.Headers) != 0 || got.Body != input {
			t.Errorf("input %q: headers=%v body=%q", input, got.Headers, got.Body)
		}
	}
}

func TestParseAcceptsTokenHeaderNames(t *testing.T) {
	got := Parse("HTTP/1.1 200 OK\nX-Trace_1: yes\nX!#$%&'*+-.^_`|~: token\n\nbody")
	if len(got.Headers) != 2 || got.Body != "body" {
		t.Fatalf("unexpected response: %#v", got)
	}
}
