// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strings"

	"github.com/harness/go-jexl/jexl/classes/java/util"
	"github.com/harness/go-jexl/jexl/coerce"
	"github.com/harness/go-jexl/jexl/internal/decimal"
)

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

func Exponent(a, b any) float64 {
	return math.Pow(
		ToFloat64(a),
		ToFloat64(b),
	)
}

// IntDivide performs integer division: exact quotients return float64 (20/4 → 5.0),
// inexact quotients return truncated int64 (7/2 → 3).
func IntDivide(a, b any) any {
	ai := ToInt64(a)
	bi := ToInt64(b)
	if ai%bi == 0 {
		return float64(ai / bi)
	}
	return ai / bi
}

func BitwiseOr(a, b any) (int64, error) {
	ai, err := toInt64Err(a)
	if err != nil {
		return 0, err
	}
	bi, err := toInt64Err(b)
	if err != nil {
		return 0, err
	}
	return ai | bi, nil
}

func BitwiseXor(a, b any) (int64, error) {
	ai, err := toInt64Err(a)
	if err != nil {
		return 0, err
	}
	bi, err := toInt64Err(b)
	if err != nil {
		return 0, err
	}
	return ai ^ bi, nil
}

func BitwiseAnd(a, b any) (int64, error) {
	ai, err := toInt64Err(a)
	if err != nil {
		return 0, err
	}
	bi, err := toInt64Err(b)
	if err != nil {
		return 0, err
	}
	return ai & bi, nil
}

func BitwiseNot(a any) (int64, error) {
	ai, err := toInt64Err(a)
	if err != nil {
		return 0, err
	}
	return ^ai, nil
}

func ShiftLeft(a, b any) (int64, error) {
	ai, err := toInt64Err(a)
	if err != nil {
		return 0, err
	}
	bi, err := toInt64Err(b)
	if err != nil {
		return 0, err
	}
	return ai << uint(bi), nil
}

func ShiftRight(a, b any) (int64, error) {
	ai, err := toInt64Err(a)
	if err != nil {
		return 0, err
	}
	bi, err := toInt64Err(b)
	if err != nil {
		return 0, err
	}
	return ai >> uint(bi), nil
}

func ShiftRightUnsigned(a, b any) (int64, error) {
	ai, err := toInt64Err(a)
	if err != nil {
		return 0, err
	}
	bi, err := toInt64Err(b)
	if err != nil {
		return 0, err
	}
	return int64(uint64(ai) >> uint(bi)), nil
}

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

func toInt64Err(a any) (int64, error) {
	switch x := a.(type) {
	case int:
		return int64(x), nil
	case int8:
		return int64(x), nil
	case int16:
		return int64(x), nil
	case int32:
		return int64(x), nil
	case int64:
		return x, nil
	case uint:
		return int64(x), nil
	case uint8:
		return int64(x), nil
	case uint16:
		return int64(x), nil
	case uint32:
		return int64(x), nil
	case uint64:
		return int64(x), nil
	case float32:
		if x != float32(int64(x)) {
			return 0, fmt.Errorf("cannot use float value %v as integer in bitwise operation", x)
		}
		return int64(x), nil
	case float64:
		if math.IsNaN(x) {
			return 0, nil
		}
		if x != float64(int64(x)) {
			return 0, fmt.Errorf("cannot use float value %v as integer in bitwise operation", x)
		}
		return int64(x), nil
	case bool:
		// JEXL coerces bool to long: true=1, false=0
		if x {
			return 1, nil
		}
		return 0, nil
	case string:
		// JEXL coerces strings to long via parsing
		var i int64
		_, err := fmt.Sscanf(x, "%d", &i)
		if err == nil {
			return i, nil
		}
		// Try as float then truncate
		var f float64
		_, err = fmt.Sscanf(x, "%g", &f)
		if err == nil {
			return int64(f), nil
		}
		return 0, fmt.Errorf("cannot convert string %q to integer in bitwise operation", x)
	default:
		return 0, fmt.Errorf("invalid operation: bitwise op on %T", a)
	}
}

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

func ToRange(min, max int) []int {
	if min <= max {
		size := max - min + 1
		rng := make([]int, size)
		for i := range rng {
			rng[i] = min + i
		}
		return rng
	}
	// descending range: min > max
	size := min - max + 1
	rng := make([]int, size)
	for i := range rng {
		rng[i] = min - i
	}
	return rng
}

func ToInt(a any) int {
	return coerce.ToInt(a)
}

func ToInt64(a any) int64 {
	return coerce.ToInt64(a)
}

func ToFloat64(a any) float64 {
	return coerce.ToFloat64(a)
}

// IsFalsy returns true when v is falsy: nil, false, zero number, or empty string.
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
		// JEXL toBoolean: empty string or "false" is falsy
		return x == "" || strings.EqualFold(x, "false")
	default:
		return false
	}
}

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
