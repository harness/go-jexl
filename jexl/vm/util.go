// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/harness/go-jexl/jexl/classes/java/util"
	"github.com/harness/go-jexl/jexl/coerce"
	"github.com/harness/go-jexl/jexl/internal/decimal"
	"github.com/harness/go-jexl/jexl/internal/eval"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// opArgLenEstimation maps opcodes to estimated arg counts.
var opArgLenEstimation = [...]int{
	OpCall1: 1,
	OpCall2: 2,
	OpCall3: 3,
	// unknown but at least 4; conservative to avoid overallocation
	OpCallN: 4,
	// unknown; 3 is a reasonable guess for fast calls
	OpCallFast: 3,
}

// addNumeric adds two values, promoting to *big.Int on overflow.
func addNumeric(a, b any) any {
	a = coerce.ToNumeric(a)
	b = coerce.ToNumeric(b)
	if ai, aOk := toInt64Safe(a); aOk {
		if bi, bOk := toInt64Safe(b); bOk {
			result := ai + bi
			if (bi > 0 && result < ai) || (bi < 0 && result > ai) {
				return new(big.Int).Add(big.NewInt(ai), big.NewInt(bi))
			}
		}
	}
	return eval.Add(a, b)
}

// applyCompoundOp applies a compound assignment operator.
func applyCompoundOp(op string, lhs, rhs any) any {
	switch op {
	case "+=":
		if _, lhsStr := lhs.(string); lhsStr {
			return lhs.(string) + coerce.ToString(rhs)
		} else if _, rhsStr := rhs.(string); rhsStr {
			return coerce.ToString(lhs) + rhs.(string)
		}
		return eval.Add(lhs, rhs)
	case "-=":
		return eval.Subtract(lhs, rhs)
	case "*=":
		return eval.Multiply(lhs, rhs)
	case "/=":
		return eval.Divide(lhs, rhs)
	case "%=":
		return eval.Modulo(lhs, rhs)
	case "&=":
		return eval.BitwiseAnd(lhs, rhs)
	case "|=":
		return eval.BitwiseOr(lhs, rhs)
	case "^=":
		return eval.BitwiseXor(lhs, rhs)
	case "<<=":
		return eval.ShiftLeft(lhs, rhs)
	case ">>=":
		return eval.ShiftRight(lhs, rhs)
	case ">>>=":
		return eval.ShiftRightUnsigned(lhs, rhs)
	default:
		panic(fmt.Sprintf("unknown compound assignment operator %q", op))
	}
}

// bigArith handles *big.Int and decimal.Decimal arithmetic.
func bigArith(a, b any, op string) (any, bool) {
	ai, aIsBigInt := coerce.ToBigInt(a)
	bi, bIsBigInt := coerce.ToBigInt(b)
	_, aIsBigDec := coerce.ToBigDecimal(a)
	_, bIsBigDec := coerce.ToBigDecimal(b)

	if !aIsBigInt && !bIsBigInt && !aIsBigDec && !bIsBigDec {
		return nil, false
	}

	if aIsBigDec || bIsBigDec {
		da := coerce.ToDecimal(a)
		db := coerce.ToDecimal(b)
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

	var ia, ib *big.Int
	if aIsBigInt {
		ia = ai
	} else {
		ia = big.NewInt(int64(coerce.ToFloat64(a)))
	}
	if bIsBigInt {
		ib = bi
	} else {
		ib = big.NewInt(int64(coerce.ToFloat64(b)))
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

// clearSlice zeroes every element of s.
func clearSlice[S ~[]E, E any](s S) {
	var zero E
	for i := range s {
		s[i] = zero
	}
}

// estimateFnArgsCount estimates the number of function arguments
// a program will need at runtime.
func estimateFnArgsCount(program *Program) int {
	var count int
	for _, op := range program.Bytecode {
		if int(op) < len(opArgLenEstimation) {
			count += opArgLenEstimation[op]
		}
	}
	return count
}

// hasAntishPrefix reports whether any key in m starts with
// prefix+".". Used to decide whether to create an AntishCursor.
func hasAntishPrefix(m map[string]any, prefix string) bool {
	dotPrefix := prefix + "."
	for k := range m {
		if strings.HasPrefix(k, dotPrefix) {
			return true
		}
	}
	return false
}

// instanceOf reports whether a is an instance of the named type.
func instanceOf(a any, typeName string) bool {
	if a == nil {
		return false
	}
	switch typeName {
	case "String":
		_, ok := a.(string)
		return ok
	case "Boolean":
		_, ok := a.(bool)
		return ok
	case "Integer", "Long", "Short", "Byte":
		switch a.(type) {
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float64:
			rv := reflect.ValueOf(a)
			switch rv.Kind() {
			case reflect.Float64, reflect.Float32:
				f := rv.Float()
				return f == float64(int64(f))
			default:
				return true
			}
		}
		return false
	case "Float", "Double":
		switch a.(type) {
		case float32, float64, int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64:
			return true
		}
		return false
	case "Number":
		switch a.(type) {
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64:
			return true
		}
		return false
	case "Map":
		return reflect.TypeOf(a).Kind() == reflect.Map
	case "List", "Array":
		k := reflect.TypeOf(a).Kind()
		return k == reflect.Slice || k == reflect.Array
	default:
		return reflect.TypeOf(a).Name() == typeName
	}
}

// toInt64Safe returns the int64 value for native signed integers.
// Used for overflow detection in addNumeric.
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

// toIterableSlice coerces v to a []any for range iteration.
func toIterableSlice(v any) any {
	if v == nil {
		return []any{}
	}
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
