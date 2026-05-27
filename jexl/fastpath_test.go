// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package jexl

import (
	"errors"
	"testing"
)

func TestIsPropertyPath(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		// valid
		{"foo", true},
		{"a.b.c", true},
		{"my_var.nested_field", true},
		{"trigger.payload.crNumber", true},
		// invalid
		{"", false},
		{"1foo", false},
		{".foo", false},
		{"foo.", false},
		{"a..b", false},
		{"a .b", false},
		{"a + b", false},
		{"a[0]", false},
		{"résumé", false},
	}
	for _, tc := range tests {
		got := isPropertyPath(tc.input)
		if got != tc.want {
			t.Fatalf("isPropertyPath(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestEvalPath(t *testing.T) {
	nested := map[string]any{
		"trigger": map[string]any{
			"payload": map[string]any{
				"crNumber": "CR-123",
			},
		},
	}
	deep := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": map[string]any{
					"d": map[string]any{
						"e": "deep",
					},
				},
			},
		},
	}
	tests := []struct {
		env  any
		path string
		want any
	}{
		// single segment
		{map[string]any{"name": "Alice"}, "name", "Alice"},
		// three-level nested
		{nested, "trigger.payload.crNumber", "CR-123"},
		// missing key returns nil
		{map[string]any{"a": 1}, "b", nil},
		// nil intermediate returns nil
		{map[string]any{"a": nil}, "a.b", nil},
		// five levels deep
		{deep, "a.b.c.d.e", "deep"},
	}
	for _, tc := range tests {
		got, err := evalPath(tc.env, tc.path)
		if err != nil {
			t.Fatalf("evalPath(%q): %v", tc.path, err)
		}
		if got != tc.want {
			t.Fatalf("evalPath(%q) = %v, want %v", tc.path, got, tc.want)
		}
	}
}

func TestEvalPath_notPropertyPath(t *testing.T) {
	// expression with operator returns ErrNotPropertyPath
	_, err := EvalPath("a + b", nil)
	if !errors.Is(err, ErrNotPropertyPath) {
		t.Fatalf("expected ErrNotPropertyPath, got %v", err)
	}
}

func TestEvalPath_emptyExpression(t *testing.T) {
	// empty string returns ErrNotPropertyPath
	_, err := EvalPath("", nil)
	if !errors.Is(err, ErrNotPropertyPath) {
		t.Fatalf("expected ErrNotPropertyPath, got %v", err)
	}
}

func TestEvalPath_validPath(t *testing.T) {
	// valid path returns correct value
	env := map[string]any{"pipeline": map[string]any{"name": "build"}}
	got, err := EvalPath("pipeline.name", env)
	if err != nil {
		t.Fatal(err)
	}
	if got != "build" {
		t.Fatalf("expected build, got %v", got)
	}
}

//
// Benchmarks
//

func BenchmarkEvalPath(b *testing.B) {
	env := map[string]any{
		"trigger": map[string]any{
			"payload": map[string]any{
				"crNumber": "CR-123",
			},
		},
	}
	var out any
	var err error
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = EvalPath("trigger.payload.crNumber", env)
	}
	b.StopTimer()
	if err != nil {
		b.Fatal(err)
	}
	if out != "CR-123" {
		b.Fatalf("expected CR-123, got %v", out)
	}
}
