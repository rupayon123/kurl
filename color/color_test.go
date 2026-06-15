package color

import (
	"os"
	"testing"
)

func TestAutoEnabledHonorsNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	if AutoEnabled(os.Stdout) {
		t.Fatal("AutoEnabled should be false when NO_COLOR is set")
	}
}

func TestWrapDisabledReturnsPlainText(t *testing.T) {
	const value = "plain"
	if got := Wrap(false, Red, value); got != value {
		t.Fatalf("Wrap(false) = %q, want %q", got, value)
	}
}
