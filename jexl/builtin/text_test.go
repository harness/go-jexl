// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import "testing"

// Ensure CharAt returns the rune at the given index.
func TestCharAt_valid(t *testing.T) {
	if CharAt("hello", 1) != 'e' {
		t.Fatal("expected 'e'")
	}
}

// Ensure CharAt returns 0 for an out-of-bounds index.
func TestCharAt_outOfBounds(t *testing.T) {
	if CharAt("hi", 10) != 0 {
		t.Fatal("expected 0")
	}
}

// Ensure CodePointAt returns the byte value at the given index.
func TestCodePointAt_valid(t *testing.T) {
	if CodePointAt("ABC", 0) != int('A') {
		t.Fatal("expected 65")
	}
}

// Ensure CodePointAt returns 0 for an out-of-bounds index.
func TestCodePointAt_outOfBounds(t *testing.T) {
	if CodePointAt("hi", 5) != 0 {
		t.Fatal("expected 0")
	}
}

// Ensure CompareTo returns -1 when a < b (int).
func TestCompareTo_intLess(t *testing.T) {
	if CompareTo(1, 2) != -1 {
		t.Fatal("expected -1")
	}
}

// Ensure CompareTo returns 1 when a > b (int).
func TestCompareTo_intGreater(t *testing.T) {
	if CompareTo(3, 2) != 1 {
		t.Fatal("expected 1")
	}
}

// Ensure CompareTo returns 0 when equal (string).
func TestCompareTo_stringEqual(t *testing.T) {
	if CompareTo("abc", "abc") != 0 {
		t.Fatal("expected 0")
	}
}

// Ensure CompareToIgnoreCase treats case-different strings as equal.
func TestCompareToIgnoreCase_equal(t *testing.T) {
	if CompareToIgnoreCase("Hello", "hello") != 0 {
		t.Fatal("expected 0")
	}
}

// Ensure Concat joins multiple arguments onto the base string.
func TestConcat_multiple(t *testing.T) {
	if Concat("a", "b", "c") != "abc" {
		t.Fatal("expected abc")
	}
}

// Ensure Contains reports true when substring is present.
func TestContains_present(t *testing.T) {
	if !Contains("foobar", "oba") {
		t.Fatal("expected true")
	}
}

// Ensure Contains reports false when substring is absent.
func TestContains_absent(t *testing.T) {
	if Contains("foobar", "xyz") {
		t.Fatal("expected false")
	}
}

// Ensure EndsWith reports true when the suffix matches.
func TestEndsWith_match(t *testing.T) {
	if !EndsWith("foobar", "bar") {
		t.Fatal("expected true")
	}
}

// Ensure Equals performs exact string comparison.
func TestEquals_same(t *testing.T) {
	if !Equals("hello", "hello") {
		t.Fatal("expected true")
	}
}

// Ensure EqualsIgnoreCase ignores case.
func TestEqualsIgnoreCase_different(t *testing.T) {
	if !EqualsIgnoreCase("HELLO", "hello") {
		t.Fatal("expected true")
	}
}

// Ensure Formatted produces a printf-style string.
func TestFormatted_basic(t *testing.T) {
	if Formatted("x=%d", 7) != "x=7" {
		t.Fatal("expected x=7")
	}
}

// Ensure IndexOf returns the correct byte offset.
func TestIndexOf_found(t *testing.T) {
	if IndexOf("foobar", "bar") != 3 {
		t.Fatal("expected 3")
	}
}

// Ensure IndexOf returns -1 when not found.
func TestIndexOf_notFound(t *testing.T) {
	if IndexOf("foobar", "xyz") != -1 {
		t.Fatal("expected -1")
	}
}

// Ensure IndexOf respects the fromIndex argument.
func TestIndexOf_fromIndex(t *testing.T) {
	if IndexOf("abcabc", "bc", 2) != 4 {
		t.Fatal("expected 4")
	}
}

// Ensure IsBlank returns true for a whitespace-only string.
func TestIsBlank_whitespace(t *testing.T) {
	if !IsBlank("   ") {
		t.Fatal("expected true")
	}
}

// Ensure IsEmpty returns true for an empty string.
func TestIsEmpty_empty(t *testing.T) {
	if !IsEmpty("") {
		t.Fatal("expected true")
	}
}

// Ensure LastIndexOf returns the last occurrence index.
func TestLastIndexOf_found(t *testing.T) {
	if LastIndexOf("abcabc", "bc") != 4 {
		t.Fatal("expected 4")
	}
}

// Ensure Length returns the byte length of the string.
func TestLength_ascii(t *testing.T) {
	if Length("hello") != 5 {
		t.Fatal("expected 5")
	}
}

// Ensure Matches performs full-string regex matching.
func TestMatches_full(t *testing.T) {
	if !Matches("hello123", `[a-z]+\d+`) {
		t.Fatal("expected true")
	}
}

// Ensure Matches returns false when pattern doesn't cover the full string.
func TestMatches_partial(t *testing.T) {
	if Matches("hello123!", `[a-z]+\d+`) {
		t.Fatal("expected false")
	}
}

// Ensure Repeat duplicates the string n times.
func TestRepeat_three(t *testing.T) {
	if Repeat("ab", 3) != "ababab" {
		t.Fatal("expected ababab")
	}
}

// Ensure Replace substitutes all occurrences.
func TestReplace_all(t *testing.T) {
	if Replace("aabbaa", "a", "x") != "xxbbxx" {
		t.Fatal("expected xxbbxx")
	}
}

// Ensure ReplaceAll substitutes all regex matches.
func TestReplaceAll_digits(t *testing.T) {
	if ReplaceAll("a1b2c3", `\d`, "X") != "aXbXcX" {
		t.Fatal("expected aXbXcX")
	}
}

// Ensure ReplaceFirst substitutes only the first regex match.
func TestReplaceFirst_single(t *testing.T) {
	if ReplaceFirst("a1b2c3", `\d`, "X") != "aXb2c3" {
		t.Fatal("expected aXb2c3")
	}
}

// Ensure Split splits on the given separator.
func TestSplit_comma(t *testing.T) {
	parts := Split("a,b,c", ",")
	if len(parts) != 3 || parts[1] != "b" {
		t.Fatalf("unexpected result: %v", parts)
	}
}

// Ensure StartsWith reports true when the prefix matches.
func TestStartsWith_match(t *testing.T) {
	if !StartsWith("foobar", "foo") {
		t.Fatal("expected true")
	}
}

// Ensure Strip removes leading and trailing whitespace.
func TestStrip_whitespace(t *testing.T) {
	if Strip("  hi  ") != "hi" {
		t.Fatal("expected hi")
	}
}

// Ensure StripLeading removes only leading whitespace.
func TestStripLeading_leading(t *testing.T) {
	if StripLeading("  hi  ") != "hi  " {
		t.Fatal("expected 'hi  '")
	}
}

// Ensure StripTrailing removes only trailing whitespace.
func TestStripTrailing_trailing(t *testing.T) {
	if StripTrailing("  hi  ") != "  hi" {
		t.Fatal("expected '  hi'")
	}
}

// Ensure Substring returns the correct slice with start and end.
func TestSubstring_startEnd(t *testing.T) {
	if Substring("hello", 1, 4) != "ell" {
		t.Fatal("expected ell")
	}
}

// Ensure Substring returns from start to end of string when no end given.
func TestSubstring_startOnly(t *testing.T) {
	if Substring("hello", 2) != "llo" {
		t.Fatal("expected llo")
	}
}

// Ensure ToCharArray splits string into single-char slices.
func TestToCharArray_basic(t *testing.T) {
	arr := ToCharArray("abc")
	if len(arr) != 3 || arr[0] != "a" || arr[2] != "c" {
		t.Fatalf("unexpected result: %v", arr)
	}
}

// Ensure ToLowerCase lowercases the string.
func TestToLowerCase_basic(t *testing.T) {
	if ToLowerCase("HELLO") != "hello" {
		t.Fatal("expected hello")
	}
}

// Ensure ToUpperCase uppercases the string.
func TestToUpperCase_basic(t *testing.T) {
	if ToUpperCase("hello") != "HELLO" {
		t.Fatal("expected HELLO")
	}
}

// Ensure Trim removes leading and trailing whitespace.
func TestTrim_whitespace(t *testing.T) {
	if Trim("\t hello \n") != "hello" {
		t.Fatal("expected hello")
	}
}

// Ensure TrimLeft removes only leading whitespace.
func TestTrimLeft_leading(t *testing.T) {
	if TrimLeft("\t hello \n") != "hello \n" {
		t.Fatal("expected 'hello \\n'")
	}
}

// Ensure TrimRight removes only trailing whitespace.
func TestTrimRight_trailing(t *testing.T) {
	if TrimRight("\t hello \n") != "\t hello" {
		t.Fatal("expected '\\t hello'")
	}
}

// Ensure SubstringBefore returns the portion before the delimiter.
func TestSubstringBefore_found(t *testing.T) {
	if SubstringBefore("foo/bar/baz", "/") != "foo" {
		t.Fatal("expected foo")
	}
}

// Ensure SubstringBefore returns the full string when delimiter is absent.
func TestSubstringBefore_notFound(t *testing.T) {
	if SubstringBefore("foobar", "/") != "foobar" {
		t.Fatal("expected foobar")
	}
}

// Ensure Reverse returns the string with characters in reverse order.
func TestReverse_basic(t *testing.T) {
	if Reverse("hello") != "olleh" {
		t.Fatal("expected olleh")
	}
}

// Ensure Reverse handles multibyte Unicode characters correctly.
func TestReverse_unicode(t *testing.T) {
	if Reverse("héllo") != "olléh" {
		t.Fatal("expected olléh")
	}
}

// Ensure Quote wraps the string in double quotes and escapes special characters.
func TestQuote_basic(t *testing.T) {
	if Quote("hello") != `"hello"` {
		t.Fatal(`expected "hello"`)
	}
}

// Ensure Quote escapes newlines and other special characters.
func TestQuote_special(t *testing.T) {
	if Quote("line1\nline2") != `"line1\nline2"` {
		t.Fatal(`expected "line1\nline2"`)
	}
}
