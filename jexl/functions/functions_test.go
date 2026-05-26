// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package functions

import (
	"testing"

	"github.com/harness/go-jexl/jexl/vm/runtime"
)

//
// Helpers
//

func mockFn(name string) *runtime.Function {
	return &runtime.Function{Name: name}
}

//
// Tests
//

// Ensure a flat function is found by its registered name.
func TestLookup_registered(t *testing.T) {
	r := New()
	r.Register("abs", mockFn("abs"))
	if _, ok := r.Lookup("abs"); !ok {
		t.Fatal("expected abs to be found")
	}
}

// Ensure Lookup returns false for a name that was
// never registered.
func TestLookup_missing(t *testing.T) {
	r := New()
	if _, ok := r.Lookup("noSuchFunc"); ok {
		t.Fatal("expected missing function to return false")
	}
}

// Ensure Lookup returns the exact function that was registered.
func TestLookup_identity(t *testing.T) {
	r := New()
	fn := mockFn("abs")
	r.Register("abs", fn)
	got, _ := r.Lookup("abs")
	if got != fn {
		t.Fatal("expected Lookup to return the registered function")
	}
}

// Ensure a namespaced function is found by its namespace
// and method name.
func TestLookupNamespace_registered(t *testing.T) {
	r := New()
	r.RegisterNamespace("base64", "encode", mockFn("base64.encode"))
	if _, ok := r.LookupNamespace("base64", "encode"); !ok {
		t.Fatal("expected base64.encode to be found")
	}
}

// Ensure LookupNamespace returns false for an unknown
// namespace.
func TestLookupNamespace_missingNamespace(t *testing.T) {
	r := New()
	if _, ok := r.LookupNamespace("noSuchNs", "method"); ok {
		t.Fatal("expected missing namespace to return false")
	}
}

// Ensure LookupNamespace returns false for an unknown
// method within a known namespace.
func TestLookupNamespace_missingMethod(t *testing.T) {
	r := New()
	r.RegisterNamespace("base64", "encode", mockFn("base64.encode"))
	if _, ok := r.LookupNamespace("base64", "noSuchMethod"); ok {
		t.Fatal("expected missing method to return false")
	}
}

// Ensure LookupNamespace returns the exact function that
// was registered.
func TestLookupNamespace_identity(t *testing.T) {
	r := New()
	fn := mockFn("base64.encode")
	r.RegisterNamespace("base64", "encode", fn)
	got, _ := r.LookupNamespace("base64", "encode")
	if got != fn {
		t.Fatal("expected LookupNamespace to return the registered function")
	}
}

// Ensure HasNamespace returns true for a namespace that
// has at least one registered method.
func TestHasNamespace_present(t *testing.T) {
	r := New()
	r.RegisterNamespace("base64", "encode", mockFn("base64.encode"))
	if !r.HasNamespace("base64") {
		t.Fatal("expected HasNamespace to return true")
	}
}

// Ensure HasNamespace returns false for a namespace that
// was never registered.
func TestHasNamespace_missing(t *testing.T) {
	r := New()
	if r.HasNamespace("noSuchNs") {
		t.Fatal("expected HasNamespace to return false")
	}
}

// Ensure multiple methods can be registered under the
// same namespace and each resolves independently.
func TestRegisterNamespace_multipleMethods(t *testing.T) {
	r := New()
	r.RegisterNamespace("base64", "encode", mockFn("base64.encode"))
	r.RegisterNamespace("base64", "decode", mockFn("base64.decode"))
	if _, ok := r.LookupNamespace("base64", "encode"); !ok {
		t.Fatal("expected base64.encode to be found")
	}
	if _, ok := r.LookupNamespace("base64", "decode"); !ok {
		t.Fatal("expected base64.decode to be found")
	}
}

// Ensure flat and namespaced registrations do not
// interfere with each other.
func TestRegister_flatAndNamespaceCoexist(t *testing.T) {
	r := New()
	r.Register("abs", mockFn("abs"))
	r.RegisterNamespace("math", "abs", mockFn("math.abs"))
	if _, ok := r.Lookup("abs"); !ok {
		t.Fatal("expected flat abs to be found")
	}
	if _, ok := r.LookupNamespace("math", "abs"); !ok {
		t.Fatal("expected math.abs to be found")
	}
}
