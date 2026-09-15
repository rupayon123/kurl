package filter

import "testing"

func TestFilterQuotedObjectFields(t *testing.T) {
	data := []byte(`{"a.b":{"x[y]":[{"q\"uote":7}]},"":9,"unicode☃":true}`)
	for _, test := range []struct{ query, want string }{
		{`.["a.b"]["x[y]"][0]["q\"uote"]`, `7`},
		{`.["a.b"]["x[y]"][]["q\"uote"]`, `[7]`},
		{`[""]`, `9`},
		{`.["unicode\u2603"]`, `true`},
	} {
		got, err := ApplyFilter(data, test.query)
		if err != nil || string(got) != test.want {
			t.Errorf("%s: got %s, %v; want %s", test.query, got, err, test.want)
		}
	}
	for _, query := range []string{`.["a"]junk`, `.["a"`, `.["a\q"]`, `.["a" 0]`} {
		if _, err := ApplyFilter(data, query); err == nil {
			t.Errorf("accepted malformed query %s", query)
		}
	}
}
