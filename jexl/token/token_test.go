// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package token

import (
	"testing"
)

// Ensure Token.String returns the kind when value is empty.
func TestToken_String_noValue(t *testing.T) {
	tok := Token{Kind: Ident}
	if tok.String() != "Identifier" {
		t.Fatalf("expected Identifier, got %q", tok.String())
	}
}

// Ensure Token.String includes the value when present.
func TestToken_String_withValue(t *testing.T) {
	tok := Token{Kind: Ident, Value: "foo"}
	if tok.String() != `Identifier("foo")` {
		t.Fatalf("unexpected: %q", tok.String())
	}
}

// Ensure Token.Is matches by kind only when no values given.
func TestToken_Is_kindOnly(t *testing.T) {
	tok := Token{Kind: Number, Value: "42"}
	if !tok.Is(Number) {
		t.Fatal("expected match on kind")
	}
	if tok.Is(Ident) {
		t.Fatal("expected no match on wrong kind")
	}
}

// Ensure Token.Is matches when the value is in the list.
func TestToken_Is_withValues(t *testing.T) {
	tok := Token{Kind: Operator, Value: "+"}
	if !tok.Is(Operator, "+", "-") {
		t.Fatal("expected match")
	}
	if tok.Is(Operator, "-", "*") {
		t.Fatal("expected no match for unlisted value")
	}
}

// Ensure Token.Is returns false when kind matches but value list does not.
func TestToken_Is_kindMatchValueMiss(t *testing.T) {
	tok := Token{Kind: Bracket, Value: "("}
	if tok.Is(Bracket, ")", "]") {
		t.Fatal("expected no match")
	}
}

// Ensure Token embeds Range correctly.
func TestToken_Location(t *testing.T) {
	tok := Token{
		Range: Range{From: 3, To: 7},
	}
	if tok.From != 3 || tok.To != 7 {
		t.Fatalf("unexpected location: %+v", tok.Range)
	}
}

// Ensure IsAlphabetic accepts letters, underscore, and dollar sign.
func TestIsAlphabetic_valid(t *testing.T) {
	for _, r := range []rune{'a', 'Z', '_', '$', 'é'} {
		if !IsAlphabetic(r) {
			t.Fatalf("expected %q to be alphabetic", r)
		}
	}
}

// Ensure IsAlphabetic rejects digits and punctuation.
func TestIsAlphabetic_invalid(t *testing.T) {
	for _, r := range []rune{'0', '9', ' ', '-', '!'} {
		if IsAlphabetic(r) {
			t.Fatalf("expected %q to not be alphabetic", r)
		}
	}
}

// Ensure IsAlphaNumeric accepts letters, digits, underscore, and dollar.
func TestIsAlphaNumeric_valid(t *testing.T) {
	for _, r := range []rune{'a', 'Z', '0', '9', '_', '$'} {
		if !IsAlphaNumeric(r) {
			t.Fatalf("expected %q to be alphanumeric", r)
		}
	}
}

// Ensure IsAlphaNumeric rejects punctuation and whitespace.
func TestIsAlphaNumeric_invalid(t *testing.T) {
	for _, r := range []rune{' ', '+', '.'} {
		if IsAlphaNumeric(r) {
			t.Fatalf("expected %q to not be alphanumeric", r)
		}
	}
}

// Ensure IsSpace accepts all ASCII whitespace variants.
func TestIsSpace(t *testing.T) {
	for _, r := range []rune{' ', '\t', '\n', '\r'} {
		if !IsSpace(r) {
			t.Fatalf("expected %q to be space", r)
		}
	}
	if IsSpace('a') {
		t.Fatal("expected 'a' to not be space")
	}
}

// Ensure IsQuote recognises single and double quotes only.
func TestIsQuote(t *testing.T) {
	if !IsQuote('\'') || !IsQuote('"') {
		t.Fatal("expected quotes to be recognised")
	}
	if IsQuote('`') || IsQuote('a') {
		t.Fatal("expected backtick and letter to not be quotes")
	}
}

// Ensure IsNumber accepts decimal digits only.
func TestIsNumber(t *testing.T) {
	for _, r := range []rune{'0', '5', '9'} {
		if !IsNumber(r) {
			t.Fatalf("expected %q to be a number", r)
		}
	}
	if IsNumber('a') || IsNumber(' ') {
		t.Fatal("expected letter/space to not be a number")
	}
}

// Ensure IsValidIdentifier accepts well-formed identifiers.
func TestIsValidIdentifier_valid(t *testing.T) {
	for _, s := range []string{"foo", "_bar", "$baz", "x1", "camelCase"} {
		if !IsValidIdentifier(s) {
			t.Fatalf("expected %q to be a valid identifier", s)
		}
	}
}

// Ensure IsValidIdentifier rejects empty strings and identifiers starting with digits.
func TestIsValidIdentifier_invalid(t *testing.T) {
	for _, s := range []string{"", "1foo", "foo bar", "foo-bar"} {
		if IsValidIdentifier(s) {
			t.Fatalf("expected %q to be invalid", s)
		}
	}
}

// Ensure IsAssignment recognises all assignment operators.
func TestIsAssignment(t *testing.T) {
	for _, op := range []string{"=", "+=", "-=", "*=", "/=", "%=", "&=", "|=", "^=", "<<=", ">>=", ">>>="} {
		if !IsAssignment(op) {
			t.Fatalf("expected %q to be an assignment operator", op)
		}
	}
	if IsAssignment("+") || IsAssignment("==") {
		t.Fatal("expected + and == to not be assignment operators")
	}
}

// Ensure IsBoolean recognises the four boolean operators.
func TestIsBoolean(t *testing.T) {
	for _, op := range []string{"and", "or", "&&", "||"} {
		if !IsBoolean(op) {
			t.Fatalf("expected %q to be boolean", op)
		}
	}
	if IsBoolean("+") || IsBoolean("==") {
		t.Fatal("expected + and == to not be boolean")
	}
}

// Ensure IsComparison recognises the four ordering operators.
func TestIsComparison(t *testing.T) {
	for _, op := range []string{"<", ">", "<=", ">="} {
		if !IsComparison(op) {
			t.Fatalf("expected %q to be a comparison", op)
		}
	}
	if IsComparison("==") || IsComparison("!=") {
		t.Fatal("expected == and != to not be comparisons")
	}
}

// Ensure all expected unary operators are present in the table.
func TestUnary_operators(t *testing.T) {
	for _, op := range []string{"not", "!", "-", "+", "~"} {
		if _, ok := Unary[op]; !ok {
			t.Fatalf("unary table missing %q", op)
		}
	}
}

// Ensure all expected binary operators are present in the table.
func TestBinary_operators(t *testing.T) {
	for _, op := range []string{"or", "and", "==", "!=", "<", ">", "in", "+", "-", "*", "/", "%", "**", "??"} {
		if _, ok := Binary[op]; !ok {
			t.Fatalf("binary table missing %q", op)
		}
	}
}

// Ensure ** has higher precedence than *.
func TestBinary_precedence(t *testing.T) {
	if Binary["**"].Precedence <= Binary["*"].Precedence {
		t.Fatal("expected ** to have higher precedence than *")
	}
}

// Ensure ** is right-associative and + is left-associative.
func TestBinary_associativity(t *testing.T) {
	if Binary["**"].Associativity != Right {
		t.Fatal("expected ** to be right-associative")
	}
	if Binary["+"].Associativity != Left {
		t.Fatal("expected + to be left-associative")
	}
}
