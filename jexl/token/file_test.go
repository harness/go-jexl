// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package token

import "testing"

// Ensure NewFile round-trips through String.
func TestFile_String(t *testing.T) {
	f := NewFile("hello world")
	if f.String() != "hello world" {
		t.Fatalf("unexpected: %q", f.String())
	}
}

// Ensure an empty File returns the empty string.
func TestFile_String_empty(t *testing.T) {
	f := NewFile("")
	if f.String() != "" {
		t.Fatalf("expected empty string, got %q", f.String())
	}
}

// Ensure Snippet returns false for an empty file.
func TestFile_Snippet_emptyFile(t *testing.T) {
	f := NewFile("")
	_, ok := f.Snippet(1)
	if ok {
		t.Fatal("expected ok=false for empty file")
	}
}

// Ensure Snippet returns the first line correctly.
func TestFile_Snippet_firstLine(t *testing.T) {
	f := NewFile("foo\nbar\nbaz")
	s, ok := f.Snippet(1)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if s != "foo" {
		t.Fatalf("expected foo, got %q", s)
	}
}

// Ensure Snippet returns a middle line correctly.
func TestFile_Snippet_middleLine(t *testing.T) {
	f := NewFile("foo\nbar\nbaz")
	s, ok := f.Snippet(2)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if s != "bar" {
		t.Fatalf("expected bar, got %q", s)
	}
}

// Ensure Snippet returns the last line when there is no trailing newline.
func TestFile_Snippet_lastLine(t *testing.T) {
	f := NewFile("foo\nbar\nbaz")
	s, ok := f.Snippet(3)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if s != "baz" {
		t.Fatalf("expected baz, got %q", s)
	}
}

// Ensure Snippet returns false for a line beyond the end of the file.
func TestFile_Snippet_beyondEnd(t *testing.T) {
	f := NewFile("foo\nbar")
	_, ok := f.Snippet(3)
	if ok {
		t.Fatal("expected ok=false for line beyond end")
	}
}

// Ensure Snippet handles a single-line file with no newline.
func TestFile_Snippet_singleLine(t *testing.T) {
	f := NewFile("only")
	s, ok := f.Snippet(1)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if s != "only" {
		t.Fatalf("expected only, got %q", s)
	}
}
