// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package ring

import "testing"

func TestNew_panicsOnZeroChunkSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for chunkSize=0")
		}
	}()
	New[int](0)
}

func TestLen_emptyRing(t *testing.T) {
	r := New[int](4)
	if r.Len() != 0 {
		t.Fatalf("expected 0, got %d", r.Len())
	}
}

func TestCap_emptyRing(t *testing.T) {
	r := New[int](4)
	if r.Cap() != 0 {
		t.Fatalf("expected 0, got %d", r.Cap())
	}
}

func TestEnqueueDequeue_fifoOrder(t *testing.T) {
	r := New[int](4)
	r.Enqueue(1)
	r.Enqueue(2)
	r.Enqueue(3)

	for _, want := range []int{1, 2, 3} {
		v, ok := r.Dequeue()
		if !ok {
			t.Fatal("expected ok=true")
		}
		if v != want {
			t.Fatalf("expected %d, got %d", want, v)
		}
	}
	if r.Len() != 0 {
		t.Fatalf("expected empty ring, got len=%d", r.Len())
	}
}

func TestDequeue_emptyRing(t *testing.T) {
	r := New[int](4)
	_, ok := r.Dequeue()
	if ok {
		t.Fatal("expected ok=false on empty dequeue")
	}
}

func TestGrow_wrappedSlice(t *testing.T) {
	// chunkSize=2 forces growth; interleave enqueue/dequeue to wrap the buffer.
	r := New[int](2)
	r.Enqueue(1)
	r.Enqueue(2)
	r.Dequeue() // back advances to index 1
	r.Enqueue(3)
	r.Enqueue(4) // triggers grow with a wrapped buffer

	for _, want := range []int{2, 3, 4} {
		v, ok := r.Dequeue()
		if !ok {
			t.Fatal("expected ok=true")
		}
		if v != want {
			t.Fatalf("expected %d, got %d", want, v)
		}
	}
}

func TestReset_clearsMemory(t *testing.T) {
	r := New[int](4)
	r.Enqueue(10)
	r.Enqueue(20)
	r.Reset()

	if r.Len() != 0 {
		t.Fatalf("expected len=0 after reset, got %d", r.Len())
	}
	if r.Cap() == 0 {
		t.Fatal("expected non-zero cap after reset")
	}
	for i, v := range r.data {
		if v != 0 {
			t.Fatalf("expected zero at index %d after reset, got %d", i, v)
		}
	}
}

func TestGrow_capIncreasesByChunkSize(t *testing.T) {
	r := New[int](4)
	r.Enqueue(1)
	r.Enqueue(2)
	r.Enqueue(3)
	r.Enqueue(4)
	if r.Cap() != 4 {
		t.Fatalf("expected cap=4, got %d", r.Cap())
	}
	r.Enqueue(5) // triggers grow
	if r.Cap() != 8 {
		t.Fatalf("expected cap=8 after grow, got %d", r.Cap())
	}
}

func TestEnqueueDequeue_pointerType(t *testing.T) {
	r := New[*string](2)
	a, b := "hello", "world"
	r.Enqueue(&a)
	r.Enqueue(&b)

	v, ok := r.Dequeue()
	if !ok || *v != "hello" {
		t.Fatalf("expected hello, got %v", v)
	}
	// Slot should be cleared (nil) after dequeue to avoid memory leaks.
	if r.data[0] != nil {
		t.Fatal("expected nil slot after dequeue of pointer type")
	}
}
