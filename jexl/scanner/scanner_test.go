// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package scanner

import (
	"io"
	"strings"
	"testing"

	"github.com/harness/go-jexl/jexl/token"
)

// scan drives the scanner to exhaustion and returns all tokens (including EOF).
// Returns (nil, err) if the scanner reports an error.
func scan(src string) ([]token.Token, error) {
	s := New()
	s.Reset(token.NewFile(src))
	var tokens []token.Token
	for {
		tok, err := s.Scan()
		if err == io.EOF {
			return tokens, nil
		}
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
		if tok.Kind == token.EOF {
			return tokens, nil
		}
	}
}

var scanTests = []struct {
	src  string
	want string
}{
	//
	// empty / whitespace
	//

	{"", "EOF"},
	{" ", "EOF"},
	{"\t\n\r", "EOF"},

	//
	// identifiers
	//

	{"foo", "IDENT(foo) EOF"},
	{"_x", "IDENT(_x) EOF"},
	{"$val", "IDENT($val) EOF"},
	{"abc123", "IDENT(abc123) EOF"},
	{"héllo", "IDENT(héllo) EOF"},
	{"a b c", "IDENT(a) IDENT(b) IDENT(c) EOF"},

	//
	// keywords → operator
	//

	{"in", "OP(in) EOF"},
	{"or", "OP(or) EOF"},
	{"and", "OP(and) EOF"},
	{"not", "OP(not) EOF"},
	{"not in", "OP(not) OP(in) EOF"},
	{"instanceof", "OP(instanceof) EOF"},
	{"!instanceof", "OP(!instanceof) EOF"},
	{"new", "OP(new) EOF"},
	{"if", "OP(if) EOF"},
	{"else", "OP(else) EOF"},
	{"for", "OP(for) EOF"},
	{"while", "OP(while) EOF"},
	{"do", "OP(do) EOF"},
	{"break", "OP(break) EOF"},
	{"continue", "OP(continue) EOF"},
	{"return", "OP(return) EOF"},
	{"var", "OP(var) EOF"},
	{"let", "OP(let) EOF"},
	{"const", "OP(const) EOF"},
	{"function", "OP(function) EOF"},
	{"empty", "OP(empty) EOF"},
	{"size", "OP(size) EOF"},
	{"try", "OP(try) EOF"},
	{"catch", "OP(catch) EOF"},
	{"throw", "OP(throw) EOF"},
	{"finally", "OP(finally) EOF"},
	{"switch", "OP(switch) EOF"},
	{"case", "OP(case) EOF"},
	{"default", "OP(default) EOF"},

	//
	// keyword aliases
	//

	{"eq", "OP(==) EOF"},
	{"ne", "OP(!=) EOF"},
	{"lt", "OP(<) EOF"},
	{"le", "OP(<=) EOF"},
	{"gt", "OP(>) EOF"},
	{"ge", "OP(>=) EOF"},
	{"mod", "OP(%) EOF"},

	//
	// identifiers that share a keyword prefix
	//

	{"nothing", "IDENT(nothing) EOF"},     // "not" prefix
	{"note", "IDENT(note) EOF"},           // "not" prefix
	{"inside", "IDENT(inside) EOF"},       // "in" prefix
	{"android", "IDENT(android) EOF"},     // "and" prefix
	{"oregon", "IDENT(oregon) EOF"},       // "or" prefix
	{"instances", "IDENT(instances) EOF"}, // "instanceof" prefix
	// "!instanceof" not matched when followed by alphanumeric
	{"!instanceofX", "OP(!) IDENT(instanceofX) EOF"},

	//
	// integer literals
	//

	{"0", "NUM(0) EOF"},
	{"42", "NUM(42) EOF"},
	{"1_000_000", "NUM(1_000_000) EOF"},
	{"0xff", "NUM(0xff) EOF"},
	{"0XFF", "NUM(0XFF) EOF"},
	{"0xDeAdBeEf", "NUM(0xDeAdBeEf) EOF"},
	{"0o77", "NUM(0o77) EOF"},
	{"0O77", "NUM(0O77) EOF"},
	{"0b1010", "NUM(0b1010) EOF"},
	{"0B1010", "NUM(0B1010) EOF"},

	//
	// integer suffixes
	//

	{"10L", "NUM(10L) EOF"},
	{"10l", "NUM(10l) EOF"},
	{"5h", "NUM(5h) EOF"},
	{"5H", "NUM(5H) EOF"},

	//
	// float literals
	//

	{"3.14", "NUM(3.14) EOF"},
	{".5", "NUM(.5) EOF"},
	{"1e10", "NUM(1e10) EOF"},
	{"1E10", "NUM(1E10) EOF"},
	{"1e+10", "NUM(1e+10) EOF"},
	{"1e-10", "NUM(1e-10) EOF"},
	{"1.5e3", "NUM(1.5e3) EOF"},

	//
	// float suffixes
	//

	{"1.5f", "NUM(1.5f) EOF"},
	{"1.5d", "NUM(1.5d) EOF"},
	{"2.5b", "NUM(2.5b) EOF"},
	{"2.5B", "NUM(2.5B) EOF"},

	//
	// range operator does not consume the second dot of ".."
	//

	{"1..5", "NUM(1) OP(..) NUM(5) EOF"},
	{"1 .. 5", "NUM(1) OP(..) NUM(5) EOF"},
	{"0..10", "NUM(0) OP(..) NUM(10) EOF"},

	//
	// string literals
	//

	{"'hello'", "STR(hello) EOF"},
	{`"world"`, "STR(world) EOF"},
	{"''", "STR() EOF"},
	{`""`, "STR() EOF"},
	{"'foo' 'bar'", "STR(foo) STR(bar) EOF"},

	//
	// string escape sequences
	//

	{`'a\nb'`, `STR(a\nb) EOF`},
	{`'a\tb'`, `STR(a\tb) EOF`},
	{`'a\rb'`, `STR(a\rb) EOF`},
	{`'a\\'`, `STR(a\) EOF`},
	{`'it\'s'`, "STR(it's) EOF"},
	{`"say \"hi\""`, `STR(say "hi") EOF`},
	{`'\a'`, `STR(\a) EOF`},
	{`'\b'`, `STR(\b) EOF`},
	{`'\f'`, `STR(\f) EOF`},
	{`'\v'`, `STR(\v) EOF`},
	{`'\x41'`, "STR(A) EOF"},
	{`'A'`, "STR(A) EOF"},
	{`'\u{41}'`, "STR(A) EOF"},
	{`'\101'`, "STR(A) EOF"},

	//
	// template literals
	//

	{"`hello`", "TMPL(hello) EOF"},
	{"``", "TMPL() EOF"},
	{"`line1\nline2`", `TMPL(line1\nline2) EOF`},
	{"`` `a``b` ``", "TMPL() TMPL(a`b) TMPL() EOF"}, // doubled backtick is an escaped backtick

	//
	// brackets
	//

	{"()", "BRK(() BRK()) EOF"},
	{"[]", "BRK([) BRK(]) EOF"},
	{"{}", "BRK({) BRK(}) EOF"},
	{"({[]})", "BRK(() BRK({) BRK([) BRK(]) BRK(}) BRK()) EOF"},

	//
	// arithmetic operators
	//

	{"+", "OP(+) EOF"},
	{"+=", "OP(+=) EOF"},
	{"++", "OP(++) EOF"},
	{"-", "OP(-) EOF"},
	{"-=", "OP(-=) EOF"},
	{"--", "OP(--) EOF"},
	{"->", "OP(->) EOF"},
	{"*", "OP(*) EOF"},
	{"*=", "OP(*=) EOF"},
	{"**", "OP(**) EOF"},
	{"%", "OP(%) EOF"},
	{"%=", "OP(%=) EOF"},

	//
	// division (not regex) after a value-producing token
	//

	{"1/2", "NUM(1) OP(/) NUM(2) EOF"},
	{"a/b", "IDENT(a) OP(/) IDENT(b) EOF"},
	{"x/=2", "IDENT(x) OP(/=) NUM(2) EOF"},
	{"(a)/b", "BRK(() IDENT(a) BRK()) OP(/) IDENT(b) EOF"},
	{"[1]/x", "BRK([) NUM(1) BRK(]) OP(/) IDENT(x) EOF"},

	//
	// comparison operators
	//

	{"==", "OP(==) EOF"},
	{"===", "OP(===) EOF"},
	{"!=", "OP(!=) EOF"},
	{"!==", "OP(!==) EOF"},
	{"<", "OP(<) EOF"},
	{"<=", "OP(<=) EOF"},
	{">", "OP(>) EOF"},
	{">=", "OP(>=) EOF"},

	//
	// logical operators
	//

	{"&&", "OP(&&) EOF"},
	{"||", "OP(||) EOF"},
	{"!", "OP(!) EOF"},

	//
	// bitwise operators
	//

	{"&", "OP(&) EOF"},
	{"&=", "OP(&=) EOF"},
	{"|", "OP(|) EOF"},
	{"|=", "OP(|=) EOF"},
	{"^", "OP(^) EOF"},
	{"^=", "OP(^=) EOF"},
	{"~", "OP(~) EOF"},
	{"<<", "OP(<<) EOF"},
	{"<<=", "OP(<<=) EOF"},
	{">>", "OP(>>) EOF"},
	{">>=", "OP(>>=) EOF"},
	{">>>", "OP(>>>) EOF"},
	{">>>=", "OP(>>>=) EOF"},

	//
	// string-match operators
	//

	{"=~", "OP(=~) EOF"},
	{"!~", "OP(!~) EOF"},
	{"=^", "OP(=^) EOF"},
	{"!^", "OP(!^) EOF"},
	{"=$", "OP(=$) EOF"},
	{"!$", "OP(!$) EOF"},

	//
	// assignment and fat-arrow
	//

	{"=", "OP(=) EOF"},
	{"=>", "OP(=>) EOF"},

	//
	// question-mark family
	//

	{"?", "OP(?) EOF"},
	{"?.", "OP(?.) EOF"},
	{"??", "OP(??) EOF"},
	// "?[" scans as two tokens: "?" and "["
	{"?[", "OP(?) BRK([) EOF"},

	//
	// colon variants
	//

	{":", "OP(:) EOF"},
	{"::", "OP(::) EOF"},

	//
	// dot variants
	//

	{".", "OP(.) EOF"},
	{"..", "OP(..) EOF"},

	//
	// comma and semicolon
	//

	{",", "OP(,) EOF"},
	{";", "OP(;) EOF"},

	//
	// hash / pointer
	//

	{"#", "OP(#) EOF"},
	{"#foo", "OP(#) IDENT(foo) EOF"},

	//
	// regex literals (value position = start of expression)
	//

	{"/foo/", "RX(foo/) EOF"},
	{"/foo/gi", "RX(foo/gi) EOF"},
	{`/fo\/o/`, `RX(fo\/o/) EOF`},
	{"(/x/)", "BRK(() RX(x/) BRK()) EOF"},
	{"1 + /x/", "NUM(1) OP(+) RX(x/) EOF"},
	{"x =~ /foo/", "IDENT(x) OP(=~) RX(foo/) EOF"},
	// after a string, "/" is division
	{"'s'/x/", "STR(s) OP(/) IDENT(x) OP(/) EOF"},

	//
	// comments
	//

	{"1 // ignored\n2", "NUM(1) NUM(2) EOF"},
	{"1 ## ignored\n2", "NUM(1) NUM(2) EOF"},
	{"1 // no newline", "NUM(1) EOF"},
	{"1 /* skip\nthis */ 2", "NUM(1) NUM(2) EOF"},
	{"a /* b */ c", "IDENT(a) IDENT(c) EOF"},

	//
	// full expression sequences
	//

	{"a + b * 2", "IDENT(a) OP(+) IDENT(b) OP(*) NUM(2) EOF"},
	{"foo(1, 2)", "IDENT(foo) BRK(() NUM(1) OP(,) NUM(2) BRK()) EOF"},
	{"a.b.c", "IDENT(a) OP(.) IDENT(b) OP(.) IDENT(c) EOF"},
	{"x ? 1 : 0", "IDENT(x) OP(?) NUM(1) OP(:) NUM(0) EOF"},
	{"1..10", "NUM(1) OP(..) NUM(10) EOF"},
	{"x -> x + 1", "IDENT(x) OP(->) IDENT(x) OP(+) NUM(1) EOF"},
	{"x => x * 2", "IDENT(x) OP(=>) IDENT(x) OP(*) NUM(2) EOF"},
	{"var x = 1", "OP(var) IDENT(x) OP(=) NUM(1) EOF"},
	{"[1, 2, 3]", "BRK([) NUM(1) OP(,) NUM(2) OP(,) NUM(3) BRK(]) EOF"},
	{`{"a": 1}`, `BRK({) STR(a) OP(:) NUM(1) BRK(}) EOF`},
	{"Math:abs(-1)", "IDENT(Math) OP(:) IDENT(abs) BRK(() OP(-) NUM(1) BRK()) EOF"},
	{"var x = 5h; ++x", "OP(var) IDENT(x) OP(=) NUM(5h) OP(;) OP(++) IDENT(x) EOF"},
	{"var x = 2.5b; ++x", "OP(var) IDENT(x) OP(=) NUM(2.5b) OP(;) OP(++) IDENT(x) EOF"},
	{"x not in [1,2]", "IDENT(x) OP(not) OP(in) BRK([) NUM(1) OP(,) NUM(2) BRK(]) EOF"},
	{"a >= 0 && a <= 10", "IDENT(a) OP(>=) NUM(0) OP(&&) IDENT(a) OP(<=) NUM(10) EOF"},
	{"x ?? y", "IDENT(x) OP(??) IDENT(y) EOF"},
	{"a?.b", "IDENT(a) OP(?.) IDENT(b) EOF"},
	{"s -> { var r = ''; r }", "IDENT(s) OP(->) BRK({) OP(var) IDENT(r) OP(=) STR() OP(;) IDENT(r) BRK(}) EOF"},
}

func TestScanner(t *testing.T) {
	for _, tt := range scanTests {
		t.Run(tt.src, func(t *testing.T) {
			tokens, err := scan(tt.src)
			if err != nil {
				t.Fatalf("scan error: %v", err)
			}
			if got := token.Print(tokens); got != tt.want {
				t.Errorf("\n got  %q\n want %q", got, tt.want)
			}
		})
	}
}

//
// Error cases
//

var errorTests = []struct {
	src     string
	message string
}{
	{`'hello`, "literal not terminated"},
	{`"hello`, "literal not terminated"},
	{"`hello", "literal not terminated"},
	{"/foo", "unterminated regex literal"},
	{"/foo\n/", "unterminated regex literal"},
	{"/* oops", "unclosed comment"},
	{"@", "unrecognized character"},
	{"0xGG", "bad number syntax"},
}

func TestScannerErrors(t *testing.T) {
	for _, tt := range errorTests {
		t.Run(tt.src, func(t *testing.T) {
			_, err := scan(tt.src)
			if err == nil {
				t.Fatalf("expected error containing %q, got none", tt.message)
			}
			if !strings.Contains(err.Error(), tt.message) {
				t.Errorf("error %q does not contain %q", err.Error(), tt.message)
			}
		})
	}
}

//
// Token ranges (rune offsets, not byte offsets)
//

var rangeTests = []struct {
	src   string
	index int // which token (0-based)
	from  int
	to    int
}{
	{"abc + 1", 0, 0, 3}, // abc
	{"abc + 1", 1, 4, 5}, // +
	{"abc + 1", 2, 6, 7}, // 1
	// unicode: "héllo" is 5 runes but 6 bytes; ranges must be rune-based
	{"héllo + 1", 0, 0, 5}, // héllo
	{"héllo + 1", 1, 6, 7}, // +
	{"héllo + 1", 2, 8, 9}, // 1
}

func TestTokenRanges(t *testing.T) {
	for _, tt := range rangeTests {
		t.Run(tt.src, func(t *testing.T) {
			s := New()
			s.Reset(token.NewFile(tt.src))
			var tok token.Token
			for i := 0; i <= tt.index; i++ {
				var err error
				tok, err = s.Scan()
				if err != nil {
					t.Fatalf("scan error at token %d: %v", i, err)
				}
			}
			if tok.Range.From != tt.from || tok.Range.To != tt.to {
				t.Errorf("token[%d] range: got [%d,%d) want [%d,%d)",
					tt.index, tok.Range.From, tok.Range.To, tt.from, tt.to)
			}
		})
	}
}

//
// Scanner reuse via reset
//

func TestScannerReset(t *testing.T) {
	s := New()

	s.Reset(token.NewFile("1"))
	tok, _ := s.Scan()
	if tok.Kind != token.Number || tok.Value != "1" {
		t.Fatalf("first scan: got %v", tok)
	}

	// Reset to a different (shorter) source — cursors must be zeroed.
	s.Reset(token.NewFile("foo"))
	tok, _ = s.Scan()
	if tok.Kind != token.Ident || tok.Value != "foo" {
		t.Fatalf("after reset: got %v", tok)
	}

	// Reset to empty source.
	s.Reset(token.NewFile(""))
	tok, _ = s.Scan()
	if tok.Kind != token.EOF {
		t.Fatalf("after reset to empty: got %v", tok)
	}
}
