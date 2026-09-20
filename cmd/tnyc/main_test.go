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

func TestRootCommands(t *testing.T) {
	want := map[string]bool{
		"api":      false,
		"database": false,
		"game":     false,
	}

	for _, command := range newRootCommand().Commands() {
		if _, ok := want[command.Name()]; ok {
			want[command.Name()] = command.Runnable()
		}
	}

	for name, runnable := range want {
		if !runnable {
			t.Errorf("root command %q is missing or not runnable", name)
		}
	}
}
