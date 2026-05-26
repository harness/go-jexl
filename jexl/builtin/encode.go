// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import (
	"encoding/base64"

	"github.com/harness/go-jexl/jexl/coerce"
)

// Base64Encode returns the base64 standard encoding of v.
func Base64Encode(v any) string {
	return base64.StdEncoding.EncodeToString(
		[]byte(coerce.ToString(v)),
	)
}

// Base64Decode returns the bytes decoded from base64 standard encoding v.
func Base64Decode(v any) ([]byte, error) {
	return base64.StdEncoding.DecodeString(
		coerce.ToString(v),
	)
}

// Base64URLEncode returns the base64 URL encoding of v.
func Base64URLEncode(v any) string {
	return base64.URLEncoding.EncodeToString(
		[]byte(coerce.ToString(v)),
	)
}

// Base64URLDecode returns the bytes decoded from base64 URL encoding v.
func Base64URLDecode(v any) ([]byte, error) {
	return base64.URLEncoding.DecodeString(
		coerce.ToString(v),
	)
}

// Base64RawEncode returns the unpadded base64 standard encoding of v.
func Base64RawEncode(v any) string {
	return base64.RawStdEncoding.EncodeToString(
		[]byte(coerce.ToString(v)),
	)
}

// Base64RawDecode returns the bytes decoded from unpadded base64 standard encoding v.
func Base64RawDecode(v any) ([]byte, error) {
	return base64.RawStdEncoding.DecodeString(
		coerce.ToString(v),
	)
}
