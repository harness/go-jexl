// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package config

import (
	"fmt"
	"reflect"

	"github.com/harness/go-jexl/jexl/builtin"
	"github.com/harness/go-jexl/jexl/checker/nature"
	"github.com/harness/go-jexl/jexl/classes"
	"github.com/harness/go-jexl/jexl/internal/deref"
	"github.com/harness/go-jexl/jexl/vm/runtime"
)

var (
	// MaxIterations is the maximium number of loop
	// iterations the virtual machine can perform.
	MaxIterations uint = 1e4

	// MaxMemory is the maximum amount of memory available
	// to the virtual machine.
	MaxMemory uint = 1e6

	// MaxNodes is the maximium number of nodes the
	// compiler can process.
	MaxNodes uint = 1e4
)

type Config struct {
	Context       nature.Nature
	Cache         nature.Cache
	Functions     map[string]*runtime.Function
	MaxIterations uint
	MaxMemory     uint
	MaxNodes      uint
	Registry      *classes.Registry
	Safe          bool
	Strict        bool
}

// New creates new config with default values.
func New() *Config {
	c := &Config{
		MaxNodes:      MaxNodes,
		MaxMemory:     MaxMemory,
		MaxIterations: MaxIterations,
		Functions:     make(map[string]*runtime.Function, len(builtin.Funcs)),
		Registry:      classes.New(),
	}
	for name, fn := range builtin.Funcs {
		c.Functions[name] = fn
	}
	return c
}

// NewCache creates new config with context.
func NewCache(ctx any) *Config {
	c := New()
	c.SetCache(ctx)
	return c
}

func (c *Config) SetCache(ctx any) {
	c.Context = contextWithCache(&c.Cache, ctx)
}

func contextWithCache(c *nature.Cache, ctx any) nature.Nature {
	if ctx == nil {
		n := c.NatureOf(map[string]any{})
		return n
	}

	v := reflect.ValueOf(ctx)
	t := v.Type()

	switch deref.Value(v).Kind() {
	case reflect.Struct:
		n := c.FromType(t)
		return n

	case reflect.Map:
		n := c.FromType(v.Type())
		if n.TypeData == nil {
			n.TypeData = new(nature.TypeData)
		}
		n.Fields = make(map[string]nature.Nature, v.Len())

		for _, key := range v.MapKeys() {
			elem := v.MapIndex(key)
			if !elem.IsValid() || !elem.CanInterface() {
				panic(fmt.Sprintf("invalid map value: %s", key))
			}

			face := elem.Interface()

			if face == nil {
				n.Fields[key.String()] = c.NatureOf(nil)
				continue
			}
			n.Fields[key.String()] = c.NatureOf(face)

		}

		return n
	}

	panic(fmt.Sprintf("unknown type %T", ctx))
}
