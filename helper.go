/*
 * File: helper.go
 * Created Date: 2022-01-26 05:48:28
 * Author: ysj
 * Description: helper func
 */
package gocopy

import (
	"reflect"
	"sync"
)

var (
	// ponytail: process-wide caches; add eviction only if dynamic plugins create unbounded types.
	deepFieldsCache sync.Map // map[reflect.Type][]reflect.StructField
	fieldCache      sync.Map // map[fieldCacheKey]fieldCacheValue
)

type fieldCacheKey struct {
	typ  reflect.Type
	name string
}

type fieldCacheValue struct {
	field reflect.StructField
	ok    bool
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
