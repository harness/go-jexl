// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package math

import (
	"fmt"
	gomath "math"
	"math/rand"

	"github.com/harness/go-jexl/jexl/coerce"
)

// MathClass is the java.lang.Math class object.
var MathClass mathClass

type mathClass struct{}

func (mathClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "PI":
		return gomath.Pi, nil
	case "E":
		return gomath.E, nil
	case "random":
		return rand.Float64(), nil
	case "abs":
		return gomath.Abs(coerce.ToFloat64(args[0])), nil
	case "ceil":
		return gomath.Ceil(coerce.ToFloat64(args[0])), nil
	case "floor":
		return gomath.Floor(coerce.ToFloat64(args[0])), nil
	case "sqrt":
		return gomath.Sqrt(coerce.ToFloat64(args[0])), nil
	case "cbrt":
		return gomath.Cbrt(coerce.ToFloat64(args[0])), nil
	case "log":
		return gomath.Log(coerce.ToFloat64(args[0])), nil
	case "log10":
		return gomath.Log10(coerce.ToFloat64(args[0])), nil
	case "log1p":
		return gomath.Log1p(coerce.ToFloat64(args[0])), nil
	case "exp":
		return gomath.Exp(coerce.ToFloat64(args[0])), nil
	case "sin":
		return gomath.Sin(coerce.ToFloat64(args[0])), nil
	case "cos":
		return gomath.Cos(coerce.ToFloat64(args[0])), nil
	case "tan":
		return gomath.Tan(coerce.ToFloat64(args[0])), nil
	case "asin":
		return gomath.Asin(coerce.ToFloat64(args[0])), nil
	case "acos":
		return gomath.Acos(coerce.ToFloat64(args[0])), nil
	case "atan":
		return gomath.Atan(coerce.ToFloat64(args[0])), nil
	case "trunc":
		return gomath.Trunc(coerce.ToFloat64(args[0])), nil
	case "signum", "sign":
		v := coerce.ToFloat64(args[0])
		switch {
		case v < 0:
			return -1.0, nil
		case v > 0:
			return 1.0, nil
		default:
			return 0.0, nil
		}
	case "round":
		return gomath.Floor(coerce.ToFloat64(args[0]) + 0.5), nil
	case "pow":
		return gomath.Pow(coerce.ToFloat64(args[0]), coerce.ToFloat64(args[1])), nil
	case "max":
		return gomath.Max(coerce.ToFloat64(args[0]), coerce.ToFloat64(args[1])), nil
	case "min":
		return gomath.Min(coerce.ToFloat64(args[0]), coerce.ToFloat64(args[1])), nil
	case "hypot":
		return gomath.Hypot(coerce.ToFloat64(args[0]), coerce.ToFloat64(args[1])), nil
	case "atan2":
		return gomath.Atan2(coerce.ToFloat64(args[0]), coerce.ToFloat64(args[1])), nil
	}
	return nil, fmt.Errorf("Math.%s: undefined", method)
}
