// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package coerce

import "testing"

// Both nil values are equal.
func TestDeepEqual_bothNil(t *testing.T) {
	if !DeepEqual(nil, nil) {
		t.Fatal("expected true")
	}
}

// One nil and one non-nil are not equal.
func TestDeepEqual_oneNil(t *testing.T) {
	if DeepEqual(nil, 1) {
		t.Fatal("expected false")
	}
}

// int and int64 with the same value are equal via int coercion.
func TestDeepEqual_intCoercion(t *testing.T) {
	if !DeepEqual(int(1), int64(1)) {
		t.Fatal("expected true")
	}
}

// int and int64 with different values are not equal.
func TestDeepEqual_intNotEqual(t *testing.T) {
	if DeepEqual(int(1), int64(2)) {
		t.Fatal("expected false")
	}
}

// uint and uint64 with the same value are equal.
func TestDeepEqual_uintCoercion(t *testing.T) {
	if !DeepEqual(uint(5), uint64(5)) {
		t.Fatal("expected true")
	}
}

// float32 and float64 with the same value are equal.
func TestDeepEqual_floatCoercion(t *testing.T) {
	if !DeepEqual(float32(1.5), float64(1.5)) {
		t.Fatal("expected true")
	}
}

// Two identical strings are equal.
func TestDeepEqual_stringEqual(t *testing.T) {
	if !DeepEqual("hello", "hello") {
		t.Fatal("expected true")
	}
}

// Two different strings are not equal.
func TestDeepEqual_stringNotEqual(t *testing.T) {
	if DeepEqual("hello", "world") {
		t.Fatal("expected false")
	}
}

// Two identical bool values are equal.
func TestDeepEqual_boolEqual(t *testing.T) {
	if !DeepEqual(true, true) {
		t.Fatal("expected true")
	}
}

// Two different bool values are not equal.
func TestDeepEqual_boolNotEqual(t *testing.T) {
	if DeepEqual(true, false) {
		t.Fatal("expected false")
	}
}

// Two maps with the same keys and values are equal.
func TestDeepEqual_mapEqual(t *testing.T) {
	a := map[string]any{"x": 1}
	b := map[string]any{"x": 1}
	if !DeepEqual(a, b) {
		t.Fatal("expected true")
	}
}

// Two maps with different values for the same key are not equal.
func TestDeepEqual_mapValueDiffers(t *testing.T) {
	a := map[string]any{"x": 1}
	b := map[string]any{"x": 2}
	if DeepEqual(a, b) {
		t.Fatal("expected false")
	}
}

// Two maps with different lengths are not equal.
func TestDeepEqual_mapLengthDiffers(t *testing.T) {
	a := map[string]any{"x": 1, "y": 2}
	b := map[string]any{"x": 1}
	if DeepEqual(a, b) {
		t.Fatal("expected false")
	}
}

// A map compared against a non-map is not equal.
func TestDeepEqual_mapVsNonMap(t *testing.T) {
	if DeepEqual(map[string]any{"x": 1}, "x") {
		t.Fatal("expected false")
	}
}

// Two slices with the same elements are equal.
func TestDeepEqual_sliceEqual(t *testing.T) {
	a := []any{1, 2, 3}
	b := []any{1, 2, 3}
	if !DeepEqual(a, b) {
		t.Fatal("expected true")
	}
}

// Two slices with different elements are not equal.
func TestDeepEqual_sliceElementDiffers(t *testing.T) {
	a := []any{1, 2, 3}
	b := []any{1, 2, 4}
	if DeepEqual(a, b) {
		t.Fatal("expected false")
	}
}

// Two slices with different lengths are not equal.
func TestDeepEqual_sliceLengthDiffers(t *testing.T) {
	a := []any{1, 2}
	b := []any{1, 2, 3}
	if DeepEqual(a, b) {
		t.Fatal("expected false")
	}
}

// A slice compared against a non-slice is not equal.
func TestDeepEqual_sliceVsNonSlice(t *testing.T) {
	if DeepEqual([]any{1}, "1") {
		t.Fatal("expected false")
	}
}

// Nested maps are compared recursively.
func TestDeepEqual_nestedMapEqual(t *testing.T) {
	a := map[string]any{"inner": map[string]any{"v": 1}}
	b := map[string]any{"inner": map[string]any{"v": 1}}
	if !DeepEqual(a, b) {
		t.Fatal("expected true")
	}
}

// Nested slices are compared recursively.
func TestDeepEqual_nestedSliceEqual(t *testing.T) {
	a := []any{[]any{1, 2}}
	b := []any{[]any{1, 2}}
	if !DeepEqual(a, b) {
		t.Fatal("expected true")
	}
}
