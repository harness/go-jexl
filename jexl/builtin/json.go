// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/harness/go-jexl/jexl/coerce"
	"github.com/tidwall/gjson"
)

// JSONMarshal returns the JSON encoding of v.
func JSONMarshal(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// JSONUnmarshal parses the JSON-encoded string v into a map.
func JSONUnmarshal(v any) (map[string]any, error) {
	var out map[string]any
	err := json.Unmarshal(
		[]byte(coerce.ToString(v)), &out,
	)
	return out, err
}

// JSONMarshalIndent returns the indented JSON encoding of v.
func JSONMarshalIndent(v any, prefix, indent any) (string, error) {
	b, err := json.MarshalIndent(
		v,
		coerce.ToString(prefix),
		coerce.ToString(indent),
	)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// JSONSelect extracts a single value from a JSON string using a gjson path expression.
// Returns the raw value (string, float64, bool, nil, or []any / map[string]any for objects/arrays).
func JSONSelect(path any, v any) (any, error) {
	result := gjson.Get(coerce.ToString(v), toGJSONPath(coerce.ToString(path)))
	if !result.Exists() {
		return nil, fmt.Errorf("path %q not found", path)
	}
	return result.Value(), nil
}

// JSONList extracts a list from a JSON string using a gjson path expression.
// Returns a []any; each element follows the same type mapping as JSONSelect.
func JSONList(path any, v any) ([]any, error) {
	result := gjson.Get(coerce.ToString(v), toGJSONPath(coerce.ToString(path)))
	if !result.Exists() {
		return nil, fmt.Errorf("path %q not found", path)
	}
	if result.IsArray() {
		items := result.Array()
		out := make([]any, len(items))
		for i, item := range items {
			out[i] = item.Value()
		}
		return out, nil
	}
	return []any{result.Value()}, nil
}

var reArrayIndex = regexp.MustCompile(`\[(\d+)\]`)

// toGJSONPath converts a JSONPath expression (e.g. $.a.b[0]) to gjson syntax.
// Paths that don't start with $ are returned unchanged.
func toGJSONPath(path string) string {
	if !strings.HasPrefix(path, "$") {
		return path
	}
	path = strings.TrimPrefix(path, "$")
	path = strings.TrimPrefix(path, ".")
	path = strings.ReplaceAll(path, "[*].", ".#.")
	path = strings.ReplaceAll(path, "[*]", "")
	path = reArrayIndex.ReplaceAllString(path, ".$1")
	path = strings.TrimPrefix(path, ".")
	return path
}
