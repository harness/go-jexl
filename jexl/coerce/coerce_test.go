// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package coerce

import "testing"

// Ensure ToInt converts numeric types to int.
func TestToInt(t *testing.T) {
	tests := []struct {
		in   any
		want int
	}{
		{int(3), 3},
		{int8(4), 4},
		{int16(5), 5},
		{int32(6), 6},
		{int64(7), 7},
		{uint(8), 8},
		{uint8(9), 9},
		{uint16(10), 10},
		{uint32(11), 11},
		{uint64(12), 12},
		{float32(1.9), 1},
		{float64(2.9), 2},
	}
	for _, tt := range tests {
		got := ToInt(tt.in)
		if got != tt.want {
			t.Fatalf("ToInt(%v) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

// Ensure ToInt64 converts numeric types to int64.
func TestToInt64(t *testing.T) {
	tests := []struct {
		in   any
		want int64
	}{
		{int(3), 3},
		{int8(4), 4},
		{int16(5), 5},
		{int32(6), 6},
		{int64(99), 99},
		{uint(8), 8},
		{uint8(9), 9},
		{uint16(10), 10},
		{uint32(11), 11},
		{uint64(12), 12},
		{float32(1.9), 1},
		{float64(4.7), 4},
	}
	for _, tt := range tests {
		got := ToInt64(tt.in)
		if got != tt.want {
			t.Fatalf("ToInt64(%v) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

// Ensure ToFloat64 converts numeric types to float64.
func TestToFloat64(t *testing.T) {
	tests := []struct {
		in   any
		want float64
	}{
		{int(3), 3.0},
		{int8(4), 4.0},
		{int16(5), 5.0},
		{int32(6), 6.0},
		{int64(7), 7.0},
		{uint(8), 8.0},
		{uint8(9), 9.0},
		{uint16(10), 10.0},
		{uint32(11), 11.0},
		{uint64(12), 12.0},
		{float32(1.5), float64(float32(1.5))},
		{float64(2.5), 2.5},
	}
	for _, tt := range tests {
		got := ToFloat64(tt.in)
		if got != tt.want {
			t.Fatalf("ToFloat64(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}
