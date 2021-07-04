package assert

import (
	"fmt"
	"path"
	"runtime"
	"testing"
	"base/lang"
)

func Equal(t *testing.T, a, b interface{}) {
	if lang.IsEqual(a, b) {
		return
	}
	_, fn, line, _ := runtime.Caller(1)
	fmt.Printf("    [FAIL] %s:%d: assert failed: '%v' != '%v'\n", path.Base(fn), line, a, b)
	t.FailNow()
}

func NotEqual(t *testing.T, a, b interface{}) {
	if !lang.IsEqual(a, b) {
		return
	}
	_, fn, line, _ := runtime.Caller(1)
	fmt.Printf("    [FAIL] %s:%d: assert failed: '%v' == '%v'\n", path.Base(fn), line, a, b)
	t.FailNow()
}

func True(t *testing.T, a interface{}) {
	if lang.IsEqual(a, true) {
		return
	}
	_, fn, line, _ := runtime.Caller(1)
	fmt.Printf("    [FAIL] %s:%d: assert failed: '%v' != '%v'\n", path.Base(fn), line, a, true)
	t.FailNow()
}
