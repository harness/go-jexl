// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import "testing"

// Ensure HexEncode returns the hex string for ASCII input.
func TestHexEncode_basic(t *testing.T) {
	if HexEncode("hi") != "6869" {
		t.Fatal("expected 6869")
	}
}

// Ensure HexDecode recovers the original string.
func TestHexDecode_basic(t *testing.T) {
	b, err := HexDecode("6869")
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "hi" {
		t.Fatalf("expected hi, got %s", b)
	}
}

// Ensure HexDecode returns an error for invalid hex input.
func TestHexDecode_invalid(t *testing.T) {
	_, err := HexDecode("zzzz")
	if err == nil {
		t.Fatal("expected error")
	}
}
