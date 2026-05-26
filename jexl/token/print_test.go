// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package token

import "testing"

var printTests = []struct {
	tokens []Token
	want   string
}{
	//
	// empty input
	//

	{nil, ""},
	{[]Token{}, ""},

	//
	// each Kind
	//

	{[]Token{{Kind: EOF}}, "EOF"},
	{[]Token{{Kind: Ident, Value: "foo"}}, "IDENT(foo)"},
	{[]Token{{Kind: Number, Value: "42"}}, "NUM(42)"},
	{[]Token{{Kind: String, Value: "hello"}}, "STR(hello)"},
	{[]Token{{Kind: Template, Value: "hi"}}, "TMPL(hi)"},
	{[]Token{{Kind: Operator, Value: "+"}}, "OP(+)"},
	{[]Token{{Kind: Bracket, Value: "("}}, "BRK(()"},
	{[]Token{{Kind: Regex, Value: "foo/"}}, "RX(foo/)"},

	//
	// unknown Kind falls through to default
	//

	{[]Token{{Kind: "Unknown", Value: "x"}}, "?(x)"},

	//
	// multiple tokens separated by spaces
	//

	{[]Token{
		{Kind: Ident, Value: "a"},
		{Kind: Operator, Value: "+"},
		{Kind: Number, Value: "1"},
		{Kind: EOF},
	}, "IDENT(a) OP(+) NUM(1) EOF"},

	//
	// control characters in String values are escaped
	//

	{[]Token{{Kind: String, Value: "\a\b\f\n\r\t\v"}}, `STR(\a\b\f\n\r\t\v)`},
	{[]Token{{Kind: Template, Value: "line1\nline2"}}, `TMPL(line1\nline2)`},
}

func TestPrint(t *testing.T) {
	for _, tt := range printTests {
		got := Print(tt.tokens)
		if got != tt.want {
			t.Errorf("Print(%v)\n got  %q\n want %q", tt.tokens, got, tt.want)
		}
	}
}
