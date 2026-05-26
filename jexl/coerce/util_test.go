// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package coerce

import "testing"

// Ensure DeepEqual returns true when both values are nil.
func TestDeepEqual_bothNil(t *testing.T) {
	if !DeepEqual(nil, nil) {
		t.Fatal("expected true")
	}
}

// Ensure DeepEqual returns false when one value is nil.
func TestDeepEqual_oneNil(t *testing.T) {
	if DeepEqual(nil, 1) {
		t.Fatal("expected false")
	}
}

// Ensure DeepEqual coerces int and int64 with the same value to equal.
func TestDeepEqual_intCoercion(t *testing.T) {
	if !DeepEqual(int(1), int64(1)) {
		t.Fatal("expected true")
	}
}

// Ensure DeepEqual returns false for int and int64 with different values.
func TestDeepEqual_intNotEqual(t *testing.T) {
	if DeepEqual(int(1), int64(2)) {
		t.Fatal("expected false")
	}
}

// Ensure DeepEqual coerces uint and uint64 with the same value to equal.
func TestDeepEqual_uintCoercion(t *testing.T) {
	if !DeepEqual(uint(5), uint64(5)) {
		t.Fatal("expected true")
	}
}

// Ensure DeepEqual coerces float32 and float64 with the same value to equal.
func TestDeepEqual_floatCoercion(t *testing.T) {
	if !DeepEqual(float32(1.5), float64(1.5)) {
		t.Fatal("expected true")
	}
}

// Ensure DeepEqual returns true for two identical strings.
func TestDeepEqual_stringEqual(t *testing.T) {
	if !DeepEqual("hello", "hello") {
		t.Fatal("expected true")
	}
}

// Ensure DeepEqual returns false for two different strings.
func TestDeepEqual_stringNotEqual(t *testing.T) {
	if DeepEqual("hello", "world") {
		t.Fatal("expected false")
	}
}

// Ensure DeepEqual returns true for two identical bool values.
func TestDeepEqual_boolEqual(t *testing.T) {
	if !DeepEqual(true, true) {
		t.Fatal("expected true")
	}
}

// Ensure DeepEqual returns false for two different bool values.
func TestDeepEqual_boolNotEqual(t *testing.T) {
	if DeepEqual(true, false) {
		t.Fatal("expected false")
	}
}

// Ensure DeepEqual returns true for two maps with the same keys and values.
func TestDeepEqual_mapEqual(t *testing.T) {
	a := map[string]any{"x": 1}
	b := map[string]any{"x": 1}
	if !DeepEqual(a, b) {
		t.Fatal("expected true")
	}
}

// Ensure DeepEqual returns false for two maps with different values for the same key.
func TestDeepEqual_mapValueDiffers(t *testing.T) {
	a := map[string]any{"x": 1}
	b := map[string]any{"x": 2}
	if DeepEqual(a, b) {
		t.Fatal("expected false")
	}
}

// Ensure DeepEqual returns false for two maps with different lengths.
func TestDeepEqual_mapLengthDiffers(t *testing.T) {
	a := map[string]any{"x": 1, "y": 2}
	b := map[string]any{"x": 1}
	if DeepEqual(a, b) {
		t.Fatal("expected false")
	}
}

// Ensure DeepEqual returns false when comparing a map against a non-map.
func TestDeepEqual_mapVsNonMap(t *testing.T) {
	if DeepEqual(map[string]any{"x": 1}, "x") {
		t.Fatal("expected false")
	}
}

// Ensure DeepEqual returns true for two slices with the same elements.
func TestDeepEqual_sliceEqual(t *testing.T) {
	a := []any{1, 2, 3}
	b := []any{1, 2, 3}
	if !DeepEqual(a, b) {
		t.Fatal("expected true")
	}
}

// Ensure DeepEqual returns false for two slices with different elements.
func TestDeepEqual_sliceElementDiffers(t *testing.T) {
	a := []any{1, 2, 3}
	b := []any{1, 2, 4}
	if DeepEqual(a, b) {
		t.Fatal("expected false")
	}
}

// Ensure DeepEqual returns false for two slices with different lengths.
func TestDeepEqual_sliceLengthDiffers(t *testing.T) {
	a := []any{1, 2}
	b := []any{1, 2, 3}
	if DeepEqual(a, b) {
		t.Fatal("expected false")
	}
}

// Ensure DeepEqual returns false when comparing a slice against a non-slice.
func TestDeepEqual_sliceVsNonSlice(t *testing.T) {
	if DeepEqual([]any{1}, "1") {
		t.Fatal("expected false")
	}
}

// Ensure DeepEqual compares nested maps recursively.
func TestDeepEqual_nestedMapEqual(t *testing.T) {
	a := map[string]any{"inner": map[string]any{"v": 1}}
	b := map[string]any{"inner": map[string]any{"v": 1}}
	if !DeepEqual(a, b) {
		t.Fatal("expected true")
	}
}

// Ensure DeepEqual compares nested slices recursively.
func TestDeepEqual_nestedSliceEqual(t *testing.T) {
	a := []any{[]any{1, 2}}
	b := []any{[]any{1, 2}}
	if !DeepEqual(a, b) {
		t.Fatal("expected true")
	}
}
