// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package deref

import (
	"reflect"
	"testing"
)

// Ensure a struct value is returned as-is (no pointer chain to follow).
func TestInterface_structValue(t *testing.T) {
	type point struct{ x, y int }
	p := point{1, 2}
	got, ok := Interface(p).(point)
	if !ok || got != p {
		t.Fatalf("expected %v, got %v", p, got)
	}
}

// Ensure a bare nil (untyped) is handled without
// panicking or falling through to the reflect code.
func TestInterface_nil(t *testing.T) {
	if Interface(nil) != nil {
		t.Fatal("expected nil")
	}
}

// Ensure a typed nil pointer returns nil.
func TestInterface_nilPointer(t *testing.T) {
	var p *int
	if Interface(p) != nil {
		t.Fatal("expected nil for nil pointer")
	}
}

// Ensure a plain non-pointer value is returned as-is.
func TestInterface_plainValue(t *testing.T) {
	if Interface(42) != 42 {
		t.Fatal("expected 42")
	}
}

// Ensure a single pointer is dereferenced to its
// underlying value.
func TestInterface_singlePointer(t *testing.T) {
	v := 42
	if Interface(&v) != 42 {
		t.Fatal("expected 42")
	}
}

// Ensure a double pointer is fully dereferenced to its
// underlying value.
func TestInterface_doublePointer(t *testing.T) {
	v := 42
	p := &v
	if Interface(&p) != 42 {
		t.Fatal("expected 42")
	}
}

// Ensure a double pointer with a nil inner pointer
// returns nil.
func TestInterface_nilDoublePointer(t *testing.T) {
	var p *int
	pp := &p
	if Interface(pp) != nil {
		t.Fatal("expected nil for nil inner pointer")
	}
}

// Ensure a nil reflect.Type is returned unchanged.
func TestType_nil(t *testing.T) {
	if Type(nil) != nil {
		t.Fatal("expected nil")
	}
}

// Test to ensure a non-pointer type is returned unchanged.
func TestType_nonPointer(t *testing.T) {
	t1 := reflect.TypeOf(0)
	if Type(t1) != t1 {
		t.Fatal("expected int type unchanged")
	}
}

// Test to ensure a single pointer type is dereferenced to its element type.
func TestType_singlePointer(t *testing.T) {
	t1 := reflect.TypeOf((*int)(nil))
	if Type(t1) != reflect.TypeOf(0) {
		t.Fatal("expected int type after deref")
	}
}

// Test to ensure a double pointer type is fully dereferenced to its element type.
func TestType_doublePointer(t *testing.T) {
	v := 0
	p := &v
	t1 := reflect.TypeOf(&p)
	if Type(t1) != reflect.TypeOf(0) {
		t.Fatal("expected int type after double deref")
	}
}

// Test to ensure a plain int value is returned unchanged.
func TestValue_plainInt(t *testing.T) {
	v := reflect.ValueOf(42)
	got := Value(v)
	if got.Int() != 42 {
		t.Fatalf("expected 42, got %v", got)
	}
}

// Test to ensure a single pointer value is dereferenced to its underlying value.
func TestValue_singlePointer(t *testing.T) {
	n := 42
	v := reflect.ValueOf(&n)
	got := Value(v)
	if got.Int() != 42 {
		t.Fatalf("expected 42, got %v", got)
	}
}

// Test to ensure a nil pointer value is returned as-is without panicking.
func TestValue_nilPointer(t *testing.T) {
	var p *int
	v := reflect.ValueOf(p)
	got := Value(v)
	if !got.IsNil() {
		t.Fatal("expected nil value")
	}
}

// Test to ensure a double pointer value is fully dereferenced to its underlying value.
func TestValue_doublePointer(t *testing.T) {
	n := 42
	p := &n
	v := reflect.ValueOf(&p)
	got := Value(v)
	if got.Int() != 42 {
		t.Fatalf("expected 42, got %v", got)
	}
}

// Test to ensure a non-pointer type is returned with changed=false.
func TestTypeKind_nonPointer(t *testing.T) {
	t1 := reflect.TypeOf(0)
	got, k, changed := TypeKind(t1, t1.Kind())
	if changed {
		t.Fatal("expected changed=false")
	}
	if k != reflect.Int {
		t.Fatalf("expected int kind, got %v", k)
	}
	if got != t1 {
		t.Fatal("expected same type")
	}
}

// Test to ensure a single pointer type is dereferenced and changed=true is returned.
func TestTypeKind_singlePointer(t *testing.T) {
	t1 := reflect.TypeOf((*int)(nil))
	got, k, changed := TypeKind(t1, t1.Kind())
	if !changed {
		t.Fatal("expected changed=true")
	}
	if k != reflect.Int {
		t.Fatalf("expected int kind, got %v", k)
	}
	if got != reflect.TypeOf(0) {
		t.Fatal("expected int type after deref")
	}
}

// Test to ensure a double pointer type is fully dereferenced and changed=true is returned.
func TestTypeKind_doublePointer(t *testing.T) {
	n := 0
	p := &n
	t1 := reflect.TypeOf(&p)
	got, k, changed := TypeKind(t1, t1.Kind())
	if !changed {
		t.Fatal("expected changed=true")
	}
	if k != reflect.Int {
		t.Fatalf("expected int kind, got %v", k)
	}
	if got != reflect.TypeOf(0) {
		t.Fatal("expected int type after double deref")
	}
}
