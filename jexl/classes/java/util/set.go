// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package util

// Set is implemented by types that behave like java.util.Set.
type Set interface {
	Contains(v any) bool
	Size() int
	ToArray() []any
}
