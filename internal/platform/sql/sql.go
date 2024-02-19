package sql

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

func DynamicQuery(columns []string, structInstance any) (string, []any) {
	query := make([]string, 0)
	values := make([]any, 0)

	reflectValue := reflect.ValueOf(structInstance)
	reflectType := reflect.TypeOf(structInstance)

	for i := range reflectValue.NumField() {
		fieldType := reflectType.Field(i)
		fieldValue := reflectValue.Field(i)
		dbTagName := fieldType.Tag.Get("db")

		if !fieldValue.IsZero() && slices.Contains(columns, dbTagName) {
			query = append(query, fmt.Sprintf("%s = ?", dbTagName))

			switch fieldValue.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				values = append(values, fieldValue.Int())
			case reflect.Float32, reflect.Float64:
				values = append(values, fieldValue.Float())
			case reflect.Bool:
				values = append(values, fieldValue.Bool())
			case reflect.String:
				values = append(values, fieldValue.String())
			default:
			}
		}
	}

	return strings.Join(query, ", "), values
}
