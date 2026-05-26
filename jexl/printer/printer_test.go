// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package printer

import (
	"testing"

	"github.com/harness/go-jexl/jexl/parser"
)

// TestSprint parses a simple arithmetic expression and verifies
// that Sprint reproduces it faithfully.
func TestSprint(t *testing.T) {
	input := "a + b * 2"
	tree, err := parser.Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	output, err := Sprint(tree.Node)
	if err != nil {
		t.Fatal(err)
	}
	if output != input {
		t.Fatalf("got %q, want %q", output, input)
	}
}
