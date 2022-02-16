/*
 * File: copyStruct2Map.go
 * Created Date: 2022-02-16 12:54:13
 * Author: ysj
 * Description:
 */

package gocopy

import (
	"reflect"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/iancoleman/strcase"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var caseFunc = map[string]func(s string) string{
	"Camel":          strcase.ToCamel,
	"LowerCamel":     strcase.ToLowerCamel,
	"Snake":          strcase.ToSnake,
	"ScreamingSnake": strcase.ToScreamingSnake,
	"Kebab":          strcase.ToKebab,
	"ScreamingKebab": strcase.ToScreamingKebab,
}

func copyStruct2Map(toValue, fromValue reflect.Value, opt *Option) {
	fromValue = indirectValue(fromValue)
	toValue = indirectValue(toValue)
	// type
	fromType, _ := indirectType(fromValue.Type())
	toType, _ := indirectType(toValue.Type())

	fromFields := deepFields(fromType)

	// avoid tomap is nil or not append
	if toValue.IsNil() || !opt.Append {
		toNewMap := reflect.MakeMapWithSize(toType, len(fromFields))
		toValue.Set(toNewMap)
	}

	for i := 0; i < len(fromFields); i++ {
		fromField := fromFields[i]

		// field case
		if opt.ToCase == "" {
			opt.ToCase = "LowerCamel"
		}
		toFieldName := caseFunc[opt.ToCase](fromField.Name)
		// the specified toName has the highest priority
		toName, ok := opt.NameFromTo[fromField.Name]
		if ok {
			toFieldName = toName
		}

		// target field name
		toKey := reflect.ValueOf(toFieldName)

		// whether ignore zero value
		fromFV := fromValue.FieldByName(fromField.Name)
		if fromFV.IsZero() && opt.IgnoreZero {
			continue
		}
		// ignore invalid source field
		fromFieldValue := indirectValue(fromFV)
		if !fromFieldValue.IsValid() {
			continue
		}

		fromFieldType, fromVIsPtr := indirectType(fromFV.Type())
		toV := indirectValue(toValue.MapIndex(toKey))
		toVType := fromFieldType
		toVIsPtr := fromVIsPtr
		if toV.IsValid() {
			toVType = reflect.TypeOf(toV.Interface())
			toVType, toVIsPtr = indirectType(toVType)
		}

		// specially handle bson.ObjectId to string and vice versa
		if objectIdType, ok := opt.ObjectIdToString[fromField.Name]; ok {
			if objectIdType == "mgo" {
				objectId, ok := fromFieldValue.Interface().(bson.ObjectId)
				if ok {
					if toVIsPtr {
						toFV := indirectValue(reflect.New(toVType))
						toFV.Set(reflect.ValueOf(objectId.Hex()))
						toValue.SetMapIndex(toKey, toFV.Addr())
					} else {
						toValue.SetMapIndex(toKey, reflect.ValueOf(objectId.Hex()))
					}
					continue
				}
			} else if objectIdType == "official" {
				objectId, ok := fromFieldValue.Interface().(primitive.ObjectID)
				if ok {
					if toVIsPtr {
						toFV := indirectValue(reflect.New(toVType))
						toFV.Set(reflect.ValueOf(objectId.Hex()))
						toValue.SetMapIndex(toKey, toFV.Addr())
					} else {
						toValue.SetMapIndex(toKey, reflect.ValueOf(objectId.Hex()))
					}
					continue
				}
			}
		}
		if objectIdType, ok := opt.StringToObjectId[fromField.Name]; ok {
			if objectIdType == "mgo" {
				objectId := bson.ObjectIdHex(fromFieldValue.String())
				if toVIsPtr {
					toFV := indirectValue(reflect.New(toVType))
					toFV.Set(reflect.ValueOf(objectId))
					toValue.SetMapIndex(toKey, toFV.Addr())
				} else {
					toValue.SetMapIndex(toKey, reflect.ValueOf(objectId))
				}
				continue
			} else if objectIdType == "official" {
				if objectId, err := primitive.ObjectIDFromHex(fromFieldValue.String()); err == nil {
					if toVIsPtr {
						toFV := indirectValue(reflect.New(toVType))
						toFV.Set(reflect.ValueOf(objectId))
						toValue.SetMapIndex(toKey, toFV.Addr())
					} else {
						toValue.SetMapIndex(toKey, reflect.ValueOf(objectId))
					}
					continue
				}
			}
		}

		// specially handle time.Time to string and vice versa
		if timeFieldMap, ok := opt.TimeToString[fromField.Name]; ok {
			timeString := ""
			if timeFieldMap == nil {
				location, err := time.LoadLocation(defaultTimeLoc)
				if err != nil {
					panic(err)
				}
				timeString = fromFieldValue.Interface().(time.Time).In(location).Format(defaultTimeLayout)
			} else {
				loc, ok := timeFieldMap["loc"]
				if !ok {
					loc = defaultTimeLoc
				}
				layout, ok := timeFieldMap["layout"]
				if !ok {
					layout = defaultTimeLayout
				}
				location, err := time.LoadLocation(loc)
				if err != nil {
					panic(err)
				}
				timeString = fromFieldValue.Interface().(time.Time).In(location).Format(layout)
			}
			if toVIsPtr {
				toFV := indirectValue(reflect.New(toVType))
				toFV.Set(reflect.ValueOf(timeString))
				toValue.SetMapIndex(toKey, toFV.Addr())
			} else {
				toValue.SetMapIndex(toKey, reflect.ValueOf(timeString))
			}
			continue
		}
		if stringFieldMap, ok := opt.StringToTime[fromField.Name]; ok {
			timeTime := time.Now()
			if stringFieldMap == nil {
				location, err := time.LoadLocation(defaultTimeLoc)
				if err != nil {
					panic(err)
				}
				timeTime, err = time.ParseInLocation(defaultTimeLayout, fromFieldValue.Interface().(string), location)
				if err != nil {
					panic(err)
				}
			} else {
				loc, ok := stringFieldMap["loc"]
				if !ok {
					loc = defaultTimeLoc
				}
				layout, ok := stringFieldMap["layout"]
				if !ok {
					layout = defaultTimeLayout
				}
				location, err := time.LoadLocation(loc)
				if err != nil {
					panic(err)
				}
				timeTime, err = time.ParseInLocation(layout, fromFieldValue.Interface().(string), location)
				if err != nil {
					panic(err)
				}
			}
			if toVIsPtr {
				toFV := indirectValue(reflect.New(toVType))
				toFV.Set(reflect.ValueOf(timeTime))
				toValue.SetMapIndex(toKey, toFV.Addr())
			} else {
				toValue.SetMapIndex(toKey, reflect.ValueOf(timeTime))
			}
			continue
		}

		switch fromFieldValue.Kind() {
		// slice to slice, need to avoid zero slice
		case reflect.Slice:
			if !toV.IsValid() { // zero slice
				toV = indirectValue(reflect.New(toVType))
			}
			dest := indirectValue(reflect.New(toVType))
			dest.Set(indirectValue(reflect.ValueOf(toV.Interface())))
			copySlice(dest, fromFieldValue, opt)
			if toVIsPtr {
				dest = dest.Addr()
			}
			toValue.SetMapIndex(toKey, dest)

		// map set kv, need to avoid nil map
		case reflect.Map:
			if !toV.IsValid() { // zero map
				toV = indirectValue(reflect.New(toVType))
			}
			dest := indirectValue(reflect.New(toVType))
			dest.Set(indirectValue(reflect.ValueOf(toV.Interface())))
			copyMap(dest, fromFieldValue, opt)
			if toVIsPtr {
				dest = dest.Addr()
			}
			toValue.SetMapIndex(toKey, dest)
		default:
			if toVIsPtr {
				toV = indirectValue(reflect.New(toVType))
				toV.Set(fromFieldValue)
				toV = toV.Addr()
				toValue.SetMapIndex(toKey, toV)
			} else {
				toValue.SetMapIndex(toKey, fromFieldValue)
			}
		}
	}
}
