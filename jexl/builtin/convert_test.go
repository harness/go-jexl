// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import "testing"

// Ensure BooleanValue converts a truthy string.
func TestBooleanValue_trueString(t *testing.T) {
	if !BooleanValue("true") {
		t.Fatal("expected true")
	}
}

// Ensure BooleanValue converts a falsy value.
func TestBooleanValue_zero(t *testing.T) {
	if BooleanValue(0) {
		t.Fatal("expected false")
	}
}

// Ensure ByteValue returns v as int8.
func TestByteValue_basic(t *testing.T) {
	if ByteValue(127) != 127 {
		t.Fatal("expected 127")
	}
}

// Ensure DoubleValue returns v as float64.
func TestDoubleValue_basic(t *testing.T) {
	if DoubleValue("3.14") != 3.14 {
		t.Fatal("expected 3.14")
	}
}

// Ensure FloatValue returns v as float32.
func TestFloatValue_basic(t *testing.T) {
	if FloatValue("1.5") != float32(1.5) {
		t.Fatal("expected 1.5")
	}
}

// Ensure IntValue returns v as int.
func TestIntValue_basic(t *testing.T) {
	if IntValue("42") != 42 {
		t.Fatal("expected 42")
	}
}

// Ensure LongValue returns v as int64.
func TestLongValue_basic(t *testing.T) {
	if LongValue("9999999999") != 9999999999 {
		t.Fatal("expected 9999999999")
	}
}

// Ensure ShortValue returns v as int16.
func TestShortValue_basic(t *testing.T) {
	if ShortValue(1000) != 1000 {
		t.Fatal("expected 1000")
	}
}

// Ensure ToString converts an integer to string.
func TestToString_int(t *testing.T) {
	if ToString(42) != "42" {
		t.Fatal("expected 42")
	}
}

// Ensure DefaultValue returns the value when it is non-nil and non-zero.
func TestDefaultValue_nonNil(t *testing.T) {
	if DefaultValue("hello", "fallback") != "hello" {
		t.Fatal("expected hello")
	}
}

// Ensure DefaultValue returns fallback when value is nil.
func TestDefaultValue_nil(t *testing.T) {
	if DefaultValue(nil, "fallback") != "fallback" {
		t.Fatal("expected fallback")
	}
}

// Ensure DefaultValue returns fallback when value is the zero string.
func TestDefaultValue_zeroString(t *testing.T) {
	if DefaultValue("", "fallback") != "fallback" {
		t.Fatal("expected fallback")
	}
}
