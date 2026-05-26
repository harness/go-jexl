// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package util

import (
	"fmt"
	"strings"
)

// HashMapClass is the java.util.HashMap class object.
var HashMapClass hashMapClass

// TreeMapClass aliases HashMapClass.
var TreeMapClass = HashMapClass

type hashMapClass struct{}

func (hashMapClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "new":
		return make(HashMap), nil
	}
	return nil, fmt.Errorf("HashMap.%s: undefined", method)
}

// HashMap mirrors java.util.HashMap.
type HashMap map[any]any

// FetchKey returns the value for key (for VM key access).
func (m HashMap) FetchKey(key any) any { return m[key] }

// Put stores val under key and returns the old value (nil if absent).
func (m HashMap) Put(key, val any) any {
	old := m[key]
	m[key] = val
	return old
}

// Get returns the value for key, or nil if not present.
func (m HashMap) Get(key any) any {
	return m[key]
}

// Remove deletes key and returns its old value.
func (m HashMap) Remove(key any) any {
	old := m[key]
	delete(m, key)
	return old
}

// ContainsKey returns true if the key exists.
func (m HashMap) ContainsKey(key any) bool {
	_, ok := m[key]
	return ok
}

// ContainsValue returns true if val is present (fmt.Sprintf comparison).
func (m HashMap) ContainsValue(val any) bool {
	s := fmt.Sprintf("%v", val)
	for _, v := range m {
		if fmt.Sprintf("%v", v) == s {
			return true
		}
	}
	return false
}

// Size returns the number of entries.
func (m HashMap) Size() int { return len(m) }

// IsEmpty returns true if the map has no entries.
func (m HashMap) IsEmpty() bool { return len(m) == 0 }

// Clear removes all entries.
func (m HashMap) Clear() {
	for k := range m {
		delete(m, k)
	}
}

// KeySet returns a slice of all keys.
func (m HashMap) KeySet() []any {
	keys := make([]any, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of all values.
func (m HashMap) Values() []any {
	vals := make([]any, 0, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals
}

// EntrySet returns a slice of map[any]any{"key": k, "value": v} entries.
func (m HashMap) EntrySet() []any {
	entries := make([]any, 0, len(m))
	for k, v := range m {
		entries = append(entries, map[any]any{"key": k, "value": v})
	}
	return entries
}

// GetOrDefault returns the value for key, or def if key is absent.
func (m HashMap) GetOrDefault(key, def any) any {
	if v, ok := m[key]; ok {
		return v
	}
	return def
}

// PutIfAbsent stores val only if key is absent; returns the existing or new value.
func (m HashMap) PutIfAbsent(key, val any) any {
	if v, ok := m[key]; ok {
		return v
	}
	m[key] = val
	return val
}

// PutAll copies all entries from other into m.
func (m HashMap) PutAll(other HashMap) {
	for k, v := range other {
		m[k] = v
	}
}

// ToString returns a string like "{k1=v1, k2=v2}".
func (m HashMap) ToString() string {
	parts := make([]string, 0, len(m))
	for k, v := range m {
		parts = append(parts, fmt.Sprintf("%v=%v", k, v))
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

// Call dispatches instance methods.
func (m HashMap) Call(method string, args ...any) (any, error) {
	switch method {
	case "put":
		if len(args) != 2 {
			return nil, fmt.Errorf("HashMap.put: expected 2 arguments")
		}
		return m.Put(args[0], args[1]), nil
	case "get":
		if len(args) != 1 {
			return nil, fmt.Errorf("HashMap.get: expected 1 argument")
		}
		return m.Get(args[0]), nil
	case "remove":
		if len(args) != 1 {
			return nil, fmt.Errorf("HashMap.remove: expected 1 argument")
		}
		return m.Remove(args[0]), nil
	case "containsKey":
		if len(args) != 1 {
			return nil, fmt.Errorf("HashMap.containsKey: expected 1 argument")
		}
		return m.ContainsKey(args[0]), nil
	case "containsValue":
		if len(args) != 1 {
			return nil, fmt.Errorf("HashMap.containsValue: expected 1 argument")
		}
		return m.ContainsValue(args[0]), nil
	case "size":
		return m.Size(), nil
	case "isEmpty":
		return m.IsEmpty(), nil
	case "clear":
		m.Clear()
		return nil, nil
	case "keySet":
		return m.KeySet(), nil
	case "values":
		return m.Values(), nil
	case "entrySet":
		return m.EntrySet(), nil
	case "getOrDefault":
		if len(args) != 2 {
			return nil, fmt.Errorf("HashMap.getOrDefault: expected 2 arguments")
		}
		return m.GetOrDefault(args[0], args[1]), nil
	case "putIfAbsent":
		if len(args) != 2 {
			return nil, fmt.Errorf("HashMap.putIfAbsent: expected 2 arguments")
		}
		return m.PutIfAbsent(args[0], args[1]), nil
	case "putAll":
		if len(args) != 1 {
			return nil, fmt.Errorf("HashMap.putAll: expected 1 argument")
		}
		other, ok := args[0].(HashMap)
		if !ok {
			return nil, fmt.Errorf("HashMap.putAll: expected HashMap argument")
		}
		m.PutAll(other)
		return nil, nil
	case "toString":
		return m.ToString(), nil
	}
	return nil, fmt.Errorf("HashMap instance: undefined method %q", method)
}
