// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package util

import (
	"fmt"
	"strings"

	"github.com/harness/go-jexl/jexl/coerce"
)

// ArrayListClass is the java.util.ArrayList class object.
var ArrayListClass arrayListClass

// LinkedListClass aliases ArrayListClass.
var LinkedListClass = ArrayListClass

type arrayListClass struct{}

func (arrayListClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "new":
		l := make(ArrayList, 0)
		return &l, nil
	}
	return nil, fmt.Errorf("ArrayList.%s: undefined", method)
}

// ArrayList mirrors java.util.ArrayList.
type ArrayList []any

// Size returns the number of items.
func (l ArrayList) Size() int { return len(l) }

// IsEmpty returns true if the list has no items.
func (l ArrayList) IsEmpty() bool { return len(l) == 0 }

// Get returns the item at the given index.
func (l ArrayList) Get(idx int) (any, error) {
	if idx < 0 || idx >= len(l) {
		return nil, fmt.Errorf("ArrayList.get: index %d out of bounds", idx)
	}
	return l[idx], nil
}

// Contains returns true if val is in the list.
func (l ArrayList) Contains(val any) bool {
	s := fmt.Sprintf("%v", val)
	for _, item := range l {
		if fmt.Sprintf("%v", item) == s {
			return true
		}
	}
	return false
}

// IndexOf returns the first index of val, or -1 if not found.
func (l ArrayList) IndexOf(val any) int {
	s := fmt.Sprintf("%v", val)
	for i, item := range l {
		if fmt.Sprintf("%v", item) == s {
			return i
		}
	}
	return -1
}

// ToArray returns a copy of the underlying slice.
func (l ArrayList) ToArray() []any {
	out := make([]any, len(l))
	copy(out, l)
	return out
}

// ToString returns a string like "[a, b, c]".
func (l ArrayList) ToString() string {
	parts := make([]string, len(l))
	for i, item := range l {
		parts[i] = fmt.Sprintf("%v", item)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// Add appends a value (1 arg) or inserts at index (2 args: idx, val).
func (l *ArrayList) Add(args ...any) (any, error) {
	switch len(args) {
	case 1:
		*l = append(*l, args[0])
		return true, nil
	case 2:
		idx := coerce.ToInt(args[0])
		val := args[1]
		if idx < 0 || idx > len(*l) {
			return nil, fmt.Errorf("ArrayList.add: index %d out of bounds", idx)
		}
		*l = append(*l, nil)
		copy((*l)[idx+1:], (*l)[idx:])
		(*l)[idx] = val
		return nil, nil
	default:
		return nil, fmt.Errorf("ArrayList.add: expected 1 or 2 arguments")
	}
}

// Set replaces the item at idx with val, returning the old value.
func (l *ArrayList) Set(idx int, val any) (any, error) {
	if idx < 0 || idx >= len(*l) {
		return nil, fmt.Errorf("ArrayList.set: index %d out of bounds", idx)
	}
	old := (*l)[idx]
	(*l)[idx] = val
	return old, nil
}

// Remove removes and returns the item at idx.
func (l *ArrayList) Remove(idx int) (any, error) {
	if idx < 0 || idx >= len(*l) {
		return nil, fmt.Errorf("ArrayList.remove: index %d out of bounds", idx)
	}
	old := (*l)[idx]
	*l = append((*l)[:idx], (*l)[idx+1:]...)
	return old, nil
}

// Clear removes all items.
func (l *ArrayList) Clear() { *l = (*l)[:0] }

// AddAll appends all items from other, returns true if list changed.
func (l *ArrayList) AddAll(other ArrayList) bool {
	if len(other) == 0 {
		return false
	}
	*l = append(*l, other...)
	return true
}

// SubList returns a new ArrayList with items from from (inclusive) to to (exclusive).
func (l ArrayList) SubList(from, to int) *ArrayList {
	sub := make(ArrayList, to-from)
	copy(sub, l[from:to])
	return &sub
}

// AddFirst inserts val at the beginning.
func (l *ArrayList) AddFirst(val any) {
	*l = append(ArrayList{val}, *l...)
}

// AddLast appends val at the end.
func (l *ArrayList) AddLast(val any) {
	*l = append(*l, val)
}

// GetFirst returns the first item or an error if empty.
func (l ArrayList) GetFirst() (any, error) {
	if len(l) == 0 {
		return nil, fmt.Errorf("ArrayList.getFirst: list is empty")
	}
	return l[0], nil
}

// GetLast returns the last item or an error if empty.
func (l ArrayList) GetLast() (any, error) {
	if len(l) == 0 {
		return nil, fmt.Errorf("ArrayList.getLast: list is empty")
	}
	return l[len(l)-1], nil
}

// Peek returns the first item or nil if empty.
func (l ArrayList) Peek() any {
	if len(l) == 0 {
		return nil
	}
	return l[0]
}

// PeekLast returns the last item or nil if empty.
func (l ArrayList) PeekLast() any {
	if len(l) == 0 {
		return nil
	}
	return l[len(l)-1]
}

// Poll removes and returns the first item, or nil if empty.
func (l *ArrayList) Poll() any {
	if len(*l) == 0 {
		return nil
	}
	val := (*l)[0]
	*l = (*l)[1:]
	return val
}

// PollLast removes and returns the last item, or nil if empty.
func (l *ArrayList) PollLast() any {
	if len(*l) == 0 {
		return nil
	}
	last := len(*l) - 1
	val := (*l)[last]
	*l = (*l)[:last]
	return val
}

// Call dispatches instance methods.
func (l *ArrayList) Call(method string, args ...any) (any, error) {
	switch method {
	case "add":
		return l.Add(args...)
	case "get":
		if len(args) != 1 {
			return nil, fmt.Errorf("ArrayList.get: expected 1 argument")
		}
		return l.Get(coerce.ToInt(args[0]))
	case "set":
		if len(args) != 2 {
			return nil, fmt.Errorf("ArrayList.set: expected 2 arguments")
		}
		return l.Set(coerce.ToInt(args[0]), args[1])
	case "remove":
		if len(args) != 1 {
			return nil, fmt.Errorf("ArrayList.remove: expected 1 argument")
		}
		return l.Remove(coerce.ToInt(args[0]))
	case "size", "length":
		return l.Size(), nil
	case "isEmpty":
		return l.IsEmpty(), nil
	case "contains":
		if len(args) != 1 {
			return nil, fmt.Errorf("ArrayList.contains: expected 1 argument")
		}
		return l.Contains(args[0]), nil
	case "indexOf":
		if len(args) != 1 {
			return nil, fmt.Errorf("ArrayList.indexOf: expected 1 argument")
		}
		return l.IndexOf(args[0]), nil
	case "clear":
		l.Clear()
		return nil, nil
	case "addAll":
		if len(args) != 1 {
			return nil, fmt.Errorf("ArrayList.addAll: expected 1 argument")
		}
		other, ok := args[0].(*ArrayList)
		if !ok {
			return nil, fmt.Errorf("ArrayList.addAll: expected *ArrayList argument")
		}
		return l.AddAll(*other), nil
	case "subList":
		if len(args) != 2 {
			return nil, fmt.Errorf("ArrayList.subList: expected 2 arguments")
		}
		return l.SubList(coerce.ToInt(args[0]), coerce.ToInt(args[1])), nil
	case "toArray":
		return l.ToArray(), nil
	case "toString":
		return l.ToString(), nil
	case "addFirst":
		if len(args) != 1 {
			return nil, fmt.Errorf("ArrayList.addFirst: expected 1 argument")
		}
		l.AddFirst(args[0])
		return nil, nil
	case "addLast":
		if len(args) != 1 {
			return nil, fmt.Errorf("ArrayList.addLast: expected 1 argument")
		}
		l.AddLast(args[0])
		return nil, nil
	case "getFirst":
		return l.GetFirst()
	case "getLast":
		return l.GetLast()
	case "iterator":
		return l, nil
	case "hasNext":
		return len(*l) > 0, nil
	case "next":
		return l.Poll(), nil
	case "peek":
		return l.Peek(), nil
	case "peekLast":
		return l.PeekLast(), nil
	case "poll":
		return l.Poll(), nil
	case "pollLast":
		return l.PollLast(), nil
	}
	return nil, fmt.Errorf("ArrayList instance: undefined method %q", method)
}
