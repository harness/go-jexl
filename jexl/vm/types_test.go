// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package vm

import (
	"reflect"
	"testing"
)

// Ensure Item returns from the int fast path.
func TestScope_Item_ints(t *testing.T) {
	s := &Scope{Ints: []int{10, 20, 30}, Len: 3, Index: 1}
	if s.Item() != 20 {
		t.Fatalf("expected 20, got %v", s.Item())
	}
}

// Ensure Item returns from the float fast path.
func TestScope_Item_floats(t *testing.T) {
	s := &Scope{Floats: []float64{1.1, 2.2, 3.3}, Len: 3, Index: 2}
	if s.Item() != 3.3 {
		t.Fatalf("expected 3.3, got %v", s.Item())
	}
}

// Ensure Item returns from the string fast path.
func TestScope_Item_strings(t *testing.T) {
	s := &Scope{Strings: []string{"a", "b", "c"}, Len: 3, Index: 0}
	if s.Item() != "a" {
		t.Fatalf("expected a, got %v", s.Item())
	}
}

// Ensure Item returns from the any fast path.
func TestScope_Item_anys(t *testing.T) {
	s := &Scope{Anys: []any{"x", 42, true}, Len: 3, Index: 1}
	if s.Item() != 42 {
		t.Fatalf("expected 42, got %v", s.Item())
	}
}

// Ensure Item falls back to reflect.Value when no typed slice is set.
func TestScope_Item_reflectFallback(t *testing.T) {
	arr := []string{"foo", "bar"}
	s := &Scope{Array: reflect.ValueOf(arr), Len: 2, Index: 1}
	if s.Item() != "bar" {
		t.Fatalf("expected bar, got %v", s.Item())
	}
}

// Ensure ThrownValue.Error formats the value.
func TestThrownValue_Error(t *testing.T) {
	tv := &ThrownValue{Value: "something went wrong"}
	if tv.Error() != "thrown: something went wrong" {
		t.Fatalf("unexpected: %q", tv.Error())
	}
}

// Ensure ThrownValue.Error handles non-string values.
func TestThrownValue_Error_nonString(t *testing.T) {
	tv := &ThrownValue{Value: 42}
	if tv.Error() != "thrown: 42" {
		t.Fatalf("unexpected: %q", tv.Error())
	}
}
