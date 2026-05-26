// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lang

import (
	"testing"
)

// Ensure new Byte() with no args returns zero value.
func TestByte_newNoArgs(t *testing.T) {
	inst, err := ByteClass.Call("new")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if inst != NewByte(0) {
		t.Errorf("expected Byte(0), got %v", inst)
	}
}

// Ensure new Byte(127) returns Byte(127).
func TestByte_newWithArg(t *testing.T) {
	inst, err := ByteClass.Call("new", int64(127))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if inst != NewByte(127) {
		t.Errorf("expected Byte(127), got %v", inst)
	}
}

// Ensure Byte.valueOf returns Byte(-1) for int64(-1).
func TestByte_valueOf(t *testing.T) {
	got, err := ByteClass.Call("valueOf", int64(-1))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewByte(-1) {
		t.Errorf("expected Byte(-1), got %v", got)
	}
}

// Ensure Byte.toString returns the string representation.
func TestByte_classToString(t *testing.T) {
	got, err := ByteClass.Call("toString", int64(42))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "42" {
		t.Errorf("expected \"42\", got %v", got)
	}
}

// Ensure Byte.compare returns -1 when first arg is less than second.
func TestByte_classCompare(t *testing.T) {
	got, err := ByteClass.Call("compare", int64(10), int64(20))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != -1 {
		t.Errorf("expected -1, got %v", got)
	}
}

// Ensure Byte.MAX_VALUE returns Byte(127).
func TestByte_maxValue(t *testing.T) {
	got, err := ByteClass.Call("MAX_VALUE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != Byte(127) {
		t.Errorf("expected Byte(127), got %v", got)
	}
}

// Ensure Byte.MIN_VALUE returns Byte(-128).
func TestByte_minValue(t *testing.T) {
	got, err := ByteClass.Call("MIN_VALUE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != Byte(-128) {
		t.Errorf("expected Byte(-128), got %v", got)
	}
}

// Ensure Byte.SIZE returns int64(8).
func TestByte_size(t *testing.T) {
	got, err := ByteClass.Call("SIZE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(8) {
		t.Errorf("expected int64(8), got %v", got)
	}
}

// Ensure Byte.BYTES returns int64(1).
func TestByte_bytes(t *testing.T) {
	got, err := ByteClass.Call("BYTES")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(1) {
		t.Errorf("expected int64(1), got %v", got)
	}
}

// Ensure byteValue on Byte(42) returns int8(42).
func TestByte_instanceByteValue(t *testing.T) {
	b := NewByte(42)
	got, err := b.Call("byteValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int8(42) {
		t.Errorf("expected int8(42), got %v", got)
	}
}

// Ensure shortValue on Byte(10) returns int16(10).
func TestByte_instanceShortValue(t *testing.T) {
	b := NewByte(10)
	got, err := b.Call("shortValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int16(10) {
		t.Errorf("expected int16(10), got %v", got)
	}
}

// Ensure intValue on Byte(10) returns int32(10).
func TestByte_instanceIntValue(t *testing.T) {
	b := NewByte(10)
	got, err := b.Call("intValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int32(10) {
		t.Errorf("expected int32(10), got %v", got)
	}
}

// Ensure longValue on Byte(10) returns int64(10).
func TestByte_instanceLongValue(t *testing.T) {
	b := NewByte(10)
	got, err := b.Call("longValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(10) {
		t.Errorf("expected int64(10), got %v", got)
	}
}

// Ensure doubleValue on Byte(10) returns float64(10).
func TestByte_instanceDoubleValue(t *testing.T) {
	b := NewByte(10)
	got, err := b.Call("doubleValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != float64(10) {
		t.Errorf("expected float64(10), got %v", got)
	}
}

// Ensure toString on Byte(42) returns "42".
func TestByte_instanceToString(t *testing.T) {
	b := NewByte(42)
	got, err := b.Call("toString")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "42" {
		t.Errorf("expected \"42\", got %v", got)
	}
}

// Ensure compareTo returns -1 when receiver is less than arg.
func TestByte_instanceCompareTo(t *testing.T) {
	b := NewByte(10)
	got, err := b.Call("compareTo", int64(20))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != -1 {
		t.Errorf("expected -1, got %v", got)
	}
}

// Ensure compareTo returns 1 when receiver is greater than arg.
func TestByte_instanceCompareToGreater(t *testing.T) {
	b := NewByte(20)
	got, err := b.Call("compareTo", int64(10))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 1 {
		t.Errorf("expected 1, got %v", got)
	}
}

// Ensure equals returns true when values match.
func TestByte_instanceEquals(t *testing.T) {
	b := NewByte(5)
	got, err := b.Call("equals", int64(5))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure an unknown class method returns an error.
func TestByte_unknownClassMethod(t *testing.T) {
	if _, err := ByteClass.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown class method")
	}
}

// Ensure an unknown instance method returns an error.
func TestByte_unknownInstanceMethod(t *testing.T) {
	b := NewByte(1)
	if _, err := b.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown instance method")
	}
}

// Ensure valueOf with wrong arg count returns an error.
func TestByte_valueOfArgCount(t *testing.T) {
	if _, err := ByteClass.Call("valueOf"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure toString with wrong arg count returns an error.
func TestByte_classToStringArgCount(t *testing.T) {
	if _, err := ByteClass.Call("toString"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure compare with wrong arg count returns an error.
func TestByte_classCompareArgCount(t *testing.T) {
	if _, err := ByteClass.Call("compare", int64(1)); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure instance compareTo with wrong arg count returns an error.
func TestByte_instanceCompareToArgCount(t *testing.T) {
	b := NewByte(1)
	if _, err := b.Call("compareTo"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure instance equals with wrong arg count returns an error.
func TestByte_instanceEqualsArgCount(t *testing.T) {
	b := NewByte(1)
	if _, err := b.Call("equals"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure toInteger aliases intValue on Byte.
func TestByte_instanceToInteger(t *testing.T) {
	b := NewByte(5)
	got, err := b.Call("toInteger")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int32(5) {
		t.Errorf("expected 5, got %v", got)
	}
}

// Ensure default on a non-null Byte returns the byte itself, not the fallback.
func TestByte_instanceDefault(t *testing.T) {
	b := NewByte(1)
	got, err := b.Call("default", int8(0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != b {
		t.Errorf("expected %v, got %v", b, got)
	}
}
