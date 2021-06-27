package util

import (
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

func IsEqual(a, b interface{}) bool {
	if IsNull(a) && IsNull(b) {
		return true
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
