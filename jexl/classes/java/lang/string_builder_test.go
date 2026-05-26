// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lang

import (
	"testing"
)

// Ensure new StringBuilder() with no args returns an empty builder.
func TestStringBuilder_newNoArgs(t *testing.T) {
	inst, err := StringBuilderClass.Call("new")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	sb, ok := inst.(*StringBuilder)
	if !ok {
		t.Errorf("expected *StringBuilder, got %T", inst)
		return
	}
	if sb.ToString() != "" {
		t.Errorf("expected empty string, got %q", sb.ToString())
	}
}

// Ensure new StringBuilder("hello") initialises the builder with "hello".
func TestStringBuilder_newWithString(t *testing.T) {
	inst, err := StringBuilderClass.Call("new", "hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	sb, ok := inst.(*StringBuilder)
	if !ok {
		t.Errorf("expected *StringBuilder, got %T", inst)
		return
	}
	if sb.ToString() != "hello" {
		t.Errorf("expected \"hello\", got %q", sb.ToString())
	}
}

// Ensure append followed by toString produces the concatenated string.
func TestStringBuilder_appendAndToString(t *testing.T) {
	inst, err := StringBuilderClass.Call("new", "hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	sb, ok := inst.(*StringBuilder)
	if !ok {
		t.Errorf("expected *StringBuilder, got %T", inst)
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

// Ensure string alias for toString works.
func TestStringBuilder_stringAlias(t *testing.T) {
	inst, _ := StringBuilderClass.Call("new", "abc")
	sb := inst.(*StringBuilder)
	got, err := sb.Call("string")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "abc" {
		t.Errorf("expected \"abc\", got %v", got)
	}
}

// Ensure insert at index 2 inserts the substring at the correct position.
func TestStringBuilder_insert(t *testing.T) {
	inst, _ := StringBuilderClass.Call("new", "helo")
	sb := inst.(*StringBuilder)
	if _, err := sb.Call("insert", 3, "l"); err != nil {
		t.Errorf("unexpected error from insert: %v", err)
		return
	}
	got, _ := sb.Call("toString")
	if got != "hello" {
		t.Errorf("expected \"hello\", got %v", got)
	}
}

// Ensure delete removes the specified range.
func TestStringBuilder_delete(t *testing.T) {
	inst, _ := StringBuilderClass.Call("new", "hello world")
	sb := inst.(*StringBuilder)
	if _, err := sb.Call("delete", 5, 11); err != nil {
		t.Errorf("unexpected error from delete: %v", err)
		return
	}
	got, _ := sb.Call("toString")
	if got != "hello" {
		t.Errorf("expected \"hello\", got %v", got)
	}
}

// Ensure deleteCharAt removes the character at the given index.
func TestStringBuilder_deleteCharAt(t *testing.T) {
	inst, _ := StringBuilderClass.Call("new", "hello")
	sb := inst.(*StringBuilder)
	if _, err := sb.Call("deleteCharAt", 4); err != nil {
		t.Errorf("unexpected error from deleteCharAt: %v", err)
		return
	}
	got, _ := sb.Call("toString")
	if got != "hell" {
		t.Errorf("expected \"hell\", got %v", got)
	}
}

// Ensure replace substitutes the specified range.
func TestStringBuilder_replace(t *testing.T) {
	inst, _ := StringBuilderClass.Call("new", "hello world")
	sb := inst.(*StringBuilder)
	if _, err := sb.Call("replace", 6, 11, "Go"); err != nil {
		t.Errorf("unexpected error from replace: %v", err)
		return
	}
	got, _ := sb.Call("toString")
	if got != "hello Go" {
		t.Errorf("expected \"hello Go\", got %v", got)
	}
}

// Ensure charAt returns the character at the given index.
func TestStringBuilder_charAt(t *testing.T) {
	inst, _ := StringBuilderClass.Call("new", "hello")
	sb := inst.(*StringBuilder)
	got, err := sb.Call("charAt", 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != rune('e') {
		t.Errorf("expected 'e', got %v", got)
	}
}

// Ensure charAt out of bounds returns an error.
func TestStringBuilder_charAtOutOfBounds(t *testing.T) {
	inst, _ := StringBuilderClass.Call("new", "hello")
	sb := inst.(*StringBuilder)
	if _, err := sb.Call("charAt", 10); err == nil {
		t.Error("expected error for out-of-bounds charAt")
	}
}

// Ensure substring with one arg returns from index to end.
func TestStringBuilder_substring(t *testing.T) {
	inst, _ := StringBuilderClass.Call("new", "hello world")
	sb := inst.(*StringBuilder)
	got, err := sb.Call("substring", 6)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "world" {
		t.Errorf("expected \"world\", got %v", got)
	}
}

// Ensure substring with two args returns the specified range.
func TestStringBuilder_substringRange(t *testing.T) {
	inst, _ := StringBuilderClass.Call("new", "hello world")
	sb := inst.(*StringBuilder)
	got, err := sb.Call("substring", 0, 5)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "hello" {
		t.Errorf("expected \"hello\", got %v", got)
	}
}

// Ensure indexOf returns the correct position.
func TestStringBuilder_indexOf(t *testing.T) {
	inst, _ := StringBuilderClass.Call("new", "hello world")
	sb := inst.(*StringBuilder)
	got, err := sb.Call("indexOf", "world")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 6 {
		t.Errorf("expected 6, got %v", got)
	}
}

// Ensure reverse on "abc" produces "cba".
func TestStringBuilder_reverse(t *testing.T) {
	inst, err := StringBuilderClass.Call("new", "abc")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	sb, ok := inst.(*StringBuilder)
	if !ok {
		t.Errorf("expected *StringBuilder, got %T", inst)
		return
	}
	if _, err = sb.Call("reverse"); err != nil {
		t.Errorf("unexpected error from reverse: %v", err)
		return
	}
	got, err := sb.Call("toString")
	if err != nil {
		t.Errorf("unexpected error from toString: %v", err)
		return
	}
	if got != "cba" {
		t.Errorf("expected \"cba\", got %v", got)
	}
}

// Ensure length on a builder initialised with "hello" returns 5.
func TestStringBuilder_length(t *testing.T) {
	inst, err := StringBuilderClass.Call("new", "hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	sb, ok := inst.(*StringBuilder)
	if !ok {
		t.Errorf("expected *StringBuilder, got %T", inst)
		return
	}
	got, err := sb.Call("length")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 5 {
		t.Errorf("expected 5, got %v", got)
	}
}

// Ensure an unknown class method returns an error.
func TestStringBuilder_unknownClassMethod(t *testing.T) {
	if _, err := StringBuilderClass.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown class method")
	}
}

// Ensure an unknown instance method returns an error.
func TestStringBuilder_unknownInstanceMethod(t *testing.T) {
	inst, err := StringBuilderClass.Call("new")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	sb, ok := inst.(*StringBuilder)
	if !ok {
		t.Errorf("expected *StringBuilder, got %T", inst)
		return
	}
	if _, err = sb.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown instance method")
	}
}

// Ensure append with wrong arg count returns an error.
func TestStringBuilder_appendArgCount(t *testing.T) {
	inst, _ := StringBuilderClass.Call("new")
	sb := inst.(*StringBuilder)
	if _, err := sb.Call("append"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure insert with wrong arg count returns an error.
func TestStringBuilder_insertArgCount(t *testing.T) {
	inst, _ := StringBuilderClass.Call("new")
	sb := inst.(*StringBuilder)
	if _, err := sb.Call("insert", 0); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure delete with wrong arg count returns an error.
func TestStringBuilder_deleteArgCount(t *testing.T) {
	inst, _ := StringBuilderClass.Call("new")
	sb := inst.(*StringBuilder)
	if _, err := sb.Call("delete", 0); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure deleteCharAt with wrong arg count returns an error.
func TestStringBuilder_deleteCharAtArgCount(t *testing.T) {
	inst, _ := StringBuilderClass.Call("new")
	sb := inst.(*StringBuilder)
	if _, err := sb.Call("deleteCharAt"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure replace with wrong arg count returns an error.
func TestStringBuilder_replaceArgCount(t *testing.T) {
	inst, _ := StringBuilderClass.Call("new")
	sb := inst.(*StringBuilder)
	if _, err := sb.Call("replace", 0, 1); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure charAt with wrong arg count returns an error.
func TestStringBuilder_charAtArgCount(t *testing.T) {
	inst, _ := StringBuilderClass.Call("new")
	sb := inst.(*StringBuilder)
	if _, err := sb.Call("charAt"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure substring with wrong arg count returns an error.
func TestStringBuilder_substringArgCount(t *testing.T) {
	inst, _ := StringBuilderClass.Call("new")
	sb := inst.(*StringBuilder)
	if _, err := sb.Call("substring"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure indexOf with wrong arg count returns an error.
func TestStringBuilder_indexOfArgCount(t *testing.T) {
	inst, _ := StringBuilderClass.Call("new")
	sb := inst.(*StringBuilder)
	if _, err := sb.Call("indexOf"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}
