// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package ast

import (
	"reflect"

	"github.com/harness/go-jexl/jexl/checker/nature"
	"github.com/harness/go-jexl/jexl/token"
)

// Tree is the result of parsing a JEXL expression.
type Tree struct {
	Node   Node       // root of the AST
	Source token.File // the source the tree was parsed from
}

//
// Interfaces
//

// Node is the interface implemented by all AST nodes.
type Node interface {
	Location() token.Range
	SetLocation(token.Range)
	Nature() *nature.Nature
	SetNature(nature.Nature)
	Type() reflect.Type
	SetType(reflect.Type)
}

//
// Base
//

var anyType = reflect.TypeOf(new(any)).Elem()

// base is embedded by every Node implementation.
type base struct {
	loc    token.Range
	nature nature.Nature
}

func (n *base) Location() token.Range       { return n.loc }
func (n *base) SetLocation(loc token.Range) { n.loc = loc }
func (n *base) Nature() *nature.Nature      { return &n.nature }
func (n *base) SetNature(nat nature.Nature) { n.nature = nat }

func (n *base) Type() reflect.Type {
	if n.nature.Type == nil {
		return anyType
	}
	return n.nature.Type
}

func (n *base) SetType(t reflect.Type) {
	n.nature = nature.FromType(t)
}

//
// Implementations
//

// ArrayLit represents an array literal: [a, b, c].
type ArrayLit struct {
	base
	Nodes []Node // element expressions
}

// AssignStmt represents an assignment: target op= value.
type AssignStmt struct {
	base
	Target Node   // assignment target
	Op     string // assignment operator, e.g. "=", "+=", "-="
	Value  Node   // right-hand side value
}

// BinaryExpr represents a binary operation: left op right.
type BinaryExpr struct {
	base
	Operator string // the operator, e.g. "+", "==", "and"
	Left     Node   // left operand
	Right    Node   // right operand
}

// BitExpr represents a bitwise operation: left op right, or unary ~left.
type BitExpr struct {
	base
	Operator string // the operator, e.g. "&", "|", "^", "~", "<<", ">>"
	Left     Node   // left operand (or sole operand for unary ~)
	Right    Node   // right operand; or nil for unary ~
}

// BlockStmt represents a sequence of statements.
type BlockStmt struct {
	base
	Statements []Node // ordered list of statements
}

// BoolLit represents a boolean literal: true or false.
type BoolLit struct {
	base
	Value bool // the literal value
}

// BreakStmt represents a break statement.
type BreakStmt struct {
	base
}

// CallExpr represents a function or method call: callee(args...).
type CallExpr struct {
	base
	Callee    Node   // the function or method being called
	Arguments []Node // argument expressions
}

// CaseClause represents a single case in a switch statement.
type CaseClause struct {
	Values []Node // matched values for this case
	Body   Node   // statements executed when matched
}

// ChainExpr wraps a node in an optional chain context.
type ChainExpr struct {
	base
	Node Node // the chained expression
}

// ConditionalExpr represents a ternary (cond ? a : b) or
// if/else statement.
type ConditionalExpr struct {
	base
	Ternary bool // true for cond ? a : b, false for if/else
	Cond    Node // condition expression
	Exp1    Node // then branch
	Exp2    Node // else branch; or nil for if without else
}

// ConstantExpr holds a pre-evaluated constant value.
type ConstantExpr struct {
	base
	Value any // the constant; may be nil
}

// ContinueStmt represents a continue statement.
type ContinueStmt struct {
	base
}

// DoWhileStmt represents a do { body } while (cond) loop.
type DoWhileStmt struct {
	base
	Body Node // loop body, executed at least once
	Cond Node // continuation condition
}

// EmptyExpr represents the empty() predicate.
type EmptyExpr struct {
	base
	Value Node // expression being tested
}

// FloatLit represents a floating-point literal.
type FloatLit struct {
	base
	Value float64 // the literal value
}

// ForeachStmt represents a for (var x : collection) loop.
type ForeachStmt struct {
	base
	Var        string // loop variable name
	Collection Node   // iterable expression
	Body       Node   // loop body
}

// ForStmt represents a C-style for (init; cond; post) loop.
type ForStmt struct {
	base
	Init Node // initializer statement; or nil
	Cond Node // continuation condition; or nil
	Post Node // post-iteration statement; or nil
	Body Node // loop body
}

// Ident represents a variable or function name.
type Ident struct {
	base
	Value string // the identifier name
}

// IncDecStmt represents a prefix or postfix increment/decrement: ++x, x--.
type IncDecStmt struct {
	base
	Target Node   // the variable being incremented or decremented
	Op     string // "++" or "--"
	Prefix bool   // true for prefix form (++x), false for postfix (x++)
}

// IntLit represents an integer literal.
type IntLit struct {
	base
	Value int // the literal value
}

// LambdaExpr represents a lambda expression: (params) -> body.
type LambdaExpr struct {
	base
	Params []string // parameter names
	Body   Node     // expression or block body
}

// MapLit represents a map literal: {key: value, ...}.
type MapLit struct {
	base
	Pairs []Node // KeyValueExpr elements
}

// MemberExpr represents property access: node.property or node[property].
// Optional is true for ?. and ?.[] forms.
type MemberExpr struct {
	base
	Node     Node // the object being accessed
	Property Node // the property key (StringLit for dot access)
	Optional bool // true for ?. and ?.[] forms
	Method   bool // true when the property resolves to a method
}

// NamespaceCallExpr represents a namespace method call: Namespace:method(args...).
type NamespaceCallExpr struct {
	base
	Namespace string // the namespace name
	Method    string // the method name
	Args      []Node // argument expressions
}

// NewExpr represents a constructor call: new ClassName(args...).
type NewExpr struct {
	base
	ClassName string // the class being instantiated
	Args      []Node // constructor argument expressions
}

// NilLit represents the nil literal.
type NilLit struct {
	base
}

// KeyValueExpr represents a key-value pair in a map literal.
type KeyValueExpr struct {
	base
	Key   Node // key expression
	Value Node // value expression
}

// PredicateExpr wraps a node used as a predicate.
type PredicateExpr struct {
	base
	Node Node // the predicate expression
}

// RegexLit represents a regex literal: /pattern/flags.
type RegexLit struct {
	base
	Pattern string // the regex pattern, without delimiters
	Flags   string // modifier flags, e.g. "gi"
}

// ReturnStmt represents a return statement.
type ReturnStmt struct {
	base
	Value Node // return value; or nil for bare return
}

// SetLit represents a set literal: {a, b, c}.
type SetLit struct {
	base
	Elements []Node // element expressions
}

// SizeExpr represents the size() function.
type SizeExpr struct {
	base
	Value Node // expression whose size is measured
}

// StrictEqualExpr represents === or !==.
type StrictEqualExpr struct {
	base
	Negated bool // true for !==, false for ===
	Left    Node // left operand
	Right   Node // right operand
}

// StringLit represents a string literal.
type StringLit struct {
	base
	Value string // the unquoted string value
}

// SwitchStmt represents a switch statement.
type SwitchStmt struct {
	base
	Subject Node         // the expression being switched on
	Cases   []CaseClause // the case clauses, in source order
	Default Node         // default clause body; or nil
}

// ThrowStmt represents a throw statement.
type ThrowStmt struct {
	base
	Value Node // the value being thrown
}

// TryStmt represents a try/catch/finally block.
type TryStmt struct {
	base
	Body        Node   // the try body
	CatchVar    string // name of the caught exception variable
	CatchBody   Node   // the catch body
	FinallyBody Node   // the finally body; or nil
}

// UnaryExpr represents a unary operation: op node.
type UnaryExpr struct {
	base
	Operator string // the operator, e.g. "!", "-", "not"
	Node     Node   // the operand
}

// VarDecl represents a variable declaration: var/let/const name = value.
type VarDecl struct {
	base
	Name    string // the variable name
	Keyword string // "var", "let", or "const"
	Value   Node   // the initializer expression
}

// WhileStmt represents a while (cond) loop.
type WhileStmt struct {
	base
	Cond Node // continuation condition
	Body Node // loop body
}
