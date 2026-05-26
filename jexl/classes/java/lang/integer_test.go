// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lang

import (
	"math"
	"testing"
)

// Ensure new Integer() with no args returns zero value.
func TestInteger_newNoArgs(t *testing.T) {
	inst, err := IntegerClass.Call("new")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if inst != NewInteger(0) {
		t.Errorf("expected Integer(0), got %v", inst)
	}
}

// Ensure new Integer(42) returns Integer(42).
func TestInteger_newWithArg(t *testing.T) {
	inst, err := IntegerClass.Call("new", int64(42))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if inst != NewInteger(42) {
		t.Errorf("expected Integer(42), got %v", inst)
	}
}

// Ensure Integer.valueOf returns Integer(7).
func TestInteger_valueOf(t *testing.T) {
	got, err := IntegerClass.Call("valueOf", int64(7))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewInteger(7) {
		t.Errorf("expected Integer(7), got %v", got)
	}
}

// Ensure Integer.parseInt returns Integer(100) for "100".
func TestInteger_parseInt(t *testing.T) {
	got, err := IntegerClass.Call("parseInt", "100")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewInteger(100) {
		t.Errorf("expected Integer(100), got %v", got)
	}
}

// Ensure Integer.parseInt with base 2 parses binary string.
func TestInteger_parseIntWithBase(t *testing.T) {
	got, err := IntegerClass.Call("parseInt", "1010", int64(2))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewInteger(10) {
		t.Errorf("expected Integer(10), got %v", got)
	}
}

// Ensure Integer.parseInt with wrong arg count returns an error.
func TestInteger_parseIntArgCount(t *testing.T) {
	if _, err := IntegerClass.Call("parseInt"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure Integer.toString with one arg returns string form.
func TestInteger_classToStringOneArg(t *testing.T) {
	got, err := IntegerClass.Call("toString", int64(255))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "255" {
		t.Errorf("expected \"255\", got %v", got)
	}
}

// Ensure Integer.toString with base returns hex string.
func TestInteger_classToStringTwoArgs(t *testing.T) {
	got, err := IntegerClass.Call("toString", int64(255), int64(16))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "ff" {
		t.Errorf("expected \"ff\", got %v", got)
	}
}

// Ensure Integer.toString with wrong arg count returns an error.
func TestInteger_classToStringArgCount(t *testing.T) {
	if _, err := IntegerClass.Call("toString"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure Integer.compare returns -1 when first arg is less.
func TestInteger_classCompare(t *testing.T) {
	got, err := IntegerClass.Call("compare", int64(1), int64(2))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != -1 {
		t.Errorf("expected -1, got %v", got)
	}
}

// Ensure Integer.compare with wrong arg count returns an error.
func TestInteger_classCompareArgCount(t *testing.T) {
	if _, err := IntegerClass.Call("compare", int64(1)); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure Integer.max returns Integer(7) for args 3 and 7.
func TestInteger_max(t *testing.T) {
	got, err := IntegerClass.Call("max", int64(3), int64(7))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewInteger(7) {
		t.Errorf("expected Integer(7), got %v", got)
	}
}

// Ensure Integer.max returns Integer(3) when first arg is larger.
func TestInteger_maxFirstLarger(t *testing.T) {
	got, err := IntegerClass.Call("max", int64(3), int64(2))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewInteger(3) {
		t.Errorf("expected Integer(3), got %v", got)
	}
}

// Ensure Integer.max with wrong arg count returns an error.
func TestInteger_classMaxArgCount(t *testing.T) {
	if _, err := IntegerClass.Call("max", int64(1)); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure Integer.min returns Integer(3) for args 3 and 7.
func TestInteger_min(t *testing.T) {
	got, err := IntegerClass.Call("min", int64(3), int64(7))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewInteger(3) {
		t.Errorf("expected Integer(3), got %v", got)
	}
}

// Ensure Integer.min returns Integer(2) when second arg is smaller.
func TestInteger_minSecondSmaller(t *testing.T) {
	got, err := IntegerClass.Call("min", int64(3), int64(2))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewInteger(2) {
		t.Errorf("expected Integer(2), got %v", got)
	}
}

// Ensure Integer.min with wrong arg count returns an error.
func TestInteger_classMinArgCount(t *testing.T) {
	if _, err := IntegerClass.Call("min", int64(1)); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure Integer.sum returns Integer(5) for args 2 and 3.
func TestInteger_sum(t *testing.T) {
	got, err := IntegerClass.Call("sum", int64(2), int64(3))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewInteger(5) {
		t.Errorf("expected Integer(5), got %v", got)
	}
}

// Ensure Integer.sum with wrong arg count returns an error.
func TestInteger_classSumArgCount(t *testing.T) {
	if _, err := IntegerClass.Call("sum", int64(1)); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure Integer.toBinaryString returns "1010" for 10.
func TestInteger_toBinaryString(t *testing.T) {
	got, err := IntegerClass.Call("toBinaryString", int64(10))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "1010" {
		t.Errorf("expected \"1010\", got %v", got)
	}
}

// Ensure Integer.toBinaryString with wrong arg count returns an error.
func TestInteger_toBinaryStringArgCount(t *testing.T) {
	if _, err := IntegerClass.Call("toBinaryString"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure Integer.toHexString returns "ff" for 255.
func TestInteger_toHexString(t *testing.T) {
	got, err := IntegerClass.Call("toHexString", int64(255))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "ff" {
		t.Errorf("expected \"ff\", got %v", got)
	}
}

// Ensure Integer.toHexString with wrong arg count returns an error.
func TestInteger_toHexStringArgCount(t *testing.T) {
	if _, err := IntegerClass.Call("toHexString"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure Integer.toOctalString returns "17" for 15.
func TestInteger_toOctalString(t *testing.T) {
	got, err := IntegerClass.Call("toOctalString", int64(15))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "17" {
		t.Errorf("expected \"17\", got %v", got)
	}
}

// Ensure Integer.toOctalString with wrong arg count returns an error.
func TestInteger_toOctalStringArgCount(t *testing.T) {
	if _, err := IntegerClass.Call("toOctalString"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure Integer.MAX_VALUE returns Integer(math.MaxInt32).
func TestInteger_maxValue(t *testing.T) {
	got, err := IntegerClass.Call("MAX_VALUE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != Integer(math.MaxInt32) {
		t.Errorf("expected MaxInt32, got %v", got)
	}
}

// Ensure Integer.MIN_VALUE returns Integer(math.MinInt32).
func TestInteger_minValue(t *testing.T) {
	got, err := IntegerClass.Call("MIN_VALUE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != Integer(math.MinInt32) {
		t.Errorf("expected MinInt32, got %v", got)
	}
}

// Ensure Integer.SIZE returns int64(32).
func TestInteger_size(t *testing.T) {
	got, err := IntegerClass.Call("SIZE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(32) {
		t.Errorf("expected int64(32), got %v", got)
	}
}

// Ensure Integer.BYTES returns int64(4).
func TestInteger_bytes(t *testing.T) {
	got, err := IntegerClass.Call("BYTES")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(4) {
		t.Errorf("expected int64(4), got %v", got)
	}
}

// Ensure intValue on Integer(7) returns int32(7).
func TestInteger_instanceIntValue(t *testing.T) {
	i := NewInteger(7)
	got, err := i.Call("intValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int32(7) {
		t.Errorf("expected int32(7), got %v", got)
	}
}

// Ensure longValue on Integer(7) returns int64(7).
func TestInteger_instanceLongValue(t *testing.T) {
	i := NewInteger(7)
	got, err := i.Call("longValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(7) {
		t.Errorf("expected int64(7), got %v", got)
	}
}

// Ensure shortValue on Integer(7) returns int16(7).
func TestInteger_instanceShortValue(t *testing.T) {
	i := NewInteger(7)
	got, err := i.Call("shortValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int16(7) {
		t.Errorf("expected int16(7), got %v", got)
	}
}

// Ensure byteValue on Integer(7) returns int8(7).
func TestInteger_instanceByteValue(t *testing.T) {
	i := NewInteger(7)
	got, err := i.Call("byteValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int8(7) {
		t.Errorf("expected int8(7), got %v", got)
	}
}

// Ensure doubleValue on Integer(7) returns float64(7).
func TestInteger_instanceDoubleValue(t *testing.T) {
	i := NewInteger(7)
	got, err := i.Call("doubleValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != float64(7) {
		t.Errorf("expected float64(7), got %v", got)
	}
}

// Ensure toString on Integer(255) returns "255".
func TestInteger_instanceToString(t *testing.T) {
	i := NewInteger(255)
	got, err := i.Call("toString")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "255" {
		t.Errorf("expected \"255\", got %v", got)
	}
}

// Ensure compareTo returns -1 when receiver is less than arg.
func TestInteger_instanceCompareTo(t *testing.T) {
	i := NewInteger(1)
	got, err := i.Call("compareTo", int64(2))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != -1 {
		t.Errorf("expected -1, got %v", got)
	}
}

// Ensure compareTo returns 1 when receiver is greater than arg.
func TestInteger_instanceCompareToGreater(t *testing.T) {
	i := NewInteger(2)
	got, err := i.Call("compareTo", int64(1))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 1 {
		t.Errorf("expected 1, got %v", got)
	}
}

// Ensure equals returns true when values match.
func TestInteger_instanceEquals(t *testing.T) {
	i := NewInteger(5)
	got, err := i.Call("equals", int64(5))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure an unknown class method returns an error.
func TestInteger_unknownClassMethod(t *testing.T) {
	if _, err := IntegerClass.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown class method")
	}
}

// Ensure an unknown instance method returns an error.
func TestInteger_unknownInstanceMethod(t *testing.T) {
	i := NewInteger(1)
	if _, err := i.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown instance method")
	}
}

// Ensure valueOf with wrong arg count returns an error.
func TestInteger_valueOfArgCount(t *testing.T) {
	if _, err := IntegerClass.Call("valueOf"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure instance compareTo with wrong arg count returns an error.
func TestInteger_instanceCompareToArgCount(t *testing.T) {
	i := NewInteger(1)
	if _, err := i.Call("compareTo"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure instance equals with wrong arg count returns an error.
func TestInteger_instanceEqualsArgCount(t *testing.T) {
	i := NewInteger(1)
	if _, err := i.Call("equals"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure toInteger aliases intValue on Integer.
func TestInteger_instanceToInteger(t *testing.T) {
	i := NewInteger(7)
	got, err := i.Call("toInteger")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int32(7) {
		t.Errorf("expected 7, got %v", got)
	}
}

// Ensure default on a non-null Integer returns the integer itself, not the fallback.
func TestInteger_instanceDefault(t *testing.T) {
	i := NewInteger(42)
	got, err := i.Call("default", int32(0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != i {
		t.Errorf("expected %v, got %v", i, got)
	}
}
