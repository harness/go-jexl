// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package coerce

import (
	"math"
	"math/big"
	"strconv"

	"github.com/harness/go-jexl/jexl/internal/decimal"
)

// ToBigInt returns the *big.Int value if v is one.
func ToBigInt(v any) (*big.Int, bool) {
	p, ok := v.(*big.Int)
	return p, ok
}

// ToBigDecimal returns the decimal.Decimal value if v is one.
func ToBigDecimal(v any) (decimal.Decimal, bool) {
	d, ok := v.(decimal.Decimal)
	return d, ok
}

// ToDecimal converts any numeric value to a decimal.Decimal.
func ToDecimal(v any) decimal.Decimal {
	if bd, ok := ToBigDecimal(v); ok {
		return bd
	}
	if bi, ok := ToBigInt(v); ok {
		return decimal.NewFromBigInt(bi, 0)
	}
	f := ToFloat64(v)
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return decimal.Zero
	}
	return decimal.NewFromFloat(f)
}

// ToNumeric converts bool and string values to numeric types
// for arithmetic. Already-numeric values are returned unchanged.
func ToNumeric(v any) any {
	switch x := v.(type) {
	case bool:
		if x {
			return 1
		}
		return 0
	case string:
		if i, err := strconv.ParseInt(x, 10, 64); err == nil {
			return int(i)
		}
		if f, err := strconv.ParseFloat(x, 64); err == nil {
			return f
		}
	}
	return v
}
