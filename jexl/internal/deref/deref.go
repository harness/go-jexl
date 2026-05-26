// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package deref

import (
	"reflect"
)

// Interface dereferences a pointer or interface chain and
// returns the underlying value. Returns nil if p is nil or
// contains a nil pointer.
func Interface(p any) any {
	if p == nil {
		return nil
	}

	v := reflect.ValueOf(p)

	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	return v.Interface()
}

// Type dereferences a pointer type chain and returns the
// base element type.
func Type(t reflect.Type) reflect.Type {
	if t == nil {
		return nil
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// Value dereferences a pointer or interface value chain
// and returns the base value. Returns the nil pointer
// value if a nil is encountered.
func Value(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return v
		}
		v = v.Elem()
	}
	return v
}

// TypeKind dereferences a pointer type chain and returns
// the base type, its kind, and whether any dereferencing
// occurred.
func TypeKind(t reflect.Type, k reflect.Kind) (_ reflect.Type, _ reflect.Kind, changed bool) {
	for k == reflect.Pointer {
		changed = true
		t = t.Elem()
		k = t.Kind()
	}
	return t, k, changed
}
