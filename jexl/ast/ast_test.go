// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package ast

import (
	"reflect"
	"testing"

	"github.com/harness/go-jexl/jexl/checker/nature"
	"github.com/harness/go-jexl/jexl/token"
)

// Ensure Location round-trips through SetLocation.
func TestBase_Location(t *testing.T) {
	n := &Ident{}
	loc := token.Range{From: 2, To: 5}
	n.SetLocation(loc)
	if n.Location() != loc {
		t.Fatalf("expected %+v, got %+v", loc, n.Location())
	}
}

// Ensure Nature pointer returned by Nature() points to the stored value.
func TestBase_Nature(t *testing.T) {
	n := &IntLit{}
	nat := nature.FromType(reflect.TypeOf(0))
	n.SetNature(nat)
	if n.Nature().Type != nat.Type {
		t.Fatalf("unexpected nature: %+v", n.Nature())
	}
}

// Ensure Type returns anyType when no type has been set.
func TestBase_Type_default(t *testing.T) {
	n := &BoolLit{}
	if n.Type() != anyType {
		t.Fatalf("expected anyType, got %v", n.Type())
	}
}

// Ensure SetType stores the type and Type returns it.
func TestBase_SetType(t *testing.T) {
	n := &FloatLit{}
	t1 := reflect.TypeOf(float64(0))
	n.SetType(t1)
	if n.Type() != t1 {
		t.Fatalf("expected float64 type, got %v", n.Type())
	}
}
