// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package token

import "strings"

// Print returns a compact space-separated representation
// of a token sequence.
func Print(tokens []Token) string {
	parts := make([]string, 0, len(tokens))
	for _, t := range tokens {
		switch t.Kind {
		case EOF:
			parts = append(parts, "EOF")
		case Ident:
			parts = append(parts, "IDENT("+t.Value+")")
		case Number:
			parts = append(parts, "NUM("+t.Value+")")
		case String:
			parts = append(parts, "STR("+printEscape(t.Value)+")")
		case Template:
			parts = append(parts, "TMPL("+printEscape(t.Value)+")")
		case Operator:
			parts = append(parts, "OP("+t.Value+")")
		case Bracket:
			parts = append(parts, "BRK("+t.Value+")")
		case Regex:
			parts = append(parts, "RX("+t.Value+")")
		default:
			parts = append(parts, "?("+t.Value+")")
		}
	}
	return strings.Join(parts, " ")
}

// helper function replaces ASCII control characters with
// two-char escape sequences so that Print output stays on
// a single line.
func printEscape(s string) string {
	s = strings.ReplaceAll(s, "\a", `\a`)
	s = strings.ReplaceAll(s, "\b", `\b`)
	s = strings.ReplaceAll(s, "\f", `\f`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	s = strings.ReplaceAll(s, "\t", `\t`)
	s = strings.ReplaceAll(s, "\v", `\v`)
	return s
}
