/*
 * File: copy.go
 * Created Date: 2022-01-26 05:36:00
 * Author: ysj
 * Description: copy slice and map by type; copy struct by field name
 */
package gocopy

import (
	"fmt"
	"reflect"
)

type Option struct {
	// for slice,map and struct slice/map field
	Append bool
	// for struct to struct/map
	NameFromTo   map[string]string
	TimeToString map[string]map[string]string // eg. {"CreateAt":{"loc":"Asia/Shanghai","layout":"2006-01-02"}}
	StringToTime map[string]map[string]string // eg. {"CreateAt":{"loc":"Asia/Shanghai","layout":"2006-01-02"}}
	// for strcut to struct/map
	Converters map[string]func(interface{}) interface{}
	// only for struct to map
	ToCase       string   // eg. "LowerCamel"(default)|"Camel"|"Snake"|"ScreamingSnake"|"Kebab"|"ScreamingKebab"
	IgnoreZero   bool     // ignore zero value
	IgnoreFields []string // slice of ignore field name
	IgnoreLevel  int      // field in IgnoreFields that embed level is only less than or equal this value will be ignored.
	state        *copyState
}

func Copy(to, from interface{}) {
	CopyWithOption(to, from, &Option{})
}

// CopyE copies values and returns conversion errors instead of panicking.
func CopyE(to, from interface{}) (err error) {
	return CopyWithOptionE(to, from, &Option{})
}

// CopyWithOptionE is the error-returning form of CopyWithOption.
func CopyWithOptionE(to, from interface{}, opt *Option) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			if recoveredErr, ok := recovered.(error); ok {
				err = recoveredErr
			} else {
				err = fmt.Errorf("%v", recovered)
			}
		}
	}()
	CopyWithOption(to, from, opt)
	return nil
}

func CopyWithOption(to, from interface{}, opt *Option) {
	if opt == nil {
		opt = &Option{}
	}
	if opt.state != nil {
		copyWithOptionValues(reflect.ValueOf(to), reflect.ValueOf(from), opt)
		return
	}
	// Work on a private option copy so callers can safely reuse an Option,
	// including from concurrent goroutines. Nested calls receive the runtime
	// copy and reuse its state through the fast path below.
	runtimeOpt := *opt
	runtimeOpt.state = copyStatePool.Get().(*copyState)
	if len(runtimeOpt.IgnoreFields) > 0 {
		runtimeOpt.state.ignoreFields = make(map[string]int, len(runtimeOpt.IgnoreFields))
		for _, field := range runtimeOpt.IgnoreFields {
			runtimeOpt.state.ignoreFields[field] = 0
		}
	}
	if len(runtimeOpt.NameFromTo) > 0 {
		runtimeOpt.state.reverseNames = make(map[string][]string, len(runtimeOpt.NameFromTo))
		for source, destination := range runtimeOpt.NameFromTo {
			runtimeOpt.state.reverseNames[destination] = append(runtimeOpt.state.reverseNames[destination], source)
		}
	}
	defer func() {
		runtimeOpt.state.reset()
		copyStatePool.Put(runtimeOpt.state)
	}()
	copyWithOptionValues(reflect.ValueOf(to), reflect.ValueOf(from), &runtimeOpt)
}

func copyWithOptionValues(toValue, fromValue reflect.Value, opt *Option) {
	if opt == nil {
		return
	}
	// avoid copy from nil
	if !indirectValue(fromValue).IsValid() {
		return
	}

	fromType, _ := indirectType(fromValue.Type())
	toType, _ := indirectType(toValue.Type())
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
	} else if toKind == reflect.Struct && fromKind == reflect.Map {
		copyMap2Struct(toValue, fromValue, opt)
	} else {
		panic("can only copy slice to slice, map to map, struct to struct, map to struct.")
	}
}
