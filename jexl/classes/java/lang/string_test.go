// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lang

import (
	"testing"
)

// Ensure new String() with no args returns empty string.
func TestString_newNoArgs(t *testing.T) {
	inst, err := StringClass.Call("new")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if inst != NewString("") {
		t.Errorf("expected String(\"\"), got %v", inst)
	}
}

// Ensure new String("hello") returns String("hello").
func TestString_newWithArg(t *testing.T) {
	inst, err := StringClass.Call("new", "hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if inst != NewString("hello") {
		t.Errorf("expected String(\"hello\"), got %v", inst)
	}
}

// Ensure String.valueOf returns String("42") for int 42.
func TestString_valueOf(t *testing.T) {
	got, err := StringClass.Call("valueOf", 42)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewString("42") {
		t.Errorf("expected String(\"42\"), got %v", got)
	}
}

// Ensure String.format returns the formatted string.
func TestString_format(t *testing.T) {
	got, err := StringClass.Call("format", "hello %s", "world")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != String("hello world") {
		t.Errorf("expected \"hello world\", got %v", got)
	}
}

// Ensure String.valueOf with wrong arg count returns an error.
func TestString_valueOfArgCount(t *testing.T) {
	if _, err := StringClass.Call("valueOf"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure String.format with no args returns an error.
func TestString_formatArgCount(t *testing.T) {
	if _, err := StringClass.Call("format"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure an unknown class method returns an error.
func TestString_unknownClassMethod(t *testing.T) {
	if _, err := StringClass.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown class method")
	}
}

// Ensure length on String("hello") returns 5.
func TestString_instanceLength(t *testing.T) {
	s := NewString("hello")
	got, err := s.Call("length")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 5 {
		t.Errorf("expected 5, got %v", got)
	}
}

// Ensure charAt on String("hello") at index 1 returns 'e'.
func TestString_instanceCharAt(t *testing.T) {
	s := NewString("hello")
	got, err := s.Call("charAt", 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != rune('e') {
		t.Errorf("expected 'e', got %v", got)
	}
}

// Ensure codePointAt on String("hello") at index 0 returns int('h').
func TestString_instanceCodePointAt(t *testing.T) {
	s := NewString("hello")
	got, err := s.Call("codePointAt", 0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int('h') {
		t.Errorf("expected int('h'), got %v", got)
	}
}

// Ensure toUpperCase on String("hello") returns "HELLO".
func TestString_instanceToUpperCase(t *testing.T) {
	s := NewString("hello")
	got, err := s.Call("toUpperCase")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "HELLO" {
		t.Errorf("expected \"HELLO\", got %v", got)
	}
}

// Ensure toLowerCase on String("HELLO") returns "hello".
func TestString_instanceToLowerCase(t *testing.T) {
	s := NewString("HELLO")
	got, err := s.Call("toLowerCase")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "hello" {
		t.Errorf("expected \"hello\", got %v", got)
	}
}

// Ensure contains returns true when the substring is present.
func TestString_instanceContains(t *testing.T) {
	s := NewString("hello world")
	got, err := s.Call("contains", "world")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure startsWith returns true when the prefix matches.
func TestString_instanceStartsWith(t *testing.T) {
	s := NewString("hello world")
	got, err := s.Call("startsWith", "hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure endsWith returns true when the suffix matches.
func TestString_instanceEndsWith(t *testing.T) {
	s := NewString("hello world")
	got, err := s.Call("endsWith", "world")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure indexOf returns 6 for "world" in "hello world".
func TestString_instanceIndexOf(t *testing.T) {
	s := NewString("hello world")
	got, err := s.Call("indexOf", "world")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 6 {
		t.Errorf("expected 6, got %v", got)
	}
}

// Ensure indexOf with fromIndex returns correct position.
func TestString_instanceIndexOfFromIndex(t *testing.T) {
	s := NewString("abcabc")
	got, err := s.Call("indexOf", "a", 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 3 {
		t.Errorf("expected 3, got %v", got)
	}
}

// Ensure lastIndexOf returns the last occurrence index.
func TestString_instanceLastIndexOf(t *testing.T) {
	s := NewString("abcabc")
	got, err := s.Call("lastIndexOf", "a")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 3 {
		t.Errorf("expected 3, got %v", got)
	}
}

// Ensure substring with one arg returns from index to end.
func TestString_instanceSubstring(t *testing.T) {
	s := NewString("hello world")
	got, err := s.Call("substring", 6)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "world" {
		t.Errorf("expected \"world\", got %v", got)
	}
}

// Ensure substring with two args returns the specified range.
func TestString_instanceSubstringRange(t *testing.T) {
	s := NewString("hello world")
	got, err := s.Call("substring", 0, 5)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "hello" {
		t.Errorf("expected \"hello\", got %v", got)
	}
}

// Ensure replace returns the string with substitution applied.
func TestString_instanceReplace(t *testing.T) {
	s := NewString("foo bar foo")
	got, err := s.Call("replace", "foo", "baz")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "baz bar baz" {
		t.Errorf("expected \"baz bar baz\", got %v", got)
	}
}

// Ensure replaceAll replaces all regex matches.
func TestString_instanceReplaceAll(t *testing.T) {
	s := NewString("foo123bar456")
	got, err := s.Call("replaceAll", `\d+`, "X")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "fooXbarX" {
		t.Errorf("expected \"fooXbarX\", got %v", got)
	}
}

// Ensure replaceFirst replaces only the first regex match.
func TestString_instanceReplaceFirst(t *testing.T) {
	s := NewString("foo123bar456")
	got, err := s.Call("replaceFirst", `\d+`, "X")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "fooXbar456" {
		t.Errorf("expected \"fooXbar456\", got %v", got)
	}
}

// Ensure split on "a,b,c" with "," returns a []string of 3 elements.
func TestString_instanceSplit(t *testing.T) {
	s := NewString("a,b,c")
	got, err := s.Call("split", ",")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	parts, ok := got.([]string)
	if !ok {
		t.Errorf("expected []string, got %T", got)
		return
	}
	if len(parts) != 3 {
		t.Errorf("expected 3 parts, got %d", len(parts))
	}
}

// Ensure trim removes leading and trailing whitespace.
func TestString_instanceTrim(t *testing.T) {
	s := NewString("  hello  ")
	got, err := s.Call("trim")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "hello" {
		t.Errorf("expected \"hello\", got %v", got)
	}
}

// Ensure strip removes leading and trailing whitespace.
func TestString_instanceStrip(t *testing.T) {
	s := NewString("  hello  ")
	got, err := s.Call("strip")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "hello" {
		t.Errorf("expected \"hello\", got %v", got)
	}
}

// Ensure stripLeading removes only leading whitespace.
func TestString_instanceStripLeading(t *testing.T) {
	s := NewString("  hello  ")
	got, err := s.Call("stripLeading")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "hello  " {
		t.Errorf("expected \"hello  \", got %v", got)
	}
}

// Ensure stripTrailing removes only trailing whitespace.
func TestString_instanceStripTrailing(t *testing.T) {
	s := NewString("  hello  ")
	got, err := s.Call("stripTrailing")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "  hello" {
		t.Errorf("expected \"  hello\", got %v", got)
	}
}

// Ensure isEmpty returns true for an empty string.
func TestString_instanceIsEmpty(t *testing.T) {
	s := NewString("")
	got, err := s.Call("isEmpty")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure isBlank returns true for a whitespace-only string.
func TestString_instanceIsBlank(t *testing.T) {
	s := NewString("   ")
	got, err := s.Call("isBlank")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure repeat returns "aaa" for String("a").repeat(3).
func TestString_instanceRepeat(t *testing.T) {
	s := NewString("a")
	got, err := s.Call("repeat", 3)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "aaa" {
		t.Errorf("expected \"aaa\", got %v", got)
	}
}

// Ensure toString returns the Go string value.
func TestString_instanceToString(t *testing.T) {
	s := NewString("hello")
	got, err := s.Call("toString")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "hello" {
		t.Errorf("expected \"hello\", got %v", got)
	}
}

// Ensure compareTo returns 0 when receiver equals arg.
func TestString_instanceCompareToEqual(t *testing.T) {
	s := NewString("abc")
	got, err := s.Call("compareTo", "abc")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 0 {
		t.Errorf("expected 0, got %v", got)
	}
}

// Ensure compareToIgnoreCase returns 0 for case-insensitive match.
func TestString_instanceCompareToIgnoreCase(t *testing.T) {
	s := NewString("ABC")
	got, err := s.Call("compareToIgnoreCase", "abc")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 0 {
		t.Errorf("expected 0, got %v", got)
	}
}

// Ensure equals returns true when values match.
func TestString_instanceEquals(t *testing.T) {
	s := NewString("hello")
	got, err := s.Call("equals", "hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure equalsIgnoreCase returns true for case-insensitive match.
func TestString_instanceEqualsIgnoreCase(t *testing.T) {
	s := NewString("Hello")
	got, err := s.Call("equalsIgnoreCase", "hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure matches returns true for a matching pattern.
func TestString_instanceMatches(t *testing.T) {
	s := NewString("hello")
	got, err := s.Call("matches", `[a-z]+`)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure toCharArray returns a slice of single-character strings.
func TestString_instanceToCharArray(t *testing.T) {
	s := NewString("abc")
	got, err := s.Call("toCharArray")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	arr, ok := got.([]string)
	if !ok {
		t.Errorf("expected []string, got %T", got)
		return
	}
	if len(arr) != 3 || arr[0] != "a" || arr[1] != "b" || arr[2] != "c" {
		t.Errorf("expected [a b c], got %v", arr)
	}
}

// Ensure formatted returns the sprintf result.
func TestString_instanceFormatted(t *testing.T) {
	s := NewString("hello %s")
	got, err := s.Call("formatted", "world")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "hello world" {
		t.Errorf("expected \"hello world\", got %v", got)
	}
}

// Ensure concat returns the concatenated string.
func TestString_instanceConcat(t *testing.T) {
	s := NewString("hello")
	got, err := s.Call("concat", " ", "world")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "hello world" {
		t.Errorf("expected \"hello world\", got %v", got)
	}
}

// Ensure an unknown instance method returns an error.
func TestString_unknownInstanceMethod(t *testing.T) {
	s := NewString("hello")
	if _, err := s.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown instance method")
	}
}

// Ensure charAt with wrong arg count returns an error.
func TestString_charAtArgCount(t *testing.T) {
	s := NewString("hello")
	if _, err := s.Call("charAt"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure codePointAt with wrong arg count returns an error.
func TestString_codePointAtArgCount(t *testing.T) {
	s := NewString("hello")
	if _, err := s.Call("codePointAt"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure startsWith with wrong arg count returns an error.
func TestString_startsWithArgCount(t *testing.T) {
	s := NewString("hello")
	if _, err := s.Call("startsWith"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure endsWith with wrong arg count returns an error.
func TestString_endsWithArgCount(t *testing.T) {
	s := NewString("hello")
	if _, err := s.Call("endsWith"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure indexOf with wrong arg count returns an error.
func TestString_indexOfArgCount(t *testing.T) {
	s := NewString("hello")
	if _, err := s.Call("indexOf"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure lastIndexOf with wrong arg count returns an error.
func TestString_lastIndexOfArgCount(t *testing.T) {
	s := NewString("hello")
	if _, err := s.Call("lastIndexOf"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure substring with wrong arg count returns an error.
func TestString_substringArgCount(t *testing.T) {
	s := NewString("hello")
	if _, err := s.Call("substring"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure replace with wrong arg count returns an error.
func TestString_replaceArgCount(t *testing.T) {
	s := NewString("hello")
	if _, err := s.Call("replace", "h"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure replaceAll with wrong arg count returns an error.
func TestString_replaceAllArgCount(t *testing.T) {
	s := NewString("hello")
	if _, err := s.Call("replaceAll", "h"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure replaceFirst with wrong arg count returns an error.
func TestString_replaceFirstArgCount(t *testing.T) {
	s := NewString("hello")
	if _, err := s.Call("replaceFirst", "h"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure split with wrong arg count returns an error.
func TestString_splitArgCount(t *testing.T) {
	s := NewString("hello")
	if _, err := s.Call("split"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure contains with wrong arg count returns an error.
func TestString_containsArgCount(t *testing.T) {
	s := NewString("hello")
	if _, err := s.Call("contains"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure repeat with wrong arg count returns an error.
func TestString_repeatArgCount(t *testing.T) {
	s := NewString("hello")
	if _, err := s.Call("repeat"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure compareTo with wrong arg count returns an error.
func TestString_compareToArgCount(t *testing.T) {
	s := NewString("hello")
	if _, err := s.Call("compareTo"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure compareToIgnoreCase with wrong arg count returns an error.
func TestString_compareToIgnoreCaseArgCount(t *testing.T) {
	s := NewString("hello")
	if _, err := s.Call("compareToIgnoreCase"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure equals with wrong arg count returns an error.
func TestString_equalsArgCount(t *testing.T) {
	s := NewString("hello")
	if _, err := s.Call("equals"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure equalsIgnoreCase with wrong arg count returns an error.
func TestString_equalsIgnoreCaseArgCount(t *testing.T) {
	s := NewString("hello")
	if _, err := s.Call("equalsIgnoreCase"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure matches with wrong arg count returns an error.
func TestString_matchesArgCount(t *testing.T) {
	s := NewString("hello")
	if _, err := s.Call("matches"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure toInteger on a numeric string returns the parsed int32 value.
func TestString_instanceToInteger(t *testing.T) {
	s := NewString("42")
	got, err := s.Call("toInteger")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int32(42) {
		t.Errorf("expected 42, got %v", got)
	}
}

// Ensure toBoolean on "true" returns true.
func TestString_instanceToBoolean(t *testing.T) {
	s := NewString("true")
	got, err := s.Call("toBoolean")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure toInteger on a non-numeric string returns zero.
func TestString_instanceToIntegerNonNumeric(t *testing.T) {
	s := NewString("abc")
	got, err := s.Call("toInteger")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int32(0) {
		t.Errorf("expected 0, got %v", got)
	}
}

// Ensure default on a non-null String returns the string itself, not the fallback.
func TestString_instanceDefault(t *testing.T) {
	s := NewString("hello")
	got, err := s.Call("default", "fallback")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != s {
		t.Errorf("expected %v, got %v", s, got)
	}
}
