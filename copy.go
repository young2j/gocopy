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
	// for slice,map and struct slice/map field
	Append bool
	// for struct to struct/map
	NameFromTo       map[string]string
	TimeToString     map[string]map[string]string // eg. {"CreateAt":{"loc":"Asia/Shanghai","layout":"2006-01-02"}}
	StringToTime     map[string]map[string]string // eg. {"CreateAt":{"loc":"Asia/Shanghai","layout":"2006-01-02"}}
	// for strcut to struct/map
	Converters map[string]func(interface{}) interface{}
	// only for struct to map
	ToCase       string   // eg. "LowerCamel"(default)|"Camel"|"Snake"|"ScreamingSnake"|"Kebab"|"ScreamingKebab"
	IgnoreZero   bool     // ignore zero value
	IgnoreFields []string // slice of ignore field name
	IgnoreLevel  int      // field in IgnoreFields that embed level is only less than or equal this value will be ignored.
	ignoreFields map[string]int
}

func Copy(to, from interface{}) {
	CopyWithOption(to, from, &Option{})
}

func CopyWithOption(to, from interface{}, opt *Option) {
	// init option
	if len(opt.IgnoreFields) > 0 {
		opt.ignoreFields = make(map[string]int)
		for _, f := range opt.IgnoreFields {
			opt.ignoreFields[f] = 0
		}
		opt.IgnoreFields = nil
	}

	fromValue := reflect.ValueOf(from)
	toValue := reflect.ValueOf(to)
	// avoid copy from nil
	if !indirectValue(fromValue).IsValid() {
		return
	}

	fromType, _ := indirectType(reflect.TypeOf(from))
	toType, _ := indirectType(reflect.TypeOf(to))
	fromKind := fromType.Kind()
	toKind := toType.Kind()

	// 1. slice to slice
	if toKind == reflect.Slice && fromKind == reflect.Slice {
		copySlice(toValue, fromValue, opt)
		// 2. map to map
	} else if toKind == reflect.Map && fromKind == reflect.Map {
		copyMap(toValue, fromValue, opt)
		// 3. struct to struct
	} else if toKind == reflect.Struct && fromKind == reflect.Struct {
		copyStruct(toValue, fromValue, opt)
	} else if toKind == reflect.Map && fromKind == reflect.Struct {
		copyStruct2Map(toValue, fromValue, opt)
	} else {
		panic("can only copy slice to slice, map to map, struct to struct.")
	}
}
