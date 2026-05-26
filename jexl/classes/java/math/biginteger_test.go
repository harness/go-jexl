// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package math

import (
	"math/big"
	"testing"
)

// Ensure BigInteger.new() with no args returns a zero BigInteger.
func TestBigInteger_newNoArgs(t *testing.T) {
	got, err := BigIntegerClass.Call("new")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	bi, ok := got.(*BigInteger)
	if !ok {
		t.Errorf("expected *BigInteger, got %T", got)
		return
	}
	if bi.V.Cmp(new(big.Int)) != 0 {
		t.Errorf("expected 0, got %v", bi.V)
	}
}

// Ensure BigInteger.new("12345") constructs from a decimal string.
func TestBigInteger_newFromString(t *testing.T) {
	got, err := BigIntegerClass.Call("new", "12345")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	bi, ok := got.(*BigInteger)
	if !ok {
		t.Errorf("expected *BigInteger, got %T", got)
		return
	}
	if bi.V.Int64() != 12345 {
		t.Errorf("expected 12345, got %v", bi.V)
	}
}

// Ensure BigInteger.new(int64(42)) constructs from an integer.
func TestBigInteger_newFromInt(t *testing.T) {
	got, err := BigIntegerClass.Call("new", int64(42))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	bi, ok := got.(*BigInteger)
	if !ok {
		t.Errorf("expected *BigInteger, got %T", got)
		return
	}
	if bi.V.Int64() != 42 {
		t.Errorf("expected 42, got %v", bi.V)
	}
}

// Ensure add(5) on a BigInteger(10) returns intValue 15.
func TestBigInteger_add(t *testing.T) {
	inst, err := BigIntegerClass.Call("new", int64(10))
	if err != nil {
		t.Errorf("unexpected error constructing: %v", err)
		return
	}
	bi := inst.(*BigInteger)
	result, err := bi.Call("add", int64(5))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	sum, ok := result.(*BigInteger)
	if !ok {
		t.Errorf("expected *BigInteger, got %T", result)
		return
	}
	if sum.IntValue() != int64(15) {
		t.Errorf("expected 15, got %v", sum.IntValue())
	}
}

// Ensure multiply(7) on a BigInteger(6) returns intValue 42.
func TestBigInteger_multiply(t *testing.T) {
	inst, err := BigIntegerClass.Call("new", int64(6))
	if err != nil {
		t.Errorf("unexpected error constructing: %v", err)
		return
	}
	bi := inst.(*BigInteger)
	result, err := bi.Call("multiply", int64(7))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	product, ok := result.(*BigInteger)
	if !ok {
		t.Errorf("expected *BigInteger, got %T", result)
		return
	}
	if product.IntValue() != int64(42) {
		t.Errorf("expected 42, got %v", product.IntValue())
	}
}

// Ensure negate() on a BigInteger(5) returns intValue -5.
func TestBigInteger_negate(t *testing.T) {
	inst, err := BigIntegerClass.Call("new", int64(5))
	if err != nil {
		t.Errorf("unexpected error constructing: %v", err)
		return
	}
	bi := inst.(*BigInteger)
	result, err := bi.Call("negate")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	neg, ok := result.(*BigInteger)
	if !ok {
		t.Errorf("expected *BigInteger, got %T", result)
		return
	}
	if neg.IntValue() != int64(-5) {
		t.Errorf("expected -5, got %v", neg.IntValue())
	}
}

// Ensure compareTo(5) on BigInteger(10) returns 1.
func TestBigInteger_compareTo(t *testing.T) {
	inst, err := BigIntegerClass.Call("new", int64(10))
	if err != nil {
		t.Errorf("unexpected error constructing: %v", err)
		return
	}
	bi := inst.(*BigInteger)
	result, err := bi.Call("compareTo", int64(5))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if result != 1 {
		t.Errorf("expected 1, got %v", result)
	}
}

// Ensure an unknown BigInteger class method returns an error.
func TestBigInteger_unknownClassMethod(t *testing.T) {
	if _, err := BigIntegerClass.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown class method")
	}
}

// Ensure an unknown BigInteger instance method returns an error.
func TestBigInteger_unknownInstanceMethod(t *testing.T) {
	bi := &BigInteger{V: big.NewInt(0)}
	if _, err := bi.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown instance method")
	}
}
