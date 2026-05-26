// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package runtime

import "reflect"

// Function represents a callable function registered in the expression environment.
type Function struct {
	Name     string
	Fast     func(arg any) any
	Func     func(args ...any) (any, error)
	Safe     func(args ...any) (any, uint, error)
	Types    []reflect.Type
	Validate func(args []reflect.Type) (reflect.Type, error)
	Deref    func(i int, arg reflect.Type) bool
}

// Type returns the reflect.Type of the first overload, or the Func field type.
func (f *Function) Type() reflect.Type {
	if len(f.Types) > 0 {
		return f.Types[0]
	}
	return reflect.TypeOf(f.Func)
}
