/*
 * File: copyMap.go
 * Created Date: 2022-01-26 06:14:44
 * Author: ysj
 * Description:  copy map to map
 */

package gocopy

import (
	"reflect"
)

func copyMap(toValue, fromValue reflect.Value, opt *Option) {
	fromType := indirectType(fromValue.Type())
	toType := indirectType(toValue.Type())

	fromElemType := fromType.Elem()
	toElemType := toType.Elem()

	// avoid tomap is nil
	if toValue.IsNil() || !opt.Append {
		toNewMap := reflect.MakeMapWithSize(toType, fromValue.Len())
		if toValue.CanSet() {
			toValue.Set(toNewMap)
		}
	}

	// 1. if can directly assign
	if fromElemType.AssignableTo(toElemType) {
		kvIter := fromValue.MapRange()
		for kvIter.Next() {
			fromK := kvIter.Key()
			fromV := kvIter.Value()
			if !opt.Append { // not append, just directly set
				toValue.SetMapIndex(fromK, fromV)
			} else { // append mode
				fromVKind := indirectType(fromV.Type()).Kind()
				switch fromVKind {
				// slice append slice, need to avoid zero slice
				case reflect.Slice:
					toV := toValue.MapIndex(fromK)
					if !toV.IsValid() { // zero slice
						toV = indirectValue(reflect.New(fromV.Type()))
					}
					dest := indirectValue(reflect.New(fromV.Type()))
					dest.Set(toV)
					copySlice(dest, fromV, opt)
					toValue.SetMapIndex(fromK, dest)
				// map set kv, need to avoid nil map
				case reflect.Map:
					toV := reflect.MakeMapWithSize(toElemType, fromV.Len())
					copyMap(toV, fromV, opt)
					toValue.SetMapIndex(fromK, toV)
				default:
					toValue.SetMapIndex(fromK, fromV)
				}
			}
		}
		// 2. if can not directly assign but can convert
	} else if fromElemType.ConvertibleTo(toElemType) {
		kvIter := fromValue.MapRange()
		for kvIter.Next() {
			fromK := kvIter.Key()
			fromV := kvIter.Value().Convert(toElemType)
			if !opt.Append { // not append, just directly set converted value
				toValue.SetMapIndex(fromK, fromV)
			} else { // append mode
				fromVKind := indirectType(fromV.Type()).Kind()
				switch fromVKind {
				// slice append slice, need to avoid zero slice
				case reflect.Slice:
					toV := toValue.MapIndex(fromK)
					if !toV.IsValid() { // zero slice
						toV = indirectValue(reflect.New(fromV.Type()))
					}
					dest := indirectValue(reflect.New(fromV.Type()))
					dest.Set(toV)
					copySlice(dest, fromV, opt)
					toValue.SetMapIndex(fromK, dest)
				// map set kv, need to avoid nil map
				case reflect.Map:
					toV := reflect.MakeMapWithSize(toElemType, fromV.Len())
					copyMap(toV, fromV, opt)
					toValue.SetMapIndex(fromK, toV)
				default:
					toValue.SetMapIndex(fromK, fromV)
				}
			}
		}
	} else {
		// 3. value kind - if can not directly assign or convert
		fromElemKind := fromElemType.Kind()
		toElemKind := toElemType.Kind()
		// a. slice to slice
		if toElemKind == reflect.Slice && fromElemKind == reflect.Slice {
			kvIter := fromValue.MapRange()
			for kvIter.Next() {
				fromK := kvIter.Key()
				fromV := kvIter.Value() //slice
				copySlice(toValue.MapIndex(fromK), fromV, opt)
			}
			// b. map to map
		} else if toElemKind == reflect.Map && fromElemKind == reflect.Map {
			kvIter := fromValue.MapRange()
			for kvIter.Next() {
				fromK := kvIter.Key()
				fromV := kvIter.Value() // map
				toV := reflect.MakeMapWithSize(toElemType, fromV.Len())
				toValue.SetMapIndex(fromK, toV)
				copyMap(toV, fromV, opt)
			}
			// c. struct to struct
		} else if toElemKind == reflect.Struct && fromElemKind == reflect.Struct {
			kvIter := fromValue.MapRange()
			for kvIter.Next() {
				fromK := kvIter.Key()
				fromV := kvIter.Value() // struct
				copyStruct(toValue.MapIndex(fromK), fromV, opt)
			}
		}
	}
}
