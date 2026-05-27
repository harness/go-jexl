// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package jexl

import (
	"strings"

	"github.com/harness/go-jexl/jexl/vm"
)

// isIdentStart reports whether c is a valid identifier
// start byte.
func isIdentStart(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

// isIdentPart reports whether c is a valid identifier
// continuation byte.
func isIdentPart(c byte) bool {
	return isIdentStart(c) || (c >= '0' && c <= '9')
}

// isPropertyPath reports whether s is a pure dot-separated
// property path.
func isPropertyPath(s string) bool {
	if len(s) == 0 {
		return false
	}
	if !isIdentStart(s[0]) {
		return false
	}
	for i := 1; i < len(s); i++ {
		c := s[i]
		if c == '.' {
			if i+1 >= len(s) || !isIdentStart(s[i+1]) {
				return false
			}
		} else if !isIdentPart(c) {
			return false
		}
	}
	return true
}

// fetchKey returns the value for key in from without
// reflection. Handles map[string]any and map[string]string only.
func fetchKey(from any, key string) (any, bool) {
	switch m := from.(type) {
	case map[string]any:
		v, ok := m[key]
		return v, ok
	case map[string]string:
		v, ok := m[key]
		return v, ok
	}
	return nil, false
}

// evalPath walks path one segment at a time against env.
// Falls back to vm.Fetch for structs and other types.
func evalPath(env any, path string) (any, error) {
	cur := env
	for {
		seg, rest, hasDot := strings.Cut(path, ".")
		if v, ok := fetchKey(cur, seg); ok {
			if !hasDot {
				return v, nil
			}
			if v == nil {
				return nil, nil
			}
			cur = v
			path = rest
			continue
		}
		// fallback to vm.Fetch for structs and other types
		fetched := func() (val any, panicked bool) {
			defer func() {
				if r := recover(); r != nil {
					panicked = true
				}
			}()
			val = vm.Fetch(cur, seg)
			return val, false
		}
		val, panicked := fetched()
		if panicked {
			return nil, nil
		}
		if !hasDot {
			return val, nil
		}
		if val == nil {
			return nil, nil
		}
		cur = val
		path = rest
	}
}
