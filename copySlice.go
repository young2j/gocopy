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

func copySlice(toValue, fromValue reflect.Value, opt *Option) {
	// dereference ptr value
	fromValue = indirectValue(fromValue)
	toValue = indirectValue(toValue)

	// type
	fromType, _ := indirectType(fromValue.Type())
	toType, _ := indirectType(toValue.Type())

	// element type, may be ptr
	fromElemType := fromType.Elem()
	toElemType := toType.Elem()

	// if append is false
	if !opt.Append {
		toValue.Set(toValue.Slice(0, 0))
	}
	// if entire slice can assignable
	if fromType.AssignableTo(toType) {
		toValue.Set(reflect.AppendSlice(toValue, fromValue))
		// if entire slice can convertible
	} else if fromType.ConvertibleTo(toType) {
		toValue.Set(reflect.AppendSlice(toValue, fromValue.Convert(toType)))
		// handle respectively every element
	} else {
		// 1. if element can directly assign
		if fromElemType.AssignableTo(toElemType) {
			for i := 0; i < fromValue.Len(); i++ {
				if !fromValue.Index(i).IsValid() {
					continue
				}
				toValue.Set(reflect.Append(toValue, fromValue.Index(i)))
			}
			// 2. if element can not directly assign but can convert
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
			fromElemType, _ := indirectType(fromElemType)
			toElemType, elemIsPtr := indirectType(toElemType)
			fromElemKind := fromElemType.Kind()
			toElemKind := toElemType.Kind()
			// a. slice to slice
			if toElemKind == reflect.Slice && fromElemKind == reflect.Slice {
				for i := 0; i < fromValue.Len(); i++ {
					if !fromValue.Index(i).IsValid() {
						continue
					}
					toElemValue := reflect.New(toElemType)
					copySlice(toElemValue, fromValue.Index(i), opt)
					if !elemIsPtr {
						toElemValue = indirectValue(toElemValue)
					}
					toValue.Set(reflect.Append(toValue, toElemValue))
				}
				// b. map to map
			} else if toElemKind == reflect.Map && fromElemKind == reflect.Map {
				for i := 0; i < fromValue.Len(); i++ {
					if !fromValue.Index(i).IsValid() {
						continue
					}
					toElemValue := reflect.New(toElemType)
					copyMap(toElemValue, fromValue.Index(i), opt)
					if !elemIsPtr {
						toElemValue = indirectValue(toElemValue)
					}
					toValue.Set(reflect.Append(toValue, toElemValue))
				}
				// c. struct to struct
			} else if toElemKind == reflect.Struct && fromElemKind == reflect.Struct {
				for i := 0; i < fromValue.Len(); i++ {
					if !fromValue.Index(i).IsValid() {
						continue
					}
					toElemValue := reflect.New(toElemType)
					copyStruct(toElemValue, fromValue.Index(i), opt)
					if !elemIsPtr {
						toElemValue = indirectValue(toElemValue)
					}
					toValue.Set(reflect.Append(toValue, toElemValue))
				}
			}
		}
	}
}
