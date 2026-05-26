// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lang

import (
	"testing"
)

// Ensure new Short() with no args returns zero value.
func TestShort_newNoArgs(t *testing.T) {
	inst, err := ShortClass.Call("new")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if inst != NewShort(0) {
		t.Errorf("expected Short(0), got %v", inst)
	}
}

// Ensure new Short(1000) returns Short(1000).
func TestShort_newWithArg(t *testing.T) {
	inst, err := ShortClass.Call("new", int64(1000))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if inst != NewShort(1000) {
		t.Errorf("expected Short(1000), got %v", inst)
	}
}

// Ensure Short.parseShort returns Short(32767) for int64(32767).
func TestShort_parseShort(t *testing.T) {
	got, err := ShortClass.Call("parseShort", int64(32767))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewShort(32767) {
		t.Errorf("expected Short(32767), got %v", got)
	}
}

// Ensure Short.valueOf returns Short(5).
func TestShort_valueOf(t *testing.T) {
	got, err := ShortClass.Call("valueOf", int64(5))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewShort(5) {
		t.Errorf("expected Short(5), got %v", got)
	}
}

// Ensure Short.toString returns the string representation.
func TestShort_classToString(t *testing.T) {
	got, err := ShortClass.Call("toString", int64(100))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "100" {
		t.Errorf("expected \"100\", got %v", got)
	}
}

// Ensure Short.compare returns 1 when first arg is greater.
func TestShort_classCompare(t *testing.T) {
	got, err := ShortClass.Call("compare", int64(20), int64(10))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 1 {
		t.Errorf("expected 1, got %v", got)
	}
}

// Ensure Short.MAX_VALUE returns Short(32767).
func TestShort_maxValue(t *testing.T) {
	got, err := ShortClass.Call("MAX_VALUE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != Short(32767) {
		t.Errorf("expected Short(32767), got %v", got)
	}
}

// Ensure Short.MIN_VALUE returns Short(-32768).
func TestShort_minValue(t *testing.T) {
	got, err := ShortClass.Call("MIN_VALUE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != Short(-32768) {
		t.Errorf("expected Short(-32768), got %v", got)
	}
}

// Ensure Short.SIZE returns int64(16).
func TestShort_size(t *testing.T) {
	got, err := ShortClass.Call("SIZE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(16) {
		t.Errorf("expected int64(16), got %v", got)
	}
}

// Ensure Short.BYTES returns int64(2).
func TestShort_bytes(t *testing.T) {
	got, err := ShortClass.Call("BYTES")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(2) {
		t.Errorf("expected int64(2), got %v", got)
	}
}

// Ensure shortValue on Short(100) returns int16(100).
func TestShort_instanceShortValue(t *testing.T) {
	s := NewShort(100)
	got, err := s.Call("shortValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int16(100) {
		t.Errorf("expected int16(100), got %v", got)
	}
}

// Ensure byteValue on Short(10) returns int8(10).
func TestShort_instanceByteValue(t *testing.T) {
	s := NewShort(10)
	got, err := s.Call("byteValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int8(10) {
		t.Errorf("expected int8(10), got %v", got)
	}
}

// Ensure intValue on Short(10) returns int32(10).
func TestShort_instanceIntValue(t *testing.T) {
	s := NewShort(10)
	got, err := s.Call("intValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int32(10) {
		t.Errorf("expected int32(10), got %v", got)
	}
}

// Ensure longValue on Short(10) returns int64(10).
func TestShort_instanceLongValue(t *testing.T) {
	s := NewShort(10)
	got, err := s.Call("longValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int64(10) {
		t.Errorf("expected int64(10), got %v", got)
	}
}

// Ensure doubleValue on Short(10) returns float64(10).
func TestShort_instanceDoubleValue(t *testing.T) {
	s := NewShort(10)
	got, err := s.Call("doubleValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != float64(10) {
		t.Errorf("expected float64(10), got %v", got)
	}
}

// Ensure toString on Short(42) returns "42".
func TestShort_instanceToString(t *testing.T) {
	s := NewShort(42)
	got, err := s.Call("toString")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "42" {
		t.Errorf("expected \"42\", got %v", got)
	}
}

// Ensure compareTo returns -1 when receiver is less than arg.
func TestShort_instanceCompareTo(t *testing.T) {
	s := NewShort(5)
	got, err := s.Call("compareTo", int64(10))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != -1 {
		t.Errorf("expected -1, got %v", got)
	}
}

// Ensure compareTo returns 1 when receiver is greater than arg.
func TestShort_instanceCompareToGreater(t *testing.T) {
	s := NewShort(10)
	got, err := s.Call("compareTo", int64(5))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 1 {
		t.Errorf("expected 1, got %v", got)
	}
}

// Ensure equals returns true when values match.
func TestShort_instanceEquals(t *testing.T) {
	s := NewShort(7)
	got, err := s.Call("equals", int64(7))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure an unknown class method returns an error.
func TestShort_unknownClassMethod(t *testing.T) {
	if _, err := ShortClass.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown class method")
	}
}

// Ensure an unknown instance method returns an error.
func TestShort_unknownInstanceMethod(t *testing.T) {
	s := NewShort(1)
	if _, err := s.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown instance method")
	}
}

// Ensure valueOf with wrong arg count returns an error.
func TestShort_valueOfArgCount(t *testing.T) {
	if _, err := ShortClass.Call("valueOf"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure toString with wrong arg count returns an error.
func TestShort_classToStringArgCount(t *testing.T) {
	if _, err := ShortClass.Call("toString"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure compare with wrong arg count returns an error.
func TestShort_classCompareArgCount(t *testing.T) {
	if _, err := ShortClass.Call("compare", int64(1)); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure instance compareTo with wrong arg count returns an error.
func TestShort_instanceCompareToArgCount(t *testing.T) {
	s := NewShort(1)
	if _, err := s.Call("compareTo"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure instance equals with wrong arg count returns an error.
func TestShort_instanceEqualsArgCount(t *testing.T) {
	s := NewShort(1)
	if _, err := s.Call("equals"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure toInteger aliases intValue on Short.
func TestShort_instanceToInteger(t *testing.T) {
	s := NewShort(8)
	got, err := s.Call("toInteger")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int32(8) {
		t.Errorf("expected 8, got %v", got)
	}
}

// Ensure default on a non-null Short returns the short itself, not the fallback.
func TestShort_instanceDefault(t *testing.T) {
	s := NewShort(1)
	got, err := s.Call("default", int16(0))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != s {
		t.Errorf("expected %v, got %v", s, got)
	}
}
