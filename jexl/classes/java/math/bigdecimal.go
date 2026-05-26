// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package math

import (
	"fmt"

	"github.com/harness/go-jexl/jexl/coerce"
	"github.com/harness/go-jexl/jexl/internal/decimal"
)

func init() {
	decimal.DivisionPrecision = 34
}

// BigDecimalClass is the java.math.BigDecimal class object.
var BigDecimalClass bigDecimalClass

type bigDecimalClass struct{}

func (bigDecimalClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "new":
		if len(args) == 0 {
			return &BigDecimal{V: decimal.Zero}, nil
		}
		switch v := args[0].(type) {
		case string:
			d, err := decimal.NewFromString(v)
			if err != nil {
				return nil, fmt.Errorf("BigDecimal.new: %w", err)
			}
			return &BigDecimal{V: d}, nil
		case *BigDecimal:
			return v, nil
		default:
			return &BigDecimal{V: decimal.NewFromFloat(coerce.ToFloat64(v))}, nil
		}
	case "ZERO":
		return &BigDecimal{V: decimal.Zero}, nil
	case "ONE":
		return &BigDecimal{V: decimal.NewFromInt(1)}, nil
	case "TEN":
		return &BigDecimal{V: decimal.NewFromInt(10)}, nil
	case "valueOf":
		if len(args) != 1 {
			return nil, fmt.Errorf("BigDecimal.valueOf: expected 1 argument")
		}
		switch v := args[0].(type) {
		case string:
			d, err := decimal.NewFromString(v)
			if err != nil {
				return nil, fmt.Errorf("BigDecimal.valueOf: %w", err)
			}
			return &BigDecimal{V: d}, nil
		default:
			return &BigDecimal{V: decimal.NewFromFloat(coerce.ToFloat64(v))}, nil
		}
	}
	return nil, fmt.Errorf("BigDecimal.%s: undefined", method)
}

// BigDecimal mirrors java.math.BigDecimal.
type BigDecimal struct {
	V decimal.Decimal
}

// toDecimal converts a value to decimal.Decimal, unwrapping *BigDecimal if needed.
func (b *BigDecimal) toDecimal(v any) decimal.Decimal {
	if bd, ok := v.(*BigDecimal); ok {
		return bd.V
	}
	return decimal.NewFromFloat(coerce.ToFloat64(v))
}

// Add returns a new BigDecimal equal to b + other.
func (b *BigDecimal) Add(other *BigDecimal) *BigDecimal {
	return &BigDecimal{V: b.V.Add(other.V)}
}

// Subtract returns a new BigDecimal equal to b - other.
func (b *BigDecimal) Subtract(other *BigDecimal) *BigDecimal {
	return &BigDecimal{V: b.V.Sub(other.V)}
}

// Multiply returns a new BigDecimal equal to b * other.
func (b *BigDecimal) Multiply(other *BigDecimal) *BigDecimal {
	return &BigDecimal{V: b.V.Mul(other.V)}
}

// Divide returns a new BigDecimal equal to b / other, or an error on zero.
func (b *BigDecimal) Divide(other *BigDecimal) (*BigDecimal, error) {
	if other.V.IsZero() {
		return nil, fmt.Errorf("BigDecimal.divide: division by zero")
	}
	return &BigDecimal{V: b.V.Div(other.V)}, nil
}

// Remainder returns a new BigDecimal equal to b % other, or an error on zero.
func (b *BigDecimal) Remainder(other *BigDecimal) (*BigDecimal, error) {
	if other.V.IsZero() {
		return nil, fmt.Errorf("BigDecimal.remainder: division by zero")
	}
	return &BigDecimal{V: b.V.Mod(other.V)}, nil
}

// Pow returns a new BigDecimal equal to b ^ exp.
func (b *BigDecimal) Pow(exp int64) *BigDecimal {
	return &BigDecimal{V: b.V.Pow(decimal.NewFromInt(exp))}
}

// Negate returns a new BigDecimal equal to -b.
func (b *BigDecimal) Negate() *BigDecimal {
	return &BigDecimal{V: b.V.Neg()}
}

// Abs returns a new BigDecimal equal to |b|.
func (b *BigDecimal) Abs() *BigDecimal {
	return &BigDecimal{V: b.V.Abs()}
}

// CompareTo returns the ordering of b relative to other.
func (b *BigDecimal) CompareTo(other *BigDecimal) int {
	return b.V.Cmp(other.V)
}

// Equals returns true if b and other are equal.
func (b *BigDecimal) Equals(other *BigDecimal) bool {
	return b.V.Equal(other.V)
}

// Max returns the larger of b and other.
func (b *BigDecimal) Max(other *BigDecimal) *BigDecimal {
	if b.V.Cmp(other.V) >= 0 {
		return b
	}
	return other
}

// Min returns the smaller of b and other.
func (b *BigDecimal) Min(other *BigDecimal) *BigDecimal {
	if b.V.Cmp(other.V) <= 0 {
		return b
	}
	return other
}

// Scale returns the scale of b (number of digits to the right of the decimal point).
func (b *BigDecimal) Scale() int64 {
	return int64(-b.V.Exponent())
}

// SetScale returns a new BigDecimal rounded to n decimal places.
func (b *BigDecimal) SetScale(n int32) *BigDecimal {
	return &BigDecimal{V: b.V.Round(n)}
}

// StripTrailingZeros returns b with trailing zeros removed.
func (b *BigDecimal) StripTrailingZeros() *BigDecimal {
	return b
}

// Signum returns -1, 0, or 1 as the sign of b.
func (b *BigDecimal) Signum() int64 {
	return int64(b.V.Sign())
}

// IntValue returns the integer part of b as int64.
func (b *BigDecimal) IntValue() int64 {
	return b.V.IntPart()
}

// LongValue returns the integer part of b as int64.
func (b *BigDecimal) LongValue() int64 {
	return b.V.IntPart()
}

// DoubleValue returns b as float64.
func (b *BigDecimal) DoubleValue() float64 {
	f, _ := b.V.Float64()
	return f
}

// FloatValue returns b as float32.
func (b *BigDecimal) FloatValue() float32 {
	f, _ := b.V.Float64()
	return float32(f)
}

// ToString returns the string representation of b.
func (b *BigDecimal) ToString() string {
	return b.V.String()
}

// ToPlainString returns the plain string representation of b.
func (b *BigDecimal) ToPlainString() string {
	return b.V.String()
}

// Call dispatches instance methods on BigDecimal.
func (b *BigDecimal) Call(method string, args ...any) (any, error) {
	switch method {
	case "toString":
		return b.ToString(), nil
	case "toPlainString":
		return b.ToPlainString(), nil
	case "intValue", "hashCode":
		return b.IntValue(), nil
	case "longValue":
		return b.LongValue(), nil
	case "doubleValue":
		return b.DoubleValue(), nil
	case "floatValue":
		return b.FloatValue(), nil
	case "add":
		other := &BigDecimal{V: b.toDecimal(args[0])}
		return b.Add(other), nil
	case "subtract":
		other := &BigDecimal{V: b.toDecimal(args[0])}
		return b.Subtract(other), nil
	case "multiply":
		other := &BigDecimal{V: b.toDecimal(args[0])}
		return b.Multiply(other), nil
	case "divide":
		other := &BigDecimal{V: b.toDecimal(args[0])}
		return b.Divide(other)
	case "remainder":
		other := &BigDecimal{V: b.toDecimal(args[0])}
		return b.Remainder(other)
	case "pow":
		return b.Pow(coerce.ToInt64(args[0])), nil
	case "negate":
		return b.Negate(), nil
	case "abs":
		return b.Abs(), nil
	case "compareTo":
		other := &BigDecimal{V: b.toDecimal(args[0])}
		return b.CompareTo(other), nil
	case "equals":
		other := &BigDecimal{V: b.toDecimal(args[0])}
		return b.Equals(other), nil
	case "max":
		other := &BigDecimal{V: b.toDecimal(args[0])}
		return b.Max(other), nil
	case "min":
		other := &BigDecimal{V: b.toDecimal(args[0])}
		return b.Min(other), nil
	case "scale":
		return b.Scale(), nil
	case "setScale":
		return b.SetScale(coerce.ToInt32(args[0])), nil
	case "stripTrailingZeros":
		return b.StripTrailingZeros(), nil
	case "signum":
		return b.Signum(), nil
	}
	return nil, fmt.Errorf("BigDecimal instance: undefined method %q", method)
}
