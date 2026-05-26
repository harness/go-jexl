// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import "fmt"

// Sprintf returns a formatted string.
func Sprintf(format any, args ...any) string {
	return fmt.Sprintf(
		fmt.Sprint(format), args...,
	)
}

// Sprint returns the default string representation of args.
func Sprint(args ...any) string {
	return fmt.Sprint(args...)
}

