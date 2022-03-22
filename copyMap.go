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

// copyMap copy map to map
func copyMap(toValue, fromValue reflect.Value, opt *Option) {
	fromValue = indirectValue(fromValue)
	toValue = indirectValue(toValue)

	// type
	fromType, _ := indirectType(fromValue.Type())
	toType, _ := indirectType(toValue.Type())

	// element type, may be ptr
	fromElemType := fromType.Elem()
	toElemType := toType.Elem()

	// avoid tomap is nil or not append
	if toValue.IsNil() || !opt.Append {
		toNewMap := reflect.MakeMapWithSize(toType, fromValue.Len())
		// if toValue.CanSet() {
		toValue.Set(toNewMap)
		// }
	}

	// map can't directly assign or convert when append is true:
	// because sometimes slice value type of destination also need to keep own elements
	// eg. map[type][]int -> map[type][]int{1,2}

	// handle respectively every key-value pairs:
	// 1. if element can directly assign
	if fromElemType.AssignableTo(toElemType) {
		kvIter := fromValue.MapRange()
		for kvIter.Next() {
			fromK := kvIter.Key()
			fromV := kvIter.Value()
			if !opt.Append {
				// not append, just directly set
				toValue.SetMapIndex(fromK, fromV)
			} else {
				// append mode
				toVType, elemIsPtr := indirectType(toElemType)
				toVKind := toVType.Kind()
				switch toVKind {
				// slice append slice, need to avoid zero slice
				case reflect.Slice:
					toV := indirectValue(toValue.MapIndex(fromK))
					if !toV.IsValid() { // zero slice
						toV = indirectValue(reflect.New(toVType))
					}
					dest := indirectValue(reflect.New(toVType))
					dest.Set(toV)
					copySlice(dest, fromV, opt)
					if elemIsPtr {
						dest = dest.Addr()
					}
					toValue.SetMapIndex(fromK, dest)

				// map set kv, need to avoid nil map
				case reflect.Map:
					toV := indirectValue(toValue.MapIndex(fromK))
					if !toV.IsValid() { // zero map
						toV = indirectValue(reflect.New(toVType))
					}
					dest := indirectValue(reflect.New(toVType))
					dest.Set(toV)
					copyMap(dest, fromV, opt)
					if elemIsPtr {
						dest = dest.Addr()
					}
					toValue.SetMapIndex(fromK, dest)
				case reflect.Interface: // like bson.M
					toV := indirectValue(toValue.MapIndex(fromK))
					if !toV.IsValid() { // zero value
						toV = indirectValue(reflect.New(toVType))
					}
					if toV.IsNil() {
						toValue.SetMapIndex(fromK, fromV)
					} else {
						CopyWithOption(toV.Interface(), fromV.Interface(), opt)
					}
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
			if !opt.Append {
				// not append, just directly set converted value
				toValue.SetMapIndex(fromK, fromV)
			} else {
				// append mode
				toVType, elemIsPtr := indirectType(toElemType)
				toVKind := toVType.Kind()
				switch toVKind {
				// slice append slice, need to avoid zero slice
				case reflect.Slice:
					toV := toValue.MapIndex(fromK)
					if !toV.IsValid() { // zero slice
						toV = indirectValue(reflect.New(fromV.Type()))
					}
					dest := indirectValue(reflect.New(fromV.Type()))
					dest.Set(toV)
					copySlice(dest, fromV, opt)
					if elemIsPtr {
						dest = dest.Addr()
					}
					toValue.SetMapIndex(fromK, dest)
				// map set kv, need to avoid nil map
				case reflect.Map:
					toV := indirectValue(toValue.MapIndex(fromK))
					if !toV.IsValid() { // zero slice
						toV = indirectValue(reflect.New(toVType))
					}
					dest := indirectValue(reflect.New(toVType))
					dest.Set(toV)
					copyMap(dest, fromV, opt)
					if elemIsPtr {
						dest = dest.Addr()
					}
					toValue.SetMapIndex(fromK, dest)
				case reflect.Interface: // like bson.M
					toV := indirectValue(toValue.MapIndex(fromK))
					if !toV.IsValid() { // zero value
						toV = indirectValue(reflect.New(toVType))
					}
					if toV.IsNil() {
						toValue.SetMapIndex(fromK, fromV)
					} else {
						CopyWithOption(toV.Interface(), fromV.Interface(), opt)
					}
				default:
					toValue.SetMapIndex(fromK, fromV)
				}
			}
		}
	} else {
		// 3. value kind - if can not directly assign or convert
		fromElemType, _ := indirectType(fromElemType)
		toElemType, elemIsPtr := indirectType(toElemType)
		fromElemKind := fromElemType.Kind()
		toElemKind := toElemType.Kind()

		// a. slice to slice
		if toElemKind == reflect.Slice && fromElemKind == reflect.Slice {
			kvIter := fromValue.MapRange()
			for kvIter.Next() {
				fromK := kvIter.Key()
				fromV := kvIter.Value() //slice
				toV := toValue.MapIndex(fromK)
				if !toV.IsValid() {
					toV = reflect.New(toElemType)
				}
				if !elemIsPtr { //toV can't set
					dest := indirectValue(reflect.New(toElemType))
					dest.Set(indirectValue(toV))
					copySlice(dest, fromV, opt)
					toValue.SetMapIndex(fromK, dest)
				} else { // toV is ptr
					copySlice(toV, fromV, opt)
					toValue.SetMapIndex(fromK, toV)
				}
			}
			// b. map to map
		} else if toElemKind == reflect.Map && fromElemKind == reflect.Map {
			kvIter := fromValue.MapRange()
			for kvIter.Next() {
				fromK := kvIter.Key()
				fromV := kvIter.Value() // map
				toV := toValue.MapIndex(fromK)
				if !toV.IsValid() {
					toV = reflect.New(toElemType)
				}
				copyMap(toV, fromV, opt)
				if !elemIsPtr {
					toV = indirectValue(toV)
				}
				toValue.SetMapIndex(fromK, toV)
			}
			// c. struct to struct
		} else if toElemKind == reflect.Struct && fromElemKind == reflect.Struct {
			kvIter := fromValue.MapRange()
			for kvIter.Next() {
				fromK := kvIter.Key()
				fromV := kvIter.Value() // struct
				toV := toValue.MapIndex(fromK)
				if !toV.IsValid() {
					toV = reflect.New(toElemType)
				}
				copyStruct(toV, fromV, opt)
				if !elemIsPtr {
					toV = indirectValue(toV)
				}
				toValue.SetMapIndex(fromK, toV)
			}
		}
	}
}
