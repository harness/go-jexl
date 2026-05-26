// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/harness/go-jexl/jexl/coerce"
)

// CharAt returns the Unicode code point at index i in v.
func CharAt(v any, i any) rune {
	runes := []rune(coerce.ToString(v))
	n := coerce.ToInt(i)
	if n < 0 || n >= len(runes) {
		return 0
	}
	return runes[n]
}

// CodePointAt returns the code point at index i in v.
func CodePointAt(v any, i any) int {
	s := coerce.ToString(v)
	n := coerce.ToInt(i)
	if n < 0 || n >= len(s) {
		return 0
	}
	return int(s[n])
}

// CompareTo returns the ordering of v vs other.
func CompareTo(v any, other any) int {
	switch v.(type) {
	case int, int8, int16, int32, int64:
		a := coerce.ToInt64(v)
		b := coerce.ToInt64(other)
		switch {
		case a < b:
			return -1
		case a > b:
			return 1
		default:
			return 0
		}
	case float32, float64:
		a := coerce.ToFloat64(v)
		b := coerce.ToFloat64(other)
		switch {
		case a < b:
			return -1
		case a > b:
			return 1
		default:
			return 0
		}
	case bool:
		a := coerce.ToBool(v)
		b := coerce.ToBool(other)
		switch {
		case a == b:
			return 0
		case b:
			return -1
		default:
			return 1
		}
	default:
		return strings.Compare(
			coerce.ToString(v),
			coerce.ToString(other),
		)
	}
}

// CompareToIgnoreCase returns case-insensitive lexicographic ordering.
func CompareToIgnoreCase(v any, other any) int {
	return strings.Compare(
		strings.ToLower(coerce.ToString(v)),
		strings.ToLower(coerce.ToString(other)),
	)
}

// Concat returns v concatenated with each arg.
func Concat(v any, args ...any) string {
	var b strings.Builder
	b.WriteString(coerce.ToString(v))
	for _, vv := range args {
		b.WriteString(coerce.ToString(vv))
	}
	return b.String()
}

// Contains reports whether v contains sub.
func Contains(v any, sub any) bool {
	return strings.Contains(
		coerce.ToString(v),
		coerce.ToString(sub),
	)
}

// EndsWith reports whether v ends with suffix.
func EndsWith(v any, suffix any) bool {
	return strings.HasSuffix(
		coerce.ToString(v),
		coerce.ToString(suffix),
	)
}

// Equals reports whether v equals other (string comparison).
func Equals(v any, other any) bool {
	return coerce.ToString(v) == coerce.ToString(other)
}

// EqualsIgnoreCase reports case-insensitive string equality.
func EqualsIgnoreCase(v any, other any) bool {
	return strings.EqualFold(
		coerce.ToString(v),
		coerce.ToString(other),
	)
}

// Formatted returns fmt.Sprintf(v, args...).
func Formatted(v any, args ...any) string {
	return fmt.Sprintf(
		coerce.ToString(v), args...,
	)
}

// IndexOf returns the byte index of sub in v, or -1.
// Accepts an optional fromIndex to start the search.
func IndexOf(v any, sub any, args ...any) int {
	s := coerce.ToString(v)
	substr := coerce.ToString(sub)
	if len(args) == 0 {
		return strings.Index(s, substr)
	}
	start := coerce.ToIntSlice(args)
	i := start[0]
	switch {
	case i < 0:
		return -1
	case i > len(s)-1:
		return -1
	default:
		idx := strings.Index(s[i:], substr)
		if idx < 0 {
			return -1
		}
		return i + idx
	}
}

// IsBlank reports whether v is empty or whitespace-only.
func IsBlank(v any) bool {
	return strings.TrimSpace(
		coerce.ToString(v),
	) == ""
}

// IsEmpty reports whether v is the empty string.
func IsEmpty(v any) bool {
	return coerce.ToString(v) == ""
}

// LastIndexOf returns the last byte index of sub in v.
func LastIndexOf(v any, sub any) int {
	return strings.LastIndex(
		coerce.ToString(v),
		coerce.ToString(sub),
	)
}

// Length returns the byte length of v.
func Length(v any) int {
	return len(coerce.ToString(v))
}

// Matches reports whether v fully matches the regex pattern.
func Matches(v any, pattern any) bool {
	matched, err := regexp.MatchString("^(?:"+coerce.ToString(pattern)+")$", coerce.ToString(v))
	return err == nil && matched
}

// Repeat returns v repeated n times.
func Repeat(v any, n any) string {
	return strings.Repeat(
		coerce.ToString(v),
		coerce.ToInt(n),
	)
}

// Replace returns v with all occurrences of old replaced by new.
func Replace(v any, old, new any) string {
	return strings.ReplaceAll(
		coerce.ToString(v),
		coerce.ToString(old),
		coerce.ToString(new),
	)
}

// ReplaceAll replaces all matches of regex pattern in v with repl.
func ReplaceAll(v any, pattern, repl any) string {
	re, err := regexp.Compile(coerce.ToString(pattern))
	if err != nil {
		return coerce.ToString(v)
	}
	return re.ReplaceAllString(coerce.ToString(v), coerce.ToString(repl))
}

// ReplaceFirst replaces the first match of regex pattern in v with repl.
func ReplaceFirst(v any, pattern, repl any) string {
	s := coerce.ToString(v)
	re, err := regexp.Compile(coerce.ToString(pattern))
	if err != nil {
		return s
	}
	loc := re.FindStringIndex(s)
	if loc == nil {
		return s
	}
	return s[:loc[0]] + re.ReplaceAllString(s[loc[0]:loc[1]], coerce.ToString(repl)) + s[loc[1]:]
}

// Split splits v by sep.
func Split(v any, sep any) []string {
	return strings.Split(
		coerce.ToString(v),
		coerce.ToString(sep),
	)
}

// StartsWith reports whether v starts with prefix.
func StartsWith(v any, prefix any) bool {
	return strings.HasPrefix(
		coerce.ToString(v),
		coerce.ToString(prefix),
	)
}

// Strip returns v with leading and trailing whitespace removed.
func Strip(v any) string {
	return strings.TrimSpace(
		coerce.ToString(v),
	)
}

// StripLeading returns v with leading whitespace removed.
func StripLeading(v any) string {
	const newlines = " \t\n\r"
	return strings.TrimLeft(
		coerce.ToString(v),
		newlines,
	)
}

// StripTrailing returns v with trailing whitespace removed.
func StripTrailing(v any) string {
	const newlines = " \t\n\r"
	return strings.TrimRight(
		coerce.ToString(v),
		newlines,
	)
}

// Substring returns v[start:] or v[start:end].
func Substring(v any, start any, args ...any) string {
	str := coerce.ToString(v)
	s := coerce.ToInt(start)
	e := len(str)
	if len(args) > 0 {
		ends := coerce.ToIntSlice(args)
		e = ends[0]
	}
	switch {
	case s < 0:
		return ""
	case e > len(str):
		return ""
	case s > e:
		return ""
	default:
		return str[s:e]
	}
}

// ToCharArray returns v as a slice of single-character strings.
func ToCharArray(v any) []string {
	s := coerce.ToString(v)
	var out []string
	for _, r := range s {
		out = append(out, string(r))
	}
	return out
}

// ToLowerCase returns v in lower case.
func ToLowerCase(v any) string {
	return strings.ToLower(
		coerce.ToString(v),
	)
}

// ToUpperCase returns v in upper case.
func ToUpperCase(v any) string {
	return strings.ToUpper(
		coerce.ToString(v),
	)
}

// Trim returns v with leading and trailing whitespace removed.
func Trim(v any) string {
	return strings.TrimSpace(
		coerce.ToString(v),
	)
}

// TrimLeft returns v with leading whitespace removed.
func TrimLeft(v any) string {
	const newlines = " \t\n\r"
	return strings.TrimLeft(
		coerce.ToString(v), newlines)
}

// TrimRight returns v with trailing whitespace removed.
func TrimRight(v any) string {
	const newlines = " \t\n\r"
	return strings.TrimRight(
		coerce.ToString(v), newlines)
}

// SubstringBefore returns the portion of v before the first occurrence of delim.
// Returns v unchanged if delim is not found.
func SubstringBefore(v any, delim any) string {
	s := coerce.ToString(v)
	before, _, found := strings.Cut(s, coerce.ToString(delim))
	if !found {
		return s
	}
	return before
}
