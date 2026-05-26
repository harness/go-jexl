// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"
	"reflect"
)

type (
	// Function is the custom function signature.
	Function = func(params ...any) (any, error)

	// AntishCursor tracks a dotted-path lookup in a flat
	// map. In antish mode, context key "a.b" is accessed
	// via expression a.b.
	AntishCursor struct {
		Env  map[string]any
		Path string // accumulated dotted path so far
	}

	// Field is a cached struct field descriptor.
	Field struct {
		Index []int
		Path  []string
	}

	// Method is a cached method descriptor.
	Method struct {
		Index int
		Name  string
	}

	// CaptureVar describes a captured outer variable.
	CaptureVar struct {
		Name string // variable name
		Slot int    // Variables slot in the enclosing program
	}

	// Closure is a lambda: a compiled body, parameter
	// names, and captured variable values. Program is
	// stored as any to avoid a circular import; callers
	// assert *vm.Program.
	Closure struct {
		Params      []string
		Captures    map[string]any
		CaptureVars []CaptureVar // set at compile time; populated into Captures at OpLambda execution
		Program     any
	}

	// CompoundAssignOp is stored in program.Constants for
	// OpCompoundAssign. VarIndex is the Variables slot;
	// Op is the operator string (e.g. "+=").
	CompoundAssignOp struct {
		VarIndex int
		Op       string
	}

	// CompoundEnvAssignOp is stored in program.Constants
	// for OpCompoundStoreEnv.
	CompoundEnvAssignOp struct {
		Key string
		Op  string
	}
)

//
// Thrown Value
//

// ThrownValue carries a value thrown by a throw
// statement through the try-frame stack.
type ThrownValue struct {
	Value any
}

// Error converns the thrown value to a String.
func (t *ThrownValue) Error() string {
	return fmt.Sprintf("thrown: %v", t.Value)
}

//
// Scope
//

// Scope holds iteration state for a foreach loop.
type Scope struct {
	Array reflect.Value
	Index int
	Len   int
	Count int
	Acc   any
	// VarSlot is the Variables index for the foreach
	// loop variable (-1 = not a foreach scope).
	VarSlot int
	// Fast paths
	Ints    []int
	Floats  []float64
	Strings []string
	Anys    []any
}

// Item returns the current element from the scope using
// fast paths when available.
func (s *Scope) Item() any {
	switch {
	case s.Ints != nil:
		return s.Ints[s.Index]
	case s.Floats != nil:
		return s.Floats[s.Index]
	case s.Strings != nil:
		return s.Strings[s.Index]
	case s.Anys != nil:
		return s.Anys[s.Index]
	default:
		return s.Array.Index(s.Index).Interface()
	}
}
