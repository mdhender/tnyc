// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package tnyc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

const wgvcSchemaVersion = 1

// ImportWGVC decodes and validates a WGVC schema-v1 document, then converts
// its game-relevant map data to a T'Nyc world.
func ImportWGVC(r io.Reader) (*World, error) {
	var document wgvcDocument
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&document); err != nil {
		return nil, fmt.Errorf("decode WGVC JSON: %w", err)
	}
	if err := rejectTrailingJSON(decoder); err != nil {
		return nil, err
	}
	if document.SchemaVersion != wgvcSchemaVersion {
		return nil, fmt.Errorf("unsupported WGVC schema version %d (want %d)", document.SchemaVersion, wgvcSchemaVersion)
	}
	if err := document.validate(); err != nil {
		return nil, fmt.Errorf("validate WGVC world: %w", err)
	}
	return document.world(), nil
}

// ImportWGVCFile imports a WGVC schema-v1 document from path.
func ImportWGVCFile(path string) (*World, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("import WGVC world %q: %w", path, err)
	}
	defer file.Close()
	world, err := ImportWGVC(file)
	if err != nil {
		return nil, fmt.Errorf("import WGVC world %q: %w", path, err)
	}
	return world, nil
}

type wgvcDocument struct {
	SchemaVersion int                  `json:"schema_version"`
	Generation    wgvcGenerationRecord `json:"generation"`
	Bounds        wgvcBoundsRecord     `json:"bounds"`
	Islands       []wgvcIslandRecord   `json:"islands"`
	Provinces     []wgvcProvinceRecord `json:"provinces"`
	Corners       []wgvcCornerRecord   `json:"corners"`
	Edges         []wgvcEdgeRecord     `json:"edges"`
}

type wgvcGenerationRecord struct {
	Result wgvcGenerationResultRecord `json:"result"`
}

type wgvcGenerationResultRecord struct {
	ProvinceCount int `json:"province_count"`
	IslandCount   int `json:"island_count"`
}

type wgvcBoundsRecord struct {
	Minimum pointRecord `json:"minimum"`
	Maximum pointRecord `json:"maximum"`
}

type wgvcIslandRecord struct {
	ID          IslandID `json:"id"`
	ProvinceIDs []CellID `json:"province_ids"`
}

type wgvcProvinceRecord struct {
	ID        CellID     `json:"id"`
	IslandID  IslandID   `json:"island_id"`
	CornerIDs []CornerID `json:"corner_ids"`
	Terrain   Terrain    `json:"terrain"`
}

type wgvcCornerRecord struct {
	ID    CornerID    `json:"id"`
	Point pointRecord `json:"point"`
}

type wgvcEdgeRecord struct {
	ID          EdgeID     `json:"id"`
	CornerIDs   []CornerID `json:"corner_ids"`
	ProvinceIDs []CellID   `json:"province_ids"`
}

func (document wgvcDocument) validate() error {
	if document.Bounds.Maximum.X <= document.Bounds.Minimum.X || document.Bounds.Maximum.Y <= document.Bounds.Minimum.Y {
		return errors.New("bounds must have positive width and height")
	}
	if got, want := document.Generation.Result.ProvinceCount, len(document.Provinces); got != want {
		return fmt.Errorf("generation result province_count is %d, want %d", got, want)
	}
	if got, want := document.Generation.Result.IslandCount, len(document.Islands); got != want {
		return fmt.Errorf("generation result island_count is %d, want %d", got, want)
	}

	converted := document.worldDocument()
	if err := converted.validate(); err != nil {
		return err
	}
	return nil
}

func (document wgvcDocument) worldDocument() worldDocument {
	converted := worldDocument{
		SchemaVersion: worldSchemaVersion,
		Islands:       make([]islandRecord, len(document.Islands)), Cells: make([]cellRecord, len(document.Provinces)),
		Corners: make([]cornerRecord, len(document.Corners)), Edges: make([]edgeRecord, len(document.Edges)),
	}
	for index, island := range document.Islands {
		converted.Islands[index] = islandRecord{ID: island.ID, CellIDs: island.ProvinceIDs}
	}
	for index, province := range document.Provinces {
		converted.Cells[index] = cellRecord{ID: province.ID, IslandID: province.IslandID, CornerIDs: province.CornerIDs, Terrain: province.Terrain}
	}
	for index, corner := range document.Corners {
		converted.Corners[index] = cornerRecord{ID: corner.ID, Point: corner.Point}
	}
	for index, edge := range document.Edges {
		converted.Edges[index] = edgeRecord{ID: edge.ID, CornerIDs: edge.CornerIDs, CellIDs: edge.ProvinceIDs}
	}
	return converted
}

func (document wgvcDocument) world() *World { return document.worldDocument().world() }
