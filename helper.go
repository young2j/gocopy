/*
 * File: helper.go
 * Created Date: 2022-01-26 05:48:28
 * Author: ysj
 * Description: helper func
 */
package gocopy

import (
	"reflect"
	"strings"
	"sync"
)

var (
	// ponytail: process-wide caches; add eviction only if dynamic plugins create unbounded types.
	deepFieldsCache sync.Map // map[reflect.Type][]reflect.StructField
	fieldCache      sync.Map // map[fieldCacheKey]fieldCacheValue
	copyStatePool   = sync.Pool{New: func() interface{} { return &copyState{} }}
)

type fieldCacheKey struct {
	typ  reflect.Type
	name string
}

type fieldCacheValue struct {
	field reflect.StructField
	ok    bool
}

type copyState struct {
	entries  [8]visitEntry
	count    int
	overflow map[visitKey]reflect.Value
}

type visitEntry struct {
	key   visitKey
	value reflect.Value
}

func (state *copyState) reset() {
	for i := 0; i < state.count; i++ {
		state.entries[i] = visitEntry{}
	}
	state.count = 0
	state.overflow = nil
}

type visitKey struct {
	source      uintptr
	sourceKind  reflect.Kind
	destination reflect.Type
}

func beginPointerCopy(fromValue, toValue reflect.Value, opt *Option) bool {
	if opt == nil || opt.state == nil {
		return false
	}
	source, sourceKind, ok := pointerAddress(fromValue)
	if !ok {
		return false
	}
	key := visitKey{source: source, sourceKind: sourceKind, destination: toValue.Type()}
	if previous, exists := opt.state.lookup(key); exists {
		if toValue.Kind() == reflect.Ptr && toValue.CanSet() && previous.Type().AssignableTo(toValue.Type()) {
			toValue.Set(previous)
		}
		return true
	}
	if toValue.Kind() == reflect.Ptr && !toValue.IsNil() {
		opt.state.store(key, toValue)
	}
	return false
}

func (state *copyState) lookup(key visitKey) (reflect.Value, bool) {
	for i := 0; i < state.count; i++ {
		if state.entries[i].key == key {
			return state.entries[i].value, true
		}
	}
	if state.overflow == nil {
		return reflect.Value{}, false
	}
	value, ok := state.overflow[key]
	return value, ok
}

func (state *copyState) store(key visitKey, value reflect.Value) {
	if state.count < len(state.entries) {
		state.entries[state.count] = visitEntry{key: key, value: value}
		state.count++
		return
	}
	if state.overflow == nil {
		state.overflow = make(map[visitKey]reflect.Value)
	}
	state.overflow[key] = value
}

func pointerAddress(value reflect.Value) (uintptr, reflect.Kind, bool) {
	for value.IsValid() && value.Kind() == reflect.Interface {
		if value.IsNil() {
			return 0, reflect.Invalid, false
		}
		value = value.Elem()
	}
	if !value.IsValid() || (value.Kind() != reflect.Ptr && value.Kind() != reflect.Map && value.Kind() != reflect.Slice) || value.IsNil() {
		return 0, reflect.Invalid, false
	}
	return value.Pointer(), value.Kind(), true
}

func fieldTag(field reflect.StructField) (name string, omitEmpty, skip bool) {
	for _, key := range []string{"gocopy", "json", "bson"} {
		raw, ok := field.Tag.Lookup(key)
		if !ok {
			continue
		}
		parts := strings.Split(raw, ",")
		if parts[0] == "-" {
			return "", false, true
		}
		for _, option := range parts[1:] {
			if option == "omitempty" {
				omitEmpty = true
			}
		}
		if parts[0] != "" {
			name = parts[0]
		}
		if name != "" {
			return name, omitEmpty, false
		}
	}
	return "", false, false
}

func fieldByName(reflectType reflect.Type, name string) (reflect.StructField, bool) {
	key := fieldCacheKey{typ: reflectType, name: name}
	if cached, ok := fieldCache.Load(key); ok {
		value := cached.(fieldCacheValue)
		return value.field, value.ok
	}
	field, ok := reflectType.FieldByName(name)
	fieldCache.Store(key, fieldCacheValue{field: field, ok: ok})
	return field, ok
}

func indirectValue(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func indirectType(reflectType reflect.Type) (reflect.Type, bool) {
	isPtr := false
	for reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
		isPtr = true
	}
	return reflectType, isPtr
}

func deepFields(reflectType reflect.Type) []reflect.StructField {
	if reflectType, _ = indirectType(reflectType); reflectType.Kind() != reflect.Struct {
		return nil
	}
	if cached, ok := deepFieldsCache.Load(reflectType); ok {
		return cached.([]reflect.StructField)
	}

	fields := make([]reflect.StructField, 0, reflectType.NumField())
	for i := 0; i < reflectType.NumField(); i++ {
		v := reflectType.Field(i)
		// PkgPath is the package path that qualifies a lower case (unexported)
		// field name. It is empty for upper case (exported) field names.
		if v.PkgPath == "" {
			fields = append(fields, v)
			if v.Anonymous {
				// also consider fields of anonymous fields as fields of the root
				for _, nested := range deepFields(v.Type) {
					nested.Index = append(append([]int{}, v.Index...), nested.Index...)
					fields = append(fields, nested)
				}
			}
		}
	}

	if cached, loaded := deepFieldsCache.LoadOrStore(reflectType, fields); loaded {
		return cached.([]reflect.StructField)
	}
	return fields
}
