// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package tests

import (
	"testing"

	expr "github.com/harness/go-jexl/jexl"
	javalang "github.com/harness/go-jexl/jexl/classes/java/lang"
)

// TestClass_booleanLogicalAnd calls the static method java.lang.Boolean.logicalAnd via dot notation.
func TestClass_booleanLogicalAnd(t *testing.T) {
	program, err := expr.Compile(
		`java.lang.Boolean.logicalAnd(true, false)`,
		expr.WithClass("java.lang.Boolean", javalang.BooleanClass),
	)
	if err != nil {
		t.Fatal(err)
	}
	out, err := expr.Run(program, nil)
	if err != nil {
		t.Fatal(err)
	}
	if out != false {
		t.Fatalf("expected false, got %v", out)
	}
}

// TestContext_func invokes a plain function stored as a context value.
func TestContext_func(t *testing.T) {
	ctx := map[string]any{
		"secrets": map[string]any{
			"getValue": func(args ...any) (any, error) {
				return "dummy-23e4567-e89b-12d3-a456-426614174000", nil
			},
		},
	}
	program, err := expr.Compile(`secrets.getValue('account.token')`, expr.WithContext(ctx))
	if err != nil {
		t.Fatal(err)
	}
	out, err := expr.Run(program, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if out != "dummy-23e4567-e89b-12d3-a456-426614174000" {
		t.Fatalf("expected dummy-23e4567-e89b-12d3-a456-426614174000, got %v", out)
	}
}

// TestBuiltin_toUpperCase verifies that the default toUpperCase builtin is available.
func TestBuiltin_toUpperCase(t *testing.T) {
	program, err := expr.Compile(`toUpperCase("hello")`)
	if err != nil {
		t.Fatal(err)
	}
	out, err := expr.Run(program, nil)
	if err != nil {
		t.Fatal(err)
	}
	if out != "HELLO" {
		t.Fatalf("expected HELLO, got %v", out)
	}
}
