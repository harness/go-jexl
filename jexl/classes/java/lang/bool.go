// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lang

import (
	"fmt"

	"github.com/harness/go-jexl/jexl/coerce"
)

// BooleanClass is the java.lang.Boolean class object.
var BooleanClass booleanClass

type booleanClass struct{}

func (booleanClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "new":
		if len(args) == 0 {
			return NewBoolean(false), nil
		}
		switch v := args[0].(type) {
		case bool:
			return NewBoolean(v), nil
		default:
			return NewBooleanFrom(v), nil
		}
	case "parseBoolean", "valueOf":
		if len(args) != 1 {
			return nil, fmt.Errorf("Boolean.%s: expected 1 argument", method)
		}
		return NewBooleanFrom(args[0]), nil
	case "toString":
		if len(args) != 1 {
			return nil, fmt.Errorf("Boolean.toString: expected 1 argument")
		}
		return NewBooleanFrom(args[0]).ToString(), nil
	case "compare":
		if len(args) != 2 {
			return nil, fmt.Errorf("Boolean.compare: expected 2 arguments")
		}
		return NewBooleanFrom(args[0]).CompareTo(coerce.ToBool(args[1])), nil
	case "logicalAnd":
		if len(args) != 2 {
			return nil, fmt.Errorf("Boolean.logicalAnd: expected 2 arguments")
		}
		return bool(NewBooleanFrom(args[0])) && bool(NewBooleanFrom(args[1])), nil
	case "logicalOr":
		if len(args) != 2 {
			return nil, fmt.Errorf("Boolean.logicalOr: expected 2 arguments")
		}
		return bool(NewBooleanFrom(args[0])) || bool(NewBooleanFrom(args[1])), nil
	case "logicalXor":
		if len(args) != 2 {
			return nil, fmt.Errorf("Boolean.logicalXor: expected 2 arguments")
		}
		return bool(NewBooleanFrom(args[0])) != bool(NewBooleanFrom(args[1])), nil
	default:
		return nil, fmt.Errorf("Boolean.%s: undefined", method)
	}
}

// Boolean mirrors java.lang.Boolean.
type Boolean bool

// Boolean constants mirror java.lang.Boolean.
const (
	BooleanTrue  Boolean = true
	BooleanFalse Boolean = false
)

// NewBoolean wraps a Go bool as a Boolean.
func NewBoolean(b bool) Boolean {
	return Boolean(b)
}

// NewBooleanFrom coerces any value to a Boolean.
func NewBooleanFrom(v any) Boolean {
	return Boolean(
		coerce.ToBool(v),
	)
}

// BooleanValue returns the boolean value.
func (v Boolean) BooleanValue() bool {
	return bool(v)
}

// IntValue returns 1 if true, 0 if false.
func (v Boolean) IntValue() int32 {
	if bool(v) {
		return 1
	}
	return 0
}

// LongValue returns 1 if true, 0 if false.
func (v Boolean) LongValue() int64 {
	if bool(v) {
		return 1
	}
	return 0
}

// ShortValue returns 1 if true, 0 if false.
func (v Boolean) ShortValue() int16 {
	if bool(v) {
		return 1
	}
	return 0
}

// ByteValue returns 1 if true, 0 if false.
func (v Boolean) ByteValue() int8 {
	if bool(v) {
		return 1
	}
	return 0
}

// FloatValue returns 1.0 if true, 0.0 if false.
func (v Boolean) FloatValue() float32 {
	if bool(v) {
		return 1
	}
	return 0
}

// DoubleValue returns 1.0 if true, 0.0 if false.
func (v Boolean) DoubleValue() float64 {
	if bool(v) {
		return 1
	}
	return 0
}

// CompareTo returns the ordering of the booleans.
func (v Boolean) CompareTo(other bool) int {
	switch {
	case bool(v) == other:
		return 0
	case bool(v):
		return 1
	default:
		return -1
	}
}

// Equals returns true if the booleans are equal.
func (v Boolean) Equals(other bool) bool {
	return bool(v) == other
}

// ToString returns the string representation.
func (v Boolean) ToString() string {
	if v {
		return "true"
	}
	return "false"
}

// Call dispatches instance methods.
func (b Boolean) Call(method string, args ...any) (any, error) {
	switch method {
	case "booleanValue", "toBoolean":
		return b.BooleanValue(), nil
	case "intValue", "toInteger":
		return b.IntValue(), nil
	case "longValue", "toLong":
		return b.LongValue(), nil
	case "shortValue", "toShort":
		return b.ShortValue(), nil
	case "byteValue", "toByte":
		return b.ByteValue(), nil
	case "floatValue", "toFloat":
		return b.FloatValue(), nil
	case "doubleValue", "toDouble":
		return b.DoubleValue(), nil
	case "toString":
		return b.ToString(), nil
	case "compareTo":
		if len(args) != 1 {
			return nil, fmt.Errorf("compareTo: expected 1 argument")
		}
		return b.CompareTo(NewBooleanFrom(args[0]).BooleanValue()), nil
	case "equals":
		if len(args) != 1 {
			return nil, fmt.Errorf("equals: expected 1 argument")
		}
		return b.Equals(NewBooleanFrom(args[0]).BooleanValue()), nil
	case "default":
		return b, nil
	default:
		return nil, fmt.Errorf("Boolean instance: undefined method %q", method)
	}
}
