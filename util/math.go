package util

import (
	"strconv"
)

func Atoi(s string) int{
	n, _ := strconv.Atoi(s)
	return n
}

func Atoi32(s string) int32{
	n, _ := strconv.ParseInt(s, 10, 32)
	return int32(n)
}

func I32toa(u int32) string{
	return strconv.FormatInt(int64(u), 10)
}

func Atoi64(s string) int64{
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

func I64toa(u int64) string{
	return strconv.FormatInt(int64(u), 10)
}

func Atou32(s string) uint32{
	n, _ := strconv.ParseUint(s, 10, 32)
	return uint32(n)
}

func MinInt(a, b int) int{
	if a < b {
		return a
	} else {
		return b
	}
}

func MaxInt32(a, b int32) int32{
	if a > b {
		return a
	} else {
		return b
	}
}

func MinInt64(a, b int64) int64{
	if a < b {
		return a
	} else {
		return b
	}
}

func MaxInt64(a, b int64) int64{
	if a > b {
		return a
	} else {
		return b
	}
}
