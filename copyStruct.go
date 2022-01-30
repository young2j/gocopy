/*
 * File: copyStruct.go
 * Created Date: 2022-01-26 06:15:08
 * Author: ysj
 * Description:  copy struct to struct
 */
package gocopy

import (
	"reflect"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func copyStruct(toValue, fromValue reflect.Value, opt *Option) {
	fromType := indirectType(fromValue.Type())
	toType := indirectType(toValue.Type())

	// handle every struct field
	fromFields := deepFields(fromType)
	for i := 0; i < len(fromFields); i++ {
		// for i := 0; i < fromType.NumField(); i++ {
		// fromField := fromType.Field(i)
		fromField := fromFields[i]
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

		fromFieldValue := indirectValue(fromValue.FieldByName(fromField.Name))
		toFieldValue := toValue.FieldByName(toField.Name)

		// pointer value like *string *int
		if toFieldValue.Kind() == reflect.Ptr && toFieldValue.IsNil() {
			toNewValue := reflect.New(indirectType(toField.Type))
			toFieldValue.Set(toNewValue)
		}

		if !toFieldValue.CanSet() || !fromFieldValue.IsValid() {
			continue
		}

		// specially handle bson.ObjectId to string and vice versa
		if objectIdType, ok := opt.ObjectIdToString[fromField.Name]; ok {
			if objectIdType == "mgo" {
				objectId, ok := fromFieldValue.Interface().(bson.ObjectId)
				if ok {
					toFieldValue.Set(reflect.ValueOf(objectId.Hex()))
					continue
				}
			} else if objectIdType == "official" {
				objectId, ok := fromFieldValue.Interface().(primitive.ObjectID)
				if ok {
					toFieldValue.Set(reflect.ValueOf(objectId.Hex()))
					continue
				}
			}
		}
		if objectIdType, ok := opt.StringToObjectId[fromField.Name]; ok {
			if objectIdType == "mgo" {
				objectId := bson.ObjectIdHex(fromFieldValue.String())
				toFieldValue.Set(reflect.ValueOf(objectId))
				continue
			} else if objectIdType == "official" {
				if objectId, err := primitive.ObjectIDFromHex(fromFieldValue.String()); err == nil {
					toFieldValue.Set(reflect.ValueOf(objectId))
					continue
				}
			}
		}

		fromFieldType := indirectType(fromField.Type)
		toFieldType := indirectType(toField.Type)
		// can direct assign
		// 可直接赋值拷贝
		if fromFieldType.AssignableTo(toFieldType) {
			if !opt.Append { // not append
				if toFieldValue.Kind() == reflect.Ptr {
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
					if toFieldValue.Kind() == reflect.Ptr {
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
				toFieldValue.Set(convertValue) // set to converted value
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
					toFieldValue.Set(convertValue) // set to converted value
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
