// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package util

import (
	"testing"
)

// Ensure Stack.new returns a *Stack instance.
func TestStack_newNoArgs(t *testing.T) {
	inst, err := StackClass.Call("new")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if _, ok := inst.(*Stack); !ok {
		t.Errorf("expected *Stack, got %T", inst)
	}
}

// Ensure push adds an element and pop retrieves it.
func TestStack_pushAndPop(t *testing.T) {
	inst, _ := StackClass.Call("new")
	s := inst.(*Stack)
	if _, err := s.Call("push", "a"); err != nil {
		t.Errorf("unexpected error from push: %v", err)
		return
	}
	got, err := s.Call("pop")
	if err != nil {
		t.Errorf("unexpected error from pop: %v", err)
		return
	}
	if got != "a" {
		t.Errorf("expected \"a\", got %v", got)
	}
}

// Ensure peek returns the top element without removing it.
func TestStack_peek(t *testing.T) {
	inst, _ := StackClass.Call("new")
	s := inst.(*Stack)
	s.Call("push", "b")
	got, err := s.Call("peek")
	if err != nil {
		t.Errorf("unexpected error from peek: %v", err)
		return
	}
	if got != "b" {
		t.Errorf("expected \"b\", got %v", got)
	}
	if s.Size() != 1 {
		t.Errorf("expected size 1 after peek, got %d", s.Size())
	}
}

// Ensure empty returns true for a new stack and false after a push.
func TestStack_empty(t *testing.T) {
	inst, _ := StackClass.Call("new")
	s := inst.(*Stack)
	got, err := s.Call("empty")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true for new stack, got %v", got)
	}
	s.Call("push", "item")
	got, err = s.Call("empty")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != false {
		t.Errorf("expected false after push, got %v", got)
	}
}

// Ensure search returns 2 (1-based from top) for the middle element of a 3-element stack.
func TestStack_search(t *testing.T) {
	inst, _ := StackClass.Call("new")
	s := inst.(*Stack)
	s.Call("push", "a")
	s.Call("push", "b")
	s.Call("push", "c")
	got, err := s.Call("search", "b")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 2 {
		t.Errorf("expected 2, got %v", got)
	}
}

// Ensure pop on an empty stack returns an error.
func TestStack_popEmpty(t *testing.T) {
	inst, _ := StackClass.Call("new")
	s := inst.(*Stack)
	if _, err := s.Call("pop"); err == nil {
		t.Error("expected error when popping empty stack")
	}
}

// Ensure an unknown class method returns an error.
func TestStack_unknownClassMethod(t *testing.T) {
	if _, err := StackClass.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown class method")
	}
}

// Ensure an unknown method not on Stack or ArrayList returns an error.
func TestStack_unknownInstanceMethod(t *testing.T) {
	inst, _ := StackClass.Call("new")
	s := inst.(*Stack)
	if _, err := s.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown instance method")
	}
}
