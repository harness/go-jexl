// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package printer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/harness/go-jexl/jexl/ast"
	"github.com/harness/go-jexl/jexl/token"
)

// Print writes node as JEXL source to os.Stdout.
func Print(node ast.Node) error {
	return Fprint(os.Stdout, node)
}

// Fprint writes node as JEXL source to w.
func Fprint(w io.Writer, node ast.Node) error {
	p := &printer{w: w}
	p.node(node)
	return p.err
}

type printer struct {
	w   io.Writer
	err error
}

func (p *printer) write(s string) {
	if p.err != nil {
		return
	}
	_, p.err = io.WriteString(p.w, s)
}

func (p *printer) writef(format string, args ...any) {
	p.write(fmt.Sprintf(format, args...))
}

func (p *printer) node(node ast.Node) {
	if node == nil {
		p.write("nil")
		return
	}
	switch n := node.(type) {
	case *ast.NilLit:
		p.write("nil")

	case *ast.Ident:
		p.write(n.Value)

	case *ast.IntLit:
		p.writef("%d", n.Value)

	case *ast.FloatLit:
		p.writef("%v", n.Value)

	case *ast.BoolLit:
		p.writef("%t", n.Value)

	case *ast.StringLit:
		p.writef("%q", n.Value)

	case *ast.ConstantExpr:
		if n.Value == nil {
			p.write("nil")
		} else {
			b, err := json.Marshal(n.Value)
			if err != nil {
				p.err = err
				return
			}
			p.write(string(b))
		}

	case *ast.UnaryExpr:
		op := n.Operator
		if n.Operator == "not" {
			op = "not "
		}
		wrap := false
		switch b := n.Node.(type) {
		case *ast.BinaryExpr:
			if binaryPrec(b.Operator) < unaryPrec(n.Operator) {
				wrap = true
			}
		case *ast.ConditionalExpr:
			wrap = true
		}
		p.write(op)
		if wrap {
			p.write("(")
			p.node(n.Node)
			p.write(")")
		} else {
			p.node(n.Node)
		}

	case *ast.BinaryExpr:
		p.binaryNode(n)

	case *ast.StrictEqualExpr:
		p.node(n.Left)
		if n.Negated {
			p.write(" !== ")
		} else {
			p.write(" === ")
		}
		p.node(n.Right)

	case *ast.BitExpr:
		if n.Right == nil {
			p.write(n.Operator)
			p.node(n.Left)
		} else {
			p.node(n.Left)
			p.write(" ")
			p.write(n.Operator)
			p.write(" ")
			p.node(n.Right)
		}

	case *ast.LambdaExpr:
		p.write("(")
		p.write(strings.Join(n.Params, ", "))
		p.write(") -> ")
		p.node(n.Body)

	case *ast.ChainExpr:
		p.node(n.Node)

	case *ast.MemberExpr:
		if _, ok := n.Node.(*ast.BinaryExpr); ok {
			p.write("(")
			p.node(n.Node)
			p.write(")")
		} else {
			p.node(n.Node)
		}
		str, isStr := n.Property.(*ast.StringLit)
		if n.Optional {
			if isStr && token.IsValidIdentifier(str.Value) {
				p.write("?.")
				p.write(str.Value)
			} else {
				p.write("?.[")
				p.node(n.Property)
				p.write("]")
			}
		} else {
			if isStr && token.IsValidIdentifier(str.Value) {
				p.write(".")
				p.write(str.Value)
			} else {
				p.write("[")
				p.node(n.Property)
				p.write("]")
			}
		}

	case *ast.CallExpr:
		p.node(n.Callee)
		p.write("(")
		for i, arg := range n.Arguments {
			if i > 0 {
				p.write(", ")
			}
			p.node(arg)
		}
		p.write(")")

	case *ast.PredicateExpr:
		p.node(n.Node)

	case *ast.ConditionalExpr:
		p.conditionalNode(n)

	case *ast.ArrayLit:
		p.write("[")
		for i, elem := range n.Nodes {
			if i > 0 {
				p.write(", ")
			}
			p.node(elem)
		}
		p.write("]")

	case *ast.MapLit:
		p.write("{")
		for i, pair := range n.Pairs {
			if i > 0 {
				p.write(", ")
			}
			p.node(pair)
		}
		p.write("}")

	case *ast.KeyValueExpr:
		if str, ok := n.Key.(*ast.StringLit); ok {
			if token.IsValidIdentifier(str.Value) {
				p.write(str.Value)
			} else {
				p.writef("%q", str.Value)
			}
		} else {
			p.write("(")
			p.node(n.Key)
			p.write(")")
		}
		p.write(": ")
		p.node(n.Value)

	case *ast.RegexLit:
		p.writef("/%s/%s", n.Pattern, n.Flags)

	case *ast.NewExpr:
		p.writef("new %s(", n.ClassName)
		for i, a := range n.Args {
			if i > 0 {
				p.write(", ")
			}
			p.node(a)
		}
		p.write(")")

	case *ast.NamespaceCallExpr:
		p.writef("%s:%s(", n.Namespace, n.Method)
		for i, a := range n.Args {
			if i > 0 {
				p.write(", ")
			}
			p.node(a)
		}
		p.write(")")

	case *ast.EmptyExpr:
		p.write("empty ")
		p.node(n.Value)

	case *ast.SizeExpr:
		p.write("size ")
		p.node(n.Value)

	case *ast.SetLit:
		p.write("{")
		for i, e := range n.Elements {
			if i > 0 {
				p.write(", ")
			}
			p.node(e)
		}
		p.write("}")

	case *ast.BlockStmt:
		for i, s := range n.Statements {
			if i > 0 {
				p.write("; ")
			}
			p.node(s)
		}

	case *ast.VarDecl:
		p.writef("%s %s = ", n.Keyword, n.Name)
		p.node(n.Value)

	case *ast.AssignStmt:
		p.node(n.Target)
		p.writef(" %s ", n.Op)
		p.node(n.Value)

	case *ast.IncDecStmt:
		if n.Prefix {
			p.write(n.Op)
			p.node(n.Target)
		} else {
			p.node(n.Target)
			p.write(n.Op)
		}

	case *ast.ForeachStmt:
		p.write("for (var ")
		p.write(n.Var)
		p.write(" : ")
		p.node(n.Collection)
		p.write(") { ")
		p.node(n.Body)
		p.write(" }")

	case *ast.WhileStmt:
		p.write("while (")
		p.node(n.Cond)
		p.write(") { ")
		p.node(n.Body)
		p.write(" }")

	case *ast.ForStmt:
		p.write("for (")
		if n.Init != nil {
			p.node(n.Init)
		}
		p.write("; ")
		if n.Cond != nil {
			p.node(n.Cond)
		}
		p.write("; ")
		if n.Post != nil {
			p.node(n.Post)
		}
		p.write(") { ")
		p.node(n.Body)
		p.write(" }")

	case *ast.DoWhileStmt:
		p.write("do { ")
		p.node(n.Body)
		p.write(" } while (")
		p.node(n.Cond)
		p.write(")")

	case *ast.BreakStmt:
		p.write("break")

	case *ast.ContinueStmt:
		p.write("continue")

	case *ast.ReturnStmt:
		if n.Value == nil {
			p.write("return")
		} else {
			p.write("return ")
			p.node(n.Value)
		}

	case *ast.TryStmt:
		p.write("try { ")
		p.node(n.Body)
		p.write(" }")
		if n.CatchBody != nil {
			p.writef(" catch (%s) { ", n.CatchVar)
			p.node(n.CatchBody)
			p.write(" }")
		}
		if n.FinallyBody != nil {
			p.write(" finally { ")
			p.node(n.FinallyBody)
			p.write(" }")
		}

	case *ast.ThrowStmt:
		p.write("throw ")
		p.node(n.Value)

	case *ast.SwitchStmt:
		p.write("switch (")
		p.node(n.Subject)
		p.write(") {")
		for _, c := range n.Cases {
			p.write(" case ")
			for i, v := range c.Values {
				if i > 0 {
					p.write(", ")
				}
				p.node(v)
			}
			p.write(": ")
			p.node(c.Body)
		}
		if n.Default != nil {
			p.write(" default: ")
			p.node(n.Default)
		}
		p.write(" }")

	default:
		p.err = fmt.Errorf("printer: unknown node type %T", node)
	}
}

func (p *printer) binaryNode(n *ast.BinaryExpr) {
	if n.Operator == ".." {
		p.node(n.Left)
		p.write("..")
		p.node(n.Right)
		return
	}

	lwrap, rwrap := binaryWrap(n)

	if lwrap {
		p.write("(")
		p.node(n.Left)
		p.write(")")
	} else {
		p.node(n.Left)
	}

	p.write(" ")
	p.write(n.Operator)
	p.write(" ")

	if rwrap {
		p.write("(")
		p.node(n.Right)
		p.write(")")
	} else {
		p.node(n.Right)
	}
}

func (p *printer) conditionalNode(n *ast.ConditionalExpr) {
	if !n.Ternary {
		p.write("if ")
		p.node(n.Cond)
		p.write(" { ")
		p.node(n.Exp1)
		p.write(" }")
		if n.Exp2 == nil {
			return
		}
		if c2, ok := n.Exp2.(*ast.ConditionalExpr); ok && !c2.Ternary {
			p.write(" else ")
			p.conditionalNode(c2)
		} else {
			p.write(" else { ")
			p.node(n.Exp2)
			p.write(" }")
		}
		return
	}

	if _, ok := n.Cond.(*ast.ConditionalExpr); ok {
		p.write("(")
		p.node(n.Cond)
		p.write(")")
	} else {
		p.node(n.Cond)
	}
	p.write(" ? ")
	if _, ok := n.Exp1.(*ast.ConditionalExpr); ok {
		p.write("(")
		p.node(n.Exp1)
		p.write(")")
	} else {
		p.node(n.Exp1)
	}
	p.write(" : ")
	if _, ok := n.Exp2.(*ast.ConditionalExpr); ok {
		p.write("(")
		p.node(n.Exp2)
		p.write(")")
	} else {
		p.node(n.Exp2)
	}
}

// binaryWrap reports whether the left and right children of n need parentheses.
func binaryWrap(n *ast.BinaryExpr) (lwrap, rwrap bool) {
	if l, ok := n.Left.(*ast.UnaryExpr); ok {
		if unaryPrec(l.Operator) < binaryPrec(n.Operator) {
			lwrap = true
		}
	}
	if lb, ok := n.Left.(*ast.BinaryExpr); ok {
		if binaryPrec(lb.Operator) < binaryPrec(n.Operator) {
			lwrap = true
		}
		if binaryPrec(lb.Operator) == binaryPrec(n.Operator) && binaryAssoc(n.Operator) == token.Right {
			lwrap = true
		}
		if lb.Operator == "??" {
			lwrap = true
		}
		if token.IsBoolean(lb.Operator) && n.Operator != lb.Operator {
			lwrap = true
		}
	}
	if _, ok := n.Left.(*ast.ConditionalExpr); ok {
		lwrap = true
	}

	if rb, ok := n.Right.(*ast.BinaryExpr); ok {
		if binaryPrec(rb.Operator) < binaryPrec(n.Operator) {
			rwrap = true
		}
		if binaryPrec(rb.Operator) == binaryPrec(n.Operator) && binaryAssoc(n.Operator) == token.Left {
			rwrap = true
		}
		if token.IsBoolean(rb.Operator) && n.Operator != rb.Operator {
			rwrap = true
		}
	}
	if _, ok := n.Right.(*ast.ConditionalExpr); ok {
		rwrap = true
	}
	return
}

func unaryPrec(op string) int {
	return token.Unary[op].Precedence
}

func binaryPrec(op string) int {
	return token.Binary[op].Precedence
}

func binaryAssoc(op string) token.Associativity {
	if e, ok := token.Binary[op]; ok {
		return e.Associativity
	}
	return token.Left
}
