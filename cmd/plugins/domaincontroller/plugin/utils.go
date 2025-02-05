package plugin

import (
	"fmt"
	"reflect"
)

// StructToMap converts a struct to a map[string]string. This is useful for converting structs to maps for use in
// PowerShell scripts.
func structToMap(input interface{}) map[string]string {
	output := make(map[string]string)
	v := reflect.ValueOf(input)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		output[typeOfS.Field(i).Name] = fmt.Sprintf("%v", field.Interface())
	}
	return output
}
