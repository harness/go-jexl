// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package util

import (
	"fmt"
	"strings"
)

// HashSetClass is the java.util.LinkedHashSet class object (insertion-order preserving).
var HashSetClass hashSetClass

// LinkedHashSetClass aliases HashSetClass.
var LinkedHashSetClass = HashSetClass

type hashSetClass struct{}

func (hashSetClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "new":
		return NewHashSet(), nil
	}
	return nil, fmt.Errorf("HashSet.%s: undefined", method)
}

// HashSet mirrors java.util.LinkedHashSet: O(1) membership with insertion-order iteration.
type HashSet struct {
	items map[any]struct{}
	order []any
}

// NewHashSet returns an empty HashSet.
func NewHashSet() *HashSet {
	return &HashSet{items: make(map[any]struct{})}
}

// NewHashSetFrom constructs a HashSet from a slice of elements.
func NewHashSetFrom(elements []any) *HashSet {
	s := NewHashSet()
	for _, e := range elements {
		s.Add(e)
	}
	return s
}

// Add inserts v into the set. Returns true if v was not already present.
func (s *HashSet) Add(v any) bool {
	if _, exists := s.items[v]; exists {
		return false
	}
	s.items[v] = struct{}{}
	s.order = append(s.order, v)
	return true
}

// Remove deletes v from the set. Returns true if v was present.
func (s *HashSet) Remove(v any) bool {
	if _, exists := s.items[v]; !exists {
		return false
	}
	delete(s.items, v)
	order := s.order[:0]
	for _, e := range s.order {
		if e != v {
			order = append(order, e)
		}
	}
	s.order = order
	return true
}

// Contains reports whether v is a member of the set.
func (s HashSet) Contains(v any) bool {
	_, ok := s.items[v]
	return ok
}

// Size returns the number of elements.
func (s HashSet) Size() int { return len(s.items) }

// IsEmpty returns true if the set has no elements.
func (s HashSet) IsEmpty() bool { return len(s.items) == 0 }

// Clear removes all elements.
func (s *HashSet) Clear() {
	s.items = make(map[any]struct{})
	s.order = nil
}

// ToArray returns the elements in insertion order.
func (s HashSet) ToArray() []any {
	out := make([]any, len(s.order))
	copy(out, s.order)
	return out
}

// ToString returns a string like "[a, b, c]".
func (s *HashSet) ToString() string {
	parts := make([]string, len(s.order))
	for i, e := range s.order {
		parts[i] = fmt.Sprintf("%v", e)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// Call dispatches instance methods.
func (s *HashSet) Call(method string, args ...any) (any, error) {
	switch method {
	case "add":
		if len(args) != 1 {
			return nil, fmt.Errorf("HashSet.add: expected 1 argument")
		}
		return s.Add(args[0]), nil
	case "remove":
		if len(args) != 1 {
			return nil, fmt.Errorf("HashSet.remove: expected 1 argument")
		}
		return s.Remove(args[0]), nil
	case "contains":
		if len(args) != 1 {
			return nil, fmt.Errorf("HashSet.contains: expected 1 argument")
		}
		return s.Contains(args[0]), nil
	case "size":
		return s.Size(), nil
	case "isEmpty":
		return s.IsEmpty(), nil
	case "clear":
		s.Clear()
		return nil, nil
	case "toArray":
		return s.ToArray(), nil
	case "toString":
		return s.ToString(), nil
	}
	return nil, fmt.Errorf("HashSet instance: undefined method %q", method)
}
