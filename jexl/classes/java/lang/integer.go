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

// IntegerClass is the java.lang.Integer class object.
var IntegerClass integerClass

type integerClass struct{}

func (integerClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "new":
		if len(args) == 0 {
			return NewInteger(0), nil
		}
		return NewIntegerFrom(args[0]), nil
	case "valueOf":
		if len(args) != 1 {
			return nil, fmt.Errorf("Integer.valueOf: expected 1 argument")
		}
		return NewIntegerFrom(args[0]), nil
	case "parseInt":
		switch len(args) {
		case 1:
			return NewIntegerFrom(args[0]), nil
		case 2:
			n, err := strconv.ParseInt(coerce.ToString(args[0]), int(coerce.ToInt64(args[1])), 64)
			if err != nil {
				return nil, fmt.Errorf("Integer.parseInt: %w", err)
			}
			return Integer(int32(n)), nil
		default:
			return nil, fmt.Errorf("Integer.parseInt: expected 1 or 2 arguments")
		}
	case "toString":
		switch len(args) {
		case 1:
			return NewIntegerFrom(args[0]).ToString(), nil
		case 2:
			base := int(coerce.ToInt64(args[1]))
			return strconv.FormatInt(coerce.ToInt64(args[0]), base), nil
		default:
			return nil, fmt.Errorf("Integer.toString: expected 1 or 2 arguments")
		}
	case "compare":
		if len(args) != 2 {
			return nil, fmt.Errorf("Integer.compare: expected 2 arguments")
		}
		return NewIntegerFrom(args[0]).CompareTo(coerce.ToInt32(args[1])), nil
	case "max":
		if len(args) != 2 {
			return nil, fmt.Errorf("Integer.max: expected 2 arguments")
		}
		a := coerce.ToInt32(args[0])
		b := coerce.ToInt32(args[1])
		if a >= b {
			return NewIntegerFrom(a), nil
		}
		return NewIntegerFrom(b), nil
	case "min":
		if len(args) != 2 {
			return nil, fmt.Errorf("Integer.min: expected 2 arguments")
		}
		a := coerce.ToInt32(args[0])
		b := coerce.ToInt32(args[1])
		if a <= b {
			return NewIntegerFrom(a), nil
		}
		return NewIntegerFrom(b), nil
	case "sum":
		if len(args) != 2 {
			return nil, fmt.Errorf("Integer.sum: expected 2 arguments")
		}
		return NewIntegerFrom(coerce.ToInt64(args[0]) + coerce.ToInt64(args[1])), nil
	case "toBinaryString":
		if len(args) != 1 {
			return nil, fmt.Errorf("Integer.toBinaryString: expected 1 argument")
		}
		return strconv.FormatInt(coerce.ToInt64(args[0]), 2), nil
	case "toHexString":
		if len(args) != 1 {
			return nil, fmt.Errorf("Integer.toHexString: expected 1 argument")
		}
		return strconv.FormatInt(coerce.ToInt64(args[0]), 16), nil
	case "toOctalString":
		if len(args) != 1 {
			return nil, fmt.Errorf("Integer.toOctalString: expected 1 argument")
		}
		return strconv.FormatInt(coerce.ToInt64(args[0]), 8), nil
	case "MAX_VALUE":
		return Integer(math.MaxInt32), nil
	case "MIN_VALUE":
		return Integer(math.MinInt32), nil
	case "SIZE":
		return int64(32), nil
	case "BYTES":
		return int64(4), nil
	}
	return nil, fmt.Errorf("Integer.%s: undefined", method)
}

// Integer mirrors java.lang.Integer.
type Integer int32

// NewInteger wraps a Go int32 as an Integer.
func NewInteger(v int32) Integer {
	return Integer(v)
}

// NewIntegerFrom coerces any value to an Integer.
func NewIntegerFrom(v any) Integer {
	return Integer(coerce.ToInt32(v))
}

// IntValue returns the int32 value.
func (i Integer) IntValue() int32 {
	return int32(i)
}

// LongValue returns the int64 value.
func (i Integer) LongValue() int64 {
	return int64(i)
}

// ShortValue returns the int16 value.
func (i Integer) ShortValue() int16 {
	return int16(i)
}

// ByteValue returns the int8 value.
func (i Integer) ByteValue() int8 {
	return int8(i)
}

// FloatValue returns the float32 value.
func (i Integer) FloatValue() float32 {
	return float32(i)
}

// DoubleValue returns the float64 value.
func (i Integer) DoubleValue() float64 {
	return float64(i)
}

// CompareTo returns -1, 0, or 1 comparing i to other.
func (i Integer) CompareTo(other int32) int {
	switch {
	case int32(i) < other:
		return -1
	case int32(i) > other:
		return 1
	default:
		return 0
	}
}

// Equals returns true if the values are equal.
func (i Integer) Equals(other int32) bool {
	return int32(i) == other
}

// ToString returns the string representation.
func (i Integer) ToString() string {
	return strconv.FormatInt(int64(i), 10)
}

// Call dispatches instance methods.
func (i Integer) Call(method string, args ...any) (any, error) {
	switch method {
	case "intValue", "toInteger":
		return i.IntValue(), nil
	case "longValue", "toLong":
		return i.LongValue(), nil
	case "shortValue", "toShort":
		return i.ShortValue(), nil
	case "byteValue", "toByte":
		return i.ByteValue(), nil
	case "floatValue", "toFloat":
		return i.FloatValue(), nil
	case "doubleValue", "toDouble":
		return i.DoubleValue(), nil
	case "booleanValue", "toBoolean":
		return i != 0, nil
	case "toString":
		return i.ToString(), nil
	case "compareTo":
		if len(args) != 1 {
			return nil, fmt.Errorf("compareTo: expected 1 argument")
		}
		return i.CompareTo(coerce.ToInt32(args[0])), nil
	case "equals":
		if len(args) != 1 {
			return nil, fmt.Errorf("equals: expected 1 argument")
		}
		return i.Equals(coerce.ToInt32(args[0])), nil
	case "default":
		return i, nil
	}
	return nil, fmt.Errorf("Integer instance: undefined method %q", method)
}
