// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import (
	"encoding/hex"

	"github.com/harness/go-jexl/jexl/coerce"
)

// HexEncode returns the hexadecimal encoding of v.
func HexEncode(v any) string {
	return hex.EncodeToString(
		[]byte(coerce.ToString(v)),
	)
}

// HexDecode returns the bytes decoded from the hexadecimal string v.
func HexDecode(v any) ([]byte, error) {
	return hex.DecodeString(
		coerce.ToString(v),
	)
}
