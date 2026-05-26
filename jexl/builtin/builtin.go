// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import (
	"reflect"

	"github.com/harness/go-jexl/jexl/functions"
	"github.com/harness/go-jexl/jexl/vm/runtime"
)

// Funcs is the default set of built-in functions available
// in every expression.
var Funcs map[string]*runtime.Function

func init() {
	funcs := map[string]any{
		// math
		"abs":        Abs,
		"ceil":       Ceil,
		"floor":      Floor,
		"isInfinite": IsInfinite,
		"isNaN":      IsNaN,
		"log":        Log,
		"log2":       Log2,
		"log10":      Log10,
		"max":        Max,
		"min":        Min,
		"pow":        Pow,
		"round":      Round,
		"sqrt":       Sqrt,

		// text
		"charAt":              CharAt,
		"codePointAt":         CodePointAt,
		"compareTo":           CompareTo,
		"compareToIgnoreCase": CompareToIgnoreCase,
		"concat":              Concat,
		"contains":            Contains,
		"endsWith":            EndsWith,
		"equals":              Equals,
		"equalsIgnoreCase":    EqualsIgnoreCase,
		"formatted":           Formatted,
		"indexOf":             IndexOf,
		"isBlank":             IsBlank,
		"isEmpty":             IsEmpty,
		"lastIndexOf":         LastIndexOf,
		"length":              Length,
		"matches":             Matches,
		"repeat":              Repeat,
		"replace":             Replace,
		"replaceAll":          ReplaceAll,
		"replaceFirst":        ReplaceFirst,
		"split":               Split,
		"startsWith":          StartsWith,
		"strip":               Strip,
		"stripLeading":        StripLeading,
		"stripTrailing":       StripTrailing,
		"substring":           Substring,
		"toCharArray":         ToCharArray,
		"toLowerCase":         ToLowerCase,
		"toUpperCase":         ToUpperCase,
		"toUpper":             ToUpperCase, // alias
		"upper":               ToUpperCase, // alias
		"trim":                Trim,
		"trimLeft":            TrimLeft,
		"trimRight":           TrimRight,
		"substringBefore":     SubstringBefore,

		// convert
		"booleanValue": BooleanValue,
		"byteValue":    ByteValue,
		"doubleValue":  DoubleValue,
		"floatValue":   FloatValue,
		"intValue":     IntValue,
		"toInteger":    IntValue, // alias
		"longValue":    LongValue,
		"shortValue":   ShortValue,
		"toString":     ToString,
		"default":      DefaultValue,

		// date
		"currentDate": CurrentDate,
		"currentTime": CurrentTime,
		"dateFormat":  DateFormat,
		"plusMinutes": PlusMinutes,
		"plusHours":   PlusHours,
		"plusDays":    PlusDays,

		// encode
		"base64Encode":    Base64Encode,
		"base64Decode":    Base64Decode,
		"base64URLEncode": Base64URLEncode,
		"base64URLDecode": Base64URLDecode,
		"base64RawEncode": Base64RawEncode,
		"base64RawDecode": Base64RawDecode,

		// hex
		"hexEncode": HexEncode,
		"hexDecode": HexDecode,

		// fmt
		"sprintf": Sprintf,
		"sprint":  Sprint,

		// json
		"jsonMarshal":       JSONMarshal,
		"jsonUnmarshal":     JSONUnmarshal,
		"jsonMarshalIndent": JSONMarshalIndent,
		"jsonSelect":        JSONSelect,
		"jsonList":          JSONList,

		// regex
		"regexExtract": RegexExtract,

		// xml
		"xmlMarshal":       XMLMarshal,
		"xmlMarshalIndent": XMLMarshalIndent,
		"xmlSelect":        XMLSelect,
	}

	Funcs = make(map[string]*runtime.Function, len(funcs))
	for name, fn := range funcs {
		wrapped, t := functions.Wrap(fn)
		Funcs[name] = &runtime.Function{
			Name:  name,
			Types: []reflect.Type{t},
			Func:  wrapped,
		}
	}
}
