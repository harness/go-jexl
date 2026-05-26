// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lang

import (
	"fmt"
	"strconv"

	"github.com/harness/go-jexl/jexl/coerce"
)

// ByteClass is the java.lang.Byte class object.
var ByteClass byteClass

type byteClass struct{}

func (byteClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "new":
		if len(args) == 0 {
			return NewByte(0), nil
		}
		return NewByteFrom(args[0]), nil
	case "valueOf", "parseByte":
		if len(args) != 1 {
			return nil, fmt.Errorf("Byte.%s: expected 1 argument", method)
		}
		return NewByteFrom(args[0]), nil
	case "toString":
		if len(args) != 1 {
			return nil, fmt.Errorf("Byte.toString: expected 1 argument")
		}
		return NewByteFrom(args[0]).ToString(), nil
	case "compare":
		if len(args) != 2 {
			return nil, fmt.Errorf("Byte.compare: expected 2 arguments")
		}
		return NewByteFrom(args[0]).CompareTo(coerce.ToInt8(args[1])), nil
	case "MAX_VALUE":
		return Byte(127), nil
	case "MIN_VALUE":
		return Byte(-128), nil
	case "SIZE":
		return int64(8), nil
	case "BYTES":
		return int64(1), nil
	}
	return nil, fmt.Errorf("Byte.%s: undefined", method)
}

// Byte mirrors java.lang.Byte.
type Byte int8

// NewByte wraps a Go int8 as a Byte.
func NewByte(v int8) Byte {
	return Byte(v)
}

// NewByteFrom coerces any value to a Byte.
func NewByteFrom(v any) Byte {
	return Byte(coerce.ToInt8(v))
}

// ByteValue returns the int8 value.
func (b Byte) ByteValue() int8 {
	return int8(b)
}

// ShortValue returns the int16 value.
func (b Byte) ShortValue() int16 {
	return int16(b)
}

// IntValue returns the int32 value.
func (b Byte) IntValue() int32 {
	return int32(b)
}

// LongValue returns the int64 value.
func (b Byte) LongValue() int64 {
	return int64(b)
}

// FloatValue returns the float32 value.
func (b Byte) FloatValue() float32 {
	return float32(b)
}

// DoubleValue returns the float64 value.
func (b Byte) DoubleValue() float64 {
	return float64(b)
}

// CompareTo returns -1, 0, or 1 comparing b to other.
func (b Byte) CompareTo(other int8) int {
	switch {
	case int8(b) < other:
		return -1
	case int8(b) > other:
		return 1
	default:
		return 0
	}
}

// Equals returns true if the values are equal.
func (b Byte) Equals(other int8) bool {
	return int8(b) == other
}

// ToString returns the string representation.
func (b Byte) ToString() string {
	return strconv.FormatInt(int64(b), 10)
}

// Call dispatches instance methods.
func (b Byte) Call(method string, args ...any) (any, error) {
	switch method {
	case "byteValue", "toByte":
		return b.ByteValue(), nil
	case "shortValue", "toShort":
		return b.ShortValue(), nil
	case "intValue", "toInteger":
		return b.IntValue(), nil
	case "longValue", "toLong":
		return b.LongValue(), nil
	case "floatValue", "toFloat":
		return b.FloatValue(), nil
	case "doubleValue", "toDouble":
		return b.DoubleValue(), nil
	case "booleanValue", "toBoolean":
		return b != 0, nil
	case "toString":
		return b.ToString(), nil
	case "compareTo":
		if len(args) != 1 {
			return nil, fmt.Errorf("compareTo: expected 1 argument")
		}
		return b.CompareTo(coerce.ToInt8(args[0])), nil
	case "equals":
		if len(args) != 1 {
			return nil, fmt.Errorf("equals: expected 1 argument")
		}
		return b.Equals(coerce.ToInt8(args[0])), nil
	case "default":
		return b, nil
	}
	return nil, fmt.Errorf("Byte instance: undefined method %q", method)
}
