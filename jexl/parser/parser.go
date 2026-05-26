// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package parser

import (
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/harness/go-jexl/jexl/ast"
	"github.com/harness/go-jexl/jexl/config"
	"github.com/harness/go-jexl/jexl/internal/decimal"
	"github.com/harness/go-jexl/jexl/scanner"
	"github.com/harness/go-jexl/jexl/token"
)

type arg byte

const expr arg = 1 << iota

const optional arg = 1 << 7

type parser struct {
	scanner                    *scanner.Scanner
	current, stashed, stashed2 token.Token
	hasStash, hasStash2        bool
	err                        *token.Error
	config                     *config.Config
	nodeCount                  uint
	localVars                  map[string]bool
	letVars                    map[string]bool
}

func (p *parser) parse(input string, config *config.Config) (*ast.Tree, error) {
	if p.scanner == nil {
		p.scanner = scanner.New()
	}
	p.config = config
	source := token.NewFile(input)
	p.scanner.Reset(source)
	p.next()
	node := p.parseScript()

	if !p.current.Is(token.EOF) {
		p.error("unexpected token %v", p.current.String())
	}

	tree := &ast.Tree{
		Node:   node,
		Source: source,
	}
	err := p.err

	p.err = nil
	p.config = nil
	p.localVars = nil
	p.letVars = nil
	p.scanner.Reset(token.File{})

	if err != nil {
		return tree, err.Bind(source)
	}

	return tree, nil
}

func (p *parser) checkNodeLimit() error {
	p.nodeCount++
	if p.config == nil {
		if p.nodeCount > config.MaxNodes {
			p.error("compilation failed: expression exceeds maximum allowed nodes")
			return nil
		}
		return nil
	}
	if p.config.MaxNodes > 0 && p.nodeCount > p.config.MaxNodes {
		p.error("compilation failed: expression exceeds maximum allowed nodes")
		return nil
	}
	return nil
}

func (p *parser) createNode(n ast.Node, loc token.Range) ast.Node {
	if err := p.checkNodeLimit(); err != nil {
		return nil
	}
	if n == nil || p.err != nil {
		return nil
	}
	n.SetLocation(loc)
	return n
}

func (p *parser) createMemberNode(n *ast.MemberExpr, loc token.Range) *ast.MemberExpr {
	if err := p.checkNodeLimit(); err != nil {
		return nil
	}
	if n == nil || p.err != nil {
		return nil
	}
	n.SetLocation(loc)
	return n
}

func Parse(input string) (*ast.Tree, error) {
	return ParseWithConfig(input, nil)
}

func ParseWithConfig(input string, config *config.Config) (*ast.Tree, error) {
	return new(parser).parse(input, config)
}

func (p *parser) error(format string, args ...any) {
	p.errorAt(p.current, format, args...)
}

func (p *parser) errorAt(tok token.Token, format string, args ...any) {
	if p.err == nil {
		p.err = &token.Error{
			Range:   tok.Range,
			Message: fmt.Sprintf(format, args...),
		}
	}
}

func (p *parser) next() {
	if p.hasStash {
		p.current = p.stashed
		p.hasStash = false
		if p.hasStash2 {
			p.stashed = p.stashed2
			p.hasStash = true
			p.hasStash2 = false
		}
		return
	}

	tok, err := p.scanner.Scan()
	var e *token.Error
	switch {
	case err == nil:
		p.current = tok
	case errors.Is(err, io.EOF):
		p.error("unexpected end of expression")
	case errors.As(err, &e):
		p.err = e
	default:
		p.err = &token.Error{
			Range:   p.current.Range,
			Message: "unknown lexing error",
			Prev:    err,
		}
	}
}

func (p *parser) expect(k token.Kind, values ...string) {
	if p.current.Is(k, values...) {
		p.next()
		return
	}
	p.error("unexpected token %v", p.current.String())
}

func (p *parser) isNamespaceCall() bool {
	tok1, err1 := p.scanner.Scan()
	if err1 != nil || (tok1.Kind != token.Ident && tok1.Kind != token.Operator) {
		p.stashed = tok1
		p.hasStash = true
		return false
	}
	tok2, err2 := p.scanner.Scan()
	if err2 != nil || !tok2.Is(token.Bracket, "(") {
		p.stashed = tok1
		p.stashed2 = tok2
		p.hasStash = true
		p.hasStash2 = true
		return false
	}
	p.stashed = tok1
	p.stashed2 = tok2
	p.hasStash = true
	p.hasStash2 = true
	return true
}

func (p *parser) isKnownNamespace(name string) bool {
	if p.config == nil || p.config.Registry == nil {
		return true
	}
	_, ok := p.config.Registry.Lookup(name)
	return ok
}

func (p *parser) parseExpression(precedence int) ast.Node {
	if p.err != nil {
		return nil
	}

	if precedence == 0 && p.current.Is(token.Operator, "if") {
		return p.parseConditionalIf()
	}

	nodeLeft := p.parsePrimary()

	prevOperator := ""
	opToken := p.current
	for opToken.Is(token.Operator) && p.err == nil {
		negate := opToken.Is(token.Operator, "not")
		var notToken token.Token

		if negate {
			tokenBackup := p.current
			p.next()
			if p.current.Value == "in" {
				if op, ok := token.Binary[p.current.Value]; ok && op.Precedence >= precedence {
					notToken = p.current
					opToken = p.current
				} else {
					p.hasStash = true
					p.stashed = p.current
					p.current = tokenBackup
					break
				}
			} else {
				p.error("unexpected token %v", p.current.String())
				break
			}
		}

		if op, ok := token.Binary[opToken.Value]; ok && op.Precedence >= precedence {
			p.next()

			if prevOperator == "??" && opToken.Value != "??" && !opToken.Is(token.Bracket, "(") {
				p.errorAt(opToken, "Operator (%v) and coalesce expressions (??) cannot be mixed. Wrap either by parentheses.", opToken.Value)
				break
			}

			if token.IsComparison(opToken.Value) {
				nodeLeft = p.parseComparison(nodeLeft, opToken, op.Precedence)
				goto next
			}

			var nodeRight ast.Node
			if op.Associativity == token.Left {
				nodeRight = p.parseExpression(op.Precedence + 1)
			} else {
				nodeRight = p.parseExpression(op.Precedence)
			}

			if opToken.Value == "|" || opToken.Value == "^" || opToken.Value == "&" ||
				opToken.Value == "<<" || opToken.Value == ">>" || opToken.Value == ">>>" {
				nodeLeft = p.createNode(&ast.BitExpr{
					Operator: opToken.Value,
					Left:     nodeLeft,
					Right:    nodeRight,
				}, opToken.Range)
			} else if opToken.Value == "===" || opToken.Value == "!==" {
				nodeLeft = p.createNode(&ast.StrictEqualExpr{
					Negated: opToken.Value == "!==",
					Left:    nodeLeft,
					Right:   nodeRight,
				}, opToken.Range)
			} else if token.IsAssignment(opToken.Value) {
				switch nodeLeft.(type) {
				case *ast.Ident, *ast.MemberExpr:
				default:
					p.errorAt(opToken, "invalid assignment target: left-hand side must be an identifier or member access")
					return nil
				}
				nodeLeft = p.createNode(&ast.AssignStmt{
					Target: nodeLeft,
					Op:     opToken.Value,
					Value:  nodeRight,
				}, opToken.Range)
			} else {
				nodeLeft = p.createNode(&ast.BinaryExpr{
					Operator: opToken.Value,
					Left:     nodeLeft,
					Right:    nodeRight,
				}, opToken.Range)
			}
			if nodeLeft == nil {
				return nil
			}

			if negate {
				nodeLeft = p.createNode(&ast.UnaryExpr{
					Operator: "not",
					Node:     nodeLeft,
				}, notToken.Range)
				if nodeLeft == nil {
					return nil
				}
			}

			goto next
		}
		break

	next:
		prevOperator = opToken.Value
		opToken = p.current
	}

	if precedence == 0 {
		nodeLeft = p.parseConditional(nodeLeft)
	}

	return nodeLeft
}

func (p *parser) parseScriptStatement() ast.Node {
	if p.err != nil {
		return nil
	}
	var stmt ast.Node
	if p.current.Is(token.Operator, "var") || p.current.Is(token.Operator, "let") || p.current.Is(token.Operator, "const") {
		stmt = p.parseVar()
	} else if p.current.Is(token.Operator, "return") {
		stmt = p.parseReturn()
	} else if p.current.Is(token.Operator, "break") {
		loc := p.current.Range
		p.next()
		stmt = p.createNode(&ast.BreakStmt{}, loc)
	} else if p.current.Is(token.Operator, "continue") {
		loc := p.current.Range
		p.next()
		stmt = p.createNode(&ast.ContinueStmt{}, loc)
	} else if p.current.Is(token.Operator, "throw") {
		stmt = p.parseThrow()
	} else {
		stmt = p.parseExpression(0)
	}
	for p.current.Is(token.Operator, ";") && p.err == nil {
		p.next()
	}
	return stmt
}

func (p *parser) parseScript() ast.Node {
	startToken := p.current
	var stmts []ast.Node
	for !p.current.Is(token.EOF) && !p.current.Is(token.Bracket, "}") && p.err == nil {
		for p.current.Is(token.Operator, ";") && p.err == nil {
			p.next()
		}
		if p.current.Is(token.EOF) || p.current.Is(token.Bracket, "}") {
			break
		}
		var stmt ast.Node
		if p.current.Is(token.Operator, "var") || p.current.Is(token.Operator, "let") || p.current.Is(token.Operator, "const") {
			stmt = p.parseVar()
		} else if p.current.Is(token.Operator, "for") {
			stmt = p.parseForeach()
		} else if p.current.Is(token.Operator, "while") {
			stmt = p.parseWhile()
		} else if p.current.Is(token.Operator, "do") {
			stmt = p.parseDoWhile()
		} else if p.current.Is(token.Operator, "return") {
			stmt = p.parseReturn()
		} else if p.current.Is(token.Operator, "break") {
			loc := p.current.Range
			p.next()
			stmt = p.createNode(&ast.BreakStmt{}, loc)
		} else if p.current.Is(token.Operator, "continue") {
			loc := p.current.Range
			p.next()
			stmt = p.createNode(&ast.ContinueStmt{}, loc)
		} else if p.current.Is(token.Operator, "try") {
			stmt = p.parseTryCatch()
		} else if p.current.Is(token.Operator, "throw") {
			stmt = p.parseThrow()
		} else if p.current.Is(token.Operator, "switch") {
			stmt = p.parseSwitch()
		} else if p.current.Is(token.Operator, "function") {
			stmt = p.parseFunctionDecl()
		} else {
			stmt = p.parseExpression(0)
		}
		if stmt == nil || p.err != nil {
			return nil
		}
		stmts = append(stmts, stmt)
		for p.current.Is(token.Operator, ";") && p.err == nil {
			p.next()
		}
	}
	if len(stmts) == 0 && p.err == nil {
		return p.createNode(&ast.NilLit{}, startToken.Range)
	}
	if len(stmts) == 1 {
		return stmts[0]
	}
	return p.createNode(&ast.BlockStmt{Statements: stmts}, startToken.Range)
}

func (p *parser) parseVar() ast.Node {
	keyword := p.current.Value
	kwToken := p.current
	p.next()
	if !p.current.Is(token.Ident) {
		p.error("expected identifier after %s", keyword)
		return nil
	}

	var names []token.Token
	names = append(names, p.current)
	p.next()
	for p.current.Is(token.Operator, ",") {
		p.next()
		if !p.current.Is(token.Ident) {
			p.error("expected identifier in %s declaration", keyword)
			return nil
		}
		names = append(names, p.current)
		p.next()
	}

	declareOne := func(nameToken token.Token) ast.Node {
		if keyword == "let" || keyword == "const" {
			if p.letVars != nil && p.letVars[nameToken.Value] {
				p.errorAt(nameToken, "identifier '%s' has already been declared", nameToken.Value)
				return nil
			}
			if p.letVars == nil {
				p.letVars = make(map[string]bool)
			}
			p.letVars[nameToken.Value] = true
		}
		p.declareLocalVar(nameToken.Value)
		var value ast.Node
		if len(names) == 1 && p.current.Is(token.Operator, "=") {
			p.next()
			value = p.parseExpression(0)
		} else {
			value = p.createNode(&ast.NilLit{}, nameToken.Range)
		}
		return p.createNode(&ast.VarDecl{
			Name:    nameToken.Value,
			Keyword: keyword,
			Value:   value,
		}, kwToken.Range)
	}

	if len(names) == 1 {
		return declareOne(names[0])
	}
	var stmts []ast.Node
	for _, nt := range names {
		node := declareOne(nt)
		if p.err != nil {
			return nil
		}
		stmts = append(stmts, node)
	}
	return p.createNode(&ast.BlockStmt{Statements: stmts}, kwToken.Range)
}

func (p *parser) declareLocalVar(name string) {
	if p.localVars == nil {
		p.localVars = make(map[string]bool)
	}
	p.localVars[name] = true
}

func (p *parser) parseLambdaBody() ast.Node {
	if p.current.Is(token.Bracket, "{") {
		p.next()
		body := p.parseScript()
		p.expect(token.Bracket, "}")
		return body
	}
	return p.parseExpression(0)
}

func (p *parser) parseFunctionDecl() ast.Node {
	funcToken := p.current
	p.next()
	if !p.current.Is(token.Ident) {
		p.error("expected function name after 'function', got %v", p.current.String())
		return nil
	}
	nameToken := p.current
	p.next()
	p.declareLocalVar(nameToken.Value)
	p.expect(token.Bracket, "(")
	if p.err != nil {
		return nil
	}
	var params []string
	for !p.current.Is(token.Bracket, ")") && !p.current.Is(token.EOF) && p.err == nil {
		if p.current.Is(token.Operator, "var") {
			p.next()
		}
		if !p.current.Is(token.Ident) {
			p.error("expected parameter name, got %v", p.current.String())
			return nil
		}
		params = append(params, p.current.Value)
		p.next()
		if p.current.Is(token.Operator, ",") {
			p.next()
		}
	}
	p.expect(token.Bracket, ")")
	if p.err != nil {
		return nil
	}
	p.expect(token.Bracket, "{")
	if p.err != nil {
		return nil
	}
	body := p.parseScript()
	p.expect(token.Bracket, "}")
	if p.err != nil {
		return nil
	}
	lambda := p.createNode(&ast.LambdaExpr{Params: params, Body: body}, funcToken.Range)
	if lambda == nil {
		return nil
	}
	return p.createNode(&ast.VarDecl{
		Name:    nameToken.Value,
		Keyword: "var",
		Value:   lambda,
	}, funcToken.Range)
}

func (p *parser) parseForeach() ast.Node {
	forToken := p.current
	p.next()
	p.expect(token.Bracket, "(")
	if p.err != nil {
		return nil
	}

	isKeyword := p.current.Is(token.Operator, "var") || p.current.Is(token.Operator, "let") || p.current.Is(token.Operator, "const")
	if isKeyword {
		kwTok := p.current
		p.next()
		if !p.current.Is(token.Ident) {
			p.error("expected identifier after %s in for loop", kwTok.Value)
			return nil
		}
		idTok := p.current
		p.next()

		if p.current.Is(token.Operator, ":") {
			p.next()
			collection := p.parseExpression(0)
			p.expect(token.Bracket, ")")
			p.expect(token.Bracket, "{")
			body := p.parseScript()
			p.expect(token.Bracket, "}")
			return p.createNode(&ast.ForeachStmt{Var: idTok.Value, Collection: collection, Body: body}, forToken.Range)
		}
		p.declareLocalVar(idTok.Value)
		var initNode ast.Node
		if p.current.Is(token.Operator, "=") {
			p.next()
			val := p.parseExpression(0)
			initNode = p.createNode(&ast.VarDecl{Name: idTok.Value, Keyword: kwTok.Value, Value: val}, kwTok.Range)
		} else {
			initNode = p.createNode(&ast.VarDecl{Name: idTok.Value, Keyword: kwTok.Value, Value: p.createNode(&ast.NilLit{}, idTok.Range)}, kwTok.Range)
		}
		return p.parseCStyleFor(forToken, initNode)
	}

	if p.current.Is(token.Ident) {
		idTok := p.current
		p.next()
		if p.current.Is(token.Operator, ":") {
			p.next()
			collection := p.parseExpression(0)
			p.expect(token.Bracket, ")")
			p.expect(token.Bracket, "{")
			body := p.parseScript()
			p.expect(token.Bracket, "}")
			return p.createNode(&ast.ForeachStmt{Var: idTok.Value, Collection: collection, Body: body}, forToken.Range)
		}
		p.stashed2 = p.stashed
		p.hasStash2 = p.hasStash
		p.stashed = p.current
		p.hasStash = true
		p.current = idTok
		initNode := p.parseExpression(0)
		return p.parseCStyleFor(forToken, initNode)
	}

	var initNode ast.Node
	if !p.current.Is(token.Operator, ";") {
		initNode = p.parseExpression(0)
	}
	return p.parseCStyleFor(forToken, initNode)
}

func (p *parser) parseCStyleFor(forToken token.Token, init ast.Node) ast.Node {
	p.expect(token.Operator, ";")
	if p.err != nil {
		return nil
	}
	var cond ast.Node
	if !p.current.Is(token.Operator, ";") {
		cond = p.parseExpression(0)
	}
	p.expect(token.Operator, ";")
	if p.err != nil {
		return nil
	}
	var post ast.Node
	if !p.current.Is(token.Bracket, ")") {
		post = p.parseExpression(0)
	}
	p.expect(token.Bracket, ")")
	p.expect(token.Bracket, "{")
	body := p.parseScript()
	p.expect(token.Bracket, "}")
	return p.createNode(&ast.ForStmt{Init: init, Cond: cond, Post: post, Body: body}, forToken.Range)
}

func (p *parser) parseDoWhile() ast.Node {
	doToken := p.current
	p.next()
	p.expect(token.Bracket, "{")
	body := p.parseScript()
	p.expect(token.Bracket, "}")
	p.expect(token.Operator, "while")
	p.expect(token.Bracket, "(")
	cond := p.parseExpression(0)
	p.expect(token.Bracket, ")")
	return p.createNode(&ast.DoWhileStmt{Body: body, Cond: cond}, doToken.Range)
}

func (p *parser) parseWhile() ast.Node {
	whileToken := p.current
	p.next()
	p.expect(token.Bracket, "(")
	cond := p.parseExpression(0)
	p.expect(token.Bracket, ")")
	p.expect(token.Bracket, "{")
	body := p.parseScript()
	p.expect(token.Bracket, "}")
	return p.createNode(&ast.WhileStmt{
		Cond: cond,
		Body: body,
	}, whileToken.Range)
}

func (p *parser) parseReturn() ast.Node {
	retToken := p.current
	p.next()
	var value ast.Node
	if !p.current.Is(token.Operator, ";") && !p.current.Is(token.EOF) && !p.current.Is(token.Bracket, "}") {
		value = p.parseExpression(0)
	}
	return p.createNode(&ast.ReturnStmt{Value: value}, retToken.Range)
}

func (p *parser) parseTryCatch() ast.Node {
	tryToken := p.current
	p.next()
	p.expect(token.Bracket, "{")
	body := p.parseScript()
	p.expect(token.Bracket, "}")

	var catchVar string
	var catchBody ast.Node
	if p.current.Is(token.Operator, "catch") {
		p.next()
		p.expect(token.Bracket, "(")
		if p.current.Is(token.Operator, "var") || p.current.Is(token.Operator, "let") || p.current.Is(token.Operator, "const") {
			p.next()
		}
		if !p.current.Is(token.Ident) {
			p.error("expected identifier in catch clause")
			return nil
		}
		catchVarToken := p.current
		catchVar = catchVarToken.Value
		p.next()
		p.expect(token.Bracket, ")")
		p.expect(token.Bracket, "{")
		catchBody = p.parseScript()
		p.expect(token.Bracket, "}")
	}

	var finallyBody ast.Node
	if p.current.Is(token.Operator, "finally") {
		p.next()
		p.expect(token.Bracket, "{")
		finallyBody = p.parseScript()
		p.expect(token.Bracket, "}")
	}

	if catchBody == nil && finallyBody == nil {
		p.errorAt(tryToken, "try requires catch and/or finally")
		return nil
	}

	return p.createNode(&ast.TryStmt{
		Body:        body,
		CatchVar:    catchVar,
		CatchBody:   catchBody,
		FinallyBody: finallyBody,
	}, tryToken.Range)
}

func (p *parser) parseThrow() ast.Node {
	throwToken := p.current
	p.next()
	value := p.parseExpression(0)
	return p.createNode(&ast.ThrowStmt{Value: value}, throwToken.Range)
}

func (p *parser) parseSwitch() ast.Node {
	switchToken := p.current
	p.next()
	p.expect(token.Bracket, "(")
	subject := p.parseExpression(0)
	p.expect(token.Bracket, ")")
	p.expect(token.Bracket, "{")

	var cases []ast.CaseClause
	var defaultBody ast.Node
	seenCaseValues := map[string]bool{}

	for !p.current.Is(token.Bracket, "}") && !p.current.Is(token.EOF) && p.err == nil {
		if p.current.Is(token.Operator, "case") {
			p.next()
			var values []ast.Node
			values = append(values, p.parseExpression(0))
			for p.current.Is(token.Operator, ",") && p.err == nil {
				p.next()
				values = append(values, p.parseExpression(0))
			}
			for _, v := range values {
				var key string
				switch n := v.(type) {
				case *ast.IntLit:
					key = fmt.Sprintf("int:%d", n.Value)
				case *ast.FloatLit:
					key = fmt.Sprintf("float:%v", n.Value)
				case *ast.StringLit:
					key = fmt.Sprintf("str:%s", n.Value)
				case *ast.BoolLit:
					key = fmt.Sprintf("bool:%v", n.Value)
				}
				if key != "" {
					if seenCaseValues[key] {
						p.error("duplicate case value in switch")
						return nil
					}
					seenCaseValues[key] = true
				}
			}
			var body ast.Node
			if p.current.Is(token.Operator, "->") {
				p.next()
				expr := p.parseExpression(0)
				if p.current.Is(token.Operator, ";") {
					p.next()
				}
				body = expr
			} else {
				p.expect(token.Operator, ":")
				body = p.parseSwitchBody()
			}
			cases = append(cases, ast.CaseClause{Values: values, Body: body})
		} else if p.current.Is(token.Operator, "default") {
			p.next()
			if p.current.Is(token.Operator, "->") {
				p.next()
				expr := p.parseExpression(0)
				if p.current.Is(token.Operator, ";") {
					p.next()
				}
				defaultBody = expr
			} else {
				p.expect(token.Operator, ":")
				defaultBody = p.parseSwitchBody()
			}
		} else {
			p.error("expected case or default in switch body, got %v", p.current.String())
			return nil
		}
	}
	p.expect(token.Bracket, "}")
	return p.createNode(&ast.SwitchStmt{
		Subject: subject,
		Cases:   cases,
		Default: defaultBody,
	}, switchToken.Range)
}

func (p *parser) parseSwitchBody() ast.Node {
	startToken := p.current
	var stmts []ast.Node
	for !p.current.Is(token.EOF) &&
		!p.current.Is(token.Bracket, "}") &&
		!p.current.Is(token.Operator, "case") &&
		!p.current.Is(token.Operator, "default") &&
		p.err == nil {
		var stmt ast.Node
		if p.current.Is(token.Operator, "var") || p.current.Is(token.Operator, "let") || p.current.Is(token.Operator, "const") {
			stmt = p.parseVar()
		} else if p.current.Is(token.Operator, "for") {
			stmt = p.parseForeach()
		} else if p.current.Is(token.Operator, "while") {
			stmt = p.parseWhile()
		} else if p.current.Is(token.Operator, "do") {
			stmt = p.parseDoWhile()
		} else if p.current.Is(token.Operator, "return") {
			stmt = p.parseReturn()
		} else if p.current.Is(token.Operator, "break") {
			loc := p.current.Range
			p.next()
			stmt = p.createNode(&ast.BreakStmt{}, loc)
		} else if p.current.Is(token.Operator, "continue") {
			loc := p.current.Range
			p.next()
			stmt = p.createNode(&ast.ContinueStmt{}, loc)
		} else if p.current.Is(token.Operator, "try") {
			stmt = p.parseTryCatch()
		} else if p.current.Is(token.Operator, "throw") {
			stmt = p.parseThrow()
		} else if p.current.Is(token.Operator, "switch") {
			stmt = p.parseSwitch()
		} else if p.current.Is(token.Operator, "function") {
			stmt = p.parseFunctionDecl()
		} else {
			stmt = p.parseExpression(0)
		}
		if stmt == nil || p.err != nil {
			return nil
		}
		stmts = append(stmts, stmt)
		for p.current.Is(token.Operator, ";") && p.err == nil {
			p.next()
		}
	}
	if len(stmts) == 0 {
		return p.createNode(&ast.BlockStmt{Statements: nil}, startToken.Range)
	}
	if len(stmts) == 1 {
		return stmts[0]
	}
	return p.createNode(&ast.BlockStmt{Statements: stmts}, startToken.Range)
}

func (p *parser) parseConditionalIf() ast.Node {
	p.next()
	if p.err != nil {
		return nil
	}
	nodeCondition := p.parseExpression(0)
	var expr1 ast.Node
	if p.current.Is(token.Bracket, "{") {
		p.next()
		expr1 = p.parseScript()
		p.expect(token.Bracket, "}")
	} else {
		expr1 = p.parseScriptStatement()
	}

	if !p.current.Is(token.Operator, "else") {
		return &ast.ConditionalExpr{
			Cond: nodeCondition,
			Exp1: expr1,
			Exp2: nil,
		}
	}
	p.next()

	var expr2 ast.Node
	if p.current.Is(token.Operator, "if") {
		expr2 = p.parseConditionalIf()
	} else if p.current.Is(token.Bracket, "{") {
		p.next()
		expr2 = p.parseScript()
		p.expect(token.Bracket, "}")
	} else {
		expr2 = p.parseScriptStatement()
	}

	return &ast.ConditionalExpr{
		Cond: nodeCondition,
		Exp1: expr1,
		Exp2: expr2,
	}
}

func (p *parser) parseConditional(node ast.Node) ast.Node {
	var expr1, expr2 ast.Node
	for p.current.Is(token.Operator, "?") && p.err == nil {
		p.next()

		if !p.current.Is(token.Operator, ":") {
			expr1 = p.parseExpression(0)
			p.expect(token.Operator, ":")
			expr2 = p.parseExpression(0)
		} else {
			p.next()
			expr1 = node
			expr2 = p.parseExpression(0)
		}

		node = p.createNode(&ast.ConditionalExpr{
			Ternary: true,
			Cond:    node,
			Exp1:    expr1,
			Exp2:    expr2,
		}, p.current.Range)
		if node == nil {
			return nil
		}
	}
	return node
}

func (p *parser) parsePrimary() ast.Node {
	tok := p.current

	if tok.Is(token.Operator, "++", "--") {
		p.next()
		target := p.parsePrimary()
		switch target.(type) {
		case *ast.Ident, *ast.MemberExpr:
			node := p.createNode(&ast.IncDecStmt{
				Target: target,
				Op:     tok.Value,
				Prefix: true,
			}, tok.Range)
			if node == nil {
				return nil
			}
			return node
		default:
			unaryOp := string(tok.Value[0])
			inner := p.createNode(&ast.UnaryExpr{Operator: unaryOp, Node: target}, tok.Range)
			if inner == nil {
				return nil
			}
			outer := p.createNode(&ast.UnaryExpr{Operator: unaryOp, Node: inner}, tok.Range)
			if outer == nil {
				return nil
			}
			return p.parsePostfixExpression(outer)
		}
	}

	if tok.Is(token.Operator, "new") {
		p.next()
		var className string
		if p.current.Is(token.Ident) {
			className = p.current.Value
			p.next()
			args := p.parseArguments([]ast.Node{})
			node := p.createNode(&ast.NewExpr{
				ClassName: className,
				Args:      args,
			}, tok.Range)
			if node == nil {
				return nil
			}
			return p.parsePostfixExpression(node)
		} else if p.current.Is(token.Bracket, "(") {
			p.next()
			if !p.current.Is(token.String) {
				p.error("expected string class name after new(, got %v", p.current.String())
				return nil
			}
			className = p.current.Value
			p.next()
			args := []ast.Node{}
			for p.current.Is(token.Operator, ",") {
				p.next()
				arg := p.parseExpression(0)
				if arg == nil {
					return nil
				}
				args = append(args, arg)
			}
			if !p.current.Is(token.Bracket, ")") {
				p.error("expected ')' after new arguments, got %v", p.current.String())
				return nil
			}
			p.next()
			node := p.createNode(&ast.NewExpr{
				ClassName: className,
				Args:      args,
			}, tok.Range)
			if node == nil {
				return nil
			}
			return p.parsePostfixExpression(node)
		} else {
			p.error("expected class name after new, got %v", p.current.String())
			return nil
		}
	}

	if tok.Is(token.Operator, "empty") {
		p.next()
		operand := p.parseExpression(90)
		node := p.createNode(&ast.EmptyExpr{Value: operand}, tok.Range)
		if node == nil {
			return nil
		}
		return p.parsePostfixExpression(node)
	}

	if tok.Is(token.Operator, "size") {
		p.next()
		operand := p.parseExpression(90)
		node := p.createNode(&ast.SizeExpr{Value: operand}, tok.Range)
		if node == nil {
			return nil
		}
		return p.parsePostfixExpression(node)
	}

	if tok.Is(token.Operator, "switch") {
		node := p.parseSwitch()
		if node == nil {
			return nil
		}
		return p.parsePostfixExpression(node)
	}

	if tok.Is(token.Operator, "function") {
		p.next()
		if p.current.Is(token.Ident) {
			p.next()
		}
		p.expect(token.Bracket, "(")
		if p.err != nil {
			return nil
		}
		var params []string
		for !p.current.Is(token.Bracket, ")") && !p.current.Is(token.EOF) && p.err == nil {
			if p.current.Is(token.Operator, "var") {
				p.next()
			}
			if !p.current.Is(token.Ident) {
				p.error("expected parameter name, got %v", p.current.String())
				return nil
			}
			params = append(params, p.current.Value)
			p.next()
			if p.current.Is(token.Operator, ",") {
				p.next()
			}
		}
		p.expect(token.Bracket, ")")
		if p.err != nil {
			return nil
		}
		p.expect(token.Bracket, "{")
		if p.err != nil {
			return nil
		}
		body := p.parseScript()
		p.expect(token.Bracket, "}")
		if p.err != nil {
			return nil
		}
		return p.createNode(&ast.LambdaExpr{Params: params, Body: body}, tok.Range)
	}

	if tok.Is(token.Operator) {
		if op, ok := token.Unary[tok.Value]; ok {
			p.next()
			exprPrec := op.Precedence
			if tok.Value == "!" || tok.Value == "not" {
				exprPrec = 17
			}
			expr := p.parseExpression(exprPrec)
			var node ast.Node
			if tok.Value == "~" {
				if _, ok := expr.(*ast.RegexLit); ok {
					return p.parsePostfixExpression(expr)
				}
				node = p.createNode(&ast.BitExpr{
					Operator: "~",
					Left:     expr,
				}, tok.Range)
			} else {
				node = p.createNode(&ast.UnaryExpr{
					Operator: tok.Value,
					Node:     expr,
				}, tok.Range)
			}
			if node == nil {
				return nil
			}
			return p.parsePostfixExpression(node)
		}
	}

	if tok.Is(token.Bracket, "(") {
		p.next()

		if p.current.Is(token.Bracket, ")") {
			p.next()
			if p.current.Is(token.Operator, "->") || p.current.Is(token.Operator, "=>") {
				p.next()
				body := p.parseLambdaBody()
				return p.createNode(&ast.LambdaExpr{Params: []string{}, Body: body}, tok.Range)
			}
			p.error("unexpected empty parentheses")
			return nil
		}

		startsWithVar := p.current.Is(token.Operator, "var")
		if startsWithVar {
			p.next()
		}
		if p.current.Is(token.Ident) {
			firstParam := p.current
			p.next()

			if p.current.Is(token.Bracket, ")") {
				p.next()
				if p.current.Is(token.Operator, "->") || p.current.Is(token.Operator, "=>") {
					p.next()
					body := p.parseLambdaBody()
					return p.createNode(&ast.LambdaExpr{
						Params: []string{firstParam.Value},
						Body:   body,
					}, tok.Range)
				}
				if startsWithVar {
					p.error("unexpected empty parentheses after var")
					return nil
				}
				var node ast.Node
				switch firstParam.Value {
				case "true":
					node = p.createNode(&ast.BoolLit{Value: true}, firstParam.Range)
				case "false":
					node = p.createNode(&ast.BoolLit{Value: false}, firstParam.Range)
				case "nil", "null":
					node = p.createNode(&ast.NilLit{}, firstParam.Range)
				default:
					node = p.createNode(&ast.Ident{Value: firstParam.Value}, firstParam.Range)
				}
				if node == nil {
					return nil
				}
				return p.parsePostfixExpression(node)
			}

			if p.current.Is(token.Operator, ",") {
				params := []string{firstParam.Value}
				for p.current.Is(token.Operator, ",") && p.err == nil {
					p.next()
					if p.current.Is(token.Operator, "var") {
						p.next()
					}
					if !p.current.Is(token.Ident) {
						p.error("expected identifier in lambda parameter list, got %v", p.current.String())
						return nil
					}
					params = append(params, p.current.Value)
					p.next()
				}
				p.expect(token.Bracket, ")")
				if p.err != nil {
					return nil
				}
				if p.current.Is(token.Operator, "->") || p.current.Is(token.Operator, "=>") {
					p.next()
					body := p.parseLambdaBody()
					return p.createNode(&ast.LambdaExpr{Params: params, Body: body}, tok.Range)
				}
				p.error("expected -> or => after lambda parameter list")
				return nil
			}

			if startsWithVar {
				p.error("unexpected token after var parameter")
				return nil
			}
			p.stashed = p.current
			p.hasStash = true
			p.current = firstParam
			expr := p.parseExpression(0)
			p.expect(token.Bracket, ")")
			return p.parsePostfixExpression(expr)
		} else if startsWithVar {
			p.error("expected identifier after var in parameter list")
			return nil
		}

		expr := p.parseExpression(0)
		p.expect(token.Bracket, ")")
		return p.parsePostfixExpression(expr)
	}

	if tok.Is(token.Operator, "::") {
		p.next()
		tok = p.current
		p.expect(token.Ident)
		return p.parsePostfixExpression(p.parseCall(tok, []ast.Node{}, false))
	}

	return p.parseSecondary()
}

func (p *parser) parseSecondary() ast.Node {
	var node ast.Node
	tok := p.current

	switch tok.Kind {

	case token.Ident:
		p.next()
		switch tok.Value {
		case "true":
			node = p.createNode(&ast.BoolLit{Value: true}, tok.Range)
			if node == nil {
				return nil
			}
			return node
		case "false":
			node = p.createNode(&ast.BoolLit{Value: false}, tok.Range)
			if node == nil {
				return nil
			}
			return node
		case "nil", "null":
			node = p.createNode(&ast.NilLit{}, tok.Range)
			if node == nil {
				return nil
			}
			return node
		default:
			if p.current.Is(token.Bracket, "(") {
				node = p.parseCall(tok, []ast.Node{}, true)
			} else if p.current.Is(token.Operator, ":") && p.isKnownNamespace(tok.Value) && p.isNamespaceCall() {
				p.next()
				methodToken := p.current
				p.next()
				node = p.createNode(&ast.NamespaceCallExpr{
					Namespace: tok.Value,
					Method:    methodToken.Value,
					Args:      p.parseArguments([]ast.Node{}),
				}, tok.Range)
				if node == nil {
					return nil
				}
				return p.parsePostfixExpression(node)
			} else if p.current.Is(token.Operator, "->") || p.current.Is(token.Operator, "=>") {
				p.next()
				body := p.parseLambdaBody()
				node = p.createNode(&ast.LambdaExpr{
					Params: []string{tok.Value},
					Body:   body,
				}, tok.Range)
				if node == nil {
					return nil
				}
				return node
			} else {
				node = p.createNode(&ast.Ident{Value: tok.Value}, tok.Range)
				if node == nil {
					return nil
				}
			}
		}

	case token.Number:
		p.next()
		value := strings.Replace(tok.Value, "_", "", -1)
		isBigInt, isBigDec, isFloat32 := false, false, false
		if len(value) > 1 && !strings.HasPrefix(strings.ToLower(value), "0x") {
			last := value[len(value)-1]
			if last == 'H' || last == 'h' {
				isBigInt = true
				value = value[:len(value)-1]
			} else if last == 'b' {
				isBigDec = true
				value = value[:len(value)-1]
			} else if last == 'f' {
				isFloat32 = true
				value = value[:len(value)-1]
			} else if last == 'B' || last == 'L' || last == 'l' || last == 'd' {
				value = value[:len(value)-1]
			}
		}
		if isBigInt {
			bi := new(big.Int)
			if _, ok := bi.SetString(value, 10); !ok {
				p.error("invalid BigInteger literal: %s", value)
				return nil
			}
			node := p.createNode(&ast.ConstantExpr{Value: bi}, tok.Range)
			return node
		}
		if isBigDec {
			d, err := decimal.NewFromString(value)
			if err != nil {
				p.error("invalid BigDecimal literal: %s", value)
				return nil
			}
			node := p.createNode(&ast.ConstantExpr{Value: d}, tok.Range)
			return node
		}
		var node ast.Node
		valueLower := strings.ToLower(value)
		switch {
		case strings.HasPrefix(valueLower, "0x"):
			number, err := strconv.ParseInt(value, 0, 64)
			if err != nil {
				p.error("invalid hex literal: %v", err)
			}
			node = p.toIntegerNode(number)
		case strings.ContainsAny(valueLower, ".e"):
			number, err := strconv.ParseFloat(value, 64)
			if err != nil {
				p.error("invalid float literal: %v", err)
			}
			node = p.toFloatNode(number)
		case strings.HasPrefix(valueLower, "0b"):
			number, err := strconv.ParseInt(value, 0, 64)
			if err != nil {
				p.error("invalid binary literal: %v", err)
			}
			node = p.toIntegerNode(number)
		case strings.HasPrefix(valueLower, "0o"):
			number, err := strconv.ParseInt(value, 0, 64)
			if err != nil {
				p.error("invalid octal literal: %v", err)
			}
			node = p.toIntegerNode(number)
		case len(value) > 1 && value[0] == '0' && !strings.ContainsAny(valueLower, ".e"):
			number, err := strconv.ParseInt(value, 8, 64)
			if err != nil {
				p.error("invalid octal literal: %v", err)
			}
			node = p.toIntegerNode(number)
		default:
			number, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				p.error("invalid integer literal: %v", err)
			}
			node = p.toIntegerNode(number)
		}
		if node != nil {
			node.SetLocation(tok.Range)
			if isFloat32 {
				if fn, ok := node.(*ast.FloatLit); ok {
					fn.SetType(reflect.TypeOf(float32(0)))
				}
			}
		}
		return node

	case token.String:
		p.next()
		node = p.createNode(&ast.StringLit{Value: tok.Value}, tok.Range)
		if node == nil {
			return nil
		}

	case token.Template:
		p.next()
		node = p.parseTemplateString(tok)
		if node == nil {
			return nil
		}

	case token.Regex:
		p.next()
		lastSlash := strings.LastIndex(tok.Value, "/")
		var pattern, flags string
		if lastSlash < 0 {
			pattern = tok.Value
		} else {
			pattern = tok.Value[:lastSlash]
			flags = tok.Value[lastSlash+1:]
		}
		node = p.createNode(&ast.RegexLit{Pattern: pattern, Flags: flags}, tok.Range)
		if node == nil {
			return nil
		}

	default:
		if tok.Is(token.Bracket, "[") {
			node = p.parseArrayExpression(tok)
		} else if tok.Is(token.Bracket, "{") {
			node = p.parseSetOrMapExpression(tok)
		} else {
			p.error("unexpected token %v", tok.String())
		}
	}

	return p.parsePostfixExpression(node)
}

func (p *parser) toIntegerNode(number int64) ast.Node {
	if number > math.MaxInt {
		p.error("integer literal is too large")
		return nil
	}
	return p.createNode(&ast.IntLit{Value: int(number)}, p.current.Range)
}

func (p *parser) toFloatNode(number float64) ast.Node {
	if number > math.MaxFloat64 {
		p.error("float literal is too large")
		return nil
	}
	return p.createNode(&ast.FloatLit{Value: number}, p.current.Range)
}

func (p *parser) parseTemplateString(tok token.Token) ast.Node {
	raw := tok.Value
	loc := tok.Range

	type segment struct {
		text   string
		isExpr bool
	}
	var segments []segment
	i := 0
	for i < len(raw) {
		idx := strings.Index(raw[i:], "${")
		if idx < 0 {
			segments = append(segments, segment{text: raw[i:]})
			break
		}
		if idx > 0 {
			segments = append(segments, segment{text: raw[i : i+idx]})
		}
		start := i + idx + 2
		depth := 1
		j := start
		for j < len(raw) && depth > 0 {
			switch raw[j] {
			case '{':
				depth++
			case '}':
				depth--
			}
			j++
		}
		if depth != 0 {
			p.errorAt(tok, "unterminated template expression")
			return nil
		}
		segments = append(segments, segment{text: raw[start : j-1], isExpr: true})
		i = j
	}

	buildNode := func(seg segment) ast.Node {
		if !seg.isExpr {
			return p.createNode(&ast.StringLit{Value: seg.text}, loc)
		}
		tree, err := new(parser).parse(seg.text, p.config)
		if err != nil {
			p.errorAt(tok, "template expression error: %v", err)
			return nil
		}
		return tree.Node
	}

	if len(segments) == 0 {
		return p.createNode(&ast.StringLit{Value: ""}, loc)
	}

	result := buildNode(segments[0])
	if result == nil {
		return nil
	}
	for _, seg := range segments[1:] {
		right := buildNode(seg)
		if right == nil {
			return nil
		}
		result = p.createNode(&ast.BinaryExpr{Operator: "+", Left: result, Right: right}, loc)
		if result == nil {
			return nil
		}
	}
	return result
}

func (p *parser) parseCall(tok token.Token, arguments []ast.Node, checkOverrides bool) ast.Node {
	callee := p.createNode(&ast.Ident{Value: tok.Value}, tok.Range)
	if callee == nil {
		return nil
	}
	node := p.createNode(&ast.CallExpr{
		Callee:    callee,
		Arguments: p.parseArguments(arguments),
	}, tok.Range)
	if node == nil {
		return nil
	}
	return node
}

func (p *parser) parseArguments(arguments []ast.Node) []ast.Node {
	offset := len(arguments)

	p.expect(token.Bracket, "(")
	for !p.current.Is(token.Bracket, ")") && p.err == nil {
		if len(arguments) > offset {
			p.expect(token.Operator, ",")
		}
		if p.current.Is(token.Bracket, ")") {
			break
		}
		node := p.parseExpression(0)
		arguments = append(arguments, node)
	}
	p.expect(token.Bracket, ")")

	return arguments
}

func (p *parser) parseArrayExpression(tok token.Token) ast.Node {
	nodes := make([]ast.Node, 0)

	p.expect(token.Bracket, "[")
	for !p.current.Is(token.Bracket, "]") && p.err == nil {
		if len(nodes) > 0 {
			p.expect(token.Operator, ",")
			if p.current.Is(token.Bracket, "]") {
				goto end
			}
		}
		if p.current.Is(token.Operator, "..") {
			p.next()
			if p.current.Is(token.Operator, ".") {
				p.next()
			}
			goto end
		}
		node := p.parseExpression(0)
		nodes = append(nodes, node)
	}
end:
	p.expect(token.Bracket, "]")

	node := p.createNode(&ast.ArrayLit{Nodes: nodes}, tok.Range)
	if node == nil {
		return nil
	}
	return node
}

func (p *parser) parseSetOrMapExpression(tok token.Token) ast.Node {
	p.expect(token.Bracket, "{")

	if p.current.Is(token.Bracket, "}") {
		p.next()
		node := p.createNode(&ast.MapLit{Pairs: []ast.Node{}}, tok.Range)
		if node == nil {
			return nil
		}
		return node
	}

	var firstElem ast.Node
	if p.current.Is(token.Number) || p.current.Is(token.String) || p.current.Is(token.Ident) {
		firstTok := p.current
		firstElem = p.createNode(&ast.StringLit{Value: firstTok.Value}, firstTok.Range)
		if firstElem == nil {
			return nil
		}
		p.next()

		if p.current.Is(token.Operator, ":") {
			firstElem = p.tokenToMapKey(firstTok)
			if firstElem == nil {
				return nil
			}
			return p.parseMapTail(tok, firstElem)
		}

		firstElem = p.tokenToSetElement(firstTok)
		if firstElem == nil {
			return nil
		}
		return p.parseSetTail(tok, firstElem)
	} else if p.current.Is(token.Bracket, "(") {
		firstElem = p.parseExpression(0)
	} else {
		firstElem = p.parseExpression(0)
		return p.parseSetTail(tok, firstElem)
	}

	if p.current.Is(token.Operator, ":") {
		return p.parseMapTail(tok, firstElem)
	}
	return p.parseSetTail(tok, firstElem)
}

func (p *parser) parseMapTail(tok token.Token, firstKey ast.Node) ast.Node {
	p.expect(token.Operator, ":")
	firstValue := p.parseExpression(0)
	firstPair := p.createNode(&ast.KeyValueExpr{Key: firstKey, Value: firstValue}, tok.Range)
	if firstPair == nil {
		return nil
	}
	nodes := []ast.Node{firstPair}

	for !p.current.Is(token.Bracket, "}") && p.err == nil {
		p.expect(token.Operator, ",")
		if p.current.Is(token.Bracket, "}") {
			break
		}
		if p.current.Is(token.Operator, ",") {
			p.error("unexpected token %v", p.current.String())
			return nil
		}

		var key ast.Node
		if p.current.Is(token.Number) || p.current.Is(token.String) || p.current.Is(token.Ident) {
			key = p.tokenToMapKey(p.current)
			if key == nil {
				return nil
			}
			p.next()
		} else if p.current.Is(token.Bracket, "(") {
			key = p.parseExpression(0)
		} else {
			p.error("a map key must be a quoted string, a number, an identifier, or an expression enclosed in parentheses (unexpected token %v)", p.current.String())
			return nil
		}

		p.expect(token.Operator, ":")
		val := p.parseExpression(0)
		pair := p.createNode(&ast.KeyValueExpr{Key: key, Value: val}, tok.Range)
		if pair == nil {
			return nil
		}
		nodes = append(nodes, pair)
	}

	p.expect(token.Bracket, "}")
	node := p.createNode(&ast.MapLit{Pairs: nodes}, tok.Range)
	if node == nil {
		return nil
	}
	return node
}

func (p *parser) parseSetTail(tok token.Token, firstElem ast.Node) ast.Node {
	elems := []ast.Node{firstElem}

	for !p.current.Is(token.Bracket, "}") && p.err == nil {
		p.expect(token.Operator, ",")
		if p.current.Is(token.Bracket, "}") {
			break
		}
		elem := p.parseExpression(0)
		elems = append(elems, elem)
	}

	p.expect(token.Bracket, "}")
	node := p.createNode(&ast.SetLit{Elements: elems}, tok.Range)
	if node == nil {
		return nil
	}
	return node
}

func (p *parser) tokenToMapKey(tok token.Token) ast.Node {
	if tok.Is(token.Number) {
		return p.tokenToSetElement(tok)
	}
	return p.createNode(&ast.StringLit{Value: tok.Value}, tok.Range)
}

func (p *parser) tokenToSetElement(tok token.Token) ast.Node {
	switch {
	case tok.Is(token.String):
		return p.createNode(&ast.StringLit{Value: tok.Value}, tok.Range)
	case tok.Is(token.Ident):
		return p.createNode(&ast.Ident{Value: tok.Value}, tok.Range)
	case tok.Is(token.Number):
		value := strings.Replace(tok.Value, "_", "", -1)
		valueLower := strings.ToLower(value)
		switch {
		case strings.HasPrefix(valueLower, "0x"), strings.HasPrefix(valueLower, "0b"), strings.HasPrefix(valueLower, "0o"):
			number, err := strconv.ParseInt(value, 0, 64)
			if err != nil {
				p.error("invalid numeric literal: %v", err)
				return nil
			}
			return p.createNode(&ast.IntLit{Value: int(number)}, tok.Range)
		case strings.ContainsAny(valueLower, ".e"):
			number, err := strconv.ParseFloat(value, 64)
			if err != nil {
				p.error("invalid float literal: %v", err)
				return nil
			}
			return p.createNode(&ast.FloatLit{Value: number}, tok.Range)
		default:
			number, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				p.error("invalid integer literal: %v", err)
				return nil
			}
			return p.createNode(&ast.IntLit{Value: int(number)}, tok.Range)
		}
	default:
		p.error("unexpected token in set literal: %v", tok.String())
		return nil
	}
}

func (p *parser) parsePostfixExpression(node ast.Node) ast.Node {
	postfixToken := p.current
	for (postfixToken.Is(token.Operator) || postfixToken.Is(token.Bracket)) && p.err == nil {
		opt := postfixToken.Value == "?."
		if postfixToken.Is(token.Operator, "?") {
			p.next()
			if p.current.Is(token.Bracket, "[") {
				postfixToken = p.current
				opt = true
			} else {
				p.stashed = p.current
				p.hasStash = true
				p.current = postfixToken
				break
			}
		}
	parseToken:
		if postfixToken.Value == "." || postfixToken.Value == "?." {
			p.next()

			propertyToken := p.current
			if opt && propertyToken.Is(token.Bracket, "[") {
				postfixToken = propertyToken
				goto parseToken
			}
			p.next()

			if propertyToken.Kind == token.Template {
				property := p.parseTemplateString(propertyToken)
				if property == nil {
					return nil
				}
				chainNode, isChain := node.(*ast.ChainExpr)
				optional := postfixToken.Value == "?."
				if isChain {
					node = chainNode.Node
				}
				memberNode := p.createMemberNode(&ast.MemberExpr{
					Node:     node,
					Property: property,
					Optional: optional,
				}, propertyToken.Range)
				if memberNode == nil {
					return nil
				}
				if isChain || optional {
					node = p.createNode(&ast.ChainExpr{Node: memberNode}, propertyToken.Range)
					if node == nil {
						return nil
					}
				} else {
					node = memberNode
				}
				postfixToken = p.current
				continue
			}

			if propertyToken.Kind != token.Ident &&
				propertyToken.Kind != token.String &&
				(propertyToken.Kind != token.Operator || !token.IsValidIdentifier(propertyToken.Value)) {
				p.error("expected name")
			}

			property := p.createNode(&ast.StringLit{Value: propertyToken.Value}, propertyToken.Range)
			if property == nil {
				return nil
			}

			chainNode, isChain := node.(*ast.ChainExpr)
			optional := postfixToken.Value == "?."

			if isChain {
				node = chainNode.Node
			}

			memberNode := p.createMemberNode(&ast.MemberExpr{
				Node:     node,
				Property: property,
				Optional: optional,
			}, propertyToken.Range)
			if memberNode == nil {
				return nil
			}

			if p.current.Is(token.Bracket, "(") {
				memberNode.Method = true
				node = p.createNode(&ast.CallExpr{
					Callee:    memberNode,
					Arguments: p.parseArguments([]ast.Node{}),
				}, propertyToken.Range)
				if node == nil {
					return nil
				}
			} else {
				node = memberNode
			}

			if isChain || optional {
				node = p.createNode(&ast.ChainExpr{Node: node}, propertyToken.Range)
				if node == nil {
					return nil
				}
			}

		} else if postfixToken.Value == "[" {
			p.next()

			from := p.parseExpression(0)

			node = p.createNode(&ast.MemberExpr{
				Node:     node,
				Property: from,
				Optional: opt,
			}, postfixToken.Range)
			if node == nil {
				return nil
			}
			if opt {
				node = p.createNode(&ast.ChainExpr{Node: node}, postfixToken.Range)
				if node == nil {
					return nil
				}
			}
			p.expect(token.Bracket, "]")
		} else if postfixToken.Is(token.Bracket, "(") {
			node = p.createNode(&ast.CallExpr{
				Callee:    node,
				Arguments: p.parseArguments([]ast.Node{}),
			}, postfixToken.Range)
			if node == nil {
				return nil
			}
		} else if postfixToken.Is(token.Operator, "++", "--") {
			p.next()
			node = p.createNode(&ast.IncDecStmt{
				Target: node,
				Op:     postfixToken.Value,
				Prefix: false,
			}, postfixToken.Range)
			if node == nil {
				return nil
			}
			break
		} else {
			break
		}
		postfixToken = p.current
	}
	return node
}

func (p *parser) parseComparison(left ast.Node, tok token.Token, precedence int) ast.Node {
	var rootNode ast.Node
	for {
		comparator := p.parseExpression(precedence + 1)
		cmpNode := p.createNode(&ast.BinaryExpr{
			Operator: tok.Value,
			Left:     left,
			Right:    comparator,
		}, tok.Range)
		if cmpNode == nil {
			return nil
		}
		if rootNode == nil {
			rootNode = cmpNode
		} else {
			rootNode = p.createNode(&ast.BinaryExpr{
				Operator: "&&",
				Left:     rootNode,
				Right:    cmpNode,
			}, tok.Range)
			if rootNode == nil {
				return nil
			}
		}

		left = comparator
		tok = p.current
		if !(tok.Is(token.Operator) && token.IsComparison(tok.Value) && p.err == nil) {
			break
		}
		p.next()
	}
	return rootNode
}
