// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package token

// Range is a half-open interval [From, To) of rune offsets in a source file.
type Range struct {
	From int `json:"from"`
	To   int `json:"to"`
}
