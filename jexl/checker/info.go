// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package checker

import (
	"reflect"

	"github.com/harness/go-jexl/jexl/ast"
	"github.com/harness/go-jexl/jexl/checker/nature"
	"github.com/harness/go-jexl/jexl/vm"
)

// FieldIndex returns the struct field index path so the
// compiler can emit direct field access.
func FieldIndex(c *nature.Cache, env nature.Nature, node ast.Node) (bool, []int, string) {
	switch n := node.(type) {
	case *ast.Ident:
		if idx, ok := env.FieldIndex(c, n.Value); ok {
			return true, idx, n.Value
		}
	case *ast.MemberExpr:
		base := n.Node.Nature().Deref(c)
		if base.Kind == reflect.Struct {
			if prop, ok := n.Property.(*ast.StringLit); ok {
				if idx, ok := base.FieldIndex(c, prop.Value); ok {
					return true, idx, prop.Value
				}
			}
		}
	}
	return false, nil, ""
}

// MethodIndex returns the method index so the compiler
// can emit a direct method call.
func MethodIndex(c *nature.Cache, env nature.Nature, node ast.Node) (bool, int, string) {
	switch n := node.(type) {
	case *ast.Ident:
		if env.Kind == reflect.Struct {
			if m, ok := env.Get(c, n.Value); ok && m.TypeData != nil {
				return m.Method, m.MethodIndex, n.Value
			}
		}
	case *ast.MemberExpr:
		if name, ok := n.Property.(*ast.StringLit); ok {
			base := n.Node.Type()
			if base != nil && base.Kind() != reflect.Interface {
				if m, ok := base.MethodByName(name.Value); ok {
					return true, m.Index, name.Value
				}
			}
		}
	}
	return false, 0, ""
}

// TypedFuncIndex returns the vm.FuncTypes slot index for fn,
// enabling OpCallTyped dispatch (skips reflection at runtime).
// Returns false for variadic, named, or unmatched signatures.
func TypedFuncIndex(fn reflect.Type, method bool) (int, bool) {
	if fn == nil {
		return 0, false
	}
	if fn.Kind() != reflect.Func {
		return 0, false
	}
	// variadic and named func types have no slot in vm.FuncTypes
	if fn.IsVariadic() {
		return 0, false
	}
	if fn.PkgPath() != "" {
		return 0, false
	}

	fnNumIn := fn.NumIn()
	fnInOffset := 0
	if method {
		fnNumIn--
		fnInOffset = 1
	}

funcTypes:
	for i := range vm.FuncTypes {
		if i == 0 {
			continue
		}
		typed := reflect.ValueOf(vm.FuncTypes[i]).Elem().Type()
		if typed.Kind() != reflect.Func {
			continue
		}
		if typed.NumOut() != fn.NumOut() {
			continue
		}
		for j := 0; j < typed.NumOut(); j++ {
			if typed.Out(j) != fn.Out(j) {
				continue funcTypes
			}
		}
		if typed.NumIn() != fnNumIn {
			continue
		}
		for j := 0; j < typed.NumIn(); j++ {
			if typed.In(j) != fn.In(j+fnInOffset) {
				continue funcTypes
			}
		}
		return i, true
	}
	return 0, false
}

// IsFastFunc reports whether fn matches the fast-call
// signature: variadic func(...any) any.
func IsFastFunc(fn reflect.Type, method bool) bool {
	if fn == nil {
		return false
	}
	if fn.Kind() != reflect.Func {
		return false
	}
	numIn := 1
	if method {
		numIn = 2
	}
	if fn.IsVariadic() &&
		fn.NumIn() == numIn &&
		fn.NumOut() == 1 &&
		fn.Out(0).Kind() == reflect.Interface {
		rest := fn.In(fn.NumIn() - 1)
		if rest != nil &&
			rest.Kind() == reflect.Slice &&
			rest.Elem().Kind() == reflect.Interface {
			return true
		}
	}
	return false
}
