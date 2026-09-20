package main

import (
	"bytes"
	"context"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer

	if err := run(context.Background(), []string{"version"}, &stdout, &stderr); err != nil {
		t.Fatalf("run() error = %v", err)
	}
	if got := stderr.String(); got != "" {
		t.Errorf("stderr = %q, want empty", got)
	}
}
