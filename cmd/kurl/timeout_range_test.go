package main

import (
	"testing"
	"time"
)

func TestTimeoutRejectsNonFiniteAndOverflowSeconds(t *testing.T) {
	for _, value := range []string{"NaN", "Inf", "+Inf", "-Inf", "1e100", "9223372037", "9223372036.854776"} {
		t.Run(value, func(t *testing.T) {
			if got, err := parseTimeout(value); err == nil {
				t.Fatalf("parseTimeout(%q) = %s without error", value, got)
			}
		})
	}
	for value, want := range map[string]time.Duration{"0": 0, "0.5": 500 * time.Millisecond, "30": 30 * time.Second, "9223372036": 9223372036 * time.Second, "500ms": 500 * time.Millisecond} {
		if got, err := parseTimeout(value); err != nil || got != want {
			t.Errorf("parseTimeout(%q) = %s, %v; want %s", value, got, err, want)
		}
	}
}
