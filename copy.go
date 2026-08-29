// Package gocopy copies values between compatible slices, maps, and structs.
// Common conversions use cached field metadata and direct assignments to keep
// reflection overhead low.
package gocopy

import (
	"fmt"
	"reflect"
)

// Option controls how CopyWithOption maps fields and combines destination data.
type Option struct {
	// Append merges slices and maps into existing destination values.
	Append bool
	// NameFromTo renames source fields to destination fields (source -> target).
	NameFromTo map[string]string
	// TimeToString formats time.Time fields. Each value may contain loc and layout.
	TimeToString map[string]map[string]string
	// StringToTime parses string fields into time.Time. Each value may contain loc and layout.
	StringToTime map[string]map[string]string
	// Converters supplies custom source-field conversion functions.
	Converters map[string]func(interface{}) interface{}
	// ToCase controls struct-to-map keys: LowerCamel (default), Camel, Snake,
	// ScreamingSnake, Kebab, or ScreamingKebab.
	ToCase string
	// IgnoreZero omits zero-valued fields when converting a struct to a map.
	IgnoreZero bool
	// IgnoreFields lists source field names to skip.
	IgnoreFields []string
	// IgnoreLevel limits IgnoreFields handling for embedded-field nesting levels.
	IgnoreLevel int
	state       *copyState
}

// Copy copies from into to using default options. to must be a pointer to a
// writable destination (for example, &dst).
func Copy(to, from interface{}) {
	CopyWithOption(to, from, &Option{})
}

// CopyE copies values and returns conversion errors instead of panicking.
func CopyE(to, from interface{}) (err error) {
	return CopyWithOptionE(to, from, &Option{})
}

// CopyWithOptionE is the error-returning form of CopyWithOption. It converts
// conversion panics into errors so callers can handle malformed input.
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

// CopyWithOption copies from into to according to opt. Supported conversions
// are slice-to-slice, map-to-map, struct-to-struct, struct-to-map, and
// map-to-struct, including nested pointers, slices, and maps.
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
