// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lang

import (
	"math"
	"testing"
)

// Ensure new Float() with no args returns zero value.
func TestFloat_newNoArgs(t *testing.T) {
	inst, err := FloatClass.Call("new")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if inst != NewFloat(0) {
		t.Errorf("expected Float(0), got %v", inst)
	}
}

// Ensure new Float(3.14) returns Float(3.14).
func TestFloat_newWithArg(t *testing.T) {
	inst, err := FloatClass.Call("new", float32(3.14))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if inst != NewFloat(3.14) {
		t.Errorf("expected Float(3.14), got %v", inst)
	}
}

// Ensure Float.parseFloat returns Float(2.5) for "2.5".
func TestFloat_parseFloat(t *testing.T) {
	got, err := FloatClass.Call("parseFloat", "2.5")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewFloat(2.5) {
		t.Errorf("expected Float(2.5), got %v", got)
	}
}

// Ensure Float.valueOf returns Float(1.0).
func TestFloat_valueOf(t *testing.T) {
	got, err := FloatClass.Call("valueOf", float32(1.0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewFloat(1.0) {
		t.Errorf("expected Float(1.0), got %v", got)
	}
}

// Ensure Float.toString returns the string form.
func TestFloat_classToString(t *testing.T) {
	got, err := FloatClass.Call("toString", float32(1.5))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "1.5" {
		t.Errorf("expected \"1.5\", got %v", got)
	}
}

// Ensure Float.compare returns -1 when first arg is less.
func TestFloat_classCompare(t *testing.T) {
	got, err := FloatClass.Call("compare", float32(1.0), float32(2.0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != -1 {
		t.Errorf("expected -1, got %v", got)
	}
}

// Ensure Float.max returns Float(2.0) for args 1.0 and 2.0.
func TestFloat_max(t *testing.T) {
	got, err := FloatClass.Call("max", float32(1.0), float32(2.0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewFloat(2.0) {
		t.Errorf("expected Float(2.0), got %v", got)
	}
}

// Ensure Float.min returns Float(1.0) for args 1.0 and 2.0.
func TestFloat_min(t *testing.T) {
	got, err := FloatClass.Call("min", float32(1.0), float32(2.0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewFloat(1.0) {
		t.Errorf("expected Float(1.0), got %v", got)
	}
}

// Ensure Float.sum returns Float(3.0) for args 1.0 and 2.0.
func TestFloat_sum(t *testing.T) {
	got, err := FloatClass.Call("sum", float32(1.0), float32(2.0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewFloat(3.0) {
		t.Errorf("expected Float(3.0), got %v", got)
	}
}

// Ensure Float.isNaN returns true for NaN.
func TestFloat_classIsNaN(t *testing.T) {
	got, err := FloatClass.Call("isNaN", float32(math.NaN()))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure Float.isInfinite returns true for +Inf.
func TestFloat_classIsInfinite(t *testing.T) {
	got, err := FloatClass.Call("isInfinite", float32(math.Inf(1)))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure Float.isFinite returns true for a normal value.
func TestFloat_classIsFinite(t *testing.T) {
	got, err := FloatClass.Call("isFinite", float32(1.0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure Float.MAX_VALUE returns Float(math.MaxFloat32).
func TestFloat_maxValue(t *testing.T) {
	got, err := FloatClass.Call("MAX_VALUE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != Float(math.MaxFloat32) {
		t.Errorf("expected MaxFloat32, got %v", got)
	}
}

// Ensure Float.MIN_VALUE returns the smallest nonzero Float.
func TestFloat_minValue(t *testing.T) {
	got, err := FloatClass.Call("MIN_VALUE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != Float(math.SmallestNonzeroFloat32) {
		t.Errorf("expected SmallestNonzeroFloat32, got %v", got)
	}
}

// Ensure Float.POSITIVE_INFINITY returns +Inf.
func TestFloat_positiveInfinity(t *testing.T) {
	got, err := FloatClass.Call("POSITIVE_INFINITY")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	f, ok := got.(Float)
	if !ok {
		t.Errorf("expected Float, got %T", got)
		return
	}
	if !math.IsInf(float64(f), 1) {
		t.Errorf("expected +Inf, got %v", got)
	}
}

// Ensure Float.NEGATIVE_INFINITY returns -Inf.
func TestFloat_negativeInfinity(t *testing.T) {
	got, err := FloatClass.Call("NEGATIVE_INFINITY")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	f, ok := got.(Float)
	if !ok {
		t.Errorf("expected Float, got %T", got)
		return
	}
	if !math.IsInf(float64(f), -1) {
		t.Errorf("expected -Inf, got %v", got)
	}
}

// Ensure Float.NaN returns a NaN Float.
func TestFloat_nan(t *testing.T) {
	got, err := FloatClass.Call("NaN")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	f, ok := got.(Float)
	if !ok {
		t.Errorf("expected Float, got %T", got)
		return
	}
	if !math.IsNaN(float64(f)) {
		t.Errorf("expected NaN, got %v", got)
	}
}

// Ensure Float.MAX_EXPONENT returns int64(127).
func TestFloat_maxExponent(t *testing.T) {
	got, err := FloatClass.Call("MAX_EXPONENT")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(127) {
		t.Errorf("expected int64(127), got %v", got)
	}
}

// Ensure Float.MIN_EXPONENT returns int64(-126).
func TestFloat_minExponent(t *testing.T) {
	got, err := FloatClass.Call("MIN_EXPONENT")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(-126) {
		t.Errorf("expected int64(-126), got %v", got)
	}
}

// Ensure Float.SIZE returns int64(32).
func TestFloat_size(t *testing.T) {
	got, err := FloatClass.Call("SIZE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(32) {
		t.Errorf("expected int64(32), got %v", got)
	}
}

// Ensure Float.BYTES returns int64(4).
func TestFloat_bytes(t *testing.T) {
	got, err := FloatClass.Call("BYTES")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(4) {
		t.Errorf("expected int64(4), got %v", got)
	}
}

// Ensure floatValue on Float(2.0) returns float32(2.0).
func TestFloat_instanceFloatValue(t *testing.T) {
	f := NewFloat(2.0)
	got, err := f.Call("floatValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != float32(2.0) {
		t.Errorf("expected float32(2.0), got %v", got)
	}
}

// Ensure doubleValue on Float(2.0) returns float64(2.0).
func TestFloat_instanceDoubleValue(t *testing.T) {
	f := NewFloat(2.0)
	got, err := f.Call("doubleValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != float64(2.0) {
		t.Errorf("expected float64(2.0), got %v", got)
	}
}

// Ensure intValue on Float(3.9) returns int32(3).
func TestFloat_instanceIntValue(t *testing.T) {
	f := NewFloat(3.9)
	got, err := f.Call("intValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int32(3) {
		t.Errorf("expected int32(3), got %v", got)
	}
}

// Ensure longValue on Float(5.0) returns int64(5).
func TestFloat_instanceLongValue(t *testing.T) {
	f := NewFloat(5.0)
	got, err := f.Call("longValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(5) {
		t.Errorf("expected int64(5), got %v", got)
	}
}

// Ensure isNaN on Float(0) returns false.
func TestFloat_instanceIsNaN(t *testing.T) {
	f := NewFloat(0)
	got, err := f.Call("isNaN")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != false {
		t.Errorf("expected false, got %v", got)
	}
}

// Ensure isInfinite on Float(0) returns false.
func TestFloat_instanceIsInfinite(t *testing.T) {
	f := NewFloat(0)
	got, err := f.Call("isInfinite")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != false {
		t.Errorf("expected false, got %v", got)
	}
}

// Ensure isFinite on Float(1.0) returns true.
func TestFloat_instanceIsFinite(t *testing.T) {
	f := NewFloat(1.0)
	got, err := f.Call("isFinite")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure toString on Float(1.5) returns "1.5".
func TestFloat_instanceToString(t *testing.T) {
	f := NewFloat(1.5)
	got, err := f.Call("toString")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "1.5" {
		t.Errorf("expected \"1.5\", got %v", got)
	}
}

// Ensure compareTo returns -1 when receiver is less than arg.
func TestFloat_instanceCompareTo(t *testing.T) {
	f := NewFloat(1.0)
	got, err := f.Call("compareTo", float32(2.0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != -1 {
		t.Errorf("expected -1, got %v", got)
	}
}

// Ensure compareTo returns 1 when receiver is greater than arg.
func TestFloat_instanceCompareToGreater(t *testing.T) {
	f := NewFloat(2.0)
	got, err := f.Call("compareTo", float32(1.0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 1 {
		t.Errorf("expected 1, got %v", got)
	}
}

// Ensure equals returns true when values match.
func TestFloat_instanceEquals(t *testing.T) {
	f := NewFloat(1.5)
	got, err := f.Call("equals", float32(1.5))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure an unknown class method returns an error.
func TestFloat_unknownClassMethod(t *testing.T) {
	if _, err := FloatClass.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown class method")
	}
}

// Ensure an unknown instance method returns an error.
func TestFloat_unknownInstanceMethod(t *testing.T) {
	f := NewFloat(1.0)
	if _, err := f.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown instance method")
	}
}

// Ensure valueOf with wrong arg count returns an error.
func TestFloat_valueOfArgCount(t *testing.T) {
	if _, err := FloatClass.Call("valueOf"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure compare with wrong arg count returns an error.
func TestFloat_classCompareArgCount(t *testing.T) {
	if _, err := FloatClass.Call("compare", float32(1)); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure max with wrong arg count returns an error.
func TestFloat_classMaxArgCount(t *testing.T) {
	if _, err := FloatClass.Call("max", float32(1)); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure min with wrong arg count returns an error.
func TestFloat_classMinArgCount(t *testing.T) {
	if _, err := FloatClass.Call("min", float32(1)); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure sum with wrong arg count returns an error.
func TestFloat_classSumArgCount(t *testing.T) {
	if _, err := FloatClass.Call("sum", float32(1)); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure isNaN with wrong arg count returns an error.
func TestFloat_classIsNaNArgCount(t *testing.T) {
	if _, err := FloatClass.Call("isNaN"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure isInfinite with wrong arg count returns an error.
func TestFloat_classIsInfiniteArgCount(t *testing.T) {
	if _, err := FloatClass.Call("isInfinite"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure isFinite with wrong arg count returns an error.
func TestFloat_classIsFiniteArgCount(t *testing.T) {
	if _, err := FloatClass.Call("isFinite"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure instance compareTo with wrong arg count returns an error.
func TestFloat_instanceCompareToArgCount(t *testing.T) {
	f := NewFloat(1.0)
	if _, err := f.Call("compareTo"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure instance equals with wrong arg count returns an error.
func TestFloat_instanceEqualsArgCount(t *testing.T) {
	f := NewFloat(1.0)
	if _, err := f.Call("equals"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure toInteger aliases intValue on Float.
func TestFloat_instanceToInteger(t *testing.T) {
	f := NewFloat(3.9)
	got, err := f.Call("toInteger")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int32(3) {
		t.Errorf("expected 3, got %v", got)
	}
}

// Ensure default on a non-null Float returns the float itself, not the fallback.
func TestFloat_instanceDefault(t *testing.T) {
	f := NewFloat(1.5)
	got, err := f.Call("default", float32(0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != f {
		t.Errorf("expected %v, got %v", f, got)
	}
}
