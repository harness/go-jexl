// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lang

import (
	"fmt"
	"math"
	"strconv"

	"github.com/harness/go-jexl/jexl/coerce"
)

// FloatClass is the java.lang.Float class object.
var FloatClass floatClass

type floatClass struct{}

func (floatClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "new":
		if len(args) == 0 {
			return NewFloat(0), nil
		}
		return NewFloatFrom(args[0]), nil
	case "valueOf", "parseFloat":
		if len(args) != 1 {
			return nil, fmt.Errorf("Float.%s: expected 1 argument", method)
		}
		return NewFloatFrom(args[0]), nil
	case "toString":
		if len(args) != 1 {
			return nil, fmt.Errorf("Float.toString: expected 1 argument")
		}
		return NewFloatFrom(args[0]).ToString(), nil
	case "compare":
		if len(args) != 2 {
			return nil, fmt.Errorf("Float.compare: expected 2 arguments")
		}
		return NewFloatFrom(args[0]).CompareTo(coerce.ToFloat32(args[1])), nil
	case "max":
		if len(args) != 2 {
			return nil, fmt.Errorf("Float.max: expected 2 arguments")
		}
		return Float(math.Max(float64(coerce.ToFloat32(args[0])), float64(coerce.ToFloat32(args[1])))), nil
	case "min":
		if len(args) != 2 {
			return nil, fmt.Errorf("Float.min: expected 2 arguments")
		}
		return Float(math.Min(float64(coerce.ToFloat32(args[0])), float64(coerce.ToFloat32(args[1])))), nil
	case "sum":
		if len(args) != 2 {
			return nil, fmt.Errorf("Float.sum: expected 2 arguments")
		}
		return Float(coerce.ToFloat32(args[0]) + coerce.ToFloat32(args[1])), nil
	case "isNaN":
		if len(args) != 1 {
			return nil, fmt.Errorf("Float.isNaN: expected 1 argument")
		}
		return NewFloatFrom(args[0]).IsNaN(), nil
	case "isInfinite":
		if len(args) != 1 {
			return nil, fmt.Errorf("Float.isInfinite: expected 1 argument")
		}
		return NewFloatFrom(args[0]).IsInfinite(), nil
	case "isFinite":
		if len(args) != 1 {
			return nil, fmt.Errorf("Float.isFinite: expected 1 argument")
		}
		return NewFloatFrom(args[0]).IsFinite(), nil
	case "MAX_VALUE":
		return Float(math.MaxFloat32), nil
	case "MIN_VALUE":
		return Float(math.SmallestNonzeroFloat32), nil
	case "POSITIVE_INFINITY":
		return Float(float32(math.Inf(1))), nil
	case "NEGATIVE_INFINITY":
		return Float(float32(math.Inf(-1))), nil
	case "NaN":
		return Float(float32(math.NaN())), nil
	case "MAX_EXPONENT":
		return int64(127), nil
	case "MIN_EXPONENT":
		return int64(-126), nil
	case "SIZE":
		return int64(32), nil
	case "BYTES":
		return int64(4), nil
	}
	return nil, fmt.Errorf("Float.%s: undefined", method)
}

// Float mirrors java.lang.Float.
type Float float32

// NewFloat wraps a Go float32 as a Float.
func NewFloat(v float32) Float {
	return Float(v)
}

// NewFloatFrom coerces any value to a Float.
func NewFloatFrom(v any) Float {
	return Float(
		coerce.ToFloat32(v),
	)
}

// FloatValue returns the float32 value.
func (f Float) FloatValue() float32 {
	return float32(f)
}

// DoubleValue returns the float64 value.
func (f Float) DoubleValue() float64 {
	return float64(f)
}

// IntValue returns the int32 value.
func (f Float) IntValue() int32 {
	return int32(f)
}

// LongValue returns the int64 value.
func (f Float) LongValue() int64 {
	return int64(f)
}

// ShortValue returns the int16 value.
func (f Float) ShortValue() int16 {
	return int16(f)
}

// ByteValue returns the int8 value.
func (f Float) ByteValue() int8 {
	return int8(f)
}

// IsNaN reports whether the value is NaN.
func (f Float) IsNaN() bool {
	return math.IsNaN(float64(f))
}

// IsInfinite reports whether the value is infinite.
func (f Float) IsInfinite() bool {
	return math.IsInf(float64(f), 0)
}

// IsFinite reports whether the value is finite.
func (f Float) IsFinite() bool {
	return !f.IsInfinite() && !f.IsNaN()
}

// CompareTo returns -1, 0, or 1 comparing f to other.
func (f Float) CompareTo(other float32) int {
	switch {
	case float32(f) < other:
		return -1
	case float32(f) > other:
		return 1
	default:
		return 0
	}
}

// Equals reports whether f equals other.
func (f Float) Equals(other float32) bool {
	return float32(f) == other
}

// ToString returns the string representation.
func (f Float) ToString() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 32)
}

// Call dispatches instance methods.
func (f Float) Call(method string, args ...any) (any, error) {
	switch method {
	case "floatValue", "toFloat":
		return f.FloatValue(), nil
	case "doubleValue", "toDouble":
		return f.DoubleValue(), nil
	case "intValue", "toInteger":
		return f.IntValue(), nil
	case "longValue", "toLong":
		return f.LongValue(), nil
	case "shortValue", "toShort":
		return f.ShortValue(), nil
	case "byteValue", "toByte":
		return f.ByteValue(), nil
	case "booleanValue", "toBoolean":
		return f != 0, nil
	case "isNaN":
		return f.IsNaN(), nil
	case "isInfinite":
		return f.IsInfinite(), nil
	case "isFinite":
		return f.IsFinite(), nil
	case "toString":
		return f.ToString(), nil
	case "compareTo":
		if len(args) != 1 {
			return nil, fmt.Errorf("compareTo: expected 1 argument")
		}
		return f.CompareTo(coerce.ToFloat32(args[0])), nil
	case "equals":
		if len(args) != 1 {
			return nil, fmt.Errorf("equals: expected 1 argument")
		}
		return f.Equals(coerce.ToFloat32(args[0])), nil
	case "default":
		return f, nil
	}
	return nil, fmt.Errorf("Float instance: undefined method %q", method)
}
