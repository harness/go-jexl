// Copyright (c) 2014 Steve Francia
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package cast

import (
	"reflect"
	"slices"
)

var kindNames = []string{
	reflect.String:  "string",
	reflect.Bool:    "bool",
	reflect.Int:     "int",
	reflect.Int8:    "int8",
	reflect.Int16:   "int16",
	reflect.Int32:   "int32",
	reflect.Int64:   "int64",
	reflect.Uint:    "uint",
	reflect.Uint8:   "uint8",
	reflect.Uint16:  "uint16",
	reflect.Uint32:  "uint32",
	reflect.Uint64:  "uint64",
	reflect.Float32: "float32",
	reflect.Float64: "float64",
}

var kinds = map[reflect.Kind]func(reflect.Value) any{
	reflect.String:  func(v reflect.Value) any { return v.String() },
	reflect.Bool:    func(v reflect.Value) any { return v.Bool() },
	reflect.Int:     func(v reflect.Value) any { return int(v.Int()) },
	reflect.Int8:    func(v reflect.Value) any { return int8(v.Int()) },
	reflect.Int16:   func(v reflect.Value) any { return int16(v.Int()) },
	reflect.Int32:   func(v reflect.Value) any { return int32(v.Int()) },
	reflect.Int64:   func(v reflect.Value) any { return v.Int() },
	reflect.Uint:    func(v reflect.Value) any { return uint(v.Uint()) },
	reflect.Uint8:   func(v reflect.Value) any { return uint8(v.Uint()) },
	reflect.Uint16:  func(v reflect.Value) any { return uint16(v.Uint()) },
	reflect.Uint32:  func(v reflect.Value) any { return uint32(v.Uint()) },
	reflect.Uint64:  func(v reflect.Value) any { return v.Uint() },
	reflect.Float32: func(v reflect.Value) any { return float32(v.Float()) },
	reflect.Float64: func(v reflect.Value) any { return v.Float() },
}

// resolveAlias unwraps a named type (e.g. type Double float64) to its
// underlying primitive so type-switch conversions can proceed.
func resolveAlias(i any) (any, bool) {
	if i == nil {
		return nil, false
	}
	t := reflect.TypeOf(i)
	if t.Name() == "" || slices.Contains(kindNames, t.Name()) {
		return i, false
	}
	resolve, ok := kinds[t.Kind()]
	if !ok {
		return i, false
	}
	return resolve(reflect.ValueOf(i)), true
}
