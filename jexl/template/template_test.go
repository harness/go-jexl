// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package template_test

import (
	"errors"
	"testing"

	jexl "github.com/harness/go-jexl/jexl"
	"github.com/harness/go-jexl/jexl/template"
)

var env = map[string]any{
	"name":    "world",
	"count":   42,
	"ok":      true,
	"nothing": nil,
	"a":       "A",
	"b":       "B",
	"id":      "tok-123",
	"foo":     map[string]any{"name": "world"},
	"bar":     map[string]any{"baz": "name"},
}

func eval(expr string) (any, error) {
	return jexl.Eval(expr, env)
}

func TestExecString(t *testing.T) {
	var tests = []struct {
		src  string
		want string
	}{
		// no expressions
		{"", ""},
		{"hello world", "hello world"},
		{"no delimiters here", "no delimiters here"},

		// single expression
		{"<+name>", "world"},
		{"hello <+name>", "hello world"},
		{"<+name> says hi", "world says hi"},
		{"hello <+name>!", "hello world!"},
		{"count: <+count>", "count: 42"},
		{"flag: <+ok>", "flag: true"},
		{"val: <+nothing>", "val: "},

		// multiple expressions
		{"<+a> and <+b>", "A and B"},
		{"<+a><+b>", "AB"},
		{"x=<+a>, y=<+b>", "x=A, y=B"},

		// ">" inside parens is not a closer
		{"<+(1 > 0)>", "true"},
		{"result: <+count + 1>", "result: 43"},

		// ">" inside string literals is not a closer
		{`<+'hello>world'>`, "hello>world"},
		{`<+"a>b">`, "a>b"},

		// backslash escapes inside string literals
		{`<+'it\'s'>`, "it's"},
		{`<+"say \"hi\"">`, `say "hi"`},

		// ">" at bracket depth > 0 is not a closer
		{"<+[1 > 0][0]>", "true"},
		{"<+{'v': count > 0}.v>", "true"},

		// nested expressions resolved inside-out
		{"<+<+a>>", "A"},
		{"<+'result:' + <+id>>", "result:tok-123"},

		// dynamic property access: inner resolves to a key name used as dot accessor
		{"<+ foo.<+ bar.baz >>", "world"},

		// leading/trailing whitespace in expression is trimmed
		{"<+ name >", "world"},
		{"<+ count >", "42"},
	}
	for _, tt := range tests {
		t.Run(tt.src, func(t *testing.T) {
			got, err := template.New(eval).ExecString(tt.src)
			if err != nil {
				t.Fatalf("ExecString(%q) error: %v", tt.src, err)
			}
			if got != tt.want {
				t.Errorf("ExecString(%q)\n got  %q\n want %q", tt.src, got, tt.want)
			}
		})
	}
}

func TestExec(t *testing.T) {
	got, err := template.New(eval).Exec([]byte("hello <+name>"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello world" {
		t.Errorf("got %q want %q", got, "hello world")
	}
}

func TestExec_CustomDelim(t *testing.T) {
	var tests = []struct {
		left  string
		right string
		src   string
		want  string
	}{
		{"${{", "}}", "hello ${{ name }}", "hello world"},
		{"${{", "}}", "${{ a }} and ${{ b }}", "A and B"},
		{"${{", "}}", "no delimiters", "no delimiters"},
		{"${", "}", "x=${count}", "x=42"},
		{"[[", "]]", "[[name]]", "world"},
	}
	for _, tt := range tests {
		t.Run(tt.src, func(t *testing.T) {
			got, err := template.New(eval).Delim(tt.left, tt.right).ExecString(tt.src)
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			if got != tt.want {
				t.Errorf("\n got  %q\n want %q", got, tt.want)
			}
		})
	}
}

func TestExec_Error(t *testing.T) {
	errFunc := func(expr string) (any, error) {
		return nil, errors.New("eval error")
	}
	_, err := template.New(errFunc).Exec([]byte("<+err>"))
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestExec_Unclosed(t *testing.T) {
	noop := func(expr string) (any, error) {
		return nil, nil
	}
	_, err := template.New(noop).Exec([]byte("<+name"))
	if err == nil {
		t.Fatal("expected error, got none")
	}
}
