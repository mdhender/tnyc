// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package tnyc

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const datastoreSchemaVersion = "1"

// Datastore identifies the schema version and world map used by a game store.
type Datastore struct {
	Version  string
	WorldMap string
}

// CreateDatastore validates worldMap and exclusively creates tnyc.json in dir.
func CreateDatastore(dir, worldMap string) (err error) {
	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("create datastore: database path %q: %w", dir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("create datastore: database path %q is not a directory", dir)
	}
	if info.Mode().Perm()&0222 == 0 {
		return fmt.Errorf("create datastore: database path %q is not writable", dir)
	}

	databasePath := filepath.Join(dir, "tnyc.json")
	if _, err := os.Lstat(databasePath); err == nil {
		return fmt.Errorf("create datastore: database %q already exists", databasePath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("create datastore: inspect %q: %w", databasePath, err)
	}

	if _, err := LoadWorld(worldMap); err != nil {
		return fmt.Errorf("create datastore: invalid world map: %w", err)
	}

	data, err := json.MarshalIndent(Datastore{
		Version:  datastoreSchemaVersion,
		WorldMap: worldMap,
	}, "", "  ")
	if err != nil {
		return fmt.Errorf("create datastore: encode %q: %w", databasePath, err)
	}
	data = append(data, '\n')

	file, err := os.OpenFile(databasePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o666)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("create datastore: database %q already exists", databasePath)
		}
		return fmt.Errorf("create datastore: create %q: %w", databasePath, err)
	}
	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("create datastore: close %q: %w", databasePath, closeErr)
		}
		if err != nil {
			_ = os.Remove(databasePath)
		}
	}()
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("create datastore: write %q: %w", databasePath, err)
	}
	return nil
}
