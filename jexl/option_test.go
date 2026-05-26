// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package jexl

import (
	"testing"

	"github.com/harness/go-jexl/jexl/classes/java/lang"
	"github.com/harness/go-jexl/jexl/classes/java/math"
	"github.com/harness/go-jexl/jexl/config"
)

// Ensure the context is configured.
func TestWithContext(t *testing.T) {
	conf := config.New()
	WithContext(map[string]any{"x": 1})(conf)
	if conf.Context.Type == nil {
		t.Fatal("expected Context to be set")
	}
}

// Ensure strict mode is configured.
func TestWithStrict(t *testing.T) {
	conf := config.New()
	WithStrict()(conf)
	if !conf.Strict {
		t.Fatal("expected Strict=true")
	}
}

// Ensure function is registered.
func TestWithFunction(t *testing.T) {
	conf := config.New()
	fn := func(args ...any) (any, error) { return nil, nil }
	WithFunction("double", fn, func(int) int { return 0 })(conf)
	if _, ok := conf.Functions["double"]; !ok {
		t.Fatal("expected double to be registered")
	}
}

// Ensure panic if not a function.
func TestWithFunctionPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for non-function type")
		}
	}()
	fn := func(args ...any) (any, error) { return nil, nil }
	WithFunction("test", fn, 42)(config.New())
}

// Ensure the max nodes are configured.
func TestWithMaxNodes(t *testing.T) {
	conf := config.New()
	WithMaxNodes(100)(conf)
	if conf.MaxNodes != 100 {
		t.Fatalf("expected MaxNodes 100, got %d", conf.MaxNodes)
	}
}

// Ensure the max iterations are configured.
func TestWithMaxIterations(t *testing.T) {
	conf := config.New()
	WithMaxIterations(50)(conf)
	if conf.MaxIterations != 50 {
		t.Fatalf("expected MaxIterations 50, got %d", conf.MaxIterations)
	}
}

// Ensure the max memory is configured.
func TestWithMaxMemory(t *testing.T) {
	conf := config.New()
	WithMaxMemory(512)(conf)
	if conf.MaxMemory != 512 {
		t.Fatalf("expected MaxMemory 512, got %d", conf.MaxMemory)
	}
}

// Ensure the class is registered.
func TestWithClass(t *testing.T) {
	conf := config.New()
	WithClass("java.lang.String", lang.StringClass)(conf)
	if _, ok := conf.Registry.Lookup("java.lang.String"); !ok {
		t.Fatal("expected java.lang.String to be registered")
	}
}

// Ensure the namespace is registered.
func TestWithNamespace(t *testing.T) {
	conf := config.New()
	WithNamespace("Math", math.MathClass)(conf)
	if _, ok := conf.Registry.Lookup("Math"); !ok {
		t.Fatal("expected Math to be registered")
	}
}
