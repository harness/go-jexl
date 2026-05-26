// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package ast

import (
	"testing"
)

// unknownNode is a Node implementation not handled by Walk,
// used to exercise the panic branch.
type unknownNode struct{ base }

// collector accumulates visited nodes in order.
type collector struct {
	nodes []Node
}

func (c *collector) Visit(n *Node) {
	c.nodes = append(c.nodes, *n)
}

// Ensure Walk on a nil node does nothing.
func TestWalk_nil(t *testing.T) {
	var n Node
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 0 {
		t.Fatalf("expected 0 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk visits a leaf node exactly once.
func TestWalk_leaf(t *testing.T) {
	var n Node = &IntLit{Value: 1}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 1 {
		t.Fatalf("expected 1 visit, got %d", len(c.nodes))
	}
}

// Ensure Walk visits children before the parent (post-order).
func TestWalk_postOrder(t *testing.T) {
	left := &IntLit{Value: 1}
	right := &IntLit{Value: 2}
	var n Node = &BinaryExpr{Operator: "+", Left: left, Right: right}
	c := &collector{}
	Walk(&n, c)
	// order: left, right, parent
	if len(c.nodes) != 3 {
		t.Fatalf("expected 3 visits, got %d", len(c.nodes))
	}
	if c.nodes[0] != Node(left) {
		t.Fatal("expected left child first")
	}
	if c.nodes[1] != Node(right) {
		t.Fatal("expected right child second")
	}
	if c.nodes[2] != n {
		t.Fatal("expected parent last")
	}
}

// Ensure Walk visits all array elements then the array node.
func TestWalk_arrayLit(t *testing.T) {
	a := &IntLit{Value: 1}
	b := &IntLit{Value: 2}
	c := &IntLit{Value: 3}
	var n Node = &ArrayLit{Nodes: []Node{a, b, c}}
	col := &collector{}
	Walk(&n, col)
	if len(col.nodes) != 4 {
		t.Fatalf("expected 4 visits, got %d", len(col.nodes))
	}
	if col.nodes[3] != n {
		t.Fatal("expected ArrayLit last")
	}
}

// Ensure Walk handles UnaryExpr.
func TestWalk_unaryExpr(t *testing.T) {
	child := &BoolLit{Value: true}
	var n Node = &UnaryExpr{Operator: "!", Node: child}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 2 {
		t.Fatalf("expected 2 visits, got %d", len(c.nodes))
	}
	if c.nodes[0] != Node(child) {
		t.Fatal("expected child visited first")
	}
}

// Ensure Walk handles BlockStmt with multiple statements.
func TestWalk_blockStmt(t *testing.T) {
	s1 := &BreakStmt{}
	s2 := &ContinueStmt{}
	var n Node = &BlockStmt{Statements: []Node{s1, s2}}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 3 {
		t.Fatalf("expected 3 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles ConditionalExpr visiting all three children.
func TestWalk_conditionalExpr(t *testing.T) {
	cond := &BoolLit{Value: true}
	exp1 := &IntLit{Value: 1}
	exp2 := &IntLit{Value: 2}
	var n Node = &ConditionalExpr{Cond: cond, Exp1: exp1, Exp2: exp2}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 4 {
		t.Fatalf("expected 4 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles ReturnStmt with nil value.
func TestWalk_returnStmt_nilValue(t *testing.T) {
	var n Node = &ReturnStmt{}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 1 {
		t.Fatalf("expected 1 visit, got %d", len(c.nodes))
	}
}

// Ensure Walk handles ReturnStmt with a value.
func TestWalk_returnStmt_withValue(t *testing.T) {
	val := &IntLit{Value: 0}
	var n Node = &ReturnStmt{Value: val}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 2 {
		t.Fatalf("expected 2 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles ForStmt with nil init/cond/post.
func TestWalk_forStmt_nilOptionals(t *testing.T) {
	body := &BreakStmt{}
	var n Node = &ForStmt{Body: body}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 2 {
		t.Fatalf("expected 2 visits (body + ForStmt), got %d", len(c.nodes))
	}
}

// Ensure Walk handles BitExpr with nil right (unary).
func TestWalk_bitExpr_unary(t *testing.T) {
	operand := &IntLit{Value: 5}
	var n Node = &BitExpr{Operator: "~", Left: operand}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 2 {
		t.Fatalf("expected 2 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles MapLit visiting all pairs.
func TestWalk_mapLit(t *testing.T) {
	k := &StringLit{Value: "a"}
	v := &IntLit{Value: 1}
	pair := &KeyValueExpr{Key: k, Value: v}
	var n Node = &MapLit{Pairs: []Node{pair}}
	c := &collector{}
	Walk(&n, c)
	// key, value, pair, map = 4
	if len(c.nodes) != 4 {
		t.Fatalf("expected 4 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles SwitchStmt visiting subject, case
// values, case bodies, and default.
func TestWalk_switchStmt(t *testing.T) {
	subject := &Ident{Value: "x"}
	val := &IntLit{Value: 1}
	body := &BreakStmt{}
	def := &ContinueStmt{}
	var n Node = &SwitchStmt{
		Subject: subject,
		Cases:   []CaseClause{{Values: []Node{val}, Body: body}},
		Default: def,
	}
	c := &collector{}
	Walk(&n, c)
	// subject, val, body, default, switch = 5
	if len(c.nodes) != 5 {
		t.Fatalf("expected 5 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles NilLit as a leaf.
func TestWalk_nilLit(t *testing.T) {
	var n Node = &NilLit{}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 1 {
		t.Fatalf("expected 1 visit, got %d", len(c.nodes))
	}
}

// Ensure Walk handles Ident as a leaf.
func TestWalk_ident(t *testing.T) {
	var n Node = &Ident{Value: "x"}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 1 {
		t.Fatalf("expected 1 visit, got %d", len(c.nodes))
	}
}

// Ensure Walk handles FloatLit as a leaf.
func TestWalk_floatLit(t *testing.T) {
	var n Node = &FloatLit{Value: 1.5}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 1 {
		t.Fatalf("expected 1 visit, got %d", len(c.nodes))
	}
}

// Ensure Walk handles StringLit as a leaf.
func TestWalk_stringLit(t *testing.T) {
	var n Node = &StringLit{Value: "hello"}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 1 {
		t.Fatalf("expected 1 visit, got %d", len(c.nodes))
	}
}

// Ensure Walk handles ConstantExpr as a leaf.
func TestWalk_constantExpr(t *testing.T) {
	var n Node = &ConstantExpr{Value: 42}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 1 {
		t.Fatalf("expected 1 visit, got %d", len(c.nodes))
	}
}

// Ensure Walk handles RegexLit as a leaf.
func TestWalk_regexLit(t *testing.T) {
	var n Node = &RegexLit{Pattern: "foo", Flags: "i"}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 1 {
		t.Fatalf("expected 1 visit, got %d", len(c.nodes))
	}
}

// Ensure Walk handles BreakStmt as a leaf.
func TestWalk_breakStmt(t *testing.T) {
	var n Node = &BreakStmt{}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 1 {
		t.Fatalf("expected 1 visit, got %d", len(c.nodes))
	}
}

// Ensure Walk handles ContinueStmt as a leaf.
func TestWalk_continueStmt(t *testing.T) {
	var n Node = &ContinueStmt{}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 1 {
		t.Fatalf("expected 1 visit, got %d", len(c.nodes))
	}
}

// Ensure Walk handles StrictEqualExpr visiting both children.
func TestWalk_strictEqualExpr(t *testing.T) {
	l := &IntLit{Value: 1}
	r := &IntLit{Value: 1}
	var n Node = &StrictEqualExpr{Left: l, Right: r}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 3 {
		t.Fatalf("expected 3 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles BitExpr binary form visiting both children.
func TestWalk_bitExpr_binary(t *testing.T) {
	l := &IntLit{Value: 1}
	r := &IntLit{Value: 2}
	var n Node = &BitExpr{Operator: "&", Left: l, Right: r}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 3 {
		t.Fatalf("expected 3 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles LambdaExpr visiting the body.
func TestWalk_lambdaExpr(t *testing.T) {
	body := &Ident{Value: "x"}
	var n Node = &LambdaExpr{Params: []string{"x"}, Body: body}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 2 {
		t.Fatalf("expected 2 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles ChainExpr visiting the inner node.
func TestWalk_chainExpr(t *testing.T) {
	inner := &Ident{Value: "a"}
	var n Node = &ChainExpr{Node: inner}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 2 {
		t.Fatalf("expected 2 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles MemberExpr visiting node and property.
func TestWalk_memberExpr(t *testing.T) {
	obj := &Ident{Value: "a"}
	prop := &StringLit{Value: "b"}
	var n Node = &MemberExpr{Node: obj, Property: prop}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 3 {
		t.Fatalf("expected 3 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles CallExpr visiting callee and all arguments.
func TestWalk_callExpr(t *testing.T) {
	callee := &Ident{Value: "f"}
	arg1 := &IntLit{Value: 1}
	arg2 := &IntLit{Value: 2}
	var n Node = &CallExpr{Callee: callee, Arguments: []Node{arg1, arg2}}
	c := &collector{}
	Walk(&n, c)
	// callee, arg1, arg2, call = 4
	if len(c.nodes) != 4 {
		t.Fatalf("expected 4 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles PredicateExpr visiting the inner node.
func TestWalk_predicateExpr(t *testing.T) {
	inner := &BoolLit{Value: true}
	var n Node = &PredicateExpr{Node: inner}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 2 {
		t.Fatalf("expected 2 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles VarDecl visiting the value.
func TestWalk_varDecl(t *testing.T) {
	val := &IntLit{Value: 0}
	var n Node = &VarDecl{Name: "i", Value: val}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 2 {
		t.Fatalf("expected 2 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles AssignStmt visiting target and value.
func TestWalk_assignStmt(t *testing.T) {
	target := &Ident{Value: "x"}
	val := &IntLit{Value: 1}
	var n Node = &AssignStmt{Target: target, Op: "=", Value: val}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 3 {
		t.Fatalf("expected 3 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles IncDecStmt visiting the target.
func TestWalk_incDecStmt(t *testing.T) {
	target := &Ident{Value: "i"}
	var n Node = &IncDecStmt{Target: target, Op: "++"}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 2 {
		t.Fatalf("expected 2 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles ForeachStmt visiting collection and body.
func TestWalk_foreachStmt(t *testing.T) {
	coll := &Ident{Value: "items"}
	body := &BreakStmt{}
	var n Node = &ForeachStmt{Var: "x", Collection: coll, Body: body}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 3 {
		t.Fatalf("expected 3 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles WhileStmt visiting cond and body.
func TestWalk_whileStmt(t *testing.T) {
	cond := &BoolLit{Value: true}
	body := &BreakStmt{}
	var n Node = &WhileStmt{Cond: cond, Body: body}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 3 {
		t.Fatalf("expected 3 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles ForStmt with all fields set.
func TestWalk_forStmt_allFields(t *testing.T) {
	init := &VarDecl{Name: "i", Value: &IntLit{Value: 0}}
	cond := &BoolLit{Value: true}
	post := &IncDecStmt{Target: &Ident{Value: "i"}, Op: "++"}
	body := &BreakStmt{}
	var n Node = &ForStmt{Init: init, Cond: cond, Post: post, Body: body}
	c := &collector{}
	Walk(&n, c)
	// init.Value, init, cond, body, post.Target, post, ForStmt = 7
	if len(c.nodes) != 7 {
		t.Fatalf("expected 7 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles DoWhileStmt visiting body then cond.
func TestWalk_doWhileStmt(t *testing.T) {
	body := &BreakStmt{}
	cond := &BoolLit{Value: false}
	var n Node = &DoWhileStmt{Body: body, Cond: cond}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 3 {
		t.Fatalf("expected 3 visits, got %d", len(c.nodes))
	}
	if c.nodes[0] != Node(body) {
		t.Fatal("expected body visited before cond")
	}
}

// Ensure Walk handles TryStmt visiting body and catch body.
func TestWalk_tryStmt(t *testing.T) {
	body := &BreakStmt{}
	catchBody := &ContinueStmt{}
	var n Node = &TryStmt{Body: body, CatchBody: catchBody}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 3 {
		t.Fatalf("expected 3 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles ThrowStmt visiting the value.
func TestWalk_throwStmt(t *testing.T) {
	val := &StringLit{Value: "err"}
	var n Node = &ThrowStmt{Value: val}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 2 {
		t.Fatalf("expected 2 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles EmptyExpr visiting the inner value.
func TestWalk_emptyExpr(t *testing.T) {
	val := &Ident{Value: "x"}
	var n Node = &EmptyExpr{Value: val}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 2 {
		t.Fatalf("expected 2 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles SizeExpr visiting the inner value.
func TestWalk_sizeExpr(t *testing.T) {
	val := &Ident{Value: "x"}
	var n Node = &SizeExpr{Value: val}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 2 {
		t.Fatalf("expected 2 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles SetLit visiting all elements.
func TestWalk_setLit(t *testing.T) {
	a := &IntLit{Value: 1}
	b := &IntLit{Value: 2}
	var n Node = &SetLit{Elements: []Node{a, b}}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 3 {
		t.Fatalf("expected 3 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles NewExpr visiting all args.
func TestWalk_newExpr(t *testing.T) {
	arg := &StringLit{Value: "x"}
	var n Node = &NewExpr{ClassName: "Foo", Args: []Node{arg}}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 2 {
		t.Fatalf("expected 2 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk handles NamespaceCallExpr visiting all args.
func TestWalk_namespaceCallExpr(t *testing.T) {
	arg := &IntLit{Value: 1}
	var n Node = &NamespaceCallExpr{Namespace: "math", Method: "abs", Args: []Node{arg}}
	c := &collector{}
	Walk(&n, c)
	if len(c.nodes) != 2 {
		t.Fatalf("expected 2 visits, got %d", len(c.nodes))
	}
}

// Ensure Walk panics on an unknown node type.
func TestWalk_unknownNode_panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unknown node type")
		}
	}()
	var n Node = &unknownNode{}
	Walk(&n, &collector{})
}

// Ensure Walk handles SwitchStmt with no default clause.
func TestWalk_switchStmt_noDefault(t *testing.T) {
	subject := &Ident{Value: "x"}
	val := &IntLit{Value: 1}
	body := &BreakStmt{}
	var n Node = &SwitchStmt{
		Subject: subject,
		Cases:   []CaseClause{{Values: []Node{val}, Body: body}},
	}
	c := &collector{}
	Walk(&n, c)
	// subject, val, body, switch = 4
	if len(c.nodes) != 4 {
		t.Fatalf("expected 4 visits, got %d", len(c.nodes))
	}
}
