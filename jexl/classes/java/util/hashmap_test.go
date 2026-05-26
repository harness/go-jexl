// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package util

import (
	"testing"
)

// Ensure HashMap.new returns a HashMap instance.
func TestHashMap_newNoArgs(t *testing.T) {
	inst, err := HashMapClass.Call("new")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if _, ok := inst.(HashMap); !ok {
		t.Errorf("expected HashMap, got %T", inst)
	}
}

// Ensure put stores a value and get retrieves it by key.
func TestHashMap_putAndGet(t *testing.T) {
	inst, _ := HashMapClass.Call("new")
	m := inst.(HashMap)
	if _, err := m.Call("put", "key", "val"); err != nil {
		t.Errorf("unexpected error from put: %v", err)
		return
	}
	got, err := m.Call("get", "key")
	if err != nil {
		t.Errorf("unexpected error from get: %v", err)
		return
	}
	if got != "val" {
		t.Errorf("expected \"val\", got %v", got)
	}
}

// Ensure size returns 2 after putting 2 distinct entries.
func TestHashMap_size(t *testing.T) {
	inst, _ := HashMapClass.Call("new")
	m := inst.(HashMap)
	m.Call("put", "a", 1)
	m.Call("put", "b", 2)
	got, err := m.Call("size")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != 2 {
		t.Errorf("expected 2, got %v", got)
	}
}

// Ensure containsKey returns true for an existing key and false for a missing one.
func TestHashMap_containsKey(t *testing.T) {
	inst, _ := HashMapClass.Call("new")
	m := inst.(HashMap)
	m.Call("put", "x", 42)
	got, err := m.Call("containsKey", "x")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
	got, err = m.Call("containsKey", "y")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != false {
		t.Errorf("expected false, got %v", got)
	}
}

// Ensure remove returns the old value and reduces the size to 0.
func TestHashMap_remove(t *testing.T) {
	inst, _ := HashMapClass.Call("new")
	m := inst.(HashMap)
	m.Call("put", "k", "oldval")
	got, err := m.Call("remove", "k")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "oldval" {
		t.Errorf("expected \"oldval\", got %v", got)
	}
	if m.Size() != 0 {
		t.Errorf("expected size 0, got %d", m.Size())
	}
}

// Ensure isEmpty returns true for a new map and false after a put.
func TestHashMap_isEmpty(t *testing.T) {
	inst, _ := HashMapClass.Call("new")
	m := inst.(HashMap)
	got, err := m.Call("isEmpty")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != true {
		t.Errorf("expected true for new map, got %v", got)
	}
	m.Call("put", "a", 1)
	got, err = m.Call("isEmpty")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != false {
		t.Errorf("expected false after put, got %v", got)
	}
}

// Ensure getOrDefault returns the default when the key is absent.
func TestHashMap_getOrDefault(t *testing.T) {
	inst, _ := HashMapClass.Call("new")
	m := inst.(HashMap)
	got, err := m.Call("getOrDefault", "missing", "default")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got != "default" {
		t.Errorf("expected \"default\", got %v", got)
	}
}

// Ensure an unknown class method returns an error.
func TestHashMap_unknownClassMethod(t *testing.T) {
	if _, err := HashMapClass.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown class method")
	}
}

// Ensure an unknown instance method returns an error.
func TestHashMap_unknownInstanceMethod(t *testing.T) {
	inst, _ := HashMapClass.Call("new")
	m := inst.(HashMap)
	if _, err := m.Call("noSuchMethod"); err == nil {
		t.Error("expected error for unknown instance method")
	}
}
