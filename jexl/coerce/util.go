// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package coerce

import (
	"math"
	"reflect"
	"strconv"
	"strings"
)

// ToStringJexl converts a value to string following JEXL/Java conventions.
// Float64 whole numbers are displayed as "N.0" (Java Double.toString behavior).
func ToStringJexl(v any) string {
	if f, ok := v.(float64); ok {
		if math.IsNaN(f) {
			return "NaN"
		}
		if math.IsInf(f, 1) {
			return "Infinity"
		}
		if math.IsInf(f, -1) {
			return "-Infinity"
		}
		s := strconv.FormatFloat(f, 'f', -1, 64)
		if !strings.Contains(s, ".") {
			s += ".0"
		}
		return s
	}
	return ToString(v)
}

// DeepEqual reports whether a and b are equal, coercing numeric types before comparing.
func DeepEqual(a, b any) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}

	ra := reflect.ValueOf(a)
	rb := reflect.ValueOf(b)
	switch ra.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ToInt64(a) == ToInt64(b)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return ToUint64(a) == ToUint64(b)
	case reflect.Float32, reflect.Float64:
		return ToFloat64(a) == ToFloat64(b)
	case reflect.String:
		return ToString(a) == ToString(b)
	case reflect.Bool:
		return ToBool(a) == ToBool(b)
	case reflect.Map:
		switch rb.Kind() {
		case reflect.Map:
			if ra.Len() != rb.Len() {
				return false
			}
			iter := ra.MapRange()
			for iter.Next() {
				k := iter.Key()
				bv := rb.MapIndex(k)
				if !bv.IsValid() {
					return false
				}
				if !DeepEqual(
					iter.Value().Interface(),
					bv.Interface(),
				) {
					return false
				}
			}
			return true
		default:
			return false
		}
	case reflect.Slice:
		switch rb.Kind() {
		case reflect.Slice:
			if ra.Len() != rb.Len() {
				return false
			}
			for i := 0; i < ra.Len(); i++ {
				if !DeepEqual(
					ra.Index(i).Interface(),
					rb.Index(i).Interface(),
				) {
					return false
				}
			}
			return true
		default:
			return false
		}
	default:
		return reflect.DeepEqual(a, b)
	}
}
