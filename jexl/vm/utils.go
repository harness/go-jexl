// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package vm

import (
	"math"
	"math/big"
	"reflect"
	"strconv"

	"github.com/harness/go-jexl/jexl/classes/java/util"
	"github.com/harness/go-jexl/jexl/internal/decimal"
)

// coerceToNumeric converts booleans and string-numbers to numeric types for arithmetic.
func coerceToNumeric(v any) any {
	switch x := v.(type) {
	case bool:
		if x {
			return 1
		}
		return 0
	case string:
		// Try integer parse (exact, no decimals).
		if i, err := strconv.ParseInt(x, 10, 64); err == nil {
			return int(i)
		}
		// Try float parse (handles "3.14", "1.5e2", etc.).
		if f, err := strconv.ParseFloat(x, 64); err == nil {
			return f
		}
		return v
	}
	return v
}

// asBigInt returns *big.Int if v is one.
func asBigInt(v any) (*big.Int, bool) {
	p, ok := v.(*big.Int)
	return p, ok
}

// asBigDec returns decimal.Decimal if v is one.
func asBigDec(v any) (decimal.Decimal, bool) {
	d, ok := v.(decimal.Decimal)
	return d, ok
}

// isIntegral returns true if v is an integer type (not float).
func isIntegral(v any) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	}
	return false
}

// toDecimalSafe converts any numeric value to a shopspring decimal.
func toDecimalSafe(v any) decimal.Decimal {
	if bd, ok := asBigDec(v); ok {
		return bd
	}
	if bi, ok := asBigInt(v); ok {
		return decimal.NewFromBigInt(bi, 0)
	}
	f := toFloat64Safe(v)
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return decimal.Zero
	}
	return decimal.NewFromFloat(f)
}

// bigCompare compares two values where at least one is BigInteger or BigDecimal.
// Returns (cmp, true) where cmp < 0 means a<b, 0 means a==b, >0 means a>b.
func bigCompare(a, b any) (int, bool) {
	_, aIsBigInt := asBigInt(a)
	_, bIsBigInt := asBigInt(b)
	_, aIsBigDec := asBigDec(a)
	_, bIsBigDec := asBigDec(b)

	if !aIsBigInt && !bIsBigInt && !aIsBigDec && !bIsBigDec {
		return 0, false
	}

	return toDecimalSafe(a).Cmp(toDecimalSafe(b)), true
}

// toFloat64Safe converts a value to float64, handling booleans and strings.
func toInt64Safe(v any) (int64, bool) {
	switch x := v.(type) {
	case int:
		return int64(x), true
	case int8:
		return int64(x), true
	case int16:
		return int64(x), true
	case int32:
		return int64(x), true
	case int64:
		return x, true
	}
	return 0, false
}

func toFloat64Safe(v any) float64 {
	return ToFloat64(v)
}

// bigArith handles *big.Int/decimal.Decimal arithmetic.
func bigArith(a, b any, op string) (any, bool) {
	ai, aIsBigInt := asBigInt(a)
	bi, bIsBigInt := asBigInt(b)
	_, aIsBigDec := asBigDec(a)
	_, bIsBigDec := asBigDec(b)

	if !aIsBigInt && !bIsBigInt && !aIsBigDec && !bIsBigDec {
		return nil, false
	}

	// If either side is BigDecimal, use decimal arithmetic (preserves scale).
	if aIsBigDec || bIsBigDec {
		da := toDecimalSafe(a)
		db := toDecimalSafe(b)
		var result decimal.Decimal
		switch op {
		case "+":
			result = da.Add(db)
		case "-":
			result = da.Sub(db)
		case "*":
			result = da.Mul(db)
		case "/":
			if db.IsZero() {
				return decimal.Zero, true
			}
			result = da.Div(db)
		case "%":
			if db.IsZero() {
				return decimal.Zero, true
			}
			result = da.Mod(db)
		default:
			return nil, false
		}
		return result, true
	}

	// Both are BigInteger (or one is a plain int/bool/string).
	var ia, ib *big.Int
	if aIsBigInt {
		ia = ai
	} else {
		ia = big.NewInt(int64(toFloat64Safe(a)))
	}
	if bIsBigInt {
		ib = bi
	} else {
		ib = big.NewInt(int64(toFloat64Safe(b)))
	}
	var result *big.Int
	switch op {
	case "+":
		result = new(big.Int).Add(ia, ib)
	case "-":
		result = new(big.Int).Sub(ia, ib)
	case "*":
		result = new(big.Int).Mul(ia, ib)
	case "/":
		result = new(big.Int).Quo(ia, ib)
	case "%":
		result = new(big.Int).Rem(ia, ib)
	default:
		return nil, false
	}
	return result, true
}

// toIterableSlice converts builtin collections and maps to []any for iteration.
func toIterableSlice(v any) any {
	if v == nil {
		return []any{}
	}
	// java.util.List (e.g. ArrayList)
	if list, ok := v.(util.List); ok {
		n := list.Size()
		items := make([]any, n)
		for i := 0; i < n; i++ {
			val, err := list.Get(i)
			if err != nil {
				panic(err)
			}
			items[i] = val
		}
		return items
	}
	if s, ok := v.(util.Set); ok {
		return s.ToArray()
	}
	// Convert map types to a slice of values
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Map {
		items := make([]any, 0, rv.Len())
		for _, key := range rv.MapKeys() {
			items = append(items, rv.MapIndex(key).Interface())
		}
		return items
	}
	return v
}

type (
	Function     = func(params ...any) (any, error)
	SafeFunction = func(params ...any) (any, uint, error)
)

var (
	errorType = reflect.TypeOf((*error)(nil)).Elem()
)

type Scope struct {
	Array reflect.Value
	Index int
	Len   int
	Count int
	Acc   any
	// VarSlot is the Variables index for the foreach loop variable (-1 = not a foreach scope).
	VarSlot int
	// Fast paths
	Ints    []int
	Floats  []float64
	Strings []string
	Anys    []any
}

// Item returns the current element from the scope using fast paths when available.
func (s *Scope) Item() any {
	if s.Ints != nil {
		return s.Ints[s.Index]
	}
	if s.Floats != nil {
		return s.Floats[s.Index]
	}
	if s.Strings != nil {
		return s.Strings[s.Index]
	}
	if s.Anys != nil {
		return s.Anys[s.Index]
	}
	return s.Array.Index(s.Index).Interface()
}
