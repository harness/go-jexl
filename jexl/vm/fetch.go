// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package vm

import (
	"fmt"
	"math/big"
	"reflect"
	"sync"

	"github.com/harness/go-jexl/jexl/classes"
	"github.com/harness/go-jexl/jexl/classes/java/lang"
	"github.com/harness/go-jexl/jexl/classes/java/math"
	"github.com/harness/go-jexl/jexl/classes/java/util"
	"github.com/harness/go-jexl/jexl/coerce"
	"github.com/harness/go-jexl/jexl/internal/decimal"
	"github.com/harness/go-jexl/jexl/internal/deref"
	"github.com/harness/go-jexl/jexl/internal/eval"
)

var fieldCache sync.Map

type fieldCacheKey struct {
	t reflect.Type
	f string
}

// Fetch fetches a value by key or index.
func Fetch(from, i any) any {
	// Antish cursor: dotted-path lookup in flat env.
	var antishCursor *AntishCursor
	switch c := from.(type) {
	case *AntishCursor:
		antishCursor = c
	case AntishCursor:
		antishCursor = &c
	}
	if antishCursor != nil {
		key, isStr := i.(string)
		if !isStr {
			panic(fmt.Sprintf("cannot fetch %v from antish cursor", i))
		}
		dotted := antishCursor.Path + "." + key
		if val, exists := antishCursor.Env[dotted]; exists {
			return val
		}
		// Continue deeper (e.g. x.y.z where x.y.z is the key).
		return &AntishCursor{Env: antishCursor.Env, Path: dotted}
	}

	v := reflect.ValueOf(from)
	if v.Kind() == reflect.Invalid {
		// nil receiver: allow .default(fallback) to work.
		// HACK: Harness compatibility only.
		if methodName, ok := i.(string); ok && methodName == "default" {
			return func(args ...any) any {
				if len(args) > 0 {
					return args[0]
				}
				return nil
			}
		}
		panic(fmt.Sprintf("cannot fetch %v from %T", i, from))
	}

	// Dispatch JEXL collection methods on Go slices and maps.
	if methodName, ok := i.(string); ok {
		if sliceMethod, handled := fetchSliceMethod(from, v, methodName); handled {
			return sliceMethod
		}
		if mapMethod, handled := fetchMapMethod(from, v, methodName); handled {
			return mapMethod
		}
	}

	// classes.Object method/index dispatch.
	if inst, ok := from.(classes.Object); ok {
		if methodName, ok := i.(string); ok {
			return func(args ...any) any {
				result, err := inst.Call(methodName, args...)
				if err != nil {
					panic(err)
				}
				return result
			}
		}
		// Integer index access for java.util.List (e.g. ArrayList).
		if list, ok := from.(util.List); ok {
			idx := coerce.ToInt(i)
			n := list.Size()
			if idx < 0 {
				idx = n + idx
			}
			if idx < 0 || idx >= n {
				panic(fmt.Sprintf("index out of range: %v (list length is %v)", idx, n))
			}
			val, err := list.Get(idx)
			if err != nil {
				panic(err)
			}
			return val
		}
		// Key access for java.util.Map (e.g. HashMap).
		if m, ok := from.(util.Map); ok {
			return m.Get(i)
		}
	}

	// Dispatch method via java.lang (primitive types).
	if methodName, ok := i.(string); ok {
		switch from.(type) {
		case string, bool,
			int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64,
			*big.Int, decimal.Decimal:
			return func(args ...any) any {
				result, _ := valueMethod(from, methodName, args)
				return result
			}
		}
	}

	// Methods can be defined on any type.
	if v.NumMethod() > 0 {
		if methodName, ok := i.(string); ok {
			method := v.MethodByName(methodName)
			if method.IsValid() {
				return method.Interface()
			}
		}
	}

	// Deref pointers before switching on kind.
	v = deref.Value(v)

	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		index := coerce.ToInt(i)
		l := v.Len()
		if index < 0 || index >= l {
			panic(fmt.Sprintf("index out of range: %v (array length is %v)", index, l))
		}
		value := v.Index(index)
		if value.IsValid() {
			return value.Interface()
		}

	case reflect.Map:
		// Fast path for map[string]any.
		if m, ok := from.(map[string]any); ok {
			if i == nil {
				return nil
			}
			key, isStr := i.(string)
			if !isStr {
				return nil
			}
			val, exists := m[key]
			if exists {
				return val
			}
			return nil
		}
		// Fast path for map[any]any.
		if m, ok := from.(map[any]any); ok {
			if i == nil {
				return nil
			}
			if val, exists := m[i]; exists {
				return val
			}
			return nil
		}
		var value reflect.Value
		if i == nil {
			value = v.MapIndex(reflect.Zero(v.Type().Key()))
		} else {
			value = v.MapIndex(reflect.ValueOf(i))
		}
		if value.IsValid() {
			return value.Interface()
		}
		elem := reflect.TypeOf(from).Elem()
		return reflect.Zero(elem).Interface()

	case reflect.Struct:
		fieldName := i.(string)
		t := v.Type()
		key := fieldCacheKey{
			t: t,
			f: fieldName,
		}
		if cv, ok := fieldCache.Load(key); ok {
			return v.FieldByIndex(cv.([]int)).Interface()
		}
		field, ok := t.FieldByNameFunc(func(name string) bool {
			field, _ := t.FieldByName(name)
			switch field.Tag.Get("expr") {
			case "-":
				return false
			case fieldName:
				return true
			default:
				return name == fieldName
			}
		})
		if ok && field.IsExported() {
			value := v.FieldByIndex(field.Index)
			if value.IsValid() {
				fieldCache.Store(key, field.Index)
				return value.Interface()
			}
		}
	}
	panic(fmt.Sprintf("cannot fetch %v from %T", i, from))
}

// FetchField fetches a struct field via a compiled descriptor.
func FetchField(from any, field *Field) any {
	v := reflect.ValueOf(from)
	if v.Kind() != reflect.Invalid {
		v = reflect.Indirect(v)
		value := fieldByIndex(v, field)
		if value.IsValid() {
			return value.Interface()
		}
	}
	panic(fmt.Sprintf("cannot get %v from %T", field.Path[0], from))
}

// FetchMethod fetches a method via a compiled descriptor.
func FetchMethod(from any, method *Method) any {
	v := reflect.ValueOf(from)
	kind := v.Kind()
	if kind != reflect.Invalid {
		// Methods can be defined on any type, no need to dereference.
		method := v.Method(method.Index)
		if method.IsValid() {
			return method.Interface()
		}
	}
	panic(fmt.Sprintf("cannot fetch %v from %T", method.Name, from))
}

// SetIndex sets a value at a key or index in a map or slice.
func SetIndex(obj any, key any, val any) {
	// Antish cursor: assign via dotted key in flat env.
	switch cursor := obj.(type) {
	case *AntishCursor:
		dotted := cursor.Path + "." + fmt.Sprintf("%v", key)
		cursor.Env[dotted] = val
		return
	case AntishCursor:
		dotted := cursor.Path + "." + fmt.Sprintf("%v", key)
		cursor.Env[dotted] = val
		return
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Map:
		v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
	case reflect.Slice, reflect.Array:
		idx := coerce.ToInt(key)
		if idx < 0 {
			idx = v.Len() + idx
		}
		v.Index(idx).Set(reflect.ValueOf(val))
	default:
		panic(fmt.Sprintf("cannot set index on %T", obj))
	}
}

// fetchMapMethod returns a JEXL method for Go maps.
func fetchMapMethod(from any, v reflect.Value, methodName string) (any, bool) {
	kind := v.Kind()
	if kind == reflect.Ptr || kind == reflect.Interface {
		v2 := deref.Value(v)
		kind = v2.Kind()
	}
	if kind != reflect.Map {
		return nil, false
	}
	switch methodName {
	case "size":
		l := v.Len()
		return func(args ...any) any { return l }, true
	case "isEmpty":
		empty := v.Len() == 0
		return func(args ...any) any { return empty }, true
	case "containsKey":
		return func(args ...any) any {
			if len(args) == 0 {
				return false
			}
			val := v.MapIndex(reflect.ValueOf(args[0]))
			return val.IsValid()
		}, true
	case "containsValue":
		return func(args ...any) any {
			if len(args) == 0 {
				return false
			}
			for _, k := range v.MapKeys() {
				if eval.Equal(v.MapIndex(k).Interface(), args[0]) {
					return true
				}
			}
			return false
		}, true
	case "get":
		return func(args ...any) any {
			if len(args) == 0 {
				return nil
			}
			val := v.MapIndex(reflect.ValueOf(args[0]))
			if !val.IsValid() {
				return nil
			}
			return val.Interface()
		}, true
	case "put":
		return func(args ...any) any {
			if len(args) < 2 {
				return nil
			}
			key := reflect.ValueOf(args[0])
			val := reflect.ValueOf(args[1])
			v.SetMapIndex(key, val)
			return nil
		}, true
	case "remove":
		return func(args ...any) any {
			if len(args) == 0 {
				return nil
			}
			key := reflect.ValueOf(args[0])
			old := v.MapIndex(key)
			v.SetMapIndex(key, reflect.Value{})
			if old.IsValid() {
				return old.Interface()
			}
			return nil
		}, true
	case "keySet":
		return func(args ...any) any {
			keys := make([]any, 0, v.Len())
			for _, k := range v.MapKeys() {
				keys = append(keys, k.Interface())
			}
			return keys
		}, true
	case "values":
		return func(args ...any) any {
			vals := make([]any, 0, v.Len())
			for _, k := range v.MapKeys() {
				vals = append(vals, v.MapIndex(k).Interface())
			}
			return vals
		}, true
	case "entrySet":
		return func(args ...any) any {
			entries := make([]any, 0, v.Len())
			for _, k := range v.MapKeys() {
				entry := map[string]any{"key": k.Interface(), "value": v.MapIndex(k).Interface()}
				entries = append(entries, entry)
			}
			return entries
		}, true
	case "clear":
		return func(args ...any) any {
			for _, k := range v.MapKeys() {
				v.SetMapIndex(k, reflect.Value{})
			}
			return nil
		}, true
	case "putIfAbsent":
		return func(args ...any) any {
			if len(args) < 2 {
				return nil
			}
			key := reflect.ValueOf(args[0])
			if !v.MapIndex(key).IsValid() {
				v.SetMapIndex(key, reflect.ValueOf(args[1]))
			}
			return v.MapIndex(key).Interface()
		}, true
	case "getOrDefault":
		return func(args ...any) any {
			if len(args) < 2 {
				return nil
			}
			val := v.MapIndex(reflect.ValueOf(args[0]))
			if val.IsValid() {
				return val.Interface()
			}
			return args[1]
		}, true
	case "toString":
		return func(args ...any) any { return fmt.Sprintf("%v", from) }, true
	}
	return nil, false
}

// fetchSliceMethod returns a JEXL method for Go slices.
func fetchSliceMethod(from any, v reflect.Value, methodName string) (any, bool) {
	kind := v.Kind()
	if kind == reflect.Ptr || kind == reflect.Interface {
		v2 := deref.Value(v)
		kind = v2.Kind()
	}
	if kind != reflect.Slice && kind != reflect.Array {
		return nil, false
	}
	switch methodName {
	case "size", "length":
		l := v.Len()
		return func(args ...any) any { return l }, true
	case "isEmpty":
		empty := v.Len() == 0
		return func(args ...any) any { return empty }, true
	case "contains":
		return func(args ...any) any {
			if len(args) == 0 {
				return false
			}
			for i := 0; i < v.Len(); i++ {
				if eval.Equal(v.Index(i).Interface(), args[0]) {
					return true
				}
			}
			return false
		}, true
	case "get":
		return func(args ...any) any {
			if len(args) == 0 {
				panic("get requires an index argument")
			}
			idx := coerce.ToInt(args[0])
			if idx < 0 {
				idx = v.Len() + idx
			}
			if idx < 0 || idx >= v.Len() {
				panic(fmt.Sprintf("index %d out of bounds (size=%d)", idx, v.Len()))
			}
			return v.Index(idx).Interface()
		}, true
	case "set":
		return func(args ...any) any {
			if len(args) < 2 {
				return nil
			}
			idx := coerce.ToInt(args[0])
			if idx < 0 {
				idx = v.Len() + idx
			}
			if idx < 0 || idx >= v.Len() {
				panic(fmt.Sprintf("index %d out of bounds (size=%d)", idx, v.Len()))
			}
			elem := v.Index(idx)
			if elem.CanSet() {
				elem.Set(reflect.ValueOf(args[1]))
			}
			return args[1]
		}, true
	case "add":
		return func(args ...any) any { return true }, true
	case "remove":
		return func(args ...any) any { return nil }, true
	case "indexOf":
		return func(args ...any) any {
			if len(args) == 0 {
				return -1
			}
			for i := 0; i < v.Len(); i++ {
				if eval.Equal(v.Index(i).Interface(), args[0]) {
					return i
				}
			}
			return -1
		}, true
	case "toArray":
		items := make([]any, v.Len())
		for i := 0; i < v.Len(); i++ {
			items[i] = v.Index(i).Interface()
		}
		return func(args ...any) any { return items }, true
	case "subList":
		return func(args ...any) any {
			if len(args) < 2 {
				return from
			}
			from := coerce.ToInt(args[0])
			to := coerce.ToInt(args[1])
			if from < 0 {
				from = 0
			}
			if to > v.Len() {
				to = v.Len()
			}
			items := make([]any, to-from)
			for i := from; i < to; i++ {
				items[i-from] = v.Index(i).Interface()
			}
			return items
		}, true
	case "sort":
		return func(args ...any) any { return from }, true
	case "reverse":
		return func(args ...any) any { return from }, true
	case "clear":
		return func(args ...any) any { return nil }, true
	case "toString":
		return func(args ...any) any { return fmt.Sprintf("%v", from) }, true
	case "stream":
		return func(args ...any) any { return from }, true
	case "iterator":
		return func(args ...any) any { return from }, true
	}
	return nil, false
}

// fieldByIndex traverses nested struct fields by index path.
func fieldByIndex(v reflect.Value, field *Field) reflect.Value {
	if len(field.Index) == 1 {
		return v.Field(field.Index[0])
	}
	for i, x := range field.Index {
		if i > 0 {
			if v.Kind() == reflect.Ptr {
				if v.IsNil() {
					panic(fmt.Sprintf("cannot get %v from %v", field.Path[i], field.Path[i-1]))
				}
				v = v.Elem()
			}
		}
		v = v.Field(x)
	}
	return v
}

// valueMethod dispatches a method to the java.lang class.
func valueMethod(v any, methodName string, args []any) (any, bool) {
	switch bi := v.(type) {
	case *big.Int:
		result, err := (&math.BigInteger{V: bi}).Call(methodName, args)
		if err != nil {
			return nil, false
		}
		if r, ok := result.(*math.BigInteger); ok {
			return r.V, true
		}
		return result, true
	case decimal.Decimal:
		result, err := (&math.BigDecimal{V: bi}).Call(methodName, args)
		if err != nil {
			return nil, false
		}
		if r, ok := result.(*math.BigDecimal); ok {
			return r.V, true
		}
		return result, true
	}
	var obj classes.Object
	switch v.(type) {
	case string:
		obj = lang.NewStringFrom(v)
	case bool:
		obj = lang.NewBooleanFrom(v)
	case int8:
		obj = lang.NewByteFrom(v)
	case int16:
		obj = lang.NewShortFrom(v)
	case int32: // covers rune — treated as Character
		obj = lang.NewCharacterFrom(v)
	case int, int64, uint, uint8, uint16, uint32, uint64:
		obj = lang.NewLongFrom(v)
	case float32:
		obj = lang.NewFloatFrom(v)
	case float64:
		obj = lang.NewDoubleFrom(v)
	}
	if obj != nil {
		result, err := obj.Call(methodName, args...)
		if err != nil {
			return nil, false
		}
		return result, true
	}
	return nil, false
}
