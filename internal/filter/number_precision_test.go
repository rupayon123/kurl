package filter

import "testing"

func TestFiltersPreserveJSONNumbers(t *testing.T) {
	for _, tc := range []struct {
		name string
		run  func() ([]byte, error)
		want string
	}{
		{"path", func() ([]byte, error) { return ApplyFilter([]byte(`{"id":9007199254740993}`), ".id") }, `9007199254740993`},
		{"keys", func() ([]byte, error) {
			return FilterKeys([]byte(`{"id":9007199254740993,"amount":0.1234567890123456789}`), "id,amount")
		}, `{"amount":0.1234567890123456789,"id":9007199254740993}`},
		{"flatten", func() ([]byte, error) { return FlattenArray([]byte(`[[9007199254740993],[1e400]]`)) }, `[9007199254740993,1e400]`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.run()
			if err != nil || string(got) != tc.want {
				t.Fatalf("got %s, %v; want %s", got, err, tc.want)
			}
		})
	}
}

func TestFiltersStillRejectTrailingJSON(t *testing.T) {
	for _, run := range []func() ([]byte, error){
		func() ([]byte, error) { return ApplyFilter([]byte(`{"id":1} {}`), ".id") },
		func() ([]byte, error) { return FilterKeys([]byte(`{"id":1} {}`), "id") },
		func() ([]byte, error) { return FlattenArray([]byte(`[] []`)) },
	} {
		if _, err := run(); err == nil {
			t.Fatal("accepted multiple JSON values")
		}
	}
}
