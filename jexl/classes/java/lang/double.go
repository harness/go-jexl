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

// DoubleClass is the java.lang.Double class object.
var DoubleClass doubleClass

type doubleClass struct{}

func (doubleClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "new":
		if len(args) == 0 {
			return NewDouble(0), nil
		}
		return NewDoubleFrom(args[0]), nil
	case "valueOf", "parseDouble":
		if len(args) != 1 {
			return nil, fmt.Errorf("Double.%s: expected 1 argument", method)
		}
		return NewDoubleFrom(args[0]), nil
	case "toString":
		if len(args) != 1 {
			return nil, fmt.Errorf("Double.toString: expected 1 argument")
		}
		return NewDoubleFrom(args[0]).ToString(), nil
	case "compare":
		if len(args) != 2 {
			return nil, fmt.Errorf("Double.compare: expected 2 arguments")
		}
		return NewDoubleFrom(args[0]).CompareTo(coerce.ToFloat64(args[1])), nil
	case "max":
		if len(args) != 2 {
			return nil, fmt.Errorf("Double.max: expected 2 arguments")
		}
		return Double(math.Max(coerce.ToFloat64(args[0]), coerce.ToFloat64(args[1]))), nil
	case "min":
		if len(args) != 2 {
			return nil, fmt.Errorf("Double.min: expected 2 arguments")
		}
		return Double(math.Min(coerce.ToFloat64(args[0]), coerce.ToFloat64(args[1]))), nil
	case "sum":
		if len(args) != 2 {
			return nil, fmt.Errorf("Double.sum: expected 2 arguments")
		}
		return Double(coerce.ToFloat64(args[0]) + coerce.ToFloat64(args[1])), nil
	case "isNaN":
		if len(args) != 1 {
			return nil, fmt.Errorf("Double.isNaN: expected 1 argument")
		}
		return NewDoubleFrom(args[0]).IsNaN(), nil
	case "isInfinite":
		if len(args) != 1 {
			return nil, fmt.Errorf("Double.isInfinite: expected 1 argument")
		}
		return NewDoubleFrom(args[0]).IsInfinite(), nil
	case "isFinite":
		if len(args) != 1 {
			return nil, fmt.Errorf("Double.isFinite: expected 1 argument")
		}
		return NewDoubleFrom(args[0]).IsFinite(), nil
	case "MAX_VALUE":
		return Double(math.MaxFloat64), nil
	case "MIN_VALUE":
		return Double(math.SmallestNonzeroFloat64), nil
	case "POSITIVE_INFINITY":
		return Double(math.Inf(1)), nil
	case "NEGATIVE_INFINITY":
		return Double(math.Inf(-1)), nil
	case "NaN":
		return Double(math.NaN()), nil
	case "MAX_EXPONENT":
		return int64(1023), nil
	case "MIN_EXPONENT":
		return int64(-1022), nil
	case "SIZE":
		return int64(64), nil
	case "BYTES":
		return int64(8), nil
	}
	return nil, fmt.Errorf("Double.%s: undefined", method)
}

// Double mirrors java.lang.Double.
type Double float64

// NewDouble wraps a Go float64 as a Double.
func NewDouble(v float64) Double {
	return Double(v)
}

// NewDoubleFrom coerces any value to a Double.
func NewDoubleFrom(v any) Double {
	return Double(coerce.ToFloat64(v))
}

// DoubleValue returns the float64 value.
func (d Double) DoubleValue() float64 {
	return float64(d)
}

// IntValue returns the value as int64.
func (d Double) IntValue() int64 {
	return int64(d)
}

// LongValue returns the value as int64.
func (d Double) LongValue() int64 {
	return int64(d)
}

// FloatValue returns the value as float32.
func (d Double) FloatValue() float32 {
	return float32(d)
}

// ShortValue returns the int16 value.
func (d Double) ShortValue() int16 {
	return int16(d)
}

// ByteValue returns the int8 value.
func (d Double) ByteValue() int8 {
	return int8(d)
}

// IsNaN reports whether the value is NaN.
func (d Double) IsNaN() bool {
	return math.IsNaN(float64(d))
}

// IsInfinite reports whether the value is infinite.
func (d Double) IsInfinite() bool {
	return math.IsInf(float64(d), 0)
}

// IsFinite reports whether the value is finite.
func (d Double) IsFinite() bool {
	return !d.IsInfinite() && !d.IsNaN()
}

// CompareTo returns -1, 0, or 1 comparing d to other.
func (d Double) CompareTo(other float64) int {
	switch {
	case float64(d) < other:
		return -1
	case float64(d) > other:
		return 1
	default:
		return 0
	}
}

// Equals reports whether d equals other.
func (d Double) Equals(other float64) bool {
	return float64(d) == other
}

// ToString returns the string representation.
func (d Double) ToString() string {
	return strconv.FormatFloat(float64(d), 'f', -1, 64)
}

// Call dispatches instance methods.
func (d Double) Call(method string, args ...any) (any, error) {
	switch method {
	case "doubleValue", "toDouble":
		return d.DoubleValue(), nil
	case "intValue", "toInteger":
		return d.IntValue(), nil
	case "longValue", "toLong":
		return d.LongValue(), nil
	case "floatValue", "toFloat":
		return d.FloatValue(), nil
	case "shortValue", "toShort":
		return d.ShortValue(), nil
	case "byteValue", "toByte":
		return d.ByteValue(), nil
	case "booleanValue", "toBoolean":
		return d != 0, nil
	case "isNaN":
		return d.IsNaN(), nil
	case "isInfinite":
		return d.IsInfinite(), nil
	case "isFinite":
		return d.IsFinite(), nil
	case "toString":
		return d.ToString(), nil
	case "compareTo":
		if len(args) != 1 {
			return nil, fmt.Errorf("compareTo: expected 1 argument")
		}
		return d.CompareTo(coerce.ToFloat64(args[0])), nil
	case "equals":
		if len(args) != 1 {
			return nil, fmt.Errorf("equals: expected 1 argument")
		}
		return d.Equals(coerce.ToFloat64(args[0])), nil
	case "default":
		return d, nil
	}
	return nil, fmt.Errorf("Double instance: undefined method %q", method)
}
