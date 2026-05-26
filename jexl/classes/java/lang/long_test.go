// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lang

import (
	"math"
	"testing"
)

// Ensure new Long() with no args returns zero value.
func TestLong_newNoArgs(t *testing.T) {
	inst, err := LongClass.Call("new")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if inst != NewLong(0) {
		t.Errorf("expected Long(0), got %v", inst)
	}
}

// Ensure new Long(9999) returns Long(9999).
func TestLong_newWithArg(t *testing.T) {
	inst, err := LongClass.Call("new", int64(9999))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if inst != NewLong(9999) {
		t.Errorf("expected Long(9999), got %v", inst)
	}
}

// Ensure Long.valueOf returns Long(7).
func TestLong_valueOf(t *testing.T) {
	got, err := LongClass.Call("valueOf", int64(7))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewLong(7) {
		t.Errorf("expected Long(7), got %v", got)
	}
}

// Ensure Long.valueOf with wrong arg count returns an error.
func TestLong_valueOfArgCount(t *testing.T) {
	if _, err := LongClass.Call("valueOf"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure Long.parseLong returns Long(42) for "42".
func TestLong_parseLong(t *testing.T) {
	got, err := LongClass.Call("parseLong", "42")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewLong(42) {
		t.Errorf("expected Long(42), got %v", got)
	}
}

// Ensure Long.parseLong with base 16 parses hex string.
func TestLong_parseLongWithBase(t *testing.T) {
	got, err := LongClass.Call("parseLong", "ff", int64(16))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewLong(255) {
		t.Errorf("expected Long(255), got %v", got)
	}
}

// Ensure Long.parseLong with wrong arg count returns an error.
func TestLong_parseLongArgCount(t *testing.T) {
	if _, err := LongClass.Call("parseLong"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure Long.toString with one arg returns string form.
func TestLong_classToStringOneArg(t *testing.T) {
	got, err := LongClass.Call("toString", int64(42))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "42" {
		t.Errorf("expected \"42\", got %v", got)
	}
}

// Ensure Long.toString with base returns binary string.
func TestLong_classToStringTwoArgs(t *testing.T) {
	got, err := LongClass.Call("toString", int64(10), int64(2))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "1010" {
		t.Errorf("expected \"1010\", got %v", got)
	}
}

// Ensure Long.toString with wrong arg count returns an error.
func TestLong_classToStringArgCount(t *testing.T) {
	if _, err := LongClass.Call("toString"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure Long.compare returns 0 for equal values.
func TestLong_classCompareEqual(t *testing.T) {
	got, err := LongClass.Call("compare", int64(5), int64(5))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 0 {
		t.Errorf("expected 0, got %v", got)
	}
}

// Ensure Long.compare with wrong arg count returns an error.
func TestLong_classCompareArgCount(t *testing.T) {
	if _, err := LongClass.Call("compare", int64(1)); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure Long.max returns Long(7) for args 3 and 7.
func TestLong_max(t *testing.T) {
	got, err := LongClass.Call("max", int64(3), int64(7))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewLong(7) {
		t.Errorf("expected Long(7), got %v", got)
	}
}

// Ensure Long.max returns Long(3) when first arg is larger.
func TestLong_maxFirstLarger(t *testing.T) {
	got, err := LongClass.Call("max", int64(3), int64(2))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewLong(3) {
		t.Errorf("expected Long(3), got %v", got)
	}
}

// Ensure Long.max with wrong arg count returns an error.
func TestLong_classMaxArgCount(t *testing.T) {
	if _, err := LongClass.Call("max", int64(1)); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure Long.min returns Long(3) for args 3 and 7.
func TestLong_min(t *testing.T) {
	got, err := LongClass.Call("min", int64(3), int64(7))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewLong(3) {
		t.Errorf("expected Long(3), got %v", got)
	}
}

// Ensure Long.min returns Long(2) when second arg is smaller.
func TestLong_minSecondSmaller(t *testing.T) {
	got, err := LongClass.Call("min", int64(3), int64(2))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewLong(2) {
		t.Errorf("expected Long(2), got %v", got)
	}
}

// Ensure Long.min with wrong arg count returns an error.
func TestLong_classMinArgCount(t *testing.T) {
	if _, err := LongClass.Call("min", int64(1)); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure Long.toBinaryString returns "1010" for 10.
func TestLong_toBinaryString(t *testing.T) {
	got, err := LongClass.Call("toBinaryString", int64(10))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "1010" {
		t.Errorf("expected \"1010\", got %v", got)
	}
}

// Ensure Long.toBinaryString with wrong arg count returns an error.
func TestLong_toBinaryStringArgCount(t *testing.T) {
	if _, err := LongClass.Call("toBinaryString"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure Long.toHexString returns "ff" for 255.
func TestLong_toHexString(t *testing.T) {
	got, err := LongClass.Call("toHexString", int64(255))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "ff" {
		t.Errorf("expected \"ff\", got %v", got)
	}
}

// Ensure Long.toHexString with wrong arg count returns an error.
func TestLong_toHexStringArgCount(t *testing.T) {
	if _, err := LongClass.Call("toHexString"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure Long.MAX_VALUE returns Long(math.MaxInt64).
func TestLong_maxValue(t *testing.T) {
	got, err := LongClass.Call("MAX_VALUE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != Long(math.MaxInt64) {
		t.Errorf("expected MaxInt64, got %v", got)
	}
}

// Ensure Long.MIN_VALUE returns Long(math.MinInt64).
func TestLong_minValue(t *testing.T) {
	got, err := LongClass.Call("MIN_VALUE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != Long(math.MinInt64) {
		t.Errorf("expected MinInt64, got %v", got)
	}
}

// Ensure Long.SIZE returns int64(64).
func TestLong_size(t *testing.T) {
	got, err := LongClass.Call("SIZE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(64) {
		t.Errorf("expected int64(64), got %v", got)
	}
}

// Ensure Long.BYTES returns int64(8).
func TestLong_bytes(t *testing.T) {
	got, err := LongClass.Call("BYTES")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(8) {
		t.Errorf("expected int64(8), got %v", got)
	}
}

// Ensure longValue on Long(7) returns int64(7).
func TestLong_instanceLongValue(t *testing.T) {
	l := NewLong(7)
	got, err := l.Call("longValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(7) {
		t.Errorf("expected int64(7), got %v", got)
	}
}

// Ensure intValue on Long(7) returns int32(7).
func TestLong_instanceIntValue(t *testing.T) {
	l := NewLong(7)
	got, err := l.Call("intValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int32(7) {
		t.Errorf("expected int32(7), got %v", got)
	}
}

// Ensure shortValue on Long(7) returns int16(7).
func TestLong_instanceShortValue(t *testing.T) {
	l := NewLong(7)
	got, err := l.Call("shortValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int16(7) {
		t.Errorf("expected int16(7), got %v", got)
	}
}

// Ensure byteValue on Long(7) returns int8(7).
func TestLong_instanceByteValue(t *testing.T) {
	l := NewLong(7)
	got, err := l.Call("byteValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int8(7) {
		t.Errorf("expected int8(7), got %v", got)
	}
}

// Ensure doubleValue on Long(7) returns float64(7).
func TestLong_instanceDoubleValue(t *testing.T) {
	l := NewLong(7)
	got, err := l.Call("doubleValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != float64(7) {
		t.Errorf("expected float64(7), got %v", got)
	}
}

// Ensure toString on Long(42) returns "42".
func TestLong_instanceToString(t *testing.T) {
	l := NewLong(42)
	got, err := l.Call("toString")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "42" {
		t.Errorf("expected \"42\", got %v", got)
	}
}

// Ensure compareTo returns 0 when receiver equals arg.
func TestLong_instanceCompareToEqual(t *testing.T) {
	l := NewLong(5)
	got, err := l.Call("compareTo", int64(5))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 0 {
		t.Errorf("expected 0, got %v", got)
	}
}

// Ensure compareTo returns 1 when receiver is greater than arg.
func TestLong_instanceCompareToGreater(t *testing.T) {
	l := NewLong(10)
	got, err := l.Call("compareTo", int64(5))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 1 {
		t.Errorf("expected 1, got %v", got)
	}
}

// Ensure equals returns true when values match.
func TestLong_instanceEquals(t *testing.T) {
	l := NewLong(42)
	got, err := l.Call("equals", int64(42))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure an unknown class method returns an error.
func TestLong_unknownClassMethod(t *testing.T) {
	if _, err := LongClass.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown class method")
	}
}

// Ensure an unknown instance method returns an error.
func TestLong_unknownInstanceMethod(t *testing.T) {
	l := NewLong(1)
	if _, err := l.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown instance method")
	}
}

// Ensure instance compareTo with wrong arg count returns an error.
func TestLong_instanceCompareToArgCount(t *testing.T) {
	l := NewLong(1)
	if _, err := l.Call("compareTo"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure instance equals with wrong arg count returns an error.
func TestLong_instanceEqualsArgCount(t *testing.T) {
	l := NewLong(1)
	if _, err := l.Call("equals"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure toInteger aliases intValue on Long.
func TestLong_instanceToInteger(t *testing.T) {
	l := NewLong(100)
	got, err := l.Call("toInteger")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int32(100) {
		t.Errorf("expected 100, got %v", got)
	}
}

// Ensure default on a non-null Long returns the long itself, not the fallback.
func TestLong_instanceDefault(t *testing.T) {
	l := NewLong(42)
	got, err := l.Call("default", int64(0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != l {
		t.Errorf("expected %v, got %v", l, got)
	}
}
