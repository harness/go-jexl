// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package ast

import "fmt"

// Visitor is implemented by any type that walks AST nodes.
type Visitor interface {
	Visit(node *Node)
}

// Walk traverses the AST rooted at node in depth-first order,
// calling v.Visit on each node after its children have been visited.
func Walk(node *Node, v Visitor) {
	if *node == nil {
		return
	}
	switch n := (*node).(type) {
	case *NilLit:
	case *Ident:
	case *IntLit:
	case *FloatLit:
	case *BoolLit:
	case *StringLit:
	case *ConstantExpr:
	case *UnaryExpr:
		Walk(&n.Node, v)
	case *BinaryExpr:
		Walk(&n.Left, v)
		Walk(&n.Right, v)
	case *StrictEqualExpr:
		Walk(&n.Left, v)
		Walk(&n.Right, v)
	case *BitExpr:
		Walk(&n.Left, v)
		if n.Right != nil {
			Walk(&n.Right, v)
		}
	case *LambdaExpr:
		Walk(&n.Body, v)
	case *ChainExpr:
		Walk(&n.Node, v)
	case *MemberExpr:
		Walk(&n.Node, v)
		Walk(&n.Property, v)
	case *CallExpr:
		Walk(&n.Callee, v)
		for i := range n.Arguments {
			Walk(&n.Arguments[i], v)
		}
	case *PredicateExpr:
		Walk(&n.Node, v)
	case *ConditionalExpr:
		Walk(&n.Cond, v)
		Walk(&n.Exp1, v)
		Walk(&n.Exp2, v)
	case *ArrayLit:
		for i := range n.Nodes {
			Walk(&n.Nodes[i], v)
		}
	case *MapLit:
		for i := range n.Pairs {
			Walk(&n.Pairs[i], v)
		}
	case *KeyValueExpr:
		Walk(&n.Key, v)
		Walk(&n.Value, v)
	case *BlockStmt:
		for i := range n.Statements {
			Walk(&n.Statements[i], v)
		}
	case *VarDecl:
		Walk(&n.Value, v)
	case *AssignStmt:
		Walk(&n.Target, v)
		Walk(&n.Value, v)
	case *IncDecStmt:
		Walk(&n.Target, v)
	case *ForeachStmt:
		Walk(&n.Collection, v)
		Walk(&n.Body, v)
	case *WhileStmt:
		Walk(&n.Cond, v)
		Walk(&n.Body, v)
	case *ForStmt:
		if n.Init != nil {
			Walk(&n.Init, v)
		}
		if n.Cond != nil {
			Walk(&n.Cond, v)
		}
		Walk(&n.Body, v)
		if n.Post != nil {
			Walk(&n.Post, v)
		}
	case *DoWhileStmt:
		Walk(&n.Body, v)
		Walk(&n.Cond, v)
	case *BreakStmt:
	case *ContinueStmt:
	case *ReturnStmt:
		if n.Value != nil {
			Walk(&n.Value, v)
		}
	case *TryStmt:
		Walk(&n.Body, v)
		Walk(&n.CatchBody, v)
	case *ThrowStmt:
		Walk(&n.Value, v)
	case *EmptyExpr:
		Walk(&n.Value, v)
	case *SizeExpr:
		Walk(&n.Value, v)
	case *SetLit:
		for i := range n.Elements {
			Walk(&n.Elements[i], v)
		}
	case *SwitchStmt:
		Walk(&n.Subject, v)
		for i := range n.Cases {
			for j := range n.Cases[i].Values {
				Walk(&n.Cases[i].Values[j], v)
			}
			Walk(&n.Cases[i].Body, v)
		}
		if n.Default != nil {
			Walk(&n.Default, v)
		}
	case *RegexLit:
	case *NewExpr:
		for i := range n.Args {
			Walk(&n.Args[i], v)
		}
	case *NamespaceCallExpr:
		for i := range n.Args {
			Walk(&n.Args[i], v)
		}
	default:
		panic(fmt.Sprintf("undefined node type (%T)", node))
	}

	v.Visit(node)
}
