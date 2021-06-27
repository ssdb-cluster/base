package util

import (
	"testing"
	"fmt"
	"runtime"
	"path"
)

func Assert(t *testing.T, a, b interface{}) {
	eq := false
	// golang: nil interface is not nil!
	eq = (fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b))
	if !eq {
		_, fn, line, _ := runtime.Caller(1)
		fmt.Printf("    [FAIL] %s:%d: assert failed: '%v' != '%v'\n", path.Base(fn), line, a, b)
		// t.Fatalf("")
		t.FailNow()
	}
}

func AssertTrue(t *testing.T, test bool) {
	if test {
		return
	}
	_, fn, line, _ := runtime.Caller(1)
	fmt.Printf("    [FAIL] %s:%d: assert failed\n", path.Base(fn), line)
	// t.Fatalf("")
	t.FailNow()
}
