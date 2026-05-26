// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import "testing"

// Ensure Base64Encode produces the standard base64 encoding.
func TestBase64Encode_basic(t *testing.T) {
	if Base64Encode("hello") != "aGVsbG8=" {
		t.Fatal("expected aGVsbG8=")
	}
}

// Ensure Base64Decode recovers the original string.
func TestBase64Decode_basic(t *testing.T) {
	b, err := Base64Decode("aGVsbG8=")
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "hello" {
		t.Fatalf("expected hello, got %s", b)
	}
}

// Ensure Base64Decode returns an error for invalid input.
func TestBase64Decode_invalid(t *testing.T) {
	_, err := Base64Decode("!!!not-base64!!!")
	if err == nil {
		t.Fatal("expected error")
	}
}

// Ensure Base64URLEncode produces URL-safe base64 encoding.
func TestBase64URLEncode_basic(t *testing.T) {
	if Base64URLEncode("hello") != "aGVsbG8=" {
		t.Fatal("expected aGVsbG8=")
	}
}

// Ensure Base64URLDecode recovers the original string.
func TestBase64URLDecode_basic(t *testing.T) {
	b, err := Base64URLDecode("aGVsbG8=")
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "hello" {
		t.Fatalf("expected hello, got %s", b)
	}
}

// Ensure Base64RawEncode produces unpadded base64 encoding.
func TestBase64RawEncode_basic(t *testing.T) {
	if Base64RawEncode("hello") != "aGVsbG8" {
		t.Fatal("expected aGVsbG8")
	}
}

// Ensure Base64RawDecode recovers the original string.
func TestBase64RawDecode_basic(t *testing.T) {
	b, err := Base64RawDecode("aGVsbG8")
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "hello" {
		t.Fatalf("expected hello, got %s", b)
	}
}
