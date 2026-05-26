// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lang

import (
	"testing"
)

// Ensure new StringBuffer() with no args returns an empty *StringBuffer.
func TestStringBuffer_newNoArgs(t *testing.T) {
	inst, err := StringBufferClass.Call("new")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	sb, ok := inst.(*StringBuffer)
	if !ok {
		t.Errorf("expected *StringBuffer, got %T", inst)
		return
	}
	if sb.ToString() != "" {
		t.Errorf("expected empty string, got %q", sb.ToString())
	}
}

// Ensure append and toString work correctly on a StringBuffer.
func TestStringBuffer_appendAndToString(t *testing.T) {
	inst, err := StringBufferClass.Call("new", "hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	sb, ok := inst.(*StringBuffer)
	if !ok {
		t.Errorf("expected *StringBuffer, got %T", inst)
		return
	}
	if _, err = sb.Call("append", " world"); err != nil {
		t.Errorf("unexpected error from append: %v", err)
		return
	}
	got, err := sb.Call("toString")
	if err != nil {
		t.Errorf("unexpected error from toString: %v", err)
		return
	}
	if got != "hello world" {
		t.Errorf("expected \"hello world\", got %v", got)
	}
}

// Ensure an unknown class method returns an error.
func TestStringBuffer_unknownClassMethod(t *testing.T) {
	if _, err := StringBufferClass.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown class method")
	}
}
