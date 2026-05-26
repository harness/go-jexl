// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lang

import (
	"fmt"
	"strings"

	"github.com/harness/go-jexl/jexl/coerce"
)

// StringBuilderClass is the java.lang.StringBuilder class object.
var StringBuilderClass stringBuilderClass

type stringBuilderClass struct{}

func (stringBuilderClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "new":
		sb := &StringBuilder{}
		if len(args) == 1 {
			sb.buf.WriteString(coerce.ToString(args[0]))
		}
		return sb, nil
	}
	return nil, fmt.Errorf("StringBuilder.%s: undefined", method)
}

// StringBuilder mirrors java.lang.StringBuilder.
type StringBuilder struct {
	buf strings.Builder
}

// Append appends s to the builder and returns the builder.
func (sb *StringBuilder) Append(s string) *StringBuilder {
	sb.buf.WriteString(s)
	return sb
}

// Insert inserts s at position idx and returns the builder.
func (sb *StringBuilder) Insert(idx int, s string) *StringBuilder {
	runes := []rune(sb.buf.String())
	if idx < 0 {
		idx = 0
	}
	if idx > len(runes) {
		idx = len(runes)
	}
	ins := []rune(s)
	result := make([]rune, 0, len(runes)+len(ins))
	result = append(result, runes[:idx]...)
	result = append(result, ins...)
	result = append(result, runes[idx:]...)
	sb.buf.Reset()
	sb.buf.WriteString(string(result))
	return sb
}

// Delete removes characters from from (inclusive) to to (exclusive) and returns the builder.
func (sb *StringBuilder) Delete(from, to int) *StringBuilder {
	runes := []rune(sb.buf.String())
	if from < 0 {
		from = 0
	}
	if to > len(runes) {
		to = len(runes)
	}
	if from > to {
		from = to
	}
	result := make([]rune, 0, len(runes)-(to-from))
	result = append(result, runes[:from]...)
	result = append(result, runes[to:]...)
	sb.buf.Reset()
	sb.buf.WriteString(string(result))
	return sb
}

// DeleteCharAt removes the character at index i and returns the builder.
func (sb *StringBuilder) DeleteCharAt(i int) *StringBuilder {
	runes := []rune(sb.buf.String())
	if i < 0 || i >= len(runes) {
		return sb
	}
	result := make([]rune, 0, len(runes)-1)
	result = append(result, runes[:i]...)
	result = append(result, runes[i+1:]...)
	sb.buf.Reset()
	sb.buf.WriteString(string(result))
	return sb
}

// Replace replaces characters from from (inclusive) to to (exclusive) with s and returns the builder.
func (sb *StringBuilder) Replace(from, to int, s string) *StringBuilder {
	runes := []rune(sb.buf.String())
	if from < 0 {
		from = 0
	}
	if to > len(runes) {
		to = len(runes)
	}
	if from > to {
		from = to
	}
	ins := []rune(s)
	result := make([]rune, 0, from+len(ins)+len(runes)-to)
	result = append(result, runes[:from]...)
	result = append(result, ins...)
	result = append(result, runes[to:]...)
	sb.buf.Reset()
	sb.buf.WriteString(string(result))
	return sb
}

// Reverse reverses the content of the builder and returns it.
func (sb *StringBuilder) Reverse() *StringBuilder {
	runes := []rune(sb.buf.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	sb.buf.Reset()
	sb.buf.WriteString(string(runes))
	return sb
}

// Length returns the number of Unicode code points in the builder.
func (sb *StringBuilder) Length() int {
	return len([]rune(sb.buf.String()))
}

// CharAt returns the Unicode code point at index i, or an error if out of bounds.
func (sb *StringBuilder) CharAt(i int) (rune, error) {
	runes := []rune(sb.buf.String())
	if i < 0 || i >= len(runes) {
		return 0, fmt.Errorf("StringBuilder.charAt: index %d out of bounds (length %d)", i, len(runes))
	}
	return runes[i], nil
}

// Substring returns a substring from from (inclusive) to optional to (exclusive).
func (sb *StringBuilder) Substring(from int, to ...int) (string, error) {
	runes := []rune(sb.buf.String())
	end := len(runes)
	if len(to) > 0 {
		end = to[0]
	}
	if from < 0 || from > len(runes) {
		return "", fmt.Errorf("StringBuilder.substring: from index %d out of bounds (length %d)", from, len(runes))
	}
	if end < from || end > len(runes) {
		return "", fmt.Errorf("StringBuilder.substring: to index %d out of bounds (length %d)", end, len(runes))
	}
	return string(runes[from:end]), nil
}

// IndexOf returns the index of the first occurrence of s, or -1.
func (sb *StringBuilder) IndexOf(s string) int {
	return strings.Index(sb.buf.String(), s)
}

// ToString returns the string content of the builder.
func (sb *StringBuilder) ToString() string {
	return sb.buf.String()
}

// Call dispatches instance methods.
func (sb *StringBuilder) Call(method string, args ...any) (any, error) {
	switch method {
	case "append":
		if len(args) != 1 {
			return nil, fmt.Errorf("append: expected 1 argument")
		}
		return sb.Append(coerce.ToString(args[0])), nil
	case "insert":
		if len(args) != 2 {
			return nil, fmt.Errorf("insert: expected 2 arguments")
		}
		return sb.Insert(coerce.ToInt(args[0]), coerce.ToString(args[1])), nil
	case "delete":
		if len(args) != 2 {
			return nil, fmt.Errorf("delete: expected 2 arguments")
		}
		return sb.Delete(coerce.ToInt(args[0]), coerce.ToInt(args[1])), nil
	case "deleteCharAt":
		if len(args) != 1 {
			return nil, fmt.Errorf("deleteCharAt: expected 1 argument")
		}
		return sb.DeleteCharAt(coerce.ToInt(args[0])), nil
	case "replace":
		if len(args) != 3 {
			return nil, fmt.Errorf("replace: expected 3 arguments")
		}
		return sb.Replace(coerce.ToInt(args[0]), coerce.ToInt(args[1]), coerce.ToString(args[2])), nil
	case "reverse":
		return sb.Reverse(), nil
	case "length":
		return sb.Length(), nil
	case "charAt":
		if len(args) != 1 {
			return nil, fmt.Errorf("charAt: expected 1 argument")
		}
		return sb.CharAt(coerce.ToInt(args[0]))
	case "substring":
		if len(args) < 1 || len(args) > 2 {
			return nil, fmt.Errorf("substring: expected 1 or 2 arguments")
		}
		if len(args) == 2 {
			return sb.Substring(coerce.ToInt(args[0]), coerce.ToInt(args[1]))
		}
		return sb.Substring(coerce.ToInt(args[0]))
	case "indexOf":
		if len(args) != 1 {
			return nil, fmt.Errorf("indexOf: expected 1 argument")
		}
		return sb.IndexOf(coerce.ToString(args[0])), nil
	case "toString", "string":
		return sb.ToString(), nil
	}
	return nil, fmt.Errorf("StringBuilder instance: undefined method %q", method)
}
