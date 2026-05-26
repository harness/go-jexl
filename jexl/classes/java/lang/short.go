// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lang

import (
	"fmt"
	"strconv"

	"github.com/harness/go-jexl/jexl/coerce"
)

// ShortClass is the java.lang.Short class object.
var ShortClass shortClass

type shortClass struct{}

func (shortClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "new":
		if len(args) == 0 {
			return NewShort(0), nil
		}
		return NewShortFrom(args[0]), nil
	case "valueOf", "parseShort":
		if len(args) != 1 {
			return nil, fmt.Errorf("Short.%s: expected 1 argument", method)
		}
		return NewShortFrom(args[0]), nil
	case "toString":
		if len(args) != 1 {
			return nil, fmt.Errorf("Short.toString: expected 1 argument")
		}
		return NewShortFrom(args[0]).ToString(), nil
	case "compare":
		if len(args) != 2 {
			return nil, fmt.Errorf("Short.compare: expected 2 arguments")
		}
		return NewShortFrom(args[0]).CompareTo(coerce.ToInt16(args[1])), nil
	case "MAX_VALUE":
		return Short(32767), nil
	case "MIN_VALUE":
		return Short(-32768), nil
	case "SIZE":
		return int64(16), nil
	case "BYTES":
		return int64(2), nil
	}
	return nil, fmt.Errorf("Short.%s: undefined", method)
}

// Short mirrors java.lang.Short.
type Short int16

// NewShort wraps a Go int16 as a Short.
func NewShort(v int16) Short {
	return Short(v)
}

// NewShortFrom coerces any value to a Short.
func NewShortFrom(v any) Short {
	return Short(
		coerce.ToInt16(v),
	)
}

// ShortValue returns the int16 value.
func (s Short) ShortValue() int16 {
	return int16(s)
}

// ByteValue returns the int8 value.
func (s Short) ByteValue() int8 {
	return int8(s)
}

// IntValue returns the int32 value.
func (s Short) IntValue() int32 {
	return int32(s)
}

// LongValue returns the int64 value.
func (s Short) LongValue() int64 {
	return int64(s)
}

// FloatValue returns the float32 value.
func (s Short) FloatValue() float32 {
	return float32(s)
}

// DoubleValue returns the float64 value.
func (s Short) DoubleValue() float64 {
	return float64(s)
}

// CompareTo returns -1, 0, or 1 comparing s to other.
func (s Short) CompareTo(other int16) int {
	switch {
	case int16(s) < other:
		return -1
	case int16(s) > other:
		return 1
	default:
		return 0
	}
}

// Equals returns true if the values are equal.
func (s Short) Equals(other int16) bool {
	return int16(s) == other
}

// ToString returns the string representation.
func (s Short) ToString() string {
	return strconv.FormatInt(int64(s), 10)
}

// Call dispatches instance methods.
func (s Short) Call(method string, args ...any) (any, error) {
	switch method {
	case "shortValue", "toShort":
		return s.ShortValue(), nil
	case "byteValue", "toByte":
		return s.ByteValue(), nil
	case "intValue", "toInteger":
		return s.IntValue(), nil
	case "longValue", "toLong":
		return s.LongValue(), nil
	case "floatValue", "toFloat":
		return s.FloatValue(), nil
	case "doubleValue", "toDouble":
		return s.DoubleValue(), nil
	case "booleanValue", "toBoolean":
		return s != 0, nil
	case "toString":
		return s.ToString(), nil
	case "compareTo":
		if len(args) != 1 {
			return nil, fmt.Errorf("compareTo: expected 1 argument")
		}
		return s.CompareTo(coerce.ToInt16(args[0])), nil
	case "equals":
		if len(args) != 1 {
			return nil, fmt.Errorf("equals: expected 1 argument")
		}
		return s.Equals(coerce.ToInt16(args[0])), nil
	case "default":
		return s, nil
	}
	return nil, fmt.Errorf("Short instance: undefined method %q", method)
}
