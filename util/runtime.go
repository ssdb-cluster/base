package util

import (
	"strings"
	"path"
	"runtime"
)

func FuncName() string {
	pc := make([]uintptr,1)
	runtime.Callers(2,pc)
	f := runtime.FuncForPC(pc[0])
	ps := strings.Split(f.Name(), ".")
	return path.Base(ps[len(ps)-1])
}
