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

// LongClass is the java.lang.Long class object.
var LongClass longClass

type longClass struct{}

func (longClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "new":
		if len(args) == 0 {
			return NewLong(0), nil
		}
		return NewLongFrom(args[0]), nil
	case "valueOf":
		if len(args) != 1 {
			return nil, fmt.Errorf("Long.valueOf: expected 1 argument")
		}
		return NewLongFrom(args[0]), nil
	case "parseLong":
		switch len(args) {
		case 1:
			return NewLongFrom(args[0]), nil
		case 2:
			n, err := strconv.ParseInt(coerce.ToString(args[0]), int(coerce.ToInt64(args[1])), 64)
			if err != nil {
				return nil, fmt.Errorf("Long.parseLong: %w", err)
			}
			return Long(n), nil
		default:
			return nil, fmt.Errorf("Long.parseLong: expected 1 or 2 arguments")
		}
	case "toString":
		switch len(args) {
		case 1:
			return NewLongFrom(args[0]).ToString(), nil
		case 2:
			base := int(coerce.ToInt64(args[1]))
			return strconv.FormatInt(coerce.ToInt64(args[0]), base), nil
		default:
			return nil, fmt.Errorf("Long.toString: expected 1 or 2 arguments")
		}
	case "compare":
		if len(args) != 2 {
			return nil, fmt.Errorf("Long.compare: expected 2 arguments")
		}
		return NewLongFrom(args[0]).CompareTo(coerce.ToInt64(args[1])), nil
	case "max":
		if len(args) != 2 {
			return nil, fmt.Errorf("Long.max: expected 2 arguments")
		}
		a, b := coerce.ToInt64(args[0]), coerce.ToInt64(args[1])
		if a >= b {
			return NewLongFrom(a), nil
		}
		return NewLongFrom(b), nil
	case "min":
		if len(args) != 2 {
			return nil, fmt.Errorf("Long.min: expected 2 arguments")
		}
		a, b := coerce.ToInt64(args[0]), coerce.ToInt64(args[1])
		if a <= b {
			return NewLongFrom(a), nil
		}
		return NewLongFrom(b), nil
	case "toBinaryString":
		if len(args) != 1 {
			return nil, fmt.Errorf("Long.toBinaryString: expected 1 argument")
		}
		return strconv.FormatInt(coerce.ToInt64(args[0]), 2), nil
	case "toHexString":
		if len(args) != 1 {
			return nil, fmt.Errorf("Long.toHexString: expected 1 argument")
		}
		return strconv.FormatInt(coerce.ToInt64(args[0]), 16), nil
	case "MAX_VALUE":
		return Long(math.MaxInt64), nil
	case "MIN_VALUE":
		return Long(math.MinInt64), nil
	case "SIZE":
		return int64(64), nil
	case "BYTES":
		return int64(8), nil
	}
	return nil, fmt.Errorf("Long.%s: undefined", method)
}

// Long mirrors java.lang.Long.
type Long int64

// NewLong wraps a Go int64 as a Long.
func NewLong(v int64) Long {
	return Long(v)
}

// NewLongFrom coerces any value to a Long.
func NewLongFrom(v any) Long {
	return Long(
		coerce.ToInt64(v),
	)
}

// LongValue returns the int64 value.
func (l Long) LongValue() int64 {
	return int64(l)
}

// IntValue returns the int32 value.
func (l Long) IntValue() int32 {
	return int32(l)
}

// ShortValue returns the int16 value.
func (l Long) ShortValue() int16 {
	return int16(l)
}

// ByteValue returns the int8 value.
func (l Long) ByteValue() int8 {
	return int8(l)
}

// FloatValue returns the float32 value.
func (l Long) FloatValue() float32 {
	return float32(l)
}

// DoubleValue returns the float64 value.
func (l Long) DoubleValue() float64 {
	return float64(l)
}

// CompareTo returns -1, 0, or 1 comparing l to other.
func (l Long) CompareTo(other int64) int {
	switch {
	case int64(l) < other:
		return -1
	case int64(l) > other:
		return 1
	default:
		return 0
	}
}

// Equals returns true if the values are equal.
func (l Long) Equals(other int64) bool {
	return int64(l) == other
}

// ToString returns the string representation.
func (l Long) ToString() string {
	return strconv.FormatInt(int64(l), 10)
}

// Call dispatches instance methods.
func (l Long) Call(method string, args ...any) (any, error) {
	switch method {
	case "longValue", "toLong":
		return l.LongValue(), nil
	case "intValue", "toInteger":
		return l.IntValue(), nil
	case "shortValue", "toShort":
		return l.ShortValue(), nil
	case "byteValue", "toByte":
		return l.ByteValue(), nil
	case "floatValue", "toFloat":
		return l.FloatValue(), nil
	case "doubleValue", "toDouble":
		return l.DoubleValue(), nil
	case "booleanValue", "toBoolean":
		return l != 0, nil
	case "toString":
		return l.ToString(), nil
	case "compareTo":
		if len(args) != 1 {
			return nil, fmt.Errorf("compareTo: expected 1 argument")
		}
		return l.CompareTo(coerce.ToInt64(args[0])), nil
	case "equals":
		if len(args) != 1 {
			return nil, fmt.Errorf("equals: expected 1 argument")
		}
		return l.Equals(coerce.ToInt64(args[0])), nil
	case "default":
		return l, nil
	}
	return nil, fmt.Errorf("Long instance: undefined method %q", method)
}
