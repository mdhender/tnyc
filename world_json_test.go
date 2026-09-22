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

func TestImportWGVCConvertsMapData(t *testing.T) {
	world, err := ImportWGVCFile("testdata/wgvc-schema-v1.json")
	if err != nil {
		t.Fatalf("ImportWGVCFile() error = %v", err)
	}
	if got, want := len(world.Islands), 1; got != want {
		t.Errorf("len(Islands) = %d, want %d", got, want)
	}
	if got, want := world.Islands[0].ID, IslandID(1); got != want {
		t.Errorf("island ID = %d, want %d", got, want)
	}
	if got, want := world.Islands[0].CellIDs, []CellID{1}; !reflect.DeepEqual(got, want) {
		t.Errorf("island cells = %v, want %v", got, want)
	}
	if got, want := len(world.Cells), 2; got != want {
		t.Errorf("len(Cells) = %d, want %d", got, want)
	}
	if got, want := world.Cells[0].Terrain, Terrain("plains"); got != want {
		t.Errorf("land terrain = %q, want %q", got, want)
	}
	if got, want := world.Cells[0].ID, CellID(1); got != want {
		t.Errorf("land cell ID = %d, want %d", got, want)
	}
	if got, want := world.Cells[0].IslandID, IslandID(1); got != want {
		t.Errorf("land cell island ID = %d, want %d", got, want)
	}
	if got, want := world.Cells[0].CornerIDs, []CornerID{1, 2, 3}; !reflect.DeepEqual(got, want) {
		t.Errorf("ordered corners = %v, want %v", got, want)
	}
	if got := world.Cells[1].IslandID; got != NoIslandID {
		t.Errorf("water cell island ID = %d, want %d", got, NoIslandID)
	}
	if got, want := world.Edges[0].ID, EdgeID(1); got != want {
		t.Errorf("boundary edge ID = %d, want %d", got, want)
	}
	if got, want := world.Edges[0].CornerIDs, ([2]CornerID{1, 2}); got != want {
		t.Errorf("boundary edge corners = %v, want %v", got, want)
	}
	if got, want := world.Edges[0].CellIDs, []CellID{1}; !reflect.DeepEqual(got, want) {
		t.Errorf("boundary edge cells = %v, want %v", got, want)
	}
	if got, want := world.Edges[1].CellIDs, []CellID{1, 2}; !reflect.DeepEqual(got, want) {
		t.Errorf("interior edge cells = %v, want %v", got, want)
	}
	if got, want := world.Corners[3].ID, CornerID(4); got != want {
		t.Errorf("corner ID = %d, want %d", got, want)
	}
	if got, want := world.Corners[3].Point, (Point{X: 2, Y: 1}); got != want {
		t.Errorf("corner point = %#v, want %#v", got, want)
	}
}

func TestImportWGVCRejectsMalformedReferences(t *testing.T) {
	tests := []struct {
		name   string
		change func(map[string]any)
		want   string
	}{
		{"island province", func(doc map[string]any) {
			doc["islands"].([]any)[0].(map[string]any)["province_ids"] = []any{float64(9)}
		}, "out-of-range cell"},
		{"province island", func(doc map[string]any) {
			doc["provinces"].([]any)[0].(map[string]any)["island_id"] = float64(9)
		}, "out-of-range island"},
		{"province corner", func(doc map[string]any) {
			doc["provinces"].([]any)[0].(map[string]any)["corner_ids"] = []any{float64(9)}
		}, "out-of-range corner"},
		{"edge corner", func(doc map[string]any) {
			doc["edges"].([]any)[0].(map[string]any)["corner_ids"] = []any{float64(0), float64(9)}
		}, "out-of-range corner"},
		{"edge province", func(doc map[string]any) {
			doc["edges"].([]any)[0].(map[string]any)["province_ids"] = []any{float64(9)}
		}, "out-of-range cell"},
	}

	valid := readTestFile(t, "testdata/wgvc-schema-v1.json")
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
			_, err = ImportWGVC(bytes.NewReader(data))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("ImportWGVC() error = %v, want error containing %q", err, tt.want)
			}
		})
	}
}

func TestWorldRoundTripPreservesTNyCData(t *testing.T) {
	want, err := ImportWGVCFile("testdata/wgvc-schema-v1.json")
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := EncodeWorld(&output, want); err != nil {
		t.Fatalf("EncodeWorld() error = %v", err)
	}
	for _, omitted := range []string{
		"generation", "config", "result", "bounds", "provinces", "province_ids",
		"center", "elevation", "elevation_band", "heat", "heat_band", "moisture", "moisture_band", "relief",
	} {
		if strings.Contains(output.String(), `"`+omitted+`"`) {
			t.Errorf("encoded T'Nyc world contains WGVC-only field %q", omitted)
		}
	}
	for _, required := range []string{"cells", "cell_ids", "islands", "corners", "edges"} {
		if !strings.Contains(output.String(), `"`+required+`"`) {
			t.Errorf("encoded T'Nyc world lacks field %q", required)
		}
	}
	got, err := DecodeWorld(&output)
	if err != nil {
		t.Fatalf("DecodeWorld(encoded) error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("round trip lost world data: got %#v, want %#v", got, want)
	}
}

func TestDecodeWorldRejectsMalformedReferences(t *testing.T) {
	tests := []struct {
		name   string
		change func(map[string]any)
		want   string
	}{
		{"zero island id", func(doc map[string]any) {
			doc["islands"].([]any)[0].(map[string]any)["id"] = float64(0)
		}, "non-canonical id 0"},
		{"zero island cell reference", func(doc map[string]any) {
			doc["islands"].([]any)[0].(map[string]any)["cell_ids"] = []any{float64(0)}
		}, "out-of-range cell id 0"},
		{"island cell", func(doc map[string]any) {
			doc["islands"].([]any)[0].(map[string]any)["cell_ids"] = []any{float64(9)}
		}, "out-of-range cell"},
		{"inconsistent island membership", func(doc map[string]any) {
			doc["islands"].([]any)[0].(map[string]any)["cell_ids"] = []any{}
		}, "disagrees with island membership"},
		{"cell corner", func(doc map[string]any) {
			doc["cells"].([]any)[0].(map[string]any)["corner_ids"] = []any{float64(9)}
		}, "out-of-range corner"},
		{"zero cell id", func(doc map[string]any) {
			doc["cells"].([]any)[0].(map[string]any)["id"] = float64(0)
		}, "non-canonical id 0"},
		{"zero cell corner reference", func(doc map[string]any) {
			doc["cells"].([]any)[0].(map[string]any)["corner_ids"] = []any{float64(0)}
		}, "out-of-range corner id 0"},
		{"zero corner id", func(doc map[string]any) {
			doc["corners"].([]any)[0].(map[string]any)["id"] = float64(0)
		}, "non-canonical id 0"},
		{"edge endpoint count", func(doc map[string]any) {
			doc["edges"].([]any)[0].(map[string]any)["corner_ids"] = []any{float64(0)}
		}, "want 2"},
		{"zero edge id", func(doc map[string]any) {
			doc["edges"].([]any)[0].(map[string]any)["id"] = float64(0)
		}, "non-canonical id 0"},
		{"zero edge corner reference", func(doc map[string]any) {
			doc["edges"].([]any)[0].(map[string]any)["corner_ids"] = []any{float64(0), float64(2)}
		}, "out-of-range corner id 0"},
		{"zero edge cell reference", func(doc map[string]any) {
			doc["edges"].([]any)[0].(map[string]any)["cell_ids"] = []any{float64(0)}
		}, "out-of-range cell id 0"},
		{"edge cell", func(doc map[string]any) {
			doc["edges"].([]any)[0].(map[string]any)["cell_ids"] = []any{float64(9)}
		}, "out-of-range cell"},
	}

	valid := readTestFile(t, "testdata/tnyc-world-v2.json")
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

func TestCanonicalWorldMatchesWGVCExport(t *testing.T) {
	want, err := ImportWGVCFile("var/wgvc-export.json")
	if err != nil {
		t.Fatalf("ImportWGVCFile() error = %v", err)
	}
	got, err := LoadWorld("var/tnyc-world.json")
	if err != nil {
		t.Fatalf("LoadWorld() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("var/tnyc-world.json does not match the converted WGVC export")
	}
	if got, want := len(got.Cells), 4688; got != want {
		t.Errorf("len(Cells) = %d, want %d", got, want)
	}
	if got, want := len(got.Corners), 9378; got != want {
		t.Errorf("len(Corners) = %d, want %d", got, want)
	}
	if got, want := len(got.Edges), 14065; got != want {
		t.Errorf("len(Edges) = %d, want %d", got, want)
	}
}

func TestLoadAndSaveWorldIncludePathInErrors(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing", "world.json")
	if _, err := LoadWorld(missing); err == nil || !strings.Contains(err.Error(), missing) {
		t.Errorf("LoadWorld() error = %v, want path %q", err, missing)
	}
	world, err := LoadWorld("testdata/tnyc-world-v2.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := SaveWorld(missing, world); err == nil || !strings.Contains(err.Error(), missing) {
		t.Errorf("SaveWorld() error = %v, want path %q", err, missing)
	}
}

func readTestFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
