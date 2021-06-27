package util

import (
	"strconv"
	"strings"
	"reflect"
)

func InArray(s interface{}, arr interface{}) bool {
	return IndexOf(s, arr) != -1
}

func IndexOf(s interface{}, arr interface{}) int {
	vals := reflect.ValueOf(arr)
	for i := 0; i < vals.Len(); i++ {
		if vals.Index(i).Interface() == s {
			return i
		}
	}
	return -1
}

func SplitInts(s string) []int {
	var ret []int
	ps := strings.Split(s, ",")
	for _, p := range ps {
		p = strings.TrimSpace(p)
		if len(p) == 0 {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			continue
		}
		ret = append(ret, n)
	}
	return ret
}

// func remove(slice []int, s int) []int {
//     return append(slice[:s], slice[s+1:]...)
// }
