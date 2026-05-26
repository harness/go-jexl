// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import (
	"reflect"

	"github.com/harness/go-jexl/jexl/coerce"
)

// BooleanValue returns v as a bool.
func BooleanValue(v any) bool {
	return coerce.ToBool(v)
}

// ByteValue returns v as int8.
func ByteValue(v any) int8 {
	return coerce.ToInt8(v)
}

// DoubleValue returns v as float64.
func DoubleValue(v any) float64 {
	return coerce.ToFloat64(v)
}

// FloatValue returns v as float32.
func FloatValue(v any) float32 {
	return coerce.ToFloat32(v)
}

// IntValue returns v as int.
func IntValue(v any) int {
	return coerce.ToInt(v)
}

// LongValue returns v as int64.
func LongValue(v any) int64 {
	return coerce.ToInt64(v)
}

// ShortValue returns v as int16.
func ShortValue(v any) int16 {
	return coerce.ToInt16(v)
}

// ToString returns v as a string.
func ToString(v any) string {
	return coerce.ToString(v)
}

// DefaultValue returns fallback when v is nil or the zero value of its type.
func DefaultValue(v any, fallback any) any {
	if v == nil {
		return fallback
	}
	rv := reflect.ValueOf(v)
	if rv.IsZero() {
		return fallback
	}
	return v
}
