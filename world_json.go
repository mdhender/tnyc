// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package tnyc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

const worldSchemaVersion = 1

// DecodeWorld decodes and validates a wgvc schema-v1 world. Unknown fields
// are ignored so schema-v1 producers may add compatible metadata.
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

// LoadWorld opens path and delegates decoding to DecodeWorld.
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

// EncodeWorld validates world and writes it as wgvc-compatible schema-v1 JSON.
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

// SaveWorld creates or truncates path and writes world as schema-v1 JSON.
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
	SchemaVersion int              `json:"schema_version"`
	Generation    generationRecord `json:"generation"`
	Bounds        boundsRecord     `json:"bounds"`
	Islands       []islandRecord   `json:"islands"`
	Provinces     []provinceRecord `json:"provinces"`
	Corners       []cornerRecord   `json:"corners"`
	Edges         []edgeRecord     `json:"edges"`
}

type generationRecord struct {
	Config generationConfigRecord `json:"config"`
	Result generationResultRecord `json:"result"`
}

type generationConfigRecord struct {
	Seed               string    `json:"seed"`
	ProvinceCount      int       `json:"province_count"`
	IslandCount        int       `json:"island_count"`
	AspectRatio        string    `json:"aspect_ratio"`
	OceanFraction      float64   `json:"ocean_fraction"`
	EdgeBarrierWidth   float64   `json:"edge_barrier_width"`
	EdgeRamp           []float64 `json:"edge_ramp"`
	AttractantCount    int       `json:"attractant_count"`
	AttractantRamp     []float64 `json:"attractant_ramp"`
	AttractantJitter   float64   `json:"attractant_jitter"`
	SoftmaxTemperature float64   `json:"softmax_temperature"`
	ControlPenalty     float64   `json:"control_penalty"`
	MaxRounds          int       `json:"max_rounds"`
	Relaxations        int       `json:"relaxations"`
	PolarIceFraction   float64   `json:"polar_ice_fraction"`
	PeakChillFraction  float64   `json:"peak_chill_fraction"`
}

type generationResultRecord struct {
	ProvinceCount      int     `json:"province_count"`
	LandProvinceCount  int     `json:"land_province_count"`
	InitialIslandCount int     `json:"initial_island_count"`
	IslandCount        int     `json:"island_count"`
	MergeCount         int     `json:"merge_count"`
	RoundsAttempted    int     `json:"rounds_attempted"`
	OceanFraction      float64 `json:"ocean_fraction"`
}

type pointRecord struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type boundsRecord struct {
	Minimum pointRecord `json:"minimum"`
	Maximum pointRecord `json:"maximum"`
}

type islandRecord struct {
	ID          IslandID     `json:"id"`
	ProvinceIDs []ProvinceID `json:"province_ids"`
}

type provinceRecord struct {
	ID            ProvinceID    `json:"id"`
	IslandID      IslandID      `json:"island_id"`
	Center        pointRecord   `json:"center"`
	CornerIDs     []CornerID    `json:"corner_ids"`
	Terrain       Terrain       `json:"terrain"`
	Elevation     float64       `json:"elevation"`
	ElevationBand ElevationBand `json:"elevation_band"`
	Relief        float64       `json:"relief"`
	Heat          float64       `json:"heat"`
	HeatBand      HeatBand      `json:"heat_band"`
	Moisture      float64       `json:"moisture"`
	MoistureBand  MoistureBand  `json:"moisture_band"`
}

type cornerRecord struct {
	ID    CornerID    `json:"id"`
	Point pointRecord `json:"point"`
}

type edgeRecord struct {
	ID          EdgeID       `json:"id"`
	CornerIDs   []CornerID   `json:"corner_ids"`
	ProvinceIDs []ProvinceID `json:"province_ids"`
	Elevation   float64      `json:"elevation"`
}

func (document worldDocument) validate() error {
	if document.Bounds.Maximum.X <= document.Bounds.Minimum.X || document.Bounds.Maximum.Y <= document.Bounds.Minimum.Y {
		return errors.New("bounds must have positive width and height")
	}
	if got, want := document.Generation.Result.ProvinceCount, len(document.Provinces); got != want {
		return fmt.Errorf("generation result province_count is %d, want %d", got, want)
	}
	if got, want := document.Generation.Result.IslandCount, len(document.Islands); got != want {
		return fmt.Errorf("generation result island_count is %d, want %d", got, want)
	}
	for index, island := range document.Islands {
		if island.ID != IslandID(index) {
			return fmt.Errorf("island index %d has non-canonical id %d", index, island.ID)
		}
		for _, provinceID := range island.ProvinceIDs {
			if !validIndex(int(provinceID), len(document.Provinces)) {
				return fmt.Errorf("island %d has out-of-range province id %d", island.ID, provinceID)
			}
		}
	}
	for index, province := range document.Provinces {
		if province.ID != ProvinceID(index) {
			return fmt.Errorf("province index %d has non-canonical id %d", index, province.ID)
		}
		if province.IslandID != NoIslandID && !validIndex(int(province.IslandID), len(document.Islands)) {
			return fmt.Errorf("province %d has out-of-range island id %d", province.ID, province.IslandID)
		}
		for _, cornerID := range province.CornerIDs {
			if !validIndex(int(cornerID), len(document.Corners)) {
				return fmt.Errorf("province %d has out-of-range corner id %d", province.ID, cornerID)
			}
		}
	}
	for index, corner := range document.Corners {
		if corner.ID != CornerID(index) {
			return fmt.Errorf("corner index %d has non-canonical id %d", index, corner.ID)
		}
	}
	for index, edge := range document.Edges {
		if edge.ID != EdgeID(index) {
			return fmt.Errorf("edge index %d has non-canonical id %d", index, edge.ID)
		}
		if len(edge.CornerIDs) != 2 {
			return fmt.Errorf("edge %d has %d corner ids, want 2", edge.ID, len(edge.CornerIDs))
		}
		for _, cornerID := range edge.CornerIDs {
			if !validIndex(int(cornerID), len(document.Corners)) {
				return fmt.Errorf("edge %d has out-of-range corner id %d", edge.ID, cornerID)
			}
		}
		if len(edge.ProvinceIDs) < 1 || len(edge.ProvinceIDs) > 2 {
			return fmt.Errorf("edge %d has %d province ids, want 1 or 2", edge.ID, len(edge.ProvinceIDs))
		}
		for _, provinceID := range edge.ProvinceIDs {
			if !validIndex(int(provinceID), len(document.Provinces)) {
				return fmt.Errorf("edge %d has out-of-range province id %d", edge.ID, provinceID)
			}
		}
	}
	return nil
}

func validIndex(id, length int) bool { return id >= 0 && id < length }

func (document worldDocument) world() *World {
	world := &World{
		Generation: Generation{
			Config: GenerationConfig(document.Generation.Config),
			Result: GenerationResult(document.Generation.Result),
		},
		Bounds: Bounds{
			Minimum: Point(document.Bounds.Minimum),
			Maximum: Point(document.Bounds.Maximum),
		},
		Islands:   make([]Island, len(document.Islands)),
		Provinces: make([]Province, len(document.Provinces)),
		Corners:   make([]Corner, len(document.Corners)),
		Edges:     make([]Edge, len(document.Edges)),
	}
	for index, island := range document.Islands {
		world.Islands[index] = Island(island)
	}
	for index, province := range document.Provinces {
		world.Provinces[index] = Province{
			ID: province.ID, IslandID: province.IslandID, Center: Point(province.Center),
			CornerIDs: province.CornerIDs, Terrain: province.Terrain, Elevation: province.Elevation,
			ElevationBand: province.ElevationBand, Relief: province.Relief, Heat: province.Heat,
			HeatBand: province.HeatBand, Moisture: province.Moisture, MoistureBand: province.MoistureBand,
		}
	}
	for index, corner := range document.Corners {
		world.Corners[index] = Corner{ID: corner.ID, Point: Point(corner.Point)}
	}
	for index, edge := range document.Edges {
		world.Edges[index] = Edge{
			ID: edge.ID, CornerIDs: [2]CornerID{edge.CornerIDs[0], edge.CornerIDs[1]},
			ProvinceIDs: edge.ProvinceIDs, Elevation: edge.Elevation,
		}
	}
	return world
}

func newWorldDocument(world *World) worldDocument {
	document := worldDocument{
		SchemaVersion: worldSchemaVersion,
		Generation: generationRecord{
			Config: generationConfigRecord(world.Generation.Config),
			Result: generationResultRecord(world.Generation.Result),
		},
		Bounds:  boundsRecord{Minimum: pointRecord(world.Bounds.Minimum), Maximum: pointRecord(world.Bounds.Maximum)},
		Islands: make([]islandRecord, len(world.Islands)), Provinces: make([]provinceRecord, len(world.Provinces)),
		Corners: make([]cornerRecord, len(world.Corners)), Edges: make([]edgeRecord, len(world.Edges)),
	}
	for index, island := range world.Islands {
		document.Islands[index] = islandRecord(island)
	}
	for index, province := range world.Provinces {
		document.Provinces[index] = provinceRecord{
			ID: province.ID, IslandID: province.IslandID, Center: pointRecord(province.Center),
			CornerIDs: province.CornerIDs, Terrain: province.Terrain, Elevation: province.Elevation,
			ElevationBand: province.ElevationBand, Relief: province.Relief, Heat: province.Heat,
			HeatBand: province.HeatBand, Moisture: province.Moisture, MoistureBand: province.MoistureBand,
		}
	}
	for index, corner := range world.Corners {
		document.Corners[index] = cornerRecord{ID: corner.ID, Point: pointRecord(corner.Point)}
	}
	for index, edge := range world.Edges {
		document.Edges[index] = edgeRecord{
			ID: edge.ID, CornerIDs: []CornerID{edge.CornerIDs[0], edge.CornerIDs[1]},
			ProvinceIDs: edge.ProvinceIDs, Elevation: edge.Elevation,
		}
	}
	return document
}
