// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package math

import (
	"fmt"
	"math/big"

	"github.com/harness/go-jexl/jexl/coerce"
)

// BigIntegerClass is the java.math.BigInteger class object.
var BigIntegerClass bigIntegerClass

type bigIntegerClass struct{}

func (bigIntegerClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "new":
		bi := new(big.Int)
		if len(args) == 0 {
			return &BigInteger{V: bi}, nil
		}
		switch v := args[0].(type) {
		case string:
			base := 10
			if len(args) >= 2 {
				base = coerce.ToInt(args[1])
			}
			if _, ok := bi.SetString(v, base); !ok {
				return nil, fmt.Errorf("BigInteger.new: invalid value %q in base %d", v, base)
			}
			return &BigInteger{V: bi}, nil
		default:
			bi.SetInt64(coerce.ToInt64(v))
			return &BigInteger{V: bi}, nil
		}
	case "ZERO":
		return &BigInteger{V: new(big.Int)}, nil
	case "ONE":
		return &BigInteger{V: big.NewInt(1)}, nil
	case "TWO":
		return &BigInteger{V: big.NewInt(2)}, nil
	case "TEN":
		return &BigInteger{V: big.NewInt(10)}, nil
	case "valueOf":
		if len(args) != 1 {
			return nil, fmt.Errorf("BigInteger.valueOf: expected 1 argument")
		}
		return &BigInteger{V: big.NewInt(coerce.ToInt64(args[0]))}, nil
	}
	return nil, fmt.Errorf("BigInteger.%s: undefined", method)
}

// BigInteger mirrors java.math.BigInteger.
type BigInteger struct {
	V *big.Int
}

// toBigInt converts a value to *big.Int, unwrapping *BigInteger if needed.
func toBigInt(v any) *big.Int {
	if bi, ok := v.(*BigInteger); ok {
		return bi.V
	}
	return big.NewInt(coerce.ToInt64(v))
}

// Add returns a new BigInteger equal to b + other.
func (b *BigInteger) Add(other *BigInteger) *BigInteger {
	return &BigInteger{V: new(big.Int).Add(b.V, other.V)}
}

// Subtract returns a new BigInteger equal to b - other.
func (b *BigInteger) Subtract(other *BigInteger) *BigInteger {
	return &BigInteger{V: new(big.Int).Sub(b.V, other.V)}
}

// Multiply returns a new BigInteger equal to b * other.
func (b *BigInteger) Multiply(other *BigInteger) *BigInteger {
	return &BigInteger{V: new(big.Int).Mul(b.V, other.V)}
}

// Divide returns a new BigInteger equal to b / other, or an error on zero.
func (b *BigInteger) Divide(other *BigInteger) (*BigInteger, error) {
	if other.V.Sign() == 0 {
		return nil, fmt.Errorf("BigInteger.divide: division by zero")
	}
	return &BigInteger{V: new(big.Int).Quo(b.V, other.V)}, nil
}

// Remainder returns a new BigInteger equal to b % other.
func (b *BigInteger) Remainder(other *BigInteger) *BigInteger {
	return &BigInteger{V: new(big.Int).Rem(b.V, other.V)}
}

// Mod returns a new BigInteger equal to b mod other.
func (b *BigInteger) Mod(other *BigInteger) *BigInteger {
	return &BigInteger{V: new(big.Int).Mod(b.V, other.V)}
}

// Pow returns a new BigInteger equal to b ^ exp.
func (b *BigInteger) Pow(exp int64) *BigInteger {
	return &BigInteger{V: new(big.Int).Exp(b.V, big.NewInt(exp), nil)}
}

// Negate returns a new BigInteger equal to -b.
func (b *BigInteger) Negate() *BigInteger {
	return &BigInteger{V: new(big.Int).Neg(b.V)}
}

// Abs returns a new BigInteger equal to |b|.
func (b *BigInteger) Abs() *BigInteger {
	return &BigInteger{V: new(big.Int).Abs(b.V)}
}

// And returns a new BigInteger equal to b & other.
func (b *BigInteger) And(other *BigInteger) *BigInteger {
	return &BigInteger{V: new(big.Int).And(b.V, other.V)}
}

// Or returns a new BigInteger equal to b | other.
func (b *BigInteger) Or(other *BigInteger) *BigInteger {
	return &BigInteger{V: new(big.Int).Or(b.V, other.V)}
}

// Xor returns a new BigInteger equal to b ^ other.
func (b *BigInteger) Xor(other *BigInteger) *BigInteger {
	return &BigInteger{V: new(big.Int).Xor(b.V, other.V)}
}

// Not returns a new BigInteger equal to ^b.
func (b *BigInteger) Not() *BigInteger {
	return &BigInteger{V: new(big.Int).Not(b.V)}
}

// ShiftLeft returns a new BigInteger equal to b << n.
func (b *BigInteger) ShiftLeft(n uint) *BigInteger {
	return &BigInteger{V: new(big.Int).Lsh(b.V, n)}
}

// ShiftRight returns a new BigInteger equal to b >> n.
func (b *BigInteger) ShiftRight(n uint) *BigInteger {
	return &BigInteger{V: new(big.Int).Rsh(b.V, n)}
}

// CompareTo returns the ordering of b relative to other.
func (b *BigInteger) CompareTo(other *BigInteger) int {
	return b.V.Cmp(other.V)
}

// Equals returns true if b and other are equal.
func (b *BigInteger) Equals(other *BigInteger) bool {
	return b.V.Cmp(other.V) == 0
}

// Max returns the larger of b and other.
func (b *BigInteger) Max(other *BigInteger) *BigInteger {
	if b.V.Cmp(other.V) >= 0 {
		return b
	}
	return other
}

// Min returns the smaller of b and other.
func (b *BigInteger) Min(other *BigInteger) *BigInteger {
	if b.V.Cmp(other.V) <= 0 {
		return b
	}
	return other
}

// BitLength returns the number of bits in the minimal two's-complement representation.
func (b *BigInteger) BitLength() int64 {
	return int64(b.V.BitLen())
}

// Signum returns -1, 0, or 1 as the sign of b.
func (b *BigInteger) Signum() int64 {
	return int64(b.V.Sign())
}

// IntValue returns the int64 value of b.
func (b *BigInteger) IntValue() int64 {
	return b.V.Int64()
}

// LongValue returns the int64 value of b.
func (b *BigInteger) LongValue() int64 {
	return b.V.Int64()
}

// DoubleValue returns the float64 value of b.
func (b *BigInteger) DoubleValue() float64 {
	f, _ := new(big.Float).SetInt(b.V).Float64()
	return f
}

// ToString returns the string representation of b in the given base (default 10).
func (b *BigInteger) ToString(base ...int) string {
	radix := 10
	if len(base) > 0 {
		radix = base[0]
	}
	return b.V.Text(radix)
}

// Call dispatches instance methods on BigInteger.
func (b *BigInteger) Call(method string, args ...any) (any, error) {
	switch method {
	case "toString":
		if len(args) == 0 {
			return b.ToString(), nil
		}
		return b.ToString(coerce.ToInt(args[0])), nil
	case "intValue", "hashCode":
		return b.IntValue(), nil
	case "longValue":
		return b.LongValue(), nil
	case "doubleValue":
		return b.DoubleValue(), nil
	case "add":
		return b.Add(&BigInteger{V: toBigInt(args[0])}), nil
	case "subtract":
		return b.Subtract(&BigInteger{V: toBigInt(args[0])}), nil
	case "multiply":
		return b.Multiply(&BigInteger{V: toBigInt(args[0])}), nil
	case "divide":
		return b.Divide(&BigInteger{V: toBigInt(args[0])})
	case "remainder":
		return b.Remainder(&BigInteger{V: toBigInt(args[0])}), nil
	case "mod":
		return b.Mod(&BigInteger{V: toBigInt(args[0])}), nil
	case "pow":
		return b.Pow(coerce.ToInt64(args[0])), nil
	case "negate":
		return b.Negate(), nil
	case "abs":
		return b.Abs(), nil
	case "compareTo":
		return b.CompareTo(&BigInteger{V: toBigInt(args[0])}), nil
	case "equals":
		return b.Equals(&BigInteger{V: toBigInt(args[0])}), nil
	case "max":
		return b.Max(&BigInteger{V: toBigInt(args[0])}), nil
	case "min":
		return b.Min(&BigInteger{V: toBigInt(args[0])}), nil
	case "and":
		return b.And(&BigInteger{V: toBigInt(args[0])}), nil
	case "or":
		return b.Or(&BigInteger{V: toBigInt(args[0])}), nil
	case "xor":
		return b.Xor(&BigInteger{V: toBigInt(args[0])}), nil
	case "not":
		return b.Not(), nil
	case "shiftLeft":
		return b.ShiftLeft(uint(coerce.ToInt(args[0]))), nil
	case "shiftRight":
		return b.ShiftRight(uint(coerce.ToInt(args[0]))), nil
	case "bitLength":
		return b.BitLength(), nil
	case "signum":
		return b.Signum(), nil
	}
	return nil, fmt.Errorf("BigInteger instance: undefined method %q", method)
}
