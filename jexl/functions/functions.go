// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package functions

import "github.com/harness/go-jexl/jexl/vm/runtime"

// Registry maps function names and namespaced methods to their
// runtime.Function implementations.
type Registry struct {
	flat       map[string]*runtime.Function
	namespaced map[string]map[string]*runtime.Function
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{
		flat:       make(map[string]*runtime.Function),
		namespaced: make(map[string]map[string]*runtime.Function),
	}
}

// Register adds a flat function under name.
func (r *Registry) Register(name string, fn *runtime.Function) {
	r.flat[name] = fn
}

// RegisterNamespace adds a function under ns.method.
func (r *Registry) RegisterNamespace(ns, method string, fn *runtime.Function) {
	if r.namespaced[ns] == nil {
		r.namespaced[ns] = make(map[string]*runtime.Function)
	}
	r.namespaced[ns][method] = fn
}

// Lookup returns the flat function registered under name.
func (r *Registry) Lookup(name string) (*runtime.Function, bool) {
	fn, ok := r.flat[name]
	return fn, ok
}

// LookupNamespace returns the function registered under ns.method.
func (r *Registry) LookupNamespace(ns, method string) (*runtime.Function, bool) {
	methods, ok := r.namespaced[ns]
	if !ok {
		return nil, false
	}
	fn, ok := methods[method]
	return fn, ok
}

// HasNamespace reports whether ns has any registered methods.
func (r *Registry) HasNamespace(ns string) bool {
	_, ok := r.namespaced[ns]
	return ok
}
