// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package token

import (
	"fmt"
	"strings"
)

// Error is a parse or scan error anchored to a source range.
type Error struct {
	Range
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Message string `json:"message"`
	Snippet string `json:"snippet"`
	Prev    error  `json:"prev"`
}

func (e *Error) Error() string {
	return e.format()
}

var tabReplacer = strings.NewReplacer("\t", " ")

// Bind resolves the error's rune-offset range to a line, column, and snippet
// using the provided source file. It returns e so callers can chain.
func (e *Error) Bind(source File) *Error {
	src := source.String()

	var runeCount, lineStart int
	e.Line = 1
	e.Column = 0
	for i, r := range src {
		if runeCount == e.From {
			break
		}
		if r == '\n' {
			lineStart = i + 1
			e.Line++
			e.Column = 0
		} else {
			e.Column++
		}
		runeCount++
	}

	lineEnd := lineStart + strings.IndexByte(src[lineStart:], '\n')
	if lineEnd < lineStart {
		lineEnd = len(src)
	}
	if lineStart == lineEnd {
		return e
	}

	const prefix = "\n | "
	line := src[lineStart:lineEnd]
	snippet := new(strings.Builder)
	snippet.Grow(2*len(prefix) + len(line) + e.Column + 1)
	snippet.WriteString(prefix)
	tabReplacer.WriteString(snippet, line)
	snippet.WriteString(prefix)
	for i := 0; i < e.Column; i++ {
		snippet.WriteByte('.')
	}
	snippet.WriteByte('^')
	e.Snippet = snippet.String()
	return e
}

func (e *Error) Unwrap() error {
	return e.Prev
}

func (e *Error) Wrap(err error) {
	e.Prev = err
}

func (e *Error) format() string {
	if e.Snippet == "" {
		return e.Message
	}
	return fmt.Sprintf(
		"%s (%d:%d)%s",
		e.Message,
		e.Line,
		e.Column+1,
		e.Snippet,
	)
}
