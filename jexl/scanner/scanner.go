// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package scanner

import (
	"fmt"
	"io"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/harness/go-jexl/jexl/internal/ring"
	"github.com/harness/go-jexl/jexl/token"
)

const ringChunkSize = 10

const eof rune = -1

// New returns a new Scanner ready to use.
func New() *Scanner {
	return &Scanner{
		tokens: ring.New[token.Token](ringChunkSize),
	}
}

// Scanner tokenises a JEXL source string.
type Scanner struct {
	state      stateFn
	source     token.File
	tokens     *ring.Ring[token.Token]
	err        *token.Error
	start, end struct {
		byte, rune int
	}
	eof           bool
	valuePosition bool
}

// Reset prepares the Scanner to tokenise source.
func (l *Scanner) Reset(source token.File) {
	l.source = source
	l.tokens.Reset()
	l.state = root
	l.valuePosition = true
	l.start.byte = 0
	l.start.rune = 0
	l.end.byte = 0
	l.end.rune = 0
	l.eof = false
	l.err = nil
}

// Scan returns the next token. It returns io.EOF when the source is exhausted.
func (l *Scanner) Scan() (token.Token, error) {
	for l.state != nil && l.err == nil && l.tokens.Len() == 0 {
		l.state = l.state(l)
	}
	if l.err != nil {
		return token.Token{}, l.err.Bind(l.source)
	}
	if t, ok := l.tokens.Dequeue(); ok {
		return t, nil
	}
	return token.Token{}, io.EOF
}

func (l *Scanner) commit() {
	l.start = l.end
}

func (l *Scanner) next() rune {
	if l.end.byte >= len(l.source.String()) {
		l.eof = true
		return eof
	}
	r, sz := utf8.DecodeRuneInString(l.source.String()[l.end.byte:])
	l.end.rune++
	l.end.byte += sz
	return r
}

func (l *Scanner) peek() rune {
	if l.end.byte < len(l.source.String()) {
		r, _ := utf8.DecodeRuneInString(l.source.String()[l.end.byte:])
		return r
	}
	return eof
}

func (l *Scanner) backup() {
	if l.eof {
		l.eof = false
	} else if l.end.rune > 0 {
		_, sz := utf8.DecodeLastRuneInString(l.source.String()[:l.end.byte])
		l.end.byte -= sz
		l.end.rune--
	}
}

func (l *Scanner) emit(t token.Kind) {
	l.emitValue(t, l.word())
}

func (l *Scanner) emitValue(t token.Kind, value string) {
	l.tokens.Enqueue(token.Token{
		Range: token.Range{From: l.start.rune, To: l.end.rune},
		Kind:  t,
		Value: value,
	})
	l.commit()
	switch t {
	case token.Ident, token.Number, token.String, token.Regex:
		l.valuePosition = false
	default:
		if t == token.Bracket && (value == ")" || value == "]" || value == "}") {
			l.valuePosition = false
		} else {
			l.valuePosition = true
		}
	}
}

func (l *Scanner) emitEOF() {
	from := l.end.rune - 1
	if from < 0 {
		from = 0
	}
	to := l.end.rune - 0
	if to < 0 {
		to = 0
	}
	l.tokens.Enqueue(token.Token{
		Range: token.Range{From: from, To: to},
		Kind:  token.EOF,
	})
	l.commit()
}

func (l *Scanner) skip() {
	l.commit()
}

func (l *Scanner) word() string {
	return l.source.String()[l.start.byte:l.end.byte]
}

func (l *Scanner) accept(valid string) bool {
	if strings.ContainsRune(valid, l.peek()) {
		l.next()
		return true
	}
	return false
}

func (l *Scanner) acceptRun(valid string) {
	for l.accept(valid) {
	}
}

func (l *Scanner) skipSpaces() {
	l.acceptRun(" ")
	l.skip()
}

// peekNonSpace returns the first non-space rune after the current position
// without consuming any input. Used to detect call context (followed by '(').
func (l *Scanner) peekNonSpace() rune {
	s := l.source.String()[l.end.byte:]
	for _, r := range s {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			return r
		}
	}
	return eof
}

func (l *Scanner) error(format string, args ...any) stateFn {
	if l.err == nil {
		end := l.end.rune
		if l.eof {
			end++
		}
		l.err = &token.Error{
			Range: token.Range{
				From: end - 1,
				To:   end,
			},
			Message: fmt.Sprintf(format, args...),
		}
	}
	return nil
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= lower(ch) && lower(ch) <= 'f':
		return int(lower(ch) - 'a' + 10)
	}
	return 16
}

func lower(ch rune) rune { return ('a' - 'A') | ch }

func (l *Scanner) scanDigits(ch rune, base, n int) rune {
	for n > 0 && digitVal(ch) < base {
		ch = l.next()
		n--
	}
	if n > 0 {
		l.error("invalid char escape")
	}
	return ch
}

func (l *Scanner) scanEscape(quote rune) rune {
	ch := l.next()
	switch ch {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		ch = l.next()
	case '0', '1', '2', '3', '4', '5', '6', '7':
		ch = l.scanDigits(ch, 8, 3)
	case 'x':
		ch = l.scanDigits(l.next(), 16, 2)
	case 'u':
		if l.peek() == '{' {
			l.next()
			digits := 0
			for {
				p := l.peek()
				if p == '}' {
					break
				}
				if digitVal(p) >= 16 {
					l.error("invalid char escape")
					return eof
				}
				if digits >= 6 {
					l.error("invalid char escape")
					return eof
				}
				l.next()
				digits++
			}
			if l.peek() != '}' || digits == 0 {
				l.error("invalid char escape")
				return eof
			}
			l.next()
			ch = l.next()
			break
		}
		ch = l.scanDigits(l.next(), 16, 4)
	case 'U':
		ch = l.scanDigits(l.next(), 16, 8)
	default:
		ch = l.next()
	}
	return ch
}

func (l *Scanner) scanString(quote rune) (n int) {
	ch := l.next()
	for ch != quote {
		if ch == '\n' || ch == eof {
			l.error("literal not terminated")
			return
		}
		if ch == '\\' {
			ch = l.scanEscape(quote)
		} else {
			ch = l.next()
		}
		n++
	}
	return
}

func (l *Scanner) scanRawString(quote rune) (n int) {
	var escapedQuotes int
loop:
	for {
		ch := l.next()
		for ch == quote && l.peek() == quote {
			l.next()
			ch = l.next()
			escapedQuotes++
		}
		switch ch {
		case quote:
			break loop
		case eof:
			l.error("literal not terminated")
			return
		}
		n++
	}
	str := l.source.String()[l.start.byte+1 : l.end.byte-1]

	if escapedQuotes == 0 {
		l.emitValue(token.Template, str)
		return
	}

	var b strings.Builder
	var skipped bool
	b.Grow(len(str) - escapedQuotes)
	for _, r := range str {
		if r == quote {
			if !skipped {
				skipped = true
				continue
			}
			skipped = false
		}
		b.WriteRune(r)
	}
	l.emitValue(token.Template, b.String())
	return
}

type stateFn func(*Scanner) stateFn

func root(l *Scanner) stateFn {
	switch r := l.next(); {
	case r == eof:
		l.emitEOF()
		return nil
	case token.IsSpace(r):
		l.skip()
		return root
	case token.IsQuote(r):
		l.scanString(r)
		str, err := unescape(l.word())
		if err != nil {
			l.error("%v", err)
		}
		l.emitValue(token.String, str)
	case r == '`':
		l.scanRawString(r)
	case token.IsNumber(r):
		l.backup()
		return number
	case r == '?':
		return questionMark
	case r == '/':
		return slash
	case r == '#':
		if l.accept("#") {
			return singleLineComment
		}
		return pointer
	case r == '|':
		if l.accept("|") {
			l.emit(token.Operator) // ||
		} else {
			l.accept("=")
			l.emit(token.Operator) // | or |=
		}
	case r == ':':
		l.accept(":")
		l.emit(token.Operator)
	case strings.ContainsRune("([{", r):
		l.emit(token.Bracket)
	case strings.ContainsRune(")]}", r):
		l.emit(token.Bracket)
	case strings.ContainsRune(",;~", r):
		l.emit(token.Operator)
	case r == '+':
		if l.accept("+") {
			l.emit(token.Operator) // ++
		} else {
			l.accept("=")
			l.emit(token.Operator) // + or +=
		}
	case r == '-':
		if l.accept(">") {
			l.emit(token.Operator) // ->
		} else if l.accept("-") {
			l.emit(token.Operator) // --
		} else {
			l.accept("=")
			l.emit(token.Operator) // - or -=
		}
	case r == '*':
		if l.accept("*") {
			l.emit(token.Operator) // **
		} else {
			l.accept("=")
			l.emit(token.Operator) // * or *=
		}
	case r == '%':
		l.accept("=")
		l.emit(token.Operator) // % or %=
	case r == '=':
		if l.accept(">") {
			l.emit(token.Operator) // =>
		} else if l.accept("~") {
			l.emit(token.Operator) // =~
		} else if l.accept("^") {
			l.emit(token.Operator) // =^
		} else if l.accept("$") {
			l.emit(token.Operator) // =$
		} else {
			l.accept("=")
			l.accept("=") // ===
			l.emit(token.Operator)
		}
	case r == '!':
		if l.accept("~") {
			l.emit(token.Operator) // !~
		} else if l.accept("^") {
			l.emit(token.Operator) // !^
		} else if l.accept("$") {
			l.emit(token.Operator) // !$
		} else if strings.HasPrefix(l.source.String()[l.end.byte:], "instanceof") {
			rest := l.source.String()[l.end.byte+len("instanceof"):]
			var nextRune rune
			if len(rest) > 0 {
				nextRune, _ = utf8.DecodeRuneInString(rest)
			}
			if len(rest) == 0 || !token.IsAlphaNumeric(nextRune) {
				for range "instanceof" {
					l.next()
				}
				l.emit(token.Operator) // !instanceof
			} else {
				l.accept("=")
				l.accept("=") // !==
				l.emit(token.Operator)
			}
		} else {
			l.accept("=")
			l.accept("=") // !==
			l.emit(token.Operator)
		}
	case r == '&':
		if l.accept("&") {
			l.emit(token.Operator) // &&
		} else {
			l.accept("=")
			l.emit(token.Operator) // & or &=
		}
	case r == '^':
		l.accept("=")
		l.emit(token.Operator) // ^ or ^=
	case r == '<':
		if l.accept("<") {
			l.accept("=")
			l.emit(token.Operator) // << or <<=
		} else {
			l.accept("=")
			l.emit(token.Operator) // < or <=
		}
	case r == '>':
		if l.accept(">") {
			if l.accept(">") {
				l.accept("=")
				l.emit(token.Operator) // >>> or >>>=
			} else {
				l.accept("=")
				l.emit(token.Operator) // >> or >>=
			}
		} else {
			l.accept("=")
			l.emit(token.Operator) // > or >=
		}
	case r == '.':
		l.backup()
		return dot
	case token.IsAlphaNumeric(r):
		l.backup()
		return identifier
	default:
		return l.error("unrecognized character: %#U", r)
	}
	return root
}

func number(l *Scanner) stateFn {
	if !l.scanNumber() {
		return l.error("bad number syntax: %q", l.word())
	}
	l.emit(token.Number)
	return root
}

func (l *Scanner) scanNumber() bool {
	digits := "0123456789_"
	if l.accept("0") {
		if l.accept("xX") {
			digits = "0123456789abcdefABCDEF_"
		} else if l.accept("oO") {
			digits = "01234567_"
		} else if l.accept("bB") {
			digits = "01_"
		}
	}
	l.acceptRun(digits)
	end := l.end
	if l.accept(".") {
		if l.peek() == '.' {
			l.end = end
			return true
		}
		l.acceptRun(digits)
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun(digits)
	}
	l.accept("LlHhBbfd")
	if token.IsAlphaNumeric(l.peek()) {
		l.next()
		return false
	}
	return true
}

func dot(l *Scanner) stateFn {
	l.next()
	if l.accept("0123456789") {
		l.backup()
		return number
	}
	l.accept(".")
	l.emit(token.Operator)
	return root
}

func identifier(l *Scanner) stateFn {
loop:
	for {
		switch r := l.next(); {
		case token.IsAlphaNumeric(r):
			// absorb
		default:
			l.backup()
			switch l.word() {
			case "not":
				return not
			case "mod":
				l.emitValue(token.Operator, "%")
			case "div":
				l.emitValue(token.Operator, "/")
			case "eq":
				l.emitValue(token.Operator, "==")
			case "ne":
				l.emitValue(token.Operator, "!=")
			case "lt":
				l.emitValue(token.Operator, "<")
			case "le":
				l.emitValue(token.Operator, "<=")
			case "gt":
				l.emitValue(token.Operator, ">")
			case "ge":
				l.emitValue(token.Operator, ">=")
			case "in", "or", "and", "instanceof", "let", "var",
				"const", "empty", "size":
				l.emit(token.Operator)
			case "default":
				// Treat as an identifier when immediately followed by '(' so
				// that default(...) can be called as a builtin function.
				if l.peekNonSpace() == '(' {
					l.emit(token.Ident)
				} else {
					l.emit(token.Operator)
				}
			case "if", "else", "for", "while", "do",
				"break", "continue", "return", "try", "catch",
				"throw", "finally", "switch", "case", "new", "function":
				l.emit(token.Operator)
			default:
				l.emit(token.Ident)
			}
			break loop
		}
	}
	return root
}

func not(l *Scanner) stateFn {
	l.emit(token.Operator)

	l.skipSpaces()

	end := l.end

	for {
		r := l.next()
		if token.IsAlphaNumeric(r) {
			// absorb
		} else {
			l.backup()
			break
		}
	}

	switch l.word() {
	case "in":
		l.emit(token.Operator)
	default:
		l.end = end
	}
	return root
}

func questionMark(l *Scanner) stateFn {
	l.accept(".?")
	l.emit(token.Operator)
	return root
}

func slash(l *Scanner) stateFn {
	if l.accept("/") {
		return singleLineComment
	}
	if l.accept("*") {
		return multiLineComment
	}
	if l.valuePosition {
		return regexLiteral
	}
	l.accept("=")
	l.emit(token.Operator)
	return root
}

func regexLiteral(l *Scanner) stateFn {
	for {
		r := l.next()
		switch r {
		case eof, '\n':
			return l.error("unterminated regex literal")
		case '\\':
			if l.next() == eof {
				return l.error("unterminated regex literal")
			}
		case '/':
			for {
				p := l.peek()
				if (p >= 'a' && p <= 'z') || (p >= 'A' && p <= 'Z') {
					l.next()
				} else {
					break
				}
			}
			raw := l.word()
			l.emitValue(token.Regex, raw[1:])
			return root
		}
	}
}

func singleLineComment(l *Scanner) stateFn {
	for {
		r := l.next()
		if r == eof || r == '\n' {
			break
		}
	}
	l.skip()
	return root
}

func multiLineComment(l *Scanner) stateFn {
	for {
		r := l.next()
		if r == eof {
			return l.error("unclosed comment")
		}
		if r == '*' && l.accept("/") {
			break
		}
	}
	l.skip()
	return root
}

func pointer(l *Scanner) stateFn {
	l.accept("#")
	l.emit(token.Operator)
	for {
		switch r := l.next(); {
		case token.IsAlphaNumeric(r):
			// absorb
		default:
			l.backup()
			if l.word() != "" {
				l.emit(token.Ident)
			}
			return root
		}
	}
}

var newlineNormalizer = strings.NewReplacer("\r\n", "\n", "\r", "\n")

func unescape(value string) (string, error) {
	value = newlineNormalizer.Replace(value)
	n := len(value)

	if n < 2 {
		return value, fmt.Errorf("unable to unescape string")
	}

	if value[0] != value[n-1] || (value[0] != '"' && value[0] != '\'') {
		return value, fmt.Errorf("unable to unescape string")
	}

	value = value[1 : n-1]

	var runeTmp [utf8.UTFMax]byte
	size := 3 * uint64(n) / 2
	if size >= math.MaxInt {
		return "", fmt.Errorf("too large string")
	}
	buf := new(strings.Builder)
	buf.Grow(int(size))
	for len(value) > 0 {
		c, multibyte, rest, err := unescapeChar(value)
		if err != nil {
			return "", err
		}
		value = rest
		if c < utf8.RuneSelf || !multibyte {
			buf.WriteByte(byte(c))
		} else {
			n := utf8.EncodeRune(runeTmp[:], c)
			buf.Write(runeTmp[:n])
		}
	}
	return buf.String(), nil
}

func unescapeChar(s string) (value rune, multibyte bool, tail string, err error) {
	switch c := s[0]; {
	case c >= utf8.RuneSelf:
		r, size := utf8.DecodeRuneInString(s)
		return r, true, s[size:], nil
	case c != '\\':
		return rune(s[0]), false, s[1:], nil
	}

	if len(s) <= 1 {
		err = fmt.Errorf("unable to unescape string, found '\\' as last character")
		return
	}

	c := s[1]
	s = s[2:]
	switch c {
	case 'a':
		value = '\a'
	case 'b':
		value = '\b'
	case 'f':
		value = '\f'
	case 'n':
		value = '\n'
	case 'r':
		value = '\r'
	case 't':
		value = '\t'
	case 'v':
		value = '\v'
	case '\\':
		value = '\\'
	case '\'':
		value = '\''
	case '"':
		value = '"'
	case '`':
		value = '`'
	case '?':
		value = '?'

	case 'x', 'X', 'u', 'U':
		if c == 'u' && len(s) > 0 && s[0] == '{' {
			s = s[1:]
			var v rune
			digits := 0
			for len(s) > 0 && s[0] != '}' {
				x, ok := unhex(s[0])
				if !ok {
					err = fmt.Errorf("unable to unescape string")
					return
				}
				if digits >= 6 {
					err = fmt.Errorf("unable to unescape string")
					return
				}
				v = v<<4 | x
				s = s[1:]
				digits++
			}
			if len(s) == 0 || s[0] != '}' || digits == 0 {
				err = fmt.Errorf("unable to unescape string")
				return
			}
			s = s[1:]
			if v > utf8.MaxRune {
				err = fmt.Errorf("unable to unescape string")
				return
			}
			value = v
			multibyte = true
			break
		}
		n := 0
		switch c {
		case 'x', 'X':
			n = 2
		case 'u':
			n = 4
		case 'U':
			n = 8
		}
		var v rune
		if len(s) < n {
			err = fmt.Errorf("unable to unescape string")
			return
		}
		for j := 0; j < n; j++ {
			x, ok := unhex(s[j])
			if !ok {
				err = fmt.Errorf("unable to unescape string")
				return
			}
			v = v<<4 | x
		}
		s = s[n:]
		if v > utf8.MaxRune {
			err = fmt.Errorf("unable to unescape string")
			return
		}
		value = v
		multibyte = true

	case '0', '1', '2', '3':
		if len(s) < 2 {
			err = fmt.Errorf("unable to unescape octal sequence in string")
			return
		}
		v := rune(c - '0')
		for j := 0; j < 2; j++ {
			x := s[j]
			if x < '0' || x > '7' {
				err = fmt.Errorf("unable to unescape octal sequence in string")
				return
			}
			v = v*8 + rune(x-'0')
		}
		if v > utf8.MaxRune {
			err = fmt.Errorf("unable to unescape string")
			return
		}
		value = v
		s = s[2:]
		multibyte = true

	default:
		value = '\\'
		tail = string(c) + s
		return
	}

	tail = s
	return
}

func unhex(b byte) (rune, bool) {
	c := rune(b)
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
}
