// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import (
	"regexp"

	"github.com/harness/go-jexl/jexl/coerce"
)

// RegexExtract returns the first capture group from text that matches pattern.
// Returns an empty string if there is no match or no capture group.
func RegexExtract(pattern any, text any) string {
	re, err := regexp.Compile(coerce.ToString(pattern))
	if err != nil {
		return ""
	}
	m := re.FindStringSubmatch(coerce.ToString(text))
	if len(m) < 2 {
		return ""
	}
	return m[1]
}
