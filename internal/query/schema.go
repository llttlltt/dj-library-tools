package query

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ReflectSchema extracts queryable fields and their types from a struct using 'query' tags.
func ReflectSchema(v interface{}) (map[string]FieldType, []string) {
	schema := make(map[string]FieldType)
	var allowed []string

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("query")
		if tag == "" || tag == "-" {
			continue
		}

		parts := strings.Split(tag, ",")
		name := parts[0]
		allowed = append(allowed, name)

		fType := TypeString
		if len(parts) > 1 && parts[1] == "numeric" {
			fType = TypeNumeric
		}
		schema[name] = fType
	}

	return schema, allowed
}

// GetFieldValue uses reflection or a custom interface to extract a field's value as a string.
func GetFieldValue(obj interface{}, fieldName string) string {
	// 1. Try Custom Interface (for derived fields like cue counts)
	if getter, ok := obj.(interface{ GetQueryValue(string) (string, bool) }); ok {
		if val, ok := getter.GetQueryValue(fieldName); ok {
			return val
		}
	}

	// 2. Try Reflection
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("query")
		if tag == "" || tag == "-" {
			continue
		}
		parts := strings.Split(tag, ",")
		if parts[0] == fieldName {
			fVal := val.Field(i)
			return formatReflectedValue(fVal)
		}
	}

	return ""
}

func formatReflectedValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%g", v.Float())
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}
