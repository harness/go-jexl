// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lang

import (
	"testing"
	"unicode"
)

// Ensure new Character() with no args returns zero value.
func TestCharacter_newNoArgs(t *testing.T) {
	inst, err := CharacterClass.Call("new")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if inst != NewCharacter(0) {
		t.Errorf("expected Character(0), got %v", inst)
	}
}

// Ensure new Character("A") returns Character('A').
func TestCharacter_newWithArg(t *testing.T) {
	inst, err := CharacterClass.Call("new", "A")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if inst != NewCharacter('A') {
		t.Errorf("expected Character('A'), got %v", inst)
	}
}

// Ensure Character.valueOf returns Character('B').
func TestCharacter_valueOf(t *testing.T) {
	got, err := CharacterClass.Call("valueOf", "B")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewCharacter('B') {
		t.Errorf("expected Character('B'), got %v", got)
	}
}

// Ensure Character.isDigit returns true for "5".
func TestCharacter_isDigit(t *testing.T) {
	got, err := CharacterClass.Call("isDigit", "5")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure Character.isLetter returns true for "A".
func TestCharacter_isLetter(t *testing.T) {
	got, err := CharacterClass.Call("isLetter", "A")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure Character.isLetterOrDigit returns true for "3".
func TestCharacter_isLetterOrDigit(t *testing.T) {
	got, err := CharacterClass.Call("isLetterOrDigit", "3")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure Character.isUpperCase returns true for "Z".
func TestCharacter_isUpperCase(t *testing.T) {
	got, err := CharacterClass.Call("isUpperCase", "Z")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure Character.isLowerCase returns true for "z".
func TestCharacter_isLowerCase(t *testing.T) {
	got, err := CharacterClass.Call("isLowerCase", "z")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure Character.isWhitespace returns true for " ".
func TestCharacter_isWhitespace(t *testing.T) {
	got, err := CharacterClass.Call("isWhitespace", " ")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure Character.toUpperCase returns Character('A') for "a".
func TestCharacter_toUpperCase(t *testing.T) {
	got, err := CharacterClass.Call("toUpperCase", "a")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewCharacter('A') {
		t.Errorf("expected Character('A'), got %v", got)
	}
}

// Ensure Character.toLowerCase returns Character('a') for "A".
func TestCharacter_toLowerCase(t *testing.T) {
	got, err := CharacterClass.Call("toLowerCase", "A")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewCharacter('a') {
		t.Errorf("expected Character('a'), got %v", got)
	}
}

// Ensure Character.toString returns "A" for "A".
func TestCharacter_classToString(t *testing.T) {
	got, err := CharacterClass.Call("toString", "A")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "A" {
		t.Errorf("expected \"A\", got %v", got)
	}
}

// Ensure Character.compare returns 0 for equal chars.
func TestCharacter_classCompare(t *testing.T) {
	got, err := CharacterClass.Call("compare", "a", int64('a'))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 0 {
		t.Errorf("expected 0, got %v", got)
	}
}

// Ensure Character.MAX_VALUE returns Character(unicode.MaxRune).
func TestCharacter_maxValue(t *testing.T) {
	got, err := CharacterClass.Call("MAX_VALUE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != Character(unicode.MaxRune) {
		t.Errorf("expected MaxRune, got %v", got)
	}
}

// Ensure Character.MIN_VALUE returns Character(0).
func TestCharacter_minValue(t *testing.T) {
	got, err := CharacterClass.Call("MIN_VALUE")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != Character(0) {
		t.Errorf("expected Character(0), got %v", got)
	}
}

// Ensure NewCharacterFrom with empty string returns Character(0).
func TestCharacter_newFromEmptyString(t *testing.T) {
	c := NewCharacterFrom("")
	if c != Character(0) {
		t.Errorf("expected Character(0), got %v", c)
	}
}

// Ensure charValue on Character('Z') returns rune('Z').
func TestCharacter_instanceCharValue(t *testing.T) {
	c := NewCharacter('Z')
	got, err := c.Call("charValue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != rune('Z') {
		t.Errorf("expected rune('Z'), got %v", got)
	}
}

// Ensure isDigit on Character('9') returns true.
func TestCharacter_instanceIsDigit(t *testing.T) {
	c := NewCharacter('9')
	got, err := c.Call("isDigit")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure isLetter on Character('a') returns true.
func TestCharacter_instanceIsLetter(t *testing.T) {
	c := NewCharacter('a')
	got, err := c.Call("isLetter")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure isLetterOrDigit on Character('5') returns true.
func TestCharacter_instanceIsLetterOrDigit(t *testing.T) {
	c := NewCharacter('5')
	got, err := c.Call("isLetterOrDigit")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure isUpperCase on Character('Z') returns true.
func TestCharacter_instanceIsUpperCase(t *testing.T) {
	c := NewCharacter('Z')
	got, err := c.Call("isUpperCase")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure isLowerCase on Character('z') returns true.
func TestCharacter_instanceIsLowerCase(t *testing.T) {
	c := NewCharacter('z')
	got, err := c.Call("isLowerCase")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure isWhitespace on Character(' ') returns true.
func TestCharacter_instanceIsWhitespace(t *testing.T) {
	c := NewCharacter(' ')
	got, err := c.Call("isWhitespace")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure toUpperCase on Character('a') returns Character('A').
func TestCharacter_instanceToUpperCase(t *testing.T) {
	c := NewCharacter('a')
	got, err := c.Call("toUpperCase")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewCharacter('A') {
		t.Errorf("expected Character('A'), got %v", got)
	}
}

// Ensure toLowerCase on Character('A') returns Character('a').
func TestCharacter_instanceToLowerCase(t *testing.T) {
	c := NewCharacter('A')
	got, err := c.Call("toLowerCase")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != NewCharacter('a') {
		t.Errorf("expected Character('a'), got %v", got)
	}
}

// Ensure toString on Character('X') returns "X".
func TestCharacter_instanceToString(t *testing.T) {
	c := NewCharacter('X')
	got, err := c.Call("toString")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "X" {
		t.Errorf("expected \"X\", got %v", got)
	}
}

// Ensure compareTo returns 0 when receiver equals arg.
func TestCharacter_instanceCompareToEqual(t *testing.T) {
	c := NewCharacter('a')
	got, err := c.Call("compareTo", int64('a'))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 0 {
		t.Errorf("expected 0, got %v", got)
	}
}

// Ensure compareTo returns -1 when receiver is less than arg.
func TestCharacter_instanceCompareToLess(t *testing.T) {
	c := NewCharacter('a')
	got, err := c.Call("compareTo", int64('z'))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != -1 {
		t.Errorf("expected -1, got %v", got)
	}
}

// Ensure equals returns true when values match.
func TestCharacter_instanceEquals(t *testing.T) {
	c := NewCharacter('a')
	got, err := c.Call("equals", int64('a'))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure an unknown class method returns an error.
func TestCharacter_unknownClassMethod(t *testing.T) {
	if _, err := CharacterClass.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown class method")
	}
}

// Ensure an unknown instance method returns an error.
func TestCharacter_unknownInstanceMethod(t *testing.T) {
	c := NewCharacter('a')
	if _, err := c.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown instance method")
	}
}

// Ensure valueOf with wrong arg count returns an error.
func TestCharacter_valueOfArgCount(t *testing.T) {
	if _, err := CharacterClass.Call("valueOf"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure isDigit with wrong arg count returns an error.
func TestCharacter_isDigitArgCount(t *testing.T) {
	if _, err := CharacterClass.Call("isDigit"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure isLetter with wrong arg count returns an error.
func TestCharacter_isLetterArgCount(t *testing.T) {
	if _, err := CharacterClass.Call("isLetter"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure isLetterOrDigit with wrong arg count returns an error.
func TestCharacter_isLetterOrDigitArgCount(t *testing.T) {
	if _, err := CharacterClass.Call("isLetterOrDigit"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure isUpperCase with wrong arg count returns an error.
func TestCharacter_isUpperCaseArgCount(t *testing.T) {
	if _, err := CharacterClass.Call("isUpperCase"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure isLowerCase with wrong arg count returns an error.
func TestCharacter_isLowerCaseArgCount(t *testing.T) {
	if _, err := CharacterClass.Call("isLowerCase"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure isWhitespace with wrong arg count returns an error.
func TestCharacter_isWhitespaceArgCount(t *testing.T) {
	if _, err := CharacterClass.Call("isWhitespace"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure toUpperCase with wrong arg count returns an error.
func TestCharacter_toUpperCaseArgCount(t *testing.T) {
	if _, err := CharacterClass.Call("toUpperCase"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure toLowerCase with wrong arg count returns an error.
func TestCharacter_toLowerCaseArgCount(t *testing.T) {
	if _, err := CharacterClass.Call("toLowerCase"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure toString with wrong arg count returns an error.
func TestCharacter_classToStringArgCount(t *testing.T) {
	if _, err := CharacterClass.Call("toString"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure compare with wrong arg count returns an error.
func TestCharacter_classCompareArgCount(t *testing.T) {
	if _, err := CharacterClass.Call("compare", "a"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure instance compareTo with wrong arg count returns an error.
func TestCharacter_instanceCompareToArgCount(t *testing.T) {
	c := NewCharacter('a')
	if _, err := c.Call("compareTo"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure instance equals with wrong arg count returns an error.
func TestCharacter_instanceEqualsArgCount(t *testing.T) {
	c := NewCharacter('a')
	if _, err := c.Call("equals"); err == nil {
		t.Error("expected error for wrong arg count")
	}
}

// Ensure toInteger on a Character returns its Unicode code point as int32.
func TestCharacter_instanceToInteger(t *testing.T) {
	c := NewCharacter('A')
	got, err := c.Call("toInteger")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != int32(65) {
		t.Errorf("expected 65, got %v", got)
	}
}

// Ensure default on a non-null Character returns the character itself, not the fallback.
func TestCharacter_instanceDefault(t *testing.T) {
	c := NewCharacter('a')
	got, err := c.Call("default", rune('z'))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != c {
		t.Errorf("expected %v, got %v", c, got)
	}
}
