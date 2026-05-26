// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package ring

// Ring is a very simple ring buffer implementation that
// uses a slice. The internal slice will only grow, never
// shrink. Pointer and reference types can be safely used
// because memory is cleared.
type Ring[T any] struct {
	data                 []T
	back, len, chunkSize int
}

// New returns a new Ring with the given chunk size.
// Panics if chunkSize < 1.
func New[T any](chunkSize int) *Ring[T] {
	if chunkSize < 1 {
		panic("chunkSize must be greater than zero")
	}
	return &Ring[T]{
		chunkSize: chunkSize,
	}
}

// Len returns the number of items currently in the ring.
func (r *Ring[T]) Len() int {
	return r.len
}

// Cap returns the current capacity of the underlying slice.
func (r *Ring[T]) Cap() int {
	return len(r.data)
}

// Reset clears all items and zeroes memory.
// Capacity is preserved.
func (r *Ring[T]) Reset() {
	var zero T
	for i := range r.data {
		// TODO the 'clear' builtin can be used
		r.data[i] = zero // clear memory, optimized by the compiler
	}
	r.back = 0
	r.len = 0
}

// Dequeue returns the oldest value.
func (r *Ring[T]) Dequeue() (v T, ok bool) {
	if r.len == 0 {
		return v, false
	}
	v, r.data[r.back] = r.data[r.back], v // retrieve and clear mem
	r.len--
	r.back = (r.back + 1) % len(r.data)
	return v, true
}

// Enqueue adds an item to the ring.
func (r *Ring[T]) Enqueue(v T) {
	if r.len == len(r.data) {
		r.grow()
	}
	writePos := (r.back + r.len) % len(r.data)
	r.data[writePos] = v
	r.len++
}

// grow expands the underlying slice by chunkSize,
// copying existing items and unwrapping any wrapped layout.
func (r *Ring[T]) grow() {
	s := make([]T, len(r.data)+r.chunkSize)
	if r.len > 0 {
		chunk1 := r.back + r.len
		if chunk1 > len(r.data) {
			chunk1 = len(r.data)
		}
		copied := copy(s, r.data[r.back:chunk1])

		if copied < r.len { // wrapped slice
			chunk2 := r.len - copied
			copy(s[copied:], r.data[:chunk2])
		}
	}
	r.back = 0
	r.data = s
}
