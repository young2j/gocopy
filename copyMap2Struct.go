package gocopy

import (
	"reflect"
	"time"
)

// copyMap2Struct copies string-keyed map values into matching exported struct fields.
func copyMap2Struct(toValue, fromValue reflect.Value, opt *Option) {
	if beginPointerCopy(fromValue, toValue, opt) {
		return
	}
	fromValue = indirectValue(fromValue)
	toValue = indirectValue(toValue)
	if fromValue.Kind() != reflect.Map || toValue.Kind() != reflect.Struct {
		return
	}

	infos := fieldInfos(toValue.Type())
	var keyIndex map[string]reflect.Value
	if fromValue.Type().Key().Kind() == reflect.Interface && fromValue.Len() > 4 {
		keyIndex = make(map[string]reflect.Value, fromValue.Len())
		iter := fromValue.MapRange()
		for iter.Next() {
			key := iter.Key()
			if key.Kind() == reflect.Interface {
				if key.IsNil() {
					continue
				}
				key = key.Elem()
			}
			if key.Kind() == reflect.String {
				keyIndex[key.String()] = iter.Value()
			}
		}
	}
	ignored := ignoredFields(opt)
	for _, info := range infos {
		field := info.field
		if info.skip {
			continue
		}
		sourceName, sourceValue, ok := mapField(fromValue, field, info.tagName, info.lowerCamel, keyIndex, opt)
		if !ok {
			continue
		}
		if level, ignoredField := ignored[sourceName]; ignoredField {
			ignored[sourceName]++
			if level <= opt.IgnoreLevel {
				continue
			}
		}

		dest := toValue.FieldByIndex(field.Index)
		if !dest.CanSet() {
			continue
		}
		if valueIsZero(sourceValue) && opt.IgnoreZero {
			continue
		}

		if convertFunc, ok := opt.Converters[sourceName]; ok {
			setMapResult(dest, convertFunc(mapInterfaceValue(sourceValue)), opt)
			continue
		}
		if timeMap, ok := opt.StringToTime[sourceName]; ok {
			if value, ok := mapInterfaceValue(sourceValue).(string); ok && value != "" {
				setMapResult(dest, parseMapTime(value, timeMap), opt)
			}
			continue
		}
		if timeMap, ok := opt.TimeToString[sourceName]; ok {
			if value, ok := mapInterfaceValue(sourceValue).(time.Time); ok {
				setMapResult(dest, formatMapTime(value, timeMap), opt)
			}
			continue
		}
		assignMapValue(dest, sourceValue, opt)
	}
}

func mapField(fromValue reflect.Value, field reflect.StructField, tagName, lowerCamel string, keyIndex map[string]reflect.Value, opt *Option) (string, reflect.Value, bool) {
	fieldName := field.Name
	// Most copies have no rename rules. Keep the candidate names on the stack
	// and allocate only when a large rename table requires more space.
	var nameBuf [6]string
	names := nameBuf[:0]
	if len(opt.NameFromTo)*2+3 > len(nameBuf) {
		names = make([]string, 0, len(opt.NameFromTo)*2+3)
	}
	if sources := reverseNames(opt)[fieldName]; len(sources) > 0 {
		names = append(names, sources...)
	}
	if tagName != "" {
		if sources := reverseNames(opt)[tagName]; len(sources) > 0 {
			names = append(names, sources...)
		}
		names = append(names, tagName)
	}
	names = append(names, fieldName)
	caseName := lowerCamel
	if opt.ToCase != "" {
		if convertCase, ok := caseFunc[opt.ToCase]; ok {
			caseName = convertCase(fieldName)
		}
	}
	if caseName != fieldName {
		names = append(names, caseName)
	}

	for _, name := range names {
		if keyIndex != nil {
			if value, ok := keyIndex[name]; ok {
				return name, value, true
			}
		} else if value, ok := mapLookup(fromValue, name); ok {
			return name, value, true
		}
	}
	return "", reflect.Value{}, false
}

func mapLookup(fromValue reflect.Value, name string) (reflect.Value, bool) {
	keyType := fromValue.Type().Key()
	if keyType.Kind() == reflect.String {
		key := reflect.ValueOf(name)
		if key.Type() != keyType {
			key = key.Convert(keyType)
		}
		value := fromValue.MapIndex(key)
		return value, value.IsValid()
	}
	if keyType.Kind() != reflect.Interface {
		return reflect.Value{}, false
	}
	iter := fromValue.MapRange()
	for iter.Next() {
		key := iter.Key()
		if key.Kind() == reflect.Interface {
			if key.IsNil() {
				continue
			}
			key = key.Elem()
		}
		if key.Kind() == reflect.String && key.String() == name {
			return iter.Value(), true
		}
	}
	return reflect.Value{}, false
}

func assignMapValue(dest, source reflect.Value, opt *Option) bool {
	if !dest.CanSet() || !source.IsValid() {
		return false
	}
	for source.Kind() == reflect.Interface {
		if source.IsNil() {
			dest.Set(reflect.Zero(dest.Type()))
			return true
		}
		source = source.Elem()
	}
	if source.Type().AssignableTo(dest.Type()) {
		dest.Set(source)
		return true
	}
	if dest.Kind() == reflect.Ptr {
		if source.Kind() == reflect.Ptr {
			if source.IsNil() {
				dest.Set(reflect.Zero(dest.Type()))
				return true
			}
			source = source.Elem()
		}
		if dest.IsNil() {
			dest.Set(reflect.New(dest.Type().Elem()))
		}
		if source.Kind() == reflect.Map && dest.Elem().Kind() == reflect.Struct {
			copyMap2Struct(dest, source, opt)
			return true
		}
		return assignMapValue(dest.Elem(), source, opt)
	}
	if source.Kind() == reflect.Ptr {
		if source.IsNil() {
			dest.Set(reflect.Zero(dest.Type()))
			return true
		}
		return assignMapValue(dest, source.Elem(), opt)
	}
	if source.Type().ConvertibleTo(dest.Type()) {
		dest.Set(source.Convert(dest.Type()))
		return true
	}

	switch {
	case dest.Kind() == reflect.Struct && source.Kind() == reflect.Map:
		copyMap2Struct(dest, source, opt)
		return true
	case dest.Kind() == reflect.Struct && source.Kind() == reflect.Struct:
		copyStruct(dest, source, opt)
		return true
	case dest.Kind() == reflect.Map && source.Kind() == reflect.Map:
		return assignMapToMap(dest, source, opt)
	case dest.Kind() == reflect.Slice && source.Kind() == reflect.Slice:
		return assignMapToSlice(dest, source, opt)
	}
	return false
}

func setMapResult(dest reflect.Value, value interface{}, opt *Option) {
	if value == nil {
		dest.Set(reflect.Zero(dest.Type()))
		return
	}
	assignMapValue(dest, reflect.ValueOf(value), opt)
}

func assignMapToSlice(dest, source reflect.Value, opt *Option) bool {
	if source.Type().AssignableTo(dest.Type()) {
		dest.Set(source)
		return true
	}
	if opt == nil || !opt.Append {
		dest.Set(reflect.MakeSlice(dest.Type(), 0, source.Len()))
	}
	for i := 0; i < source.Len(); i++ {
		element := reflect.New(dest.Type().Elem()).Elem()
		if !assignMapValue(element, source.Index(i), opt) {
			return false
		}
		dest.Set(reflect.Append(dest, element))
	}
	return true
}

func assignMapToMap(dest, source reflect.Value, opt *Option) bool {
	if source.Type().AssignableTo(dest.Type()) {
		dest.Set(source)
		return true
	}
	if opt == nil || !opt.Append || dest.IsNil() {
		dest.Set(reflect.MakeMapWithSize(dest.Type(), source.Len()))
	}
	iter := source.MapRange()
	for iter.Next() {
		key := iter.Key()
		if !key.Type().AssignableTo(dest.Type().Key()) {
			if !key.Type().ConvertibleTo(dest.Type().Key()) {
				return false
			}
			key = key.Convert(dest.Type().Key())
		}
		element := reflect.New(dest.Type().Elem()).Elem()
		if !assignMapValue(element, iter.Value(), opt) {
			return false
		}
		dest.SetMapIndex(key, element)
	}
	return true
}

func mapInterfaceValue(value reflect.Value) interface{} {
	for value.IsValid() && value.Kind() == reflect.Interface {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}
	if !value.IsValid() {
		return nil
	}
	return value.Interface()
}

func valueIsZero(value reflect.Value) bool {
	for value.IsValid() && value.Kind() == reflect.Interface {
		if value.IsNil() {
			return true
		}
		value = value.Elem()
	}
	return value.IsValid() && value.IsZero()
}

func parseMapTime(value string, options map[string]string) time.Time {
	loc, layout := mapTimeOptions(options)
	location, err := time.LoadLocation(loc)
	if err != nil {
		panic(err)
	}
	parsed, err := time.ParseInLocation(layout, value, location)
	if err != nil {
		panic(err)
	}
	return parsed
}

func formatMapTime(value time.Time, options map[string]string) string {
	loc, layout := mapTimeOptions(options)
	location, err := time.LoadLocation(loc)
	if err != nil {
		panic(err)
	}
	return value.In(location).Format(layout)
}

func mapTimeOptions(options map[string]string) (string, string) {
	loc, ok := options["loc"]
	if !ok {
		loc = defaultTimeLoc
	}
	layout, ok := options["layout"]
	if !ok {
		layout = defaultTimeLayout
	}
	return loc, layout
}
