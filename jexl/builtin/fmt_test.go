// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import "testing"

// Ensure Sprintf formats with the given format string and args.
func TestSprintf_basic(t *testing.T) {
	if Sprintf("x=%d y=%s", 1, "hi") != "x=1 y=hi" {
		t.Fatal("expected x=1 y=hi")
	}
}

// Ensure Sprint returns default string representation of args.
func TestSprint_basic(t *testing.T) {
	if Sprint("a", "b") != "ab" {
		t.Fatal("expected ab")
	}
}
