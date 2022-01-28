/*
 * File: copyStruct.go
 * Created Date: 2022-01-26 06:15:08
 * Author: ysj
 * Description:  copy struct to struct
 */
package copy

import (
	"reflect"

	"github.com/globalsign/mgo/bson"
)

func copyStruct(fromValue, toValue reflect.Value, opt *ConvertOption) {
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

		fromFieldType := indirectType(fromField.Type)
		toFieldType := indirectType(toField.Type)
		// can direct assign
		// 可直接赋值拷贝
		if fromFieldType.AssignableTo(toFieldType) {
			if toFieldValue.Kind() == reflect.Ptr {
				// string -> *string
				toFieldValue.Set(fromFieldValue.Addr())
			} else {
				toFieldValue.Set(fromFieldValue)
			}
			// can convert field type
			// 类型可转换拷贝
		} else if fromFieldType.ConvertibleTo(toFieldType) {
			convertValue := fromFieldValue.Convert(toFieldType)
			// specially handle bson.ObjectId to string and vice versa
			if opt.ObjectIdToString == fromField.Name {
				objectId, ok := fromFieldValue.Interface().(bson.ObjectId)
				if ok {
					convertValue = reflect.ValueOf(objectId.Hex())
				}
			}
			if opt.StringToObjectId == fromField.Name {
				objectId := bson.ObjectIdHex(fromFieldValue.String())
				convertValue = reflect.ValueOf(objectId)
			}

			toFieldValue.Set(convertValue) // set to converted value

		} else {
			// can not directly assign or convert
			fromFieldKind := fromFieldType.Kind()
			toFieldKind := toFieldType.Kind()
			// 1. slice to slice
			if toFieldKind == reflect.Slice && fromFieldKind == reflect.Slice {
				copySlice(fromFieldValue, toFieldValue, opt)
				// 2. map to map
			} else if toFieldKind == reflect.Map && fromFieldKind == reflect.Map {
				copyMap(fromFieldValue, toFieldValue, opt)
				// 3. struct to struct
			} else if toFieldKind == reflect.Struct && fromFieldKind == reflect.Struct {
				copyStruct(fromFieldValue, toFieldValue, opt)
			}
		}
	}
}
