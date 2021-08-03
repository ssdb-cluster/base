// Copyright 2021 The SSDB-cluster Authors
package util

import (
	"time"
)

var startTime time.Time = time.Now()

func Sleep(second float64) {
	time.Sleep((time.Duration)(second * 1000000) * time.Microsecond)
}

// monotonic time in seconds
func Time() float64 {
	return time.Since(startTime).Seconds()
	// return float64(time.Now().UnixNano()) / 1000 / 1000 / 1000;
}

func Millitime() int64 {
	return time.Since(startTime).Milliseconds()
	// return time.Now().UnixNano() / 1000 / 1000
}

func Microtime() int64 {
	return time.Since(startTime).Microseconds()
	// return time.Now().UnixNano() / 1000
}
