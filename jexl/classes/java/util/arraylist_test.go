// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package util

import (
	"testing"
)

// Ensure ArrayList.new returns a *ArrayList instance.
func TestArrayList_newNoArgs(t *testing.T) {
	inst, err := ArrayListClass.Call("new")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if _, ok := inst.(*ArrayList); !ok {
		t.Errorf("expected *ArrayList, got %T", inst)
	}
}

// Ensure add and get work correctly for a single element.
func TestArrayList_addAndGet(t *testing.T) {
	inst, err := ArrayListClass.Call("new")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	l := inst.(*ArrayList)
	if _, err := l.Call("add", "a"); err != nil {
		t.Errorf("unexpected error from add: %v", err)
		return
	}
	got, err := l.Call("get", 0)
	if err != nil {
		t.Errorf("unexpected error from get: %v", err)
		return
	}
	if got != "a" {
		t.Errorf("expected \"a\", got %v", got)
	}
}

// Ensure size returns 3 after adding 3 items.
func TestArrayList_size(t *testing.T) {
	inst, _ := ArrayListClass.Call("new")
	l := inst.(*ArrayList)
	l.Call("add", 1)
	l.Call("add", 2)
	l.Call("add", 3)
	got, err := l.Call("size")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 3 {
		t.Errorf("expected 3, got %v", got)
	}
}

// Ensure remove returns the removed element and shrinks the list.
func TestArrayList_remove(t *testing.T) {
	inst, _ := ArrayListClass.Call("new")
	l := inst.(*ArrayList)
	l.Call("add", "x")
	got, err := l.Call("remove", 0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "x" {
		t.Errorf("expected \"x\", got %v", got)
	}
	if l.Size() != 0 {
		t.Errorf("expected size 0, got %d", l.Size())
	}
}

// Ensure contains returns true for an element that was added.
func TestArrayList_contains(t *testing.T) {
	inst, _ := ArrayListClass.Call("new")
	l := inst.(*ArrayList)
	l.Call("add", "hello")
	got, err := l.Call("contains", "hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

// Ensure isEmpty returns true for a new list and false after adding an item.
func TestArrayList_isEmpty(t *testing.T) {
	inst, _ := ArrayListClass.Call("new")
	l := inst.(*ArrayList)
	got, err := l.Call("isEmpty")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true for empty list, got %v", got)
	}
	l.Call("add", "item")
	got, err = l.Call("isEmpty")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != false {
		t.Errorf("expected false after add, got %v", got)
	}
}

// Ensure toString returns a bracket-wrapped string after adding items.
func TestArrayList_toString(t *testing.T) {
	inst, _ := ArrayListClass.Call("new")
	l := inst.(*ArrayList)
	l.Call("add", 1)
	l.Call("add", 2)
	if l.Size() != 2 {
		t.Errorf("expected size 2, got %d", l.Size())
		return
	}
	got, err := l.Call("toString")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	s, ok := got.(string)
	if !ok {
		t.Errorf("expected string, got %T", got)
		return
	}
	if len(s) == 0 || s[0] != '[' {
		t.Errorf("expected string starting with '[', got %q", s)
	}
}

// Ensure an unknown class method returns an error.
func TestArrayList_unknownClassMethod(t *testing.T) {
	if _, err := ArrayListClass.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown class method")
	}
}

// Ensure an unknown instance method returns an error.
func TestArrayList_unknownInstanceMethod(t *testing.T) {
	inst, _ := ArrayListClass.Call("new")
	l := inst.(*ArrayList)
	if _, err := l.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown instance method")
	}
}
