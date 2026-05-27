// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import (
	"math"
	"testing"
)

// Ensure Abs returns the absolute value of a negative float.
func TestAbs_negative(t *testing.T) {
	if Abs(-3.5) != 3.5 {
		t.Fatal("expected 3.5")
	}
}

// Ensure Abs returns the value unchanged when already positive.
func TestAbs_positive(t *testing.T) {
	if Abs(2.0) != 2.0 {
		t.Fatal("expected 2.0")
	}
}

// Ensure Ceil rounds up to the nearest integer.
func TestCeil_positive(t *testing.T) {
	if Ceil(1.2) != 2.0 {
		t.Fatal("expected 2.0")
	}
}

// Ensure Floor rounds down to the nearest integer.
func TestFloor_positive(t *testing.T) {
	if Floor(1.9) != 1.0 {
		t.Fatal("expected 1.0")
	}
}

// Ensure IsInfinite returns true for positive infinity.
func TestIsInfinite_posInf(t *testing.T) {
	if !IsInfinite(math.Inf(1)) {
		t.Fatal("expected true")
	}
}

// Ensure IsNaN returns true for NaN.
func TestIsNaN_nan(t *testing.T) {
	if !IsNaN(math.NaN()) {
		t.Fatal("expected true")
	}
}

// Ensure Log returns natural log of e ≈ 1.
func TestLog_e(t *testing.T) {
	if math.Abs(Log(math.E)-1.0) > 1e-9 {
		t.Fatal("expected ~1.0")
	}
}

// Ensure Log2 returns 3 for input 8.
func TestLog2_eight(t *testing.T) {
	if Log2(8) != 3.0 {
		t.Fatal("expected 3.0")
	}
}

// Ensure Log10 returns 2 for input 100.
func TestLog10_hundred(t *testing.T) {
	if Log10(100) != 2.0 {
		t.Fatal("expected 2.0")
	}
}

// Ensure Max returns the larger of two values.
func TestMax_basic(t *testing.T) {
	if Max(3.0, 7.0) != 7.0 {
		t.Fatal("expected 7.0")
	}
}

// Ensure Min returns the smaller of two values.
func TestMin_basic(t *testing.T) {
	if Min(3.0, 7.0) != 3.0 {
		t.Fatal("expected 3.0")
	}
}

// Ensure Pow returns the correct exponentiation.
func TestPow_basic(t *testing.T) {
	if Pow(2.0, 10.0) != 1024.0 {
		t.Fatal("expected 1024.0")
	}
}

// Ensure Round rounds to the nearest integer.
func TestRound_half(t *testing.T) {
	if Round(2.5) != 3.0 {
		t.Fatal("expected 3.0")
	}
}

// Ensure Sqrt returns the correct square root.
func TestSqrt_basic(t *testing.T) {
	if Sqrt(9.0) != 3.0 {
		t.Fatal("expected 3.0")
	}
}

// Ensure BitAnd returns the bitwise AND of two integers.
func TestBitAnd_basic(t *testing.T) {
	if BitAnd(0b1010, 0b1100) != 0b1000 {
		t.Fatal("expected 8")
	}
}

// Ensure BitOr returns the bitwise OR of two integers.
func TestBitOr_basic(t *testing.T) {
	if BitOr(0b1010, 0b1100) != 0b1110 {
		t.Fatal("expected 14")
	}
}

// Ensure BitXor returns the bitwise XOR of two integers.
func TestBitXor_basic(t *testing.T) {
	if BitXor(0b1010, 0b1100) != 0b0110 {
		t.Fatal("expected 6")
	}
}

// Ensure BitNot returns the bitwise NOT of an integer.
func TestBitNot_basic(t *testing.T) {
	if BitNot(0) != -1 {
		t.Fatal("expected -1")
	}
}

// Ensure BitShiftLeft shifts bits left by n positions.
func TestBitShiftLeft_basic(t *testing.T) {
	if BitShiftLeft(1, 3) != 8 {
		t.Fatal("expected 8")
	}
}

// Ensure BitShiftRight shifts bits right by n positions.
func TestBitShiftRight_basic(t *testing.T) {
	if BitShiftRight(8, 3) != 1 {
		t.Fatal("expected 1")
	}
}
