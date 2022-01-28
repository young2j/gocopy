/*
 * File: copy.go
 * Created Date: 2022-01-26 05:36:00
 * Author: ysj
 * Description: copy slice and map by type; copy struct by field name
 */
package copy

import (
	"reflect"
)

type ConvertOption struct {
	NameFromTo       map[string]string
	ObjectIdToString string
	StringToObjectId string
	Append           bool
}

func Copy(from, to interface{}) {
	CopyWithOption(from, to, &ConvertOption{Append: false})
}

func CopyWithOption(from, to interface{}, opt *ConvertOption) {
	fromValue := indirectValue(reflect.ValueOf(from))
	toValue := indirectValue(reflect.ValueOf(to))

	fromKind := indirectType(reflect.TypeOf(from)).Kind()
	toKind := indirectType(reflect.TypeOf(to)).Kind()
	// 1. slice to slice
	if toKind == reflect.Slice && fromKind == reflect.Slice {
		copySlice(fromValue, toValue, opt)
		// 2. map to map
	} else if toKind == reflect.Map && fromKind == reflect.Map {
		copyMap(fromValue, toValue, opt)
		// 3. struct to struct
	} else if toKind == reflect.Struct && fromKind == reflect.Struct {
		copyStruct(fromValue, toValue, opt)
	} else {
		panic("can only copy slice to slice, map to map, struct to struct.")
	}
}
