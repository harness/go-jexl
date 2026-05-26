// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lang

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/harness/go-jexl/jexl/coerce"
)

// StringClass is the java.lang.String class object.
var StringClass stringClass

type stringClass struct{}

func (stringClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "new":
		if len(args) == 0 {
			return NewString(""), nil
		}
		return NewStringFrom(args[0]), nil
	case "valueOf":
		if len(args) != 1 {
			return nil, fmt.Errorf("String.valueOf: expected 1 argument")
		}
		return NewStringFrom(args[0]), nil
	case "format":
		if len(args) < 1 {
			return nil, fmt.Errorf("String.format: expected at least 1 argument")
		}
		return String(fmt.Sprintf(coerce.ToString(args[0]), args[1:]...)), nil
	}
	return nil, fmt.Errorf("String.%s: undefined", method)
}

// String mirrors java.lang.String.
type String string

// NewString wraps a Go string as a String.
func NewString(v string) String {
	return String(v)
}

// NewStringFrom coerces any value to a String.
func NewStringFrom(v any) String {
	return String(coerce.ToString(v))
}

// Length returns the byte length of the string.
func (s String) Length() int {
	return len(string(s))
}

// CharAt returns the Unicode code point at index i, or 0 if out of bounds.
func (s String) CharAt(i int) rune {
	runes := []rune(string(s))
	if i < 0 || i >= len(runes) {
		return 0
	}
	return runes[i]
}

// CodePointAt returns the Unicode code point value at index i as an int, or 0 if out of bounds.
func (s String) CodePointAt(i int) int {
	runes := []rune(string(s))
	if i < 0 || i >= len(runes) {
		return 0
	}
	return int(runes[i])
}

// Contains reports whether sub is within s.
func (s String) Contains(sub string) bool {
	return strings.Contains(string(s), sub)
}

// StartsWith reports whether s begins with prefix.
func (s String) StartsWith(prefix string) bool {
	return strings.HasPrefix(string(s), prefix)
}

// EndsWith reports whether s ends with suffix.
func (s String) EndsWith(suffix string) bool {
	return strings.HasSuffix(string(s), suffix)
}

// IndexOf returns the index of the first instance of sub starting at optional
// fromIndex, or -1 if not found.
func (s String) IndexOf(sub string, fromIndex ...int) int {
	str := string(s)
	if len(fromIndex) == 0 {
		return strings.Index(str, sub)
	}
	i := fromIndex[0]
	switch {
	case i < 0:
		return -1
	case i > len(str)-1:
		return -1
	default:
		idx := strings.Index(str[i:], sub)
		if idx < 0 {
			return -1
		}
		return i + idx
	}
}

// LastIndexOf returns the index of the last instance of sub, or -1.
func (s String) LastIndexOf(sub string) int {
	return strings.LastIndex(string(s), sub)
}

// Substring returns a substring using func.go's clamping logic.
func (s String) Substring(from int, to ...int) string {
	str := string(s)
	e := len(str)
	if len(to) > 0 {
		e = to[0]
	}
	switch {
	case from < 0:
		return ""
	case e > len(str):
		return ""
	case from > e:
		return ""
	default:
		return str[from:e]
	}
}

// Replace returns a copy with all occurrences of old replaced by new.
func (s String) Replace(old, new string) string {
	return strings.ReplaceAll(string(s), old, new)
}

// Split splits s around sep and returns a []string.
func (s String) Split(sep string) []string {
	return strings.Split(string(s), sep)
}

// ToUpperCase returns s in uppercase.
func (s String) ToUpperCase() string {
	return strings.ToUpper(string(s))
}

// ToLowerCase returns s in lowercase.
func (s String) ToLowerCase() string {
	return strings.ToLower(string(s))
}

// Trim returns s with leading and trailing whitespace removed.
func (s String) Trim() string {
	return strings.TrimSpace(string(s))
}

// IsEmpty reports whether s has zero length.
func (s String) IsEmpty() bool {
	return len(string(s)) == 0
}

// IsBlank reports whether s is empty or whitespace-only.
func (s String) IsBlank() bool {
	return strings.TrimSpace(string(s)) == ""
}

// Repeat returns a new string consisting of n copies of s.
func (s String) Repeat(n int) string {
	return strings.Repeat(string(s), n)
}

// CompareTo returns -1, 0, or 1 comparing s to other lexicographically.
func (s String) CompareTo(other string) int {
	return strings.Compare(string(s), other)
}

// CompareToIgnoreCase returns case-insensitive lexicographic ordering.
func (s String) CompareToIgnoreCase(other string) int {
	return strings.Compare(strings.ToLower(string(s)), strings.ToLower(other))
}

// Equals reports whether s equals other.
func (s String) Equals(other string) bool {
	return string(s) == other
}

// EqualsIgnoreCase reports case-insensitive string equality.
func (s String) EqualsIgnoreCase(other string) bool {
	return strings.EqualFold(string(s), other)
}

// Matches reports whether s fully matches the regex pattern.
func (s String) Matches(pattern string) bool {
	matched, err := regexp.MatchString("^(?:"+pattern+")$", string(s))
	return err == nil && matched
}

// ReplaceAll replaces all matches of regex pattern in s with repl.
func (s String) ReplaceAll(pattern, repl string) string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return string(s)
	}
	return re.ReplaceAllString(string(s), repl)
}

// ReplaceFirst replaces the first match of regex pattern in s with repl.
func (s String) ReplaceFirst(pattern, repl string) string {
	str := string(s)
	re, err := regexp.Compile(pattern)
	if err != nil {
		return str
	}
	loc := re.FindStringIndex(str)
	if loc == nil {
		return str
	}
	return str[:loc[0]] + re.ReplaceAllString(str[loc[0]:loc[1]], repl) + str[loc[1]:]
}

// Strip returns s with leading and trailing whitespace removed.
func (s String) Strip() string {
	return strings.TrimSpace(string(s))
}

// StripLeading returns s with leading whitespace removed.
func (s String) StripLeading() string {
	const whitespace = " \t\n\r"
	return strings.TrimLeft(string(s), whitespace)
}

// StripTrailing returns s with trailing whitespace removed.
func (s String) StripTrailing() string {
	const whitespace = " \t\n\r"
	return strings.TrimRight(string(s), whitespace)
}

// ToCharArray returns s as a slice of single-character strings.
func (s String) ToCharArray() []string {
	var out []string
	for _, r := range string(s) {
		out = append(out, string(r))
	}
	return out
}

// Formatted returns fmt.Sprintf(s, args...).
func (s String) Formatted(args ...any) string {
	return fmt.Sprintf(string(s), args...)
}

// Concat returns s concatenated with each arg.
func (s String) Concat(args ...any) string {
	var b strings.Builder
	b.WriteString(string(s))
	for _, v := range args {
		b.WriteString(coerce.ToString(v))
	}
	return b.String()
}

// ByteValue returns the string parsed as int8.
func (s String) ByteValue() int8 {
	return coerce.ToInt8(string(s))
}

// ShortValue returns the string parsed as int16.
func (s String) ShortValue() int16 {
	return coerce.ToInt16(string(s))
}

// IntValue returns the string parsed as int32.
func (s String) IntValue() int32 {
	return coerce.ToInt32(string(s))
}

// LongValue returns the string parsed as int64.
func (s String) LongValue() int64 {
	return coerce.ToInt64(string(s))
}

// FloatValue returns the string parsed as float32.
func (s String) FloatValue() float32 {
	return coerce.ToFloat32(string(s))
}

// DoubleValue returns the string parsed as float64.
func (s String) DoubleValue() float64 {
	return coerce.ToFloat64(string(s))
}

// BooleanValue returns the string parsed as bool.
func (s String) BooleanValue() bool {
	return coerce.ToBool(string(s))
}

// ToString returns the Go string value.
func (s String) ToString() string {
	return string(s)
}

// Call dispatches instance methods.
func (s String) Call(method string, args ...any) (any, error) {
	switch method {
	case "length":
		return s.Length(), nil
	case "charAt":
		if len(args) != 1 {
			return nil, fmt.Errorf("charAt: expected 1 argument")
		}
		return s.CharAt(coerce.ToInt(args[0])), nil
	case "contains":
		if len(args) != 1 {
			return nil, fmt.Errorf("contains: expected 1 argument")
		}
		return s.Contains(coerce.ToString(args[0])), nil
	case "startsWith":
		if len(args) != 1 {
			return nil, fmt.Errorf("startsWith: expected 1 argument")
		}
		return s.StartsWith(coerce.ToString(args[0])), nil
	case "endsWith":
		if len(args) != 1 {
			return nil, fmt.Errorf("endsWith: expected 1 argument")
		}
		return s.EndsWith(coerce.ToString(args[0])), nil
	case "indexOf":
		if len(args) < 1 || len(args) > 2 {
			return nil, fmt.Errorf("indexOf: expected 1 or 2 arguments")
		}
		if len(args) == 2 {
			return s.IndexOf(coerce.ToString(args[0]), coerce.ToInt(args[1])), nil
		}
		return s.IndexOf(coerce.ToString(args[0])), nil
	case "lastIndexOf":
		if len(args) != 1 {
			return nil, fmt.Errorf("lastIndexOf: expected 1 argument")
		}
		return s.LastIndexOf(coerce.ToString(args[0])), nil
	case "substring":
		if len(args) < 1 || len(args) > 2 {
			return nil, fmt.Errorf("substring: expected 1 or 2 arguments")
		}
		if len(args) == 2 {
			return s.Substring(coerce.ToInt(args[0]), coerce.ToInt(args[1])), nil
		}
		return s.Substring(coerce.ToInt(args[0])), nil
	case "replace":
		if len(args) != 2 {
			return nil, fmt.Errorf("replace: expected 2 arguments")
		}
		return s.Replace(coerce.ToString(args[0]), coerce.ToString(args[1])), nil
	case "replaceAll":
		if len(args) != 2 {
			return nil, fmt.Errorf("replaceAll: expected 2 arguments")
		}
		return s.ReplaceAll(coerce.ToString(args[0]), coerce.ToString(args[1])), nil
	case "replaceFirst":
		if len(args) != 2 {
			return nil, fmt.Errorf("replaceFirst: expected 2 arguments")
		}
		return s.ReplaceFirst(coerce.ToString(args[0]), coerce.ToString(args[1])), nil
	case "split":
		if len(args) != 1 {
			return nil, fmt.Errorf("split: expected 1 argument")
		}
		return s.Split(coerce.ToString(args[0])), nil
	case "toUpperCase":
		return s.ToUpperCase(), nil
	case "toLowerCase":
		return s.ToLowerCase(), nil
	case "trim":
		return s.Trim(), nil
	case "strip":
		return s.Strip(), nil
	case "stripLeading":
		return s.StripLeading(), nil
	case "stripTrailing":
		return s.StripTrailing(), nil
	case "isEmpty":
		return s.IsEmpty(), nil
	case "isBlank":
		return s.IsBlank(), nil
	case "repeat":
		if len(args) != 1 {
			return nil, fmt.Errorf("repeat: expected 1 argument")
		}
		return s.Repeat(coerce.ToInt(args[0])), nil
	case "toString":
		return s.ToString(), nil
	case "compareTo":
		if len(args) != 1 {
			return nil, fmt.Errorf("compareTo: expected 1 argument")
		}
		return s.CompareTo(coerce.ToString(args[0])), nil
	case "compareToIgnoreCase":
		if len(args) != 1 {
			return nil, fmt.Errorf("compareToIgnoreCase: expected 1 argument")
		}
		return s.CompareToIgnoreCase(coerce.ToString(args[0])), nil
	case "equals":
		if len(args) != 1 {
			return nil, fmt.Errorf("equals: expected 1 argument")
		}
		return s.Equals(coerce.ToString(args[0])), nil
	case "equalsIgnoreCase":
		if len(args) != 1 {
			return nil, fmt.Errorf("equalsIgnoreCase: expected 1 argument")
		}
		return s.EqualsIgnoreCase(coerce.ToString(args[0])), nil
	case "matches":
		if len(args) != 1 {
			return nil, fmt.Errorf("matches: expected 1 argument")
		}
		return s.Matches(coerce.ToString(args[0])), nil
	case "toCharArray":
		return s.ToCharArray(), nil
	case "formatted":
		return s.Formatted(args...), nil
	case "concat":
		return s.Concat(args...), nil
	case "codePointAt":
		if len(args) != 1 {
			return nil, fmt.Errorf("codePointAt: expected 1 argument")
		}
		return s.CodePointAt(coerce.ToInt(args[0])), nil
	case "byteValue", "toByte":
		return s.ByteValue(), nil
	case "shortValue", "toShort":
		return s.ShortValue(), nil
	case "intValue", "toInteger":
		return s.IntValue(), nil
	case "longValue", "toLong":
		return s.LongValue(), nil
	case "floatValue", "toFloat":
		return s.FloatValue(), nil
	case "doubleValue", "toDouble":
		return s.DoubleValue(), nil
	case "booleanValue", "toBoolean":
		return s.BooleanValue(), nil
	case "default":
		return s, nil
	}
	return nil, fmt.Errorf("String instance: undefined method %q", method)
}
