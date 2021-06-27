package util

import (
	"time"
)

func Sleep(second float64) {
	time.Sleep((time.Duration)(second * 1000000) * time.Microsecond)
}

func Time() float64 {
	return float64(time.Now().UnixNano()) / 1000 / 1000 / 1000;
}

func Millitime() int64 {
	return time.Now().UnixNano() / 1000 / 1000
}

func Microtime() int64 {
	return time.Now().UnixNano() / 1000
}
