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
	if got, want := stdout.String(), "0.1.7-alpha\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
	if got := stderr.String(); got != "" {
		t.Errorf("stderr = %q, want empty", got)
	}
}
