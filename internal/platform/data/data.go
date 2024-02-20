package data

import (
	"reflect"
)

func OverwriteStruct(structInstance, dtoInstance any, columns []string) {
	structReflectValue := reflect.ValueOf(structInstance).Elem()
	dtoReflectValue := reflect.ValueOf(dtoInstance)

	for _, column := range columns {
		ptr := reflect.Indirect(dtoReflectValue.FieldByName(column))

		if ptr.Kind() == reflect.String ||
			ptr.Kind() == reflect.Int || ptr.Kind() == reflect.Int8 || ptr.Kind() == reflect.Int16 ||
			ptr.Kind() == reflect.Int32 || ptr.Kind() == reflect.Int64 || ptr.Kind() == reflect.Float32 ||
			ptr.Kind() == reflect.Float64 || ptr.Kind() == reflect.Bool {
			if ptr.IsZero() {
				continue
			}

			if structReflectValue.FieldByName(column).Kind() != ptr.Kind() {
				continue
			}
		} else if dtoReflectValue.FieldByName(column).IsNil() {
			continue
		}

		switch ptr.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val := ptr.Int()
			structReflectValue.FieldByName(column).SetInt(val)
		case reflect.Float32, reflect.Float64:
			val := ptr.Float()
			structReflectValue.FieldByName(column).SetFloat(val)
		case reflect.Bool:
			val := ptr.Bool()
			structReflectValue.FieldByName(column).SetBool(val)
		case reflect.String:
			val := ptr.String()
			structReflectValue.FieldByName(column).SetString(val)
		default:
		}
	}
}
