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

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	defaultTimeLoc    = "Asia/Shanghai"
	defaultTimeLayout = "2006-01-02 15:04:05"
)

// copyStruct copy struct to struct
func copyStruct(toValue, fromValue reflect.Value, opt *Option) {
	fromValue = indirectValue(fromValue)
	toValue = indirectValue(toValue)
	// type
	fromType, _ := indirectType(fromValue.Type())
	toType, _ := indirectType(toValue.Type())

	// handle every struct field
	fromFields := deepFields(fromType)
	for i := 0; i < len(fromFields); i++ {
		fromField := fromFields[i]
		// ignore field to skip copy
		if _, ok := opt.ignoreFields[fromField.Name]; ok {
			continue
		}
		// from field to field
		toFieldName, ok := opt.NameFromTo[fromField.Name]
		if !ok {
			toFieldName = fromField.Name
		}
		// whether field is in toType
		toField, ok := toType.FieldByName(toFieldName)
		if !ok {
			continue
		}
		fromFValue := fromValue.FieldByName(fromField.Name)
		fromFieldValue := indirectValue(fromFValue)
		toFieldValue := toValue.FieldByName(toField.Name)

		fromFieldType, _ := indirectType(fromField.Type)
		toFieldType, toFieldIsPtr := indirectType(toField.Type)

		// pointer value like *string *int
		if toFieldIsPtr && toFieldValue.IsNil() {
			toNewValue := reflect.New(toFieldType)
			toFieldValue.Set(toNewValue)
		}

		if !toFieldValue.CanSet() || !fromFieldValue.IsValid() {
			continue
		}

		// avoid tofield is zero slice/map, set on zero value will panic
		if toFieldValue.IsZero() {
			if toFieldIsPtr {
				toFieldValue.Set(reflect.New(toField.Type))
			} else {
				toFieldValue.Set(indirectValue(reflect.New(toField.Type)))
			}
		}

		// convert field value by customized func
		if convertFunc, ok := opt.Converters[fromField.Name]; ok {
			convertValue := convertFunc(fromFValue.Interface())
			if toFieldIsPtr {
				toFV := indirectValue(reflect.New(toFieldType))
				toFV.Set(reflect.ValueOf(convertValue))
				toFieldValue.Set(toFV.Addr())
			} else {
				toFieldValue.Set(reflect.ValueOf(convertValue))
			}
			continue
		}

		// specially handle bson.ObjectId to string and vice versa
		if objectIdType, ok := opt.ObjectIdToString[fromField.Name]; ok {
			if objectIdType == "mgo" {
				objectId, ok := fromFieldValue.Interface().(bson.ObjectId)
				if ok {
					if toFieldIsPtr {
						toFV := indirectValue(reflect.New(toFieldType))
						toFV.Set(reflect.ValueOf(objectId.Hex()))
						toFieldValue.Set(toFV.Addr())
					} else {
						toFieldValue.Set(reflect.ValueOf(objectId.Hex()))
					}
					continue
				}
			} else if objectIdType == "official" {
				objectId, ok := fromFieldValue.Interface().(primitive.ObjectID)
				if ok {
					if toFieldIsPtr {
						toFV := indirectValue(reflect.New(toFieldType))
						toFV.Set(reflect.ValueOf(objectId.Hex()))
						toFieldValue.Set(toFV.Addr())
					} else {
						toFieldValue.Set(reflect.ValueOf(objectId.Hex()))
					}
					continue
				}
			}
		}
		if objectIdType, ok := opt.StringToObjectId[fromField.Name]; ok {
			if objectIdType == "mgo" {
				objectId := bson.ObjectIdHex(fromFieldValue.String())
				if toFieldIsPtr {
					toFV := indirectValue(reflect.New(toFieldType))
					toFV.Set(reflect.ValueOf(objectId))
					toFieldValue.Set(toFV.Addr())
				} else {
					toFieldValue.Set(reflect.ValueOf(objectId))
				}
				continue
			} else if objectIdType == "official" {
				if objectId, err := primitive.ObjectIDFromHex(fromFieldValue.String()); err == nil {
					if toFieldIsPtr {
						toFV := indirectValue(reflect.New(toFieldType))
						toFV.Set(reflect.ValueOf(objectId))
						toFieldValue.Set(toFV.Addr())
					} else {
						toFieldValue.Set(reflect.ValueOf(objectId))
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
		if fromFieldType.AssignableTo(toFieldType) {
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
				fromFieldKind := fromFieldType.Kind()
				switch fromFieldKind {
				// slice append slice, need to avoid zero slice
				case reflect.Slice:
					copySlice(toFieldValue, fromFieldValue, opt)
				// map set kv, need to avoid nil map
				case reflect.Map:
					copyMap(toFieldValue, fromFieldValue, opt)
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
		} else if fromFieldType.ConvertibleTo(toFieldType) {
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
			fromFieldKind := fromFieldType.Kind()
			toFieldKind := toFieldType.Kind()
			// 1. slice to slice
			if toFieldKind == reflect.Slice && fromFieldKind == reflect.Slice {
				copySlice(toFieldValue, fromFieldValue, opt)
				// 2. map to map
			} else if toFieldKind == reflect.Map && fromFieldKind == reflect.Map {
				copyMap(toFieldValue, fromFieldValue, opt)
				// 3. struct to struct
			} else if toFieldKind == reflect.Struct && fromFieldKind == reflect.Struct {
				copyStruct(toFieldValue, fromFieldValue, opt)
			}
		}
	}
}
