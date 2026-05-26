// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package template

import (
	"fmt"
	"strings"

	"github.com/harness/go-jexl/jexl/coerce"
)

// Template evaluates JEXL expressions embedded in a string,
// replacing each expression with its stringified result.
type Template struct {
	eval  func(expr string) (any, error)
	left  string
	right string
}

// New returns a Template using eval to execute expressions.
func New(eval func(expr string) (any, error)) *Template {
	return &Template{eval: eval, left: "<+", right: ">"}
}

// Delim sets the left and right delimiters and returns t.
func (t *Template) Delim(left, right string) *Template {
	t.left = left
	t.right = right
	return t
}

// Exec replaces every delimited expression in b with its result.
func (t *Template) Exec(b []byte) ([]byte, error) {
	out, err := t.ExecString(string(b))
	if err != nil {
		return nil, err
	}
	return []byte(out), nil
}

// ExecString replaces every delimited expression in s with its result.
func (t *Template) ExecString(s string) (string, error) {
	for {
		next, err := t.execOnce(s)
		if err != nil {
			return "", err
		}
		if next == s {
			return s, nil
		}
		s = next
	}
}

// execOnce finds the innermost expression, evaluates it,
// and splices the result back.
func (t *Template) execOnce(s string) (string, error) {
	start := strings.LastIndex(s, t.left)
	if start < 0 {
		return s, nil
	}
	body, end, err := t.scanToClose(s, start+len(t.left))
	if err != nil {
		return "", err
	}
	result, err := t.eval(body)
	if err != nil {
		return "", err
	}
	// When the splice point is inside an outer open expression,
	// string results must be single-quoted. Note that JEXL supports
	// quoted property access: foo.'name' and foo['name']
	nested := strings.LastIndex(s[:start], t.left) > strings.LastIndex(s[:start], t.right)
	var replacement string
	if nested {
		if str, ok := result.(string); ok {
			replacement = "'" + strings.ReplaceAll(str, "'", "\\'") + "'"
		} else {
			replacement = coerce.ToString(result)
		}
	} else {
		replacement = coerce.ToString(result)
	}
	return s[:start] + replacement + s[end+len(t.right):], nil
}

// helper function scans forward from pos for the closing
// right delimiter at bracket/paren/string depth 0. Returns
// the body and the index of the closer.
func (t *Template) scanToClose(s string, pos int) (body string, closeIdx int, err error) {
	depth := 0
	inSingle := false
	inDouble := false
	i := pos
	for i < len(s) {
		if inSingle {
			if s[i] == '\\' && i+1 < len(s) {
				i += 2
				continue
			}
			if s[i] == '\'' {
				inSingle = false
			}
			i++
			continue
		}
		if inDouble {
			if s[i] == '\\' && i+1 < len(s) {
				i += 2
				continue
			}
			if s[i] == '"' {
				inDouble = false
			}
			i++
			continue
		}
		if depth == 0 && strings.HasPrefix(s[i:], t.right) {
			return s[pos:i], i, nil
		}
		switch s[i] {
		case '\'':
			inSingle = true
		case '"':
			inDouble = true
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			depth--
		}
		i++
	}
	err = fmt.Errorf("template: unclosed expression: %q:%d", t.left, pos-len(t.left))
	return
}
