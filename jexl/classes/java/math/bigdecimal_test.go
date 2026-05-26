// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package math

import (
	"testing"

	"github.com/harness/go-jexl/jexl/internal/decimal"
)

// Ensure BigDecimal.new() with no args returns a zero BigDecimal.
func TestBigDecimal_newNoArgs(t *testing.T) {
	got, err := BigDecimalClass.Call("new")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	bd, ok := got.(*BigDecimal)
	if !ok {
		t.Errorf("expected *BigDecimal, got %T", got)
		return
	}
	if !bd.V.Equal(decimal.Zero) {
		t.Errorf("expected 0, got %v", bd.V)
	}
}

// Ensure BigDecimal.new("3.14") constructs from a decimal string.
func TestBigDecimal_newFromString(t *testing.T) {
	got, err := BigDecimalClass.Call("new", "3.14")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	bd, ok := got.(*BigDecimal)
	if !ok {
		t.Errorf("expected *BigDecimal, got %T", got)
		return
	}
	expected, _ := decimal.NewFromString("3.14")
	if !bd.V.Equal(expected) {
		t.Errorf("expected 3.14, got %v", bd.V)
	}
}

// Ensure add("2.5") on a BigDecimal("1.5") returns doubleValue 4.0.
func TestBigDecimal_add(t *testing.T) {
	inst, err := BigDecimalClass.Call("new", "1.5")
	if err != nil {
		t.Errorf("unexpected error constructing: %v", err)
		return
	}
	bd := inst.(*BigDecimal)

	other, err := BigDecimalClass.Call("new", "2.5")
	if err != nil {
		t.Errorf("unexpected error constructing other: %v", err)
		return
	}

	result, err := bd.Call("add", other)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	sum, ok := result.(*BigDecimal)
	if !ok {
		t.Errorf("expected *BigDecimal, got %T", result)
		return
	}
	if sum.DoubleValue() != 4.0 {
		t.Errorf("expected 4.0, got %v", sum.DoubleValue())
	}
}

// Ensure multiply("3.0") on a BigDecimal("2.0") returns doubleValue 6.0.
func TestBigDecimal_multiply(t *testing.T) {
	inst, err := BigDecimalClass.Call("new", "2.0")
	if err != nil {
		t.Errorf("unexpected error constructing: %v", err)
		return
	}
	bd := inst.(*BigDecimal)

	other, err := BigDecimalClass.Call("new", "3.0")
	if err != nil {
		t.Errorf("unexpected error constructing other: %v", err)
		return
	}

	result, err := bd.Call("multiply", other)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	product, ok := result.(*BigDecimal)
	if !ok {
		t.Errorf("expected *BigDecimal, got %T", result)
		return
	}
	if product.DoubleValue() != 6.0 {
		t.Errorf("expected 6.0, got %v", product.DoubleValue())
	}
}

// Ensure compareTo(BigDecimal("5")) on BigDecimal("10") returns 1.
func TestBigDecimal_compareTo(t *testing.T) {
	inst, err := BigDecimalClass.Call("new", "10")
	if err != nil {
		t.Errorf("unexpected error constructing: %v", err)
		return
	}
	bd := inst.(*BigDecimal)

	other, err := BigDecimalClass.Call("new", "5")
	if err != nil {
		t.Errorf("unexpected error constructing other: %v", err)
		return
	}

	result, err := bd.Call("compareTo", other)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if result != 1 {
		t.Errorf("expected 1, got %v", result)
	}
}

// Ensure an unknown BigDecimal class method returns an error.
func TestBigDecimal_unknownClassMethod(t *testing.T) {
	if _, err := BigDecimalClass.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown class method")
	}
}

// Ensure an unknown BigDecimal instance method returns an error.
func TestBigDecimal_unknownInstanceMethod(t *testing.T) {
	bd := &BigDecimal{V: decimal.Zero}
	if _, err := bd.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown instance method")
	}
}
