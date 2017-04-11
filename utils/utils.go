package utils

import (
	"reflect"
	"runtime"
	"strconv"
	"time"
)

// GetFunctionName uses reflection to get the name of a function as a string
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// GetTypeName uses reflection to get the name of a type as a string
func GetTypeName(i interface{}) string {
	return reflect.TypeOf(i).Name()
}

func MapKeys(m map[string]interface{}) []string {
	keys := reflect.ValueOf(m).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}
	return strkeys
}
