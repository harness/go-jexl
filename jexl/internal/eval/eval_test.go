// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package eval

import (
	"math"
	"math/big"
	"testing"

	"github.com/harness/go-jexl/jexl/classes/java/util"
	"github.com/harness/go-jexl/jexl/internal/decimal"
)

// Ensure Negate negates each numeric type; !bool for booleans.
func TestNegate(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want any
	}{
		{"int", int(5), int(-5)},
		{"int8", int8(3), int8(-3)},
		{"int16", int16(10), int16(-10)},
		{"int32", int32(7), int32(-7)},
		{"int64", int64(100), int64(-100)},
		{"uint", uint(4), func() any { var v uint = 4; return -v }()},
		{"uint8", uint8(2), func() any { var v uint8 = 2; return -v }()},
		{"uint16", uint16(3), func() any { var v uint16 = 3; return -v }()},
		{"uint32", uint32(5), func() any { var v uint32 = 5; return -v }()},
		{"uint64", uint64(6), func() any { var v uint64 = 6; return -v }()},
		{"float32", float32(1.5), float32(-1.5)},
		{"float64", float64(2.5), float64(-2.5)},
		{"bool true", true, false},
		{"bool false", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Negate(tt.in)
			if got != tt.want {
				t.Fatalf("Negate(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

// Ensure Negate coerces unknown types via ToInt.
func TestNegateDefault(t *testing.T) {
	// default case coerces via cast.ToInt
	got := Negate("5")
	if got != -5 {
		t.Fatalf("Negate(\"5\") = %v, want -5", got)
	}
}

// Ensure Exponent returns a**b as float64.
func TestExponent(t *testing.T) {
	tests := []struct {
		a, b any
		want float64
	}{
		{2, 3, 8},
		{int64(3), int64(2), 9},
		{float64(2.5), float64(2), 6.25},
		{10, 0, 1},
	}
	for _, tt := range tests {
		got := Exponent(tt.a, tt.b)
		if got != tt.want {
			t.Fatalf("Exponent(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

// Ensure DivideInt returns float64 for exact, int64 for inexact.
func TestIntDivide(t *testing.T) {
	tests := []struct {
		name string
		a, b any
		want any
	}{
		{"exact", int(20), int(4), float64(5)},
		{"truncated", int64(7), int64(2), int64(3)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DivideInt(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf("DivideInt(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// Ensure BitwiseOr returns a | b.
func TestBitwiseOr(t *testing.T) {
	if got := BitwiseOr(int64(0b1010), int64(0b0101)); got != 0b1111 {
		t.Fatalf("BitwiseOr = %v", got)
	}
}

// Ensure BitwiseXor returns a ^ b.
func TestBitwiseXor(t *testing.T) {
	if got := BitwiseXor(int64(0b1111), int64(0b0101)); got != 0b1010 {
		t.Fatalf("BitwiseXor = %v", got)
	}
}

// Ensure BitwiseAnd returns a & b.
func TestBitwiseAnd(t *testing.T) {
	if got := BitwiseAnd(int64(0b1110), int64(0b1011)); got != 0b1010 {
		t.Fatalf("BitwiseAnd = %v", got)
	}
}

// Ensure BitwiseNot returns ^a.
func TestBitwiseNot(t *testing.T) {
	if got := BitwiseNot(int64(0)); got != -1 {
		t.Fatalf("BitwiseNot = %v", got)
	}
}

// Ensure ShiftLeft returns a << b.
func TestShiftLeft(t *testing.T) {
	if got := ShiftLeft(int64(1), int64(3)); got != 8 {
		t.Fatalf("ShiftLeft = %v", got)
	}
}

// Ensure ShiftRight returns a >> b (arithmetic).
func TestShiftRight(t *testing.T) {
	if got := ShiftRight(int64(16), int64(2)); got != 4 {
		t.Fatalf("ShiftRight = %v", got)
	}
}

// Ensure ShiftRightUnsigned returns a >>> b (zero-filling).
func TestShiftRightUnsigned(t *testing.T) {
	if got := ShiftRightUnsigned(int64(-1), int64(1)); got != math.MaxInt64 {
		t.Fatalf("ShiftRightUnsigned = %v", got)
	}
}

// Ensure StrictEqual requires matching type and value.
func TestStrictEqual(t *testing.T) {
	tests := []struct {
		name string
		a, b any
		want bool
	}{
		{"both nil", nil, nil, true},
		{"nil vs value", nil, 1, false},
		{"value vs nil", 1, nil, false},
		{"same int", 42, 42, true},
		{"diff type same val", int(1), int64(1), false},
		{"same string", "foo", "foo", true},
		{"diff string", "foo", "bar", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StrictEqual(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf("StrictEqual(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// Ensure Increment adds one, preserving type including big.Int.
func TestIncrement(t *testing.T) {
	bigI := &big.Int{}
	bigI.SetInt64(5)

	tests := []struct {
		name string
		in   any
		want any
	}{
		{"int", int(4), int(5)},
		{"int8", int8(1), int8(2)},
		{"int16", int16(2), int16(3)},
		{"int32", int32(3), int32(4)},
		{"int64", int64(9), int64(10)},
		{"uint", uint(1), uint(2)},
		{"uint8", uint8(1), uint8(2)},
		{"uint16", uint16(1), uint16(2)},
		{"uint32", uint32(1), uint32(2)},
		{"uint64", uint64(1), uint64(2)},
		{"float32", float32(1.0), float32(2.0)},
		{"float64", float64(1.5), float64(2.5)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Increment(tt.in)
			if err != nil || got != tt.want {
				t.Fatalf("Increment(%v) = %v, %v; want %v", tt.in, got, err, tt.want)
			}
		})
	}
	// *big.Int
	gotBip, err := Increment(big.NewInt(1))
	if err != nil {
		t.Fatalf("Increment(*big.Int) err=%v", err)
	}
	_ = gotBip
	// decimal.Decimal
	gotBdp, err := Increment(decimal.NewFromInt(1))
	if err != nil {
		t.Fatalf("Increment(decimal.Decimal) err=%v", err)
	}
	_ = gotBdp
	// error case
	_, err = Increment("bad")
	if err == nil {
		t.Fatal("expected error for string")
	}
}

// Ensure Decrement subtracts one, preserving type incl. big.Int.
func TestDecrement(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want any
	}{
		{"int", int(5), int(4)},
		{"int8", int8(2), int8(1)},
		{"int16", int16(3), int16(2)},
		{"int32", int32(4), int32(3)},
		{"int64", int64(10), int64(9)},
		{"uint", uint(2), uint(1)},
		{"uint8", uint8(2), uint8(1)},
		{"uint16", uint16(2), uint16(1)},
		{"uint32", uint32(2), uint32(1)},
		{"uint64", uint64(2), uint64(1)},
		{"float32", float32(2.0), float32(1.0)},
		{"float64", float64(2.5), float64(1.5)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decrement(tt.in)
			if err != nil || got != tt.want {
				t.Fatalf("Decrement(%v) = %v, %v; want %v", tt.in, got, err, tt.want)
			}
		})
	}

	// *big.Int
	gotBip, err := Decrement(big.NewInt(2))
	if err != nil {
		t.Fatalf("Decrement(*big.Int) err=%v", err)
	}
	_ = gotBip
	// decimal.Decimal
	gotBdp, err := Decrement(decimal.NewFromInt(2))
	if err != nil {
		t.Fatalf("Decrement(decimal.Decimal) err=%v", err)
	}
	_ = gotBdp
	// error case
	_, err = Decrement("bad")
	if err == nil {
		t.Fatal("expected error for string")
	}
}

// Ensure ToRange produces ascending and descending slices.
func TestToRange(t *testing.T) {
	tests := []struct {
		min, max int
		want     []int
	}{
		{1, 3, []int{1, 2, 3}},
		{3, 1, []int{3, 2, 1}},
		{5, 5, []int{5}},
	}
	for _, tt := range tests {
		got := ToRange(tt.min, tt.max)
		if len(got) != len(tt.want) {
			t.Fatalf("ToRange(%d, %d) len=%d, want %d", tt.min, tt.max, len(got), len(tt.want))
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Fatalf("ToRange(%d, %d)[%d] = %d, want %d", tt.min, tt.max, i, got[i], tt.want[i])
			}
		}
	}
}

// Ensure IsFalsy returns true for nil, false, zero, and "false".
func TestIsFalsy(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want bool
	}{
		{"nil", nil, true},
		{"false", false, true},
		{"true", true, false},
		{"zero int", int(0), true},
		{"nonzero int", int(1), false},
		{"zero int8", int8(0), true},
		{"zero int16", int16(0), true},
		{"zero int32", int32(0), true},
		{"zero int64", int64(0), true},
		{"zero uint", uint(0), true},
		{"nonzero uint", uint(1), false},
		{"zero uint8", uint8(0), true},
		{"zero uint16", uint16(0), true},
		{"zero uint32", uint32(0), true},
		{"zero uint64", uint64(0), true},
		{"zero float32", float32(0), true},
		{"nonzero float32", float32(0.1), false},
		{"zero float64", float64(0), true},
		{"nonzero float64", float64(0.1), false},
		{"empty string", "", true},
		{"false string", "false", true},
		{"FALSE string", "FALSE", true},
		{"nonempty string", "hello", false},
		{"struct (default false)", struct{}{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsFalsy(tt.in)
			if got != tt.want {
				t.Fatalf("IsFalsy(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

// Ensure IsNil returns true for typed nils (ptr, chan, slice).
func TestIsNil(t *testing.T) {
	var p *int
	var ch chan int
	var fn func()
	var sl []int
	tests := []struct {
		name string
		in   any
		want bool
	}{
		{"nil", nil, true},
		{"typed nil ptr", p, true},
		{"typed nil chan", ch, true},
		{"typed nil func", fn, true},
		{"typed nil slice", sl, true},
		{"non-nil int", 42, false},
		{"non-nil string", "x", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsNil(tt.in)
			if got != tt.want {
				t.Fatalf("IsNil(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

// Ensure Empty returns true for nil, zero, NaN, and empty maps.
func TestEmpty(t *testing.T) {
	nilPtr := (*int)(nil)
	tests := []struct {
		name string
		in   any
		want bool
	}{
		{"nil", nil, true},
		{"empty string", "", true},
		{"nonempty string", "x", false},
		{"zero int", int(0), true},
		{"nonzero int", int(1), false},
		{"zero int8", int8(0), true},
		{"zero int16", int16(0), true},
		{"zero int32", int32(0), true},
		{"zero int64", int64(0), true},
		{"zero uint", uint(0), true},
		{"zero uint8", uint8(0), true},
		{"zero uint16", uint16(0), true},
		{"zero uint32", uint32(0), true},
		{"zero uint64", uint64(0), true},
		{"zero float32", float32(0), true},
		{"nonzero float32", float32(0.1), false},
		{"zero float64", float64(0), true},
		{"NaN float", math.NaN(), true},
		{"nonzero float64", float64(0.1), false},
		{"empty slice", []any{}, true},
		{"nonempty slice", []any{1}, false},
		{"empty map", map[string]any{}, true},
		{"nonempty map", map[string]any{"k": "v"}, false},
		{"bool false", false, false},
		{"bool true", true, false},
		{"nil ptr", nilPtr, true},
		{"struct (other)", struct{}{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Empty(tt.in)
			if got != tt.want {
				t.Fatalf("Empty(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

// Ensure Size returns the element count for strings, slices,
// maps, sets, and pointer indirection.
func TestSize(t *testing.T) {
	s := []any{1, 2, 3}
	var nilSlicePtr *[]any
	tests := []struct {
		name string
		in   any
		want int64
	}{
		{"nil", nil, 0},
		{"string", "hello", 5},
		{"slice", []any{1, 2, 3}, 3},
		{"map", map[string]any{"a": 1, "b": 2}, 2},
		{"set", util.NewHashSetFrom([]any{1, 2, 3}), 3},
		{"ptr to slice", &s, 3},
		{"nil ptr", nilSlicePtr, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Size(tt.in)
			if got != tt.want {
				t.Fatalf("Size(%v) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

// Ensure Size panics for unsupported types.
func TestSizePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	Size(42)
}
