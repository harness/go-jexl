// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package functions

import "reflect"

var variadicAnyError = reflect.TypeOf((*func(...any) (any, error))(nil)).Elem()

// Wrap adapts any typed Go function into the func(...any)(any,error)
// signature the vm expects. If fn is already that type it is returned as-is.
func Wrap(fn any) (func(...any) (any, error), reflect.Type) {
	t := reflect.TypeOf(fn)
	if t == variadicAnyError {
		return fn.(func(...any) (any, error)), t
	}
	rv := reflect.ValueOf(fn)
	isVariadic := t.IsVariadic()
	numIn := t.NumIn()
	wrapped := func(args ...any) (any, error) {
		in := make([]reflect.Value, len(args))
		for i, arg := range args {
			var paramType reflect.Type
			if isVariadic && i >= numIn-1 {
				paramType = t.In(numIn - 1).Elem()
			} else {
				paramType = t.In(i)
			}
			if arg == nil {
				in[i] = reflect.Zero(paramType)
			} else {
				v := reflect.ValueOf(arg)
				if v.Type().ConvertibleTo(paramType) && v.Type() != paramType {
					v = v.Convert(paramType)
				}
				in[i] = v
			}
		}
		out := rv.Call(in)
		if len(out) == 2 {
			if err, ok := out[1].Interface().(error); ok && err != nil {
				return nil, err
			}
		}
		return out[0].Interface(), nil
	}
	return wrapped, t
}
