// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package util

// List is implemented by types that behave like java.util.List.
type List interface {
	Size() int
	Get(i int) (any, error)
}
