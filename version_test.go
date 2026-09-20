package tnyc

import "testing"

func TestVersion(t *testing.T) {
	if got, want := Version().String(), "0.1.7-alpha"; got != want {
		t.Fatalf("Version().String() = %q, want %q", got, want)
	}
}
