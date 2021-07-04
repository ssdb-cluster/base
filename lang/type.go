package lang

import (
	"fmt"
	"reflect"
)

func IsNull(a interface{}) bool {
	if a == nil {
		return true
	}
	v := reflect.ValueOf(a)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer:
		return v.IsNil()
	case reflect.Interface, reflect.Slice:
		return v.IsNil()
	}
	return false
}

func IsInt(a interface{}) bool {
	v := reflect.ValueOf(a)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		 reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		 return true
	}
	return false
}

func IsFloat(a interface{}) bool {
	v := reflect.ValueOf(a)
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		 return true
	}
	return false
}

func IsEqual(a, b interface{}) bool {
	if IsNull(a) && IsNull(b) {
		return true
	}
	if IsInt(a) && IsInt(b) {
		return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
	}
	if IsFloat(a) && IsFloat(b) {
		return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
	}

	v1 := reflect.ValueOf(a)
	v2 := reflect.ValueOf(b)
	if v1.Kind() != v2.Kind() {
		return false
	}

	switch v1.Kind() {
	case reflect.Slice:
		if v1.Len() == 0 && v2.Len() == 0 {
			return true
		}
	}
	return reflect.DeepEqual(a, b)
}
