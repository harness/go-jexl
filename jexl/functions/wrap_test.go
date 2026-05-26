// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package functions

import (
	"fmt"
	"testing"
)

// Ensure a func(...any)(any,error) is returned as-is.
func TestWrap_variadicAny(t *testing.T) {
	fn := func(args ...any) (any, error) { return args[0], nil }
	wrapped, _ := Wrap(fn)
	got, err := wrapped("hello")
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello" {
		t.Fatalf("expected hello, got %v", got)
	}
}

// Ensure a typed function with no return error is wrapped correctly.
func TestWrap_typedNoError(t *testing.T) {
	fn := func(x int) int { return x * 2 }
	wrapped, _ := Wrap(fn)
	got, err := wrapped(3)
	if err != nil {
		t.Fatal(err)
	}
	if got != 6 {
		t.Fatalf("expected 6, got %v", got)
	}
}

// Ensure a typed function that returns an error propagates it.
func TestWrap_typedWithError(t *testing.T) {
	fn := func(x string) (string, error) {
		if x == "" {
			return "", fmt.Errorf("empty")
		}
		return x, nil
	}
	wrapped, _ := Wrap(fn)
	if _, err := wrapped(""); err == nil {
		t.Fatal("expected error, got nil")
	}
	got, err := wrapped("ok")
	if err != nil {
		t.Fatal(err)
	}
	if got != "ok" {
		t.Fatalf("expected ok, got %v", got)
	}
}

// Ensure a typed variadic function is wrapped correctly.
func TestWrap_typedVariadic(t *testing.T) {
	fn := func(args ...string) string {
		result := ""
		for _, a := range args {
			result += a
		}
		return result
	}
	wrapped, _ := Wrap(fn)
	got, err := wrapped("a", "b", "c")
	if err != nil {
		t.Fatal(err)
	}
	if got != "abc" {
		t.Fatalf("expected abc, got %v", got)
	}
}

// Ensure a nil argument is converted to the zero value of the parameter type.
func TestWrap_nilArgument(t *testing.T) {
	fn := func(x string) string { return "[" + x + "]" }
	wrapped, _ := Wrap(fn)
	got, err := wrapped(nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "[]" {
		t.Fatalf("expected [], got %v", got)
	}
}

// Ensure Wrap returns the concrete reflect.Type of the original function.
func TestWrap_returnsType(t *testing.T) {
	fn := func(x int) int { return x }
	_, typ := Wrap(fn)
	if typ.NumIn() != 1 || typ.In(0).Kind().String() != "int" {
		t.Fatalf("unexpected type: %v", typ)
	}
}
