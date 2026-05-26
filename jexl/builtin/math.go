// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import (
	"math"

	"github.com/harness/go-jexl/jexl/coerce"
)

// Abs returns the absolute value of v.
func Abs(v any) float64 {
	return math.Abs(coerce.ToFloat64(v))
}

// Ceil returns the least integer value >= v.
func Ceil(v any) float64 {
	return math.Ceil(coerce.ToFloat64(v))
}

// Floor returns the greatest integer value <= v.
func Floor(v any) float64 {
	return math.Floor(coerce.ToFloat64(v))
}

// IsInfinite reports whether v is positive or negative infinity.
func IsInfinite(v any) bool {
	return math.IsInf(
		coerce.ToFloat64(v), 0)
}

// IsNaN reports whether v is NaN.
func IsNaN(v any) bool {
	return math.IsNaN(
		coerce.ToFloat64(v),
	)
}

// Log returns the natural logarithm of v.
func Log(v any) float64 {
	return math.Log(coerce.ToFloat64(v))
}

// Log2 returns the binary logarithm of v.
func Log2(v any) float64 {
	return math.Log2(coerce.ToFloat64(v))
}

// Log10 returns the decimal logarithm of v.
func Log10(v any) float64 {
	return math.Log10(coerce.ToFloat64(v))
}

// Max returns the larger of a and b.
func Max(a any, b any) float64 {
	return math.Max(coerce.ToFloat64(a), coerce.ToFloat64(b))
}

// Min returns the smaller of a and b.
func Min(a any, b any) float64 {
	return math.Min(coerce.ToFloat64(a), coerce.ToFloat64(b))
}

// Pow returns base raised to the power of exp.
func Pow(base any, exp any) float64 {
	return math.Pow(coerce.ToFloat64(base), coerce.ToFloat64(exp))
}

// Round returns the nearest integer to v.
func Round(v any) float64 {
	return math.Round(coerce.ToFloat64(v))
}

// Sqrt returns the square root of v.
func Sqrt(v any) float64 {
	return math.Sqrt(coerce.ToFloat64(v))
}
