/*
 * File: copySlice.go
 * Created Date: 2022-01-26 06:14:15
 * Author: ysj
 * Description:  slice to slice copy
 */

package gocopy

import (
	"reflect"
)

func copySlice(fromValue, toValue reflect.Value, opt *Option) {
	fromType := indirectType(fromValue.Type())
	toType := indirectType(toValue.Type())

	fromElemType := fromType.Elem()
	toElemType := toType.Elem()

	// if append is false
	if !opt.Append {
		toValue.Set(toValue.Slice(0, 0))
	}

	// 1. if can directly assign
	if fromElemType.AssignableTo(toElemType) {
		for i := 0; i < fromValue.Len(); i++ {
			if !fromValue.Index(i).IsValid() {
				continue
			}
			toValue.Set(reflect.Append(toValue, fromValue.Index(i)))
		}
		// 2. if can not directly assign but can convert
	} else if fromElemType.ConvertibleTo(toElemType) {
		for i := 0; i < fromValue.Len(); i++ {
			if !fromValue.Index(i).IsValid() {
				continue
			}
			convertValue := fromValue.Index(i).Convert(toElemType)
			toValue.Set(reflect.Append(toValue, convertValue))
		}
	} else {
		// 3. ElementKind - if can not directly assign or convert
		fromElemKind := fromElemType.Kind()
		toElemKind := toElemType.Kind()
		// a. slice to slice
		if toElemKind == reflect.Slice && fromElemKind == reflect.Slice {
			for i := 0; i < fromValue.Len(); i++ {
				if !fromValue.Index(i).IsValid() {
					continue
				}
				toElemValue := reflect.New(toElemType)
				copySlice(fromValue.Index(i), toElemValue, opt)
				toValue = reflect.Append(toValue, toElemValue)
			}
			// b. map to map
		} else if toElemKind == reflect.Map && fromElemKind == reflect.Map {
			for i := 0; i < fromValue.Len(); i++ {
				if !fromValue.Index(i).IsValid() {
					continue
				}
				toElemValue := reflect.New(toElemType)
				copyMap(fromValue.Index(i), toElemValue, opt)
				toValue = reflect.Append(toValue, toElemValue)
			}
			// c. struct to struct
		} else if toElemKind == reflect.Struct && fromElemKind == reflect.Struct {
			for i := 0; i < fromValue.Len(); i++ {
				if !fromValue.Index(i).IsValid() {
					continue
				}
				toElemValue := reflect.New(toElemType)
				copyStruct(fromValue.Index(i), toElemValue, opt)
				toValue = reflect.Append(toValue, toElemValue)
			}
		}
	}

}
