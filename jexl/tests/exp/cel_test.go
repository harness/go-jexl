// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package cel_tests

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/harness/go-jexl/jexl/coerce"

	"github.com/harness/go-jexl/jexl"
	"github.com/harness/go-jexl/jexl/builtin"
)

// Test structure holds unit tests
type Test struct {
	ID          string         `json:"id"`
	Category    string         `json:"category"`
	Description string         `json:"description"`
	Context     map[string]any `json:"context"`
	Expr        string         `json:"expr"`
	Result      any            `json:"expected"`
	Error       string         `json:"error"`
	Skip        bool           `json:"skip"`
}

func TestSpec(t *testing.T) {

	data, err := os.ReadFile("./cel.json")
	if err != nil {
		t.Fatal(err)
	}

	var cases []*Test
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}

	for _, test_ := range cases {

		t.Run(test_.ID, func(t *testing.T) {

			if test_.Skip {
				t.Skip()
				return
			}

			opts := []jexl.Option{
				// base64.* namespace (CEL base64 extensions)
				jexl.WithFunctionNamespace("base64", "encode", builtin.Base64Encode),
				jexl.WithFunctionNamespace("base64", "decode", builtin.Base64Decode),

				// math.* namespace (CEL math extensions)
				jexl.WithFunctionNamespace("math", "greatest", celGreatest),
				jexl.WithFunctionNamespace("math", "least", celLeast),
				jexl.WithFunctionNamespace("math", "ceil", builtin.Ceil),
				jexl.WithFunctionNamespace("math", "floor", builtin.Floor),
				jexl.WithFunctionNamespace("math", "round", builtin.Round),
				jexl.WithFunctionNamespace("math", "trunc", celMathTrunc),
				jexl.WithFunctionNamespace("math", "abs", builtin.Abs),
				jexl.WithFunctionNamespace("math", "sign", celMathSign),
				jexl.WithFunctionNamespace("math", "isNaN", builtin.IsNaN),
				jexl.WithFunctionNamespace("math", "isInf", builtin.IsInfinite),
				jexl.WithFunctionNamespace("math", "isFinite", celMathIsFinite),
				jexl.WithFunctionNamespace("math", "bitAnd", builtin.BitAnd),
				jexl.WithFunctionNamespace("math", "bitOr", builtin.BitOr),
				jexl.WithFunctionNamespace("math", "bitXor", builtin.BitXor),
				jexl.WithFunctionNamespace("math", "bitNot", builtin.BitNot),
				jexl.WithFunctionNamespace("math", "bitShiftLeft", builtin.BitShiftLeft),
				jexl.WithFunctionNamespace("math", "bitShiftRight", builtin.BitShiftRight),

				// strings.* namespace (CEL string extensions)
				jexl.WithFunctionNamespace("strings", "quote", builtin.Quote),

				// CEL type conversion and reflection functions
				jexl.WithFunction("double", celDouble),
				jexl.WithFunction("int", celInt),
				jexl.WithFunction("string", celString),
				jexl.WithFunction("uint", celUint),
				jexl.WithFunction("bool", celBool),
				jexl.WithFunction("dyn", celDyn),
				jexl.WithFunction("type", celTypeof),
			}

			// add context to options
			opts = append(opts, jexl.WithContext(test_.Context))

			program, err := jexl.Compile(test_.Expr, opts...)
			if err != nil {
				if test_.Error == "" {
					t.Error(err)
					t.Logf("name: %s", test_.ID)
					t.Logf("category: %s", test_.Category)
					t.Logf("description: %s", test_.Description)
					t.Logf("expression: %v", test_.Expr)
					return
				} else {
					return // expect error
				}
			}

			out, err := jexl.Run(program, test_.Context)
			if err != nil {
				if test_.Error == "" {
					t.Error(err)
					t.Logf("name: %s", test_.ID)
					t.Logf("category: %s", test_.Category)
					t.Logf("description: %s", test_.Description)
					t.Logf("expression: %v", test_.Expr)
					return
				} else {
					return // expect error
				}
			}

			// ensure equality
			if !coerce.DeepEqual(test_.Result, out) {
				t.Errorf("Unexpected result")
				t.Logf("name: %s", test_.ID)
				t.Logf("category: %s", test_.Category)
				t.Logf("description: %s", test_.Description)
				t.Logf("expression: %v", test_.Expr)
				t.Logf("got: %v", out)
				t.Logf("want: %v", test_.Result)
				return
			}

		})

	}
}

//
// Custom Functions
//

// celGreatest returns the maximum value from variadic args or a single list.
func celGreatest(args ...any) (any, error) {
	elems, err := celFlattenArgs("math.greatest", args)
	if err != nil {
		return nil, err
	}
	best := elems[0]
	for _, a := range elems[1:] {
		if coerce.ToFloat64(a) > coerce.ToFloat64(best) {
			best = a
		}
	}
	return best, nil
}

// celLeast returns the minimum value from variadic args or a single list.
func celLeast(args ...any) (any, error) {
	elems, err := celFlattenArgs("math.least", args)
	if err != nil {
		return nil, err
	}
	best := elems[0]
	for _, a := range elems[1:] {
		if coerce.ToFloat64(a) < coerce.ToFloat64(best) {
			best = a
		}
	}
	return best, nil
}

// celFlattenArgs unpacks a single []any arg into its elements, or validates
// that at least one arg was provided.
func celFlattenArgs(name string, args []any) ([]any, error) {
	if len(args) == 1 {
		if slice, ok := args[0].([]any); ok {
			if len(slice) == 0 {
				return nil, fmt.Errorf("%s: empty list", name)
			}
			return slice, nil
		}
		return args, nil
	}
	if len(args) == 0 {
		return nil, fmt.Errorf("%s: no arguments", name)
	}
	return args, nil
}

// celMathSign returns -1, 0, or 1 for negative, zero, positive inputs.
func celMathSign(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("math.sign: expected 1 argument")
	}
	f := coerce.ToFloat64(args[0])
	switch {
	case f > 0:
		return int64(1), nil
	case f < 0:
		return int64(-1), nil
	default:
		return int64(0), nil
	}
}

// celMathIsFinite reports whether v is a finite float64.
func celMathIsFinite(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("math.isFinite: expected 1 argument")
	}
	f := coerce.ToFloat64(args[0])
	return !math.IsNaN(f) && !math.IsInf(f, 0), nil
}

// celMathTrunc truncates a float64 toward zero.
func celMathTrunc(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("math.trunc: expected 1 argument")
	}
	return math.Trunc(coerce.ToFloat64(args[0])), nil
}

// celDouble converts a value to float64.
func celDouble(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("double: expected 1 argument")
	}
	v := args[0]
	if s, ok := v.(string); ok {
		switch s {
		case "NaN":
			return math.NaN(), nil
		case "Infinity":
			return math.Inf(1), nil
		case "-Infinity":
			return math.Inf(-1), nil
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("double: invalid string %q", s)
		}
		return f, nil
	}
	return coerce.ToFloat64(v), nil
}

// celInt converts a value to int64, truncating floats toward zero.
func celInt(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("int: expected 1 argument")
	}
	v := args[0]
	switch val := v.(type) {
	case float64:
		if math.IsNaN(val) || math.IsInf(val, 0) {
			return nil, fmt.Errorf("int: cannot convert %v to int", val)
		}
		// float64(math.MaxInt64) rounds up to 9223372036854775808, so use it as the boundary
		if val >= float64(math.MaxInt64) || val < float64(math.MinInt64) {
			return nil, fmt.Errorf("int: %v out of range", val)
		}
		return int64(val), nil
	case float32:
		return int64(val), nil
	case string:
		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("int: invalid string %q", val)
		}
		return n, nil
	default:
		return coerce.ToInt64(v), nil
	}
}

// celString converts a value to its string representation.
func celString(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string: expected 1 argument")
	}
	v := args[0]
	if v == nil {
		return "null", nil
	}
	switch val := v.(type) {
	case float64:
		if math.IsNaN(val) {
			return "NaN", nil
		}
		if math.IsInf(val, 1) {
			return "Infinity", nil
		}
		if math.IsInf(val, -1) {
			return "-Infinity", nil
		}
		s := strconv.FormatFloat(val, 'g', -1, 64)
		return s, nil
	case bool:
		if val {
			return "true", nil
		}
		return "false", nil
	default:
		return coerce.ToString(v), nil
	}
}

// celUint converts a value to uint64.
func celUint(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("uint: expected 1 argument")
	}
	v := args[0]
	switch val := v.(type) {
	case float64:
		if math.IsNaN(val) || math.IsInf(val, 0) {
			return nil, fmt.Errorf("uint: cannot convert %v", val)
		}
		return uint64(val), nil
	case float32:
		return uint64(val), nil
	case string:
		n, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("uint: invalid string %q", val)
		}
		return n, nil
	default:
		return coerce.ToUint64(v), nil
	}
}

// celBool converts a value to bool, supporting CEL's string-to-bool rules.
func celBool(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("bool: expected 1 argument")
	}
	v := args[0]
	if s, ok := v.(string); ok {
		switch strings.ToLower(s) {
		case "true", "1", "t":
			return true, nil
		case "false", "0", "f":
			return false, nil
		}
		return nil, fmt.Errorf("bool: invalid string %q", s)
	}
	return coerce.ToBool(v), nil
}

// celDyn is the CEL dyn() identity function.
func celDyn(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("dyn: expected 1 argument")
	}
	return args[0], nil
}

// celTypeof returns the CEL type name of its argument.
func celTypeof(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("type: expected 1 argument")
	}
	v := args[0]
	if v == nil {
		return "null_type", nil
	}
	switch v.(type) {
	case bool:
		return "bool", nil
	case int, int8, int16, int32, int64:
		return "int", nil
	case uint, uint8, uint16, uint32, uint64:
		return "uint", nil
	case float32, float64:
		return "double", nil
	case string:
		return "string", nil
	case []any:
		return "list", nil
	case map[string]any, map[any]any:
		return "map", nil
	}
	// type-of-a-type: if it's already a type-name string, return "type"
	if _, ok := v.(string); ok {
		return "string", nil
	}
	return "dyn", nil
}
