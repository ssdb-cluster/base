package assert

import (
	"fmt"
	"path"
	"runtime"
	"testing"
	"reflect"
)

func Equal(t *testing.T, a, b interface{}) {
	if isEqual(a, b) {
		return
	}
	_, fn, line, _ := runtime.Caller(1)
	fmt.Printf("    [FAIL] %s:%d: assert failed: '%v' != '%v'\n", path.Base(fn), line, a, b)
	t.FailNow()
}

func NotEqual(t *testing.T, a, b interface{}) {
	if !isEqual(a, b) {
		return
	}
	_, fn, line, _ := runtime.Caller(1)
	fmt.Printf("    [FAIL] %s:%d: assert failed: '%v' == '%v'\n", path.Base(fn), line, a, b)
	t.FailNow()
}

func True(t *testing.T, a interface{}) {
	if isEqual(a, true) {
		return
	}
	_, fn, line, _ := runtime.Caller(1)
	fmt.Printf("    [FAIL] %s:%d: assert failed: '%v' != '%v'\n", path.Base(fn), line, a, true)
	t.FailNow()
}

func isNull(a interface{}) bool {
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

func isEqual(a, b interface{}) bool {
	if isNull(a) && isNull(b) {
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
