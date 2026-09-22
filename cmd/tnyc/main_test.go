package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mdhender/tnyc"
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
		"world":    false,
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

func TestWorldCommandDefaults(t *testing.T) {
	cmd := newWorldCommand()
	if !strings.Contains(cmd.Short, "Convert") || !strings.Contains(cmd.Short, "WGVC") || !strings.Contains(cmd.Short, "T'Nyc") {
		t.Errorf("Short = %q, want WGVC-to-T'Nyc conversion description", cmd.Short)
	}
	for name, want := range map[string]string{
		"input":  "var/wgvc-export.json",
		"output": "var/tnyc-world.json",
	} {
		flag := cmd.Flags().Lookup(name)
		if flag == nil {
			t.Fatalf("flag %q is missing", name)
		}
		if got := flag.DefValue; got != want {
			t.Errorf("flag %q default = %q, want %q", name, got, want)
		}
	}
}

func TestWorldCommandWritesJSON(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "tnyc-world.json")
	var stdout, stderr bytes.Buffer
	if err := run(context.Background(), []string{
		"world",
		"--input", filepath.Join("..", "..", "testdata", "wgvc-schema-v1.json"),
		"--output", outputPath,
	}, &stdout, &stderr); err != nil {
		t.Fatalf("run() error = %v", err)
	}
	if got := stderr.String(); got != "" {
		t.Errorf("stderr = %q, want empty", got)
	}
	if !strings.Contains(stdout.String(), outputPath) {
		t.Errorf("stdout = %q, want output path %q", stdout.String(), outputPath)
	}
	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("output file: %v", err)
	}
	if _, err := tnyc.LoadWorld(outputPath); err != nil {
		t.Fatalf("LoadWorld(output) error = %v", err)
	}
}

func TestDatabaseCreateCommand(t *testing.T) {
	databasePath := t.TempDir()
	worldMap := filepath.Join("..", "..", "testdata", "tnyc-world-v2.json")
	var stdout, stderr bytes.Buffer
	if err := run(context.Background(), []string{
		"database", "create",
		"--db-path", databasePath,
		"--world-map", worldMap,
	}, &stdout, &stderr); err != nil {
		t.Fatalf("run() error = %v", err)
	}
	if got := stderr.String(); got != "" {
		t.Errorf("stderr = %q, want empty", got)
	}
	if !strings.Contains(stdout.String(), databasePath) {
		t.Errorf("stdout = %q, want database path %q", stdout.String(), databasePath)
	}
	if _, err := os.Stat(filepath.Join(databasePath, "tnyc.json")); err != nil {
		t.Fatalf("database file: %v", err)
	}
}

func TestDatabaseCreateCommandDefaults(t *testing.T) {
	database := newDatabaseCommand()
	create, _, err := database.Find([]string{"create"})
	if err != nil {
		t.Fatal(err)
	}
	for name, want := range map[string]string{
		"db-path":   "var/",
		"world-map": "var/tnyc-world.json",
	} {
		flag := create.Flags().Lookup(name)
		if flag == nil {
			t.Fatalf("flag %q is missing", name)
		}
		if got := flag.DefValue; got != want {
			t.Errorf("flag %q default = %q, want %q", name, got, want)
		}
	}
}
