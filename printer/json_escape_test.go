package printer

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestPrettyJSONEscapesDecodedStrings(t *testing.T) {
	for _, input := range []string{`{"a\"b":"line\nnext\t\\end"}`, `["\b\f\r\t","\u0000"]`, `{"雪":"<hello>&"}`} {
		var out bytes.Buffer
		if _, err := PrettyJSON(&out, strings.NewReader(input), false); err != nil {
			t.Fatal(err)
		}
		var want, got interface{}
		if err := json.Unmarshal([]byte(input), &want); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(out.Bytes(), &got); err != nil {
			t.Fatalf("formatter produced invalid JSON: %q: %v", out.String(), err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("round trip = %#v; want %#v", got, want)
		}
	}
}
