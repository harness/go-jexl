// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lang

import (
	"testing"
)

// Ensure new Boolean() with no args returns false.
func TestBoolean_newNoArgs(t *testing.T) {
	inst, err := BooleanClass.Call("new")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if inst != BooleanFalse {
		t.Errorf("expected false, got %v", inst)
	}
}

// Ensure new Boolean(true) returns true.
func TestBoolean_newBoolTrue(t *testing.T) {
	inst, err := BooleanClass.Call("new", true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if inst != BooleanTrue {
		t.Errorf("expected true, got %v", inst)
	}
}

// Ensure new Boolean("true") parses the string case-insensitively.
func TestBoolean_newStringTrue(t *testing.T) {
	inst, err := BooleanClass.Call("new", "TRUE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if inst != BooleanTrue {
		t.Errorf("expected true, got %v", inst)
	}
}

// Ensure Boolean.parseBoolean returns true for "true" (case-insensitive).
func TestBoolean_parseBoolean(t *testing.T) {
	got, err := BooleanClass.Call("parseBoolean", "True")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != BooleanTrue {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure Boolean.valueOf returns the correct Boolean for a bool arg.
func TestBoolean_valueOf(t *testing.T) {
	got, err := BooleanClass.Call("valueOf", false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != BooleanFalse {
		t.Errorf("expected false, got %v", got)
	}
}

// Ensure Boolean.toString returns "true" for true.
func TestBoolean_toString(t *testing.T) {
	got, err := BooleanClass.Call("toString", true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "true" {
		t.Errorf("expected \"true\", got %v", got)
	}
}

// Ensure Boolean.compare returns 0 for equal values.
func TestBoolean_compareEqual(t *testing.T) {
	got, err := BooleanClass.Call("compare", true, true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 0 {
		t.Errorf("expected 0, got %v", got)
	}
}

// Ensure Boolean.compare returns 1 when first arg is true and second is false.
func TestBoolean_compareGreater(t *testing.T) {
	got, err := BooleanClass.Call("compare", true, false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 1 {
		t.Errorf("expected 1, got %v", got)
	}
}

// Ensure Boolean.compare returns -1 when first arg is false and second is true.
func TestBoolean_compareLess(t *testing.T) {
	got, err := BooleanClass.Call("compare", false, true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != -1 {
		t.Errorf("expected -1, got %v", got)
	}
}

// Ensure Boolean.logicalAnd returns the logical AND of two bools.
func TestBoolean_logicalAnd(t *testing.T) {
	got, err := BooleanClass.Call("logicalAnd", true, false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != false {
		t.Errorf("expected false, got %v", got)
	}
}

// Ensure Boolean.logicalOr returns the logical OR of two bools.
func TestBoolean_logicalOr(t *testing.T) {
	got, err := BooleanClass.Call("logicalOr", true, false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure Boolean.logicalXor returns true when inputs differ.
func TestBoolean_logicalXor(t *testing.T) {
	got, err := BooleanClass.Call("logicalXor", true, false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure Boolean.logicalXor returns false when inputs are equal.
func TestBoolean_logicalXorEqual(t *testing.T) {
	got, err := BooleanClass.Call("logicalXor", true, true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != false {
		t.Errorf("expected false, got %v", got)
	}
}

// Ensure booleanValue on an instance returns the underlying bool.
func TestBoolean_instanceBooleanValue(t *testing.T) {
	inst, err := BooleanClass.Call("new", true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	b, ok := inst.(Boolean)
	if !ok {
		t.Error("expected Boolean instance")
		return
	}
	got, err := b.Call("booleanValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure toString on an instance returns the string representation.
func TestBoolean_instanceToString(t *testing.T) {
	inst, err := BooleanClass.Call("new", false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	b, ok := inst.(Boolean)
	if !ok {
		t.Error("expected Boolean instance")
		return
	}
	got, err := b.Call("toString")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "false" {
		t.Errorf("expected \"false\", got %v", got)
	}
}

// Ensure compareTo returns 0 when values are equal.
func TestBoolean_instanceCompareToEqual(t *testing.T) {
	b := NewBoolean(true)
	got, err := b.Call("compareTo", true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 0 {
		t.Errorf("expected 0, got %v", got)
	}
}

// Ensure compareTo returns 1 when receiver is true and arg is false.
func TestBoolean_instanceCompareToGreater(t *testing.T) {
	b := NewBoolean(true)
	got, err := b.Call("compareTo", false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 1 {
		t.Errorf("expected 1, got %v", got)
	}
}

// Ensure equals returns true when values match.
func TestBoolean_instanceEqualsTrue(t *testing.T) {
	b := NewBoolean(false)
	got, err := b.Call("equals", false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure equals returns false when values differ.
func TestBoolean_instanceEqualsFalse(t *testing.T) {
	b := NewBoolean(true)
	got, err := b.Call("equals", false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != false {
		t.Errorf("expected false, got %v", got)
	}
}

// Ensure an unknown class method returns an error.
func TestBoolean_unknownClassMethod(t *testing.T) {
	if _, err := BooleanClass.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown class method")
	}
}

// Ensure an unknown instance method returns an error.
func TestBoolean_unknownInstanceMethod(t *testing.T) {
	b := NewBoolean(true)
	if _, err := b.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown instance method")
	}
}

// Ensure parseBoolean with wrong arg count returns an error.
func TestBoolean_parseBooleanArgCount(t *testing.T) {
	if _, err := BooleanClass.Call("parseBoolean"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure toString with wrong arg count returns an error.
func TestBoolean_classToStringArgCount(t *testing.T) {
	if _, err := BooleanClass.Call("toString"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure compare with wrong arg count returns an error.
func TestBoolean_classCompareArgCount(t *testing.T) {
	if _, err := BooleanClass.Call("compare", true); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure logicalAnd with wrong arg count returns an error.
func TestBoolean_classLogicalAndArgCount(t *testing.T) {
	if _, err := BooleanClass.Call("logicalAnd", true); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure logicalOr with wrong arg count returns an error.
func TestBoolean_classLogicalOrArgCount(t *testing.T) {
	if _, err := BooleanClass.Call("logicalOr", true); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure logicalXor with wrong arg count returns an error.
func TestBoolean_classLogicalXorArgCount(t *testing.T) {
	if _, err := BooleanClass.Call("logicalXor", true); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure instance compareTo with wrong arg count returns an error.
func TestBoolean_instanceCompareToArgCount(t *testing.T) {
	b := NewBoolean(true)
	if _, err := b.Call("compareTo"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure instance equals with wrong arg count returns an error.
func TestBoolean_instanceEqualsArgCount(t *testing.T) {
	b := NewBoolean(true)
	if _, err := b.Call("equals"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure toInteger on true returns 1.
func TestBoolean_instanceToInteger(t *testing.T) {
	b := NewBoolean(true)
	got, err := b.Call("toInteger")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int32(1) {
		t.Errorf("expected 1, got %v", got)
	}
}

// Ensure toBoolean aliases booleanValue.
func TestBoolean_instanceToBoolean(t *testing.T) {
	b := NewBoolean(true)
	got, err := b.Call("toBoolean")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure default on a non-null Boolean returns the boolean itself, not the fallback.
func TestBoolean_instanceDefault(t *testing.T) {
	b := NewBoolean(true)
	got, err := b.Call("default", false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != b {
		t.Errorf("expected %v, got %v", b, got)
	}
}
