package graphql

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExecuteGraphQLDoesNotDuplicateContentType(t *testing.T) {
	for _, tc := range []struct {
		headers []string
		want    string
	}{
		{nil, "application/json"},
		{[]string{"content-type: application/json; charset=utf-8"}, "application/json; charset=utf-8"},
	} {
		received := make(chan []string, 1)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			received <- r.Header.Values("Content-Type")
			_, _ = w.Write([]byte(`{"data":{}}`))
		}))
		result, err := ExecuteGraphQL(Options{URL: server.URL, Query: "{ hello }", Headers: tc.headers})
		if err != nil {
			server.Close()
			t.Fatal(err)
		}
		result.Response.Body.Close()
		server.Close()
		values := <-received
		if len(values) != 1 || values[0] != tc.want {
			t.Errorf("got Content-Type %v, want only %q", values, tc.want)
		}
	}
}
