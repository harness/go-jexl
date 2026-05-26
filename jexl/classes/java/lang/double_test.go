// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lang

import (
	"math"
	"testing"
)

// Ensure new Double() with no args returns zero value.
func TestDouble_newNoArgs(t *testing.T) {
	inst, err := DoubleClass.Call("new")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if inst != NewDouble(0) {
		t.Errorf("expected Double(0), got %v", inst)
	}
}

// Ensure new Double(3.14) returns Double(3.14).
func TestDouble_newWithArg(t *testing.T) {
	inst, err := DoubleClass.Call("new", float64(3.14))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if inst != NewDouble(3.14) {
		t.Errorf("expected Double(3.14), got %v", inst)
	}
}

// Ensure Double.parseDouble returns Double(2.5) for "2.5".
func TestDouble_parseDouble(t *testing.T) {
	got, err := DoubleClass.Call("parseDouble", "2.5")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewDouble(2.5) {
		t.Errorf("expected Double(2.5), got %v", got)
	}
}

// Ensure Double.valueOf returns Double(1.0).
func TestDouble_valueOf(t *testing.T) {
	got, err := DoubleClass.Call("valueOf", float64(1.0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewDouble(1.0) {
		t.Errorf("expected Double(1.0), got %v", got)
	}
}

// Ensure Double.toString returns the string form.
func TestDouble_classToString(t *testing.T) {
	got, err := DoubleClass.Call("toString", float64(1.5))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "1.5" {
		t.Errorf("expected \"1.5\", got %v", got)
	}
}

// Ensure Double.compare returns -1 when first arg is less.
func TestDouble_classCompare(t *testing.T) {
	got, err := DoubleClass.Call("compare", float64(1.0), float64(2.0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != -1 {
		t.Errorf("expected -1, got %v", got)
	}
}

// Ensure Double.max returns Double(2.0) for args 1.0 and 2.0.
func TestDouble_max(t *testing.T) {
	got, err := DoubleClass.Call("max", float64(1.0), float64(2.0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewDouble(2.0) {
		t.Errorf("expected Double(2.0), got %v", got)
	}
}

// Ensure Double.min returns Double(1.0) for args 1.0 and 2.0.
func TestDouble_min(t *testing.T) {
	got, err := DoubleClass.Call("min", float64(1.0), float64(2.0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewDouble(1.0) {
		t.Errorf("expected Double(1.0), got %v", got)
	}
}

// Ensure Double.sum returns Double(3.0) for args 1.0 and 2.0.
func TestDouble_sum(t *testing.T) {
	got, err := DoubleClass.Call("sum", float64(1.0), float64(2.0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewDouble(3.0) {
		t.Errorf("expected Double(3.0), got %v", got)
	}
}

// Ensure Double.isNaN returns true for NaN.
func TestDouble_classIsNaN(t *testing.T) {
	got, err := DoubleClass.Call("isNaN", math.NaN())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure Double.isInfinite returns true for +Inf.
func TestDouble_classIsInfinite(t *testing.T) {
	got, err := DoubleClass.Call("isInfinite", math.Inf(1))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure Double.isFinite returns true for a normal value.
func TestDouble_classIsFinite(t *testing.T) {
	got, err := DoubleClass.Call("isFinite", float64(1.0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure Double.MAX_VALUE returns a large positive Double.
func TestDouble_maxValue(t *testing.T) {
	got, err := DoubleClass.Call("MAX_VALUE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != Double(math.MaxFloat64) {
		t.Errorf("expected MaxFloat64, got %v", got)
	}
}

// Ensure Double.MIN_VALUE returns the smallest nonzero Double.
func TestDouble_minValue(t *testing.T) {
	got, err := DoubleClass.Call("MIN_VALUE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != Double(math.SmallestNonzeroFloat64) {
		t.Errorf("expected SmallestNonzeroFloat64, got %v", got)
	}
}

// Ensure Double.POSITIVE_INFINITY returns +Inf.
func TestDouble_positiveInfinity(t *testing.T) {
	got, err := DoubleClass.Call("POSITIVE_INFINITY")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	d, ok := got.(Double)
	if !ok {
		t.Errorf("expected Double, got %T", got)
		return
	}
	if !math.IsInf(float64(d), 1) {
		t.Errorf("expected +Inf, got %v", got)
	}
}

// Ensure Double.NEGATIVE_INFINITY returns -Inf.
func TestDouble_negativeInfinity(t *testing.T) {
	got, err := DoubleClass.Call("NEGATIVE_INFINITY")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	d, ok := got.(Double)
	if !ok {
		t.Errorf("expected Double, got %T", got)
		return
	}
	if !math.IsInf(float64(d), -1) {
		t.Errorf("expected -Inf, got %v", got)
	}
}

// Ensure Double.NaN returns a NaN Double.
func TestDouble_nan(t *testing.T) {
	got, err := DoubleClass.Call("NaN")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	d, ok := got.(Double)
	if !ok {
		t.Errorf("expected Double, got %T", got)
		return
	}
	if !math.IsNaN(float64(d)) {
		t.Errorf("expected NaN, got %v", got)
	}
}

// Ensure Double.MAX_EXPONENT returns int64(1023).
func TestDouble_maxExponent(t *testing.T) {
	got, err := DoubleClass.Call("MAX_EXPONENT")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(1023) {
		t.Errorf("expected int64(1023), got %v", got)
	}
}

// Ensure Double.MIN_EXPONENT returns int64(-1022).
func TestDouble_minExponent(t *testing.T) {
	got, err := DoubleClass.Call("MIN_EXPONENT")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(-1022) {
		t.Errorf("expected int64(-1022), got %v", got)
	}
}

// Ensure Double.SIZE returns int64(64).
func TestDouble_size(t *testing.T) {
	got, err := DoubleClass.Call("SIZE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(64) {
		t.Errorf("expected int64(64), got %v", got)
	}
}

// Ensure Double.BYTES returns int64(8).
func TestDouble_bytes(t *testing.T) {
	got, err := DoubleClass.Call("BYTES")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(8) {
		t.Errorf("expected int64(8), got %v", got)
	}
}

// Ensure doubleValue on Double(3.0) returns float64(3.0).
func TestDouble_instanceDoubleValue(t *testing.T) {
	d := NewDouble(3.0)
	got, err := d.Call("doubleValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != float64(3.0) {
		t.Errorf("expected float64(3.0), got %v", got)
	}
}

// Ensure intValue on Double(3.9) returns int64(3).
func TestDouble_instanceIntValue(t *testing.T) {
	d := NewDouble(3.9)
	got, err := d.Call("intValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(3) {
		t.Errorf("expected int64(3), got %v", got)
	}
}

// Ensure longValue on Double(5.0) returns int64(5).
func TestDouble_instanceLongValue(t *testing.T) {
	d := NewDouble(5.0)
	got, err := d.Call("longValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(5) {
		t.Errorf("expected int64(5), got %v", got)
	}
}

// Ensure floatValue on Double(2.0) returns float32(2.0).
func TestDouble_instanceFloatValue(t *testing.T) {
	d := NewDouble(2.0)
	got, err := d.Call("floatValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != float32(2.0) {
		t.Errorf("expected float32(2.0), got %v", got)
	}
}

// Ensure isNaN on Double(0) returns false.
func TestDouble_instanceIsNaN(t *testing.T) {
	d := NewDouble(0)
	got, err := d.Call("isNaN")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != false {
		t.Errorf("expected false, got %v", got)
	}
}

// Ensure isInfinite on Double(0) returns false.
func TestDouble_instanceIsInfinite(t *testing.T) {
	d := NewDouble(0)
	got, err := d.Call("isInfinite")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != false {
		t.Errorf("expected false, got %v", got)
	}
}

// Ensure isFinite on Double(1.0) returns true.
func TestDouble_instanceIsFinite(t *testing.T) {
	d := NewDouble(1.0)
	got, err := d.Call("isFinite")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure toString on Double(1.5) returns "1.5".
func TestDouble_instanceToString(t *testing.T) {
	d := NewDouble(1.5)
	got, err := d.Call("toString")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "1.5" {
		t.Errorf("expected \"1.5\", got %v", got)
	}
}

// Ensure compareTo returns -1 when receiver is less than arg.
func TestDouble_instanceCompareTo(t *testing.T) {
	d := NewDouble(1.0)
	got, err := d.Call("compareTo", float64(2.0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != -1 {
		t.Errorf("expected -1, got %v", got)
	}
}

// Ensure compareTo returns 1 when receiver is greater than arg.
func TestDouble_instanceCompareToGreater(t *testing.T) {
	d := NewDouble(2.0)
	got, err := d.Call("compareTo", float64(1.0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 1 {
		t.Errorf("expected 1, got %v", got)
	}
}

// Ensure equals returns true when values match.
func TestDouble_instanceEquals(t *testing.T) {
	d := NewDouble(1.5)
	got, err := d.Call("equals", float64(1.5))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure an unknown class method returns an error.
func TestDouble_unknownClassMethod(t *testing.T) {
	if _, err := DoubleClass.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown class method")
	}
}

// Ensure an unknown instance method returns an error.
func TestDouble_unknownInstanceMethod(t *testing.T) {
	d := NewDouble(1.0)
	if _, err := d.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown instance method")
	}
}

// Ensure valueOf with wrong arg count returns an error.
func TestDouble_valueOfArgCount(t *testing.T) {
	if _, err := DoubleClass.Call("valueOf"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure compare with wrong arg count returns an error.
func TestDouble_classCompareArgCount(t *testing.T) {
	if _, err := DoubleClass.Call("compare", float64(1)); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure max with wrong arg count returns an error.
func TestDouble_classMaxArgCount(t *testing.T) {
	if _, err := DoubleClass.Call("max", float64(1)); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure min with wrong arg count returns an error.
func TestDouble_classMinArgCount(t *testing.T) {
	if _, err := DoubleClass.Call("min", float64(1)); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure sum with wrong arg count returns an error.
func TestDouble_classSumArgCount(t *testing.T) {
	if _, err := DoubleClass.Call("sum", float64(1)); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure isNaN with wrong arg count returns an error.
func TestDouble_classIsNaNArgCount(t *testing.T) {
	if _, err := DoubleClass.Call("isNaN"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure isInfinite with wrong arg count returns an error.
func TestDouble_classIsInfiniteArgCount(t *testing.T) {
	if _, err := DoubleClass.Call("isInfinite"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure isFinite with wrong arg count returns an error.
func TestDouble_classIsFiniteArgCount(t *testing.T) {
	if _, err := DoubleClass.Call("isFinite"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure instance compareTo with wrong arg count returns an error.
func TestDouble_instanceCompareToArgCount(t *testing.T) {
	d := NewDouble(1.0)
	if _, err := d.Call("compareTo"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure instance equals with wrong arg count returns an error.
func TestDouble_instanceEqualsArgCount(t *testing.T) {
	d := NewDouble(1.0)
	if _, err := d.Call("equals"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure toInteger aliases intValue on Double.
func TestDouble_instanceToInteger(t *testing.T) {
	d := NewDouble(9.9)
	got, err := d.Call("toInteger")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(9) {
		t.Errorf("expected 9, got %v", got)
	}
}

// Ensure default on a non-null Double returns the double itself, not the fallback.
func TestDouble_instanceDefault(t *testing.T) {
	d := NewDouble(1.5)
	got, err := d.Call("default", float64(0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != d {
		t.Errorf("expected %v, got %v", d, got)
	}
}
