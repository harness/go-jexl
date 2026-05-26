// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package eval

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"regexp"
	"strings"

	"github.com/harness/go-jexl/jexl/classes/java/util"
	"github.com/harness/go-jexl/jexl/coerce"
	"github.com/harness/go-jexl/jexl/internal/decimal"
)

// BigCompare returns the sign of a vs b when either is
// BigInteger or BigDecimal.
func BigCompare(a, b any) (int, bool) {
	_, aIsBigInt := coerce.ToBigInt(a)
	_, bIsBigInt := coerce.ToBigInt(b)
	_, aIsBigDec := coerce.ToBigDecimal(a)
	_, bIsBigDec := coerce.ToBigDecimal(b)

	if !aIsBigInt && !bIsBigInt && !aIsBigDec && !bIsBigDec {
		return 0, false
	}

	return coerce.ToDecimal(a).Cmp(coerce.ToDecimal(b)), true
}

// BitwiseAnd returns a & b as int64.
func BitwiseAnd(a, b any) int64 {
	return coerce.ToInt64(a) & coerce.ToInt64(b)
}

// BitwiseNot returns ^a as int64.
func BitwiseNot(a any) int64 {
	return ^coerce.ToInt64(a)
}

// BitwiseOr returns a | b as int64.
func BitwiseOr(a, b any) int64 {
	return coerce.ToInt64(a) | coerce.ToInt64(b)
}

// BitwiseXor returns a ^ b as int64.
func BitwiseXor(a, b any) int64 {
	return coerce.ToInt64(a) ^ coerce.ToInt64(b)
}

// Decrement returns a decremented by one, preserving
// the original type.
func Decrement(a any) (any, error) {
	switch v := a.(type) {
	case int:
		return v - 1, nil
	case int8:
		return v - 1, nil
	case int16:
		return v - 1, nil
	case int32:
		return v - 1, nil
	case int64:
		return v - 1, nil
	case uint:
		return v - 1, nil
	case uint8:
		return v - 1, nil
	case uint16:
		return v - 1, nil
	case uint32:
		return v - 1, nil
	case uint64:
		return v - 1, nil
	case float32:
		return v - 1, nil
	case float64:
		return v - 1, nil
	case *big.Int:
		return new(big.Int).Sub(v, big.NewInt(1)), nil
	case decimal.Decimal:
		return v.Sub(decimal.NewFromInt(1)), nil
	default:
		return nil, fmt.Errorf("invalid operation: -- on %T", a)
	}
}

// DivideInt performs integer division. Exact quotients
// return float64; inexact return int64.
func DivideInt(a, b any) any {
	ai := coerce.ToInt64(a)
	bi := coerce.ToInt64(b)
	if ai%bi == 0 {
		return float64(ai / bi)
	}
	return ai / bi
}

// Empty reports whether x is nil, zero, or an empty collection.
func Empty(x any) bool {
	if x == nil {
		return true
	}
	switch v := x.(type) {
	case string:
		return v == ""
	case bool:
		return false
	case int:
		return v == 0
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0
	case uint:
		return v == 0
	case uint8:
		return v == 0
	case uint16:
		return v == 0
	case uint32:
		return v == 0
	case uint64:
		return v == 0
	case float32:
		return v == 0
	case float64:
		return v == 0 || math.IsNaN(v)
	default:
		r := reflect.ValueOf(x)
		switch r.Kind() {
		case reflect.Array, reflect.Slice, reflect.Map:
			return r.Len() == 0
		case reflect.Ptr, reflect.Interface:
			return r.IsNil()
		}
		return false
	}
}

// Exponent returns a**b as float64.
func Exponent(a, b any) float64 {
	return math.Pow(
		coerce.ToFloat64(a),
		coerce.ToFloat64(b),
	)
}

// In reports whether needle is in array.
func In(needle any, array any) bool {
	if array == nil {
		return false
	}
	if s, ok := array.(util.Set); ok {
		return s.Contains(needle)
	}
	v := reflect.ValueOf(array)
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			value := v.Index(i)
			if value.IsValid() {
				if Equal(value.Interface(), needle) {
					return true
				}
			}
		}
		return false
	case reflect.Map:
		var value reflect.Value
		if needle == nil {
			value = v.MapIndex(reflect.Zero(v.Type().Key()))
		} else {
			value = v.MapIndex(reflect.ValueOf(needle))
		}
		return value.IsValid()
	case reflect.Struct:
		n := reflect.ValueOf(needle)
		if !n.IsValid() || n.Kind() != reflect.String {
			panic(fmt.Sprintf("cannot use %T as field name of %T", needle, array))
		}
		field, ok := v.Type().FieldByName(n.String())
		if !ok || !field.IsExported() || field.Tag.Get("expr") == "-" {
			return false
		}
		return v.FieldByIndex(field.Index).IsValid()
	case reflect.Ptr:
		value := v.Elem()
		if value.IsValid() {
			return In(needle, value.Interface())
		}
		return false
	}
	panic(fmt.Sprintf(`operator "in" not defined on %T`, array))
}

// Increment returns a incremented by one, preserving the
// original type.
func Increment(a any) (any, error) {
	switch v := a.(type) {
	case int:
		return v + 1, nil
	case int8:
		return v + 1, nil
	case int16:
		return v + 1, nil
	case int32:
		return v + 1, nil
	case int64:
		return v + 1, nil
	case uint:
		return v + 1, nil
	case uint8:
		return v + 1, nil
	case uint16:
		return v + 1, nil
	case uint32:
		return v + 1, nil
	case uint64:
		return v + 1, nil
	case float32:
		return v + 1, nil
	case float64:
		return v + 1, nil
	case *big.Int:
		return new(big.Int).Add(v, big.NewInt(1)), nil
	case decimal.Decimal:
		return v.Add(decimal.NewFromInt(1)), nil
	default:
		return nil, fmt.Errorf("invalid operation: ++ on %T", a)
	}
}

// IsFalsy reports whether v is falsy: nil, false, zero,
// empty string, or "false".
func IsFalsy(v any) bool {
	if v == nil {
		return true
	}
	switch x := v.(type) {
	case bool:
		return !x
	case int:
		return x == 0
	case int8:
		return x == 0
	case int16:
		return x == 0
	case int32:
		return x == 0
	case int64:
		return x == 0
	case uint:
		return x == 0
	case uint8:
		return x == 0
	case uint16:
		return x == 0
	case uint32:
		return x == 0
	case uint64:
		return x == 0
	case float32:
		return x == 0
	case float64:
		return x == 0
	case string:
		return x == "" || strings.EqualFold(x, "false")
	default:
		return false
	}
}

// IsNil reports whether v is nil, including typed
// nils (pointer, slice, map, etc.).
func IsNil(v any) bool {
	if v == nil {
		return true
	}
	r := reflect.ValueOf(v)
	switch r.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return r.IsNil()
	default:
		return false
	}
}

// LooseEqual reports whether a equals b using JEXL
// loose-equality rules.
func LooseEqual(a, b any) bool {
	if cmp, ok := BigCompare(a, b); ok {
		return cmp == 0
	}
	aStr, aIsStr := a.(string)
	bStr, bIsStr := b.(string)
	if aIsStr && !bIsStr && b != nil {
		return aStr == coerce.ToString(b)
	}
	if bIsStr && !aIsStr && a != nil {
		return coerce.ToString(a) == bStr
	}
	return Equal(a, b)
}

// Match reports whether a matches b, where b is a
// *regexp.Regexp or pattern string.
func Match(a, b any) bool {
	switch r := b.(type) {
	case *regexp.Regexp:
		fullPattern := "^(?:" + r.String() + ")$"
		fr, err := regexp.Compile(fullPattern)
		if err != nil {
			fr = r
		}
		if s, ok := a.(string); ok {
			return fr.MatchString(s)
		}
		if bs, ok := a.([]byte); ok {
			return fr.Match(bs)
		}
		return false
	case string:
		return MatchFull(a, r)
	}
	return false
}

// MatchFull reports whether a fully matches the regex
// pattern string.
func MatchFull(a any, pattern string) bool {
	anchored := "^(?:" + pattern + ")$"
	r, err := regexp.Compile(anchored)
	if err != nil {
		panic(err)
	}
	if s, ok := a.(string); ok {
		return r.MatchString(s)
	}
	if b, ok := a.([]byte); ok {
		return r.Match(b)
	}
	return false
}

// Negate returns -i, or logical NOT if i is bool.
func Negate(i any) any {
	switch v := i.(type) {
	case float32:
		return -v
	case float64:
		return -v
	case int:
		return -v
	case int8:
		return -v
	case int16:
		return -v
	case int32:
		return -v
	case int64:
		return -v
	case uint:
		return -v
	case uint8:
		return -v
	case uint16:
		return -v
	case uint32:
		return -v
	case uint64:
		return -v
	case bool:
		return !v
	default:
		return -coerce.ToInt(i)
	}
}

// ShiftLeft returns a << b as int64.
func ShiftLeft(a, b any) int64 {
	return coerce.ToInt64(a) << uint(coerce.ToInt64(b))
}

// ShiftRight returns a >> b as int64.
func ShiftRight(a, b any) int64 {
	return coerce.ToInt64(a) >> uint(coerce.ToInt64(b))
}

// ShiftRightUnsigned returns a >>> b as int64.
func ShiftRightUnsigned(a, b any) int64 {
	return int64(uint64(coerce.ToInt64(a)) >> uint(coerce.ToInt64(b)))
}

// Size returns the length of x.
func Size(x any) int64 {
	if x == nil {
		return 0
	}
	if s, ok := x.(util.Set); ok {
		return int64(s.Size())
	}
	r := reflect.ValueOf(x)
	switch r.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return int64(r.Len())
	case reflect.Ptr:
		elem := r.Elem()
		if elem.IsValid() {
			return Size(elem.Interface())
		}
		return 0
	default:
		panic(fmt.Sprintf("invalid argument for size (type %T)", x))
	}
}

// Slice returns array[from:to] with negative-index support.
func Slice(array, from, to any) any {
	v := reflect.ValueOf(array)
	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		length := v.Len()
		a, b := coerce.ToInt(from), coerce.ToInt(to)
		if a < 0 {
			a = length + a
		}
		if a < 0 {
			a = 0
		}
		if b < 0 {
			b = length + b
		}
		if b < 0 {
			b = 0
		}
		if b > length {
			b = length
		}
		if a > b {
			a = b
		}
		if v.Kind() == reflect.Array && !v.CanAddr() {
			newValue := reflect.New(v.Type()).Elem()
			newValue.Set(v)
			v = newValue
		}
		value := v.Slice(a, b)
		if value.IsValid() {
			return value.Interface()
		}
	case reflect.Ptr:
		value := v.Elem()
		if value.IsValid() {
			return Slice(value.Interface(), from, to)
		}
	}
	panic(fmt.Sprintf("cannot slice %v", from))
}

// StrictEqual reports whether a equals b in both
// type and value.
func StrictEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}
	return reflect.DeepEqual(a, b)
}

// ToRange returns an ascending or descending int
// slice from min to max inclusive.
func ToRange(min, max int) []int {
	if min <= max {
		size := max - min + 1
		rng := make([]int, size)
		for i := range rng {
			rng[i] = min + i
		}
		return rng
	}
	size := min - max + 1
	rng := make([]int, size)
	for i := range rng {
		rng[i] = min - i
	}
	return rng
}
