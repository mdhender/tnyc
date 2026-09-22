// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package tnyc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

const worldSchemaVersion = 2

// DecodeWorld decodes and validates a T'Nyc world document.
func DecodeWorld(r io.Reader) (*World, error) {
	var document worldDocument
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&document); err != nil {
		return nil, fmt.Errorf("decode world JSON: %w", err)
	}
	if err := rejectTrailingJSON(decoder); err != nil {
		return nil, err
	}
	if document.SchemaVersion != worldSchemaVersion {
		return nil, fmt.Errorf("unsupported world schema version %d (want %d)", document.SchemaVersion, worldSchemaVersion)
	}
	if err := document.validate(); err != nil {
		return nil, fmt.Errorf("validate world: %w", err)
	}
	return document.world(), nil
}

// LoadWorld loads a T'Nyc world document from path.
func LoadWorld(path string) (*World, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("load world %q: %w", path, err)
	}
	defer file.Close()
	world, err := DecodeWorld(file)
	if err != nil {
		return nil, fmt.Errorf("load world %q: %w", path, err)
	}
	return world, nil
}

// EncodeWorld validates world and writes it as T'Nyc schema-v2 JSON.
func EncodeWorld(w io.Writer, world *World) error {
	if world == nil {
		return errors.New("encode world: nil world")
	}
	document := newWorldDocument(world)
	if err := document.validate(); err != nil {
		return fmt.Errorf("encode world: validate world: %w", err)
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(document); err != nil {
		return fmt.Errorf("encode world JSON: %w", err)
	}
	return nil
}

// SaveWorld creates or truncates path and writes a T'Nyc world document.
func SaveWorld(path string, world *World) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("save world %q: %w", path, err)
	}
	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("save world %q: %w", path, closeErr)
		}
	}()
	if err := EncodeWorld(file, world); err != nil {
		return fmt.Errorf("save world %q: %w", path, err)
	}
	return nil
}

func rejectTrailingJSON(decoder *json.Decoder) error {
	var trailing any
	err := decoder.Decode(&trailing)
	if errors.Is(err, io.EOF) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("decode trailing world JSON: %w", err)
	}
	return errors.New("decode world JSON: trailing JSON value")
}

type worldDocument struct {
	SchemaVersion int            `json:"schema_version"`
	Islands       []islandRecord `json:"islands"`
	Cells         []cellRecord   `json:"cells"`
	Corners       []cornerRecord `json:"corners"`
	Edges         []edgeRecord   `json:"edges"`
}

type pointRecord struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type islandRecord struct {
	ID      IslandID `json:"id"`
	CellIDs []CellID `json:"cell_ids"`
}

type cellRecord struct {
	ID        CellID     `json:"id"`
	IslandID  IslandID   `json:"island_id"`
	CornerIDs []CornerID `json:"corner_ids"`
	Terrain   Terrain    `json:"terrain"`
}

type cornerRecord struct {
	ID    CornerID    `json:"id"`
	Point pointRecord `json:"point"`
}

type edgeRecord struct {
	ID        EdgeID     `json:"id"`
	CornerIDs []CornerID `json:"corner_ids"`
	CellIDs   []CellID   `json:"cell_ids"`
}

func (document worldDocument) validate() error {
	if len(document.Cells) == 0 {
		return errors.New("world must contain cells")
	}
	if len(document.Corners) == 0 {
		return errors.New("world must contain corners")
	}
	if len(document.Edges) == 0 {
		return errors.New("world must contain edges")
	}
	cellIslands := make([]IslandID, len(document.Cells))
	for index := range cellIslands {
		cellIslands[index] = NoIslandID
	}
	for index, island := range document.Islands {
		if island.ID != IslandID(index+1) {
			return fmt.Errorf("island index %d has non-canonical id %d", index, island.ID)
		}
		for _, cellID := range island.CellIDs {
			if !validID(int(cellID), len(document.Cells)) {
				return fmt.Errorf("island %d has out-of-range cell id %d", island.ID, cellID)
			}
			if cellIslands[cellID-1] != NoIslandID {
				return fmt.Errorf("cell %d belongs to multiple islands", cellID)
			}
			cellIslands[cellID-1] = island.ID
		}
	}
	for index, cell := range document.Cells {
		if cell.ID != CellID(index+1) {
			return fmt.Errorf("cell index %d has non-canonical id %d", index, cell.ID)
		}
		if cell.IslandID != NoIslandID && !validID(int(cell.IslandID), len(document.Islands)) {
			return fmt.Errorf("cell %d has out-of-range island id %d", cell.ID, cell.IslandID)
		}
		if cell.IslandID != cellIslands[index] {
			return fmt.Errorf("cell %d island id %d disagrees with island membership %d", cell.ID, cell.IslandID, cellIslands[index])
		}
		for _, cornerID := range cell.CornerIDs {
			if !validID(int(cornerID), len(document.Corners)) {
				return fmt.Errorf("cell %d has out-of-range corner id %d", cell.ID, cornerID)
			}
		}
	}
	for index, corner := range document.Corners {
		if corner.ID != CornerID(index+1) {
			return fmt.Errorf("corner index %d has non-canonical id %d", index, corner.ID)
		}
	}
	for index, edge := range document.Edges {
		if edge.ID != EdgeID(index+1) {
			return fmt.Errorf("edge index %d has non-canonical id %d", index, edge.ID)
		}
		if len(edge.CornerIDs) != 2 {
			return fmt.Errorf("edge %d has %d corner ids, want 2", edge.ID, len(edge.CornerIDs))
		}
		for _, cornerID := range edge.CornerIDs {
			if !validID(int(cornerID), len(document.Corners)) {
				return fmt.Errorf("edge %d has out-of-range corner id %d", edge.ID, cornerID)
			}
		}
		if len(edge.CellIDs) < 1 || len(edge.CellIDs) > 2 {
			return fmt.Errorf("edge %d has %d cell ids, want 1 or 2", edge.ID, len(edge.CellIDs))
		}
		for _, cellID := range edge.CellIDs {
			if !validID(int(cellID), len(document.Cells)) {
				return fmt.Errorf("edge %d has out-of-range cell id %d", edge.ID, cellID)
			}
		}
	}
	return nil
}

func validID(id, count int) bool { return id > 0 && id <= count }

func (document worldDocument) world() *World {
	world := &World{
		Islands: make([]Island, len(document.Islands)), Cells: make([]Cell, len(document.Cells)),
		Corners: make([]Corner, len(document.Corners)), Edges: make([]Edge, len(document.Edges)),
	}
	for index, island := range document.Islands {
		world.Islands[index] = Island(island)
	}
	for index, cell := range document.Cells {
		world.Cells[index] = Cell(cell)
	}
	for index, corner := range document.Corners {
		world.Corners[index] = Corner{ID: corner.ID, Point: Point(corner.Point)}
	}
	for index, edge := range document.Edges {
		world.Edges[index] = Edge{ID: edge.ID, CornerIDs: [2]CornerID{edge.CornerIDs[0], edge.CornerIDs[1]}, CellIDs: edge.CellIDs}
	}
	return world
}

func newWorldDocument(world *World) worldDocument {
	document := worldDocument{
		SchemaVersion: worldSchemaVersion,
		Islands:       make([]islandRecord, len(world.Islands)), Cells: make([]cellRecord, len(world.Cells)),
		Corners: make([]cornerRecord, len(world.Corners)), Edges: make([]edgeRecord, len(world.Edges)),
	}
	for index, island := range world.Islands {
		document.Islands[index] = islandRecord(island)
	}
	for index, cell := range world.Cells {
		document.Cells[index] = cellRecord(cell)
	}
	for index, corner := range world.Corners {
		document.Corners[index] = cornerRecord{ID: corner.ID, Point: pointRecord(corner.Point)}
	}
	for index, edge := range world.Edges {
		document.Edges[index] = edgeRecord{ID: edge.ID, CornerIDs: []CornerID{edge.CornerIDs[0], edge.CornerIDs[1]}, CellIDs: edge.CellIDs}
	}
	return document
}
