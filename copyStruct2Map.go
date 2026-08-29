/*
 * File: copyStruct2Map.go
 * Created Date: 2022-02-16 12:54:13
 * Author: ysj
 * Description:
 */

package gocopy

import (
	"reflect"
	"time"

	"github.com/iancoleman/strcase"
)

var caseFunc = map[string]func(s string) string{
	"Camel":          strcase.ToCamel,
	"LowerCamel":     strcase.ToLowerCamel,
	"Snake":          strcase.ToSnake,
	"ScreamingSnake": strcase.ToScreamingSnake,
	"Kebab":          strcase.ToKebab,
	"ScreamingKebab": strcase.ToScreamingKebab,
}

// copyStruct2Map copy struct to map
func copyStruct2Map(toValue, fromValue reflect.Value, opt *Option) {
	if beginPointerCopy(fromValue, toValue, opt) {
		return
	}
	fromValue = indirectValue(fromValue)
	toValue = indirectValue(toValue)
	// type
	fromType, _ := indirectType(fromValue.Type())
	toType, _ := indirectType(toValue.Type())
	if toType == mapStringInterfaceType && plainValueType(fromType) && canDirectStructCopy(opt) {
		copyPlainStructToMap(toValue, fromValue, fromType, toType, opt)
		return
	}

	caseName := opt.ToCase
	if caseName == "" {
		caseName = "LowerCamel"
	}
	var fromFields []fieldInfo
	var plans []mapFieldPlan
	if len(opt.NameFromTo) == 0 {
		plans = mapFieldPlans(fromType, toType, caseName)
	} else {
		fromFields = fieldInfos(fromType)
	}
	nativeMap := map[string]interface{}(nil)

	// avoid tomap is nil or not append
	if toValue.IsNil() || (!opt.Append && toValue.CanSet()) {
		fieldCapacity := len(fromFields)
		if plans != nil {
			fieldCapacity = len(plans)
		}
		toNewMap := reflect.MakeMapWithSize(toType, fieldCapacity)
		toValue.Set(toNewMap)
	}
	if toType == mapStringInterfaceType && canDirectStructCopy(opt) {
		nativeMap = toValue.Interface().(map[string]interface{})
	}

	fieldCount := len(fromFields)
	if plans != nil {
		fieldCount = len(plans)
	}
	ignored := ignoredFields(opt)
	for i := 0; i < fieldCount; i++ {
		info := fieldInfo{}
		toFieldName := ""
		toKey := reflect.Value{}
		if plans != nil {
			info = plans[i].info
			toFieldName = plans[i].name
			toKey = plans[i].key
		} else {
			info = fromFields[i]
		}
		fromField := info.field
		// ignore field to skip copy
		if level, ok := ignored[fromField.Name]; ok {
			ignored[fromField.Name]++
			if level <= opt.IgnoreLevel {
				continue
			}
		}

		// field case
		tagName, omitEmpty := info.tagName, info.omitEmpty
		if info.skip {
			continue
		}
		toFieldName = tagName
		if plans != nil {
			toFieldName = plans[i].name
		} else if toFieldName == "" {
			toFieldName = caseFunc[caseName](fromField.Name)
		}
		// the specified toName has the highest priority
		toName, ok := opt.NameFromTo[fromField.Name]
		if ok {
			toFieldName = toName
		}

		// target field name
		if !toKey.IsValid() {
			toKey = reflect.ValueOf(toFieldName)
			if toKey.Type() != toType.Key() {
				toKey = toKey.Convert(toType.Key())
			}
		}

		// whether ignore zero value
		fromFV := fromValue.FieldByIndex(fromField.Index)
		if fromFV.IsZero() && (opt.IgnoreZero || omitEmpty) {
			continue
		}
		// ignore invalid source field
		fromFieldValue := indirectValue(fromFV)
		if !fromFieldValue.IsValid() {
			continue
		}

		fromFieldType, fromVIsPtr := indirectType(fromFV.Type())
		toV := indirectValue(toValue.MapIndex(toKey))

		toVType := fromFieldType
		toVIsPtr := fromVIsPtr
		if toV.IsValid() {
			toVType = reflect.TypeOf(toV.Interface())
			toVType, toVIsPtr = indirectType(toVType)
		}

		// convert field value by customized func
		if convertFunc, ok := opt.Converters[fromField.Name]; ok {
			convertValue := convertFunc(fromFV.Interface())
			if toVIsPtr && toV.IsValid() {
				toFV := indirectValue(reflect.New(toVType))
				toFV.Set(reflect.ValueOf(convertValue))
				toValue.SetMapIndex(toKey, toFV.Addr())
			} else {
				// sometimes convertValue maybe nil
				if convertValue == nil {
					toValue.SetMapIndex(toKey, reflect.Zero(toValue.Type()))
				} else {
					toValue.SetMapIndex(toKey, reflect.ValueOf(convertValue))
				}
			}
			continue
		}

		// specially handle time.Time to string and vice versa
		if timeFieldMap, ok := opt.TimeToString[fromField.Name]; ok {
			timeString := ""
			if timeFieldMap == nil {
				location, err := time.LoadLocation(defaultTimeLoc)
				if err != nil {
					panic(err)
				}
				timeString = fromFieldValue.Interface().(time.Time).In(location).Format(defaultTimeLayout)
			} else {
				loc, ok := timeFieldMap["loc"]
				if !ok {
					loc = defaultTimeLoc
				}
				layout, ok := timeFieldMap["layout"]
				if !ok {
					layout = defaultTimeLayout
				}
				location, err := time.LoadLocation(loc)
				if err != nil {
					panic(err)
				}
				timeString = fromFieldValue.Interface().(time.Time).In(location).Format(layout)
			}
			if toVIsPtr {
				toFV := indirectValue(reflect.New(toVType))
				toFV.Set(reflect.ValueOf(timeString))
				toValue.SetMapIndex(toKey, toFV.Addr())
			} else {
				if nativeMap != nil {
					nativeMap[toFieldName] = timeString
				} else {
					toValue.SetMapIndex(toKey, reflect.ValueOf(timeString))
				}
			}
			continue
		}
		if stringFieldMap, ok := opt.StringToTime[fromField.Name]; ok {
			if fromFieldValue.IsZero() { // ""
				continue
			}
			var timeTime = time.Now()
			if stringFieldMap == nil {
				location, err := time.LoadLocation(defaultTimeLoc)
				if err != nil {
					panic(err)
				}
				timeTime, err = time.ParseInLocation(defaultTimeLayout, fromFieldValue.Interface().(string), location)
				if err != nil {
					panic(err)
				}
			} else {
				loc, ok := stringFieldMap["loc"]
				if !ok {
					loc = defaultTimeLoc
				}
				layout, ok := stringFieldMap["layout"]
				if !ok {
					layout = defaultTimeLayout
				}
				location, err := time.LoadLocation(loc)
				if err != nil {
					panic(err)
				}
				timeTime, err = time.ParseInLocation(layout, fromFieldValue.Interface().(string), location)
				if err != nil {
					panic(err)
				}
			}
			if toVIsPtr {
				toFV := indirectValue(reflect.New(toVType))
				toFV.Set(reflect.ValueOf(timeTime))
				toValue.SetMapIndex(toKey, toFV.Addr())
			} else {
				if nativeMap != nil {
					nativeMap[toFieldName] = timeTime
				} else {
					toValue.SetMapIndex(toKey, reflect.ValueOf(timeTime))
				}
			}
			continue
		}

		switch fromFieldValue.Kind() {
		// slice to slice, need to avoid zero slice
		case reflect.Slice:
			// The common struct-to-map case stores a slice value directly in an
			// interface-valued map. Clone it once instead of allocating a
			// temporary destination and routing through the reflective copier.
			// Append mode still needs the general path to merge existing values.
			fromElemType, _ := indirectType(fromFieldType.Elem())
			shallowElem := fromElemType.Kind() != reflect.Struct && fromElemType.Kind() != reflect.Map && fromElemType.Kind() != reflect.Slice
			if !opt.Append && !fromFieldValue.IsNil() && shallowElem && fromFieldType.AssignableTo(toVType) {
				cloned := reflect.MakeSlice(fromFieldValue.Type(), fromFieldValue.Len(), fromFieldValue.Len())
				reflect.Copy(cloned, fromFieldValue)
				toValue.SetMapIndex(toKey, cloned)
				continue
			}
			toElemType, elemIsPtr := indirectType(toVType.Elem())
			fromElemIsStruct := fromElemType.Kind() == reflect.Struct
			toElemIsStruct := toElemType.Kind() == reflect.Struct

			targetType := toVType
			if fromElemIsStruct || toElemIsStruct {
				if elemIsPtr {
					targetType = reflect.SliceOf(reflect.New(toType).Type())
					// slice of struct to slice of map
				} else {
					targetType = reflect.SliceOf(toType)
				}
			}
			if !opt.Append && fromFieldValue.IsNil() {
				nilSlice := reflect.Zero(targetType)
				if toVIsPtr {
					wrapped := reflect.New(targetType)
					wrapped.Elem().Set(nilSlice)
					toValue.SetMapIndex(toKey, wrapped)
				} else {
					toValue.SetMapIndex(toKey, nilSlice)
				}
				continue
			}
			if !toV.IsValid() { // zero slice
				toV = indirectValue(reflect.New(targetType))
				// toVType is slice of struct, need to transfrom every elem to map/*map
			} else if toElemIsStruct {
				toV = reflect.ValueOf(toV.Interface())
				toSliceOfMap := reflect.MakeSlice(targetType, toV.Len(), toV.Len())
				out := 0
				for i := 0; i < toV.Len(); i++ {
					if !toV.Index(i).IsValid() {
						continue
					}
					toElemV := reflect.New(toType) // map
					if !elemIsPtr {
						toElemV = reflect.MakeMapWithSize(toType, len(fieldInfos(toV.Index(i).Type())))
					}
					copyStruct2Map(toElemV, toV.Index(i), opt)
					if !elemIsPtr {
						toElemV = indirectValue(toElemV)
					}
					toSliceOfMap.Index(out).Set(toElemV)
					out++
				}
				if out != toSliceOfMap.Len() {
					toSliceOfMap = toSliceOfMap.Slice(0, out)
				}
				toV = toSliceOfMap
			}

			dest := indirectValue(reflect.New(targetType))            // slice
			dest.Set(indirectValue(reflect.ValueOf(toV.Interface()))) // slice append slice

			if !fromElemIsStruct {
				copySlice(dest, fromFV, opt)
			} else {
				start := dest.Len()
				fromLen := fromFieldValue.Len()
				if dest.Cap() < start+fromLen {
					grown := reflect.MakeSlice(targetType, start, start+fromLen)
					reflect.Copy(grown, dest)
					dest = grown
				}
				dest = dest.Slice(0, start+fromLen)
				out := start
				for i := 0; i < fromFieldValue.Len(); i++ {
					if !fromFieldValue.Index(i).IsValid() {
						continue
					}
					toElemValue := reflect.New(toType) // map
					if !elemIsPtr {
						toElemValue = reflect.MakeMapWithSize(toType, len(fieldInfos(fromElemType)))
					}
					copyStruct2Map(toElemValue, fromFieldValue.Index(i), opt)
					if !elemIsPtr {
						toElemValue = indirectValue(toElemValue)
					}
					dest.Index(out).Set(toElemValue)
					out++
				}
				if out != start+fromLen {
					dest = dest.Slice(0, out)
				}
			}

			if toVIsPtr {
				dest = dest.Addr()
			}
			toValue.SetMapIndex(toKey, dest)

		// map set kv, need to avoid nil map
		case reflect.Map:
			// A nil source map has no entries to recursively copy. Preserve the
			// existing empty-map result with one allocation instead of creating
			// temporary maps that copyMap immediately replaces.
			if !opt.Append && fromFieldValue.IsNil() {
				empty := reflect.MakeMapWithSize(toType, 0)
				if toVIsPtr {
					wrapped := reflect.New(toType)
					wrapped.Elem().Set(empty)
					toValue.SetMapIndex(toKey, wrapped)
				} else {
					toValue.SetMapIndex(toKey, empty)
				}
				continue
			}
			fromElemType, _ := indirectType(fromFieldType.Elem())
			toElemType, elemIsPtr := indirectType(toVType.Elem())
			fromElemIsStruct := fromElemType.Kind() == reflect.Struct
			toElemIsStruct := toElemType.Kind() == reflect.Struct
			if !opt.Append && !toV.IsValid() && fromElemIsStruct {
				result := reflect.New(toType)
				resultMap := indirectValue(result)
				resultMap.Set(reflect.MakeMapWithSize(toType, fromFieldValue.Len()))
				iter := fromFieldValue.MapRange()
				for iter.Next() {
					value := reflect.New(toType)
					copyStruct2Map(value, iter.Value(), opt)
					if !elemIsPtr {
						value = indirectValue(value)
					}
					resultMap.SetMapIndex(iter.Key(), value)
				}
				if toVIsPtr {
					toValue.SetMapIndex(toKey, result)
				} else {
					toValue.SetMapIndex(toKey, resultMap)
				}
				continue
			}

			if !toV.IsValid() { // zero map
				toV = indirectValue(reflect.New(toType))
				if toV.IsNil() { // nil map
					toNewV := reflect.MakeMapWithSize(toType, toV.Len())
					toV.Set(toNewV)
				}
			} else if toElemIsStruct {
				toMap := indirectValue(reflect.New(toType))
				toV = reflect.ValueOf(toV.Interface())
				if toMap.IsNil() {
					toNewMap := reflect.MakeMapWithSize(toType, toV.Len())
					toMap.Set(toNewMap)
				}
				toVIter := toV.MapRange()
				for toVIter.Next() {
					toVK := toVIter.Key()
					toVV := toVIter.Value()            // struct
					toElemValue := reflect.New(toType) // map
					copyStruct2Map(toElemValue, toVV, opt)
					if !elemIsPtr {
						toElemValue = indirectValue(toElemValue)
					}
					toMap.SetMapIndex(toVK, toElemValue)
				}
				toV = toMap
			}

			dest := indirectValue(reflect.New(toType))
			// avoid tomap is nil
			if dest.IsNil() {
				toV = reflect.ValueOf(toV.Interface())
				toNewDest := reflect.MakeMapWithSize(toType, (toV.Len() + fromFieldValue.Len()))
				dest.Set(toNewDest)
			}
			toVIter := toV.MapRange()
			for toVIter.Next() {
				toK := toVIter.Key()
				toV := toVIter.Value()
				dest.SetMapIndex(toK, toV)
			}

			if !fromElemIsStruct {
				copyMap(dest, fromFV, opt)
			} else {
				fromKVIter := fromFieldValue.MapRange()
				for fromKVIter.Next() {
					fromK := fromKVIter.Key()
					fromV := fromKVIter.Value()        // struct
					toElemValue := reflect.New(toType) // map
					copyStruct2Map(toElemValue, fromV, opt)
					if !elemIsPtr {
						toElemValue = indirectValue(toElemValue)
					}
					dest.SetMapIndex(fromK, toElemValue)
				}
			}
			if toVIsPtr {
				dest = dest.Addr()
			}
			toValue.SetMapIndex(toKey, dest)

		case reflect.Struct:
			// if time.Time field
			_, ok := fromFieldValue.Interface().(time.Time)
			if ok {
				if toVIsPtr {
					toV = indirectValue(reflect.New(toVType))
					toV.Set(fromFieldValue)
					toV = toV.Addr()
					toValue.SetMapIndex(toKey, toV)
				} else {
					toValue.SetMapIndex(toKey, fromFieldValue)
				}
				continue
			}
			// other struct field
			if !toV.IsValid() { // zero map
				toV = indirectValue(reflect.New(toType))
			} else if toVType.Kind() == reflect.Struct {
				toNewV := reflect.New(toType) // map
				copyStruct2Map(toNewV, reflect.ValueOf(toV.Interface()), opt)
				toV = toNewV
			} else if toVType.Kind() == reflect.Map {
				toNewV := reflect.New(toType) // map
				copyMap(toNewV, reflect.ValueOf(toV.Interface()), opt)
				toV = toNewV
			}

			dest := indirectValue(reflect.New(toType))
			dest.Set(indirectValue(reflect.ValueOf(toV.Interface())))
			copyStruct2Map(dest, fromFV, opt)
			if toVIsPtr {
				dest = dest.Addr()
			}
			toValue.SetMapIndex(toKey, dest)

		default:
			if toVIsPtr {
				toV = indirectValue(reflect.New(toVType))
				toV.Set(fromFieldValue)
				toV = toV.Addr()
				toValue.SetMapIndex(toKey, toV)
			} else if nativeMap != nil {
				nativeMap[toFieldName] = fromFieldValue.Interface()
			} else {
				toValue.SetMapIndex(toKey, fromFieldValue)
			}
		}
	}
}

func copyPlainStructToMap(toValue, fromValue reflect.Value, fromType, toType reflect.Type, opt *Option) {
	caseName := opt.ToCase
	if caseName == "" {
		caseName = "LowerCamel"
	}
	plans := mapFieldPlans(fromType, toType, caseName)
	if toValue.IsNil() || (!opt.Append && toValue.CanSet()) {
		toValue.Set(reflect.MakeMapWithSize(toType, len(plans)))
	}
	destination := toValue.Interface().(map[string]interface{})
	for _, plan := range plans {
		if plan.info.skip {
			continue
		}
		value := fromValue.FieldByIndex(plan.info.field.Index)
		if value.IsValid() {
			if plan.info.omitEmpty && value.IsZero() {
				continue
			}
			destination[plan.name] = value.Interface()
		}
	}
}
