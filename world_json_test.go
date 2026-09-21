package tnyc

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoadWorld(t *testing.T) {
	world, err := LoadWorld("testdata/wgvc-schema-v1.json")
	if err != nil {
		t.Fatalf("LoadWorld() error = %v", err)
	}
	if got, want := world.Generation.Config.Seed, "0xffffffffffffffff"; got != want {
		t.Errorf("seed = %q, want %q", got, want)
	}
	if got, want := len(world.Islands), 1; got != want {
		t.Errorf("len(Islands) = %d, want %d", got, want)
	}
	if got, want := len(world.Provinces), 2; got != want {
		t.Errorf("len(Provinces) = %d, want %d", got, want)
	}
	if got, want := len(world.Corners), 4; got != want {
		t.Errorf("len(Corners) = %d, want %d", got, want)
	}
	if got, want := len(world.Edges), 3; got != want {
		t.Errorf("len(Edges) = %d, want %d", got, want)
	}
	if got := world.Provinces[1].IslandID; got != NoIslandID {
		t.Errorf("water province island ID = %d, want %d", got, NoIslandID)
	}
	if got, want := world.Edges[0].ProvinceIDs, []ProvinceID{0}; !equalProvinceIDs(got, want) {
		t.Errorf("boundary edge provinces = %v, want %v", got, want)
	}
	if got, want := world.Edges[1].ProvinceIDs, []ProvinceID{0, 1}; !equalProvinceIDs(got, want) {
		t.Errorf("interior edge provinces = %v, want %v", got, want)
	}
	if got, want := world.Provinces[0].Terrain, Terrain("plains"); got != want {
		t.Errorf("terrain = %q, want %q", got, want)
	}
	if got, want := world.Provinces[1].ElevationBand, ElevationBand("deep-water"); got != want {
		t.Errorf("elevation band = %q, want %q", got, want)
	}
}

func TestLoadWGVCExport(t *testing.T) {
	world, err := LoadWorld("var/wgvc-export.json")
	if err != nil {
		t.Fatalf("LoadWorld() error = %v", err)
	}
	if got, want := len(world.Provinces), 4688; got != want {
		t.Errorf("len(Provinces) = %d, want %d", got, want)
	}
	if got, want := len(world.Corners), 9378; got != want {
		t.Errorf("len(Corners) = %d, want %d", got, want)
	}
	if got, want := len(world.Edges), 14065; got != want {
		t.Errorf("len(Edges) = %d, want %d", got, want)
	}
}

func TestDecodeWorldRejectsInvalidDocuments(t *testing.T) {
	valid := readFixture(t)
	tests := []struct {
		name   string
		change func(map[string]any)
		want   string
	}{
		{"missing schema", func(doc map[string]any) { delete(doc, "schema_version") }, "schema version 0"},
		{"zero schema", func(doc map[string]any) { doc["schema_version"] = float64(0) }, "schema version 0"},
		{"unsupported schema", func(doc map[string]any) { doc["schema_version"] = float64(2) }, "schema version 2"},
		{"invalid bounds", func(doc map[string]any) { doc["bounds"].(map[string]any)["maximum"].(map[string]any)["x"] = float64(0) }, "positive width"},
		{"island id", func(doc map[string]any) { doc["islands"].([]any)[0].(map[string]any)["id"] = float64(1) }, "non-canonical id"},
		{"province id", func(doc map[string]any) { doc["provinces"].([]any)[1].(map[string]any)["id"] = float64(8) }, "non-canonical id"},
		{"corner id", func(doc map[string]any) { doc["corners"].([]any)[1].(map[string]any)["id"] = float64(8) }, "non-canonical id"},
		{"edge id", func(doc map[string]any) { doc["edges"].([]any)[1].(map[string]any)["id"] = float64(8) }, "non-canonical id"},
		{"island province reference", func(doc map[string]any) {
			doc["islands"].([]any)[0].(map[string]any)["province_ids"] = []any{float64(9)}
		}, "out-of-range province"},
		{"province island reference", func(doc map[string]any) { doc["provinces"].([]any)[0].(map[string]any)["island_id"] = float64(9) }, "out-of-range island"},
		{"province corner reference", func(doc map[string]any) {
			doc["provinces"].([]any)[0].(map[string]any)["corner_ids"] = []any{float64(9)}
		}, "out-of-range corner"},
		{"edge endpoint count", func(doc map[string]any) { doc["edges"].([]any)[0].(map[string]any)["corner_ids"] = []any{float64(0)} }, "want 2"},
		{"edge endpoint reference", func(doc map[string]any) {
			doc["edges"].([]any)[0].(map[string]any)["corner_ids"] = []any{float64(0), float64(9)}
		}, "out-of-range corner"},
		{"edge incidence empty", func(doc map[string]any) { doc["edges"].([]any)[0].(map[string]any)["province_ids"] = []any{} }, "want 1 or 2"},
		{"edge incidence three", func(doc map[string]any) {
			doc["edges"].([]any)[0].(map[string]any)["province_ids"] = []any{float64(0), float64(1), float64(0)}
		}, "want 1 or 2"},
		{"edge province reference", func(doc map[string]any) { doc["edges"].([]any)[0].(map[string]any)["province_ids"] = []any{float64(9)} }, "out-of-range province"},
		{"province result count", func(doc map[string]any) { result(doc)["province_count"] = float64(3) }, "province_count"},
		{"island result count", func(doc map[string]any) { result(doc)["island_count"] = float64(2) }, "island_count"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var document map[string]any
			if err := json.Unmarshal(valid, &document); err != nil {
				t.Fatal(err)
			}
			tt.change(document)
			data, err := json.Marshal(document)
			if err != nil {
				t.Fatal(err)
			}
			_, err = DecodeWorld(bytes.NewReader(data))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("DecodeWorld() error = %v, want error containing %q", err, tt.want)
			}
		})
	}
}

func TestDecodeWorldRejectsMalformedAndTrailingJSON(t *testing.T) {
	for _, test := range []struct {
		name string
		data []byte
		want string
	}{
		{"invalid", []byte(`{"schema_version":`), "unexpected EOF"},
		{"trailing value", append(readFixture(t), []byte("\n{}")...), "trailing JSON value"},
		{"invalid trailing data", append(readFixture(t), []byte("\n{")...), "trailing world JSON"},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := DecodeWorld(bytes.NewReader(test.data))
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Errorf("DecodeWorld() error = %v, want error containing %q", err, test.want)
			}
		})
	}
}

func TestDecodeWorldAcceptsUnknownFields(t *testing.T) {
	var document map[string]any
	if err := json.Unmarshal(readFixture(t), &document); err != nil {
		t.Fatal(err)
	}
	document["future_metadata"] = map[string]any{"value": true}
	data, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DecodeWorld(bytes.NewReader(data)); err != nil {
		t.Fatalf("DecodeWorld() error = %v", err)
	}
}

func TestEncodeWorldRoundTrip(t *testing.T) {
	want, err := LoadWorld("testdata/wgvc-schema-v1.json")
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := EncodeWorld(&output, want); err != nil {
		t.Fatalf("EncodeWorld() error = %v", err)
	}
	got, err := DecodeWorld(&output)
	if err != nil {
		t.Fatalf("DecodeWorld(encoded) error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("round trip lost world data: got %#v, want %#v", got, want)
	}
}

func TestLoadAndSaveWorldIncludePathInErrors(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing", "world.json")
	if _, err := LoadWorld(missing); err == nil || !strings.Contains(err.Error(), missing) {
		t.Errorf("LoadWorld() error = %v, want path %q", err, missing)
	}
	directory := t.TempDir()
	if _, err := LoadWorld(directory); err == nil || !strings.Contains(err.Error(), directory) {
		t.Errorf("LoadWorld(directory) error = %v, want path %q", err, directory)
	}
	world, err := LoadWorld("testdata/wgvc-schema-v1.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := SaveWorld(missing, world); err == nil || !strings.Contains(err.Error(), missing) {
		t.Errorf("SaveWorld() error = %v, want path %q", err, missing)
	}
}

func readFixture(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/wgvc-schema-v1.json")
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func result(document map[string]any) map[string]any {
	return document["generation"].(map[string]any)["result"].(map[string]any)
}

func equalProvinceIDs(a, b []ProvinceID) bool {
	if len(a) != len(b) {
		return false
	}
	for index := range a {
		if a[index] != b[index] {
			return false
		}
	}
	return true
}
