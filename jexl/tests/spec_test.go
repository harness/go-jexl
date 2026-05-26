// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package tests

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
	"github.com/harness/go-jexl/jexl/classes"
	javalang "github.com/harness/go-jexl/jexl/classes/java/lang"
	javamath "github.com/harness/go-jexl/jexl/classes/java/math"
	javautil "github.com/harness/go-jexl/jexl/classes/java/util"
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

	data1, err1 := os.ReadFile("./testdata/integration.json")
	if err1 != nil {
		t.Fatal(err1)
	}

	data2, err2 := os.ReadFile("./testdata/synthetic.json")
	if err2 != nil {
		t.Fatal(err2)
	}

	var cases1 []Test
	if err := json.Unmarshal(data1, &cases1); err != nil {
		t.Fatal(err)
	}
	var cases2 []Test
	if err := json.Unmarshal(data2, &cases2); err != nil {
		t.Fatal(err)
	}
	cases := append(cases1, cases2...)

	for _, test_ := range cases {

		t.Run(test_.ID, func(t *testing.T) {

			if test_.Skip {
				t.Skip()
				return
			}

			opts := []jexl.Option{
				// Custom classes
				jexl.WithClass("Empty", emptyClass{}),
				jexl.WithClass("Point", pointClass{}),
				jexl.WithClass("org.apache.commons.jexl3.ConformanceTest$Empty", emptyClass{}),
				jexl.WithClass("org.apache.commons.jexl3.ConformanceTest$Point", pointClass{}),

				// Custom namespaces
				jexl.WithNamespace("String", StringNS),
				jexl.WithNamespace("Integer", IntegerNS),
				jexl.WithNamespace("Math", javamath.MathClass),

				// Custom functions
				jexl.WithFunctionNamespace("secrets", "getValue", func(args ...any) (any, error) {
					return "dummy-23e4567-e89b-12d3-a456-426614174000", nil
				}),

				// Custom builtin-functions (aliases to existing functions to match Harness jexl)
				jexl.WithFunctionNamespace("json", "format", builtin.JSONMarshal),
				jexl.WithFunctionNamespace("json", "stringify", builtin.JSONMarshal),
				jexl.WithFunctionNamespace("json", "object", builtin.JSONUnmarshal),
				jexl.WithFunctionNamespace("json", "select", builtin.JSONSelect),
				jexl.WithFunctionNamespace("json", "list", builtin.JSONList),
				jexl.WithFunctionNamespace("regex", "extract", builtin.RegexExtract),
				jexl.WithFunctionNamespace("regex", "replace", builtin.ReplaceAll),
				jexl.WithFunctionNamespace("regex", "match", builtin.Matches),
				jexl.WithFunctionNamespace("datetime", "currentDate", builtin.CurrentDate),
				jexl.WithFunctionNamespace("datetime", "currentTime", builtin.CurrentTime),
				jexl.WithFunctionNamespace("datetime", "format", builtin.DateFormat),
				jexl.WithFunctionNamespace("datetime", "plusMinutes", builtin.PlusMinutes),
				jexl.WithFunctionNamespace("xml", "select", builtin.XMLSelect),
				jexl.WithFunctionNamespace("base64", "encode", builtin.Base64Encode),
				jexl.WithFunctionNamespace("base64", "decode", builtin.Base64Decode),

				// Java standard library
				jexl.WithClass("java.lang.Boolean", javalang.BooleanClass),
				jexl.WithClass("java.lang.Byte", javalang.ByteClass),
				jexl.WithClass("java.lang.Short", javalang.ShortClass),
				jexl.WithClass("java.lang.Integer", javalang.IntegerClass),
				jexl.WithClass("java.lang.Long", javalang.LongClass),
				jexl.WithClass("java.lang.Double", javalang.DoubleClass),
				jexl.WithClass("java.lang.Character", javalang.CharacterClass),
				jexl.WithClass("java.lang.String", javalang.StringClass),
				jexl.WithClass("java.lang.StringBuilder", javalang.StringBuilderClass),
				jexl.WithClass("java.lang.StringBuffer", javalang.StringBufferClass),
				jexl.WithClass("java.lang.Math", javamath.MathClass),
				jexl.WithClass("java.math.BigInteger", javamath.BigIntegerClass),
				jexl.WithClass("java.math.BigDecimal", javamath.BigDecimalClass),
				jexl.WithClass("java.util.ArrayList", javautil.ArrayListClass),
				jexl.WithClass("java.util.LinkedList", javautil.LinkedListClass),
				jexl.WithClass("java.util.HashMap", javautil.HashMapClass),
				jexl.WithClass("java.util.TreeMap", javautil.TreeMapClass),
				jexl.WithClass("java.util.Stack", javautil.StackClass),
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

			if test_.Error != "" {
				t.Errorf("Expect error got %v", out)
				t.Logf("name: %s", test_.ID)
				t.Logf("category: %s", test_.Category)
				t.Logf("description: %s", test_.Description)
				t.Logf("expression: %v", test_.Expr)
				return
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
// Custom Classes for synthetic tests
//

type pointClass struct{}

func (pointClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "new":
		m := map[string]any{}
		if len(args) >= 1 {
			m["x"] = args[0]
		}
		if len(args) >= 2 {
			m["y"] = args[1]
		}
		return m, nil
	}
	return nil, fmt.Errorf("Point.%s: undefined", method)
}

type emptyClass struct{}

func (emptyClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "new":
		return map[string]any{}, nil
	}
	return nil, fmt.Errorf("Empty.%s: undefined", method)
}

//
// Custom Classes for namespace unit tests
//

// stringNS is the String namespace used in spec tests.
type stringNS struct{}

func (stringNS) Call(method string, args ...any) (any, error) {
	switch method {
	case "length":
		if len(args) != 1 {
			return nil, fmt.Errorf("String.length: expected 1 argument")
		}
		return len([]rune(coerce.ToString(args[0]))), nil
	case "toUpperCase":
		if len(args) != 1 {
			return nil, fmt.Errorf("String.toUpperCase: expected 1 argument")
		}
		return strings.ToUpper(coerce.ToString(args[0])), nil
	case "toLowerCase":
		if len(args) != 1 {
			return nil, fmt.Errorf("String.toLowerCase: expected 1 argument")
		}
		return strings.ToLower(coerce.ToString(args[0])), nil
	case "trim":
		if len(args) != 1 {
			return nil, fmt.Errorf("String.trim: expected 1 argument")
		}
		return strings.TrimSpace(coerce.ToString(args[0])), nil
	case "includes":
		if len(args) != 2 {
			return nil, fmt.Errorf("String.includes: expected 2 arguments")
		}
		return strings.Contains(coerce.ToString(args[0]), coerce.ToString(args[1])), nil
	case "startsWith":
		if len(args) != 2 {
			return nil, fmt.Errorf("String.startsWith: expected 2 arguments")
		}
		return strings.HasPrefix(coerce.ToString(args[0]), coerce.ToString(args[1])), nil
	case "endsWith":
		if len(args) != 2 {
			return nil, fmt.Errorf("String.endsWith: expected 2 arguments")
		}
		return strings.HasSuffix(coerce.ToString(args[0]), coerce.ToString(args[1])), nil
	case "substring":
		if len(args) < 2 || len(args) > 3 {
			return nil, fmt.Errorf("String.substring: expected 2 or 3 arguments")
		}
		runes := []rune(coerce.ToString(args[0]))
		from := coerce.ToInt(args[1])
		to := len(runes)
		if len(args) == 3 {
			to = coerce.ToInt(args[2])
		}
		if from < 0 || to > len(runes) || from > to {
			return nil, fmt.Errorf("String.substring: index out of range")
		}
		return string(runes[from:to]), nil
	case "indexOf":
		if len(args) != 2 {
			return nil, fmt.Errorf("String.indexOf: expected 2 arguments")
		}
		return strings.Index(coerce.ToString(args[0]), coerce.ToString(args[1])), nil
	case "lastIndexOf":
		if len(args) != 2 {
			return nil, fmt.Errorf("String.lastIndexOf: expected 2 arguments")
		}
		return strings.LastIndex(coerce.ToString(args[0]), coerce.ToString(args[1])), nil
	case "replace":
		if len(args) != 3 {
			return nil, fmt.Errorf("String.replace: expected 3 arguments")
		}
		return strings.ReplaceAll(coerce.ToString(args[0]), coerce.ToString(args[1]), coerce.ToString(args[2])), nil
	case "repeat":
		if len(args) != 2 {
			return nil, fmt.Errorf("String.repeat: expected 2 arguments")
		}
		return strings.Repeat(coerce.ToString(args[0]), coerce.ToInt(args[1])), nil
	case "concat":
		var b strings.Builder
		for _, a := range args {
			b.WriteString(coerce.ToString(a))
		}
		return b.String(), nil
	case "fromCharCode":
		var b strings.Builder
		for _, a := range args {
			b.WriteRune(rune(coerce.ToInt32(a)))
		}
		return b.String(), nil
	case "charCodeAt":
		if len(args) != 2 {
			return nil, fmt.Errorf("String.charCodeAt: expected 2 arguments")
		}
		runes := []rune(coerce.ToString(args[0]))
		i := coerce.ToInt(args[1])
		if i < 0 || i >= len(runes) {
			return nil, fmt.Errorf("String.charCodeAt: index %d out of range", i)
		}
		return int(runes[i]), nil
	case "padStart":
		if len(args) != 3 {
			return nil, fmt.Errorf("String.padStart: expected 3 arguments")
		}
		str := coerce.ToString(args[0])
		target := coerce.ToInt(args[1])
		pad := coerce.ToString(args[2])
		runes := []rune(str)
		padRunes := []rune(pad)
		if len(runes) > target {
			return string(runes[len(runes)-target:]), nil
		}
		need := target - len(runes)
		var prefix []rune
		for len(prefix) < need {
			prefix = append(prefix, padRunes...)
		}
		return string(prefix[:need]) + str, nil
	case "padEnd":
		if len(args) != 3 {
			return nil, fmt.Errorf("String.padEnd: expected 3 arguments")
		}
		str := coerce.ToString(args[0])
		target := coerce.ToInt(args[1])
		pad := coerce.ToString(args[2])
		runes := []rune(str)
		padRunes := []rune(pad)
		if len(runes) > target {
			return string(runes[:target]), nil
		}
		need := target - len(runes)
		var suffix []rune
		for len(suffix) < need {
			suffix = append(suffix, padRunes...)
		}
		return str + string(suffix[:need]), nil
	}
	return nil, fmt.Errorf("String.%s: undefined", method)
}

// integerNS is the Integer namespace used in spec tests.
type integerNS struct{}

func (integerNS) Call(method string, args ...any) (any, error) {
	switch method {
	case "parseInt":
		switch len(args) {
		case 1:
			return coerce.ToInt64(args[0]), nil
		case 2:
			n, err := strconv.ParseInt(coerce.ToString(args[0]), coerce.ToInt(args[1]), 64)
			if err != nil {
				return nil, fmt.Errorf("Integer.parseInt: %w", err)
			}
			return n, nil
		default:
			return nil, fmt.Errorf("Integer.parseInt: expected 1 or 2 arguments")
		}
	case "parseFloat":
		if len(args) != 1 {
			return nil, fmt.Errorf("Integer.parseFloat: expected 1 argument")
		}
		return strconv.ParseFloat(strings.TrimSpace(coerce.ToString(args[0])), 64)
	case "isNaN":
		if len(args) != 1 {
			return nil, fmt.Errorf("Integer.isNaN: expected 1 argument")
		}
		return math.IsNaN(coerce.ToFloat64(args[0])), nil
	case "isFinite":
		if len(args) != 1 {
			return nil, fmt.Errorf("Integer.isFinite: expected 1 argument")
		}
		f := coerce.ToFloat64(args[0])
		return !math.IsInf(f, 0) && !math.IsNaN(f), nil
	case "MAX_VALUE", "maxValue":
		return int64(math.MaxInt32), nil
	case "MIN_VALUE", "minValue":
		return int64(math.MinInt32), nil
	}
	return nil, fmt.Errorf("Integer.%s: undefined", method)
}

var (
	// StringNS is the String namespace for spec tests.
	StringNS classes.Object = stringNS{}
	// IntegerNS is the Integer namespace for spec tests.
	IntegerNS classes.Object = integerNS{}
)
