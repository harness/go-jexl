// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package math

import (
	"testing"

	gomath "math"
)

// Ensure Math.abs(-2.0) returns 2.0.
func TestMath_abs(t *testing.T) {
	got, err := MathClass.Call("abs", -2.0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 2.0 {
		t.Errorf("expected 2.0, got %v", got)
	}
}

// Ensure Math.PI returns the Go math.Pi constant.
func TestMath_PI(t *testing.T) {
	got, err := MathClass.Call("PI")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != gomath.Pi {
		t.Errorf("expected %v, got %v", gomath.Pi, got)
	}
}

// Ensure Math.pow(2.0, 10.0) returns 1024.0.
func TestMath_pow(t *testing.T) {
	got, err := MathClass.Call("pow", 2.0, 10.0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 1024.0 {
		t.Errorf("expected 1024.0, got %v", got)
	}
}

// Ensure Math.round(2.5) returns 3.0.
func TestMath_round(t *testing.T) {
	got, err := MathClass.Call("round", 2.5)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 3.0 {
		t.Errorf("expected 3.0, got %v", got)
	}
}

// Ensure Math.max(3.0, 7.0) returns 7.0.
func TestMath_max(t *testing.T) {
	got, err := MathClass.Call("max", 3.0, 7.0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 7.0 {
		t.Errorf("expected 7.0, got %v", got)
	}
}

// Ensure an unknown Math method returns an error.
func TestMath_unknownMethod(t *testing.T) {
	if _, err := MathClass.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown method")
	}
}
