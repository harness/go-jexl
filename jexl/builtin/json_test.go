// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import "testing"

// Ensure JSONMarshal encodes a map to JSON.
func TestJSONMarshal_map(t *testing.T) {
	out, err := JSONMarshal(map[string]any{"k": "v"})
	if err != nil {
		t.Fatal(err)
	}
	if out != `{"k":"v"}` {
		t.Fatalf("unexpected output: %s", out)
	}
}

// Ensure JSONUnmarshal parses a JSON string into a map.
func TestJSONUnmarshal_basic(t *testing.T) {
	m, err := JSONUnmarshal(`{"a":1}`)
	if err != nil {
		t.Fatal(err)
	}
	if m["a"] != float64(1) {
		t.Fatalf("expected 1, got %v", m["a"])
	}
}

// Ensure JSONUnmarshal returns an error for invalid JSON.
func TestJSONUnmarshal_invalid(t *testing.T) {
	_, err := JSONUnmarshal(`not json`)
	if err == nil {
		t.Fatal("expected error")
	}
}

// Ensure JSONMarshalIndent produces indented JSON output.
func TestJSONMarshalIndent_basic(t *testing.T) {
	out, err := JSONMarshalIndent(map[string]any{"k": "v"}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	want := "{\n  \"k\": \"v\"\n}"
	if out != want {
		t.Fatalf("expected %q, got %q", want, out)
	}
}

// Ensure JSONSelect extracts a scalar value by plain path.
func TestJSONSelect_scalar(t *testing.T) {
	out, err := JSONSelect("name", `{"name":"alice","age":30}`)
	if err != nil {
		t.Fatal(err)
	}
	if out != "alice" {
		t.Fatalf("expected alice, got %v", out)
	}
}

// Ensure JSONSelect accepts a JSONPath expression with $ prefix.
func TestJSONSelect_jsonpath(t *testing.T) {
	out, err := JSONSelect("$.name", `{"name":"alice"}`)
	if err != nil {
		t.Fatal(err)
	}
	if out != "alice" {
		t.Fatalf("expected alice, got %v", out)
	}
}

// Ensure toGJSONPath strips the leading $. and leaves plain paths unchanged.
func TestToGJSONPath_plain(t *testing.T) {
	if got := toGJSONPath("name"); got != "name" {
		t.Fatalf("expected name, got %s", got)
	}
}

// Ensure toGJSONPath converts $.a.b to a.b.
func TestToGJSONPath_dotPath(t *testing.T) {
	if got := toGJSONPath("$.a.b"); got != "a.b" {
		t.Fatalf("expected a.b, got %s", got)
	}
}

// Ensure toGJSONPath converts array index syntax.
func TestToGJSONPath_arrayIndex(t *testing.T) {
	if got := toGJSONPath("$.items[0]"); got != "items.0" {
		t.Fatalf("expected items.0, got %s", got)
	}
}

// Ensure toGJSONPath converts wildcard array syntax.
func TestToGJSONPath_wildcard(t *testing.T) {
	if got := toGJSONPath("$.items[*].name"); got != "items.#.name" {
		t.Fatalf("expected items.#.name, got %s", got)
	}
}

// Ensure JSONSelect returns an error when the path is not found.
func TestJSONSelect_notFound(t *testing.T) {
	_, err := JSONSelect("missing", `{"name":"alice"}`)
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

// Ensure JSONList extracts an array by path.
func TestJSONList_array(t *testing.T) {
	out, err := JSONList("tags", `{"tags":["a","b","c"]}`)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(out))
	}
	if out[0] != "a" {
		t.Fatalf("expected a, got %v", out[0])
	}
}

// Ensure JSONList wraps a scalar result in a single-element slice.
func TestJSONList_scalar(t *testing.T) {
	out, err := JSONList("name", `{"name":"alice"}`)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0] != "alice" {
		t.Fatalf("expected [alice], got %v", out)
	}
}
