/*
 * File: copy.go
 * Created Date: 2022-01-26 05:36:00
 * Author: ysj
 * Description: copy slice and map by type; copy struct by field name
 */
package gocopy

import (
	"reflect"
)

type Option struct {
	NameFromTo       map[string]string
	ObjectIdToString map[string]string // eg. {"Id": "mgo"}
	StringToObjectId map[string]string // eg. {"Id": "official"}
	Append           bool
}

func Copy(to, from interface{}) {
	CopyWithOption(to, from, &Option{})
}

func CopyWithOption(to, from interface{}, opt *Option) {
	fromValue := indirectValue(reflect.ValueOf(from))
	toValue := indirectValue(reflect.ValueOf(to))

	fromKind := indirectType(reflect.TypeOf(from)).Kind()
	toKind := indirectType(reflect.TypeOf(to)).Kind()
	// 1. slice to slice
	if toKind == reflect.Slice && fromKind == reflect.Slice {
		copySlice(toValue, fromValue, opt)
		// 2. map to map
	} else if toKind == reflect.Map && fromKind == reflect.Map {
		copyMap(toValue, fromValue, opt)
		// 3. struct to struct
	} else if toKind == reflect.Struct && fromKind == reflect.Struct {
		copyStruct(toValue, fromValue, opt)
	} else {
		panic("can only copy slice to slice, map to map, struct to struct.")
	}
}
