// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package util

import (
	"fmt"
)

// StackClass is the java.util.Stack class object.
var StackClass stackClass

type stackClass struct{}

func (stackClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "new":
		return &Stack{}, nil
	}
	return nil, fmt.Errorf("Stack.%s: undefined", method)
}

// Stack mirrors java.util.Stack.
type Stack struct {
	items []any
}

func (s *Stack) Push(val any) any {
	s.items = append(s.items, val)
	return val
}

func (s *Stack) Pop() (any, error) {
	if len(s.items) == 0 {
		return nil, fmt.Errorf("EmptyStackException")
	}
	last := len(s.items) - 1
	val := s.items[last]
	s.items = s.items[:last]
	return val, nil
}

func (s *Stack) Peek() (any, error) {
	if len(s.items) == 0 {
		return nil, fmt.Errorf("EmptyStackException")
	}
	return s.items[len(s.items)-1], nil
}

func (s *Stack) Empty() bool { return len(s.items) == 0 }

func (s *Stack) Size() int { return len(s.items) }

func (s *Stack) Get(idx int) (any, error) {
	if idx < 0 || idx >= len(s.items) {
		return nil, fmt.Errorf("Stack.get: index %d out of bounds", idx)
	}
	return s.items[idx], nil
}

func (s *Stack) Search(val any) int {
	sv := fmt.Sprintf("%v", val)
	for i := len(s.items) - 1; i >= 0; i-- {
		if fmt.Sprintf("%v", s.items[i]) == sv {
			return len(s.items) - i
		}
	}
	return -1
}

func (s *Stack) Call(method string, args ...any) (any, error) {
	switch method {
	case "push":
		if len(args) != 1 {
			return nil, fmt.Errorf("Stack.push: expected 1 argument")
		}
		return s.Push(args[0]), nil
	case "pop":
		return s.Pop()
	case "peek":
		return s.Peek()
	case "empty", "isEmpty":
		return s.Empty(), nil
	case "size":
		return s.Size(), nil
	case "get":
		if len(args) != 1 {
			return nil, fmt.Errorf("Stack.get: expected 1 argument")
		}
		idx, ok := args[0].(int)
		if !ok {
			return nil, fmt.Errorf("Stack.get: expected int argument")
		}
		return s.Get(idx)
	case "search":
		if len(args) != 1 {
			return nil, fmt.Errorf("Stack.search: expected 1 argument")
		}
		return s.Search(args[0]), nil
	}
	return nil, fmt.Errorf("Stack instance: undefined method %q", method)
}
