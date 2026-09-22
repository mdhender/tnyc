package tnyc

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateDatastore(t *testing.T) {
	dir := t.TempDir()
	worldMap := "testdata/tnyc-world-v2.json"

	if err := CreateDatastore(dir, worldMap); err != nil {
		t.Fatalf("CreateDatastore() error = %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "tnyc.json"))
	if err != nil {
		t.Fatal(err)
	}
	var got Datastore
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("decode datastore: %v", err)
	}
	if got.Version != "1" {
		t.Errorf("Version = %q, want %q", got.Version, "1")
	}
	if got.WorldMap != worldMap {
		t.Errorf("WorldMap = %q, want %q", got.WorldMap, worldMap)
	}
}

func TestCreateDatastoreRejectsInvalidInputs(t *testing.T) {
	worldMap := "testdata/tnyc-world-v2.json"
	tests := []struct {
		name      string
		prepare   func(*testing.T) (string, string)
		wantError string
	}{
		{
			name: "missing database path",
			prepare: func(t *testing.T) (string, string) {
				return filepath.Join(t.TempDir(), "missing"), worldMap
			},
			wantError: "no such file or directory",
		},
		{
			name: "database path is a file",
			prepare: func(t *testing.T) (string, string) {
				path := filepath.Join(t.TempDir(), "file")
				if err := os.WriteFile(path, nil, 0o600); err != nil {
					t.Fatal(err)
				}
				return path, worldMap
			},
			wantError: "not a directory",
		},
		{
			name: "database path is not writable",
			prepare: func(t *testing.T) (string, string) {
				path := t.TempDir()
				if err := os.Chmod(path, 0o500); err != nil {
					t.Fatal(err)
				}
				return path, worldMap
			},
			wantError: "not writable",
		},
		{
			name: "missing world map",
			prepare: func(t *testing.T) (string, string) {
				return t.TempDir(), filepath.Join(t.TempDir(), "missing.json")
			},
			wantError: "invalid world map",
		},
		{
			name: "invalid world map",
			prepare: func(t *testing.T) (string, string) {
				mapPath := filepath.Join(t.TempDir(), "world.json")
				if err := os.WriteFile(mapPath, []byte(`{"schema_version": 1}`), 0o600); err != nil {
					t.Fatal(err)
				}
				return t.TempDir(), mapPath
			},
			wantError: "invalid world map",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, mapPath := tt.prepare(t)
			err := CreateDatastore(dir, mapPath)
			if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Errorf("CreateDatastore() error = %v, want error containing %q", err, tt.wantError)
			}
			if info, statErr := os.Stat(dir); statErr == nil && info.IsDir() {
				if _, statErr := os.Stat(filepath.Join(dir, "tnyc.json")); !os.IsNotExist(statErr) {
					t.Errorf("tnyc.json exists after rejected input; stat error = %v", statErr)
				}
			}
		})
	}
}

func TestCreateDatastoreDoesNotOverwriteExistingDatabase(t *testing.T) {
	dir := t.TempDir()
	databasePath := filepath.Join(dir, "tnyc.json")
	want := []byte("existing database")
	if err := os.WriteFile(databasePath, want, 0o600); err != nil {
		t.Fatal(err)
	}

	err := CreateDatastore(dir, "testdata/tnyc-world-v2.json")
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("CreateDatastore() error = %v, want already exists error", err)
	}
	got, err := os.ReadFile(databasePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Errorf("existing database changed to %q, want %q", got, want)
	}
}
