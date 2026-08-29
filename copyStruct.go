/*
 * File: copyStruct.go
 * Created Date: 2022-01-26 06:15:08
 * Author: ysj
 * Description:  copy struct to struct
 */
package gocopy

import (
	"reflect"
	"time"
)

var (
	defaultTimeLoc    = "Asia/Shanghai"
	defaultTimeLayout = "2006-01-02 15:04:05"
)

// copyStruct copies exported fields by name, applying rename, time, converter,
// ignore, and append options before recursively handling nested values.
func copyStruct(toValue, fromValue reflect.Value, opt *Option) {
	if beginPointerCopy(fromValue, toValue, opt) {
		return
	}
	fromValue = indirectValue(fromValue)
	toValue = indirectValue(toValue)
	// type
	fromType, _ := indirectType(fromValue.Type())
	toType, _ := indirectType(toValue.Type())
	if fromType == toType && plainValueType(fromType) && canDirectStructCopy(opt) && toValue.CanSet() {
		toValue.Set(fromValue)
		return
	}

	// Handle every struct field. The common path has no rename rules, so use a
	// cached source/target field pairing instead of a per-field FieldByName.
	var pairs []fieldPair
	if len(opt.NameFromTo) == 0 {
		pairs = fieldPairs(fromType, toType)
	}
	fromFields := deepFields(fromType)
	ignored := ignoredFields(opt)
	for i := 0; i < len(fromFields); i++ {
		fromField := fromFields[i]
		// ignore field to skip copy
		if level, ok := ignored[fromField.Name]; ok {
			ignored[fromField.Name]++
			if level <= opt.IgnoreLevel {
				continue
			}
		}
		// from field to field
		var toField reflect.StructField
		var ok bool
		if len(opt.NameFromTo) == 0 {
			if i >= len(pairs) || !pairs[i].ok {
				continue
			}
			toField = pairs[i].to
			ok = true
		} else {
			toFieldName, mapped := opt.NameFromTo[fromField.Name]
			if !mapped {
				toFieldName = fromField.Name
			}
			toField, ok = fieldByName(toType, toFieldName)
		}
		if !ok {
			continue
		}
		fromFValue := fromValue.FieldByIndex(fromField.Index)
		fromFieldValue := indirectValue(fromFValue)
		toFieldValue := toValue.FieldByIndex(toField.Index)

		fromFieldType, toFieldType := fromField.Type, toField.Type
		toFieldIsPtr := false
		mode := fieldCopyNested
		if len(opt.NameFromTo) == 0 && i < len(pairs) && pairs[i].ok {
			fromFieldType = pairs[i].fromType
			toFieldType = pairs[i].toType
			toFieldIsPtr = pairs[i].toPtr
			mode = pairs[i].mode
		} else {
			fromFieldType, _ = indirectType(fromFieldType)
			toFieldType, toFieldIsPtr = indirectType(toFieldType)
			if fromFieldType.AssignableTo(toFieldType) {
				mode = fieldCopyAssignable
			} else if fromFieldType.ConvertibleTo(toFieldType) {
				mode = fieldCopyConvertible
			}
		}

		if !toFieldValue.CanSet() || !fromFieldValue.IsValid() {
			continue
		}

		// convert field value by customized func
		if convertFunc, ok := opt.Converters[fromField.Name]; ok {
			convertValue := convertFunc(fromFValue.Interface())

			if toFieldIsPtr {
				toFV := indirectValue(reflect.New(toFieldType))
				toFV.Set(reflect.ValueOf(convertValue))
				toFieldValue.Set(toFV.Addr())
			} else {
				// sometimes convertValue maybe nil
				if convertValue == nil {
					toFieldValue.Set(reflect.Zero(toFieldValue.Type()))
				} else {
					toFieldValue.Set(reflect.ValueOf(convertValue))
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
			if toFieldIsPtr {
				toFV := indirectValue(reflect.New(toFieldType))
				toFV.Set(reflect.ValueOf(timeString))
				toFieldValue.Set(toFV.Addr())
			} else {
				toFieldValue.Set(reflect.ValueOf(timeString))
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
			if toFieldIsPtr {
				toFV := indirectValue(reflect.New(toFieldType))
				toFV.Set(reflect.ValueOf(timeTime))
				toFieldValue.Set(toFV.Addr())
			} else {
				toFieldValue.Set(reflect.ValueOf(timeTime))
			}
			continue
		}

		// can direct assign
		// 可直接赋值拷贝
		if mode == fieldCopyAssignable {
			if !opt.Append { // not append
				if toFieldIsPtr {
					// like string -> *string
					if fromFieldValue.CanAddr() {
						toFieldValue.Set(fromFieldValue.Addr())
					} else {
						fromFV := indirectValue(reflect.New(fromFieldType))
						fromFV.Set(fromFieldValue)
						toFieldValue.Set(fromFV.Addr())
					}
				} else {
					toFieldValue.Set(fromFieldValue)
				}
			} else { // append mode
				if toFieldIsPtr && toFieldValue.IsNil() {
					toFieldValue.Set(reflect.New(toFieldType))
				}
				fromFieldKind := fromFieldType.Kind()
				switch fromFieldKind {
				// slice append slice, need to avoid zero slice
				case reflect.Slice:
					copySlice(toFieldValue, fromFValue, opt)
				// map set kv, need to avoid nil map
				case reflect.Map:
					copyMap(toFieldValue, fromFValue, opt)
				default:
					if toFieldIsPtr {
						// like string -> *string
						if fromFieldValue.CanAddr() {
							toFieldValue.Set(fromFieldValue.Addr())
						} else {
							fromFV := indirectValue(reflect.New(fromFieldType))
							fromFV.Set(fromFieldValue)
							toFieldValue.Set(fromFV.Addr())
						}
					} else {
						toFieldValue.Set(fromFieldValue)
					}
				}
			}
			// can convert field type
			// 类型可转换拷贝
		} else if mode == fieldCopyConvertible {
			convertValue := fromFieldValue.Convert(toFieldType)
			if !opt.Append { // not append
				if toFieldIsPtr {
					// like string -> *string
					if convertValue.CanAddr() {
						toFieldValue.Set(convertValue.Addr())
					} else {
						convertFV := indirectValue(reflect.New(toFieldType))
						convertFV.Set(convertValue)
						toFieldValue.Set(convertFV.Addr())
					}
				} else {
					toFieldValue.Set(convertValue) // set to converted value
				}
			} else { // append mode
				if toFieldIsPtr && toFieldValue.IsNil() {
					toFieldValue.Set(reflect.New(toFieldType))
				}
				fromFieldKind := fromFieldType.Kind()
				switch fromFieldKind {
				// slice append slice, need to avoid zero slice
				case reflect.Slice:
					copySlice(toFieldValue, convertValue, opt)
				// map set kv, need to avoid nil map
				case reflect.Map:
					copyMap(toFieldValue, convertValue, opt)
				default:
					if toFieldIsPtr {
						// like string -> *string
						if convertValue.CanAddr() {
							toFieldValue.Set(convertValue.Addr())
						} else {
							convertFV := indirectValue(reflect.New(toFieldType))
							convertFV.Set(convertValue)
							toFieldValue.Set(convertFV.Addr())
						}
					} else {
						toFieldValue.Set(convertValue) // set to converted value
					}
				}
			}
		} else {
			// can not directly assign or convert
			if toFieldIsPtr && toFieldValue.IsNil() {
				toFieldValue.Set(reflect.New(toFieldType))
			}
			fromFieldKind := fromFieldType.Kind()
			toFieldKind := toFieldType.Kind()
			// 1. slice to slice
			if toFieldKind == reflect.Slice && fromFieldKind == reflect.Slice {
				copySlice(toFieldValue, fromFValue, opt)
				// 2. map to map
			} else if toFieldKind == reflect.Map && fromFieldKind == reflect.Map {
				copyMap(toFieldValue, fromFValue, opt)
				// 3. struct to struct
			} else if toFieldKind == reflect.Struct && fromFieldKind == reflect.Struct {
				copyStruct(toFieldValue, fromFValue, opt)
			}
		}
	}
}
