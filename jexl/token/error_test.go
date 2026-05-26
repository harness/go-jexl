// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package token

import (
	"errors"
	"testing"
)

// Ensure Error.Error returns the message when no snippet is present.
func TestError_Error_noSnippet(t *testing.T) {
	e := &Error{Message: "something went wrong"}
	if e.Error() != "something went wrong" {
		t.Fatalf("unexpected: %q", e.Error())
	}
}

// Ensure Error.Error includes line, column, and snippet when bound.
func TestError_Error_withSnippet(t *testing.T) {
	e := &Error{Range: Range{From: 4, To: 5}, Message: "bad token"}
	e.Bind(NewFile("foo bar"))
	msg := e.Error()
	if msg == "bad token" {
		t.Fatal("expected snippet to be included in message")
	}
	// should contain line:col and the source line
	if e.Line != 1 || e.Column != 4 {
		t.Fatalf("unexpected line=%d col=%d", e.Line, e.Column)
	}
}

// Ensure Bind sets line and column for a single-line source.
func TestError_Bind_lineColumn(t *testing.T) {
	e := &Error{Range: Range{From: 6, To: 7}, Message: "oops"}
	e.Bind(NewFile("hello world"))
	if e.Line != 1 {
		t.Fatalf("expected line 1, got %d", e.Line)
	}
	if e.Column != 6 {
		t.Fatalf("expected column 6, got %d", e.Column)
	}
}

// Ensure Bind sets the correct line for a multi-line source.
func TestError_Bind_multiLine(t *testing.T) {
	// "foo\nbar\nbaz" — 'b' of "bar" is rune offset 4
	e := &Error{Range: Range{From: 4, To: 5}, Message: "oops"}
	e.Bind(NewFile("foo\nbar\nbaz"))
	if e.Line != 2 {
		t.Fatalf("expected line 2, got %d", e.Line)
	}
	if e.Column != 0 {
		t.Fatalf("expected column 0, got %d", e.Column)
	}
}

// Ensure Bind produces a non-empty snippet for a line with content.
func TestError_Bind_snippet(t *testing.T) {
	e := &Error{Range: Range{From: 2, To: 3}, Message: "oops"}
	e.Bind(NewFile("hello"))
	if e.Snippet == "" {
		t.Fatal("expected non-empty snippet")
	}
}

// Ensure Bind returns without a snippet when the line is empty.
func TestError_Bind_emptyLine(t *testing.T) {
	// source is just a newline — line 1 has no content
	e := &Error{Range: Range{From: 0, To: 1}, Message: "oops"}
	e.Bind(NewFile("\n"))
	if e.Snippet != "" {
		t.Fatalf("expected empty snippet for empty line, got %q", e.Snippet)
	}
}

// Ensure Unwrap returns the wrapped error.
func TestError_Unwrap(t *testing.T) {
	inner := errors.New("inner")
	e := &Error{Prev: inner}
	if !errors.Is(e, inner) {
		t.Fatal("expected errors.Is to find inner error")
	}
}

// Ensure Wrap sets the previous error.
func TestError_Wrap(t *testing.T) {
	inner := errors.New("cause")
	e := &Error{}
	e.Wrap(inner)
	if e.Prev != inner {
		t.Fatal("expected Prev to be set")
	}
}
