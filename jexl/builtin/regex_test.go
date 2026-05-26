// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import "testing"

// Ensure RegexExtract returns the first capture group.
func TestRegexExtract_match(t *testing.T) {
	got := RegexExtract(`v(\d+\.\d+)`, "release v1.42 build")
	if got != "1.42" {
		t.Fatalf("expected 1.42, got %q", got)
	}
}

// Ensure RegexExtract returns empty string when there is no match.
func TestRegexExtract_noMatch(t *testing.T) {
	got := RegexExtract(`v(\d+)`, "no version here")
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

// Ensure RegexExtract returns empty string for an invalid pattern.
func TestRegexExtract_invalidPattern(t *testing.T) {
	got := RegexExtract(`[invalid`, "text")
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}
