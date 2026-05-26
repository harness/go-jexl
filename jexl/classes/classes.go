// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package classes

import "strings"

// Object is the single interface implemented by all
// classes and class instances.
type Object interface {
	Call(method string, args ...any) (any, error)
}

// Registry maps class names to their Object implementations.
type Registry struct {
	items      map[string]Object
	classNames map[string]struct{}
}

// New returns an empty Registry.
func New() *Registry {
	r := new(Registry)
	r.items = make(map[string]Object)
	r.classNames = make(map[string]struct{})
	return r
}

// Register adds a class under the given name.
func (r *Registry) Register(name string, obj Object) *Registry {
	r.items[name] = obj
	if strings.Contains(name, ".") {
		r.classNames[name] = struct{}{}
		short := name[strings.LastIndexByte(name, '.')+1:]
		if _, exists := r.items[short]; !exists {
			r.items[short] = obj
			r.classNames[short] = struct{}{}
		}
	} else {
		delete(r.classNames, name)
	}
	return r
}

// Lookup returns the Object registered under name.
func (r *Registry) Lookup(name string) (Object, bool) {
	obj, ok := r.items[name]
	return obj, ok
}

// LookupClass returns the Object registered under name, if name is a class.
func (r *Registry) LookupClass(name string) (Object, bool) {
	if _, ok := r.classNames[name]; !ok {
		return nil, false
	}
	obj, ok := r.items[name]
	return obj, ok
}
