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

	"github.com/globalsign/mgo/bson"
	"github.com/iancoleman/strcase"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	fromValue = indirectValue(fromValue)
	toValue = indirectValue(toValue)
	// type
	fromType, _ := indirectType(fromValue.Type())
	toType, _ := indirectType(toValue.Type())

	fromFields := deepFields(fromType)

	// avoid tomap is nil or not append
	if toValue.IsNil() || !opt.Append {
		toNewMap := reflect.MakeMapWithSize(toType, len(fromFields))
		toValue.Set(toNewMap)
	}

	for i := 0; i < len(fromFields); i++ {
		fromField := fromFields[i]
		// ignore field to skip copy
		if level, ok := opt.ignoreFields[fromField.Name]; ok {
			opt.ignoreFields[fromField.Name]++
			if level <= opt.IgnoreLevel {
				continue
			}
		}

		// field case
		if opt.ToCase == "" {
			opt.ToCase = "LowerCamel"
		}
		toFieldName := caseFunc[opt.ToCase](fromField.Name)
		// the specified toName has the highest priority
		toName, ok := opt.NameFromTo[fromField.Name]
		if ok {
			toFieldName = toName
		}

		// target field name
		toKey := reflect.ValueOf(toFieldName)

		// whether ignore zero value
		fromFV := fromValue.FieldByName(fromField.Name)
		if fromFV.IsZero() && opt.IgnoreZero {
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

		// specially handle bson.ObjectId to string and vice versa
		if objectIdType, ok := opt.ObjectIdToString[fromField.Name]; ok {
			if objectIdType == "mgo" {
				objectId, ok := fromFieldValue.Interface().(bson.ObjectId)
				if ok {
					if toVIsPtr {
						toFV := indirectValue(reflect.New(toVType))
						toFV.Set(reflect.ValueOf(objectId.Hex()))
						toValue.SetMapIndex(toKey, toFV.Addr())
					} else {
						toValue.SetMapIndex(toKey, reflect.ValueOf(objectId.Hex()))
					}
					continue
				}
			} else if objectIdType == "official" {
				objectId, ok := fromFieldValue.Interface().(primitive.ObjectID)
				if ok {
					if toVIsPtr {
						toFV := indirectValue(reflect.New(toVType))
						toFV.Set(reflect.ValueOf(objectId.Hex()))
						toValue.SetMapIndex(toKey, toFV.Addr())
					} else {
						toValue.SetMapIndex(toKey, reflect.ValueOf(objectId.Hex()))
					}
					continue
				}
			}
		}
		if objectIdType, ok := opt.StringToObjectId[fromField.Name]; ok {
			if objectIdType == "mgo" && bson.IsObjectIdHex(fromFieldValue.String()) {
				objectId := bson.ObjectIdHex(fromFieldValue.String())
				if toVIsPtr {
					toFV := indirectValue(reflect.New(toVType))
					toFV.Set(reflect.ValueOf(objectId))
					toValue.SetMapIndex(toKey, toFV.Addr())
				} else {
					toValue.SetMapIndex(toKey, reflect.ValueOf(objectId))
				}
				continue
			} else if objectIdType == "official" {
				if objectId, err := primitive.ObjectIDFromHex(fromFieldValue.String()); err == nil {
					if toVIsPtr {
						toFV := indirectValue(reflect.New(toVType))
						toFV.Set(reflect.ValueOf(objectId))
						toValue.SetMapIndex(toKey, toFV.Addr())
					} else {
						toValue.SetMapIndex(toKey, reflect.ValueOf(objectId))
					}
					continue
				}
			}
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
				toValue.SetMapIndex(toKey, reflect.ValueOf(timeString))
			}
			continue
		}
		if stringFieldMap, ok := opt.StringToTime[fromField.Name]; ok {
			if fromFieldValue.IsZero() { // ""
				continue
			}
			timeTime := time.Now()
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
				toValue.SetMapIndex(toKey, reflect.ValueOf(timeTime))
			}
			continue
		}

		switch fromFieldValue.Kind() {
		// slice to slice, need to avoid zero slice
		case reflect.Slice:
			fromElemType, _ := indirectType(fromFieldType.Elem())
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
			if !toV.IsValid() { // zero slice
				toV = indirectValue(reflect.New(targetType))
				// toVType is slice of struct, need to transfrom every elem to map/*map
			} else if toElemIsStruct {
				toSliceOfMap := indirectValue(reflect.New(targetType))
				toV = reflect.ValueOf(toV.Interface())
				for i := 0; i < toV.Len(); i++ {
					if !toV.Index(i).IsValid() {
						continue
					}
					toElemV := reflect.New(toType) // map
					copyStruct2Map(toElemV, toV.Index(i), opt)
					if !elemIsPtr {
						toElemV = indirectValue(toElemV)
					}
					toSliceOfMap.Set(reflect.Append(toSliceOfMap, toElemV))
				}
				toV = toSliceOfMap
			}

			dest := indirectValue(reflect.New(targetType))            // slice
			dest.Set(indirectValue(reflect.ValueOf(toV.Interface()))) // slice append slice

			if !fromElemIsStruct {
				copySlice(dest, fromFieldValue, opt)
			} else {
				for i := 0; i < fromFieldValue.Len(); i++ {
					if !fromFieldValue.Index(i).IsValid() {
						continue
					}
					toElemValue := reflect.New(toType) // map
					copyStruct2Map(toElemValue, fromFieldValue.Index(i), opt)
					if !elemIsPtr {
						toElemValue = indirectValue(toElemValue)
					}
					dest.Set(reflect.Append(dest, toElemValue))
				}
			}

			if toVIsPtr {
				dest = dest.Addr()
			}
			toValue.SetMapIndex(toKey, dest)

		// map set kv, need to avoid nil map
		case reflect.Map:
			fromElemType, _ := indirectType(fromFieldType.Elem())
			toElemType, elemIsPtr := indirectType(toVType.Elem())
			fromElemIsStruct := fromElemType.Kind() == reflect.Struct
			toElemIsStruct := toElemType.Kind() == reflect.Struct

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
				copyMap(dest, fromFieldValue, opt)
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
			copyStruct2Map(dest, fromFieldValue, opt)
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
			} else {
				toValue.SetMapIndex(toKey, fromFieldValue)
			}
		}
	}
}
