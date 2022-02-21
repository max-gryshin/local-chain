package utils

import (
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
