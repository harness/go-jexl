// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package jexl

import (
	"fmt"
	"reflect"

	"github.com/harness/go-jexl/jexl/classes"
	"github.com/harness/go-jexl/jexl/config"
	"github.com/harness/go-jexl/jexl/functions"
	"github.com/harness/go-jexl/jexl/vm/runtime"
)

// Option configures a jexl compiler option.
type Option func(*config.Config)

// WithClass registers a class under a given name.
func WithClass(name string, obj classes.Object) Option {
	return func(c *config.Config) {
		c.Registry.Register(name, obj)
	}
}

// WithContext returns an option to set the context for
// type checking. Struct fields and map keys are treated
// as variables. Methods on the type are available as
// functions.
func WithContext(ctx any) Option {
	return func(c *config.Config) {
		c.SetCache(ctx)
	}
}

// WithFunction returns an option to register a named
// function for use in expressions.
func WithFunction(name string, fn any) Option {
	return func(c *config.Config) {
		t := reflect.TypeOf(fn)
		if t == nil || t.Kind() != reflect.Func {
			panic(fmt.Sprintf("jexl: %s is not a function", name))
		}
		wrapped, typ := functions.Wrap(fn)
		c.Functions.Register(name, &runtime.Function{
			Name:  name,
			Func:  wrapped,
			Types: []reflect.Type{typ},
		})
	}
}

// WithFunctionNamespace registers a single function under a dot-separated
// namespace, making it callable as ns.method(...) in expressions.
func WithFunctionNamespace(ns, method string, fn any) Option {
	return func(c *config.Config) {
		t := reflect.TypeOf(fn)
		if t == nil || t.Kind() != reflect.Func {
			panic(fmt.Sprintf("jexl: %s.%s is not a function", ns, method))
		}
		wrapped, typ := functions.Wrap(fn)
		c.Functions.RegisterNamespace(ns, method, &runtime.Function{
			Name:  ns + "." + method,
			Types: []reflect.Type{typ},
			Func:  wrapped,
		})
	}
}

// WithMaxIterations returns an option to set the maximum
// number of loop iterations. A value of zero disables the limit.
func WithMaxIterations(n uint) Option {
	return func(c *config.Config) {
		c.MaxIterations = n
	}
}

// WithMaxNodes returns an option to set the maximum number
// of nodes allowed in the expression. A value of zero
// disables the limit.
func WithMaxNodes(n uint) Option {
	return func(c *config.Config) {
		c.MaxNodes = n
	}
}

// MaxMemory returns an option to set the maximum amount
// of memory allocated to the virtual machine. A value of
// zero disables the limit.
func WithMaxMemory(n uint) Option {
	return func(c *config.Config) {
		c.MaxMemory = n
	}
}

// WithNamespace registers a class under a given namespace.
func WithNamespace(name string, obj classes.Object) Option {
	return func(c *config.Config) {
		c.Registry.Register(name, obj)
	}
}

// WithStrict returns an option to enable strict mode.
// This forces the compiler to error when it detects undefined
// variables in the expression.
func WithStrict() Option {
	return func(c *config.Config) {
		c.Strict = true
	}
}

// WithSafe returns an option to enable safe mode.
// This prevents expressions from mutating the context passed by the caller.
func WithSafe() Option {
	return func(c *config.Config) {
		c.Safe = true
	}
}
