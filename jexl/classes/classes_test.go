// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package classes

import (
	"fmt"
	"testing"
)

//
// Mocks
//

// mock class
type mockClass struct {
	name string
}

func (s mockClass) Call(method string, args ...any) (any, error) {
	if method == "new" {
		return &mockInstance{class: s.name}, nil
	} else {
		return nil, fmt.Errorf("%s: unknown method %q", s.name, method)
	}
}

// mock class instance
type mockInstance struct {
	class string
}

func (s *mockInstance) Call(method string, args ...any) (any, error) {
	if method == "name" {
		return s.class, nil
	} else {
		return nil, fmt.Errorf("instance: unknown method %q", method)
	}
}

//
// Tests
//

// Ensure a class registered under a fully-qualified name
// is found by that name.
func TestLookup_fullyQualifiedName(t *testing.T) {
	r := New()
	r.Register("java.util.ArrayList", mockClass{"ArrayList"})
	if _, ok := r.Lookup("java.util.ArrayList"); !ok {
		t.Error("expected fully-qualified name to resolve")
	}
}

// Ensure a class registered under a fully-qualified
// name is also found by its short name.
func TestLookup_shortNameAlias(t *testing.T) {
	r := New()
	obj := mockClass{"ArrayList"}
	r.Register("java.util.ArrayList", obj)
	got, ok := r.Lookup("ArrayList")
	if !ok {
		t.Error("expected short name to resolve")
	}
	if got != obj {
		t.Error("expected short name to return same object as full name")
	}
}

// Ensure a short name registered directly is not
// overwritten by a later fully-qualified registration
// that shares the same short segment.
func TestLookup_shortNameNotOverwritten(t *testing.T) {
	r := New()
	first := mockClass{"first"}
	r.Register("ArrayList", first)
	r.Register("java.util.ArrayList", mockClass{"second"})
	got, _ := r.Lookup("ArrayList")
	if got != first {
		t.Error("expected original short-name registration to be preserved")
	}
}

// Ensure Lookup returns false for a name that was never
// registered.
func TestLookup_missing(t *testing.T) {
	r := New()
	if _, ok := r.Lookup("NoSuchClass"); ok {
		t.Error("expected missing class to return false")
	}
}

// Ensure calling new on a registered class returns an
// instance that itself implements Object and can dispatch
// further method calls.
func TestRegister_constructAndDispatch(t *testing.T) {
	r := New()
	r.Register("ArrayList", mockClass{"ArrayList"})

	cls, _ := r.Lookup("ArrayList")
	inst, err := cls.Call("new", nil)
	if err != nil {
		t.Errorf("construction failed: %v", err)
		return
	}
	obj, ok := inst.(Object)
	if !ok {
		t.Error("expected instance to implement Object")
		return
	}
	name, err := obj.Call("name", nil)
	if err != nil {
		t.Errorf("dispatch failed: %v", err)
		return
	}
	if name != "ArrayList" {
		t.Errorf("expected ArrayList, got %v", name)
	}
}
