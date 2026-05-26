// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package coerce

import (
	"math"
	"math/big"
	"testing"

	"github.com/harness/go-jexl/jexl/internal/decimal"
)

// Ensure ToBigInt recognises a *big.Int value.
func TestToBigInt_recognized(t *testing.T) {
	v := big.NewInt(42)
	got, ok := ToBigInt(v)
	if !ok || got != v {
		t.Fatal("expected *big.Int to be recognized")
	}
}

// Ensure ToBigInt returns false for non-big.Int values.
func TestToBigInt_notBigInt(t *testing.T) {
	_, ok := ToBigInt(42)
	if ok {
		t.Fatal("expected false for plain int")
	}
}

// Ensure ToBigDecimal recognises a decimal.Decimal value.
func TestToBigDecimal_recognized(t *testing.T) {
	v := decimal.NewFromFloat(3.14)
	got, ok := ToBigDecimal(v)
	if !ok || !got.Equal(v) {
		t.Fatal("expected decimal.Decimal to be recognized")
	}
}

// Ensure ToBigDecimal returns false for non-decimal values.
func TestToBigDecimal_notDecimal(t *testing.T) {
	_, ok := ToBigDecimal(3.14)
	if ok {
		t.Fatal("expected false for float64")
	}
}

// Ensure ToDecimal passes a decimal.Decimal through unchanged.
func TestToDecimal_decimal(t *testing.T) {
	v := decimal.NewFromFloat(1.5)
	got := ToDecimal(v)
	if !got.Equal(v) {
		t.Fatalf("ToDecimal(decimal) = %v, want %v", got, v)
	}
}

// Ensure ToDecimal converts *big.Int to the equivalent decimal.
func TestToDecimal_bigInt(t *testing.T) {
	v := big.NewInt(100)
	got := ToDecimal(v)
	want := decimal.NewFromInt(100)
	if !got.Equal(want) {
		t.Fatalf("ToDecimal(*big.Int) = %v, want %v", got, want)
	}
}

// Ensure ToDecimal converts float64 to the equivalent decimal.
func TestToDecimal_float64(t *testing.T) {
	got := ToDecimal(float64(2.5))
	want := decimal.NewFromFloat(2.5)
	if !got.Equal(want) {
		t.Fatalf("ToDecimal(float64) = %v, want %v", got, want)
	}
}

// Ensure ToDecimal returns zero for NaN.
func TestToDecimal_nan(t *testing.T) {
	got := ToDecimal(math.NaN())
	if !got.Equal(decimal.Zero) {
		t.Fatalf("ToDecimal(NaN) = %v, want zero", got)
	}
}

// Ensure ToDecimal returns zero for Inf.
func TestToDecimal_inf(t *testing.T) {
	got := ToDecimal(math.Inf(1))
	if !got.Equal(decimal.Zero) {
		t.Fatalf("ToDecimal(Inf) = %v, want zero", got)
	}
}

// Ensure ToNumeric converts bool true to 1.
func TestToNumeric_boolTrue(t *testing.T) {
	got := ToNumeric(true)
	if got != 1 {
		t.Fatalf("ToNumeric(true) = %v, want 1", got)
	}
}

// Ensure ToNumeric converts bool false to 0.
func TestToNumeric_boolFalse(t *testing.T) {
	got := ToNumeric(false)
	if got != 0 {
		t.Fatalf("ToNumeric(false) = %v, want 0", got)
	}
}

// Ensure ToNumeric converts an integer string to int.
func TestToNumeric_stringInt(t *testing.T) {
	got := ToNumeric("42")
	if got != int(42) {
		t.Fatalf("ToNumeric(\"42\") = %v (%T), want int(42)", got, got)
	}
}

// Ensure ToNumeric converts a float string to float64.
func TestToNumeric_stringFloat(t *testing.T) {
	got := ToNumeric("3.14")
	if got != float64(3.14) {
		t.Fatalf("ToNumeric(\"3.14\") = %v (%T), want float64(3.14)", got, got)
	}
}

// Ensure ToNumeric returns non-numeric strings unchanged.
func TestToNumeric_stringInvalid(t *testing.T) {
	got := ToNumeric("abc")
	if got != "abc" {
		t.Fatalf("ToNumeric(\"abc\") = %v, want \"abc\"", got)
	}
}

// Ensure ToNumeric returns already-numeric values unchanged.
func TestToNumeric_numeric(t *testing.T) {
	got := ToNumeric(7)
	if got != 7 {
		t.Fatalf("ToNumeric(7) = %v, want 7", got)
	}
}
