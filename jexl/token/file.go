// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package token

import "strings"

// File holds the raw source text of an expression.
type File struct {
	raw string
}

// NewFile returns a File wrapping the given source text.
func NewFile(contents string) File {
	return File{raw: contents}
}

// String returns the raw source text.
func (f File) String() string {
	return f.raw
}

// Snippet returns the line-th line of the source (1-based). Returns ("", false)
// if the line does not exist.
func (f File) Snippet(line int) (string, bool) {
	if f.raw == "" {
		return "", false
	}
	var start int
	for i := 1; i < line; i++ {
		pos := strings.IndexByte(f.raw[start:], '\n')
		if pos < 0 {
			return "", false
		}
		start += pos + 1
	}
	end := start + strings.IndexByte(f.raw[start:], '\n')
	if end < start {
		end = len(f.raw)
	}
	return f.raw[start:end], true
}
