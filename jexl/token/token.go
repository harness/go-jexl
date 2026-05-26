// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package token

import (
	"fmt"
	"slices"
	"unicode"
)

// Kind is the lexical category of a token.
type Kind string

const (
	Ident    Kind = "Identifier"
	Number   Kind = "Number"
	String   Kind = "String"
	Template Kind = "Template"
	Operator Kind = "Operator"
	Bracket  Kind = "Bracket"
	Regex    Kind = "Regex"
	EOF      Kind = "EOF"
)

// Token is a lexical token produced by the scanner.
type Token struct {
	Range
	Kind  Kind
	Value string
}

// String implements fmt.Stringer.
func (t Token) String() string {
	if t.Value == "" {
		return string(t.Kind)
	}
	return fmt.Sprintf("%s(%#v)", t.Kind, t.Value)
}

// Is reports whether t matches kind k and, if values are
// given, whether t.Value is one of them.
func (t Token) Is(k Kind, values ...string) bool {
	if k != t.Kind {
		return false
	}
	if slices.Contains(values, t.Value) {
		return true
	}
	return len(values) == 0
}

// Associativity of a binary operator.
type Associativity int

const (
	Left Associativity = iota + 1
	Right
)

// Op holds the precedence and associativity of
// an operator.
type Op struct {
	Precedence    int
	Associativity Associativity
}

// Unary is the table of unary operators.
var Unary = map[string]Op{
	"not": {50, Left},
	"!":   {50, Left},
	"-":   {90, Left},
	"+":   {90, Left},
	"~":   {90, Left},
}

// Binary is the table of binary operators.
var Binary = map[string]Op{
	"|":           {22, Left},
	"or":          {10, Left},
	"||":          {10, Left},
	"and":         {15, Left},
	"&&":          {15, Left},
	"==":          {20, Left},
	"!=":          {20, Left},
	"===":         {20, Left},
	"!==":         {20, Left},
	"<":           {20, Left},
	">":           {20, Left},
	">=":          {20, Left},
	"<=":          {20, Left},
	"in":          {20, Left},
	"instanceof":  {20, Left},
	"!instanceof": {20, Left},
	"=~":          {20, Left},
	"!~":          {20, Left},
	"=^":          {20, Left},
	"!^":          {20, Left},
	"=$":          {20, Left},
	"!$":          {20, Left},
	"..":          {25, Left},
	"+":           {30, Left},
	"-":           {30, Left},
	"*":           {60, Left},
	"/":           {60, Left},
	"%":           {60, Left},
	"**":          {100, Right},
	"^":           {24, Left},
	"&":           {26, Left},
	"<<":          {28, Left},
	">>":          {28, Left},
	">>>":         {28, Left},
	"??":          {500, Left},
	"=":           {1, Right},
	"+=":          {1, Right},
	"-=":          {1, Right},
	"*=":          {1, Right},
	"/=":          {1, Right},
	"%=":          {1, Right},
	"&=":          {1, Right},
	"|=":          {1, Right},
	"^=":          {1, Right},
	"<<=":         {1, Right},
	">>=":         {1, Right},
	">>>=":        {1, Right},
}

// IsAssignment reports whether op is an assignment operator.
func IsAssignment(op string) bool {
	switch op {
	case "=", "+=", "-=", "*=", "/=", "%=", "&=", "|=", "^=", "<<=", ">>=", ">>>=":
		return true
	}
	return false
}

// IsBoolean reports whether op is a boolean operator.
func IsBoolean(op string) bool {
	return op == "and" || op == "or" || op == "&&" || op == "||"
}

// IsComparison reports whether op is an ordering comparison operator.
func IsComparison(op string) bool {
	return op == "<" || op == ">" || op == ">=" || op == "<="
}

// IsValidIdentifier reports whether s is a valid JEXL identifier.
func IsValidIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !IsAlphabetic(r) {
				return false
			}
		} else if !IsAlphaNumeric(r) {
			return false
		}
	}
	return true
}

// IsAlphabetic reports whether r can start an identifier.
func IsAlphabetic(r rune) bool {
	return r == '_' || r == '$' || unicode.IsLetter(r)
}

// IsAlphaNumeric reports whether r can appear inside an identifier.
func IsAlphaNumeric(r rune) bool {
	return IsAlphabetic(r) || unicode.IsDigit(r)
}

// IsSpace reports whether r is ASCII whitespace.
func IsSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

// IsQuote reports whether r is a string-literal quote character.
func IsQuote(r rune) bool {
	return r == '\'' || r == '"'
}

// IsNumber reports whether r is an ASCII decimal digit.
func IsNumber(r rune) bool {
	return r >= '0' && r <= '9'
}
