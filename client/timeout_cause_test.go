package client

import (
	"context"
	"errors"
	"net/url"
	"testing"
	"time"
)

func TestAnnotatedTimeoutPreservesCause(t *testing.T) {
	original := &url.Error{Op: "Get", URL: "http://example.test", Err: context.DeadlineExceeded}
	err := annotateTimeout(original, time.Second)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal("deadline cause was lost")
	}
	var urlErr *url.Error
	if !errors.As(err, &urlErr) {
		t.Fatal("URL error metadata was lost")
	}
}
