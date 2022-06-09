package utils

import (
	"fmt"
	"reflect"
)

// GetTagValue returns tags values from structure fields
func GetTagValue(e interface{}, tagName string) []interface{} {
	rez := make([]interface{}, 0)
	t := reflect.TypeOf(e)
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get(tagName)
		if len(tag) > 0 {
			rez = append(rez, tag)
		}
	}

	return rez
}

func GetMapFromModel(e interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	v := reflect.ValueOf(e)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("function GetMapFromModel only accepts models; got %T", v)
	}

	typeField := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeField.Field(i)

		if tagVal := field.Tag.Get("db"); tagVal != "" {
			val := v.Field(i).Interface()
			typeVal := reflect.TypeOf(val)
			if skipField := field.Tag.Get("goqu"); skipField != "" {
				continue
			}
			if typeVal.Kind().String() == "string" && val.(string) == "" ||
				typeVal.Kind().String() == "int" && val == 0 ||
				typeVal.Kind().String() == "float" && val == 0.0 {
				continue
			}
			out[tagVal] = val
		}
	}

	return out, nil
}
