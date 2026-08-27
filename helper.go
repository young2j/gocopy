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
	deepFieldsCache   sync.Map // map[reflect.Type][]reflect.StructField
	fieldInfoCache    sync.Map // map[reflect.Type][]fieldInfo
	fieldCache        sync.Map // map[fieldCacheKey]fieldCacheValue
	fieldPairCache    sync.Map // map[fieldPairCacheKey][]fieldPair
	mapFieldPlanCache sync.Map // map[mapFieldPlanKey][]mapFieldPlan
	plainTypeCache    sync.Map // map[reflect.Type]bool
	copyStatePool     = sync.Pool{New: func() interface{} { return &copyState{} }}
)

type fieldCacheKey struct {
	typ  reflect.Type
	name string
}

type fieldCacheValue struct {
	field reflect.StructField
	ok    bool
}

type fieldPairCacheKey struct {
	from reflect.Type
	to   reflect.Type
}

type fieldPair struct {
	from     reflect.StructField
	to       reflect.StructField
	ok       bool
	fromType reflect.Type
	toType   reflect.Type
	toPtr    bool
	mode     fieldCopyMode
}

type fieldCopyMode uint8

const (
	fieldCopyNested fieldCopyMode = iota
	fieldCopyAssignable
	fieldCopyConvertible
)

type fieldInfo struct {
	field      reflect.StructField
	tagName    string
	lowerCamel string
	omitEmpty  bool
	skip       bool
}

type mapFieldPlanKey struct {
	from     reflect.Type
	to       reflect.Type
	caseName string
}

type mapFieldPlan struct {
	info fieldInfo
	name string
	key  reflect.Value
}

type copyState struct {
	entries      [8]visitEntry
	count        int
	overflow     map[visitKey]reflect.Value
	ignoreFields map[string]int
	reverseNames map[string][]string
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
	state.ignoreFields = nil
	state.reverseNames = nil
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

func fieldInfos(reflectType reflect.Type) []fieldInfo {
	reflectType, _ = indirectType(reflectType)
	if cached, ok := fieldInfoCache.Load(reflectType); ok {
		return cached.([]fieldInfo)
	}
	fields := deepFields(reflectType)
	infos := make([]fieldInfo, len(fields))
	for i, field := range fields {
		name, omitEmpty, skip := fieldTag(field)
		infos[i] = fieldInfo{field: field, tagName: name, lowerCamel: caseFunc["LowerCamel"](field.Name), omitEmpty: omitEmpty, skip: skip}
	}
	if cached, loaded := fieldInfoCache.LoadOrStore(reflectType, infos); loaded {
		return cached.([]fieldInfo)
	}
	return infos
}

func mapFieldPlans(fromType, toType reflect.Type, caseName string) []mapFieldPlan {
	fromType, _ = indirectType(fromType)
	toType, _ = indirectType(toType)
	key := mapFieldPlanKey{from: fromType, to: toType, caseName: caseName}
	if cached, ok := mapFieldPlanCache.Load(key); ok {
		return cached.([]mapFieldPlan)
	}
	keyType := toType.Key()
	infos := fieldInfos(fromType)
	plans := make([]mapFieldPlan, 0, len(infos))
	convertCase := caseFunc[caseName]
	for _, info := range infos {
		if info.skip {
			plans = append(plans, mapFieldPlan{info: info})
			continue
		}
		name := info.tagName
		if name == "" {
			if caseName == "LowerCamel" {
				name = info.lowerCamel
			} else {
				name = convertCase(info.field.Name)
			}
		}
		mapKey := reflect.ValueOf(name)
		if mapKey.Type() != keyType {
			mapKey = mapKey.Convert(keyType)
		}
		plans = append(plans, mapFieldPlan{info: info, name: name, key: mapKey})
	}
	if cached, loaded := mapFieldPlanCache.LoadOrStore(key, plans); loaded {
		return cached.([]mapFieldPlan)
	}
	return plans
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

func fieldPairs(fromType, toType reflect.Type) []fieldPair {
	fromType, _ = indirectType(fromType)
	toType, _ = indirectType(toType)
	key := fieldPairCacheKey{from: fromType, to: toType}
	if cached, ok := fieldPairCache.Load(key); ok {
		return cached.([]fieldPair)
	}

	fields := deepFields(fromType)
	pairs := make([]fieldPair, 0, len(fields))
	for _, from := range fields {
		if to, ok := fieldByName(toType, from.Name); ok {
			fromValueType, _ := indirectType(from.Type)
			toValueType, toPtr := indirectType(to.Type)
			mode := fieldCopyNested
			if fromValueType.AssignableTo(toValueType) {
				mode = fieldCopyAssignable
			} else if fromValueType.ConvertibleTo(toValueType) {
				mode = fieldCopyConvertible
			}
			pairs = append(pairs, fieldPair{from: from, to: to, ok: true, fromType: fromValueType, toType: toValueType, toPtr: toPtr, mode: mode})
		} else {
			pairs = append(pairs, fieldPair{from: from})
		}
	}
	if cached, loaded := fieldPairCache.LoadOrStore(key, pairs); loaded {
		return cached.([]fieldPair)
	}
	return pairs
}

func plainValueType(reflectType reflect.Type) bool {
	if cached, ok := plainTypeCache.Load(reflectType); ok {
		return cached.(bool)
	}
	plain := true
	switch reflectType.Kind() {
	case reflect.Array:
		plain = plainValueType(reflectType.Elem())
	case reflect.Struct:
		for i := 0; i < reflectType.NumField(); i++ {
			if !plainValueType(reflectType.Field(i).Type) {
				plain = false
				break
			}
		}
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String,
		reflect.UnsafePointer:
	default:
		plain = false
	}
	if cached, loaded := plainTypeCache.LoadOrStore(reflectType, plain); loaded {
		return cached.(bool)
	}
	return plain
}

func canDirectStructCopy(opt *Option) bool {
	return opt != nil && !opt.Append && !opt.IgnoreZero &&
		len(opt.NameFromTo) == 0 && len(opt.Converters) == 0 &&
		len(opt.TimeToString) == 0 && len(opt.StringToTime) == 0 &&
		len(opt.IgnoreFields) == 0
}

func ignoredFields(opt *Option) map[string]int {
	if opt == nil || opt.state == nil {
		return nil
	}
	return opt.state.ignoreFields
}

func reverseNames(opt *Option) map[string][]string {
	if opt == nil || opt.state == nil {
		return nil
	}
	return opt.state.reverseNames
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
